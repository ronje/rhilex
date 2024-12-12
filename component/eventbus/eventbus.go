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
	"strings"
	"sync"
)

type EventMessage struct {
	Topic   string
	From    string
	Type    string
	Event   string
	Ts      uint64
	Payload interface{}
}

func (E EventMessage) String() string {
	return fmt.Sprintf("Event Message@ Payload: %v", E.Payload)
}

type Subscriber struct {
	id       string                               // 随机生成的ID
	Callback func(topic string, msg EventMessage) // 回调函数
}

type TrieNode struct {
	Subscribers map[string]*Subscriber
	Children    map[string]*TrieNode
	Wildcard    *TrieNode // 用于处理 "*" 通配符
}

type EventBus struct {
	root    *TrieNode
	mutex   sync.RWMutex
	stopped bool
	cancel  context.CancelFunc
	ctx     context.Context
}

func NewEventBus() *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventBus{
		root: &TrieNode{
			Subscribers: make(map[string]*Subscriber),
			Children:    make(map[string]*TrieNode),
		},
		cancel: cancel,
		ctx:    ctx,
	}
}

// Subscribe 订阅某个主题（支持通配符）
func (eb *EventBus) Subscribe(topic string, sub *Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	if eb.stopped {
		return
	}

	nodes := strings.Split(topic, ".")
	node := eb.root
	for _, part := range nodes {
		if part == "*" {
			if node.Wildcard == nil {
				node.Wildcard = &TrieNode{
					Subscribers: make(map[string]*Subscriber),
					Children:    make(map[string]*TrieNode),
				}
			}
			node = node.Wildcard
		} else {
			if _, exists := node.Children[part]; !exists {
				node.Children[part] = &TrieNode{
					Subscribers: make(map[string]*Subscriber),
					Children:    make(map[string]*TrieNode),
				}
			}
			node = node.Children[part]
		}
	}

	node.Subscribers[sub.id] = sub
}

// UnSubscribe 取消订阅
func (eb *EventBus) UnSubscribe(topic string, sub *Subscriber) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	if eb.stopped {
		return
	}

	nodes := strings.Split(topic, ".")
	eb.unsubscribeRecursive(eb.root, nodes, sub, 0)
}

func (eb *EventBus) unsubscribeRecursive(node *TrieNode, nodes []string, sub *Subscriber, index int) bool {
	if index == len(nodes) {
		delete(node.Subscribers, sub.id)
		return len(node.Subscribers) == 0 && len(node.Children) == 0 && node.Wildcard == nil
	}

	part := nodes[index]
	if part == "*" {
		if node.Wildcard != nil && eb.unsubscribeRecursive(node.Wildcard, nodes, sub, index+1) {
			node.Wildcard = nil
		}
	} else {
		if child, exists := node.Children[part]; exists {
			if eb.unsubscribeRecursive(child, nodes, sub, index+1) {
				delete(node.Children, part)
			}
		}
	}

	return len(node.Subscribers) == 0 && len(node.Children) == 0 && node.Wildcard == nil
}

// Publish 发布消息
func (eb *EventBus) Publish(topic string, msg EventMessage) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()
	if eb.stopped {
		return
	}

	nodes := strings.Split(topic, ".")
	eb.publishRecursive(eb.root, nodes, msg, 0)
}

func (eb *EventBus) publishRecursive(node *TrieNode, nodes []string, msg EventMessage, index int) {
	// 调用所有订阅者
	for _, sub := range node.Subscribers {
		go sub.Callback(msg.Topic, msg)
	}

	if index == len(nodes) {
		return
	}

	part := nodes[index]
	if child, exists := node.Children[part]; exists {
		eb.publishRecursive(child, nodes, msg, index+1)
	}

	if node.Wildcard != nil {
		eb.publishRecursive(node.Wildcard, nodes, msg, index+1)
	}
}

// Stop 停止 EventBus，并释放所有资源
func (eb *EventBus) Stop() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	if eb.stopped {
		return
	}

	eb.stopped = true
	eb.cancel()          // 取消上下文
	eb.root = &TrieNode{ // 清空订阅树
		Subscribers: make(map[string]*Subscriber),
		Children:    make(map[string]*TrieNode),
	}
}
