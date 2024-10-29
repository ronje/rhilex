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

package cjt1882004

import (
	"bytes"
	"fmt"
	"strconv"
)

// CJT188Frame represents the structure of a CJ-T188 protocol frame.
type CJT188Frame0x01 struct {
	Start        byte    // 帧起始符
	MeterType    byte    // 仪表类型
	Address      [7]byte // 地址域
	CtrlCode     byte    // 控制码
	DataLength   byte    // 数据长度域
	DataType     [2]byte // 数据长度域
	DataArea     []byte  // 数据域
	SerialNumber byte
	CheckSum     byte // 校验码
	End          byte // 结束符
}

func (frame CJT188Frame0x01) String() string {
	var result string
	result += fmt.Sprintf("CJT188Frame0x01:\n=======\nStart: 0x%02x ", frame.Start)
	result += fmt.Sprintf("\nMeterType: 0x%02x ", frame.MeterType)
	result += "\nAddress: "
	for _, b := range frame.Address {
		result += fmt.Sprintf("0x%02x ", b)
	}
	result += fmt.Sprintf("\nCtrlCode: 0x%02x ", frame.CtrlCode)
	result += fmt.Sprintf("\nDataLength: 0x%02x ", frame.DataLength)
	result += "\nDataType: "
	if len(frame.DataType) == 0 {
		result += "[]"
	} else {
		for _, b := range frame.DataType {
			result += fmt.Sprintf("0x%02x ", b)
		}
	}
	result += fmt.Sprintf("\nSerialNumber: 0x%02x ", frame.SerialNumber)
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

// Encode encodes the CJT188Frame into a byte slice.
func (frame CJT188Frame0x01) Encode() ([]byte, error) {
	nFrame := new(bytes.Buffer)
	nFrame.WriteByte(frame.Start)
	nFrame.WriteByte(frame.MeterType)
	nFrame.Write(frame.Address[:])
	nFrame.WriteByte(frame.CtrlCode)
	nFrame.WriteByte(frame.DataLength)
	nFrame.Write(frame.DataType[:])
	nFrame.Write(frame.DataArea[:])
	nFrame.WriteByte(frame.SerialNumber)
	frame.CheckSum = crc8(nFrame.Bytes())
	nFrame.WriteByte(frame.CheckSum)
	nFrame.WriteByte(frame.End)
	return nFrame.Bytes(), nil
}

/**
 * 解析数据
 *
 */
func (frame CJT188Frame0x01) GetData() (int64, error) {
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

// CJT188Frame0x01Response
type CJT188Frame0x01Response struct {
	Start        byte      // 帧起始符
	MeterType    byte      // 仪表类型
	Address      [7]byte   // 地址域
	CtrlCode     byte      // 控制码
	DataLength   byte      // 数据长度域
	DataType     [2]byte   // 数据长度域
	SerialNumber byte      // 序列号
	Data         [][5]byte // 数据域 [4 byte]+1unit
	Time         [7]byte   // 时间
	S0           byte      // 状态0
	S1           byte      // 状态1
	CheckSum     byte      // 校验码
	End          byte      // 结束符
}

func (frame CJT188Frame0x01Response) String() string {
	var result string
	result += fmt.Sprintf("CJT188Frame0x01:\n=======\nStart: 0x%02x ", frame.Start)
	result += fmt.Sprintf("\nMeterType: 0x%02x ", frame.MeterType)
	result += "\nAddress: "
	for _, b := range frame.Address {
		result += fmt.Sprintf("0x%02x ", b)
	}
	result += fmt.Sprintf("\nCtrlCode: 0x%02x ", frame.CtrlCode)
	result += fmt.Sprintf("\nDataLength: 0x%02x ", frame.DataLength)
	result += "\nDataType: "
	if len(frame.DataType) == 0 {
		result += "[]"
	} else {
		for _, b := range frame.DataType {
			result += fmt.Sprintf("0x%02x ", b)
		}
	}
	result += fmt.Sprintf("\nSerialNumber: 0x%02x ", frame.SerialNumber)
	result += "\nData: "
	if len(frame.Data) == 0 {
		result += "[]"
	} else {
		for i, row := range frame.Data {
			result += fmt.Sprintf("\n\tframe.Data[%d] ", i)
			for _, v := range row {
				result += fmt.Sprintf("0x%02x ", v)
			}
		}
	}
	result += "\nDataType: "
	if len(frame.Data) == 0 {
		result += "[]"
	} else {
		for _, b := range frame.DataType {
			result += fmt.Sprintf("0x%02x ", b)
		}
	}
	Value, _ := frame.GetData()
	result += fmt.Sprintf("\nData Value: %v ", Value)
	result += fmt.Sprintf("\nSN: 0x%02x ", frame.SerialNumber)
	result += fmt.Sprintf("\nTime: 0x%02x ", frame.Time)
	result += fmt.Sprintf("\nS0: 0x%02x ", frame.S0)
	result += fmt.Sprintf("\nS1: 0x%02x ", frame.S1)
	result += fmt.Sprintf("\nCheckSum: 0x%02x ", frame.CheckSum)
	result += fmt.Sprintf("\nEnd: 0x%02x\n=======\n", frame.End)
	return result
}

/**
 * 解析数据
 *
 */
func (frame CJT188Frame0x01Response) GetData() ([]int64, error) {
	result := []int64{}
	for _, row := range frame.Data {
		v0 := BCDBytesToDecimal(ByteReverse(row[:4])) // Unit
		result = append(result, v0)
	}
	return result, nil
}
