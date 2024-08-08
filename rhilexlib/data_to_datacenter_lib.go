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
package rhilexlib

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/component/datacenter"
	"github.com/hootrhino/rhilex/component/dataschema"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

// GenInsertSql 生成 INSERT SQL 语句
func GenInsertSql(tableName string, rowList []kvp) (string, error) {
	if len(rowList) == 0 {
		return "", fmt.Errorf("no rows to insert")
	}
	columnStr := strings.Join(buildColumns(rowList), ", ")
	values := make([]string, len(rowList))
	for i, Row := range rowList {
		value := Row.V
		valueStr, err := formatValue(value)
		if err != nil {
			return "", err
		}
		values[i] = valueStr
	}
	insertSql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnStr, strings.Join(values, ", "))
	return insertSql, nil
}

// buildColumns 构建一个包含列名的字符串列表
func buildColumns(rowList []kvp) []string {
	var columnList []string
	for _, Row := range rowList {
		column := Row.K
		columnList = append(columnList, fmt.Sprintf("`%s`", column))
	}
	return columnList
}

// formatValue 根据值的类型格式化 SQL 值
func formatValue(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		if v {
			return "1", nil
		} else {
			return "0", nil
		}
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")), nil
	default:
		return "", fmt.Errorf("unsupported type: %s", reflect.TypeOf(val))
	}
}

/*
*
* 数据写入Sqlite，需要验证规则,验证规则的方式:
* 1 从全局缓存里面拿出来属性和验证器
* 2 在这里执行Validator
* 问题： 如何避免性能问题？
 */

func InsertToDataCenterTable(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		schema_uuid := l.ToString(2)
		kvs := l.ToTable(3)
		if kvs == nil {
			l.Push(lua.LString("missing table fields"))
			return 1
		}
		RowList := []kvp{}
		// create_at
		RowList = append(RowList, kvp{
			"create_at", time.Now(),
		})
		kvs.ForEach(func(k, v lua.LValue) {
			Row := kvp{}
			// K 只能String
			if k.Type() == lua.LTString {
				// create_at id : 不允许用户填写
				if Row.K != "create_at" && Row.K != "id" {
					switch v.Type() {
					case lua.LTString:
						Row.K = lua.LVAsString(k)
						Row.V = lua.LVAsString(v)
					case lua.LTNumber:
						Row.K = lua.LVAsString(k)
						Row.V = float64(lua.LVAsNumber(v))
					case lua.LTBool:
						Row.K = lua.LVAsString(k)
						Row.V = bool(lua.LVAsBool(v))
					case lua.LTNil:
						Row.K = lua.LVAsString(k)
						Row.V = ""
					default:
						Row.K = lua.LVAsString(k)
						Row.V = "" // 不支持其他类型
					}
					RowList = append(RowList, Row)
				}
			}
		})
		// TODO: checkRule
		if errCheckRule := checkRule(RowList); errCheckRule != nil {
			glogger.GLogger.Error("checkRule error:", errCheckRule)
			l.Push(lua.LString(errCheckRule.Error()))
			return 1
		}
		TableName := fmt.Sprintf("data_center_%s", schema_uuid)
		if errSave := saveToDataCenter(TableName, RowList); errSave != nil {
			l.Push(lua.LString(errSave.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}

/*
*
* 键值对
*
 */
type kvp struct {
	K string
	V interface{}
}

/*
*
* 检查规则
*
 */
func checkRule(RowList []kvp) error {
	for _, Row := range RowList {
		if Row.K == "create_at" || Row.K == "id" {
			continue
		}
		IoTProperty, ok := dataschema.GetDataSchemaCache(Row.K)
		if ok {
			IoTProperty.Value = Row.V
			err := IoTProperty.ValidateRule()
			if err != nil {
				return fmt.Errorf("filed '%s' invalid, %s", Row.K, err.Error())
			}
		}
	}
	return nil
}

// Save to local DataCenter
func saveToDataCenter(schema_uuid string, RowList []kvp) error {
	Sql, err0 := GenInsertSql(schema_uuid, RowList)
	glogger.GLogger.Debug(Sql)
	if err0 != nil {
		return err0
	}
	err1 := datacenter.DB().Exec(Sql).Error
	if err1 != nil {
		return err1
	}
	return nil
}

/*
*
* 检查属性约束
*
 */
func CheckSchemaConsist(schema_uuid string, RowList [][2]interface{}) error {
	return nil
}

// Query List
// GET /api/v1/datacenter/queryDataList?uuid=<uuid>&current=<current>&size=<size>&order=<order>&select=<select> HTTP/1.1
// Host: 127.0.0.1:2580
func QueryDataCenterList(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		schema_uuid := l.ToString(2)
		page := l.ToNumber(3)
		size := l.ToNumber(4)
		fields := l.ToString(5)
		if err := queryDataCenterList(fmt.Sprintf("data_center_%s", schema_uuid),
			int(page), int(size), fields); err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func queryDataCenterList(schema_uuid string, page, size int, fields string) error {
	glogger.GLogger.Debug(schema_uuid, page, size, fields)
	return nil
}

/*
*
* last data
*
 */
func QueryDataCenterLast(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		schema_uuid := l.ToString(2)
		fields := l.ToString(3)
		if err := queryDataCenterLast(schema_uuid, fields); err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func queryDataCenterLast(schema_uuid string, fields string) error {
	glogger.GLogger.Debug(schema_uuid, fields)
	return nil
}

/*
*
* 更新最后一行，而不是插入
 */
func UpdateDataCenterLast(rx typex.Rhilex, uuid string) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		schema_uuid := l.ToString(2)
		kvs := l.ToTable(3)
		if kvs == nil {
			l.Push(lua.LString("missing table fields"))
			return 1
		}
		RowList := []kvp{}
		kvs.ForEach(func(k, v lua.LValue) {
			Row := kvp{}
			// K 只能String
			if k.Type() == lua.LTString {
				// create_at 不允许用户填写
				if Row.K != "create_at" && Row.K != "id" {
					switch v.Type() {
					case lua.LTString:
						Row.K = lua.LVAsString(k)
						Row.V = lua.LVAsString(v)
					case lua.LTNumber:
						Row.K = lua.LVAsString(k)
						Row.V = float64(lua.LVAsNumber(v))
					case lua.LTBool:
						Row.K = lua.LVAsString(k)
						Row.V = bool(lua.LVAsBool(v))
					case lua.LTNil:
						Row.K = lua.LVAsString(k)
						Row.V = ""
					default:
						Row.K = lua.LVAsString(k)
						Row.V = "" // 不支持其他类型
					}
					RowList = append(RowList, Row)
				}
			}
		})
		// SELECT id FROM data_center_%s ORDER BY id DESC LIMIT 1;
		id := -1
		if err := datacenter.DB().
			Raw(fmt.Sprintf("SELECT id FROM data_center_%s ORDER BY id DESC LIMIT 1;",
				schema_uuid)).Scan(&id).Error; err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		if id < 0 {
			if err := saveToDataCenter(fmt.Sprintf("data_center_%s", schema_uuid), RowList); err != nil {
				l.Push(lua.LString(err.Error()))
			} else {
				l.Push(lua.LNil)
			}
		} else {
			if err := updateLast(fmt.Sprintf("data_center_%s", schema_uuid), RowList); err != nil {
				l.Push(lua.LString(err.Error()))
			} else {
				l.Push(lua.LNil)
			}
		}
		return 1
	}
}

// GenUpdateSql 生成 Update SQL 语句
// -- 1. 添加自增的主键列
//     ALTER TABLE your_table ADD COLUMN id INTEGER PRIMARY KEY AUTOINCREMENT;
// -- 2. 找到最后一行数据
//     SELECT id FROM your_table ORDER BY id DESC LIMIT 1;
// -- 3. 使用找到的最后一行的 id 进行更新
//        UPDATE your_table
//        SET K1 = 'new_value1',
//            K2 = 'new_value2',
//            K3 = 'new_value3'
//        WHERE condition;

func GenUpdateSql(tableName string, rowList []kvp) (string, error) {
	if len(rowList) == 0 {
		return "", fmt.Errorf("no rows to update")
	}
	fieldValuePairs := []string{}
	for _, Row := range rowList {
		fieldValuePairs = append(fieldValuePairs, fmt.Sprintf("%v=%v", Row.K, Row.V))
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE id=%v;", tableName, strings.Join(fieldValuePairs, ","), 1)
	return sql, nil
}

// Save to local DataCenter
func updateLast(schema_uuid string, RowList []kvp) error {
	Sql, err0 := GenUpdateSql(schema_uuid, RowList)
	glogger.GLogger.Debug(Sql)
	if err0 != nil {
		return err0
	}
	err1 := datacenter.DB().Exec(Sql).Error
	if err1 != nil {
		return err1
	}
	return nil
}
