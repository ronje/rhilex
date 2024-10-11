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

package lostcache

import (
	"runtime"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/typex"

	"github.com/hootrhino/rhilex/glogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const __DEFAULT_DB_PATH string = "./rhilex_lostcache.db?cache=shared&mode=rwc"

var __Sqlite *SqliteDAO

type SqliteDAO struct {
	engine typex.Rhilex
	name   string
	db     *gorm.DB
}

/*
*
* 初始化DAO
*
 */
func Init(engine typex.Rhilex) error {
	__Sqlite = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.AppDebugMode {
		__Sqlite.db, err = gorm.Open(sqlite.Open(__DEFAULT_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__Sqlite.db, err = gorm.Open(sqlite.Open(__DEFAULT_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}

	__Sqlite.db.Exec("VACUUM;")
	return err
}

/*
*
* 停止
*
 */
func Stop() {
	__Sqlite.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func DB() *gorm.DB {
	return __Sqlite.db
}

/*
*
* 返回名称
*
 */
func Name() string {
	return __Sqlite.name
}

/*
*
* 注册数据模型
*
 */
func RegisterModel(dist ...interface{}) {
	__Sqlite.db.AutoMigrate(dist...)
}
