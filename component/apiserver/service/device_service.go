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

func AllDevices() []model.MDevice {
	devices := []model.MDevice{}
	interdb.DB().Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func GetMDeviceWithUUID(uuid string) (*model.MDevice, error) {
	m := new(model.MDevice)
	return m, interdb.DB().Where("uuid=?", uuid).First(m).Error
}

// 检查名称是否重复
func CheckDeviceCount(T string) int64 {
	Count := int64(0)
	interdb.DB().Model(model.MDevice{}).Where("type=?", T).Count(&Count)
	return Count
}

// 检查名称是否重复
func CheckDeviceNameDuplicate(name string) bool {
	Count := int64(0)
	interdb.DB().Model(model.MDevice{}).Where("name=?", name).Count(&Count)
	return Count > 0
}

// 删除设备
func DeleteDevice(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MDevice{}).Error
}

// 创建设备
func InsertDevice(o *model.MDevice) error {
	return interdb.DB().Table("m_devices").Create(o).Error
}

// 更新设备信息
func UpdateDevice(uuid string, o *model.MDevice) error {
	m := model.MDevice{}
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
func FindDeviceByGroup(gid string) []model.MDevice {
	sql := `
WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN
		m_generic_group_relations ON (m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC;`

	m := []model.MDevice{}
	interdb.DB().Raw(`SELECT * FROM m_devices `+sql, gid).Find(&m)
	return m

}
func PageDevice(current, size int) (int64, []model.MDevice) {
	sql := `SELECT * FROM m_devices ORDER BY created_at DESC limit ? offset ?;`
	MDevices := []model.MDevice{}
	offset := (current - 1) * size
	interdb.DB().Raw(sql, size, offset).Find(&MDevices)
	var count int64
	interdb.DB().Model(&model.MDevice{}).Count(&count)
	return count, MDevices
}

/*
*
* 新增的分页获取
*
 */
func PageDeviceByGroup(current, size int, gid string) (int64, []model.MDevice) {
	sql := `
SELECT * FROM m_devices WHERE uuid IN (
	SELECT m_generic_group_relations.rid
	  FROM m_generic_groups
		LEFT JOIN m_generic_group_relations ON
		(m_generic_groups.uuid = m_generic_group_relations.gid)
	WHERE type = 'DEVICE' AND gid = ?
) ORDER BY created_at DESC limit ? offset ?;`
	MDevices := []model.MDevice{}
	offset := (current - 1) * size
	interdb.DB().Raw(sql, gid, size, offset).Find(&MDevices)
	var count int64
	countSql := `SELECT count(id)
FROM m_devices
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
	return count, MDevices
}
