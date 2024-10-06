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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/component/uartctrl"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"go.bug.st/serial"
	"gorm.io/gorm"
)

type UartConfigDto struct {
	Timeout  int
	Uart     string
	BaudRate int
	DataBits int
	Parity   string
	StopBits int
}
type UartDto struct {
	UUID     string
	Name     string        // 接口名称
	Type     string        // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias    string        // 别名
	Busy     bool          // 是否被占
	OccupyBy string        // 被谁占用了
	Config   UartConfigDto // 配置, 串口配置、或者网卡、USB等
}

func (u UartConfigDto) JsonString() string {
	if bytes, err := json.Marshal(u); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

/*
*
* 所有的接口列表配置
*
 */
func AllUart() ([]model.MUart, error) {
	ports := []model.MUart{}
	return ports, interdb.DB().
		Model(&model.MUart{}).Find(&ports).Error
}

/*
*
* 配置WIFI Uart
*
 */
func UpdateUartConfig(MUart model.MUart) error {
	Model := model.MUart{}
	return interdb.DB().
		Model(Model).
		Where("uuid=?", MUart.UUID).
		Updates(MUart).Error
}

/*
*
* 获取Uart的配置信息
*
 */
func GetUartConfig(uuid string) (model.MUart, error) {
	MUart := model.MUart{}
	err := interdb.DB().
		Where("uuid=?", uuid).
		Find(&MUart).Error
	return MUart, err
}

/*
*
* 扫描
*
 */
func ReScanUartConfig() error {
	ClearNullPort()
	for _, portName := range GetOsPort() {
		count := int64(-1)
		interdb.DB().Model(model.MUart{}).Where("name=?", portName).Count(&count)
		if count > 0 {
			continue
		}
		NewPort := model.MUart{
			UUID: portName,
			Name: portName,
			Type: "UART",
			Alias: func() string {
				if portName == "/dev/ttyS1" {
					return "RS458-1"
				}
				if portName == "/dev/ttyS2" {
					return "RS458-2"
				}
				return portName
			}(),
			Description: func() string {
				// Alias Ext
				if typex.DefaultVersionInfo.Product == "RHILEXG1" {
					if portName == "/dev/ttyS1" {
						return "RS458-1"
					}
					if portName == "/dev/ttyS2" {
						return "RS458-2"
					}
				}
				return portName
			}(),
		}
		uartCfg := UartConfigDto{
			Timeout:  3000,
			Uart:     portName,
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		}
		NewPort.Config = uartCfg.JsonString()
		err1 := interdb.DB().Model(NewPort).
			Where("name", portName).
			FirstOrCreate(&NewPort).Error
		if err1 != nil {
			return err1
		}
		uartctrl.SetUart(uartctrl.SystemUart{
			UUID: portName,
			Name: portName,
			Type: "UART",
			Alias: func() string {
				if portName == "/dev/ttyS1" {
					return "RS458-1"
				}
				if portName == "/dev/ttyS2" {
					return "RS458-2"
				}
				return portName
			}(),
			Config: uartctrl.UartConfig{
				Timeout:  3000,
				Uart:     portName,
				BaudRate: 9600,
				DataBits: 8,
				Parity:   "N",
				StopBits: 1,
			},
			Description: func() string {
				// Alias Ext
				if typex.DefaultVersionInfo.Product == "RHILEXG1" {
					if portName == "/dev/ttyS1" {
						return "RS458-1"
					}
					if portName == "/dev/ttyS2" {
						return "RS458-2"
					}
				}
				return portName
			}(),
		})
	}
	return nil
}

/*
*
* 清除内存里面没用的串口
*
 */
func ClearNullPort() {
	DbPorts, _ := AllUart()
	OsPorts := GetOsPort()
	TotalPort := []string{}
	for _, DbPort := range DbPorts {
		TotalPort = append(TotalPort, DbPort.Name)
	}
	// 清除缓存里面的数据
	for _, portName := range complement(TotalPort, OsPorts) {
		uartctrl.RemovePort(portName)
		interdb.DB().Model(model.MUart{}).Where("name=?", portName).Delete(model.MUart{})
	}
}

// R = A U B
// 遍历a，找出不在b中的元素
func complement(a, b []string) []string {
	bMap := make(map[interface{}]bool)
	for _, item := range b {
		bMap[item] = true
	}
	var result []string
	for _, item := range a {
		if _, found := bMap[item]; !found {
			result = append(result, item)
		}
	}
	return result
}

/*
*
* 重置
*
 */
func ResetUartConfig() error {
	DbTx := interdb.DB()
	errDbTx := DbTx.Transaction(func(tx *gorm.DB) error {
		err0 := tx.Session(&gorm.Session{
			AllowGlobalUpdate: true,
		}).Delete(&model.MUart{}).Error
		if err0 != nil {
			return err0
		}
		return nil
	})
	if errDbTx != nil {
		return errDbTx
	}
	return InitUartConfig()
}

/*
*
* 初始化网卡配置参数
*
 */
func InitUartConfig() error {
	for _, portName := range GetOsPort() {
		Port := model.MUart{
			UUID: portName,
			Name: portName,
			Type: "UART",
			Alias: func() string {
				// Alias Ext
				if typex.DefaultVersionInfo.Product == "RHILEXG1" {
					if portName == "/dev/ttyS1" {
						return "RS458-1"
					}
					if portName == "/dev/ttyS2" {
						return "RS458-2"
					}
				}
				return portName
			}(),
			Description: func() string {
				// Alias Ext
				if typex.DefaultVersionInfo.Product == "RHILEXG1" {
					if portName == "/dev/ttyS1" {
						return "RS458-1"
					}
					if portName == "/dev/ttyS2" {
						return "RS458-2"
					}
				}
				return portName
			}(),
		}
		uartCfg := UartConfigDto{
			Timeout:  3000,
			Uart:     portName,
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		}
		Port.Config = uartCfg.JsonString()
		err1 := interdb.DB().
			Model(Port).Where("uuid", portName).
			FirstOrCreate(&Port).Error
		if err1 != nil {
			return err1
		}
	}
	return nil
}

/*
*
* 获取系统串口, 这个接口比较特殊，当运行在特殊硬件上的时候，某些系统占用的直接不显示
* 这个接口需要兼容各类特殊硬件
 */
func GetOsPort() []string {
	var ports []string
	if runtime.GOOS == "linux" {
		ports, _ = GetPortsList()
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

// GetPortsList: 获取系统中所有可用的串口设备
func GetPortsList() ([]string, error) {
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
