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

package upgrader

import (
	"testing"
)

func Test_upgrader(t *testing.T) {
	// 测试 CompareVersion 函数
	t.Log(CompareVersion("v1.2.3", "v1.2.3")) // 应该返回 false，因为它们相等
	t.Log(CompareVersion("v2.2.3", "v1.2.3")) // 应该返回 true，因为第一个版本更新
	t.Log(CompareVersion("v1.3.3", "v1.2.3")) // 应该返回 true，因为第一个版本更新
	t.Log(CompareVersion("v1.2.4", "v1.2.3")) // 应该返回 true，因为第一个版本更新

}
