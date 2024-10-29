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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RhilexConfig

// Init config, First to run!
func InitGlobalConfig(path string) typex.RhilexConfig {
	log.Println("Init rhilex config:", path)
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		log.Fatalf("Fail to read config file: %v", err)
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
		AppDebugMode:          false,
		LogLevel:              "info",
		LogPath:               "rhilex-running-log",
		LogMaxSize:            5,     // MB
		LogMaxBackups:         5,     // Per
		LogMaxAge:             7,     // days
		LogCompress:           true,  // Compress
		MaxKvStoreSize:        1024,  // 20MB
		MaxLostCacheSize:      10000, // 10000 lines
		ExtLibs:               []string{},
		DataSchemaSecret:      []string{"rhilex-secret"},
	}
	if err := cfg.Section("app").MapTo(&GlobalConfig); err != nil {
		log.Fatalf("Fail to map config file: %v", err)
		os.Exit(1)
	}
	log.Println("rhilex config init successfully")
	return GlobalConfig
}

/*
*
* 设置go的线程，通常=0 不需要配置
*
 */
func SetGomaxProcs(GomaxProcs int) {
	if GomaxProcs > 0 {
		if GlobalConfig.GomaxProcs < runtime.NumCPU() {
			runtime.GOMAXPROCS(GlobalConfig.GomaxProcs)
		}
	}
}

/*
*
* 设置性能，通常用来Debug用，生产环境建议关闭
*
 */
func SetDebugMode(EnablePProf bool) {

	//------------------------------------------------------
	// pprof: https://segmentfault.com/a/1190000016412013
	//------------------------------------------------------
	if EnablePProf {
		log.Println("Start PProf debug at: 0.0.0.0:6060")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:6060", nil)
	}
	if EnablePProf {
		go func() {
			readyDebug := false
			for {
				select {
				case <-context.Background().Done():
					{
						glogger.GLogger.Info("PProf exited")
						return
					}
				default:
					{
						time.Sleep(3 * time.Second)
						if !readyDebug {
							fmt.Printf("HeapObjects,\tHeapAlloc,\tTotalAlloc,\tHeapSys")
							fmt.Printf(",\tHeapIdle,\tHeapReleased,\tHeapIdle-HeapReleased")
							fmt.Println()
						}
						readyDebug = true
						TraceMemStats()
					}
				}
			}

		}()

	}
}

/*
*
* DEBUG使用
*
 */
func TraceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	var info [7]float64
	info[0] = float64(ms.HeapObjects)
	info[1] = BtoMB(ms.HeapAlloc)
	info[2] = BtoMB(ms.TotalAlloc)
	info[3] = BtoMB(ms.HeapSys)
	info[4] = BtoMB(ms.HeapIdle)
	info[5] = BtoMB(ms.HeapReleased)
	info[6] = BtoMB(ms.HeapIdle - ms.HeapReleased)

	for _, v := range info {
		fmt.Printf("%v,\t", v)
	}
	fmt.Println()
}
func BtoMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

/*
*
* Byte to Mbyte
*
 */
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
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
