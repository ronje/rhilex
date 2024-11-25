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

package common

/*
*
* 通用的含有主机:端口的这类配置
*
 */
type HostConfig struct {
	Host    string `json:"host" validate:"required" title:"服务地址"`
	Port    int    `json:"port" validate:"required" title:"服务端口"`
	Timeout int    `json:"timeout,omitempty" title:"连接超时"`
}

func (c *HostConfig) Validate() error {
	return nil
}
