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
	"encoding/binary"
	"errors"
	"io"

	"github.com/hootrhino/rhilex/typex"
)

type dlt645SerialTransporter struct {
	port typex.GenericRWC
}

func (dlt *dlt645SerialTransporter) SendFrame(aduRequest []byte) (aduResponse []byte, err error) {
	if _, err = dlt.port.Write(aduRequest); err != nil {
		return nil, err
	}
	return dlt.ReadFrame(dlt.port)
}

/**
 * 读请求
 *
 */
func (dlt *dlt645SerialTransporter) ReadFrame(rwc io.ReadWriteCloser) (aduResponse []byte, err error) {
	// 0x68
	var start1 byte
	for {
		var b byte
		if err = binary.Read(rwc, binary.BigEndian, &b); err != nil {
			return nil, err
		}
		if b != 0xFE {
			start1 = b
			break
		}

	}
	if start1 != 0x68 {
		return nil, errors.New("invalid start1 byte")
	}
	// 读取地址域 [6]byte
	var address [6]byte
	if err = binary.Read(rwc, binary.BigEndian, &address); err != nil {
		return nil, err
	}
	// 0x68
	var start2 byte
	if err = binary.Read(rwc, binary.BigEndian, &start2); err != nil {
		return nil, err
	}
	if start2 != 0x68 {
		return nil, errors.New("invalid start2 byte")
	}
	// 读取控制码
	var c byte
	if err = binary.Read(rwc, binary.BigEndian, &c); err != nil {
		return nil, err
	}

	// 读取数据域长度
	var l uint8
	if err = binary.Read(rwc, binary.BigEndian, &l); err != nil {
		return nil, err
	}

	// 读取数据域
	data := make([]byte, l)
	if _, err = io.ReadFull(rwc, data); err != nil {
		return nil, err
	}

	// 读取校验码
	var cs byte
	if err = binary.Read(rwc, binary.BigEndian, &cs); err != nil {
		return nil, err
	}
	// 读取帧结束符
	var end byte
	if err = binary.Read(rwc, binary.BigEndian, &end); err != nil {
		return nil, err
	}
	if end != 0x16 {
		return nil, errors.New("invalid end byte")
	}
	aduResponse = append(aduResponse, start1)
	aduResponse = append(aduResponse, address[:]...)
	aduResponse = append(aduResponse, start2)
	aduResponse = append(aduResponse, c)
	aduResponse = append(aduResponse, l)
	aduResponse = append(aduResponse, data...)
	aduResponse = append(aduResponse, cs)
	aduResponse = append(aduResponse, end)

	return aduResponse, nil
}
