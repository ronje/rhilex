// Copyright (C) 2024 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/*
*
// ## 生成模板属性
// - 温湿度计          (TEMP_HUMIDITY)
// - 开关状态          (SWITCH_STATUS)
// - 水质传感器        (WATER_QUALITY)
// - 空气质量传感器     (AIR_QUALITY)
// - 动作传感器        (MOTION_SENSOR)
// - 智能电表          (SMART_METER)
// - 土壤湿度传感器     (SOIL_MOISTURE)
// - GPS追踪器         (GPS_TRACKER)
// - 烟雾探测器         (SMOKE_DETECTOR)
// - 智能锁            (SMART_LOCK)
// - 六轴加速度计       (SIX_AXIS_ACCELEROMETER)
*/
func GenerateSchemaTemplate(c *gin.Context, ruleEngine typex.Rhilex) {
	schemaId, _ := c.GetQuery("schemaId")
	templateId, _ := c.GetQuery("templateId")
	_, err := service.GetDataSchemaWithUUID(schemaId)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if templateId == "TEMP_HUMIDITY" {
		if err := _Template_TEMP_HUMIDITY(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SWITCH_STATUS" {
		if err := _Template_SWITCH_STATUS(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "WATER_QUALITY" {
		if err := _Template_WATER_QUALITY(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "AIR_QUALITY" {
		if err := _Template_AIR_QUALITY(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "MOTION_SENSOR" {
		if err := _Template_MOTION_SENSOR(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SMART_METER" {
		if err := _Template_SMART_METER(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SOIL_MOISTURE" {
		if err := _Template_SOIL_MOISTURE(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "GPS_TRACKER" {
		if err := _Template_GPS_TRACKER(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SMOKE_DETECTOR" {
		if err := _Template_SMOKE_DETECTOR(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SMART_LOCK" {
		if err := _Template_SMART_LOCK(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	if templateId == "SIX_AXIS_ACCELEROMETER" {
		if err := _Template_SIX_AXIS_ACCELEROMETER(schemaId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("template note exists:"+templateId))
}

// TEMP_HUMIDITY
func _Template_TEMP_HUMIDITY(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name =?)",
		schemaId, "temperature", "humidity").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[temperature, humidity]")
	}
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":64,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '温度', 'temperature', 'FLOAT', 'R', '℃',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '温度值'),
(CURRENT_TIMESTAMP, '%s', '%s', '湿度', 'humidity', 'FLOAT', 'R', '%%',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '湿度值');
`
	return tx.Exec(fmt.Sprintf(sql, schemaId, uuid1, schemaId, uuid2, schemaId, uuid3)).Error
}

// SWITCH_STATUS
func _Template_SWITCH_STATUS(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and name=?",
		schemaId, "status").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[status]")
	}
	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '开关状态', 'status', 'BOOL', 'RW', '-',
    '{"defaultValue":"0","max":1000,"min":0,"trueLabel":"开","falseLabel":"关","round":2}', 'true 为开, false 为关');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql, schemaId, uuid1, schemaId, uuid2)).Error

}

// WATER_QUALITY
func _Template_WATER_QUALITY(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name =? or name =? or name =? or name=?)",
		schemaId, "ph", "turbidity", "turbidity", "dissolved_oxygen", "conductivity", "temperature").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[ph, turbidity, dissolved_oxygen, conductivity, temperature]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', 'pH值', 'ph', 'FLOAT', 'R', 'mol/μl',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'pH值'),
(CURRENT_TIMESTAMP, '%s', '%s', '浊度', 'turbidity', 'FLOAT', 'R', 'NTU',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '浊度'),
(CURRENT_TIMESTAMP, '%s', '%s', '溶解氧', 'dissolved_oxygen', 'FLOAT', 'R', 'mg/L',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '溶解氧'),
(CURRENT_TIMESTAMP, '%s', '%s', '电导率', 'conductivity', 'FLOAT', 'R', 'μS/cm',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电导率'),
(CURRENT_TIMESTAMP, '%s', '%s', '水温', 'temperature', 'FLOAT', 'R', '摄氏度',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '水温');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	uuid6 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5,
		schemaId, uuid6)).Error
}

// AIR_QUALITY
func _Template_AIR_QUALITY(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name =? or name =? or name =? or name=?)",
		schemaId, "pm25", "pm10", "co2", "tvoc", "temperature", "humidity").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[pm25, pm10, co2, tvoc, temperature, humidity]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', 'PM2.5浓度', 'pm25', 'FLOAT', 'R', 'μg/m³',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , 'PM2.5浓度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'PM10浓度', 'pm10', 'FLOAT', 'R', 'μg/m³',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , 'PM10浓度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'CO2浓度', 'co2', 'INT', 'R', 'ppm',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , 'CO2浓度'),
(CURRENT_TIMESTAMP, '%s', '%s', '总挥发性有机化合物', 'tvoc', 'FLOAT', 'R', 'ppb',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , '总挥发性有机化合物'),
(CURRENT_TIMESTAMP, '%s', '%s', '温度', 'temperature', 'FLOAT', 'R', '℃',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , '温度'),
(CURRENT_TIMESTAMP, '%s', '%s', '湿度', 'humidity', 'FLOAT', 'R', '%%',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}' , '湿度');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	uuid6 := utils.MakeUUID("PROP")
	uuid7 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5,
		schemaId, uuid6,
		schemaId, uuid7)).Error
}

// MOTION_SENSOR
func _Template_MOTION_SENSOR(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name =?)",
		schemaId, "detected", "intensity", "battery").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[detected, intensity, battery]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '运动状态', 'detected', 'BOOL', 'R', '-',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '是否检测到运动'),
(CURRENT_TIMESTAMP, '%s', '%s', '运动强度', 'intensity', 'INTEGER', 'R', '-',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '运动强度'),
(CURRENT_TIMESTAMP, '%s', '%s', '电池电量', 'battery', 'INTEGER', 'R', '%%',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电池电量');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4)).Error
}

// SMART_METER
func _Template_SMART_METER(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name =? or name =?)",
		schemaId, "energy_consumption", "current", "voltage", "power_factor").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[energy_consumption, current, voltage, power_factor]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '能源消耗', 'energy_consumption', 'FLOAT', 'R', 'kWh',
   '{"defaultValue":"0","max":1000,"min":0,"round":2}', '能源消耗'),
(CURRENT_TIMESTAMP, '%s', '%s', '电流', 'current', 'FLOAT', 'R', 'A',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电流'),
(CURRENT_TIMESTAMP, '%s', '%s', '电压', 'voltage', 'FLOAT', 'R', 'V',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电压'),
(CURRENT_TIMESTAMP, '%s', '%s', '功率因数', 'power_factor', 'FLOAT', 'R', '-',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '功率因数');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5)).Error
}

// SOIL_MOISTURE
func _Template_SOIL_MOISTURE(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name =?)",
		schemaId, "moisture", "temperature", "ec").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[moisture, temperature, ec]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-',
	'{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '土壤湿度', 'moisture', 'FLOAT', 'R', '%%',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '土壤湿度'),
(CURRENT_TIMESTAMP, '%s', '%s', '土壤温度', 'temperature', 'FLOAT', 'R', '℃',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '土壤温度'),
(CURRENT_TIMESTAMP, '%s', '%s', '电导率', 'ec', 'FLOAT', 'R', 'mS/cm',
    '{"defaultValue":"0","max":1000,"min":0,"round":2}', '土壤电导率'),
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4)).Error
}

// GPS_TRACKER
func _Template_GPS_TRACKER(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name=? or name=? or name=?)",
		schemaId, "latitude", "longitude", "altitude", "speed", "battery").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[latitude, longitude, altitude, speed, battery]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '纬度', 'latitude', 'FLOAT', 'R', '°', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备的纬度'),
(CURRENT_TIMESTAMP, '%s', '%s', '经度', 'longitude', 'FLOAT', 'R', '°', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备的经度'),
(CURRENT_TIMESTAMP, '%s', '%s', '海拔', 'altitude', 'FLOAT', 'R', '米', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备的海拔高度'),
(CURRENT_TIMESTAMP, '%s', '%s', '速度', 'speed', 'FLOAT', 'R', 'km/h', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备移动速度'),
(CURRENT_TIMESTAMP, '%s', '%s', '电池电量百分比', 'battery', 'FLOAT', 'R', '%%', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电池电量百分比');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	uuid6 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5,
		schemaId, uuid6)).Error
}

// SMOKE_DETECTOR
func _Template_SMOKE_DETECTOR(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name=? or name=? or name=?)",
		schemaId, "smoke_detected", "co_level", "battery", "last_tested").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[smoke_detected, co_level, battery, last_tested]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '烟雾检测状态', 'smoke_detected', 'BOOLEAN', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '是否检测到烟雾'),
(CURRENT_TIMESTAMP, '%s', '%s', '一氧化碳水平', 'co_level', 'FLOAT', 'R', 'ppm', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '一氧化碳浓度'),
(CURRENT_TIMESTAMP, '%s', '%s', '电池电量百分比', 'battery', 'FLOAT', 'R', '%%', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电池电量百分比'),
(CURRENT_TIMESTAMP, '%s', '%s', '上次测试时间', 'last_tested', 'TIMESTAMP', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '上次测试时间');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5)).Error
}

// SMART_LOCK
func _Template_SMART_LOCK(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name=? or name=? or name=?)",
		schemaId, "lock_status", "access_method", "user_id", "battery").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[lock_status, access_method, user_id, battery]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', '锁状态', 'lock_status', 'BOOLEAN', 'RW', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '锁的状态'),
(CURRENT_TIMESTAMP, '%s', '%s', '访问方法', 'access_method', 'VARCHAR(20)', 'RW', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '访问方式'),
(CURRENT_TIMESTAMP, '%s', '%s', '操作用户ID', 'user_id', 'STRING', 'RW', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '操作用户ID'),
(CURRENT_TIMESTAMP, '%s', '%s', '电池电量百分比', 'battery', 'FLOAT', 'R', '%%', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电池电量百分比');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5)).Error
}

// SIX_AXIS_ACCELEROMETER
func _Template_SIX_AXIS_ACCELEROMETER(schemaId string) error {
	tx := interdb.InterDb()
	count := int64(0)
	tx.Model(model.MIotProperty{}).Where("schema_id = ? and (name=? or name=? or name=? or name=? or name=? or name=? or name=?)",
		schemaId, "accel_x", "accel_y", "accel_z", "gyro_x", "gyro_y", "gyro_z", "temperature", "battery").Count(&count)
	if count > 0 {
		return fmt.Errorf("Already exists fields :[accel_x, accel_y, accel_z, gyro_x, gyro_y, gyro_z, temperature, battery]")
	}

	sql := `
INSERT INTO m_iot_properties (created_at, schema_id, uuid, label, name, type, rw, unit, rule, description)
VALUES
(CURRENT_TIMESTAMP, '%s', '%s', '设备id', 'device_id', 'STRING', 'R', '-', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '设备id'),
(CURRENT_TIMESTAMP, '%s', '%s', 'X轴加速度', 'accel_x', 'FLOAT', 'R', 'g', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'X轴方向的加速度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'Y轴加速度', 'accel_y', 'FLOAT', 'R', 'g', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'Y轴方向的加速度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'Z轴加速度', 'accel_z', 'FLOAT', 'R', 'g', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'Z轴方向的加速度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'X轴角速度', 'gyro_x', 'FLOAT', 'R', 'deg/s', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'X轴方向的角速度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'Y轴角速度', 'gyro_y', 'FLOAT', 'R', 'deg/s', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'Y轴方向的角速度'),
(CURRENT_TIMESTAMP, '%s', '%s', 'Z轴角速度', 'gyro_z', 'FLOAT', 'R', 'deg/s', '{"defaultValue":"0","max":1000,"min":0,"round":2}', 'Z轴方向的角速度'),
(CURRENT_TIMESTAMP, '%s', '%s', '传感器温度', 'temperature', 'FLOAT', 'R', '℃', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '传感器的温度'),
(CURRENT_TIMESTAMP, '%s', '%s', '电池电量百分比', 'battery', 'FLOAT', 'R', '%%', '{"defaultValue":"0","max":1000,"min":0,"round":2}', '电池电量百分比');
`
	uuid1 := utils.MakeUUID("PROP")
	uuid2 := utils.MakeUUID("PROP")
	uuid3 := utils.MakeUUID("PROP")
	uuid4 := utils.MakeUUID("PROP")
	uuid5 := utils.MakeUUID("PROP")
	uuid6 := utils.MakeUUID("PROP")
	uuid7 := utils.MakeUUID("PROP")
	uuid8 := utils.MakeUUID("PROP")
	uuid9 := utils.MakeUUID("PROP")
	return tx.Exec(fmt.Sprintf(sql,
		schemaId, uuid1,
		schemaId, uuid2,
		schemaId, uuid3,
		schemaId, uuid4,
		schemaId, uuid5,
		schemaId, uuid6,
		schemaId, uuid7,
		schemaId, uuid8,
		schemaId, uuid9)).Error
}

// Field struct represents each field's name, type, and comment in the template.
type TemplateField struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

func GetTemplateFields(c *gin.Context, ruleEngine typex.Rhilex) {

	// Define the map of templates with different schema categories
	var templates = map[string][]TemplateField{
		"TEMP_HUMIDITY": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "temperature", Type: "FLOAT", Comment: "温度值，精确到小数点后两位"},
			{Name: "humidity", Type: "FLOAT", Comment: "湿度值，精确到小数点后两位"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
		"SWITCH_STATUS": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "status", Type: "BOOLEAN", Comment: "开关状态：true为开，false为关"},
			{Name: "last_changed", Type: "TIMESTAMP", Comment: "最后一次状态改变的时间"},
			{Name: "changed_by", Type: "STRING", Comment: "触发状态改变的用户或系统标识"},
		},
		"WATER_QUALITY": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "ph", Type: "FLOAT", Comment: "pH值，范围通常在0-14之间"},
			{Name: "turbidity", Type: "FLOAT", Comment: "浊度，单位可能是NTU"},
			{Name: "dissolved_oxygen", Type: "FLOAT", Comment: "溶解氧含量，单位mg/L"},
			{Name: "conductivity", Type: "FLOAT", Comment: "电导率，单位可能是μS/cm"},
			{Name: "temperature", Type: "FLOAT", Comment: "水温，单位摄氏度"},
		},
		"AIR_QUALITY": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "pm25", Type: "FLOAT", Comment: "PM2.5浓度，单位μg/m³"},
			{Name: "pm10", Type: "FLOAT", Comment: "PM10浓度，单位μg/m³"},
			{Name: "co2", Type: "INTEGER", Comment: "CO2浓度，单位ppm"},
			{Name: "tvoc", Type: "FLOAT", Comment: "总挥发性有机化合物，单位ppb"},
			{Name: "temperature", Type: "FLOAT", Comment: "温度，单位摄氏度"},
			{Name: "humidity", Type: "FLOAT", Comment: "湿度，百分比"},
		},
		"MOTION_SENSOR": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "motion_detected", Type: "BOOLEAN", Comment: "是否检测到运动"},
			{Name: "intensity", Type: "INTEGER", Comment: "运动强度，可选"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
		"SMART_METER": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "energy_consumption", Type: "FLOAT", Comment: "能源消耗，单位kWh"},
			{Name: "current", Type: "FLOAT", Comment: "电流，单位A"},
			{Name: "voltage", Type: "FLOAT", Comment: "电压，单位V"},
			{Name: "power_factor", Type: "FLOAT", Comment: "功率因数"},
		},
		"SOIL_MOISTURE": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "moisture_level", Type: "FLOAT", Comment: "土壤湿度水平，百分比"},
			{Name: "temperature", Type: "FLOAT", Comment: "土壤温度，单位摄氏度"},
			{Name: "ec", Type: "FLOAT", Comment: "土壤电导率，单位mS/cm"},
			{Name: "ph", Type: "FLOAT", Comment: "土壤pH值"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
		"GPS_TRACKER": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "latitude", Type: "FLOAT", Comment: "纬度"},
			{Name: "longitude", Type: "FLOAT", Comment: "经度"},
			{Name: "altitude", Type: "FLOAT", Comment: "海拔，单位米"},
			{Name: "speed", Type: "FLOAT", Comment: "速度，单位km/h"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
		"SMOKE_DETECTOR": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "smoke_detected", Type: "BOOLEAN", Comment: "是否检测到烟雾"},
			{Name: "co_level", Type: "FLOAT", Comment: "一氧化碳水平，单位ppm"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
			{Name: "last_tested", Type: "TIMESTAMP", Comment: "上次测试时间"},
		},
		"SMART_LOCK": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "lock_status", Type: "BOOLEAN", Comment: "锁状态：true为锁定，false为解锁"},
			{Name: "access_method", Type: "VARCHAR(20)", Comment: "访问方法：PIN、指纹、NFC等"},
			{Name: "user_id", Type: "STRING", Comment: "操作用户ID"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
		"SIX_AXIS_ACCELEROMETER": {
			{Name: "device_id", Type: "STRING", Comment: "设备唯一标识符"},
			{Name: "accel_x", Type: "FLOAT", Comment: "X轴加速度，单位g"},
			{Name: "accel_y", Type: "FLOAT", Comment: "Y轴加速度，单位g"},
			{Name: "accel_z", Type: "FLOAT", Comment: "Z轴加速度，单位g"},
			{Name: "gyro_x", Type: "FLOAT", Comment: "X轴角速度，单位deg/s"},
			{Name: "gyro_y", Type: "FLOAT", Comment: "Y轴角速度，单位deg/s"},
			{Name: "gyro_z", Type: "FLOAT", Comment: "Z轴角速度，单位deg/s"},
			{Name: "temperature", Type: "FLOAT", Comment: "传感器温度，单位摄氏度"},
			{Name: "battery", Type: "INTEGER", Comment: "电池电量百分比"},
		},
	}

	c.JSON(common.HTTP_OK, common.OkWithData(templates))

}

type TemplateSensorType struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func GetTemplates(c *gin.Context, ruleEngine typex.Rhilex) {
	sensorTypes := []TemplateSensorType{
		{Label: "温湿度传感器", Value: "TEMP_HUMIDITY"},
		{Label: "开关状态", Value: "SWITCH_STATUS"},
		{Label: "水质传感器", Value: "WATER_QUALITY"},
		{Label: "空气质量传感器", Value: "AIR_QUALITY"},
		{Label: "动作传感器", Value: "MOTION_SENSOR"},
		{Label: "智能电表", Value: "SMART_METER"},
		{Label: "土壤湿度传感器", Value: "SOIL_MOISTURE"},
		{Label: "GPS追踪器", Value: "GPS_TRACKER"},
		{Label: "烟雾探测器", Value: "SMOKE_DETECTOR"},
		{Label: "智能锁", Value: "SMART_LOCK"},
		{Label: "六轴加速度计", Value: "SIX_AXIS_ACCELEROMETER"},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(sensorTypes))
}
