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
	"fmt"
)

type PresentationLayer struct {
}

func (layer PresentationLayer) DecimalToBCD(decimal int) []byte {
	return DecimalToBCD(decimal)
}
func (layer PresentationLayer) BCDToDecimal(bcd []byte) int {
	return BCDToDecimal(bcd)
}

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

// Decimal: 12345 -> BCD: [1 2 3 4 5]
// BCD: [1 2 3 4 5] -> Decimal: 12345
func TestBCD() {
	decimal := 12345
	bcd := DecimalToBCD(decimal)
	fmt.Printf("Decimal: %d -> BCD: %v\n", decimal, bcd)

	decoded := BCDToDecimal(bcd)
	fmt.Printf("BCD: %v -> Decimal: %d\n", bcd, decoded)
}
