// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package modbusscanner

import (
	"encoding/binary"

	"github.com/hootrhino/rhilex/typex"
)

/*
*
* CRC 计算
*
 */

func calculateCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF

	for _, b := range data {
		crc ^= uint16(b)

		for i := 0; i < 8; i++ {
			lsb := crc & 0x0001
			crc >>= 1

			if lsb == 1 {
				crc ^= 0xA001 // 0xA001 是Modbus CRC16多项式的表示
			}
		}
	}

	return crc
}
func uint16ToBytes(val uint16) [2]byte {
	bytes := [2]byte{}
	binary.LittleEndian.PutUint16(bytes[:], val)
	return bytes
}

/*
*
* 服务调用接口
*
 */
func (ms *modbusScanner) Service(arg typex.ServiceArg) typex.ServiceResult {
	if ms.busying {
		if arg.Name == "stop" {
			if ms.cancel != nil {
				ms.cancel()
				ms.busying = false
				return typex.ServiceResult{Out: []map[string]interface{}{
					{"error": "Device busying now"},
				}}
			}
		}
		return typex.ServiceResult{Out: []map[string]interface{}{
			{"error": "Modbus Scanner Busing now"},
		}}
	}

	if arg.Name == "scan" {
		ms.busying = true
	}
	return typex.ServiceResult{Out: []map[string]interface{}{}}
}
