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

	"github.com/hootrhino/rhilex/device/dlt6452007"
	"github.com/sirupsen/logrus"

	serial "github.com/hootrhino/goserial"
)

// SimpleReadWriteCloser 是一个简单的 ReadWriteCloser 实现
type SimpleReadWriteCloser struct {
	buffer *bytes.Buffer
}

// NewSimpleReadWriteCloser 创建一个新的 SimpleReadWriteCloser 实例
func NewSimpleReadWriteCloser() *SimpleReadWriteCloser {
	return &SimpleReadWriteCloser{
		buffer: new(bytes.Buffer),
	}
}

// Read 从 buffer 中读取数据
func (s *SimpleReadWriteCloser) Read(p []byte) (n int, err error) {
	return s.buffer.Read(p)
}

// Write 向 buffer 中写入数据
func (s *SimpleReadWriteCloser) Write(p []byte) (n int, err error) {
	return s.buffer.Write(p)
}

// Close 清空 buffer 并模拟关闭连接
func (s *SimpleReadWriteCloser) Close() error {
	s.buffer.Reset()
	return nil
}

// go test -timeout 30s -run ^TestCodec_DLT645_2007_Frame github.com/hootrhino/rhilex/test -v -count=1
// 68 45 92 66 23 00 10 68 11 04 33 34 34 35 25 16
// 68 45 92 66 23 00 10 68 91 06 33 34 34 35 66 55 62 16
// ============================= 33 34 34 35 66 55
// ============================= 00 01 01 02 33 22
func TestCodec_DLT645_2007_Frame(t *testing.T) {
	client := dlt6452007.NewDLT645ClientHandler(NewSimpleReadWriteCloser())
	frame := dlt6452007.DLT645Frame0x11{
		Start:      dlt6452007.CTRL_CODE_FRAME_START,
		Address:    [6]byte{0x45, 0x92, 0x66, 0x23, 0x00, 0x10},
		CtrlCode:   dlt6452007.CTRL_CODE_READ_DATA,
		DataLength: 0x04,
		DataType:   [4]byte{0x33, 0x34, 0x34, 0x35},
		End:        dlt6452007.CTRL_CODE_FRAME_END,
	}
	packedFrame, err := frame.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("frame.Encode():")
	println()
	for _, v := range packedFrame {
		fmt.Printf(" %x", v)
	}
	println()
	DLT645Frame0x11, err2 := client.DecodeDLT645Frame0x11(packedFrame)
	if err != nil {
		panic(err2)
	}
	t.Log("dlt6452007.DecodeDLT645Frame0x11:", DLT645Frame0x11.String())
	responseFrame, err := client.DecodeDLT645Frame0x11Response([]byte{
		0x68,
		0x45, 0x92, 0x66, 0x23, 0x00, 0x10,
		0x68,
		0x91,
		0x06, 0x33, 0x34, 0x34, 0x35, 0x66, 0x55, 0x62,
		0x16,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("DecodeDLT645Frame0x11Response: %+v\n", responseFrame.String())
	Data1, err1 := responseFrame.GetDataType()
	if err1 != nil {
		panic(err1)
	}
	t.Log(Data1)
	Data2, err2 := responseFrame.GetData()
	if err2 != nil {
		panic(err2)
	}
	t.Log(Data2)
}

// go test -timeout 30s -run ^TestCodec_DLT645_2007_Meter github.com/hootrhino/rhilex/test -v -count=1

func TestCodec_DLT645_2007_Meter(t *testing.T) {
	port, err := serial.Open(&serial.Config{
		Address:  "COM9",
		BaudRate: 2400,
		DataBits: 8,
		StopBits: 1,
		Parity:   "E",
		// Timeout:  utils.GiveMeSeconds(3),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer port.Close()
	client := dlt6452007.NewDLT645ClientHandler(port)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	client.SetLogger(logger)
	// 68 45 92 66 23 00 10 68 11 04 33 34 34 35 25 16
	frame := dlt6452007.DLT645Frame0x11{
		Start:      dlt6452007.CTRL_CODE_FRAME_START,
		Address:    [6]byte{0x45, 0x92, 0x66, 0x23, 0x00, 0x10},
		CtrlCode:   dlt6452007.CTRL_CODE_READ_DATA,
		DataLength: 0x04,
		DataType:   [4]byte{0x33, 0x34, 0x34, 0x35},
		End:        dlt6452007.CTRL_CODE_FRAME_END,
	}
	packedFrame, err := frame.Encode()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("packedFrame=", packedFrame)
	resp, err2 := client.Request(packedFrame)
	if err2 != nil {
		t.Fatal(err2)
	}
	t.Log("client.Request == ", resp)
	frameResp, err4 := client.DecodeDLT645Frame0x11Response(resp)
	if err4 != nil {
		t.Fatal(err4)
	}
	// 68 45 92 66 23 00 10 68 91 06 33 34 34 35 A3 54 9E
	// 0x68 0x45 0x92 0x66 0x23 0x00 0x10 0x68 0x91 0x06 0x33 0x34 0x34 0x35 0xa3 0x54 0x9e
	t.Log("client.DecodeDLT645Frame0x11Response == ", frameResp.String())
	t.Log(frameResp.GetData())
}
