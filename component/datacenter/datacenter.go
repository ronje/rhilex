// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package datacenter

import (
	"github.com/hootrhino/rhilex/typex"
	"gorm.io/gorm"
)

var __DefaultDataCenter *DataCenter

/*
*
* 留着未来扩充数据中心的功能
*
 */
type DataCenter struct {
	Sqlite *SqliteDb
	rhilex typex.Rhilex
}

func InitDataCenter(rhilex typex.Rhilex) {
	__DefaultDataCenter = new(DataCenter)
	__DefaultDataCenter.rhilex = rhilex
	__DefaultDataCenter.Sqlite = InitSqliteDb(rhilex)
	go StartDataCenterCron()
}

func DB() *gorm.DB {
	return __DefaultDataCenter.Sqlite.db
}
func VACUUM() {
	DB().Exec("VACUUM;")
}
