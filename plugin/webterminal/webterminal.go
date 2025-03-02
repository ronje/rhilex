// Copyright (C) 2025 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package webterminal

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"gopkg.in/ini.v1"
)

type WebTerminalConfig struct {
	ListenPort int `ini:"listen_port" json:"listen_port"`
}
type WebTerminal struct {
	terminalPty *os.File
	httpServer  *http.Server
	upgrader    websocket.Upgrader
	busy        bool
	mu          sync.Mutex // 用于保护 busy 标志
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mainConfig  WebTerminalConfig
}

func NewWebTerminal() *WebTerminal {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebTerminal{
		terminalPty: nil,
		busy:        false,
		ctx:         ctx,
		cancel:      cancel,
		mainConfig: WebTerminalConfig{
			ListenPort: 2579,
		},
	}
}

func (wt *WebTerminal) Init(config *ini.Section) error {
	glogger.GLogger.Debug("Init web terminal")
	if err := utils.InIMapToStruct(config, &wt.mainConfig); err != nil {
		return err
	}
	return nil
}

func (wt *WebTerminal) Start(rhilex typex.Rhilex) error {
	glogger.GLogger.Debug("Start web terminal")
	bashCmd := exec.Command("/bin/bash")
	terminalPty, errStart := pty.Start(bashCmd)
	if errStart != nil {
		glogger.GLogger.Error(errStart)
		return errStart
	}
	wt.terminalPty = terminalPty

	// websocket
	wt.upgrader = websocket.Upgrader{
		ReadBufferSize:    1024 * 10,
		WriteBufferSize:   1024 * 10,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			glogger.GLogger.Error("websocket error:", status, reason)
		},
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", wt.handleTerminal)
	wt.httpServer = &http.Server{
		Addr:    ":2579",
		Handler: serverMux,
	}

	go func() {
		glogger.GLogger.Debug("WebTerminal Server Started on: 2579")
		if errListenAndServe := wt.httpServer.ListenAndServe(); errListenAndServe != nil && errListenAndServe != http.ErrServerClosed {
			glogger.GLogger.Error(errListenAndServe)
		}
	}()

	return nil
}

func (wt *WebTerminal) Stop() error {
	glogger.GLogger.Debug("Stop web terminal")
	// 取消上下文
	wt.cancel()
	// 等待所有 goroutine 完成
	wt.wg.Wait()

	if wt.terminalPty != nil {
		if err := wt.terminalPty.Close(); err != nil {
			glogger.GLogger.Error(err)
		}
	}
	if wt.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := wt.httpServer.Shutdown(ctx); err != nil {
			glogger.GLogger.Error(err)
		}
	}
	return nil
}

func (wt *WebTerminal) Restart(rhilex typex.Rhilex) error {
	glogger.GLogger.Debug("Restart web terminal")
	// 先停止服务
	if err := wt.Stop(); err != nil {
		glogger.GLogger.Error("Failed to stop web terminal during restart:", err)
		return err
	}
	// 重新创建上下文
	wt.ctx, wt.cancel = context.WithCancel(context.Background())
	// 再启动服务
	return wt.Start(rhilex)
}

func (wt *WebTerminal) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "WEB-TERMINAL",
		Name:        "WebTerminal",
		Version:     "v0.0.1",
		Description: "A simple web terminal",
	}
}

func (wt *WebTerminal) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

func (wt *WebTerminal) handleTerminal(w http.ResponseWriter, r *http.Request) {
	wsCon, err := wt.upgrader.Upgrade(w, r, nil)
	if err != nil {
		glogger.GLogger.Error(err)
		return
	}
	glogger.GLogger.Debug("websocket client connected:", wsCon.RemoteAddr().String())
	wt.mu.Lock()
	if wt.busy {
		wt.mu.Unlock()
		if err := wsCon.WriteMessage(websocket.TextMessage, []byte("Web Terminal is busy now!")); err != nil {
			glogger.GLogger.Error(err)
		}
		wsCon.Close()
		return
	}
	wt.busy = true
	wt.mu.Unlock()

	defer func() {
		glogger.GLogger.Debug("websocket client disconnected:", wsCon.RemoteAddr().String())
		wt.mu.Lock()
		wt.busy = false
		wt.mu.Unlock()
		wsCon.Close()
	}()

	wt.wg.Add(2)

	// 发送 Ping 消息
	go func() {
		defer wt.wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-wt.ctx.Done():
				return
			case <-ticker.C:
				glogger.GLogger.Debug("websocket send Ping Message")
				if err := wsCon.WriteMessage(websocket.PingMessage, nil); err != nil {
					glogger.GLogger.Error(err)
					return
				}
			}
		}
	}()

	// 从Bash读到的数据重定向到websocket
	go func() {
		defer wt.wg.Done()
		buf := make([]byte, 1024*10)
		for {
			select {
			case <-wt.ctx.Done():
				return
			default:
			}
			n, err := wt.terminalPty.Read(buf)
			if err != nil {
				glogger.GLogger.Error(err)
				return
			}
			if err := wsCon.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				glogger.GLogger.Error(err)
				return
			}
		}
	}()

	// HTML发来的数据，一般是terminal输入
	for {
		select {
		case <-wt.ctx.Done():
			return
		default:
		}
		_, message, err := wsCon.ReadMessage()
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
		// 定向到bash进程
		if _, err := wt.terminalPty.Write(message); err != nil {
			glogger.GLogger.Error(err)
			return
		}
	}
}
