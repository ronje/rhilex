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

package service

import (
	"fmt"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/datacenter"
	"github.com/hootrhino/rhilex/component/interdb"
	"gorm.io/gorm"
)

// 获取DataSchema列表
func AllDataSchema() []model.MIotSchema {
	m := []model.MIotSchema{}
	interdb.DB().Model(model.MIotSchema{}).Find(&m)
	return m

}
func GetDataSchemaWithUUID(uuid string) (model.MIotSchema, error) {
	m := model.MIotSchema{}
	return m, interdb.DB().Model(model.MIotSchema{}).Where("uuid=?", uuid).First(&m).Error
}

/*
*
* 重置数据中心
*
 */
func ResetSchema(schemaUuid string) error {
	return interdb.DB().Model(model.MIotSchema{}).Transaction(func(tx *gorm.DB) error {
		MIotSchema := model.MIotSchema{}
		if err := tx.Where("uuid=?", schemaUuid).
			First(&MIotSchema).Error; err != nil {
			return err
		}
		if err := tx.Where("uuid=?", schemaUuid).
			Update("published", new(bool)).Error; err != nil {
			return err
		}
		return datacenter.DB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS data_center_%s;", schemaUuid)).Error
	})
}

// 删除DataSchema
func DeleteDataSchemaAndProperty(schemaUuid string) error {
	MIotSchema := model.MIotSchema{}
	if err := interdb.DB().Model(model.MIotSchema{}).Where("uuid=?", schemaUuid).
		First(&MIotSchema).Error; err != nil {
		return err
	}
	// 未发布的情况
	if !*MIotSchema.Published {
		// Only Delete Schema
		err2 := interdb.DB().Model(model.MIotSchema{}).Where("uuid=?", schemaUuid).Delete(&model.MIotSchema{}).Error
		if err2 != nil {
			return err2
		}
		if CountIotSchemaProperty(MIotSchema.Name, MIotSchema.UUID) > 0 {
			return fmt.Errorf("Schema Have Already Binding Properties")
		}
		return nil
	}
	// 已经发布了，清空RHILEX数据库
	return interdb.DB().Transaction(func(tx *gorm.DB) error {
		// Delete Schema
		err2 := tx.Model(model.MIotSchema{}).Where("uuid=?", schemaUuid).Delete(&model.MIotSchema{}).Error
		if err2 != nil {
			return err2
		}
		// Delete All IotProperty
		err1 := tx.Model(model.MIotProperty{}).Where("schema_id=?", schemaUuid).Delete(model.MIotProperty{}).Error
		if err1 != nil {
			return err1
		}
		// 清空数据中心的表
		err1Exec := datacenter.DB().Exec(fmt.Sprintf("DROP TABLE IF EXISTS data_center_%s;", schemaUuid)).Error
		if err1Exec != nil {
			return err1Exec
		}
		return nil
	})
}

// 创建DataSchema
func InsertDataSchema(DataSchema model.MIotSchema) error {
	return interdb.DB().Model(model.MIotSchema{}).Create(&DataSchema).Error
}

// 更新DataSchema
func UpdateDataSchema(DataSchema model.MIotSchema) error {
	return interdb.DB().
		Model(DataSchema).
		Where("uuid=?", DataSchema.UUID).
		Updates(&DataSchema).Error
}

// 更新DataSchema
func UpdateIotSchemaProperty(MIotProperty model.MIotProperty) error {
	return interdb.DB().
		Model(MIotProperty).
		Where("uuid=?", MIotProperty.UUID).
		Updates(&MIotProperty).Error
}

// 查找DataSchema
func FindIotSchemaProperty(uuid string) (model.MIotProperty, error) {
	MIotProperty := model.MIotProperty{}
	return MIotProperty,
		interdb.DB().
			Model(model.MIotProperty{}).
			Where("uuid=?", uuid).Find(&MIotProperty).Error
}

// 统计DataSchema
func CountIotSchemaProperty(name, schema_id string) int64 {
	var count int64
	interdb.DB().Model(model.MIotProperty{}).
		Where("name=? and schema_id=?", name, schema_id).Count(&count)
	return count
}

// 创建DataSchema
func InsertIotSchemaProperty(MIotProperty model.MIotProperty) error {
	return interdb.DB().
		Model(model.MIotProperty{}).Create(&MIotProperty).Error
}

// 删除
func DeleteIotSchemaProperty(uuid string) error {
	return interdb.DB().
		Model(model.MIotProperty{}).Where("uuid=?", uuid).Delete(model.MIotProperty{}).Error
}
