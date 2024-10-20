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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type TransporterConfig struct {
	Port         GenericPort
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Logger       *logrus.Logger
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

func (h Header) String() string {
	return fmt.Sprintf("Type: %X, Length: %d", h.Type, binary.BigEndian.Uint16(h.Length[:]))
}

type AppLayerFrame struct {
	Header  Header
	Payload []byte
}

func (frame AppLayerFrame) String() string {
	return fmt.Sprintf("Header: %s, Payload: %X", frame.Header.String(), frame.Payload)
}

func (frame AppLayerFrame) Encode() ([]byte, error) {
	payloadLength := len(frame.Payload)
	if payloadLength > 65535 {
		return nil, errors.New("payload length exceeds maximum")
	}
	encodedFrame := new(bytes.Buffer)
	encodedFrame.WriteByte(frame.Header.Type[0])
	encodedFrame.WriteByte(frame.Header.Type[1])
	encodedFrame.WriteByte(frame.Header.Length[0])
	encodedFrame.WriteByte(frame.Header.Length[1])
	encodedFrame.Write(frame.Payload)
	return encodedFrame.Bytes(), nil
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
