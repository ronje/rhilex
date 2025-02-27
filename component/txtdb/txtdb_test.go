// Copyright (C) 2025 wwhai
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

package txtdb

import (
	"testing"
)

func Test_txtdb_test(t *testing.T) {
	// 创建一个新的文本数据库实例
	db := NewTextDB("test_txtdb_data.txt")

	// 添加数据
	err := db.Add("key1", "value1")
	if err != nil {
		t.Log("Error adding data:", err)
	}

	// 获取数据
	value, err := db.Get("key1")
	if err != nil {
		t.Log("Error getting data:", err)
	} else {
		t.Log("Value for key1:", value)
	}

	// 更新数据
	err = db.Update("key1", "newvalue1")
	if err != nil {
		t.Log("Error updating data:", err)
	}

	// 再次获取数据
	value, err = db.Get("key1")
	if err != nil {
		t.Log("Error getting data:", err)
	} else {
		t.Log("Updated value for key1:", value)
	}

	// 删除数据
	err = db.Delete("key1")
	if err != nil {
		t.Log("Error deleting data:", err)
	}

	// 尝试获取已删除的数据
	value, err = db.Get("key1")
	if err != nil {
		t.Log("Error getting data:", err)
	} else {
		t.Log("Value for key1:", value)
	}
}
