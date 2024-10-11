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
	"encoding/binary"
	"io"
)

/*
*
* 固定报文长度协议：Header[4] | Data[N]
*
 */
func StartFixLengthReceive(Ctx context.Context,
	OutChannel chan []byte,
	InputIO io.ReadWriteCloser) error {
	receiveBuffer := make([]byte, 256)
	OneFrame := make([]byte, 256)
	bytesCursor := uint32(0)
	packetHeaderLength := uint32(4)
	segmentData := false
	for {
		select {
		case <-Ctx.Done():
			return nil
		default:
		}
		N, errR := InputIO.Read(receiveBuffer[bytesCursor:])
		if errR != nil {
			continue
		}
		bytesCursor += uint32(N)
		if bytesCursor < 4 {
			continue
		}
		if bytesCursor >= 256 {
			for i := 0; i < 256; i++ {
				receiveBuffer[i] = 0
			}
			bytesCursor = 0
			packetHeaderLength = 0
			segmentData = false
			continue
		}
	PARSE_PACKET:
		BinaryLength := binary.BigEndian.Uint32(receiveBuffer[:4])
		if BinaryLength > bytesCursor {
			// 解决死循环
			if segmentData {
				goto SEGMENT
			} else {
				continue
			}
		}
	SEGMENT:
		onePacketBytesCount := packetHeaderLength + BinaryLength
		if onePacketBytesCount > bytesCursor {
			continue
		}
		copiedBytesCount := copy(OneFrame, receiveBuffer[:onePacketBytesCount])
		DataPacket := [256]byte{}
		DataPacketN := copy(DataPacket[:0], OneFrame[:copiedBytesCount])
		OutChannel <- DataPacket[:DataPacketN]
		leastMoreBytesCount := copy(receiveBuffer[0:], receiveBuffer[onePacketBytesCount:bytesCursor])
		bytesCursor -= onePacketBytesCount
		if leastMoreBytesCount > 4 {
			segmentData = true
			goto PARSE_PACKET
		} else {
			segmentData = false
			goto PARSE_PACKET
		}
	}
}
