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

package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type TransporterConfig struct {
	Port         GenericPort
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GenericPort interface {
	io.ReadWriteCloser
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type Header struct {
	Type   [2]byte
	Length [2]byte
}

type AppLayerFrame struct {
	Header  Header
	Payload []byte
}

func (frame AppLayerFrame) Encode() ([]byte, error) {
	payloadLength := len(frame.Payload)
	if payloadLength > 65535 {
		return nil, errors.New("payload length exceeds maximum")
	}
	encodedLength := 4 + payloadLength + 1
	encodedFrame := make([]byte, encodedLength)
	encodedFrame[0] = frame.Header.Type[0]
	encodedFrame[1] = frame.Header.Type[1]
	binary.BigEndian.PutUint16(encodedFrame[2:], uint16(payloadLength))
	copy(encodedFrame[4:], frame.Payload)
	return encodedFrame, nil
}

func Decode(encodedFrame []byte) (AppLayerFrame, error) {
	if len(encodedFrame) < 6 {
		return AppLayerFrame{}, errors.New("encoded frame is too short")
	}
	headerType := [2]byte{encodedFrame[0], encodedFrame[1]}
	payloadLength := int(binary.BigEndian.Uint16(encodedFrame[2:4]))
	if len(encodedFrame)-5 != payloadLength {
		return AppLayerFrame{}, errors.New("invalid payload length")
	}
	payload := make([]byte, payloadLength)
	copy(payload, encodedFrame[4:4+payloadLength])
	return AppLayerFrame{
		Header:  Header{Type: headerType},
		Payload: payload,
	}, nil
}
