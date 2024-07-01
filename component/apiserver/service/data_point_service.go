package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

func BatchDataPointCreate(list []model.MDataPoint) error {
	return interdb.DB().Create(list).Error
}

func BatchDataPointUpdate(list []model.MDataPoint) error {
	return interdb.DB().Updates(list).Error
}

func ListDataPointByUuid(deviceUuid string) ([]model.MDataPoint, error) {
	var records []model.MDataPoint
	tx := interdb.DB().Where("device_uuid=?", deviceUuid).Find(&records)
	return records, tx.Error
}

func BatchDeleteDataPointByUuids(deviceUuid string, uuids []string) error {
	return interdb.DB().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MDataPoint{}).Error
}

func BatchDeleteDataPointByDeviceUuid(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MDataPoint{}).Error
}
