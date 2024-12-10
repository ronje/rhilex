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
package eventbus

import (
	"context"
	"fmt"
	"sync"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var __DefaultEventBus *EventBus

type EventMessage struct {
	Payload string
}

func (E EventMessage) String() string {
	return fmt.Sprintf("Event Message@ Payload: %s", E.Payload)
}

type Topic struct {
	Topic       string
	channel     chan EventMessage
	ctx         context.Context
	cancel      context.CancelFunc
	Subscribers map[string]*Subscriber
}

type Subscriber struct {
	id       string
	Callback func(Topic string, Msg EventMessage)
}

type EventBus struct {
	// Topic, chan EventMessage
	// 给每个订阅者分配一个 Channel，实现消息订阅
	// Topic 一样的会挂在同一个树上
	Topics sync.Map // 订阅树: MAP<Topic>[]Subscribers
}

func InitEventBus() *EventBus {
	__DefaultEventBus = &EventBus{}
	return __DefaultEventBus
}

func (eb *EventBus) createTopic(topic string) *Topic {
	t := &Topic{
		channel:     make(chan EventMessage, 100),
		Subscribers: make(map[string]*Subscriber),
		Topic:       topic,
	}
	ctx, cancel := typex.NewCCTX()
	t.ctx = ctx
	t.cancel = cancel
	go eb.topicWorker(t)
	return t
}

func (eb *EventBus) topicWorker(t *Topic) {
	for {
		select {
		case <-t.ctx.Done():
			return
		case msg := <-t.channel:
			for _, sub := range t.Subscribers {
				if sub.Callback != nil {
					sub.Callback(t.Topic, msg)
				}
			}
		}
	}
}

func (eb *EventBus) deleteTopic(topic string) {
	if value, ok := eb.Topics.Load(topic); ok {
		t := value.(*Topic)
		t.cancel()
		eb.Topics.Delete(topic)
	}
}

func (eb *EventBus) getTopic(topic string) *Topic {
	if value, ok := eb.Topics.Load(topic); ok {
		return value.(*Topic)
	}
	return nil
}

func (eb *EventBus) ensureTopic(topic string) *Topic {
	t := eb.getTopic(topic)
	if t == nil {
		t = eb.createTopic(topic)
		eb.Topics.Store(topic, t)
	}
	return t
}

func (eb *EventBus) removeSubscriber(topic string, sub *Subscriber) {
	if t := eb.getTopic(topic); t != nil {
		delete(t.Subscribers, sub.id)
		if len(t.Subscribers) == 0 {
			eb.deleteTopic(topic)
		}
	}
}

// Subscribe 订阅
func Subscribe(topic string, sub *Subscriber) {
	NewUUID := utils.MakeUUID("SUB")
	sub.id = NewUUID
	eb := InitEventBus()
	t := eb.ensureTopic(topic)
	t.Subscribers[sub.id] = sub
}

// UnSubscribe 取消订阅
func UnSubscribe(topic string, sub *Subscriber) {
	eb := InitEventBus()
	eb.removeSubscriber(topic, sub)
}

// Publish 发布
func Publish(topic string, msg EventMessage) {
	eb := InitEventBus()
	if t := eb.getTopic(topic); t != nil {
		t.channel <- msg
	}
}

// Flush 释放所有
func Flush() {
	eb := InitEventBus()
	eb.Topics.Range(func(key, value interface{}) bool {
		t := value.(*Topic)
		t.cancel()
		return true
	})
	eb.Topics = sync.Map{}
}
