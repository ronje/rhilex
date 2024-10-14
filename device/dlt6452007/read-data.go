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
	"strconv"
)

// http://www.csmhkj.com/images/DLT645-2007%E9%80%9A%E4%BF%A1%E8%A7%84%E7%BA%A6.pdf
type DLT645Frame0x11 struct {
	Start      byte   // 起始符，固定为0x68
	Address    []byte // 地址域
	CtrlCode   byte   // 控制码
	DataLength byte   // 数据长度
	DataType   [4]byte
	DataArea   []byte // 数据域
	CheckSum   byte   // 校验和
	End        byte   // 结束符，固定为0x16
}

func (frame DLT645Frame0x11) String() string {
	var result string
	result += fmt.Sprintf("DLT645Frame0x11:\n=======\nStart: 0x%02x ", frame.Start)
	result += "\nAddress: "
	for _, b := range frame.Address {
		result += fmt.Sprintf("0x%02x ", b)
	}
	result += fmt.Sprintf("\nCtrlCode: 0x%02x ", frame.CtrlCode)
	result += fmt.Sprintf("\nDataLength: 0x%02x ", frame.DataLength)
	result += "\nDataType: "
	for _, b := range frame.DataType {
		result += fmt.Sprintf("0x%02x ", b)
	}
	result += "\nDataArea: "
	if len(frame.DataArea) == 0 {
		result += "[]"
	} else {
		for _, b := range frame.DataArea {
			result += fmt.Sprintf("0x%02x ", b)
		}
	}

	result += fmt.Sprintf("\nCheckSum: 0x%02x ", frame.CheckSum)
	result += fmt.Sprintf("\nEnd: 0x%02x\n=======\n", frame.End)
	return result
}

// Pack 打包DLT645协议帧
func (frame DLT645Frame0x11) Encode() ([]byte, error) {
	if len(frame.Address) != 6 {
		return nil, errors.New("address length must be 6 bytes")
	}
	if frame.DataLength != byte(len(frame.DataType)+len(frame.DataArea)) {
		return nil, errors.New("data length mismatch")
	}

	nFrame := new(bytes.Buffer)
	nFrame.WriteByte(frame.Start)
	nFrame.Write(frame.Address)
	nFrame.WriteByte(frame.Start)
	nFrame.WriteByte(frame.CtrlCode)
	nFrame.WriteByte(frame.DataLength)
	nFrame.Write(frame.DataType[:])
	nFrame.Write(frame.DataArea)
	frame.CheckSum = crc8(nFrame.Bytes())
	nFrame.WriteByte(frame.CheckSum)
	nFrame.WriteByte(frame.End)
	return nFrame.Bytes(), nil
}

/**
 * 获取数据类型
 *
 */
func (frame DLT645Frame0x11) GetDataType() (int64, error) {
	BCD := []byte{}
	for _, b := range frame.DataType {
		BCD = append(BCD, b-0x33)
	}
	for i, j := 0, len(BCD)-1; i < j; i, j = i+1, j-1 {
		BCD[i], BCD[j] = BCD[j], BCD[i]
	}
	str := ""
	for _, v := range BCD {
		str += fmt.Sprintf("%x", v)
	}
	return strconv.ParseInt(str, 10, len(BCD)*8)
}

/**
 * 解析数据
 *
 */
func (frame DLT645Frame0x11) GetData() (int64, error) {
	BCD := []byte{}
	for _, b := range frame.DataArea {
		BCD = append(BCD, b-0x33)
	}
	for i, j := 0, len(BCD)-1; i < j; i, j = i+1, j-1 {
		BCD[i], BCD[j] = BCD[j], BCD[i]
	}
	str := ""
	for _, v := range BCD {
		str += fmt.Sprintf("%x", v)
	}
	return strconv.ParseInt(str, 10, len(BCD)*8)
}
