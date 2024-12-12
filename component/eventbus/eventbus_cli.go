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

package eventbus

import "github.com/hootrhino/rhilex/typex"

var __DefaultEventBus *EventBus

func InitEventBus(r typex.Rhilex) {
	__DefaultEventBus = NewEventBus()
}

// Subscribe 订阅
func Subscribe(topic string, sub *Subscriber) {
	__DefaultEventBus.Subscribe(topic, sub)
}

// UnSubscribe 取消订阅
func UnSubscribe(topic string, sub *Subscriber) {
	__DefaultEventBus.UnSubscribe(topic, sub)
}

// Publish 发布
func Publish(topic string, msg EventMessage) {
	__DefaultEventBus.Publish(topic, msg)
}

// Flush 释放所有
func Stop() {
	__DefaultEventBus.Stop()
}
