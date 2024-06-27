package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

func DeleteByUuids(deviceUuid string, uuids []string) error {
	return interdb.DB().
		Where("uuid IN ? AND device_uuid=?", uuids, deviceUuid).
		Delete(&model.MDataPoint{}).Error
}

func DeleteAllByDeviceUuid(deviceUuid string) error {
	return interdb.DB().
		Where("device_uuid=?", deviceUuid).
		Delete(&model.MDataPoint{}).Error
}
