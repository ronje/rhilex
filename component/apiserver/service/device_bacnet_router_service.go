package service

import (
	"fmt"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

func BatchInsertBacnetRouterPoint(list []model.MBacnetRouterDataPoint) error {
	m := model.MBacnetRouterDataPoint{}
	return interdb.InterDb().Model(m).Create(list).Error
}

func InsertBacnetRouterPoint(dataPoint model.MBacnetRouterDataPoint) error {
	IgnoreUUID := dataPoint.UUID
	Count := int64(0)
	dataPoint.UUID = ""
	interdb.InterDb().Model(dataPoint).Where(dataPoint).Count(&Count)
	if Count > 0 {
		return fmt.Errorf("already exists same record:%s", IgnoreUUID)
	}
	dataPoint.UUID = IgnoreUUID
	return interdb.InterDb().Model(dataPoint).Create(&dataPoint).Error
}

func UpdateBacnetRouterPoint(dataPoint model.MBacnetRouterDataPoint) error {
	return interdb.InterDb().Model(&model.MBacnetRouterDataPoint{}).
		Where("device_uuid=? AND uuid=?",
			dataPoint.DeviceUuid, dataPoint.UUID).
		Updates(dataPoint).Error
}

func BatchDeleteBacnetRouterPoint(uuids []string, deviceUuid string) error {
	return interdb.InterDb().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MBacnetRouterDataPoint{}).Error
}

func DeleteAllBacnetRouterPointByDeviceUuid(deviceUuid string) error {
	return interdb.InterDb().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MBacnetRouterDataPoint{}).Error
}
