package modbus

import (
	"github.com/hootrhino/rhilex/component/apiserver/dto"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/xuri/excelize/v2"
)

type ModbusValidator struct {
}

func (v ModbusValidator) Validate(dto dto.DataPointCreateOrUpdateDTO) (model.MDataPoint, error) {
	//TODO implement me
	panic("implement me")
}

func (v ModbusValidator) Import(file *excelize.File) ([]model.MDataPoint, error) {
	//TODO implement me
	panic("implement me")
}

func (v ModbusValidator) Export(list []model.MDataPoint) (*excelize.File, error) {
	//TODO implement me
	panic("implement me")
}
