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

package dataschema

import (
	"encoding/json"
	"time"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 将所有规则加载进去
*
 */
func InitDataSchemaCache(e typex.Rhilex) {
	// interdb.DB().Exec(fmt.Sprintf(`
	// CREATE TRIGGER IF NOT EXISTS limit_m_internal_notifies
	// AFTER INSERT ON m_internal_notifies
	// BEGIN
	// 	DELETE FROM m_internal_notifies
	// 	WHERE id IN (
	// 		SELECT id FROM m_internal_notifies
	// 		ORDER BY id ASC
	// 		LIMIT (SELECT COUNT(*) - %d FROM m_internal_notifies)
	// 	);
	// END;`, 300))
	intercache.RegisterSlot("__DataSchema")
	MIotSchemas := []model.MIotSchema{}
	interdb.DB().Model(model.MIotSchema{}).Find(&MIotSchemas)
	for _, MIotSchema := range MIotSchemas {
		MIotProperties := []model.MIotProperty{}
		interdb.DB().Model(model.MIotProperty{}).
			Where("schema_id=?", MIotSchema.UUID).Find(&MIotProperties)
		for _, MIotProperty := range MIotProperties {
			CacheIoTProperty := &IoTProperty{
				UUID:        MIotProperty.UUID,
				Label:       MIotProperty.Label,
				Name:        MIotProperty.Name,
				Type:        IoTPropertyType(MIotProperty.Type),
				Rw:          MIotProperty.Rw,
				Unit:        MIotProperty.Unit,
				Description: MIotProperty.Description,
			}
			errUnmarshal := json.Unmarshal([]byte(MIotProperty.Rule), &CacheIoTProperty.Rule)
			if errUnmarshal != nil {
				glogger.GLogger.Error(errUnmarshal)
				continue
			}
			if CacheIoTProperty.Name == "create_at" || CacheIoTProperty.Name == "id" {
				continue
			}
			errHold := CacheIoTProperty.HoldValidator()
			if errHold != nil {
				glogger.GLogger.Error(errUnmarshal)
				continue
			}
			intercache.SetValue("__DataSchema", MIotProperty.Name, intercache.CacheValue{
				UUID:          MIotProperty.UUID,
				Status:        0,
				LastFetchTime: uint64(time.Now().UnixMilli()),
				Value:         CacheIoTProperty,
			})

		}

	}
}

/*
*
* 释放资源
*
 */
func FlushDataSchemaCache() {
	intercache.UnRegisterSlot("__DataSchema")
}

/*
*
* 更新值
*
 */
func UpdateDataSchemaCache(schemaUUID, PropertyUUID string, CachedIoTProperty IoTProperty) {
	Slot := intercache.GetSlot(schemaUUID)
	if Slot != nil {
		Slot[PropertyUUID] = intercache.CacheValue{
			UUID:          PropertyUUID,
			Status:        0,
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         CachedIoTProperty,
		}
	}
}

/*
*
* 获取缓存
*
 */
func GetDataSchemaCache(PropertyName string) (IoTProperty, bool) {
	Slot := intercache.GetSlot("__DataSchema")
	if Slot != nil {
		V, Ok := Slot[PropertyName]
		if Ok {
			if V.Value != nil {
				switch T := V.Value.(type) {
				case IoTProperty:
					return T, Ok
				case *IoTProperty:
					return *T, Ok
				}
			}
		}
	}
	return IoTProperty{}, false
}
