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

import "testing"

// go test -timeout 30s -run ^TestFetchIthingsSchema$ github.com/hootrhino/rhilex/cecolla/ithings -v -count=1
func TestFetchIthingsSchema(t *testing.T) {
	// demo.ithings.net.cn
	// 01D
	// 基站1
	// 01D&基站1;12010126;Q1AMS;1889253787453
	// 4ae2b2175209fe7dc5fc1ef57468bc6a950046de;hmacsha1
	Resp, err := FetchIthingsSchema("demo.ithings.net.cn",
		"01D", "基站1", "01D&基站1;12010126;Q1AMS;1889253787453",
		"4ae2b2175209fe7dc5fc1ef57468bc6a950046de;hmacsha1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(Resp.String())
}
