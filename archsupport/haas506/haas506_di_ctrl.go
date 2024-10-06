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

const HAAS506_DI_SYSDEV_PATH = "/sys/class/gpio/gpio%d/value"

//-----------------------------------------------
// 这是HAAS506网关的DI-DI支持库
//-----------------------------------------------
/*
/sys/class/gpio/gpio56
/sys/class/gpio/gpio57
/sys/class/gpio/gpio58
/sys/class/gpio/gpio59
*/
const (
	// DI
	HAAS506_DI1 string = "56"
	HAAS506_DI2 string = "57"
	HAAS506_DI3 string = "58"
	HAAS506_DI4 string = "59"
)

func _HAAS506_DI_Init() int {
	gpio56 := "/sys/class/gpio/gpio56/value"
	gpio57 := "/sys/class/gpio/gpio57/value"
	gpio58 := "/sys/class/gpio/gpio58/value"
	gpio59 := "/sys/class/gpio/gpio59/value"
	_, err1 := os.Stat(gpio56)
	_, err2 := os.Stat(gpio57)
	_, err3 := os.Stat(gpio58)
	_, err4 := os.Stat(gpio59)
	if err1 != nil {
		if strings.Contains(err1.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DI1, HAAS506_In)
			fmt.Println("HAAS506_GPIOAllInit DI1 Out Mode Ok")
		}
	}
	if err2 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DI2, HAAS506_In)
			fmt.Println("HAAS506_GPIOAllInit DI2 Out Mode Ok")
		}
	}
	if err3 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DI3, HAAS506_In)
			fmt.Println("HAAS506_GPIOAllInit DI2 Out Mode Ok")
		}
	}
	if err4 != nil {
		if strings.Contains(err2.Error(), "no such file or directory") {
			_HAAS506_GPIOInit(HAAS506_DI4, HAAS506_In)
			fmt.Println("HAAS506_GPIOAllInit DI2 Out Mode Ok")
		}
	}
	return 1
}

/*
*
* 新版本的文件读取形式获取GPIO状态
*
 */
func HAAS506_GPIOGetDI1() (int, error) {
	return HAAS506_GPIOGetByFile(56)
}
func HAAS506_GPIOGetDI2() (int, error) {
	return HAAS506_GPIOGetByFile(57)
}
func HAAS506_GPIOGetDI3() (int, error) {
	return HAAS506_GPIOGetByFile(58)
}
func HAAS506_GPIOGetDI4() (int, error) {
	return HAAS506_GPIOGetByFile(59)
}
