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

import "fmt"

type DLT6452007DataLinkLayer struct {
}

func (dlt DLT6452007DataLinkLayer) CheckCrc(A []byte, B byte) error {
	crcCode := crc8(A)
	if crcCode != B {
		return fmt.Errorf("Crc8 (%v, CRC=%v) Check Failed: %v", A, crcCode, B)
	}
	return nil
}
