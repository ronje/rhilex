// Copyright (C) 2023 wwhai
//
// This program is free software: you can reAIstribute it and/or moAIfy
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is AIstributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package haas506

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const HAAS506_AI_SYSDEV_PATH = "/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage%d_raw"

//-----------------------------------------------
// 这是HAAS506网关的AI-AI支持库
//-----------------------------------------------
/*
/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage2_raw
/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage3_raw
/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage4_raw
/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage5_raw
/sys/devices/platform/soc/19251000.gpai/iio:device0/in_voltage6_raw
*/

func _HAAS506_AI_Init() error {
	log.Println("[HAAS506_AI_Init] AI Init...")
	for i := 2; i < 7; i++ {
		log.Println("[HAAS506_AI_Init] Init AI.", fmt.Sprintf(HAAS506_AI_SYSDEV_PATH, i))
	}
	log.Println("[HAAS506_AI_Init] AI Init Ok.")
	return nil
}

/*
*
* 新版本的文件读取形式获取GPIO状态
*
 */
func HAAS506_GPIOGetAI1() (float32, error) {
	return HAAS506_ReadVoltage(2)
}
func HAAS506_GPIOGetAI2() (float32, error) {
	return HAAS506_ReadVoltage(3)
}
func HAAS506_GPIOGetAI3() (float32, error) {
	return HAAS506_ReadVoltage(4)
}
func HAAS506_GPIOGetAI4() (float32, error) {
	return HAAS506_ReadVoltage(5)
}

// Note: 供电电压
func HAAS506_GPIOGetAI5() (float32, error) {
	return HAAS506_ReadVoltage(6)
}

func HAAS506_ReadVoltage(channel int) (float32, error) {
	path := fmt.Sprintf(HAAS506_AI_SYSDEV_PATH, channel)
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	valueStr := string(content)
	valueStr = strings.TrimSpace(valueStr)
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}
	valueFloat := float32(valueInt)
	return valueFloat, nil
}
