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
	"fmt"
	"time"
)

// GenericReadWriteCloser 是一个简单的 ReadWriteCloser 实现
type GenericReadWriteCloser struct {
	buffer *bytes.Buffer
}

func NewGenericReadWriteCloser() *GenericReadWriteCloser {
	return &GenericReadWriteCloser{
		buffer: new(bytes.Buffer),
	}
}
func (s *GenericReadWriteCloser) Read(p []byte) (n int, err error) {
	header := []byte{0xAB, 0xAB}
	testData := []byte{
		0x00, 0x05, // Length (5 bytes)
		0x68, 0x65, 0x6C, 0x6C, 0x6F, // Payload ("hello")
		0x1C, 0x2F, // CRC16 (0x1C2F)
	}
	tail := []byte{0xBA, 0xBA}
	packet := bytes.Join([][]byte{header, testData, tail}, []byte{})
	p = append(p, packet...)
	fmt.Println("GenericReadWriteCloser.Read:", p)
	return len(packet), nil
}
func (s *GenericReadWriteCloser) Write(p []byte) (n int, err error) {
	fmt.Println("GenericReadWriteCloser.Write:", p)
	return s.buffer.Write(p)
}
func (s *GenericReadWriteCloser) Close() error {
	s.buffer.Reset()
	return nil
}

func (s *GenericReadWriteCloser) SetReadDeadline(t time.Time) error {
	return nil

}
func (s *GenericReadWriteCloser) SetWriteDeadline(t time.Time) error {
	return nil

}
