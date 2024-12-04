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
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type TransportLayer struct {
	config TransporterConfig
}

func NewTransportLayer(config TransporterConfig) *TransportLayer {
	return &TransportLayer{config: config}
}

func (transport *TransportLayer) Write(data []byte) error {
	transport.config.Logger.Debug("TransportLayer.Write=", ByteDumpHexString(data))
	transport.config.Port.SetWriteDeadline(time.Now().Add(
		transport.config.WriteTimeout * time.Millisecond))
	if _, err := transport.config.Port.Write(data); err != nil {
		return err
	}
	return transport.config.Port.SetWriteDeadline(time.Time{})
}

func (transport *TransportLayer) Read() ([]byte, error) {
	const (
		maxPacketSize = 512
		startByte1    = 0xAA
		startByte2    = 0xBB
	)

	var buffer []byte
	tempBuffer := make([]byte, 1)
	errorBuffer := make([]byte, 0)

	// 设置读取超时
	transport.config.Port.SetReadDeadline(time.Now().Add(transport.config.ReadTimeout))

	for {
		// 如果 errorBuffer 有数据，则优先恢复数据到缓冲区
		if len(errorBuffer) > 0 {
			buffer = append(buffer, errorBuffer...)
			errorBuffer = errorBuffer[:0]
		}

		// 同步寻找包头
		for {
			_, err := transport.config.Port.Read(tempBuffer)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("connection closed")
				}
				return nil, fmt.Errorf("read error: %w", err)
			}

			// 识别包头起始标志
			if len(buffer) == 0 && tempBuffer[0] == startByte1 {
				buffer = append(buffer, tempBuffer[0])
				continue
			} else if len(buffer) == 1 && tempBuffer[0] == startByte2 {
				buffer = append(buffer, tempBuffer[0])
				break
			} else {
				// 丢弃不符合包头的数据
				buffer = buffer[:0]
			}
		}

		// 读取包头剩余部分
		headerBytes := make([]byte, 4)
		_, err := transport.config.Port.Read(headerBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read full header: %w", err)
		}
		buffer = append(buffer, headerBytes...)

		// 解码 Header 并获取包长
		header, err := DecodeHeader(buffer[:6])
		if err != nil {
			return nil, fmt.Errorf("failed to decode header: %w", err)
		}
		packetLength := binary.BigEndian.Uint16(header.Length[:])

		// 检查长度是否合法
		if packetLength == 0 || int(packetLength) > maxPacketSize {
			buffer = buffer[:0] // 重置缓冲区并继续寻找新帧
			continue
		}

		// 读取 payload、校验码和分隔符
		payloadCrcDelimiterSize := int(packetLength) + 4 + 2 // payload + crc + reverse delimiter
		data := make([]byte, payloadCrcDelimiterSize)
		n, err := io.ReadFull(transport.config.Port, data)
		if err != nil || n < payloadCrcDelimiterSize {
			// 数据不完整，缓存错误数据
			errorBuffer = append(errorBuffer, buffer...)
			errorBuffer = append(errorBuffer, data[:n]...)
			buffer = buffer[:0]
			continue
		}

		// 检查下一个帧的标志符来验证长度
		for i := int(packetLength); i < len(data)-1; i++ {
			if data[i] == startByte1 && data[i+1] == startByte2 {
				// 如果找到下一个包的开始标志符，丢弃当前帧
				errorBuffer = append(errorBuffer, buffer...)
				errorBuffer = append(errorBuffer, data[:n]...)
				buffer = buffer[:0]
				continue
			}
		}

		// 拼接完整数据帧，完成解析
		buffer = append(buffer, data...)
		frame, err := DecodeAppLayerFrame(buffer)
		if err != nil {
			return nil, fmt.Errorf("failed to decode frame: %w", err)
		}

		// 成功解析帧后返回 payload
		return frame.Payload, nil
	}
}

func (transport *TransportLayer) Status() error {
	if transport.config.Port == nil {
		return fmt.Errorf("invalid port")
	}
	_, err := transport.config.Port.Write([]byte{0})
	return err
}

func (transport *TransportLayer) Close() error {
	return transport.config.Port.Close()
}
