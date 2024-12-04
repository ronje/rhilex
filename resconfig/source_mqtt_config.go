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

import "errors"

type GenericMqttConfig struct {
	Host      string   `json:"host" validate:"required" title:"服务地址"`
	Port      int      `json:"port" validate:"required" title:"服务端口"`
	ClientId  string   `json:"clientId" validate:"required" title:"客户端ID"`
	Username  string   `json:"username" validate:"required" title:"连接账户"`
	Password  string   `json:"password" validate:"required" title:"连接密码"`
	Qos       int      `json:"qos" validate:"required" title:"数据质量"`
	SubTopics []string `json:"subTopics" title:"订阅topic组"`
}

func (cfg *GenericMqttConfig) Validate() error {
	if cfg.Host == "" {
		return errors.New("mqtt config error: host cannot be empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("mqtt config error: port must be a valid number between 1 and 65535")
	}
	if cfg.ClientId == "" {
		return errors.New("mqtt config error: client ID cannot be empty")
	}
	if cfg.Qos < 0 || cfg.Qos > 2 {
		return errors.New("mqtt config error: QoS must be 0, 1, or 2")
	}
	if len(cfg.SubTopics) == 0 {
		return errors.New("mqtt config error: at least one subscription topic is required")
	}
	for _, topic := range cfg.SubTopics {
		if topic == "" {
			return errors.New("mqtt config error: subscription topics cannot be empty")
		}
	}
	return nil
}
