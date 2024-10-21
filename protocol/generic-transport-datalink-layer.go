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
	"fmt"
)

type DataLinkLayer struct {
	transport *TransportLayer
}

func NewDataLinkLayer(config TransporterConfig) *DataLinkLayer {
	return &DataLinkLayer{transport: NewTransportLayer(config)}
}

func (dll *DataLinkLayer) Write(data []byte) error {
	crc := dll.checksumCrc8(data)
	data = append(data, crc)
	return dll.transport.Write(data)
}

func (dll *DataLinkLayer) Read() ([]byte, error) {
	Bytes, errRead := dll.transport.Read()
	if errRead != nil {
		return nil, errRead
	}
	ByteLen := len(Bytes)
	if ByteLen < 5 {
		return nil, fmt.Errorf("Invalid data length")
	}
	Sum1 := Bytes[ByteLen-1]
	Data := Bytes[:ByteLen-1]
	Sum2 := dll.checksumCrc8(Data)
	if Sum1 != Sum2 {
		return nil, fmt.Errorf("Check sum error, expected:%d, checked: %d", Sum1, Sum2)
	}
	return Data, nil
}

func (dll *DataLinkLayer) Status() error {
	return dll.transport.Status()
}

func (dll *DataLinkLayer) Close() error {
	return dll.transport.Close()
}

// CRC8多项式: x^8 + x^2 + x + 1
func (dll *DataLinkLayer) checksumCrc8(data []byte) byte {
	crc := byte(0x00)
	for _, b := range data {
		crc ^= b
		for i := 8; i > 0; i-- {
			if crc&0x80 == 0x80 {
				crc = (crc << 1) ^ 0x07
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
