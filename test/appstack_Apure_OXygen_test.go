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

	"github.com/hootrhino/rhilex/component/appstack"

	"github.com/hootrhino/rhilex/glogger"

	"testing"
)

// go test -timeout 30s -run ^Test_Run_Apure_OXygen_parse github.com/hootrhino/rhilex/test -v -count=1

func Test_Run_Apure_OXygen_parse(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	app := appstack.NewApplication(
		"JustForTest-UUID",
		"JustForTest-Name",
		"JustForTest-Version",
	)
	luaSource :=
		`
function Main(arg)
    local Value = apure:ParseDOxygen("00010001")
    Debug("===== apure:ParseDOxygen(00010001): " .. Value)
    return 0
end

`
	if err := appstack.LoadApp(app, luaSource); err != nil {
		glogger.GLogger.Fatal(err)
	}
	if err := appstack.StartApp("JustForTest-UUID"); err != nil {
		glogger.GLogger.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	engine.Stop()
}
