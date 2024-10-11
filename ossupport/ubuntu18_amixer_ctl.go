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

package ossupport

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
)

/*
* amixer 设置音量, 输入参数是个数值, 每次增加或者减少1%
*        amixer set 'Line Out' 1 | grep 'Front Left:' | awk -F '[][]' '{print $2}'
*
 */
func SetVolume(v int) (string, error) {
	shellCmd := "amixer set 'Line Out' %s | grep 'Front Left:' | awk -F '[][]' '{print $2}'"
	if v > -100 && v < 100 {
		var cmd *exec.Cmd
		if v < 0 {
			cmd = exec.Command("sh", "-c", fmt.Sprintf(shellCmd, fmt.Sprintf("%v%%-", math.Abs(float64(v)))))
		}
		if v > 0 {
			cmd = exec.Command("sh", "-c", fmt.Sprintf(shellCmd, fmt.Sprintf("%v%%+", math.Abs(float64(v)))))
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
		}
		volume := strings.TrimSpace(string(output))
		return volume, nil
	}
	return "", fmt.Errorf("Invalid volume:%v, must be in range [0,100]", v)

}

/*
*
  - 获取音量百分比 20%
    amixer get Master | grep 'Front Left:' | awk -F '[][]' '{print $2}'

*
*/
func GetVolume() (string, error) {
	// 创建一个 Command 对象，执行多个命令通过管道连接
	cmd := exec.Command("sh", "-c",
		"amixer get 'Line Out' | grep 'Front Left:' | awk -F '[][]' '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	volume := strings.TrimSpace(string(output))
	return volume, nil
}
