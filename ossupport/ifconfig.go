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
	"bytes"
	"errors"
	"os/exec"
	"runtime"
)

// Ifconfig 执行系统的 ifconfig 或 ipconfig 命令，并返回输出结果
func Ifconfig() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ipconfig", "/all")
	} else {
		cmd = exec.Command("ifconfig", "-a")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	result := out.String()
	if result == "" {
		return "", errors.New("no output from ifconfig/ipconfig command")
	}

	return result, nil
}
