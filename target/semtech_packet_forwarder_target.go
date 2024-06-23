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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rhilex/glogger"
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
		GwMac: "00010203AABBCCDD",
	}
	ht.mac = [8]byte{}
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
	ht.mac[0] = GwMacByte[0]
	ht.mac[1] = GwMacByte[1]
	ht.mac[2] = GwMacByte[2]
	ht.mac[3] = GwMacByte[3]
	ht.mac[4] = GwMacByte[4]
	ht.mac[5] = GwMacByte[5]
	ht.mac[6] = GwMacByte[6]
	ht.mac[7] = GwMacByte[7]
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
		SemtechPushMessage := NewSemtechPushMessage(ht.mac, T)
		SemtechPushMessageByte, err1 := SemtechPushMessage.Encode()
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

// | Bytes | Function                                          |
// | :---: | ------------------------------------------------- |
// |   0   | protocol version = 2                              |
// |  1-2  | same token as the PUSH_DATA packet to acknowledge |
// |   3   | PUSH_ACK identifier 0x01                          |

/*
*
* 赛门铁克UDP数据包
*
 */
//  | Bytes  | Function                                                   |
//  | :----: | ---------------------------------------------------------- |
//  |   0    | protocol version = 2                                       |
//  |  1-2   | random token                                               |
//  |   3    | PUSH_DATA identifier 0x00                                  |
//  |  4-11  | Gateway unique identifier (MAC address)                    |
//  | 12-end | JSON object, starting with {, ending with }, see section 4 |

type SemtechPushMessage struct {
	Version         byte            `json:"-"`                 // 02
	TokenH          byte            `json:"-"`                 // 00
	TokenL          byte            `json:"-"`                 // 00
	Identifier      byte            `json:"-"`                 // 00
	Mac             [8]byte         `json:"-"`                 // AA BB CC DD EE FF 00 11
	PushDataPayload PushDataPayload `json:"push_data_payload"` // {...}
}

func NewSemtechPushMessage(Mac [8]byte, Payload string) SemtechPushMessage {
	// Token := genToken()
	currentTime := time.Now().UTC()
	return SemtechPushMessage{
		Version:    2,
		TokenH:     0,
		TokenL:     0,
		Identifier: 0,
		Mac:        Mac,
		PushDataPayload: PushDataPayload{
			Rxpk: []rxpk{
				{
					Time: currentTime.UTC().Format("2006-01-02T15:04:05.999999Z"),
					Tmst: uint32(time.Now().UnixMilli()),
					Chan: 1,
					Rfch: 1,
					Freq: 868.1,
					Stat: 1,
					Modu: "LORA",
					Datr: "SF7BW125",
					Codr: "4/5",
					Rssi: -50,
					Lsnr: 7.5,
					Size: len(Payload),
					Data: Payload,
				},
			},
		},
	}
}

/*
*
* 编码
*
 */
func (M SemtechPushMessage) Encode() ([]byte, error) {
	Packet := []byte{}
	Packet = append(Packet, M.Version)
	Packet = append(Packet, M.TokenH)
	Packet = append(Packet, M.TokenL)
	Packet = append(Packet, M.Identifier)
	if bytes, err := json.Marshal(M.PushDataPayload); err != nil {
		Packet = append(Packet, '{')
		Packet = append(Packet, '}')
		return Packet, err
	} else {
		Packet = append(Packet, M.Mac[0])
		Packet = append(Packet, M.Mac[1])
		Packet = append(Packet, M.Mac[2])
		Packet = append(Packet, M.Mac[3])
		Packet = append(Packet, M.Mac[4])
		Packet = append(Packet, M.Mac[5])
		Packet = append(Packet, M.Mac[6])
		Packet = append(Packet, M.Mac[7])
		Packet = append(Packet, bytes...)
		return Packet, nil
	}
}

// | Bytes | Function                                          |
// | :---: | ------------------------------------------------- |
// |   0   | protocol version = 2                              |
// |  1-2  | same token as the PUSH_DATA packet to acknowledge |
// |   3   | PUSH_ACK identifier 0x01                          |
type SemtechPushMessageAck struct {
	Version    byte
	Token      uint16
	Identifier byte
}

/*
*
* 上传RF数据格式
*
 */
type rxpk struct {
	Time string  `json:"time"`
	Tmst uint32  `json:"tmst"`
	Chan int     `json:"chan"`
	Rfch int     `json:"rfch"`
	Freq float32 `json:"freq"`
	Stat int     `json:"stat"`
	Modu string  `json:"modu"`
	Datr string  `json:"datr"`
	Codr string  `json:"codr"`
	Rssi int     `json:"rssi"`
	Lsnr float64 `json:"lsnr"`
	Size int     `json:"size"`
	Data string  `json:"data"`
}

/*
*
* UDP协议
*
 */
type PushDataPayload struct {
	Rxpk []rxpk `json:"rxpk"`
}

func genToken() [2]byte {
	var b [2]byte
	rand.Read(b[:])
	return b
}
