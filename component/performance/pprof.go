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

package performance

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 设置go的线程，通常=0 不需要配置
*
 */
func SetGomaxProcs(GomaxProcs int) {
	if GomaxProcs > 0 {
		if GomaxProcs < runtime.NumCPU() {
			runtime.GOMAXPROCS(GomaxProcs)
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
		log.Println("Start PProf debugger at: http://0.0.0.0:60600")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:60600", nil)
	}
	if EnablePProf {
		go func() {
			readyDebug := false
			for {
				select {
				case <-typex.GCTX.Done():
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
