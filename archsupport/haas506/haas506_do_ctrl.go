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

package haas506

import (
	"fmt"
	"os"
	"strings"
)

const HAAS506_DO_SYSDEV_PATH = "/sys/class/gpio/gpio%d/value"

//-----------------------------------------------
// 这是HAAS506网关的DI-DO支持库
//-----------------------------------------------
/*
/sys/class/gpio/gpio47
/sys/class/gpio/gpio48
/sys/class/gpio/gpio49
/sys/class/gpio/gpio50
*/
const (
	// DO
	HAAS506_DO1 string = "47"
	HAAS506_DO2 string = "48"
	HAAS506_DO3 string = "49"
	HAAS506_DO4 string = "50"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		_HAAS506_DO_Init()
	}
}

func _HAAS506_DO_Init() int {
	gpio47 := "/sys/class/gpio/gpio47/value"
	gpio48 := "/sys/class/gpio/gpio48/value"
	gpio49 := "/sys/class/gpio/gpio49/value"
	gpio50 := "/sys/class/gpio/gpio50/value"
	_, err1 := os.Stat(gpio47)
	_, err2 := os.Stat(gpio48)
	_, err3 := os.Stat(gpio49)
	_, err4 := os.Stat(gpio50)
	if err1 != nil {
		if strings.Contains(err1.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DO1, HAAS506_Out)
			fmt.Println("HAAS506_GPIOAllInit DO1 Out Mode Ok")
		}
	}
	if err2 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DO2, HAAS506_Out)
			fmt.Println("HAAS506_GPIOAllInit DO2 Out Mode Ok")
		}
	}
	if err3 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DO3, HAAS506_Out)
			fmt.Println("HAAS506_GPIOAllInit DO2 Out Mode Ok")
		}
	}
	if err4 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DO4, HAAS506_Out)
			fmt.Println("HAAS506_GPIOAllInit DO2 Out Mode Ok")
		}
	}
	return 1
}

/*
*
* 新版本的文件读取形式获取GPIO状态
*
 */
func HAAS506_GPIOGetDO1() (int, error) {
	return HAAS506_GPIOGetByFile(47)
}
func HAAS506_GPIOGetDO2() (int, error) {
	return HAAS506_GPIOGetByFile(48)
}
func HAAS506_GPIOGetDO3() (int, error) {
	return HAAS506_GPIOGetByFile(49)
}
func HAAS506_GPIOGetDO4() (int, error) {
	return HAAS506_GPIOGetByFile(50)
}

// Set

func HAAS506_GPIOSetDO1(value int) error {
	return HAAS506_GPIOSetByFile(47, value)
}
func HAAS506_GPIOSetDO2(value int) error {
	return HAAS506_GPIOSetByFile(48, value)
}
func HAAS506_GPIOSetDO3(value int) error {
	return HAAS506_GPIOSetByFile(49, value)
}
func HAAS506_GPIOSetDO4(value int) error {
	return HAAS506_GPIOSetByFile(50, value)
}
