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
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hootrhino/rhilex/protocol"
	"github.com/sirupsen/logrus"
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
	v := []byte{0x00, 0x01, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x0B}
	copy(p, v)
	fmt.Println("GenericReadWriteCloser.Read2:", p)
	return 8, nil
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

// go test -timeout 30s -run ^TestGenericProtocolTest github.com/hootrhino/rhilex/test -v -count=1
func TestGenericProtocolTest(t *testing.T) {
	config := protocol.TransporterConfig{
		Port:         NewGenericReadWriteCloser(),
		ReadTimeout:  2000,
		WriteTimeout: 2000,
	}
	ProtocolHandler := protocol.NewGenericProtocolHandler(config)
	appLayerFrame := protocol.AppLayerFrame{
		Header: protocol.Header{
			Type:   [2]byte{0x00, 0x01},
			Length: [2]byte{0x00, 0x04},
		},
		Payload: []byte{0, 1, 2, 3},
	}
	response, err := ProtocolHandler.Request(appLayerFrame)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

// go test -timeout 30s -run ^TestGenericProtocolSlaverTest github.com/hootrhino/rhilex/test -v -count=1
// 00 01 00 05 01 02 03 04 10
func TestGenericProtocolSlaverTest(t *testing.T) {
	Listener, err := net.Listen("tcp", ":7799")
	if err != nil {
		t.Fatal(err)
	}
	defer Listener.Close()
	t.Log("Server listening on port 7799")
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		fmt.Println("Accepting connection:", conn.RemoteAddr())
		Logger := logrus.StandardLogger()
		Logger.SetLevel(logrus.DebugLevel)
		config := protocol.TransporterConfig{
			Port:         conn,
			ReadTimeout:  5000,
			WriteTimeout: 5000,
			Logger:       Logger,
		}
		ctx, cancel := context.WithCancel(context.Background())
		TransportSlaver := protocol.NewGenericProtocolSlaver(ctx, cancel, config)
		go TransportSlaver.StartLoop(func(AppLayerFrame protocol.AppLayerFrame, err error) {
			if err != nil {
				t.Log(err)
			} else {
				t.Log(AppLayerFrame.String())
			}
		})
	}
}

// go test -timeout 30s -run ^TestParseBinaryData github.com/hootrhino/rhilex/test -v -count=1

func TestParseBinaryData(t *testing.T) {
	{
		expr := "ID:32:int:BE; Name:40:string:BE; Age:16:int:LE"
		data := []byte{0x00, 0x00, 0x00, 0x01, 'A', 'l', 'i', 'c', 'e', 0x00, 0x20}
		t.Log("解析:", data)
		parsedData, err := protocol.ParseBinary(expr, data)
		if err != nil {
			t.Fatal("解析失败:", err)
			return
		}

		t.Log("解析结果:", parsedData.String())
	}
	{
		expr := "ID:32:int:BE; Name:40:string:BE; Age:16:int:LE"
		data := []byte{0x00, 0x00, 0x00, 0x01, 'A', 'l', 'i', 'c', 'e', 0x00, 0x20}
		t.Log("解析:", data)
		parsedData, err := protocol.ParseBinary(expr, data)
		if err != nil {
			t.Fatal("解析失败:", err)
			return
		}

		t.Log("解析结果:", parsedData.String())
	}
	{
		expr := "ID:32:int:BE; Name:40:string:BE; Age:16:int:LE"
		data := []byte{}
		t.Log("解析:", data)
		parsedData, err := protocol.ParseBinary(expr, data)
		if err != nil {
			t.Fatal("解析失败:", err)
			return
		}

		t.Log("解析结果:", parsedData.String())
	}
}
