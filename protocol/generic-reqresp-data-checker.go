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

import "errors"

// 一个简单的实现 DataChecker 接口的示例（可以根据实际需要替换）
type SimpleChecker struct{}

func (c *SimpleChecker) CheckData(data []byte) error {

	if len(data) == 0 {
		return errors.New("data is empty")
	}
	return nil
}

// 实现CRC16 checker
type Crc16Checker struct {
}

func (c *Crc16Checker) CheckData(data []byte) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	return nil
}
