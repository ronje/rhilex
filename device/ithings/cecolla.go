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

package ithings

import (
	"encoding/json"
	"errors"
	"strings"
)

type SubDeviceParam struct {
	Timestamp int64  `json:"timestamp"`
	ProductId string `json:"productID"`
	DeviceId  string `json:"deviceID"`
	Param     string `json:"param"`
	Value     any    `json:"value"`
}

func (O SubDeviceParam) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

// ParseProductInfo 尝试解析产品信息字符串。
// 它期望字符串格式为 "第一个字符串:第二个字符串"。
// 如果格式正确，返回两个字符串；如果格式不正确，返回错误。
func ParseProductInfo(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("invalid 'alias' filed: expected 'productId:deviceId'")
	}
	return parts[0], parts[1], nil
}
