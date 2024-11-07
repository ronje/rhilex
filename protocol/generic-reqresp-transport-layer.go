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
	if transport.config.WriteTimeout != 0 {
		if err := transport.config.Port.SetWriteDeadline(time.Now().Add(
			transport.config.WriteTimeout * time.Millisecond)); err != nil {
			return err
		}
	}

	if _, err := transport.config.Port.Write(data); err != nil {
		return err
	}
	return transport.config.Port.SetWriteDeadline(time.Time{})
}

func (transport *TransportLayer) Read() ([]byte, error) {

	maxPacketSize := uint16(512)
	var buffer []byte
	temp := make([]byte, 1)
	errorBuffer := make([]byte, 0)
	transport.config.Port.SetReadDeadline(time.Now().Add(transport.config.ReadTimeout * time.Millisecond))
	for {
		// 如果 errorBuffer 有数据，则先合并
		if len(errorBuffer) > 0 {
			buffer = append(buffer, errorBuffer...)
			errorBuffer = errorBuffer[:0]
		}

		// 寻找包头同步
		for {
			_, err := transport.config.Port.Read(temp)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("connection closed")
				}
				return nil, fmt.Errorf("read error: %w", err)
			}

			// 如果缓冲区为空且匹配包头的第一字节，继续读取
			if len(buffer) == 0 && temp[0] == 0xAA { // 假设包头是 0xAA
				buffer = append(buffer, temp[0])
				continue
			} else if len(buffer) == 1 && temp[0] == 0xBB { // 假设包头第二字节是 0xBB
				buffer = append(buffer, temp[0])
				break
			} else {
				// 丢弃不符合包头的数据
				buffer = buffer[:0]
			}
		}

		// 读取包头
		header := Header{}
		_, err := transport.config.Port.Read(header.Type[:])
		if err != nil {
			return nil, fmt.Errorf("failed to read packet type: %w", err)
		}

		_, err = transport.config.Port.Read(header.Length[:])
		if err != nil {
			return nil, fmt.Errorf("failed to read packet length: %w", err)
		}

		// 获取数据包长度
		packetLength := binary.BigEndian.Uint16(header.Length[:])
		if packetLength > maxPacketSize || packetLength == 0 {
			// 如果包的长度不合法，丢弃当前包并恢复
			fmt.Println("Error: Invalid packet length, resynchronizing.")
			buffer = buffer[:0]
			continue
		}

		// 读取数据内容
		data := make([]byte, packetLength+4) // 包含 CRC 校验的长度
		n, err := io.ReadFull(transport.config.Port, data)
		if err != nil || n < int(packetLength+4) {
			// 读取不完整的包，缓存错误数据
			fmt.Println("Error: Incomplete packet received, buffering error data.")
			errorBuffer = append(errorBuffer, header.Type[:]...)
			errorBuffer = append(errorBuffer, header.Length[:]...)
			errorBuffer = append(errorBuffer, data[:n]...)
			continue
		}

		// 校验 CRC
		receivedData := data[:packetLength]


		// 如果没有错误，返回数据
		buffer = append(buffer, header.Type[:]...)
		buffer = append(buffer, header.Length[:]...)
		buffer = append(buffer, receivedData...)
		break
	}
	transport.config.Port.SetReadDeadline(time.Time{})
	transport.config.Logger.Debug("TransportLayer.Read=", ByteDumpHexString(buffer))
	return buffer, nil
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
