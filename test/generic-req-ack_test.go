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
	"testing"

	"github.com/hootrhino/rhilex/protocol"
)

func Test_SimpleReqAck(t *testing.T) {
	// 模拟接收到的数据包
	receivedPacketData := []byte{
		0xFF, 0x01, 0x01, 0x00, 0x01, 0x0C,
		'R', 'e', 'p', 'o', 'r', 't', ' ', 'D', 'a', 't', 'a',
		0xAB, 0xCD, 0xFE}

	// 解析数据包
	parsedPacket, err := protocol.ParsePacket(receivedPacketData)
	if err != nil {
		fmt.Println("Error parsing packet:", err)
		return
	}
	parsedPacket.PrintPacket()
	// 构造主动上报报文
	reportPacket := protocol.ConstructReportPacket(0x01, 0x0001, "Report Data")
	reportPacket.PrintPacket()

	// 构造网关下发命令报文
	commandPacket := protocol.ConstructCommandPacket(0x01, 0x0002, "Command Data")
	commandPacket.PrintPacket()

	// 构造设备回复报文
	replyPacket := protocol.ConstructReplyPacket(0x01, 0x0002, "Reply Data")
	replyPacket.PrintPacket()
}
