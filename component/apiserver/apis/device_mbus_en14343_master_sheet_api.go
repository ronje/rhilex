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

func InitMBusRoute() {
	Api := server.RouteGroup(server.ContextUrl("/mbus_master_sheet"))
	{
		Api.POST(("/sheetImport"), server.AddRoute(MBusMasterSheetImport))
		Api.GET(("/sheetExport"), server.AddRoute(MBusMasterPointsExport))
		Api.GET(("/list"), server.AddRoute(MBusMasterSheetPageList))
		Api.POST(("/update"), server.AddRoute(MBusMasterSheetUpdate))
		Api.DELETE(("/delIds"), server.AddRoute(MBusMasterSheetDelete))
		Api.DELETE(("/delAll"), server.AddRoute(MBusMasterSheetDeleteAll))
	}
}

type MBusMasterPointVo struct {
	UUID          string      `json:"uuid"`
	DeviceUuid    string      `json:"device_uuid"`
	SlaverId      string      `json:"slaverId"`
	Type          string      `json:"type"`
	Tag           string      `json:"tag"`
	Alias         string      `json:"alias"`
	Manufacturer  string      `json:"manufacturer"`
	Frequency     *uint64     `json:"frequency"`
	DataLength    *uint64     `json:"dataLength"`
	Weight        *float64    `json:"weight"`
	Status        int         `json:"status"`        // 运行时数据
	LastFetchTime uint64      `json:"lastFetchTime"` // 运行时数据
	Value         interface{} `json:"value"`         // 运行时数据
	ErrMsg        string      `json:"errMsg"`        // 运行时数据

}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/mbus_data_sheet/export
 */

// MBusMasterPoints 获取MBus_excel类型的点位数据
func MBusMasterPointsExport(c *gin.Context, ruleEngine typex.Rhilex) {
	deviceUuid, _ := c.GetQuery("device_uuid")

	var records []model.MMBusDataPoint
	result := interdb.InterDb().Table("m_mbus_data_points").
		Where("device_uuid=?", deviceUuid).Find(&records)
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{"SlaverId", "Type", "Manufacturer", "Tag", "Alias", "Frequency", "DataLength", "Weight"}

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
			record.SlaverId,
			record.Type,
			record.Manufacturer,
			record.Tag,
			record.Alias,
			fmt.Sprintf("%d", *record.Frequency),
			fmt.Sprintf("%d", *record.DataLength),
			fmt.Sprintf("%.2f", *record.Weight),
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
// SELECT * FROM WHERE
// `m_mbus_data_points`.`device_uuid` = "UUID"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func MBusMasterSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.InterDb()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.InterDb().Model(&model.MMBusDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MMBusDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MMBusDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []MBusMasterPointVo{}
	// "SlaverId", "Type", "Manufacturer", "Tag", "Alias", "Frequency", "DataLength", "Weight"
	for _, record := range records {
		Slot := intercache.GetSlot(deviceUuid)
		value, ok := Slot[record.UUID]
		Vo := MBusMasterPointVo{
			UUID:          record.UUID,
			DeviceUuid:    record.DeviceUuid,
			SlaverId:      record.SlaverId,
			Type:          record.Type,
			Manufacturer:  record.Manufacturer,
			Tag:           record.Tag,
			Alias:         record.Alias,
			Frequency:     record.Frequency,
			DataLength:    record.DataLength,
			Weight:        record.Weight.ToFloat64(),
			LastFetchTime: value.LastFetchTime,
			Value:         value.Value,
			ErrMsg:        value.ErrMsg,
		}
		if ok {
			Vo.Status = func() int {
				if value.Value == "" {
					return 0
				}
				return 1
			}()
			Vo.LastFetchTime = value.LastFetchTime
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
func MBusMasterSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllMBusByDevice(form.DeviceUUID)
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
func MBusMasterSheetDelete(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteMBusPointByDevice(form.UUIDs, form.DeviceUUID)
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

func CheckMBusMasterDataPoints(M MBusMasterPointVo) error {
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

	return nil
}

/*
*
* 更新点位表
*
 */
func MBusMasterSheetUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID           string              `json:"device_uuid"`
		MBusMasterDataPoints []MBusMasterPointVo `json:"data_points"`
	}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, MBusMasterDataPoint := range form.MBusMasterDataPoints {
		if err := CheckMBusMasterDataPoints(MBusMasterDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if MBusMasterDataPoint.UUID == "" ||
			MBusMasterDataPoint.UUID == "new" ||
			MBusMasterDataPoint.UUID == "copy" {
			NewRow := model.MMBusDataPoint{
				UUID:         utils.MBusPointUUID(),
				DeviceUuid:   MBusMasterDataPoint.DeviceUuid,
				SlaverId:     MBusMasterDataPoint.SlaverId,
				Type:         MBusMasterDataPoint.Type,
				Manufacturer: MBusMasterDataPoint.Manufacturer,
				Tag:          MBusMasterDataPoint.Tag,
				Alias:        MBusMasterDataPoint.Alias,
				Frequency:    MBusMasterDataPoint.Frequency,
				DataLength:   MBusMasterDataPoint.DataLength,
				Weight:       model.NewDecimal(*MBusMasterDataPoint.Weight),
			}
			err0 := service.InsertMBusPoint(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MMBusDataPoint{
				UUID:         MBusMasterDataPoint.UUID,
				DeviceUuid:   MBusMasterDataPoint.DeviceUuid,
				SlaverId:     MBusMasterDataPoint.SlaverId,
				Type:         MBusMasterDataPoint.Type,
				Manufacturer: MBusMasterDataPoint.Manufacturer,
				Tag:          MBusMasterDataPoint.Tag,
				Alias:        MBusMasterDataPoint.Alias,
				Frequency:    MBusMasterDataPoint.Frequency,
				DataLength:   MBusMasterDataPoint.DataLength,
				Weight:       model.NewDecimal(*MBusMasterDataPoint.Weight),
			}
			err0 := service.UpdateMBusPoint(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

type MbusDeviceDto struct {
	UUID   string
	Name   string
	Type   string
	Config string
}

func (md MbusDeviceDto) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

// MBusMasterSheetImport 上传Excel文件
func MBusMasterSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
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
	if Device.Type != typex.GENERIC_MBUS_EN13433_MASTER.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import MBusMaster Device"))
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
	list, err := parseMBusMasterPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertMBusPoints(list); err != nil {
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

func parseMBusMasterPointExcel(r io.Reader, sheetName string,
	deviceUuid string) (list []model.MMBusDataPoint, err error) {
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
	// "SlaverId", "Type", "Manufacturer", "Tag", "Alias", "Frequency", "DataLength", "Weight"
	err1 := errors.New(" Invalid Sheet Header")
	if len(rows[0]) < 8 {
		return nil, err1
	}
	// 严格检查表结构
	if rows[0][0] != "SlaverId" ||
		rows[0][1] != "Type" ||
		rows[0][2] != "Manufacturer" ||
		rows[0][3] != "Tag" ||
		rows[0][4] != "Alias" ||
		rows[0][5] != "Frequency" ||
		rows[0][6] != "DataLength" ||
		rows[0][7] != "Weight" {
		return nil, err1
	}

	list = make([]model.MMBusDataPoint, 0)
	// "SlaverId", "Type", "Manufacturer", "Tag", "Alias", "Frequency", "DataLength", "Weight"
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		slaverId := row[0]
		Type := row[1]
		manufacturer := row[2]
		Tag := row[3]
		Alias := row[4]
		Frequency, _ := strconv.ParseUint(row[5], 10, 64)
		DataLength, _ := strconv.ParseUint(row[6], 10, 64)
		Weight, _ := strconv.ParseFloat(row[7], 64)
		limitedWeight := float64(int(Weight*100)) / 100.0
		if err := CheckMBusMasterDataPoints(MBusMasterPointVo{
			SlaverId:     slaverId,
			Type:         Type,
			Manufacturer: manufacturer,
			Tag:          Tag,
			Alias:        Alias,
			Frequency:    &Frequency,
			DataLength:   &DataLength,
			Weight:       &limitedWeight,
		}); err != nil {
			return nil, err
		}
		//
		model := model.MMBusDataPoint{
			UUID:         utils.MBusPointUUID(),
			DeviceUuid:   deviceUuid,
			SlaverId:     slaverId,
			Type:         Type,
			Manufacturer: manufacturer,
			Tag:          Tag,
			Alias:        Alias,
			Frequency:    &Frequency,
			DataLength:   &DataLength,
			Weight:       model.NewDecimal(Weight),
		}
		list = append(list, model)
	}
	return list, nil
}
