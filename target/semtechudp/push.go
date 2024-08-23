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

package semtechudp

import (
	"encoding/binary"
	"encoding/json"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"
)

// PushDataPacket type is used by the gateway mainly to forward the RF packets
// received, and associated metadata, to the server.
type PushDataPacket struct {
	ProtocolVersion uint8
	RandomToken     uint16
	GatewayMAC      lorawan.EUI64
	Payload         PushDataPayload
}

// MarshalBinary encodes the packet into binary form compatible with the
// Semtech UDP protocol.
func (p PushDataPacket) MarshalBinary() ([]byte, error) {
	pb, err := json.Marshal(&p.Payload)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 4, len(pb)+12)
	out[0] = p.ProtocolVersion
	binary.LittleEndian.PutUint16(out[1:3], p.RandomToken)
	out[3] = byte(PushData)
	out = append(out, p.GatewayMAC[0:len(p.GatewayMAC)]...)
	out = append(out, pb...)
	return out, nil
}

// UnmarshalBinary decodes the packet from Semtech UDP binary form.
func (p *PushDataPacket) UnmarshalBinary(data []byte) error {
	if len(data) < 13 {
		return errors.New("backend/semtechudp/packets: at least 13 bytes are expected")
	}
	if data[3] != byte(PushData) {
		return errors.New("backend/semtechudp/packets: identifier mismatch (PUSH_DATA expected)")
	}

	if !protocolSupported(data[0]) {
		return ErrInvalidProtocolVersion
	}

	p.ProtocolVersion = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])
	for i := 0; i < 8; i++ {
		p.GatewayMAC[i] = data[4+i]
	}

	return json.Unmarshal(data[12:], &p.Payload)
}

// PushDataPayload represents the upstream JSON data structure.
type PushDataPayload struct {
	RXPK []RXPK `json:"rxpk,omitempty"`
	Stat *Stat  `json:"stat,omitempty"`
}

// Stat contains the status of the gateway.
type Stat struct {
	Time ExpandedTime      `json:"time"` // UTC 'system' time of the gateway, ISO 8601 'expanded' format (e.g 2014-01-12 08:59:28 GMT)
	Lati float64           `json:"lati"` // GPS latitude of the gateway in degree (float, N is +)
	Long float64           `json:"long"` // GPS latitude of the gateway in degree (float, E is +)
	Alti int32             `json:"alti"` // GPS altitude of the gateway in meter RX (integer)
	RXNb uint32            `json:"rxnb"` // Number of radio packets received (unsigned integer)
	RXOK uint32            `json:"rxok"` // Number of radio packets received with a valid PHY CRC
	RXFW uint32            `json:"rxfw"` // Number of radio packets forwarded (unsigned integer)
	ACKR float64           `json:"ackr"` // Percentage of upstream datagrams that were acknowledged
	DWNb uint32            `json:"dwnb"` // Number of downlink datagrams received (unsigned integer)
	TXNb uint32            `json:"txnb"` // Number of packets emitted (unsigned integer)
	Meta map[string]string `json:"meta"` // Custom meta-data (Optional, not part of PROTOCOL.TXT)
}

// RXPK contain a RF packet and associated metadata.
type RXPK struct {
	Time  *CompactTime      `json:"time"`  // UTC time of pkt RX, us precision, ISO 8601 'compact' format (e.g. 2013-03-31T16:21:17.528002Z)
	Tmms  *int64            `json:"tmms"`  // GPS time of pkt RX, number of milliseconds since 06.Jan.1980
	Tmst  uint32            `json:"tmst"`  // Internal timestamp of "RX finished" event (32b unsigned)
	FTime *uint32           `json:"ftime"` // Fine timestamp, number of nanoseconds since last PPS [0..999999999] (Optional)
	AESK  uint8             `json:"aesk"`  // AES key index used for encrypting fine timestamps
	Chan  uint8             `json:"chan"`  // Concentrator "IF" channel used for RX (unsigned integer)
	RFCh  uint8             `json:"rfch"`  // Concentrator "RF chain" used for RX (unsigned integer)
	Stat  int8              `json:"stat"`  // CRC status: 1 = OK, -1 = fail, 0 = no CRC
	Freq  float64           `json:"freq"`  // RX central frequency in MHz (unsigned float, Hz precision)
	Brd   uint32            `json:"brd"`   // Concentrator board used for RX (unsigned integer)
	RSSI  int16             `json:"rssi"`  // RSSI in dBm (signed integer, 1 dB precision)
	Size  uint16            `json:"size"`  // RF packet payload size in bytes (unsigned integer)
	DatR  DatR              `json:"datr"`  // LoRa datarate identifier (eg. SF12BW500) || FSK datarate (unsigned, in bits per second)
	Modu  string            `json:"modu"`  // Modulation identifier "LORA" or "FSK"
	CodR  string            `json:"codr"`  // LoRa ECC coding rate identifier
	LSNR  float64           `json:"lsnr"`  // Lora SNR ratio in dB (signed float, 0.1 dB precision)
	HPW   uint8             `json:"hpw"`   // LR-FHSS hopping grid number of steps.
	Data  []byte            `json:"data"`  // Base64 encoded RF packet payload, padded
	RSig  []RSig            `json:"rsig"`  // Received signal information, per antenna (Optional)
	Meta  map[string]string `json:"meta"`  // Custom meta-data (Optional, not part of PROTOCOL.TXT)
}

// RSig contains the received signal information per antenna.
type RSig struct {
	Ant   uint8   `json:"ant"`   // Antenna number on which signal has been received
	Chan  uint8   `json:"chan"`  // Concentrator "IF" channel used for RX (unsigned integer)
	RSSIC int16   `json:"rssic"` // RSSI in dBm of the channel (signed integer, 1 dB precision)
	LSNR  float32 `json:"lsnr"`  // Lora SNR ratio in dB (signed float, 0.1 dB precision)
	ETime []byte  `json:"etime"` // Encrypted 'main' fine timestamp, ns precision [0..999999999] (Optional)
}
