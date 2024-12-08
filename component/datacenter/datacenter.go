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

package datacenter

import (
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/typex"
)

var __DefaultDataCenter *DataCenter

/*
*
* 留着未来扩充数据中心的功能
*
 */
type DataCenter struct {
	rhilex  typex.Rhilex
	secrets map[string]bool
}

func InitDataCenter(rhilex typex.Rhilex) {
	__DefaultDataCenter = new(DataCenter)
	__DefaultDataCenter.rhilex = rhilex
	secrets := map[string]bool{}
	for _, v := range core.GlobalConfig.DataSchemaSecret {
		secrets[v] = true
	}
	loadSecrets(secrets)
	go StartClearDataCenterCron()
}
func loadSecrets(secrets map[string]bool) {
	__DefaultDataCenter.secrets = secrets
}
func CheckSecrets(secret string) bool {
	_, ok := __DefaultDataCenter.secrets[secret]
	return ok
}
