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

	"github.com/hootrhino/rhilex/component/crontask"
)

var crontabExpressions = []string{
	"45 23 * * 6",                // 每周六晚上 23 点 45 分执行
	"5 0 * 8 *",                  // 每年 8 月的每天 0 点 5 分执行
	"0 12 15 5 *",                // 每年 5 月 15 日中午 12 点执行
}

// go test -timeout 30s -run ^TestParseCronTab github.com/hootrhino/rhilex/test -v -count=1
func TestParseCronTab(t *testing.T) {
	for _, expr := range crontabExpressions {
		parsed, err := crontask.ParseCronExpr(expr)
		if err != nil {
			t.Fatal("Error:", err)
		}
		t.Logf("Parsed Cron Expression: %+v\n", parsed)
	}

}
