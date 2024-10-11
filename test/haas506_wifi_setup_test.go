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

import (
	"testing"
	"time"

	haas506 "github.com/hootrhino/rhilex/archsupport/haas506"
)

// go test -timeout 30s -run ^TestSetWifi github.com/hootrhino/rhilex/test -v -count=1

func TestSetWifi(t *testing.T) {
	// 示例用法
	iface := "wlan0"
	ssid := "YourWiFiSSID"
	psk := "YourWiFiPassword"
	timeout := 30 * time.Second

	err := haas506.SetWifi(iface, ssid, psk, timeout)
	if err != nil {
		t.Logf("Error setting up Wi-Fi: %v\n", err)
	} else {
		t.Log("Wi-Fi setup successful")
	}
}
