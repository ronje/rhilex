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

package test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hootrhino/rhilex/device/cjt1882004"
	"github.com/sirupsen/logrus"
)

// CjtReadWriteCloser 是一个简单的 ReadWriteCloser 实现
type CjtReadWriteCloser struct {
	buffer *bytes.Buffer
}

// NewCjtReadWriteCloser 创建一个新的 CjtReadWriteCloser 实例
func NewCjtReadWriteCloser() *CjtReadWriteCloser {
	return &CjtReadWriteCloser{
		buffer: new(bytes.Buffer),
	}
}

// Read 从 buffer 中读取数据
func (s *CjtReadWriteCloser) Read(p []byte) (n int, err error) {
	return s.buffer.Read(p)
}

// Write 向 buffer 中写入数据
func (s *CjtReadWriteCloser) Write(p []byte) (n int, err error) {
	fmt.Println("CjtReadWriteCloser.Write:", p)
	return s.buffer.Write(p)
}

// Close 清空 buffer 并模拟关闭连接
func (s *CjtReadWriteCloser) Close() error {
	s.buffer.Reset()
	return nil
}

// go test -timeout 30s -run ^TestCodec_CJT188_2007_Frame github.com/hootrhino/rhilex/test -v -count=1
func TestCodec_CJT188_2007_Frame(t *testing.T) {
	client := cjt1882004.NewCJT188ClientHandler(NewCjtReadWriteCloser())
	client.SetLogger(logrus.StandardLogger())
	frame := cjt1882004.CJT188Frame0x01{
		Start:        cjt1882004.CTRL_CODE_FRAME_START,
		MeterType:    0x10,
		Address:      [7]byte{0x01, 0x00, 0x00, 0x05, 0x08, 0x00, 0x00},
		CtrlCode:     cjt1882004.CTRL_CODE_READ_DATA,
		DataLength:   0x03,
		DataType:     [2]byte{0x90, 0x1F},
		DataArea:     []byte{},
		SerialNumber: 0x00,
		End:          cjt1882004.CTRL_CODE_FRAME_END,
	}
	t.Log(frame.String())
	packedFrame, errEncode := frame.Encode()
	if errEncode != nil {
		t.Fatal(errEncode)
	}
	t.Log("frame.Encode():")
	println()
	for _, v := range packedFrame {
		fmt.Printf(" 0x%x", v)
	}
	println()
	client.Request(packedFrame)
	CJT188Frame0x01, err1 := client.DecodeCJT188Frame0x01(packedFrame)
	if err1 != nil {
		t.Fatal(err1)
	}
	t.Log(CJT188Frame0x01.String())
	// 68 10 01 00 00 05 08 00 00 81 09 90 1F 00 00 23 01 00 00 FF E2 16
	var data = []byte{0x68, 0x10, 0x01, 0x00, 0x00, 0x05, 0x08, 0x00, 0x00, 0x81,
		0x09, 0x90, 0x1F, 0x00, 0x00, 0x23, 0x01, 0x00, 0x00, 0xFF, 0xE2, 0x16}
	Frame0x01, err2 := client.DecodeCJT188Frame0x01Response(data)
	if err2 != nil {
		t.Fatal(err2)
	}
	t.Log(Frame0x01.String())
}
