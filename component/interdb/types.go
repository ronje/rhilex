package interdb

import (
	"github.com/hootrhino/rhilex/typex"
	"gorm.io/gorm"
)

/*
*
* Sqlite 数据持久层
*
 */
type SqliteDAO struct {
	engine typex.Rhilex
	name   string   // 框架可以根据名称来选择不同的数据库驱动,为以后扩展准备
	db     *gorm.DB // Sqlite 驱动
}
