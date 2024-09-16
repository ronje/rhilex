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
	"fmt"
	"log"
	"time"

	"testing"

	"github.com/hootrhino/gombus"
	serial "github.com/hootrhino/goserial"
	"github.com/sirupsen/logrus"
)

// NKEPacket 定义了NKE包的结构（示例）
type NKEPacket struct {
	Header  byte
	Payload []byte
}

// SingleCharFrame 定义了SingleCharFrame响应格式
type SingleCharFrame struct {
	Header  byte
	Payload byte
}

// handleSerialPort 处理串口通信中的NKE数据
func StartSerialModeServer(path string) {
	port17, err := gombus.OpenSerial(serial.Config{
		Address:  path,
		BaudRate: 9600,
		DataBits: 8,
		Parity:   "N",
		StopBits: 1,
	})
	if err != nil {
		log.Fatal("Error opening serial port:", err)
	}
	defer port17.Close()
	buffer := make([]byte, 5)
	for {
		n, err := port17.Read(buffer)
		if err != nil {
			log.Println("Error reading from serial port:", err)
			continue
		}
		log.Println("port17.Read === ", buffer[:n])
		fixedResponse := []byte{0x68, 0x04, 0x04, 0x68, 0x08, 0x32, 0xFF, 0xFF, 0xE5, 0x16}
		log.Println("Single CharFrame Response === ", fixedResponse)
		_, err = port17.Write(fixedResponse)
		if err != nil {
			log.Println("Error writing response to serial port:", err)
		}

	}
}

var primaryID byte = 1

// go test -timeout 30s -run ^TestSerialMode github.com/hootrhino/rhilex/test -v -count=1
func TestSerialMode(t *testing.T) {
	go StartSerialModeServer("COM17")
	serial_transport_test("COM16")
}

// go test -timeout 30s -run ^TestTcpMode github.com/hootrhino/rhilex/test -v -count=1
func TestTcpMode(t *testing.T) {
	tcp_transport_test()
}
func serial_transport_test(path string) {

	conn, err := gombus.OpenSerial(serial.Config{
		Address:  path,
		BaudRate: 9600,
		DataBits: 8,
		Parity:   "N",
		StopBits: 1,
		Timeout:  3 * time.Second,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	defer conn.Close()

	frame := gombus.SndNKE(uint8(primaryID))
	log.Println("gombus.SndNKE: ", frame)
	_, err = conn.Write(frame)
	if err != nil {
		logrus.Error(err)
		return
	}
	_, err = gombus.ReadSingleCharFrame(conn)
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
		frame = gombus.RequestUD2(uint8(primaryID))
		if !lastFCB {
			frame.SetFCB()
			frame.SetChecksum()
		}
		lastFCB = frame.C().FCB()

		fmt.Printf("sending: % x\n", frame)
		fmt.Printf("sending C: % b\n", frame.C())
		_, err = conn.Write(frame)
		if err != nil {
			logrus.Error(err)
			return
		}

		resp, err := gombus.ReadLongFrame(conn)
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
		logrus.Info("read total values: ", len(respFrame.DataRecords))
	}

	logrus.Info("read total frames: ", frames)
}
func tcp_transport_test() {

	conn, err := gombus.OpenTCP("192.168.13.42:10001")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer conn.Close()

	frame := gombus.SndNKE(uint8(primaryID))
	fmt.Printf("sending nke: % x\n", frame)
	_, err = conn.Write(frame)
	if err != nil {
		logrus.Error(err)
		return
	}
	_, err = gombus.ReadSingleCharFrame(conn)
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
		frame = gombus.RequestUD2(uint8(primaryID))
		if !lastFCB {
			frame.SetFCB()
			frame.SetChecksum()
		}
		lastFCB = frame.C().FCB()

		fmt.Printf("sending: % x\n", frame)
		fmt.Printf("sending C: % b\n", frame.C())
		_, err = conn.Write(frame)
		if err != nil {
			logrus.Error(err)
			return
		}

		resp, err := gombus.ReadLongFrame(conn)
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
		logrus.Info("read total values: ", len(respFrame.DataRecords))
	}

	logrus.Info("read total frames: ", frames)
}
