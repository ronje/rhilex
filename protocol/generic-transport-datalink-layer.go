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

	"github.com/hootrhino/rhilex/utils"
)

type DataLinkLayer struct {
	transport *TransportLayer
}

func NewDataLinkLayer(config TransporterConfig) *DataLinkLayer {
	return &DataLinkLayer{transport: NewTransportLayer(config)}
}

func (dll *DataLinkLayer) DataLinkerLayerHandle(data []byte) ([]byte, error) {
	crc := dll.checksumCrc8(data)
	data = append(data, crc)
	errSend := dll.transport.SendRequest(data)
	if errSend != nil {
		return nil, errSend
	}
	Bytes, errRead := dll.transport.ReadResponse()
	if errRead != nil {
		return nil, errRead
	}
	ByteLen := len(Bytes)
	if ByteLen < 5 {
		return nil, fmt.Errorf("Invalid data length")
	}
	Sum1 := Bytes[ByteLen-1]
	Sum2 := dll.checksumCrc8(Bytes[:ByteLen-1])
	if Sum1 != Sum2 {
		return nil, fmt.Errorf("Check sum error, expected:%d, checked: %d", Sum1, Sum2)
	}
	return Bytes, nil
}

func (dll *DataLinkLayer) Status() error {
	return dll.transport.Status()
}

func (dll *DataLinkLayer) Close() error {
	return dll.transport.Close()
}

func (dll *DataLinkLayer) checksumCrc8(data []byte) byte {
	return utils.ChecksumCrc8(data)
}
