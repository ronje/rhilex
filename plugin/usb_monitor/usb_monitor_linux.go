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

package usbmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/component/eventbus"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"golang.org/x/sys/unix"
	"gopkg.in/ini.v1"
)

// USB 热插拔监控器, 方便观察USB状态, 本插件只支持Linux！！！
type usbMonitor struct {
	uuid string
}

func NewUsbMonitor() typex.XPlugin {
	return &usbMonitor{uuid: "USB_EVENT_MONITOR"}
}

func (usbm *usbMonitor) Init(_ *ini.Section) error {
	return nil
}

type _info struct {
	Type   string `json:"type"`
	Device string `json:"device"`
}

func (usbm *usbMonitor) Start(_ typex.Rhilex) error {
	// 为了减小问题, 直接把Windows给限制了不支持, 实际上大部分情况下都是Arm-Linux场景
	if runtime.GOOS == "windows" {
		return errors.New("USB monitor plugin not support windows")
	}

	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, unix.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		glogger.GLogger.Error(fmt.Sprintf("Failed to create socket: %v", err))
		return err
	}

	err = unix.Bind(fd, &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: 1,
		Pid:    0,
	})
	if err != nil {
		glogger.GLogger.Error(fmt.Sprintf("Failed to bind socket: %v", err))
		// 关闭已创建的fd，避免资源泄漏
		_ = unix.Close(fd)
		return err
	}
	defer func() {
		// 在函数结束时关闭fd
		if err := unix.Close(fd); err != nil {
			glogger.GLogger.Error(fmt.Sprintf("Failed to close socket: %v", err))
		}
	}()

	go func(ctx context.Context) {
		data := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 这里default分支可以去掉，因为没有实际操作
			}
			n, _, err := unix.Recvfrom(fd, data, 0)
			if err != nil {
				glogger.GLogger.Error(fmt.Sprintf("Failed to receive data: %v", err))
				continue
			}
			if n > 16 {
				Msg := parseType(data, n)
				if len(Msg) > 0 {
					glogger.GLogger.Info(Msg)
					lineS := "system.usb.event." + usbm.uuid
					eventbus.Publish(lineS, eventbus.EventMessage{
						Topic:   lineS,
						From:    "usb-monitor",
						Type:    "HARDWARE",
						Event:   lineS,
						Ts:      uint64(time.Now().UnixMilli()),
						Payload: Msg,
					})
				}
			}
		}
	}(context.Background())
	return nil
}

func parseType(data []byte, len int) string {
	if strings.HasPrefix(string(data), "add@") {
		return parseMsg("add", data[4:], len-4)
	}
	if strings.HasPrefix(string(data), "remove@") {
		return parseMsg("remove", data[7:], len-7)
	}
	return ""
}

// 只监控串口"/dev/tty*"设备, U盘不管
func parseMsg(Type string, data []byte, offset int) string {
	msg := string(data[:offset])
	if !strings.Contains(msg, "tty") {
		return ""
	}

	nameTokens := strings.Split(msg, "/")
	info := _info{Type: Type}

	switch len(nameTokens) {
	case 1:
		info.Device = nameTokens[0]
	case 2:
		info.Device = nameTokens[1]
	case 3:
		if nameTokens[0] != nameTokens[2] {
			info.Device = nameTokens[2]
		} else {
			return ""
		}
	default:
		return ""
	}

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		glogger.GLogger.Error(fmt.Sprintf("Failed to marshal JSON: %v", err))
		return ""
	}
	return string(jsonBytes)
}

func (usbm *usbMonitor) Stop() error {
	// 这里可以添加关闭资源的逻辑，如果有必要的话
	return nil
}

func (usbm *usbMonitor) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        usbm.uuid,
		Name:        "USB Monitor",
		Version:     "v0.0.1",
		Description: "USB Hot Plugin Monitor",
	}
}

// 服务调用接口
func (cs *usbMonitor) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
