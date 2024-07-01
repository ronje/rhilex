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
	return model.MDataPoint{}, nil
}

func (v ModbusValidator) ParseImportFile(file *excelize.File) ([]model.MDataPoint, error) {
	//TODO implement me
	return nil, nil
}

func (v ModbusValidator) Export(file *excelize.File, list []model.MDataPoint) error {
	//TODO implement me
	return nil
}
