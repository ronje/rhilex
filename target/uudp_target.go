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
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/**
 * UDP
 *
 */

type UdpHostConfig struct {
	AllowPing        *bool  `json:"allowPing" validate:"required"`
	CacheOfflineData *bool  `json:"cacheOfflineData" validate:"required"`
	Host             string `json:"host" validate:"required"`
	Port             int    `json:"port" validate:"required"`
	Timeout          int    `json:"timeout" validate:"required"`
	DataMode         string `json:"dataMode" validate:"required"`
	PingPacket       string `json:"pingPacket" validate:"required"`
}

type UdpHostMainConfig struct {
	UdpHostConfig UdpHostConfig `json:"commonConfig" validate:"required"`
}

/*
*
* 数据推到UDP
*
 */
type UUdpTarget struct {
	typex.XStatus
	mainConfig UdpHostMainConfig
	status     typex.SourceState
}

func NewUUdpTarget(e typex.Rhilex) typex.XTarget {
	ut := new(UUdpTarget)
	ut.RuleEngine = e
	ut.mainConfig = UdpHostMainConfig{
		UdpHostConfig: UdpHostConfig{
			Host:             "127.0.0.1",
			Port:             6502,
			DataMode:         "RAW_STRING",
			PingPacket:       "rhilex\r\n",
			Timeout:          3000,
			AllowPing:        new(bool),
			CacheOfflineData: new(bool),
		},
	}
	ut.status = typex.SOURCE_DOWN
	return ut
}

func (ut *UUdpTarget) Init(outEndId string, configMap map[string]any) error {
	ut.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &ut.mainConfig); err != nil {
		return err
	}
	return nil

}
func (ut *UUdpTarget) Start(cctx typex.CCTX) error {
	ut.Ctx = cctx.Ctx
	ut.CancelCTX = cctx.CancelCTX
	if *ut.mainConfig.UdpHostConfig.AllowPing {
		go func(ht *UUdpTarget) {
			for {
				select {
				case <-ht.Ctx.Done():
					return
				default:
				}
				socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
					IP:   net.ParseIP(ut.mainConfig.UdpHostConfig.Host),
					Port: ut.mainConfig.UdpHostConfig.Port,
				})
				if err != nil {
					glogger.GLogger.Error(err)
					ut.status = typex.SOURCE_DOWN
					return
				}
				socket.Close()
				time.Sleep(5 * time.Second)
			}
		}(ut)
	}
	ut.status = typex.SOURCE_UP
	// 补发数据
	if *ut.mainConfig.UdpHostConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(ut.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				ut.To(data.Data)
				{
					lostcache.DeleteLostCacheData(ut.PointId, data.ID)
				}
			}
		}
	}

	glogger.GLogger.Info("UUdpTarget started")
	return nil
}

func (ut *UUdpTarget) Status() typex.SourceState {
	if err := ut.UdpStatus(fmt.Sprintf("%s:%d",
		ut.mainConfig.UdpHostConfig.Host, ut.mainConfig.UdpHostConfig.Port)); err != nil {
		return typex.SOURCE_DOWN
	}
	return ut.status

}

func (ut *UUdpTarget) To(data any) (any, error) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ut.mainConfig.UdpHostConfig.Host),
		Port: ut.mainConfig.UdpHostConfig.Port,
	})
	if err != nil {
		return 0, err
	}
	defer socket.Close()
	switch T := data.(type) {
	case string:

		socket.SetReadDeadline(
			time.Now().Add((time.Duration(ut.mainConfig.UdpHostConfig.Timeout) *
				time.Millisecond)),
		)
		_, err0 := socket.Write([]byte(T + "\r\n"))
		socket.SetReadDeadline(time.Time{})
		if err0 != nil {
			if *ut.mainConfig.UdpHostConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ut.PointId, lostcache.CacheDataDto{
					TargetId: ut.PointId,
					Data:     T,
				})
			}
			return 0, err0
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("only support string format")
	}
}

func (ut *UUdpTarget) Stop() {
	ut.status = typex.SOURCE_DOWN
	if ut.CancelCTX != nil {
		ut.CancelCTX()
	}
}
func (ut *UUdpTarget) Details() *typex.OutEnd {
	return ut.RuleEngine.GetOutEnd(ut.PointId)
}
func (ut *UUdpTarget) UdpStatus(serverAddr string) error {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return fmt.Errorf("UDP connection failed: %v", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte{})
	if err != nil {
		return fmt.Errorf("failed to send data over UDP: %v", err)
	}
	return nil
}
