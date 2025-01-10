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

package upgrader

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// 返回最新版本JSON格式
func GetNewVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 从HTTP地址下载Zip压缩包
func FetchNewestPackage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CompareVersion 比较两个版本号
// 如果 version1 是比 version2 更新的版本，则返回 true
func CompareVersion(version1, version2 string) bool {
	v1Parts := strings.SplitN(version1[1:], ".", 3)
	v2Parts := strings.SplitN(version2[1:], ".", 3)

	for i := 0; i < 3; i++ {
		part1, err1 := strconv.Atoi(v1Parts[i])
		part2, err2 := strconv.Atoi(v2Parts[i])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting version parts to integers")
			return false
		}
		if part1 > part2 {
			return true
		} else if part1 < part2 {
			return false
		}
	}
	return false
}
