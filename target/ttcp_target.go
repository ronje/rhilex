// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package target

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/**
 * TCP
 *
 */
type TcpHostConfig struct {
	AllowPing        *bool  `json:"allowPing"`
	DataMode         string `json:"dataMode"`
	Host             string `json:"host"`
	PingPacket       string `json:"pingPacket"`
	Port             int    `json:"port"`
	Timeout          int    `json:"timeout"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}

type TTcpTarget struct {
	typex.XStatus
	client     *net.TCPConn
	mainConfig TcpHostConfig
	status     typex.SourceState
}

/*
*
* TCP 透传模式
*
 */
func NewTTcpTarget(e typex.Rhilex) typex.XTarget {
	ht := new(TTcpTarget)
	ht.RuleEngine = e
	ht.mainConfig = TcpHostConfig{
		Host:       "127.0.0.1",
		Port:       6502,
		DataMode:   "RAW_STRING",
		PingPacket: "rhilex\r\n",
		Timeout:    3000,
		AllowPing:  new(bool),
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *TTcpTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}
	return nil

}
func (ht *TTcpTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	var err error
	host := fmt.Sprintf("%s:%d", ht.mainConfig.Host, ht.mainConfig.Port)
	serverAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}
	ht.client, err = net.DialTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}, serverAddr)
	if err != nil {
		return err
	}
	if *ht.mainConfig.AllowPing {
		go func(ht *TTcpTarget) {
			for {
				select {
				case <-ht.Ctx.Done():
					return
				default:
				}
				ht.client.SetWriteDeadline(
					time.Now().Add((time.Duration(ht.mainConfig.Timeout) *
						time.Millisecond)),
				)
				_, err1 := ht.client.Write([]byte(ht.mainConfig.PingPacket))
				ht.client.SetWriteDeadline(time.Time{})
				if err1 != nil {
					glogger.GLogger.Error("TTcpTarget Ping Error:", err1)
					ht.status = typex.SOURCE_DOWN
					return
				}
				time.Sleep(5 * time.Second)
			}
		}(ht)
	}
	ht.status = typex.SOURCE_UP
	// 补发数据
	if CacheData, err1 := lostcache.GetLostCacheData(ht.PointId); err1 != nil {
		glogger.GLogger.Error(err1)
	} else {
		for _, data := range CacheData {
			_, errTo := ht.To(data.Data)
			if errTo == nil {
				lostcache.DeleteLostCacheData(data.ID)
			}
		}
	}
	glogger.GLogger.Info("TTcpTarget  success connect to:", host)
	return nil
}

func (ht *TTcpTarget) Status() typex.SourceState {
	if ht.client == nil {
		return typex.SOURCE_DOWN
	}
	_, err := ht.client.Write([]byte{})
	if err != nil {
		return typex.SOURCE_DOWN
	}
	return ht.status
}

type TcpOutEndTargetOutputData struct {
	Label string `json:"label"`
	Body  string `json:"body"`
}

func (O TcpOutEndTargetOutputData) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

/*
*
* 透传模式：字符串和十六进制
*
 */
func (ht *TTcpTarget) To(data interface{}) (interface{}, error) {
	if ht.client != nil {
		switch s := data.(type) {
		case string:
			ht.client.SetReadDeadline(
				time.Now().Add((time.Duration(ht.mainConfig.Timeout) *
					time.Millisecond)),
			)
			outputData := TcpOutEndTargetOutputData{
				Label: ht.mainConfig.PingPacket,
				Body:  s,
			}
			_, err0 := ht.client.Write([]byte(outputData.String() + "\r\n"))
			ht.client.SetReadDeadline(time.Time{})
			if err0 != nil {
				return 0, err0
			}
			return len(s), nil
		default:
			return 0, fmt.Errorf("only support string format")
		}
	}
	return 0, fmt.Errorf("tcp already disconnected")

}

func (ht *TTcpTarget) Stop() {
	ht.status = typex.SOURCE_DOWN
	if ht.CancelCTX != nil {
		ht.CancelCTX()
	}
	if ht.client != nil {
		ht.client.Close()
	}
}
func (ht *TTcpTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}
