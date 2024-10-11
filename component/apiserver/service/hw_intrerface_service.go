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

package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"go.bug.st/serial"
)

/**
 * 兼容获取Windows和Linux系统的串口
 *
 */
func GetOsPort() []string {
	var ports []string
	if runtime.GOOS == "linux" {
		ports, _ = GetLinuxPortsList()
	} else {
		ports, _ = serial.GetPortsList()
	}
	List := []string{}
	for _, port := range ports {
		if typex.DefaultVersionInfo.Product == "RHILEXG1" {
			// RHILEXG1的下列串口被系统占用
			if utils.SContains([]string{
				"/dev/ttyS0",
				"/dev/ttyS3",
				"/dev/ttyS4",   // Linux System
				"/dev/ttyS5",   // Linux System
				"/dev/ttyS6",   // Linux System
				"/dev/ttyS7",   // Linux System
				"/dev/ttyUSB0", // 4G
				"/dev/ttyUSB1", // 4G
				"/dev/ttyUSB2", // 4G
			}, port) {
				continue
			}
		}
		List = append(List, port)
	}
	return List
}

// GetLinuxPortsList: 获取系统中所有可用的串口设备
func GetLinuxPortsList() ([]string, error) {
	var availablePorts []string
	serialFile := "/proc/tty/driver/serial"
	if _, err := os.Stat(serialFile); os.IsNotExist(err) {
		fmt.Printf("%s does not exist, skipping this check.\n", serialFile)
	} else {
		file, err := os.Open(serialFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %v", serialFile, err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) > 1 && strings.HasPrefix(fields[0], "0:") ||
				strings.HasPrefix(fields[0], "1:") || strings.HasPrefix(fields[0], "2:") {
				index := strings.TrimSuffix(fields[0], ":")
				uartType := fields[1]
				if uartType != "unknown" {
					device := fmt.Sprintf("/dev/ttyS%s", index)
					availablePorts = append(availablePorts, device)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
	}
	devDir := "/dev/"
	files, err := os.ReadDir(devDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s directory: %v", devDir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		// 检查是否是 ttyS*、ttyUSB* 或 tty485_*
		if strings.HasPrefix(name, "ttyS") ||
			strings.HasPrefix(name, "ttyUSB") ||
			strings.HasPrefix(name, "tty232") ||
			strings.HasPrefix(name, "tty485") {
			availablePorts = append(availablePorts, filepath.Join(devDir, name))
		}
	}
	return availablePorts, nil
}
