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
	"errors"

	"github.com/brocaar/lorawan"
)

// TXACKPacket is used by the gateway to send a feedback to the server
// to inform if a downlink request has been accepted or rejected by the
// gateway.
type TXACKPacket struct {
	ProtocolVersion uint8
	RandomToken     uint16
	GatewayMAC      lorawan.EUI64
	Payload         *TXACKPayload
}

// MarshalBinary marshals the object into binary form.
func (p TXACKPacket) MarshalBinary() ([]byte, error) {
	var pb []byte
	if p.Payload != nil {
		var err error
		pb, err = json.Marshal(p.Payload)
		if err != nil {
			return nil, err
		}
	}

	out := make([]byte, 4, len(pb)+12)
	out[0] = p.ProtocolVersion
	binary.LittleEndian.PutUint16(out[1:3], p.RandomToken)
	out[3] = byte(TXACK)
	out = append(out, p.GatewayMAC[:]...)
	out = append(out, pb...)
	return out, nil
}

// UnmarshalBinary decodes the object from binary form.
func (p *TXACKPacket) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return errors.New("gateway: at least 12 bytes of data are expected")
	}
	if data[3] != byte(TXACK) {
		return errors.New("gateway: identifier mismatch (TXACK expected)")
	}
	if !protocolSupported(data[0]) {
		return ErrInvalidProtocolVersion
	}
	p.ProtocolVersion = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])
	for i := 0; i < 8; i++ {
		p.GatewayMAC[i] = data[4+i]
	}
	if len(data) > 13 { // the min payload + the length of at least "{}"
		p.Payload = &TXACKPayload{}
		return json.Unmarshal(data[12:], p.Payload)
	}
	return nil
}

// TXACKPayload contains the TXACKPacket payload.
type TXACKPayload struct {
	TXPKACK TXPKACK `json:"txpk_ack"`
}

// TXPKACK contains the status information of the associated PULL_RESP
// packet.
type TXPKACK struct {
	Error string `json:"error"`
}
