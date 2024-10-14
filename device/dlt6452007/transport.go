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
)

type dlt645SerialTransporter struct {
	port io.ReadWriteCloser
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
	for {
		var b byte
		if err = binary.Read(rwc, binary.BigEndian, &b); err != nil {
			return nil, err
		}
		if b != 0xFE {
			break
		}
	}

	// 读取帧起始符
	var start byte
	if err = binary.Read(rwc, binary.BigEndian, &start); err != nil {
		return nil, err
	}
	if start != 0x68 {
		return nil, errors.New("invalid start byte")
	}

	// 读取地址域
	var address [6]byte
	if err = binary.Read(rwc, binary.BigEndian, &address); err != nil {
		return nil, err
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
	// crc:= dlt.DataLinkLayer
	// 读取帧结束符
	var end byte
	if err = binary.Read(rwc, binary.BigEndian, &end); err != nil {
		return nil, err
	}
	if end != 0x16 {
		return nil, errors.New("invalid end byte")
	}
	aduResponse = append(aduResponse, start)
	aduResponse = append(aduResponse, address[:]...)
	aduResponse = append(aduResponse, c)
	aduResponse = append(aduResponse, l)
	aduResponse = append(aduResponse, data...)
	aduResponse = append(aduResponse, cs)
	aduResponse = append(aduResponse, end)

	return aduResponse, nil
}
