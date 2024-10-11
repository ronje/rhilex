// Copyright (C) 2024 wwhai
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

package target

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rhilex/component/lostcache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/target/semtechudp"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/**
 *
 *
 */
type SemtechUdpForwarderConfig struct {
	GwMac            string `json:"mac"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	CacheOfflineData *bool  `json:"cacheOfflineData" title:"离线缓存"`
}
type SemtechUdpForwarder struct {
	typex.XStatus
	mainConfig SemtechUdpForwarderConfig
	status     typex.SourceState
	addr       *net.UDPAddr
	mac        [8]byte
}

func NewSemtechUdpForwarder(e typex.Rhilex) typex.XTarget {
	ht := new(SemtechUdpForwarder)
	ht.RuleEngine = e
	ht.mainConfig = SemtechUdpForwarderConfig{
		Host:             "127.0.0.1",
		Port:             1700,
		GwMac:            "0102030405060708",
		CacheOfflineData: new(bool),
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *SemtechUdpForwarder) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId
	lostcache.CreateLostDataTable(outEndId)
	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}
	GwMacByte, err1 := hex.DecodeString(ht.mainConfig.GwMac)
	if err1 != nil {
		return err1
	}
	if len(GwMacByte) != 8 {
		return fmt.Errorf("invalid mac addr:%s", ht.mainConfig.GwMac)
	}
	copy(ht.mac[:], GwMacByte)
	Ip := net.ParseIP(ht.mainConfig.Host)
	if Ip == nil {
		return fmt.Errorf("invalid host format:%v", ht.mainConfig.Host)
	}
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ht.mainConfig.Host, ht.mainConfig.Port))
	if err != nil {
		return err
	}
	ht.addr = addr
	return nil

}
func (ht *SemtechUdpForwarder) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	//
	ht.status = typex.SOURCE_UP
	// 补发数据
	if *ht.mainConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(ht.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				ht.To(data.Data)
				{
					lostcache.DeleteLostCacheData(ht.PointId, data.ID)
				}
			}
		}
	}

	glogger.GLogger.Info("Semtech Udp Forwarder started")
	return nil
}

func (ht *SemtechUdpForwarder) Status() typex.SourceState {
	return ht.status
}

/*
*
* 数据转发
*
 */
func (ht *SemtechUdpForwarder) To(data interface{}) (interface{}, error) {
	switch T := data.(type) {
	case string:
		SemtechPushMessage := NewSemtechPushMessage(ht.mac, []byte(T))
		SemtechPushMessageByte, err1 := SemtechPushMessage.MarshalBinary()
		if err1 != nil {
			glogger.GLogger.Error(err1)
			if *ht.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ht.PointId, lostcache.CacheDataDto{
					TargetId: ht.PointId,
					Data:     T,
				})
			}
			return nil, err1
		}
		glogger.GLogger.Debug("Semtech Udp Forwarder:", string(SemtechPushMessageByte))
		PushAck, errAck := ht.SendUdpData(SemtechPushMessageByte)
		if errAck != nil {
			glogger.GLogger.Error(errAck)
			if *ht.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ht.PointId, lostcache.CacheDataDto{
					TargetId: ht.PointId,
					Data:     T,
				})
			}
			return nil, errAck
		}
		if errCheckAck := checkAck(PushAck); errCheckAck != nil {
			glogger.GLogger.Error(errCheckAck)
			if *ht.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(ht.PointId, lostcache.CacheDataDto{
					TargetId: ht.PointId,
					Data:     T,
				})
			}
			return nil, errCheckAck
		}
	default:
		return 0, fmt.Errorf("invalid data type: %v", data)
	}
	return 0, nil
}

func (ht *SemtechUdpForwarder) Stop() {
	ht.status = typex.SOURCE_DOWN
	if ht.CancelCTX != nil {
		ht.CancelCTX()
	}
}
func (ht *SemtechUdpForwarder) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}

/*
*
* 向 SemTech UDP Forwarder发送UDP包
*
 */
func (ht *SemtechUdpForwarder) SendUdpData(data []byte) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, ht.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	Ack := [4]byte{}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	N, err := conn.Read(Ack[:])
	if err != nil {
		return nil, err
	}
	return Ack[:N], nil

}

/*
*
* 从UDP转发器拉取数据
*
 */
func (ht *SemtechUdpForwarder) PullDownlinkData() {
	for {
		select {
		case <-ht.Ctx.Done():
			return
		default:
		}
		PullDataPacket := semtechudp.PullDataPacket{
			ProtocolVersion: 2,
			RandomToken:     0x0305,
			GatewayMAC:      ht.mac,
		}
		PullDataPacketBytes, err := PullDataPacket.MarshalBinary()
		if err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		PullAck, errAck := ht.SendUdpData(PullDataPacketBytes)
		if errAck != nil {
			glogger.GLogger.Error(errAck)
			continue
		}
		if errCheckAck := checkAck(PullAck); errCheckAck != nil {
			glogger.GLogger.Error(errCheckAck)
			continue
		}
		PullACKPacket := semtechudp.PullACKPacket{
			ProtocolVersion: 2,
			RandomToken:     0x0305,
		}
		RespACKPacketBytes, errPull := PullACKPacket.MarshalBinary()
		if errPull != nil {
			glogger.GLogger.Error(errPull)
			continue
		}
		RespAck, errRespAck := ht.SendUdpData(RespACKPacketBytes)
		if errRespAck != nil {
			glogger.GLogger.Error(errRespAck)
			continue
		}
		PullRespPacket := semtechudp.PullRespPacket{}
		errUnmarshal := PullRespPacket.UnmarshalBinary(RespAck)
		if errUnmarshal != nil {
			glogger.GLogger.Error(errUnmarshal)
			continue
		}
		glogger.GLogger.Debug("semtechudp Forwarder PullResp:", PullRespPacket.String())
		time.Sleep(1 * time.Second)
	}

}

// 3125 -> 2531
func checkAck(ack []byte) error {
	if len(ack) == 4 {
		Version := ack[0]
		TokenH := ack[2] // 大端 -> 小端
		TokenL := ack[1] // 大端 -> 小端
		PushID := ack[3]
		glogger.GLogger.Debug("checkAck:", Version, TokenH, TokenL, PushID)
		return nil
	}
	return fmt.Errorf("error ack:%v", ack)
}
func NewSemtechPushMessage(Mac [8]byte, Payload []byte) semtechudp.PushDataPacket {
	currentTime := time.Now().UTC()
	GatewayMAC := [8]byte{}
	return semtechudp.PushDataPacket{
		ProtocolVersion: 2,
		RandomToken:     0x0305,
		GatewayMAC:      GatewayMAC,
		Payload: semtechudp.PushDataPayload{
			RXPK: []semtechudp.RXPK{
				{
					Time: (*semtechudp.CompactTime)(&currentTime),
					Tmst: uint32(currentTime.UnixMilli()),
					Tmms: new(int64),
					Chan: 1,
					RFCh: 1,
					Freq: 868.1, //EU868
					Stat: 1,
					Modu: "LORA",
					DatR: semtechudp.DatR{LoRa: "SF7BW125"},
					CodR: "4/5",
					RSSI: -50,
					LSNR: 7.5,
					RSig: []semtechudp.RSig{},
					Size: uint16(len(Payload)),
					Data: Payload,
					Meta: map[string]string{
						"name": "rhilex",
						"gateway_id": fmt.Sprintf("%X%X%X%X%X%X%X%X",
							Mac[0], Mac[1], Mac[2], Mac[3], Mac[4], Mac[5], Mac[6], Mac[7]),
					},
				},
			},
		},
	}
}
