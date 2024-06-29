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
	"encoding/binary"
	"fmt"

	"github.com/hootrhino/rhilex/utils"
)

// 定义报文结构体
type Packet struct {
	StartFlag   byte
	MessageType byte
	DeviceID    byte
	MessageID   uint16
	DataLength  byte
	Data        []byte
	Crc         uint16
	EndFlag     byte
}

// 假设已经有一个名为crc16的函数可用
func crc16(data []byte) uint16 {
	return utils.CRC16(data)
}

// 初始化报文结构体的函数
func NewPacket(t byte, dID byte, mID uint16, d []byte, dLen int) Packet {
	p := Packet{}
	p.StartFlag = 0xFF
	p.MessageType = t
	p.DeviceID = dID
	p.MessageID = mID
	p.DataLength = byte(dLen)
	p.Data = make([]byte, dLen)
	copy(p.Data, d)
	p.Crc = crc16(append([]byte{p.StartFlag, p.MessageType, p.DeviceID, byte(p.MessageID >> 8),
		byte(p.MessageID & 0xFF), p.DataLength}, p.Data...))
	p.EndFlag = 0xFE
	return p
}

// 打印报文内容的函数
func (p Packet) PrintPacket() {
	fmt.Printf("Start Flag: 0x%02X\n", p.StartFlag)
	fmt.Printf("Message Type: %d\n", p.MessageType)
	fmt.Printf("Device ID: %d\n", p.DeviceID)
	fmt.Printf("Message ID: %d\n", p.MessageID)
	fmt.Printf("Data Length: %d\n", p.DataLength)
	fmt.Printf("Data: %s\n", string(p.Data))
	fmt.Printf("CRC: 0x%04X\n", p.Crc)
	fmt.Printf("End Flag: 0x%02X\n", p.EndFlag)
}

// 构造主动上报报文的函数
func ConstructReportPacket(deviceID byte, messageID uint16, data string) Packet {
	return NewPacket(0x01, deviceID, messageID, []byte(data), len(data))

}

// 构造网关下发命令报文的函数
func ConstructCommandPacket(deviceID byte, messageID uint16, data string) Packet {
	return NewPacket(0x02, deviceID, messageID, []byte(data), len(data))
}

// 构造设备回复报文的函数
func ConstructReplyPacket(deviceID byte, messageID uint16, data string) Packet {
	return NewPacket(0x03, deviceID, messageID, []byte(data), len(data))
}

// 解析数据包的函数
func ParsePacket(packetData []byte) (Packet, error) {
	if len(packetData) < 10 { // 确保至少有足够的空间容纳基本字段
		return Packet{}, fmt.Errorf("invalid packet data length: %d", len(packetData))
	}

	p := Packet{}
	buf := bytes.NewBuffer(packetData)

	// 读取字段值
	err := binary.Read(buf, binary.BigEndian, &p.StartFlag)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.MessageType)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.DeviceID)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.MessageID)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.DataLength)
	if err != nil {
		return Packet{}, err
	}
	p.Data = make([]byte, p.DataLength)
	_, err = buf.Read(p.Data)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.Crc)
	if err != nil {
		return Packet{}, err
	}
	err = binary.Read(buf, binary.BigEndian, &p.EndFlag)
	if err != nil {
		return Packet{}, err
	}

	// 验证CRC
	expectedCrc := crc16(packetData[:len(packetData)-2]) // 不包括CRC和EndFlag字段
	if expectedCrc != p.Crc {
		return Packet{}, fmt.Errorf("CRC check failed: expected 0x%04X, got 0x%04X", expectedCrc, p.Crc)
	}

	// 验证结束标志
	if p.EndFlag != 0xFE {
		return Packet{}, fmt.Errorf("invalid end flag: 0x%02X", p.EndFlag)
	}

	return p, nil
}
func TestSRA() {
	// 模拟接收到的数据包
	receivedPacketData := []byte{
		0xFF, 0x01, 0x01, 0x00, 0x01, 0x0C,
		'R', 'e', 'p', 'o', 'r', 't', ' ', 'D', 'a', 't', 'a',
		0xAB, 0xCD, 0xFE}

	// 解析数据包
	parsedPacket, err := ParsePacket(receivedPacketData)
	if err != nil {
		fmt.Println("Error parsing packet:", err)
		return
	}
	parsedPacket.PrintPacket()
	// 构造主动上报报文
	reportPacket := ConstructReportPacket(0x01, 0x0001, "Report Data")
	reportPacket.PrintPacket()

	// 构造网关下发命令报文
	commandPacket := ConstructCommandPacket(0x01, 0x0002, "Command Data")
	commandPacket.PrintPacket()

	// 构造设备回复报文
	replyPacket := ConstructReplyPacket(0x01, 0x0002, "Reply Data")
	replyPacket.PrintPacket()
}
