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
	"runtime"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/typex"

	"github.com/hootrhino/rhilex/glogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const __NOTIFY_DB_PATH string = "./rhilex_internotify.db?cache=shared&mode=rwc"

var __InterNotifySqlite *SqliteDAO

/*
*
* 初始化DAO
*
 */
func InitInterNotifyDb(engine typex.Rhilex) error {
	__InterNotifySqlite = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.DebugMode {
		__InterNotifySqlite.db, err = gorm.Open(sqlite.Open(__NOTIFY_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__InterNotifySqlite.db, err = gorm.Open(sqlite.Open(__NOTIFY_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	__InterNotifySqlite.db.Exec("VACUUM;")
	InitInterNotifyModel(__InterNotifySqlite.db)
	return err
}

/*
*
* 停止
*
 */
func StopInterNotify() {
	__InterNotifySqlite.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func InterNotifyDb() *gorm.DB {
	return __InterNotifySqlite.db
}

/*
*
* 注册数据模型
*
 */
func InterNotifyRegisterModel(dist ...interface{}) {
	__InterNotifySqlite.db.AutoMigrate(dist...)
}
func InitInterNotifyModel(db *gorm.DB) {
	db.AutoMigrate(&MInternalNotify{})
	sql := `
CREATE TRIGGER IF NOT EXISTS limit_m_internal_notifies
AFTER INSERT ON m_internal_notifies
WHEN ((SELECT COUNT(*) FROM m_internal_notifies) / 100) * 100 = (SELECT COUNT(*) FROM m_internal_notifies)
AND (SELECT COUNT(*) FROM m_internal_notifies) > 1000
BEGIN
    DELETE FROM m_internal_notifies
    WHERE id IN (
        SELECT id FROM m_internal_notifies
        ORDER BY id ASC
        LIMIT 100
    );
END;
`
	if errTrigger := InterNotifyDb().Exec(sql).Error; errTrigger != nil {
		glogger.GLogger.Error(errTrigger)
	}
}
