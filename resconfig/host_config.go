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
* 通用的含有主机:端口的这类配置
*
 */
type HostConfig struct {
	Host    string `json:"host" validate:"required"`
	Port    int    `json:"port" validate:"required"`
	Timeout int    `json:"timeout" validate:"required"`
}

// Validate checks the HostConfig for valid values.
func (hc *HostConfig) Validate() error {
	if hc.Host == "" {
		return errors.New("host config error: host cannot be empty")
	}
	if net.ParseIP(hc.Host) == nil {
		return errors.New("host config error: host must be a valid IP address")
	}
	if hc.Port <= 0 || hc.Port > 65535 {
		return errors.New("host config error: port must be a valid number between 1 and 65535")
	}
	if hc.Timeout < 0 {
		return errors.New("host config error: timeout must be a non-negative number")
	}
	return nil
}
