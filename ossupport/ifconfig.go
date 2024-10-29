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
	"io"
	"os/exec"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	if runtime.GOOS == "windows" {
		// 将 GBK 转换为 UTF-8
		reader := transform.NewReader(
			bytes.NewReader(out.Bytes()), // 使用 bytes.NewReader
			simplifiedchinese.GBK.NewDecoder(),
		)
		decodedOutput, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return string(decodedOutput), nil

	}
	return out.String(), nil
}
