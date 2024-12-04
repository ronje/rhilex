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

import (
	"errors"
	"net"
)

/*
*
* SNMP 配置
*
 */
type GenericSnmpConfig struct {
	// Target is an ipv4 address.
	Target string `json:"target" validate:"required"`
	// Port is a port.
	Port uint16 `json:"port" validate:"required"`
	// Transport is the transport protocol to use ("udp" or "tcp"); if unset "udp" will be used.
	Transport string `json:"transport" validate:"required"`
	// Community is an SNMP Community string.
	Community string `json:"community" validate:"required"`
	// 1 2 3
	Version uint8 `json:"version" validate:"required"`
}

// Validate checks the GenericSnmpConfig for valid values.
func (cfg *GenericSnmpConfig) Validate() error {
	if cfg.Target == "" {
		return errors.New("snmp config error: target cannot be empty")
	}
	if net.ParseIP(cfg.Target) == nil {
		return errors.New("snmp config error: target must be a valid IPv4 address")
	}
	if cfg.Port == 0 || cfg.Port > uint16(65535) {
		return errors.New("snmp config error: port must be a valid number between 1 and 65535")
	}
	if cfg.Transport != "udp" && cfg.Transport != "tcp" {
		return errors.New("snmp config error: transport must be 'udp' or 'tcp'")
	}
	if cfg.Community == "" {
		return errors.New("snmp config error: community cannot be empty")
	}
	if cfg.Version < 1 || cfg.Version > 3 {
		return errors.New("snmp config error: version must be 1, 2, or 3")
	}
	return nil
}
