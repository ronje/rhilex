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

package szy2062016

import (
	"bytes"
	"errors"
	"fmt"
)

// 当前阶段仅仅实现查询数据报文
// 传输方向位 DIR =0 DIV=0 FCB=FF FUNC
type SZY206Frame0x00 struct {
	Start1     byte    // 起始符，固定为0x68
	DataLength byte    // 数据长度
	Start2     byte    // 起始符，固定为0x68
	CtrlCode   byte    // 控制码：DIR[7] DIV [6] FCB[5 4] FUNC [3 2 1]
	Address    [5]byte // 地址域
	DataArea   []byte  // 数据域
	CheckSum   byte    // 校验和
	End        byte    // 结束符，固定为0x16
}

func (frame SZY206Frame0x00) String() string {
	var result string
	result += fmt.Sprintf("SZY206Frame0x00:\n=======\nStart: 0x%02x ", frame.Start1)
	result += "\nAddress: "
	for _, b := range frame.Address {
		result += fmt.Sprintf("0x%02x ", b)
	}
	result += fmt.Sprintf("\nCtrlCode: 0x%02x ", frame.CtrlCode)
	result += fmt.Sprintf("\nDataLength: 0x%02x ", frame.DataLength)
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

// Pack 打包SZY206协议帧
func (frame SZY206Frame0x00) Encode() ([]byte, error) {
	if len(frame.Address) != 5 {
		return nil, errors.New("address length must be 6 bytes")
	}
	nFrame := new(bytes.Buffer)
	nFrame.WriteByte(frame.Start1)
	nFrame.WriteByte(frame.DataLength)
	nFrame.WriteByte(frame.Start2)
	nFrame.WriteByte(frame.CtrlCode)
	nFrame.Write(frame.Address[:])
	nFrame.Write(frame.DataArea)
	frame.CheckSum = crc(nFrame.Bytes())
	nFrame.WriteByte(frame.CheckSum)
	nFrame.WriteByte(frame.End)
	return nFrame.Bytes(), nil
}

/**
 * 解析数据
 *
 */
func (frame SZY206Frame0x00) GetData() (int64, error) {
	return 0, nil
}
