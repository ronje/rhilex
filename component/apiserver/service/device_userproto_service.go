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

func InsertUserProtocolPoints(list []model.MUserProtocolDataPoint) error {
	m := model.MUserProtocolDataPoint{}
	return interdb.DB().Model(m).Create(list).Error
}

func InsertUserProtocolPoint(P model.MUserProtocolDataPoint) error {
	IgnoreUUID := P.UUID
	Count := int64(0)
	P.UUID = ""
	interdb.DB().Model(P).Where(P).Count(&Count)
	if Count > 0 {
		return fmt.Errorf("already exists same record:%s", IgnoreUUID)
	}
	P.UUID = IgnoreUUID
	return interdb.DB().Model(P).Create(&P).Error
}

func DeleteUserProtocolPointByDevice(uuids []string, deviceUuid string) error {
	return interdb.DB().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MUserProtocolDataPoint{}).Error
}

func DeleteAllMUserProtocolByDevice(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MUserProtocolDataPoint{}).Error
}

func UpdateUserProtocolPoint(MUserProtocolDataPoint model.MUserProtocolDataPoint) error {
	return interdb.DB().Model(model.MUserProtocolDataPoint{}).
		Where("device_uuid=? AND uuid=?",
			MUserProtocolDataPoint.DeviceUuid, MUserProtocolDataPoint.UUID).
		Updates(MUserProtocolDataPoint).Error
}
