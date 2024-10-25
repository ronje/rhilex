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

func InitSzy2062016Route() {
	Api := server.RouteGroup(server.ContextUrl("/szy2062016_master_sheet"))
	{
		Api.POST(("/sheetImport"), server.AddRoute(Szy2062016MasterSheetImport))
		Api.GET(("/sheetExport"), server.AddRoute(Szy2062016MasterPointsExport))
		Api.GET(("/list"), server.AddRoute(Szy2062016MasterSheetPageList))
		Api.POST(("/update"), server.AddRoute(Szy2062016MasterSheetUpdate))
		Api.DELETE(("/delIds"), server.AddRoute(Szy2062016MasterSheetDelete))
		Api.DELETE(("/delAll"), server.AddRoute(Szy2062016MasterSheetDeleteAll))
	}
}

type Szy2062016MasterPointVo struct {
	DeviceUuid    string      `json:"device_uuid"`
	UUID          string      `json:"uuid"`
	MeterId       string      `json:"meterId"`
	MeterType     string      `json:"meterType"`
	Tag           string      `json:"tag"`
	Alias         string      `json:"alias"`
	Frequency     uint64      `json:"frequency"`
	Status        int         `json:"status"`        // 运行时数据
	LastFetchTime uint64      `json:"lastFetchTime"` // 运行时数据
	Value         interface{} `json:"value"`         // 运行时数据
	ErrMsg        string      `json:"errMsg"`        // 运行时数据
}

/*
*
* 特殊设备需要和外界交互，这里主要就是一些设备的点位表导入导出等支持
*  http://127.0.0.1:2580/api/v1/Szy2062016_data_sheet/export
 */

// Szy2062016MasterPoints 获取Szy2062016_excel类型的点位数据
func Szy2062016MasterPointsExport(c *gin.Context, ruleEngine typex.Rhilex) {
	deviceUuid, _ := c.GetQuery("device_uuid")
	var records []model.MSzy2062016DataPoint
	result := interdb.DB().Table("m_szy2062016_data_points").
		Where("device_uuid=?", deviceUuid).Find(&records)
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	Headers := []string{"MeterId", "MeterType", "Tag", "Alias", "Frequency"}
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
			record.MeterId,
			record.MeterType,
			record.Tag,
			record.Alias,
			fmt.Sprintf("%d", record.Frequency),
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
// `m_szy2062016_data_points`.`device_uuid` = "UUID"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func Szy2062016MasterSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MSzy2062016DataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MSzy2062016DataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MSzy2062016DataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []Szy2062016MasterPointVo{}
	// "MeterId", "MeterType", "Tag", "Alias", "Frequency"
	for _, record := range records {
		Slot := intercache.GetSlot(deviceUuid)
		value, ok := Slot[record.UUID]
		Vo := Szy2062016MasterPointVo{
			UUID:          record.UUID,
			DeviceUuid:    record.DeviceUuid,
			MeterId:       record.MeterId,
			MeterType:     record.MeterType,
			Tag:           record.Tag,
			Alias:         record.Alias,
			Frequency:     record.Frequency,
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
func Szy2062016MasterSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllMSzy2062016ByDevice(form.DeviceUUID)
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
func Szy2062016MasterSheetDelete(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteSzy2062016PointByDevice(form.UUIDs, form.DeviceUUID)
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

func CheckSzy2062016MasterDataPoints(M Szy2062016MasterPointVo) error {
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
	if M.Frequency < 1 {
		return fmt.Errorf("'frequency' must be greater than 50ms")
	}
	if M.Frequency > 100000 {
		return fmt.Errorf("'frequency' must be less than 100s")
	}
	return nil
}

/*
*
* 更新点位表
*
 */
func Szy2062016MasterSheetUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID                 string                    `json:"device_uuid"`
		Szy2062016MasterDataPoints []Szy2062016MasterPointVo `json:"data_points"`
	}
	//  Szy2062016MasterDataPoints := [] Szy2062016MasterPointVo{}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, Szy2062016MasterDataPoint := range form.Szy2062016MasterDataPoints {
		if err := CheckSzy2062016MasterDataPoints(Szy2062016MasterDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// "MeterId", "MeterType", "Tag", "Alias", "Frequency"
		if Szy2062016MasterDataPoint.UUID == "" ||
			Szy2062016MasterDataPoint.UUID == "new" ||
			Szy2062016MasterDataPoint.UUID == "copy" {
			NewRow := model.MSzy2062016DataPoint{
				UUID:       utils.Szy2062016PointUUID(),
				DeviceUuid: Szy2062016MasterDataPoint.DeviceUuid,
				MeterId:    Szy2062016MasterDataPoint.MeterId,
				MeterType:  Szy2062016MasterDataPoint.MeterType,
				Tag:        Szy2062016MasterDataPoint.Tag,
				Alias:      Szy2062016MasterDataPoint.Alias,
				Frequency:  Szy2062016MasterDataPoint.Frequency,
			}
			err0 := service.InsertSzy2062016Point(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MSzy2062016DataPoint{
				UUID:       Szy2062016MasterDataPoint.UUID,
				DeviceUuid: Szy2062016MasterDataPoint.DeviceUuid,
				MeterId:    Szy2062016MasterDataPoint.MeterId,
				MeterType:  Szy2062016MasterDataPoint.MeterType,
				Tag:        Szy2062016MasterDataPoint.Tag,
				Alias:      Szy2062016MasterDataPoint.Alias,
				Frequency:  Szy2062016MasterDataPoint.Frequency,
			}
			err0 := service.UpdateSzy2062016Point(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

type Szy2062016DeviceDto struct {
	UUID   string
	Name   string
	Type   string
	Config string
}

func (md Szy2062016DeviceDto) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

// Szy2062016MasterSheetImport 上传Excel文件
func Szy2062016MasterSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
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
	if Device.Type != typex.SZY2062016_MASTER.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Szy2062016Master Device"))
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
	list, err := parseSzy2062016MasterPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertSzy2062016Points(list); err != nil {
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

func parseSzy2062016MasterPointExcel(r io.Reader, sheetName string,
	deviceUuid string) (list []model.MSzy2062016DataPoint, err error) {
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
	// "MeterId", "MeterType", "Tag", "Alias", "Frequency"
	err1 := errors.New(" Invalid Sheet Header")
	if len(rows[0]) < 5 {
		return nil, err1
	}
	// "MeterId", "MeterType", "Tag", "Alias", "Frequency"

	// 严格检查表结构
	if rows[0][0] != "MeterId" ||
		rows[0][1] != "MeterType" ||
		rows[0][2] != "Tag" ||
		rows[0][3] != "Alias" ||
		rows[0][4] != "Frequency" {
		return nil, err1
	}

	list = make([]model.MSzy2062016DataPoint, 0)
	// "MeterId", "MeterType", "Tag", "Alias", "Frequency"
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		MeterId := row[0]
		MeterType := row[1]
		Tag := row[2]
		Alias := row[3]
		Frequency, _ := strconv.ParseUint(row[4], 10, 64)
		// "MeterId", "MeterType", "Tag", "Alias", "Frequency"
		if err := CheckSzy2062016MasterDataPoints(Szy2062016MasterPointVo{
			MeterId:   MeterId,
			MeterType: MeterType,
			Tag:       Tag,
			Alias:     Alias,
			Frequency: Frequency,
		}); err != nil {
			return nil, err
		}
		//
		model := model.MSzy2062016DataPoint{
			UUID:       utils.Szy2062016PointUUID(),
			DeviceUuid: deviceUuid,
			MeterId:    MeterId,
			MeterType:  MeterType,
			Tag:        Tag,
			Alias:      Alias,
			Frequency:  Frequency,
		}
		list = append(list, model)
	}
	return list, nil
}
