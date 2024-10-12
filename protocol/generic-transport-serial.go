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
	"fmt"
	"net"
	"time"
)

type SerialTransport struct {
	conn net.Conn
}

// Connect 连接到指定地址并设置连接超时
func (t *SerialTransport) Connect(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to device: %v", err)
	}
	t.conn = conn
	return nil
}

// SendRequest 发送请求数据，并设置写超时
func (t *SerialTransport) SendRequest(data []byte, timeout time.Duration) error {
	if t.conn == nil {
		return fmt.Errorf("no connection established")
	}
	t.conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err := t.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	return nil
}

// ReadResponse 读取带有固定包头和不定长数据的响应，支持设置读超时
func (t *SerialTransport) ReadResponse(headerLength int, timeout time.Duration) ([]byte, error) {
	if t.conn == nil {
		return nil, fmt.Errorf("no connection established")
	}

	// 1. 先读取固定长度的包头
	header := make([]byte, headerLength)
	t.conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := t.conn.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}
	if n != headerLength {
		return nil, fmt.Errorf("incomplete header read: expected %d bytes, got %d", headerLength, n)
	}

	// 2. 假设包头中最后两个字节表示数据长度，协议具体实现
	dataLength := int(header[headerLength-2])<<8 | int(header[headerLength-1])

	// 3. 根据解析出的数据长度读取剩余的数据
	data := make([]byte, dataLength)
	totalRead := 0
	for totalRead < dataLength {
		t.conn.SetReadDeadline(time.Now().Add(timeout))
		n, err = t.conn.Read(data[totalRead:])
		if err != nil {
			return nil, fmt.Errorf("failed to read data: %v", err)
		}
		totalRead += n
	}

	// 4. 返回完整的包头和数据
	return append(header, data...), nil
}

// Status 返回连接状态，如果连接正常则返回 nil
func (t *SerialTransport) Status() error {
	if t.conn == nil {
		return fmt.Errorf("no connection established")
	}
	return nil
}

// Close 关闭连接
func (t *SerialTransport) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
