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
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
)

type Header struct {
	Id     [2]byte
	Type   [2]byte
	Length [2]byte
}

type AppLayerFrame struct {
	Delimiter        [2]byte
	Header           Header
	Payload          []byte
	CrcSum           [4]byte
	ReverseDelimiter [2]byte
}

func (f AppLayerFrame) String() string {
	bytes, _ := json.Marshal(f)
	return string(bytes)
}

// Encode 将 Header 编码为字节数组
func (h *Header) Encode() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.BigEndian, h.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Id: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, h.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Type: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, h.Length)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Length: %w", err)
	}

	return buf.Bytes(), nil
}

// DecodeHeader 从字节数组解码出 Header
func DecodeHeader(data []byte) (Header, error) {
	if len(data) < 6 {
		return Header{}, errors.New("data too short to decode header")
	}

	var header Header
	buf := bytes.NewReader(data[:6])

	err := binary.Read(buf, binary.BigEndian, &header.Id)
	if err != nil {
		return Header{}, fmt.Errorf("failed to decode Id: %w", err)
	}

	err = binary.Read(buf, binary.BigEndian, &header.Type)
	if err != nil {
		return Header{}, fmt.Errorf("failed to decode Type: %w", err)
	}

	err = binary.Read(buf, binary.BigEndian, &header.Length)
	if err != nil {
		return Header{}, fmt.Errorf("failed to decode Length: %w", err)
	}

	return header, nil
}

// Encode 将 AppLayerFrame 编码为字节数组
func (f *AppLayerFrame) Encode() ([]byte, error) {
	var buf bytes.Buffer

	// 写入 Delimiter
	err := binary.Write(&buf, binary.BigEndian, f.Delimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Delimiter: %w", err)
	}

	// 写入 Header
	headerBytes, err := f.Header.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode Header: %w", err)
	}
	buf.Write(headerBytes)

	// 写入 Payload
	buf.Write(f.Payload)

	// 计算并写入 CRC 校验和
	crc := crc32.ChecksumIEEE(f.Payload)
	binary.BigEndian.PutUint32(f.CrcSum[:], crc)
	buf.Write(f.CrcSum[:])

	// 写入 ReverseDelimiter
	err = binary.Write(&buf, binary.BigEndian, f.ReverseDelimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ReverseDelimiter: %w", err)
	}

	return buf.Bytes(), nil
}

// DecodeAppLayerFrame 从字节数组解码出 AppLayerFrame
func DecodeAppLayerFrame(data []byte) (AppLayerFrame, error) {
	if len(data) < 16 {
		return AppLayerFrame{}, errors.New("data too short to decode AppLayerFrame")
	}

	frame := AppLayerFrame{}
	buf := bytes.NewReader(data)

	// 读取 Delimiter
	err := binary.Read(buf, binary.BigEndian, &frame.Delimiter)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to decode Delimiter: %w", err)
	}

	// 读取 Header
	headerBytes := make([]byte, 6)
	_, err = buf.Read(headerBytes)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to read Header bytes: %w", err)
	}
	frame.Header, err = DecodeHeader(headerBytes)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to decode Header: %w", err)
	}

	// 获取 Payload 长度并读取 Payload
	payloadLength := binary.BigEndian.Uint16(frame.Header.Length[:])
	frame.Payload = make([]byte, payloadLength)
	_, err = buf.Read(frame.Payload)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to read Payload: %w", err)
	}

	// 读取 CRC 校验和
	err = binary.Read(buf, binary.BigEndian, &frame.CrcSum)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to decode CRC: %w", err)
	}

	// 校验 CRC
	calculatedCrc := crc32.ChecksumIEEE(frame.Payload)
	if calculatedCrc != binary.BigEndian.Uint32(frame.CrcSum[:]) {
		return AppLayerFrame{}, errors.New("CRC check failed")
	}

	// 读取 ReverseDelimiter
	err = binary.Read(buf, binary.BigEndian, &frame.ReverseDelimiter)
	if err != nil {
		return AppLayerFrame{}, fmt.Errorf("failed to decode ReverseDelimiter: %w", err)
	}

	return frame, nil
}
