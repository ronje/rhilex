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

	"github.com/hootrhino/rhilex/utils"
)

/*
*
* 检查报文的CRC值:[EE EF ········ CRC-H CRC-L \r\n]
*
 */
func CheckDataCrc16(buffer []byte) ([]byte, error) {
	Len := len(buffer)
	if Len < 4 {
		return nil, fmt.Errorf("Invalid packet:%v", buffer)
	}
	crcByte := [2]byte{buffer[Len-4], buffer[Len-3]}
	crcCheckedValue := uint16(crcByte[0])<<8 | uint16(crcByte[1])
	crcCalculatedValue := utils.CRC16(buffer[2 : Len-4])
	if crcCalculatedValue != crcCheckedValue {
		return nil, fmt.Errorf("CRC Check Error: (Checked=%d,Calculated=%d), data=%v",
			crcCheckedValue, crcCalculatedValue, buffer)
	}
	return buffer[2 : Len-4], nil
}
