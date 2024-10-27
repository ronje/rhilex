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

	"github.com/hootrhino/rhilex/component/orderedmap"
)

// go test -timeout 30s -run ^TestOrderedMap github.com/hootrhino/rhilex/test -v -count=1
func TestOrderedMap(t *testing.T) {
	om := orderedmap.NewOrderedMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	for _, k := range om.Keys() {
		v, _ := om.Get(k)
		t.Log(k, v)
	}
	t.Log("Keys:", om.Keys())
	t.Log("Values:", om.Values())

	om.Delete("b")

	t.Log("After deleting 'b':")
	t.Log("Keys:", om.Keys())
	t.Log("Values:", om.Values())

	value, exists := om.Get("a")
	if exists {
		t.Log("Value of 'a':", value)
	}
}
