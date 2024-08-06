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

const MAX_BUFFER_SIZE = 1024 * 10 * 10

/*
*
* EE EF .... \r \n
*
 */
func StartDelimiterReceive(Ctx context.Context,
	OutChannel chan []byte,
	InputIO io.ReadWriteCloser) error {
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
					if currentByte == 0xEF && buffer[i-1] == 0xEE {
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

/*
*
* 固定包格式
*
 */
type BinaryPacket struct {
	_type  uint8     // 数据包类型
	length uint32    // 数据包长度
	data   [256]byte // 数据体
}

func NewBinaryPacket(data []byte) BinaryPacket {
	B := BinaryPacket{}
	copy(B.data[:], data)
	return B
}
func (B BinaryPacket) Type() {

}
func (B BinaryPacket) Length() uint32 {
	return binary.BigEndian.Uint32(B.data[:4])
}

/*
*
* Header[4]| Data[N]
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

/*
*
* EE EF ········ \r \n
*
 */
func StartEdgeSymReceive(Ctx context.Context,
	OutChannel chan []byte,
	InputIO io.ReadWriteCloser) error {
	edge1Sym := byte(0xEE)
	edge2Sym := byte(0xEF)
	edge3Sym := byte(0x0D)
	edge4Sym := byte(0x0A)
	MAX_BUFFER_SIZE := 1024 * 10 * 10
	buffer := make([]byte, MAX_BUFFER_SIZE)
	byteACC := 0
	edgeSignal1 := false
	edgeSignal2 := false
	expectPacket := make([]byte, MAX_BUFFER_SIZE)
	for {
		select {
		case <-Ctx.Done():
			return nil
		default:
		}
		N, errR := InputIO.Read(buffer[byteACC:])
		// 读取异常，重启
		if errR != nil {
			continue
		}
		byteACC += N
		if byteACC < 4 {
			continue
		}
		if byteACC > 256 {
			for i := 0; i < byteACC; i++ {
				buffer[i] = '\x00'
			}
			byteACC = 0
			edgeSignal1 = false
			edgeSignal2 = false
			continue
		}
	PARSE:
		expectPacketACC := 0
		expectPacketLength := 0
		for i, currentByte := range buffer[:byteACC] {
			expectPacketACC++
			if !edgeSignal1 {
				if expectPacketACC == 2 {
					if currentByte == edge2Sym && buffer[i-1] == edge1Sym {
						edgeSignal1 = true
					}
				}
			}
			if !edgeSignal1 || expectPacketACC < 4 {
				continue
			}
			if edgeSignal1 {
				if currentByte == edge4Sym && buffer[expectPacketACC-2] == edge3Sym {
					expectPacketLength = copy(expectPacket, buffer[2:expectPacketACC-2])
					edgeSignal2 = true
				}
			}
			if edgeSignal1 && edgeSignal2 {
				DataPacket := [256]byte{}
				DataPacketN := copy(DataPacket[0:], expectPacket[:expectPacketLength])
				OutChannel <- DataPacket[:DataPacketN]
				edgeSignal1 = false
				edgeSignal2 = false
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

/*
*
* ········\r\n
*
 */
func StartNewLineLoopReceive(Ctx context.Context,
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
