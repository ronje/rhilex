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
	"time"

	"github.com/hootrhino/rhilex/component/applet"
	"github.com/hootrhino/rhilex/typex"

	"github.com/hootrhino/rhilex/glogger"

	"testing"
)

// go test -timeout 30s -run ^Test_ModbusSlaverF5 github.com/hootrhino/rhilex/test -v -count=1

func Test_ModbusSlaverF5(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()
	Slaver := typex.NewDevice(typex.GENERIC_MODBUS_SLAVER,
		"GENERIC_MODBUS_SLAVER", "GENERIC_MODBUS_SLAVER",
		map[string]interface{}{
			"commonConfig": map[string]interface{}{
				"mode": "TCP",
			},
			"hostConfig": map[string]interface{}{
				"host": "127.0.0.1",
				"port": 1501,
			},
		})
	ctx, cancelF := typex.NewCCTX()
	Slaver.UUID = "JustForTest-UUID"
	if err := engine.LoadDeviceWithCtx(Slaver, ctx, cancelF); err != nil {
		t.Fatal(err)
	}
	app := applet.NewApplication(
		"JustForTest-UUID",
		"JustForTest-Name",
		"JustForTest-Version",
	)
	luaSource :=
		`
function Main(arg)
    local Value1 = modbus_slaver:F5("JustForTest-UUID", 1, 0)
    Debug("========= ModbusSlaver:F5: " .. Value1)
    local Value1 = modbus_slaver:F5("JustForTest-UUID", 1, 1)
    Debug("========= ModbusSlaver:F5: " .. Value1)
    local Value2 = modbus_slaver:F6("JustForTest-UUID", 1, 0xFFFF)
    Debug("========= ModbusSlaver:F6: " .. Value2)
    local Value2 = modbus_slaver:F6("JustForTest-UUID", 1, 0xAABB)
    Debug("========= ModbusSlaver:F6: " .. Value2)
    return 0
end

`
	if err := applet.LoadApp(app, luaSource); err != nil {
		glogger.GLogger.Fatal(err)
	}
	if err := applet.StartApp("JustForTest-UUID"); err != nil {
		glogger.GLogger.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	engine.Stop()
}
