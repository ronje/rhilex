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

const __LOSTCACHE_DB_PATH string = "./rhilex_lostcache.db?cache=shared&mode=rwc"

var __LostCache *SqliteDAO

/*
*
* 初始化DAO
*
 */
func InitLostCacheDb(engine typex.Rhilex) error {
	__LostCache = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.DebugMode {
		__LostCache.db, err = gorm.Open(sqlite.Open(__LOSTCACHE_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__LostCache.db, err = gorm.Open(sqlite.Open(__LOSTCACHE_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}

	__LostCache.db.Exec("VACUUM;")
	return err
}

/*
*
* 停止
*
 */
func StopLostCacheDb() {
	__LostCache.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func LostCacheDb() *gorm.DB {
	return __LostCache.db
}

/*
*
* 注册数据模型
*
 */
func LostCacheDbRegisterModel(dist ...interface{}) {
	__LostCache.db.AutoMigrate(dist...)
}
