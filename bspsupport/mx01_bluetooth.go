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

package archsupport

import (
	"fmt"
	"os"
	"time"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "RHILEXG1" {
		fmt.Println("Init Mx01 Bluetooth")
		__EC200AInitMx01Bluetooth()
		fmt.Println("Init Mx01 Bluetooth Ok.")
	}
}
func __EC200AInitMx01Bluetooth() {

}

/*
*
* MX-01模块基本信息
*
 */
type Mx01BaseInfo struct {
	Mac             string //
	Name            string //
	BaudRate        int    //
	BroadCastStatus int    //
}

/*
*
* 获取模块信息
*
 */
func GetMx01BaseInfo() (Mx01BaseInfo, error) {
	Mx01BaseInfo := Mx01BaseInfo{}
	{
		// → AT+MAC?
		// ← +MAC:FF2310184034
		Mac, err := __Mx01_AT("AT+MAC?\r\n", 100)
		if err != nil {
			return Mx01BaseInfo, err
		}
		Mx01BaseInfo.Mac = Mac[6:]
	}
	{
		// → AT+NAME?
		// ← +NAME:FF2310184034
		Name, err := __Mx01_AT("AT+NAME?\r\n", 100)
		if err != nil {
			return Mx01BaseInfo, err
		}
		Mx01BaseInfo.Mac = Name[7:]
	}
	{
		// → AT+UART?
		// ← +UART:0
		// 0:9600; 1:14400; 2:19200; 3:38400; 4:57600; 5:115200;
		BaudRate, err := __Mx01_AT("AT+UART?\r\n", 100)
		if err != nil {
			return Mx01BaseInfo, err
		}
		Mx01BaseInfo.Mac = BaudRate[7:]
	}
	return Mx01BaseInfo, nil
}

/*
*
* MX-01蓝牙指令集
*
 */
func __Mx01_AT(command string, timeout time.Duration) (string, error) {

	return string(""), nil
}
