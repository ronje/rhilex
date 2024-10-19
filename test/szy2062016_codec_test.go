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
	"time"

	"github.com/hootrhino/rhilex/device/szy2062016"
	"github.com/sirupsen/logrus"
)

// SzyReadWriteCloser 是一个简单的 ReadWriteCloser 实现
type SzyReadWriteCloser struct {
	buffer *bytes.Buffer
}

// NewSzyReadWriteCloser 创建一个新的 SzyReadWriteCloser 实例
func NewSzyReadWriteCloser() *SzyReadWriteCloser {
	return &SzyReadWriteCloser{
		buffer: new(bytes.Buffer),
	}
}

// Read 从 buffer 中读取数据
func (s *SzyReadWriteCloser) Read(p []byte) (n int, err error) {
	return s.buffer.Read(p)
}

// Write 向 buffer 中写入数据
func (s *SzyReadWriteCloser) Write(p []byte) (n int, err error) {
	fmt.Println("SzyReadWriteCloser.Write:", p)
	return s.buffer.Write(p)
}

// Close 清空 buffer 并模拟关闭连接
func (s *SzyReadWriteCloser) Close() error {
	s.buffer.Reset()
	return nil
}

func (s *SzyReadWriteCloser) SetReadDeadline(t time.Time) error {
	return nil

}
func (s *SzyReadWriteCloser) SetWriteDeadline(t time.Time) error {
	return nil

}

// go test -timeout 30s -run ^TestCodec_SZY206_2016_Frame github.com/hootrhino/rhilex/test -v -count=1
func TestCodec_SZY206_2016_Frame(t *testing.T) {
	client := szy2062016.NewSZY206ClientHandler(NewSzyReadWriteCloser())
	client.SetLogger(logrus.StandardLogger())
	frame := szy2062016.SZY206Frame0x00{
		Start1:     szy2062016.CTRL_CODE_FRAME_START,
		DataLength: 0x07,
		Start2:     szy2062016.CTRL_CODE_FRAME_START,
		CtrlCode:   szy2062016.CtrlCodeRainfall(),
		Address:    [5]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xFF},
		DataArea:   []byte{},
		End:        szy2062016.CTRL_CODE_FRAME_END,
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
	CJT188Frame0x01, err1 := client.DecodeSZY206Frame0x00(packedFrame)
	if err1 != nil {
		t.Fatal(err1)
	}
	t.Log(CJT188Frame0x01.String())
	// // 68 10 01 00 00 05 08 00 00 81 09 90 1F 00 00 23 01 00 00 FF E2 16
	// var data = []byte{0x68, 0x10, 0x01, 0x00, 0x00, 0x05, 0x08, 0x00, 0x00, 0x81,
	// 	0x09, 0x90, 0x1F, 0x00, 0x00, 0x23, 0x01, 0x00, 0x00, 0xFF, 0xE2, 0x16}
	// Frame0x01, err2 := client.DecodeSZY206Frame0x00Response(data)
	// if err2 != nil {
	// 	t.Fatal(err2)
	// }
	// t.Log(Frame0x01.String())
}
