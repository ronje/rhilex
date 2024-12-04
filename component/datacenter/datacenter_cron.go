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

package datacenter

import (
	"fmt"
	"time"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

// period
// - `-n day`：表示当前日期减去n天
// - `-n month`：表示当前日期减去n个月
// - `-n year`：表示当前日期减去n年
// - `-n hour`：表示当前日期减去n小时
// - `-n minute`：表示当前日期减去n分钟
// - `-n second`：表示当前日期减去n秒
func StartClearDataCenterCron() {
	for {
		select {
		case <-typex.GCTX.Done():
			return
		default:
		}
		if core.GlobalConfig.DebugMode {
			execDataCenterCron(`-1 day`)
			time.Sleep(24 * time.Second) // For test
		} else {
			execDataCenterCron(`-1 day`)
			time.Sleep(24 * time.Hour)
		}
	}
}
func execDataCenterCron(period string) {
	sql := `SELECT name FROM sqlite_master WHERE type = 'table' AND name LIKE 'data_center_%';`
	tables := []string{}
	err := DB().Raw(sql).Scan(&tables).Error
	if err != nil {
		glogger.GLogger.Error(err)
		return
	}
	glogger.GLogger.Debug("ExecDataCenterCron:", sql)

	//
	// 删除period之前的数据
	// DELETE FROM %s WHERE create_at < date('now', '$s');
	//
	for _, table := range tables {
		deleteSql := fmt.Sprintf("DELETE FROM %s WHERE create_at < date('now', '%s');", table, period)
		ExecError := DB().Exec(deleteSql).Error
		if ExecError != nil {
			glogger.GLogger.Error(ExecError)
		}
		glogger.GLogger.Debug("ExecDataCenterCron:", deleteSql)
	}
	// DB().Exec("VACUUM;")
}
