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

package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/hootrhino/rhilex/glogger"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/dto"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/xuri/excelize/v2"
)

func InitModbusRoute() {
	// Modbus 点位表
	modbusMasterApi := server.RouteGroup(server.ContextUrl("/modbus_master_sheet"))
	{
		modbusMasterApi.POST(("/sheetImport"), server.AddRoute(ModbusMasterSheetImport))
		modbusMasterApi.GET(("/sheetExport"), server.AddRoute(ModbusMasterPointsExport))
		modbusMasterApi.GET(("/list"), server.AddRoute(ModbusMasterSheetPageList))
		modbusMasterApi.POST(("/update"), server.AddRoute(ModbusMasterSheetUpdate))
		modbusMasterApi.DELETE(("/delIds"), server.AddRoute(ModbusMasterSheetDelete))
		modbusMasterApi.DELETE(("/delAll"), server.AddRoute(ModbusMasterSheetDeleteAll))
		modbusMasterApi.POST(("/writeModbusSheet"), server.AddRoute(WriteModbusSheet))
	}
}

type ModbusMasterPointVo struct {
	UUID          string      `json:"uuid,omitempty"`
	DeviceUUID    string      `json:"device_uuid"`
	Tag           string      `json:"tag"`
	Alias         string      `json:"alias"`
	Function      *int        `json:"function"`
	SlaverId      *byte       `json:"slaverId"`
	Address       *uint16     `json:"address"`
	Frequency     *uint64     `json:"frequency"`
	Quantity      *uint16     `json:"quantity"`
	DataType      string      `json:"dataType"`      // 数据类型
	DataOrder     string      `json:"dataOrder"`     // 字节序
	Weight        *float64    `json:"weight"`        // 权重
	Status        int         `json:"status"`        // 运行时数据
	LastFetchTime uint64      `json:"lastFetchTime"` // 运行时数据
	Value         interface{} `json:"value"`         // 运行时数据
	ErrMsg        string      `json:"errMsg"`        // 运行时数据

}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/modbus_data_sheet/export
 */

// ModbusMasterPoints 获取modbus_excel类型的点位数据
func ModbusMasterPointsExport(c *gin.Context, ruleEngine typex.Rhilex) {
	deviceUuid, _ := c.GetQuery("device_uuid")

	var records []model.MModbusDataPoint
	result := interdb.DB().Table("m_modbus_data_points").
		Where("device_uuid=?", deviceUuid).Find(&records)
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{
		"tag", "alias",
		"function", "frequency",
		"slaverId", "address",
		"quality", "type",
		"order", "weight",
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			glogger.GLogger.Errorf("close excel file, err=%v", err)
		}
	}()
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	xlsx.SetSheetRow("Sheet1", cell, &Headers)
	for idx, record := range records[0:] {
		Row := []string{
			record.Tag,
			record.Alias,
			fmt.Sprintf("%d", *record.Function),
			fmt.Sprintf("%d", *record.Frequency),
			fmt.Sprintf("%d", *record.SlaverId),
			fmt.Sprintf("%d", *record.Address),
			fmt.Sprintf("%d", *record.Quantity),
			record.DataType,
			record.DataOrder,
			fmt.Sprintf("%f", *record.Weight),
		}
		cell, _ = excelize.CoordinatesToCellName(1, idx+2)
		xlsx.SetSheetRow("Sheet1", cell, &Row)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%v.xlsx",
		time.Now().UnixMilli()))
	xlsx.WriteTo(c.Writer)
}

// 分页获取
// SELECT * FROM `m_modbus_data_points` WHERE
// `m_modbus_data_points`.`device_uuid` = "DEVICEDQNLO8"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func ModbusMasterSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MModbusDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MModbusDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MModbusDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []ModbusMasterPointVo{}

	for _, record := range records {
		Slot := intercache.GetSlot(deviceUuid)
		value, ok := Slot[record.UUID]
		Vo := ModbusMasterPointVo{
			UUID:          record.UUID,
			DeviceUUID:    record.DeviceUuid,
			Tag:           record.Tag,
			Alias:         record.Alias,
			Function:      record.Function,
			SlaverId:      record.SlaverId,
			Address:       record.Address,
			Frequency:     record.Frequency,
			Quantity:      record.Quantity,
			DataType:      record.DataType,
			DataOrder:     record.DataOrder,
			Weight:        record.Weight,
			LastFetchTime: value.LastFetchTime, // 运行时
			Value:         value.Value,         // 运行时
			ErrMsg:        value.ErrMsg,        // 运行时
		}
		if ok {
			Vo.Status = func() int {
				if value.Value == "" {
					return 0
				}
				return 1
			}() // 运行时
			Vo.LastFetchTime = value.LastFetchTime // 运行时
			types, _ := utils.IsArrayAndGetValueList(value.Value)
			Vo.Value = types
			recordsVo = append(recordsVo, Vo)
		} else {
			recordsVo = append(recordsVo, Vo)
		}
	}
	Result := service.WrapPageResult(*pager, recordsVo, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 删除单行
*
 */
func ModbusMasterSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllModbusPointByDevice(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
*删除
*
 */
func ModbusMasterSheetDelete(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteModbusPointByDevice(form.UUIDs, form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 检查点位合法性
*
 */

func CheckModbusMasterDataPoints(M ModbusMasterPointVo) error {
	// Helper function to check string length
	checkStringLength := func(value string, paramName string, maxLength int) error {
		if value == "" {
			return fmt.Errorf("missing required param '%s'", paramName)
		}
		if len(value) > maxLength {
			return fmt.Errorf("'%s' length must be in the range of 1-%d", paramName, maxLength)
		}
		return nil
	}

	// Check required string fields
	if err := checkStringLength(M.Tag, "tag", 64); err != nil {
		return err
	}
	if err := checkStringLength(M.Alias, "alias", 64); err != nil {
		return err
	}

	// Check Address
	if M.Address == nil {
		return fmt.Errorf("missing required param 'address'")
	}
	if *M.Address > uint16(65535) {
		return fmt.Errorf("'address' must be in the range of 0-65535")
	}

	// Check Function
	if M.Function == nil {
		return fmt.Errorf("missing required param 'function'")
	}
	if *M.Function < int(1) || *M.Function > int(4) {
		return fmt.Errorf("'function' only supports values of 1, 2, 3, or 4")
	}

	// Check SlaverId
	if M.SlaverId == nil {
		return fmt.Errorf("missing required param 'slaverId'")
	}
	if *M.SlaverId > uint8(255) {
		return fmt.Errorf("'slaverId' must be in the range of 0-255")
	}

	// Check Frequency
	if M.Frequency == nil {
		return fmt.Errorf("missing required param 'frequency'")
	}
	if *M.Frequency < uint64(1) {
		return fmt.Errorf("'frequency' must be greater than 50ms")
	}
	if *M.Frequency > uint64(100000) {
		return fmt.Errorf("'frequency' must be less than 100s")
	}

	// Check Quantity
	if M.Quantity == nil {
		return fmt.Errorf("missing required param 'quantity'")
	}

	// Validate DataOrder for different DataTypes
	dataOrderMap := map[string][]string{
		"UTF8":     {"BIG_ENDIAN", "LITTLE_ENDIAN"},
		"I":        {"A"},
		"Q":        {"A"},
		"BYTE":     {"A"},
		"BOOL":     {"A"},
		"INT16":    {"AB", "BA"},
		"UINT16":   {"AB", "BA"},
		"RAW":      {"ABCD", "DCBA", "CDAB"},
		"INT":      {"ABCD", "DCBA", "CDAB"},
		"INT32":    {"ABCD", "DCBA", "CDAB"},
		"UINT":     {"ABCD", "DCBA", "CDAB"},
		"UINT32":   {"ABCD", "DCBA", "CDAB"},
		"FLOAT":    {"ABCD", "DCBA", "CDAB"},
		"FLOAT32":  {"ABCD", "DCBA", "CDAB"},
		"UFLOAT32": {"ABCD", "DCBA", "CDAB"},
	}

	validOrders, exists := dataOrderMap[M.DataType]
	if !exists {
		return fmt.Errorf("invalid data type '%s'", M.DataType)
	}
	if !utils.SContains(validOrders, M.DataOrder) {
		return fmt.Errorf("invalid '%s' order '%s'", M.DataType, M.DataOrder)
	}

	// Check Weight
	if M.Weight == nil {
		return fmt.Errorf("invalid weight value: %v", M.Weight)
	}

	// Validate Tag Name
	if !utils.IsValidColumnName(M.Tag) {
		return fmt.Errorf("invalid tag name: %s", M.Tag)
	}

	return nil
}

/*
*
* 更新点位表
*
 */
func ModbusMasterSheetUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID             string                `json:"device_uuid"`
		ModbusMasterDataPoints []ModbusMasterPointVo `json:"data_points"`
	}
	//  ModbusMasterDataPoints := [] ModbusMasterPointVo{}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, ModbusMasterDataPoint := range form.ModbusMasterDataPoints {
		if err := CheckModbusMasterDataPoints(ModbusMasterDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if ModbusMasterDataPoint.UUID == "" ||
			ModbusMasterDataPoint.UUID == "new" ||
			ModbusMasterDataPoint.UUID == "copy" {
			NewRow := model.MModbusDataPoint{
				UUID:       utils.ModbusPointUUID(),
				DeviceUuid: ModbusMasterDataPoint.DeviceUUID,
				Tag:        ModbusMasterDataPoint.Tag,
				Alias:      ModbusMasterDataPoint.Alias,
				Function:   ModbusMasterDataPoint.Function,
				SlaverId:   ModbusMasterDataPoint.SlaverId,
				Address:    ModbusMasterDataPoint.Address,
				Frequency:  ModbusMasterDataPoint.Frequency,
				Quantity:   ModbusMasterDataPoint.Quantity,
				DataType:   ModbusMasterDataPoint.DataType,
				DataOrder:  ModbusMasterDataPoint.DataOrder,
				Weight:     utils.HandleZeroValue(ModbusMasterDataPoint.Weight),
			}
			err0 := service.InsertModbusPointPosition(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MModbusDataPoint{
				UUID:       ModbusMasterDataPoint.UUID,
				DeviceUuid: ModbusMasterDataPoint.DeviceUUID,
				Tag:        ModbusMasterDataPoint.Tag,
				Alias:      ModbusMasterDataPoint.Alias,
				Function:   ModbusMasterDataPoint.Function,
				SlaverId:   ModbusMasterDataPoint.SlaverId,
				Address:    ModbusMasterDataPoint.Address,
				Frequency:  ModbusMasterDataPoint.Frequency,
				Quantity:   ModbusMasterDataPoint.Quantity,
				DataType:   ModbusMasterDataPoint.DataType,
				DataOrder:  ModbusMasterDataPoint.DataOrder,
				Weight:     utils.HandleZeroValue(ModbusMasterDataPoint.Weight),
			}
			err0 := service.UpdateModbusPoint(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

// ModbusMasterSheetImport 上传Excel文件
func ModbusMasterSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
	// 解析 multipart/form-data 类型的请求体
	err := c.Request.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	defer file.Close()
	deviceUuid := c.Request.Form.Get("device_uuid")

	Device := dto.RhilexDeviceDto{}
	errDb := interdb.DB().Table("m_devices").
		Where("uuid=?", deviceUuid).Find(&Device).Error
	if errDb != nil {
		c.JSON(common.HTTP_OK, common.Error400(errDb))
		return
	}
	if Device.Type == "" {
		c.JSON(common.HTTP_OK,
			common.Error("Device Not Exists"))
		return
	}
	if Device.Type != typex.GENERIC_MODBUS_MASTER.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import ModbusMaster Device"))
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		contentType != "application/vnd.ms-excel" {
		c.JSON(common.HTTP_OK, common.Error("File Must be Excel Sheet"))
		return
	}
	// 判断文件大小是否符合要求（10MB）
	if header.Size > 1024*1024*10 {
		c.JSON(common.HTTP_OK, common.Error("Excel file size cannot be greater than 10MB"))
		return
	}
	// 只取第一张表，而且名字必须是Sheet1
	list, err := parseModbusMasterPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertModbusPointPositions(list); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(deviceUuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 解析表格
*
 */

func parseModbusMasterPointExcel(r io.Reader, sheetName string,
	deviceUuid string) (list []model.MModbusDataPoint, err error) {
	excelFile, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()
	// 读取表格
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	// 判断首行标头
	// tag, alias, function, frequency, slaverId, address, quality
	err1 := errors.New(" Invalid Sheet Header")
	if len(rows[0]) < 10 {
		return nil, err1
	}

	// 严格检查表结构
	if rows[0][0] != "tag" ||
		rows[0][1] != "alias" ||
		rows[0][2] != "function" ||
		rows[0][3] != "frequency" ||
		rows[0][4] != "slaverId" ||
		rows[0][5] != "address" ||
		rows[0][6] != "quality" ||
		rows[0][7] != "type" ||
		rows[0][8] != "order" ||
		rows[0][9] != "weight" {
		return nil, err1
	}

	list = make([]model.MModbusDataPoint, 0)
	// tag, alias, function, frequency, slaverId, address, quality
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		tag := row[0]
		alias := row[1]
		function, _ := strconv.ParseUint(row[2], 10, 8)
		frequency, _ := strconv.ParseUint(row[3], 10, 64)
		slaverId, _ := strconv.ParseUint(row[4], 10, 8)
		address, _ := strconv.ParseUint(row[5], 10, 16)
		quantity, _ := strconv.ParseUint(row[6], 10, 16)
		Type := row[7]
		Order := row[8]
		Weight, _ := strconv.ParseFloat(row[9], 32)
		if Weight == 0 {
			Weight = 1 // 防止解析异常的时候系数0
		}
		Function := int(function)
		SlaverId := byte(slaverId)
		Address := uint16(address)
		Frequency := uint64(frequency)
		Quantity := uint16(quantity)

		if err := CheckModbusMasterDataPoints(ModbusMasterPointVo{
			Tag:       tag,
			Alias:     alias,
			Function:  &Function,
			SlaverId:  &SlaverId,
			Address:   &Address,
			Frequency: &Frequency, //ms
			Quantity:  &Quantity,
			DataType:  Type,
			DataOrder: utils.GetDefaultDataOrder(Type, Order),
			Weight:    &Weight,
		}); err != nil {
			return nil, err
		}
		//
		model := model.MModbusDataPoint{
			UUID:       utils.ModbusPointUUID(),
			DeviceUuid: deviceUuid,
			Tag:        tag,
			Alias:      alias,
			Function:   &Function,
			SlaverId:   &SlaverId,
			Address:    &Address,
			Frequency:  &Frequency, //ms
			Quantity:   &Quantity,
			DataType:   Type,
			DataOrder:  utils.GetDefaultDataOrder(Type, Order),
			Weight:     &Weight,
		}
		list = append(list, model)
	}
	return list, nil
}

/**
 * 给某个云边发指令
 *
 */
// POST -> temp , 0x0001
type CtrlCmd struct {
	UUID    string `json:"uuid"`    // 设备UUID
	PointId string `json:"pointId"` // 点位Point Id
	Tag     string `json:"tag"`     // 点位表的Tag
	Value   string `json:"value"`   // 写的值
}

func (O CtrlCmd) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}
func WriteModbusSheet(c *gin.Context, ruleEngine typex.Rhilex) {
	form := CtrlCmd{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	device := ruleEngine.GetDevice(form.UUID)
	if device != nil {
		_, errOnCtrl := device.Device.OnCtrl([]byte("WriteToSheetRegister"), []byte(form.String()))
		if errOnCtrl != nil {
			c.JSON(common.HTTP_OK, common.Error400(errOnCtrl))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("Device not exists:"+form.UUID))
}
