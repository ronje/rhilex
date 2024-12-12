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

package alarmcenter

import (
	"fmt"
	"runtime"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/typex"

	"github.com/hootrhino/rhilex/glogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const __ALARM_DB_PATH string = "./rhilex_alarmcenter.db?cache=shared&mode=rwc"

var __AlarmSqlite *SqliteDAO

/*
*
* 初始化DAO
*
 */
func InitAlarmDb(engine typex.Rhilex) error {
	__AlarmSqlite = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.DebugMode {
		__AlarmSqlite.db, err = gorm.Open(sqlite.Open(__ALARM_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__AlarmSqlite.db, err = gorm.Open(sqlite.Open(__ALARM_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	__AlarmSqlite.db.Exec("VACUUM;")
	InitAlarmDbModel(__AlarmSqlite.db)
	return err
}

/*
*
* 停止
*
 */
func StopAlarmDb() {
	__AlarmSqlite.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func AlarmDb() *gorm.DB {
	return __AlarmSqlite.db
}

/*
*
* 注册数据模型
*
 */
func AlarmDbRegisterModel(dist ...interface{}) {
	__AlarmSqlite.db.AutoMigrate(dist...)
}
func InitAlarmDbModel(db *gorm.DB) {
	db.AutoMigrate(&MAlarmLog{})
	sql := `
CREATE TRIGGER IF NOT EXISTS limit_m_internal_notifies
AFTER INSERT ON m_alarm_logs
WHEN (SELECT COUNT(*) FROM m_alarm_logs) > %d
BEGIN
	DELETE FROM m_alarm_logs
	WHERE id IN (
		SELECT id FROM m_alarm_logs
		ORDER BY id ASC
		LIMIT (SELECT COUNT(*) - %d FROM m_alarm_logs)
	);
END;
`
	if errTrigger := AlarmDb().Exec(fmt.Sprintf(sql, 1000, 1000)).Error; errTrigger != nil {
		glogger.GLogger.Error(errTrigger)
	}
}
