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
	"flag"
	"fmt"
	"time"

	"testing"

	"github.com/hootrhino/rhilex/device/mbus"
	gombus "github.com/hootrhino/rhilex/device/mbus"
	"github.com/sirupsen/logrus"
)

// MbusReadWriteCloser 是一个简单的 ReadWriteCloser 实现
type MbusReadWriteCloser struct {
	buffer *bytes.Buffer
}

// NewMbusReadWriteCloser 创建一个新的 MbusReadWriteCloser 实例
func NewMbusReadWriteCloser() *MbusReadWriteCloser {
	return &MbusReadWriteCloser{
		buffer: new(bytes.Buffer),
	}
}

// Read 从 buffer 中读取数据
func (s *MbusReadWriteCloser) Read(p []byte) (n int, err error) {
	return s.buffer.Read(p)
}

// Write 向 buffer 中写入数据
func (s *MbusReadWriteCloser) Write(p []byte) (n int, err error) {
	fmt.Println("MbusReadWriteCloser.Write:", p)
	return s.buffer.Write(p)
}

// Close 清空 buffer 并模拟关闭连接
func (s *MbusReadWriteCloser) Close() error {
	s.buffer.Reset()
	return nil
}

func (s *MbusReadWriteCloser) SetReadDeadline(t time.Time) error {
	return nil

}
func (s *MbusReadWriteCloser) SetWriteDeadline(t time.Time) error {
	return nil

}

var primaryID = flag.Int("id", 1, "primaryID to fetch data from")

// go test -timeout 30s -run ^TestMbusCodec github.com/hootrhino/rhilex/test -v -count=1
func TestMbusCodec(t *testing.T) {
	client := mbus.NewMbusClientHandler(NewDltReadWriteCloser())
	client.SetLogger(logrus.StandardLogger())

	frame := client.SndNKE(uint8(*primaryID))
	fmt.Printf("sending nke: % x\n", frame)
	var err error
	_, err = client.Request(frame)
	if err != nil {
		logrus.Error(err)
		return
	}
	_, err = client.ReadSingleCharFrame()
	if err != nil {
		logrus.Error(err)
		return
	}

	// frame := gombus.SetPrimaryUsingPrimary(0, 3)
	respFrame := &gombus.DecodedFrame{}
	lastFCB := true
	frames := 0
	for respFrame.HasMoreRecords() || frames == 0 {
		frames++
		// frame := gombus.SetPrimaryUsingPrimary(0, 3)
		frame = client.RequestUD2(uint8(*primaryID))
		if !lastFCB {
			frame.SetFCB()
			frame.SetChecksum()
		}
		lastFCB = frame.C().FCB()

		fmt.Printf("sending: % x\n", frame)
		fmt.Printf("sending C: % b\n", frame.C())
		_, err = client.Request(frame)
		if err != nil {
			logrus.Error(err)
			return
		}

		resp, err := client.ReadLongFrame()
		if err != nil {
			logrus.Error(err)
			return
		}

		fmt.Printf("read: % x\n", resp)
		fmt.Printf("read C: % b\n", resp.C())

		respFrame, err = resp.Decode()
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("read total values: ", respFrame)
	}

	logrus.Info("read total frames: ", frames)
}
