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
	"context"
	"io"
)

/**
 *
 *
 * @description: 特殊标记结尾的协议
 *
 * @param param description
 * @return return description
 */

func StartSpecialEndSymbolReceive(Ctx context.Context,
	OutChannel chan []byte,
	InputIO io.ReadWriteCloser, E1, E2 uint8) error {
	buffer := make([]byte, MAX_BUFFER_SIZE)
	byteACC := 0                      // 计数器而不是下标
	edgeSignal1 := false              // 两个边沿
	edgeSignal2 := false              // 两个边沿
	expectPacket := make([]byte, 256) // 默认最大包长256字节
	for {
		select {
		case <-Ctx.Done():
			return nil
		default:
		}
		N, errR := InputIO.Read(buffer[byteACC:])
		// 读取异常，重启
		if errR != nil {
			return errR
		}
		if N == 0 {
			continue
		}
		byteACC += N
		if byteACC > 256 { // 单个包最大256字节
			if !edgeSignal1 || !edgeSignal2 {
				for i := 0; i < byteACC; i++ {
					buffer[i] = '\x00'
				}
				byteACC = 0
				edgeSignal1 = false
				edgeSignal2 = false
			}
			continue
		}
		expectPacketACC := 0
		expectPacketLength := 0
		for i, currentByte := range buffer[:byteACC] {
			expectPacketACC++
			if !edgeSignal1 {
				if expectPacketACC >= 2 {
					if currentByte == E2 && buffer[i-1] == E1 {
						edgeSignal1 = true
					}
				}
			}
			if !edgeSignal1 || expectPacketACC < 4 {
				continue
			}
			if edgeSignal1 {
				if currentByte == 0x0A && buffer[expectPacketACC-2] == 0x0D {
					expectPacketLength = copy(expectPacket, buffer[:expectPacketACC-1])
					edgeSignal2 = true
				}
			}
			if !edgeSignal1 || !edgeSignal2 {
				continue
			}
			if edgeSignal1 && edgeSignal2 {
				DataPacket := [256]byte{}
				DataPacketN := copy(DataPacket[0:], expectPacket[2:expectPacketLength-1])
				OutChannel <- DataPacket[:DataPacketN]
			}
			if expectPacketACC < byteACC {
				if !edgeSignal1 || !edgeSignal2 {
					copy(buffer[0:], buffer[expectPacketACC-1:byteACC])
					byteACC = byteACC - expectPacketACC
				}
			} else {
				byteACC = 0
			}
			expectPacketLength = 0
			expectPacketACC = 0
			edgeSignal1 = false
			edgeSignal2 = false
		}
	}
}
