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

package core

import (
	"encoding/json"
	"fmt"

	"github.com/hootrhino/rhilex/component/intercache"
)

/*
*
* 从全局缓存器获取设备的配置
*
 */
func GetDeviceConfigMap(deviceUuid string) map[string]interface{} {
	Slot := intercache.GetSlot("__DeviceConfigMap")
	Value, ok := Slot[deviceUuid]
	if !ok {
		return nil
	}
	configMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(Value.Value), &configMap)
	if err != nil {
		return nil
	}
	return configMap
}
