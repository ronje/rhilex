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
 * @description: 换行符协议
 *
 */

func StartNewLineReceive(Ctx context.Context,
	OutChannel chan []byte,
	InputIO io.ReadWriteCloser) error {
	edge0D := byte(0x0D)
	edge0A := byte(0x0A)
	MAX_BUFFER_SIZE := 1024 * 10 * 10
	buffer := make([]byte, MAX_BUFFER_SIZE)
	byteACC := 0                      // 计数器而不是下标
	expectPacket := make([]byte, 256) // 默认最大包长256字节
	for {
		select {
		case <-Ctx.Done():
			return nil
		default:
		}
		N, errR := InputIO.Read(buffer[byteACC:])
		if errR != nil {
			continue
		}
		if N == 0 {
			continue
		}
		byteACC += N
		if byteACC > 256 {
			for i := 0; i < byteACC; i++ {
				buffer[i] = '\x00'
			}
			byteACC = 0
			continue
		}
	PARSE:
		expectPacketACC := 0
		expectPacketLength := 0
		for _, currentByte := range buffer[:byteACC] {
			expectPacketACC++
			if expectPacketACC < 4 {
				continue
			}
			if currentByte == edge0A && buffer[expectPacketACC-2] == edge0D {
				expectPacketLength = copy(expectPacket, buffer[:expectPacketACC-2])
				DataPacket := [256]byte{}
				DataPacketN := copy(DataPacket[0:], expectPacket[:expectPacketLength])
				OutChannel <- DataPacket[:DataPacketN]
				lessMoreBytesCount := byteACC - expectPacketACC
				if lessMoreBytesCount < byteACC {
					copiedCount := copy(buffer[0:], buffer[expectPacketACC:byteACC])
					byteACC = copiedCount
					if lessMoreBytesCount > 4 {
						goto PARSE
					}
				} else {
					byteACC = 0
					expectPacketACC = 0
					expectPacketLength = 0
				}
			}
		}
	}
}
