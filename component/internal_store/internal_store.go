// Copyright (C) 2025 wwhai
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

package internal_store

import (
	"sync"

	"github.com/hootrhino/rhilex/glogger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InternalStore 数据存储模块结构体
type InternalStore struct {
	db *gorm.DB
}

// 单例实例
var instance *InternalStore
var once sync.Once

func GetInstance(path string) *InternalStore {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})

		if err != nil {
			glogger.GLogger.Fatal(err)
		}
		instance = &InternalStore{db: db}
	})
	return instance
}

// GetByID 根据 ID 查询数据
func (is *InternalStore) GetByID(model interface{}, id string) error {
	return is.db.First(model, "id = ?", id).Error
}

// GetAll 查询所有数据
func (is *InternalStore) GetAll(model interface{}) error {
	return is.db.Find(model).Error
}

// Create 创建数据
func (is *InternalStore) Create(model interface{}) error {
	return is.db.Create(model).Error
}

// UpdateByID 根据 ID 更新数据
func (is *InternalStore) UpdateByID(model interface{}, id string) error {
	return is.db.Model(model).Where("id = ?", id).Updates(model).Error
}

// DeleteByID 根据 ID 删除数据
func (is *InternalStore) DeleteByID(model interface{}, id string) error {
	return is.db.Delete(model, "id = ?", id).Error
}

// Count 统计数据数量
func (is *InternalStore) Count(model interface{}) (int64, error) {
	var count int64
	err := is.db.Model(model).Count(&count).Error
	return count, err
}

// FindByCondition 根据条件查询数据
func (is *InternalStore) FindByCondition(model interface{}, condition string, values ...interface{}) error {
	return is.db.Where(condition, values...).Find(model).Error
}

// Migrate 迁移库表
func (is *InternalStore) Migrate(models ...interface{}) error {
	return is.db.AutoMigrate(models...)
}

// Paginate 分页查询数据
func (is *InternalStore) Paginate(model interface{}, page, pageSize int) error {
	offset := (page - 1) * pageSize
	return is.db.Offset(offset).Limit(pageSize).Find(model).Error
}
