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

// go test -timeout 30s -run ^Test_alarm_Center_Normal$ github.com/hootrhino/rhilex/test -v -count=1
func Test_alarm_Center_Normal(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	alarmcenter.InitAlarmCenter(engine)
	alarmcenter.LoadAlarmRule("test", alarmcenter.AlarmRule{})
	{
		for i := 0; i < 20; i++ {
			R, err := alarmcenter.RunExpr("test", "test", map[string]any{
				"temp": 12,
				"humi": 10,
				"oxy":  1,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Log(R)
		}
	}

	alarmcenter.RemoveExpr("test")
	{
		R, err := alarmcenter.RunExpr("test", "test", map[string]any{
			"age": 120,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(R)
	}
	alarmcenter.StopAlarmCenter()
}

// go test -timeout 30s -run ^Test_alarm_Center_Not_Effect$ github.com/hootrhino/rhilex/test -v -count=1

func Test_alarm_Center_Not_Effect(t *testing.T) {
	alarmcenter.InitAlarmCenter(nil)
	alarmcenter.LoadAlarmRule("test", alarmcenter.AlarmRule{})
	{
		for i := 0; i < 2000; i++ {
			R, err := alarmcenter.RunExpr("test", "test", map[string]any{
				"temp": 12,
				"humi": 10,
				"oxy":  1,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Log(R)
		}
	}

	alarmcenter.RemoveExpr("test")
	{
		R, err := alarmcenter.RunExpr("test", "test", map[string]any{
			"age": 120,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(R)
	}
	alarmcenter.StopAlarmCenter()
}
