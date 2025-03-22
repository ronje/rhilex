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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/glogger"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/jinzhu/copier"
	"github.com/xuri/excelize/v2"
)

func InitSiemensS7Route() {
	SIEMENS_PLC := server.RouteGroup(server.ContextUrl("/s1200_data_sheet"))
	{
		SIEMENS_PLC.POST(("/sheetImport"), server.AddRoute(SiemensSheetImport))
		SIEMENS_PLC.GET(("/sheetExport"), server.AddRoute(SiemensPointsExport))
		SIEMENS_PLC.GET(("/list"), server.AddRoute(SiemensSheetPageList))
		SIEMENS_PLC.POST(("/update"), server.AddRoute(SiemensSheetUpdate))
		SIEMENS_PLC.DELETE(("/delIds"), server.AddRoute(SiemensSheetDelete))
		SIEMENS_PLC.DELETE(("/delAll"), server.AddRoute(SiemensSheetDeleteAll))
	}
}

type SiemensPointVo struct {
	UUID           string   `json:"uuid"`
	DeviceUUID     string   `json:"device_uuid"`
	SiemensAddress string   `json:"siemensAddress"` // 西门子的地址字符串
	Tag            string   `json:"tag"`
	Alias          string   `json:"alias"`
	DataOrder      string   `json:"dataOrder"` // 字节序
	DataType       string   `json:"dataType"`
	Frequency      *uint64  `json:"frequency"`
	Weight         *float64 `json:"weight"`        // 权重
	Status         int      `json:"status"`        // 运行时数据
	LastFetchTime  uint64   `json:"lastFetchTime"` // 运行时数据
	Value          any      `json:"value"`         // 运行时数据
	ErrMsg         string   `json:"errMsg"`        // 运行时数据

}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/Siemens_data_sheet/export
 */

// SiemensPoints 获取Siemens_excel类型的点位数据
func SiemensPointsExport(c *gin.Context, ruleEngine typex.Rhilex) {
	deviceUuid, _ := c.GetQuery("device_uuid")

	var records []model.MSiemensDataPoint
	result := interdb.InterDb().Table("m_siemens_data_points").
		Where("device_uuid=?", deviceUuid).Find(&records)
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{
		"address", "tag", "alias", "type", "order", "weight", "frequency",
	}
	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			glogger.GLogger.Errorf("close excel file, err=%v", err)
		}
	}()
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	xlsx.SetSheetRow("Sheet1", cell, &Headers)

	if len(records) >= 1 {
		for idx, record := range records[0:] {
			Row := []string{
				record.SiemensAddress,
				record.Tag,
				record.Alias,
				record.DataBlockType,
				record.DataBlockOrder,
				fmt.Sprintf("%f", *record.Weight),
				fmt.Sprintf("%d", *record.Frequency),
			}
			cell, _ = excelize.CoordinatesToCellName(1, idx+2)
			xlsx.SetSheetRow("Sheet1", cell, &Row)
		}
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%v.xlsx",
		time.Now().UnixMilli()))
	xlsx.WriteTo(c.Writer)

}

// 分页获取
// SELECT * FROM `m_Siemens_data_points` WHERE
// `m_Siemens_data_points`.`device_uuid` = "DEVICEDQNLO8"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func SiemensSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.InterDb()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.InterDb().Model(&model.MSiemensDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MSiemensDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MSiemensDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []SiemensPointVo{}
	for _, record := range records {
		Slot := intercache.GetSlot(deviceUuid)
		value, ok := Slot[record.UUID]
		Vo := SiemensPointVo{
			UUID:           record.UUID,
			DeviceUUID:     record.DeviceUuid,
			SiemensAddress: record.SiemensAddress,
			Tag:            record.Tag,
			Alias:          record.Alias,
			Frequency:      record.Frequency,
			DataType:       record.DataBlockType,
			DataOrder:      record.DataBlockOrder,
			Weight:         record.Weight.ToFloat64(),
			LastFetchTime:  value.LastFetchTime, // 运行时
			Value:          value.Value,         // 运行时
			ErrMsg:         value.ErrMsg,
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
func SiemensSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllSiemensPointByDevice(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}
func SiemensSheetDelete(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteSiemensPointByDevice(form.UUIDs, form.DeviceUUID)
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
func CheckSiemensDataPoints(M SiemensPointVo) error {
	// Helper function to check string length
	checkStringLength := func(value, paramName string, maxLength int) error {
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
	if err := checkStringLength(M.SiemensAddress, "address", 64); err != nil { // Assuming a max length of 256 for SiemensAddress
		return err
	}

	// Check Frequency
	if M.Frequency == nil {
		return fmt.Errorf("missing required param 'frequency'")
	}
	if *M.Frequency < 1 {
		return fmt.Errorf("'frequency' must be greater than 50ms")
	}
	if *M.Frequency > 100000 {
		return fmt.Errorf("'frequency' must be less than 100s")
	}

	// Validate DataOrder for different DataTypes
	dataOrderMap := map[string][]string{
		"I":        {"A"},
		"Q":        {"A"},
		"BYTE":     {"A"},
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
		return fmt.Errorf("invalid Weight value: %v", M.Weight)
	}

	// Validate Tag Name
	if !utils.IsValidColumnName(M.Tag) {
		return fmt.Errorf("invalid Tag Name: %s", M.Tag)
	}

	return nil
}

/*
*
* 更新点位表
*
 */
func SiemensSheetUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID        string           `json:"device_uuid"`
		SiemensDataPoints []SiemensPointVo `json:"data_points"`
	}
	form := Form{}
	// SiemensDataPoints := []SiemensPointVo{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, SiemensDataPoint := range form.SiemensDataPoints {
		if err := CheckSiemensDataPoints(SiemensDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if SiemensDataPoint.UUID == "" ||
			SiemensDataPoint.UUID == "new" ||
			SiemensDataPoint.UUID == "copy" {
			NewRow := model.MSiemensDataPoint{}
			copier.Copy(&NewRow, &SiemensDataPoint)
			NewRow.DeviceUuid = SiemensDataPoint.DeviceUUID
			NewRow.UUID = utils.SiemensPointUUID()
			NewRow.DataBlockType = SiemensDataPoint.DataType
			NewRow.DataBlockOrder = SiemensDataPoint.DataOrder
			NewRow.Weight = model.NewDecimal(*SiemensDataPoint.Weight)
			err0 := service.InsertSiemensPointPosition(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MSiemensDataPoint{}
			copier.Copy(&OldRow, &SiemensDataPoint)
			OldRow.DeviceUuid = SiemensDataPoint.DeviceUUID
			OldRow.UUID = SiemensDataPoint.UUID
			OldRow.DataBlockType = SiemensDataPoint.DataType
			OldRow.DataBlockOrder = SiemensDataPoint.DataOrder
			OldRow.Weight = model.NewDecimal(*utils.HandleZeroValue(SiemensDataPoint.Weight))

			err0 := service.UpdateSiemensPoint(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

// SiemensSheetImport 上传Excel文件
func SiemensSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
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
	type DeviceDto struct {
		UUID string
		Name string
		Type string
	}
	Device := DeviceDto{}
	errDb := interdb.InterDb().Table("m_devices").
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
	if Device.Type != typex.SIEMENS_PLC.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Siemens Device"))
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
	list, err := parseSiemensPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertSiemensPointPositions(list); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(deviceUuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

func parseSiemensPointExcel(
	r io.Reader,
	sheetName string,
	deviceUuid string) (list []model.MSiemensDataPoint, err error) {
	excelFile, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		excelFile.Close()
	}()
	// 读取表格
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	// 判断首行标头
	//
	err1 := errors.New("Invalid Sheet Header")
	if len(rows[0]) < 7 {
		return nil, err1
	}
	// Address Tag Alias Type Order Frequency

	if strings.ToLower(rows[0][0]) != "address" ||
		strings.ToLower(rows[0][1]) != "tag" ||
		strings.ToLower(rows[0][2]) != "alias" ||
		strings.ToLower(rows[0][3]) != "type" ||
		strings.ToLower(rows[0][4]) != "order" ||
		strings.ToLower(rows[0][5]) != "weight" ||
		strings.ToLower(rows[0][6]) != "frequency" {
		return nil, err1
	}

	list = make([]model.MSiemensDataPoint, 0)
	// Address Tag Alias Type Order Frequency
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		SiemensAddress := row[0]
		Tag := row[1]
		Alias := row[2]
		Type := row[3]
		Order := row[4]
		Weight, _ := strconv.ParseFloat(row[5], 32)
		limitedWeight := float64(int(Weight*100)) / 100.0
		if Weight == 0 {
			Weight = 1 // 防止解析异常的时候系数0
		}
		frequency, _ := strconv.ParseUint(row[6], 10, 64)
		Frequency := uint64(frequency)
		_, errParse1 := utils.ParseSiemensDB(SiemensAddress)
		if errParse1 != nil {
			return nil, errParse1
		}
		model := model.MSiemensDataPoint{
			UUID:           utils.SiemensPointUUID(),
			DeviceUuid:     deviceUuid,
			SiemensAddress: SiemensAddress,
			Tag:            Tag,
			Alias:          Alias,
			DataBlockType:  Type,
			DataBlockOrder: utils.GetDefaultDataOrder(Type, Order),
			Frequency:      &Frequency,
			Weight:         model.NewDecimal(limitedWeight),
		}
		list = append(list, model)
	}
	return list, nil
}
