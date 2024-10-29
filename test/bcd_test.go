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

package test

import "testing"

func DecimalToBCD(decimal int) []byte {
	bcd := []byte{}

	for decimal > 0 {
		digit := decimal % 10
		bcd = append([]byte{byte(digit)}, bcd...) // 将当前数字放到BCD数组前面
		decimal /= 10
	}

	return bcd
}

func BCDToDecimal(bcd []byte) int {
	decimal := 0
	multiplier := 1

	for i := len(bcd) - 1; i >= 0; i-- {
		decimal += int(bcd[i]) * multiplier
		multiplier *= 10
	}

	return decimal
}

// go test -timeout 30s -run ^TestBCDToDecimal github.com/hootrhino/rhilex/test -v -count=1
func TestBCDToDecimal(t *testing.T) {
	t.Log("BCDToDecimal == ", BCDToDecimal([]byte{0x22, 0x33}))
	t.Log("DecimalToBCD == ", DecimalToBCD(220))
}
