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

import (
	"fmt"
	"testing"
)

// 示例代码
func Test_EventBus_PubSub(t *testing.T) {
	eb := NewEventBus()

	// 创建订阅者
	sub1 := &Subscriber{
		id: "sub1",
		Callback: func(topic string, msg EventMessage) {
			fmt.Printf("Received message on %s: %v\n", topic, msg)
		},
	}

	sub2 := &Subscriber{
		id: "sub2",
		Callback: func(topic string, msg EventMessage) {
			fmt.Printf("Wildcard received message on %s: %v\n", topic, msg)
		},
	}

	// 订阅主题
	eb.Subscribe("a.b.c", sub1)
	eb.Subscribe("a.b.*", sub2)
	eb.Subscribe("a.*.d", sub2)
	eb.Subscribe("*.b.c", sub2)

	// 发布消息
	eb.Publish("a.b.c", EventMessage{Payload: "Hello, Exact Match!"})
	eb.Publish("a.b.d", EventMessage{Payload: "Hello, Wildcard Match!"})
	eb.Publish("a.x.d", EventMessage{Payload: "Hello, Multi-Level Wildcard Match!"})
	eb.Publish("x.b.c", EventMessage{Payload: "Hello, Leading Wildcard Match!"})
}
