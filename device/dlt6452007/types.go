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

package dlt6452007

import (
	"bytes"
	"errors"
	"fmt"
)

// DLT645Frame 表示一个DLT645协议帧
// http://www.csmhkj.com/images/DLT645-2007%E9%80%9A%E4%BF%A1%E8%A7%84%E7%BA%A6.pdf
type DLT645Frame struct {
	START      byte   // 起始符，固定为0x68
	ADDRESS    []byte // 地址域，通常是12个字节
	CtrlCode   byte   // 控制码
	DataLength byte   // 数据长度
	Data       []byte // 数据域
	CheckSum   byte   // 校验和
	END        byte   // 结束符，固定为0x16
}

// Pack 打包DLT645协议帧
func (f *DLT645Frame) Pack() ([]byte, error) {
	if len(f.ADDRESS) != 6 {
		return nil, errors.New("address length must be 6 bytes")
	}
	if f.DataLength != byte(len(f.Data)) {
		return nil, errors.New("data length mismatch")
	}

	frame := new(bytes.Buffer)
	frame.WriteByte(f.START)

	// 地址域，包括地址域标识和实际地址
	frame.WriteByte(0x01) // 地址域标识
	frame.Write(f.ADDRESS)

	frame.WriteByte(f.CtrlCode)
	frame.WriteByte(f.DataLength)
	frame.Write(f.Data)

	// 计算校验和
	var sum byte
	for _, b := range frame.Bytes()[1:] { // 跳过起始符
		sum += b
	}
	frame.WriteByte(sum)

	frame.WriteByte(f.END)

	return frame.Bytes(), nil
}

// Unpack 解包DLT645协议帧
func Unpack(data []byte) (*DLT645Frame, error) {
	if len(data) < 13 { // 至少包含起始符、地址域、控制码、数据长度、校验和和结束符
		return nil, errors.New("invalid frame length")
	}

	frame := &DLT645Frame{
		START:      data[0],
		ADDRESS:    data[2:8], // 跳过地址域标识
		CtrlCode:   data[8],
		DataLength: data[9],
	}

	if frame.DataLength > 0 {
		frame.Data = data[10 : 10+frame.DataLength]
	}

	frame.CheckSum = data[len(data)-2]
	frame.END = data[len(data)-1]

	// 校验校验和
	var sum byte
	for _, b := range data[1 : len(data)-2] { // 跳过起始符和结束符
		sum += b
	}
	if sum != frame.CheckSum {
		return nil, errors.New("checksum error")
	}

	return frame, nil
}

func TestCodec() {
	// 示例：构造一个DLT645请求帧
	address := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB}
	ctrlCode := byte(0x01)
	data := []byte{0x00, 0x00, 0x00, 0x00} // 示例数据

	frame := &DLT645Frame{
		START:      0x68,
		ADDRESS:    address,
		CtrlCode:   ctrlCode,
		DataLength: byte(len(data)),
		Data:       data,
		END:        0x16,
	}

	packedFrame, err := frame.Pack()
	if err != nil {
		fmt.Println("Pack error:", err)
		return
	}

	fmt.Printf("Packed frame: %x\n", packedFrame)

	// 示例：解析一个DLT645响应帧
	responseFrame, err := Unpack(packedFrame)
	if err != nil {
		fmt.Println("Unpack error:", err)
		return
	}

	fmt.Printf("Unpacked frame: %+v\n", responseFrame)
}
