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

package internotify

import (
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/component/interdb"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/**
 *  开启定时清理全局通知
 *
 */
func StartClearInterNotifyCron() {
	for {
		select {
		case <-typex.GCTX.Done():
			return
		default:
		}
		if core.GlobalConfig.AppDebugMode {
			execInterNotifyCron(`-1 day`)
			time.Sleep(60 * time.Second) // For test
		} else {
			execInterNotifyCron(`-1 day`)
			time.Sleep(24 * time.Hour)
		}
	}
}
func execInterNotifyCron(period string) {
	deleteSql := fmt.Sprintf(`
DELETE FROM m_internal_notifies
WHERE created_at < date('now', '%s')
AND EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='m_internal_notifies');
`, period)
	ExecError := interdb.DB().Exec(deleteSql).Error
	if ExecError != nil {
		glogger.GLogger.Error(ExecError)
	}
	// glogger.GLogger.Debug("ExecDataCenterCron:", deleteSql)
}
