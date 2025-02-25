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

package target

import (
	"fmt"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

// RabbitMQConfig 定义 RabbitMQ 配置结构体
type RabbitMQConfig struct {
	Address   string `json:"address"`
	Exchange  string `json:"exchange"`
	QueueName string `json:"queueName"`
}

// RabbitMqTarget 实现 XSource 接口的 RabbitMQ 客户端
type RabbitMqTarget struct {
	typex.XStatus
	mainConfig RabbitMQConfig
	isActive   bool
}

func NewRabbitMqTarget(e typex.Rhilex) typex.XTarget {
	ht := new(RabbitMqTarget)
	ht.RuleEngine = e
	ht.mainConfig = RabbitMQConfig{
		Address:   "amqp://rhilex:rhilex@localhost:5672/",
		Exchange:  "rhilex_exchange",
		QueueName: "rhilex_queue",
	}
	ht.SourceState = typex.SOURCE_DOWN
	return ht
}

// Init 实现 Init 方法
func (r *RabbitMqTarget) Init(inEndId string, configMap map[string]any) error {
	r.PointId = inEndId

	return nil
}

// Start 实现 Start 方法
func (r *RabbitMqTarget) Start(ctx typex.CCTX) error {
	if r.isActive {
		glogger.GLogger.Errorf("RabbitMqTarget is already active")
		return fmt.Errorf("RabbitMqTarget is already active")
	}
	r.isActive = true
	glogger.GLogger.Infof("RabbitMqTarget with ID %s started", r.PointId)
	return nil
}
func (r *RabbitMqTarget) Details() *typex.OutEnd {
	return r.RuleEngine.GetOutEnd(r.PointId)
}
func (r *RabbitMqTarget) Status() typex.SourceState {
	return r.SourceState
}
func (r *RabbitMqTarget) To(data any) (any, error) {
	switch T := data.(type) {
	case string:
		glogger.GLogger.Debugf("RabbitMqTarget with ID %s sending message: %s", r.PointId, T)
		return T, nil
	default:
		return nil, fmt.Errorf("email content must plain txt type")
	}
}

// Stop 实现 Stop 方法
func (r *RabbitMqTarget) Stop() {
	r.SourceState = typex.SOURCE_DOWN
	if r.CancelCTX != nil {
		r.CancelCTX()
	}

	glogger.GLogger.Infof("RabbitMqTarget with ID %s stopped", r.PointId)

}
