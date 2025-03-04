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
package licensemanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hootrhino/rhilex/utils"
)

// ReadLicense 从指定文件路径读取并解析许可证
func ReadLocalFileLicense(filePath string) (*utils.LocalLicense, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("[LOAD LICENSE] load license file failed:%s", err)
	}
	// 解析 JSON 数据
	var license utils.LocalLicense
	err = json.Unmarshal(data, &license)
	if err != nil {
		return nil, fmt.Errorf("[LOAD LICENSE] Unmarshal license failed:%s", err)
	}
	return &license, nil
}
