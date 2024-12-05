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

	"github.com/hootrhino/rhilex/component/alarmcenter"
)

// go test -timeout 30s -run ^Test_alarmcenter$ github.com/hootrhino/rhilex/component/alarmcenter -v -count=1
func Test_alarmcenter(t *testing.T) {
	// engine := RunTestEngine()
	// engine.Start()
	alarmcenter.InitAlarmCenter(nil)
	alarmcenter.LoadExpr("test", "temp > 10 && humi == 10 && oxy > 0")
	{
		R, err := alarmcenter.RunExpr("test", map[string]any{
			"temp": 12,
			"humi": 10,
			"oxy":  1,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(R)
	}

	alarmcenter.RemoveExpr("test")
	{
		R, err := alarmcenter.RunExpr("test", map[string]any{
			"age": 120,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(R)
	}
	alarmcenter.FlushAlarmCenter()
}
