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

func InsertSzy2062016Points(list []model.MSzy2062016DataPoint) error {
	m := model.MSzy2062016DataPoint{}
	return interdb.InterDb().Model(m).Create(list).Error
}

func InsertSzy2062016Point(P model.MSzy2062016DataPoint) error {
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

func DeleteSzy2062016PointByDevice(uuids []string, deviceUuid string) error {
	return interdb.InterDb().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MSzy2062016DataPoint{}).Error
}

func DeleteAllMSzy2062016ByDevice(deviceUuid string) error {
	return interdb.InterDb().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MSzy2062016DataPoint{}).Error
}

func UpdateSzy2062016Point(MSzy2062016DataPoint model.MSzy2062016DataPoint) error {
	return interdb.InterDb().Model(model.MSzy2062016DataPoint{}).
		Where("device_uuid=? AND uuid=?",
			MSzy2062016DataPoint.DeviceUuid, MSzy2062016DataPoint.UUID).
		Updates(MSzy2062016DataPoint).Error
}
