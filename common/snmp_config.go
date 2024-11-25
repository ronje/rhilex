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
* SNMP 配置
*
 */
type GenericSnmpConfig struct {
	// Target is an ipv4 address.
	Target string `json:"target" validate:"required" title:"Target" info:"Target"`
	// Port is a port.
	Port uint16 `json:"port" validate:"required" title:"Port" info:"Port"`
	// Transport is the transport protocol to use ("udp" or "tcp"); if unset "udp" will be used.
	Transport string `json:"transport" validate:"required" title:"Transport" info:"Transport"`
	// Community is an SNMP Community string.
	Community string `json:"community" validate:"required" title:"Community" info:"Community"`
	// 1 2 3
	Version uint8 `json:"version" validate:"required" title:"Community" info:"Community"`
}

func (c *GenericSnmpConfig) Validate() error {
	return nil
}
