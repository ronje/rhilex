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

package core

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/typex"

	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RhilexConfig

// Init config, First to run!
func InitGlobalConfig(path string) typex.RhilexConfig {
	log.Println("[RHILEX INIT] Init config:", path)
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		log.Fatalf("[RHILEX INIT] Load config failed: %v Make sure your config path is valid", err)
		os.Exit(1)
	}
	GlobalConfig = typex.RhilexConfig{
		AppId:                 "rhilex",
		IniPath:               path,
		MaxQueueSize:          10240,
		SourceRestartInterval: 5000,
		GomaxProcs:            0,
		EnablePProf:           false,
		EnableConsole:         false,
		DebugMode:             false,
		LogLevel:              "info",
		LogMaxSize:            5,    // MB
		LogMaxBackups:         5,    // Per
		LogMaxAge:             7,    // days
		LogCompress:           true, // Compress
		MaxKvStoreSize:        1024, // 20MB
		ExtLibs:               []string{},
		DataSchemaSecret:      []string{"rhilex-secret"},
	}
	if err := cfg.Section("main").MapTo(&GlobalConfig); err != nil {
		log.Fatalf("[RHILEX INIT] Fail to map config file: %v", err)
		os.Exit(1)
	}
	log.Println("[RHILEX INIT] RHILEX config load successfully:", path)
	return GlobalConfig
}

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
	switch T := Value.Value.(type) {
	case []byte:
		err := json.Unmarshal(T, &configMap)
		if err != nil {
			return nil
		}
	}
	return configMap
}
