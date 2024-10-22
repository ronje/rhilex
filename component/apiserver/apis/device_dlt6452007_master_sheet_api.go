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

func InitDlt6452007Route() {
	Api := server.RouteGroup(server.ContextUrl("/dlt6452007_master_sheet"))
	{
		Api.POST(("/sheetImport"), server.AddRoute(Dlt6452007MasterSheetImport))
		Api.GET(("/sheetExport"), server.AddRoute(Dlt6452007MasterPointsExport))
		Api.GET(("/list"), server.AddRoute(Dlt6452007MasterSheetPageList))
		Api.POST(("/update"), server.AddRoute(Dlt6452007MasterSheetUpdate))
		Api.DELETE(("/delIds"), server.AddRoute(Dlt6452007MasterSheetDelete))
		Api.DELETE(("/delAll"), server.AddRoute(Dlt6452007MasterSheetDeleteAll))
	}
}

type Dlt6452007MasterPointVo struct {
	DeviceUuid    string      `json:"device_uuid"`
	UUID          string      `json:"uuid"`
	MeterId       string      `json:"meterId"`
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
*  http://127.0.0.1:2580/api/v1/dlt6452007_data_sheet/export
 */

// Dlt6452007MasterPoints 获取Dlt6452007_excel类型的点位数据
func Dlt6452007MasterPointsExport(c *gin.Context, ruleEngine typex.Rhilex) {
	deviceUuid, _ := c.GetQuery("device_uuid")
	var records []model.MSzy2062016DataPoint
	result := interdb.DB().Table("m_dlt6452007_data_points").
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
// `m_dlt6452007_data_points`.`device_uuid` = "UUID"
// ORDER BY
// created_at DESC LIMIT 2 OFFSET 10
func Dlt6452007MasterSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MDlt6452007DataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MDlt6452007DataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MDlt6452007DataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVo := []Dlt6452007MasterPointVo{}
	// "MeterId", "Tag", "Alias", "Frequency"
	for _, record := range records {
		Slot := intercache.GetSlot(deviceUuid)
		value, ok := Slot[record.UUID]
		Vo := Dlt6452007MasterPointVo{
			UUID:          record.UUID,
			DeviceUuid:    record.DeviceUuid,
			MeterId:       record.MeterId,
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
			Vo.Value = value.Value
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
func Dlt6452007MasterSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllMDlt6452007ByDevice(form.DeviceUUID)
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
func Dlt6452007MasterSheetDelete(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteDlt6452007PointByDevice(form.UUIDs, form.DeviceUUID)
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

func CheckDlt6452007MasterDataPoints(M Dlt6452007MasterPointVo) error {
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
func Dlt6452007MasterSheetUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID                 string                    `json:"device_uuid"`
		Dlt6452007MasterDataPoints []Dlt6452007MasterPointVo `json:"data_points"`
	}
	//  Dlt6452007MasterDataPoints := [] Dlt6452007MasterPointVo{}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, Dlt6452007MasterDataPoint := range form.Dlt6452007MasterDataPoints {
		if err := CheckDlt6452007MasterDataPoints(Dlt6452007MasterDataPoint); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// "MeterId", "Tag", "Alias", "Frequency"
		if Dlt6452007MasterDataPoint.UUID == "" ||
			Dlt6452007MasterDataPoint.UUID == "new" ||
			Dlt6452007MasterDataPoint.UUID == "copy" {
			NewRow := model.MDlt6452007DataPoint{
				UUID:       utils.Dlt6452007PointUUID(),
				DeviceUuid: Dlt6452007MasterDataPoint.DeviceUuid,
				MeterId:    Dlt6452007MasterDataPoint.MeterId,
				Tag:        Dlt6452007MasterDataPoint.Tag,
				Alias:      Dlt6452007MasterDataPoint.Alias,
				Frequency:  Dlt6452007MasterDataPoint.Frequency,
			}
			err0 := service.InsertDlt6452007Point(NewRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		} else {
			OldRow := model.MDlt6452007DataPoint{
				UUID:       Dlt6452007MasterDataPoint.UUID,
				DeviceUuid: Dlt6452007MasterDataPoint.DeviceUuid,
				MeterId:    Dlt6452007MasterDataPoint.MeterId,
				Tag:        Dlt6452007MasterDataPoint.Tag,
				Alias:      Dlt6452007MasterDataPoint.Alias,
				Frequency:  Dlt6452007MasterDataPoint.Frequency,
			}
			err0 := service.UpdateDlt6452007Point(OldRow)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
		}
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())

}

type Dlt6452007DeviceDto struct {
	UUID   string
	Name   string
	Type   string
	Config string
}

func (md Dlt6452007DeviceDto) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

// Dlt6452007MasterSheetImport 上传Excel文件
func Dlt6452007MasterSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
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
	if Device.Type != typex.DLT6452007_MASTER.String() {
		c.JSON(common.HTTP_OK,
			common.Error("Invalid Device Type, Only Support Import Dlt6452007Master Device"))
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
	list, err := parseDlt6452007MasterPointExcel(file, "Sheet1", deviceUuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err = service.InsertDlt6452007Points(list); err != nil {
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

func parseDlt6452007MasterPointExcel(r io.Reader, sheetName string,
	deviceUuid string) (list []model.MDlt6452007DataPoint, err error) {
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
	// "MeterId", "Tag", "Alias", "Frequency"
	err1 := errors.New(" Invalid Sheet Header")
	if len(rows[0]) < 4 {
		return nil, err1
	}
	// "MeterId", "Tag", "Alias", "Frequency"

	// 严格检查表结构
	if rows[0][0] != "MeterId" ||
		rows[0][1] != "Tag" ||
		rows[0][2] != "Alias" ||
		rows[0][3] != "Frequency" {
		return nil, err1
	}

	list = make([]model.MDlt6452007DataPoint, 0)
	// "MeterId", "Tag", "Alias", "Frequency"
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		MeterId := row[0]
		Tag := row[1]
		Alias := row[2]
		Frequency, _ := strconv.ParseUint(row[3], 10, 64)
		// "MeterId", "Tag", "Alias", "Frequency"
		if err := CheckDlt6452007MasterDataPoints(Dlt6452007MasterPointVo{
			MeterId:   MeterId,
			Tag:       Tag,
			Alias:     Alias,
			Frequency: Frequency,
		}); err != nil {
			return nil, err
		}
		//
		model := model.MDlt6452007DataPoint{
			UUID:       utils.Dlt6452007PointUUID(),
			DeviceUuid: deviceUuid,
			MeterId:    MeterId,
			Tag:        Tag,
			Alias:      Alias,
			Frequency:  Frequency,
		}
		list = append(list, model)
	}
	return list, nil
}
