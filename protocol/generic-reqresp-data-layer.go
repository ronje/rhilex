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
	"fmt"
	"hash/crc32"
)

type DataLayer struct {
	errTxCount int32 // 错误包计数器
	errRxCount int32 // 错误包计数器
	transport  *TransportLayer
}

func NewDataLayer(config TransporterConfig) *DataLayer {
	return &DataLayer{errTxCount: 0, errRxCount: 0, transport: NewTransportLayer(config)}
}

func (dll *DataLayer) Write(data []byte) error {
	expectedCrc := crc32.ChecksumIEEE(data)
	dataWithCrc := append(data, make([]byte, 4)...)
	binary.BigEndian.PutUint32(dataWithCrc[len(data):], expectedCrc)
	err := dll.transport.Write(dataWithCrc)
	if err != nil {
		dll.errTxCount++
		return fmt.Errorf("failed to write data: %w", err)
	}
	return nil
}

func (dll *DataLayer) Read() ([]byte, error) {
	Bytes, errRead := dll.transport.Read()
	if errRead != nil {
		dll.errRxCount++
		return nil, errRead
	}
	if err := ChecksumIEEE(Bytes); err != nil {
		dll.errRxCount++
		return nil, err
	}
	return Bytes, nil
}

func (dll *DataLayer) Status() error {
	return dll.transport.Status()
}
func (dll *DataLayer) GetErrTxCount() int32 {
	return dll.errTxCount
}
func (dll *DataLayer) GetErrRxCount() int32 {
	return dll.errRxCount
}
func (dll *DataLayer) Close() error {
	return dll.transport.Close()
}

func ChecksumIEEE(buffer []byte) error {
	data := buffer[:len(buffer)-4]
	expectedCrc := binary.BigEndian.Uint32(buffer[len(buffer)-4:])
	actualCrc := crc32.ChecksumIEEE(data)
	if expectedCrc != actualCrc {
		return fmt.Errorf("CRC Check Failed: Expected CRC = 0x%08X, Current CRC = 0x%08X", expectedCrc, actualCrc)
	}
	return nil
}
