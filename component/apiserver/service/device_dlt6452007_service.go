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
	"github.com/hootrhino/rhilex/component/interdb"
)

func InsertDlt6452007Points(list []model.MDlt6452007DataPoint) error {
	m := model.MDlt6452007DataPoint{}
	return interdb.InterDb().Model(m).Create(list).Error
}

func InsertDlt6452007Point(P model.MDlt6452007DataPoint) error {
	IgnoreUUID := P.UUID
	Count := int64(0)
	P.UUID = ""
	interdb.InterDb().Model(P).Where(P).Count(&Count)
	if Count > 0 {
		return fmt.Errorf("already exists same record:%s", IgnoreUUID)
	}
	P.UUID = IgnoreUUID
	return interdb.InterDb().Model(P).Create(&P).Error
}

func DeleteDlt6452007PointByDevice(uuids []string, deviceUuid string) error {
	return interdb.InterDb().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MDlt6452007DataPoint{}).Error
}

func DeleteAllMDlt6452007ByDevice(deviceUuid string) error {
	return interdb.InterDb().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MDlt6452007DataPoint{}).Error
}

func UpdateDlt6452007Point(MDlt6452007DataPoint model.MDlt6452007DataPoint) error {
	return interdb.InterDb().Model(model.MDlt6452007DataPoint{}).
		Where("device_uuid=? AND uuid=?",
			MDlt6452007DataPoint.DeviceUuid, MDlt6452007DataPoint.UUID).
		Updates(MDlt6452007DataPoint).Error
}
