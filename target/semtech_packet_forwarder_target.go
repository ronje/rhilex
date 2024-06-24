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
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/target/semtechudp"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type SemtechUdpForwarderConfig struct {
	GwMac string `json:"mac"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
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
		Host:  "127.0.0.1",
		Port:  1700,
		GwMac: "0102030405060708",
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *SemtechUdpForwarder) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId

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
			return nil, err1
		}
		if err2 := ht.SendUdpData(SemtechPushMessageByte); err2 != nil {
			glogger.GLogger.Error(err2)
			return nil, err2
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
func (ht *SemtechUdpForwarder) SendUdpData(data []byte) error {
	conn, err := net.DialUDP("udp", nil, ht.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	return nil
	// Ack := [6]byte{}
	// conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	// N, err := conn.Read(Ack[:])
	// if err != nil {
	// 	return err
	// }
	// if N >= 6 {
	// 	Version := Ack[0]
	// 	TokenH := Ack[1]
	// 	TokenL := Ack[2]
	// 	PushID := Ack[3]
	// 	if data[0] == Version &&
	// 		data[1] == TokenH &&
	// 		data[2] == TokenL &&
	// 		data[3] == PushID {
	// 		return nil
	// 	}
	// }
	// return fmt.Errorf("invalid response:%v", Ack[:N])
}
func NewSemtechPushMessage(Mac [8]byte, Payload []byte) semtechudp.PushDataPacket {
	currentTime := time.Now().UTC()
	GatewayMAC := lorawan.EUI64(Mac)
	return semtechudp.PushDataPacket{
		ProtocolVersion: 2,
		RandomToken:     0x1234,
		GatewayMAC:      GatewayMAC,
		Payload: semtechudp.PushDataPayload{
			RXPK: []semtechudp.RXPK{
				{
					Time: (*semtechudp.CompactTime)(&currentTime),
					Tmst: uint32(currentTime.UnixMilli()),
					Chan: 1,
					RFCh: 1,
					Freq: 868.1,
					Stat: 1,
					Modu: "LORA",
					DatR: semtechudp.DatR{LoRa: "SF12BW500"},
					CodR: "4/5",
					RSSI: -50,
					LSNR: 7.5,
					RSig: []semtechudp.RSig{},
					Size: uint16(len(Payload)),
					Data: Payload,
					Meta: map[string]string{
						"gateway_name": "test-gateway",
					},
				},
			},
		},
	}
}
