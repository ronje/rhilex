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
	"encoding/base64"
	"fmt"

	"time"
)

const (
	T_HmacSha256 = "hmacsha256"
	T_HmacSha1   = "hmacsha1"
)

// API
func GenSecretDeviceInfo(hmacType string, productID string, deviceName string, deviceSecret string) (
	clientID, userName, password string) {
	var (
		connID = Random(5, 1)
		expiry = time.Now().AddDate(10, 10, 10).Unix()
		token  string
		pwd, _ = base64.StdEncoding.DecodeString(deviceSecret)
	)
	clientID = productID + "&" + deviceName
	userName = fmt.Sprintf("%s;12010126;%s;%d", clientID, connID, expiry)
	if hmacType == T_HmacSha1 {
		token = HmacSha1(userName, pwd)
		password = token + ";hmacsha1"
	} else {
		token = HmacSha256(userName, pwd)
		password = token + ";hmacsha256"
	}
	return
}
