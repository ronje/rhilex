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

package resconfig

type GenericMqttConfig struct {
	Host      string   `json:"host" validate:"required" title:"服务地址"`
	Port      int      `json:"port" validate:"required" title:"服务端口"`
	ClientId  string   `json:"clientId" validate:"required" title:"客户端ID"`
	Username  string   `json:"username" validate:"required" title:"连接账户"`
	Password  string   `json:"password" validate:"required" title:"连接密码"`
	Qos       int      `json:"qos" validate:"required" title:"数据质量"`
	SubTopics []string `json:"subTopics" title:"订阅topic组"`
}

func (c *GenericMqttConfig) Validate() error {
	return nil
}
