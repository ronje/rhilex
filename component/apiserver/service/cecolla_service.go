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

package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

func AllCecollas() []model.MCecolla {
	Cecollas := []model.MCecolla{}
	interdb.DB().Find(&Cecollas)
	return Cecollas
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func GetMCecollaWithUUID(uuid string) (*model.MCecolla, error) {
	m := new(model.MCecolla)
	return m, interdb.DB().Where("uuid=?", uuid).First(m).Error
}

// 检查名称是否重复
func CheckCecollaCount(T string) int64 {
	Count := int64(0)
	interdb.DB().Model(model.MCecolla{}).Where("type=?", T).Count(&Count)
	return Count
}

// 检查名称是否重复
func CheckCecollaNameDuplicate(name string) bool {
	Count := int64(0)
	interdb.DB().Model(model.MCecolla{}).Where("name=?", name).Count(&Count)
	return Count > 0
}

// 删除设备
func DeleteCecolla(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MCecolla{}).Error
}

// 创建设备
func InsertCecolla(o *model.MCecolla) error {
	return interdb.DB().Table("m_Cecollas").Create(o).Error
}

// 更新设备信息
func UpdateCecolla(uuid string, o *model.MCecolla) error {
	m := model.MCecolla{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*o)
		return nil
	}
}

/*
*
* 查询分组下的设备
*
 */
func FindCecollaByGroup(gid string) []model.MCecolla {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC;`

	m := []model.MCecolla{}
	interdb.DB().Raw(`SELECT * FROM m_cecollas `+sql, gid).Find(&m)
	return m

}
func PageCecolla(current, size int) (int64, []model.MCecolla) {
	sql := `SELECT * FROM m_cecollas ORDER BY created_at DESC limit ? offset ?;`
	MCecollas := []model.MCecolla{}
	offset := (current - 1) * size
	interdb.DB().Raw(sql, size, offset).Find(&MCecollas)
	var count int64
	interdb.DB().Model(&model.MCecolla{}).Count(&count)
	return count, MCecollas
}

/*
*
* 新增的分页获取
*
 */
func PageCecollaByGroup(current, size int, gid string) (int64, []model.MCecolla) {
	sql := `
SELECT * FROM m_cecollas WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN m_generic_group_relations ON
		(m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC limit ? offset ?;`
	MCecollas := []model.MCecolla{}
	offset := (current - 1) * size
	interdb.DB().Raw(sql, gid, size, offset).Find(&MCecollas)
	var count int64
	countSql := `SELECT count(id)
FROM m_cecollas
WHERE uuid IN (
SELECT m_generic_group_relations.rid
FROM m_generic_groups
LEFT JOIN
m_generic_group_relations
ON (m_generic_groups.uuid = m_generic_group_relations.gid)
WHERE type = 'DEVICE' AND
gid = ?
);
`
	interdb.DB().Raw(countSql, gid).Scan(&count)
	return count, MCecollas
}
