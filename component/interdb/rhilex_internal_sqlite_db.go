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

package interdb

import (
	"runtime"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/typex"

	"github.com/hootrhino/rhilex/glogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const __INTERNAL_DB_PATH string = "./rhilex.db?cache=shared&mode=rwc"

var __InternalSqlite *SqliteDAO

/*
*
* 初始化DAO
*
 */
func InitInterDb(engine typex.Rhilex) error {
	__InternalSqlite = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.DebugMode {
		__InternalSqlite.db, err = gorm.Open(sqlite.Open(__INTERNAL_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__InternalSqlite.db, err = gorm.Open(sqlite.Open(__INTERNAL_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	__InternalSqlite.db.Exec("VACUUM;")
	return err
}

/*
*
* 停止
*
 */
func StopInterDb() {
	__InternalSqlite.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func InterDb() *gorm.DB {
	return __InternalSqlite.db
}

/*
*
* 注册数据模型
*
 */
func InterDbRegisterModel(dist ...interface{}) {
	__InternalSqlite.db.AutoMigrate(dist...)
}
