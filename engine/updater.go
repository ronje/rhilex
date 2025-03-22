// Copyright (C) 2025 wwhai
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

package engine

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hootrhino/rhilex/typex"
)

// 服务器版本号的 URL
const versionURL = "http://127.0.0.1/versions"

// 比较两个版本号的大小
// 返回 true 表示服务器版本号大于本地版本号，否则返回 false
// 如果版本号格式不正确，则返回 false
// 版本号格式为 x.y.z，其中 x、y、z 为整数
func compareVersions(local, server string) bool {
	localParts := strings.Split(local, ".")
	serverParts := strings.Split(server, ".")

	for i := 0; i < len(localParts) && i < len(serverParts); i++ {
		localNum := parseInt(localParts[i])
		serverNum := parseInt(serverParts[i])

		if serverNum > localNum {
			return true
		} else if serverNum < localNum {
			return false
		}
	}

	// 如果前面的部分都相同，比较长度
	return len(serverParts) > len(localParts)
}

// 将字符串转换为整数，如果转换失败则返回 0
func parseInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

// 检查最新版本
func CheckNewestVersion() (string, error) {
	resp, err := http.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("failed to get version information: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	serverVersion := strings.TrimSpace(string(body))
	if compareVersions(typex.MainVersion, serverVersion) {
		return serverVersion, nil
	}
	return typex.MainVersion, fmt.Errorf("Current version is up-to-date:%s", typex.MainVersion)
}
