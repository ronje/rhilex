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
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

func TestPerformance() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	availableMemory := memStats.Sys / (1024 * 1024) // MB
	totalMemory := getTotalMemory()                 // MB
	cpuFrequency := getCPUFrequency()               // MHz
	diskSpace := getDiskSpace()                     // GB
	V := calculateV(availableMemory, totalMemory, cpuFrequency, diskSpace)
	printVReferenceTable(V, getLevel(V))
}

// getTotalMemory returns the total amount of memory in bytes.
func getTotalMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.TotalAlloc
}

// getCPUFrequency returns the current CPU frequency in MHz.
func getCPUFrequency() uint64 {
	freq, _ := cpu.Info()
	if len(freq) > 0 {
		return uint64(freq[0].Mhz)
	}
	return 0
}

// getDiskSpace returns the total disk space in bytes.
func getDiskSpace() uint64 {
	// Here we use "/" for the root directory, you can change it to the specific path you want.
	partitions, _ := disk.Partitions(true)
	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		return usage.Total
	}
	return 0
}

// 计算指标 V
// \[
// V = 0.3 \times \left( \frac{\text{可用内存} - 70}{\text{推荐可用内存} - 70} \times 100 \right) +
// 0.3 \times \left( \frac{\text{总内存} - 128}{\text{推荐总内存} - 128} \times 100 \right) +
// 0.2 \times \left( \frac{\text{CPU频率} - 500}{\text{推荐CPU频率} - 500} \times 100 \right) +
// 0.2 \times \left( \frac{\text{磁盘空间} - 8}{\text{推荐磁盘空间} - 8} \times 100 \right)
// \]
func calculateV(availableMemory, totalMemory, cpuFrequency, diskSpace uint64) float64 {
	var score float64

	score += 0.3 * float64(max(0, int(availableMemory)-70)) / float64(70) * 100
	score += 0.3 * float64(max(0, int(totalMemory)-128)) / float64(128) * 100
	score += 0.2 * float64(max(0, int(cpuFrequency)-500)) / float64(500) * 100
	score += 0.2 * float64(max(0, int(diskSpace)-8)) / float64(8) * 100

	return score
}

// 获取等级
func getLevel(V float64) int {
	switch {
	case V < 50:
		return 1
	case V < 70:
		return 2
	case V < 85:
		return 3
	case V < 96:
		return 4
	default:
		return 5
	}
}

// 获取最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 输出计算指标 V 的参考表格
func printVReferenceTable(V1 float64, V2 int) {
	fmt.Println("# Performance Bench Test Result")
	fmt.Printf("- Score: ** %2f **\n- Level: ** %d **\n", V1, V2)
	fmt.Println(`### Grade Level Reference
- 【Level 1】: (Lowest) : Total score 0-49
- 【Level 2】: (Common1): Total score 50-69
- 【Level 3】: (Common2): Total score 70-84
- 【Level 4】: (Common3): Total score 85-95
- 【Level 5】: (Highest): Total score 96-100`)
}
