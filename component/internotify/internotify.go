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

package internotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var __DefaultInternalEventBus *InternalEventBus

// ---------------------------------------------------------
// Type
// ---------------------------------------------------------
// - SOURCE: 南向事件
// - DEVICE: 设备事件
// - TARGET: 北向事件
// - SYSTEM: 系统内部事件
// - HARDWARE: 硬件事件

type BaseEvent struct {
	Type    string
	Event   string
	Ts      uint64
	Summary string
	Info    interface{}
}

func (be BaseEvent) String() string {
	return fmt.Sprintf(`Event: [%s], [%s], %s`, be.Type, be.Event, be.Info)
}

/*
*
* Push
*
 */
func Push(e BaseEvent) error {
	if len(__DefaultInternalEventBus.Queue)+1 > __DefaultInternalEventBus.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			__DefaultInternalEventBus.GetSize(), len(__DefaultInternalEventBus.Queue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		__DefaultInternalEventBus.Queue <- e
		return nil
	}
}

/*
*
* 内部事件总线
*
 */
type InternalEventBus struct {
	Queue       chan BaseEvent
	rhilex      typex.Rhilex
	Subscribers map[string]chan BaseEvent
}

func (q *InternalEventBus) GetSize() int {
	return cap(q.Queue)
}

/*
*
* 取消外部订阅
*
 */
func RemoveSubscriber(name string) {
	if Channel, Ok := __DefaultInternalEventBus.Subscribers[name]; Ok {
		if _, open := <-Channel; open {
			close(Channel) // 一定要记住关闭这个channel
		}
		delete(__DefaultInternalEventBus.Subscribers, name)
	}
}

/*
*
* 加入一个外部订阅者
*
 */
func AddSubscriber(name string, channel chan BaseEvent) {
	if Channel, Ok := __DefaultInternalEventBus.Subscribers[name]; !Ok {
		__DefaultInternalEventBus.Subscribers[name] = Channel
	}
}
func GetQueue() chan BaseEvent {
	return __DefaultInternalEventBus.Queue
}

/*
*
  - 内部事件，例如资源挂了或者设备离线、超时等等,该资源是单例模式,
    维护一个channel来接收各种事件，将收到的消息吐到InterQueue即可

*
*/
func InitInternalEventBus(r typex.Rhilex, MaxQueueSize int) *InternalEventBus {
	__DefaultInternalEventBus = new(InternalEventBus)
	__DefaultInternalEventBus.Queue = make(chan BaseEvent, 1024)
	__DefaultInternalEventBus.Subscribers = map[string]chan BaseEvent{}
	__DefaultInternalEventBus.rhilex = r
	return __DefaultInternalEventBus
}

/*
*
* 监控chan
*
 */
type MInternalNotify struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UUID      string    `gorm:"not null"` // UUID
	Type      string    `gorm:"not null"` // INFO | ERROR | WARNING
	Status    int       `gorm:"not null"` // 1 未读 2 已读
	Event     string    `gorm:"not null"` // 字符串
	Ts        uint64    `gorm:"not null"` // 时间戳
	Summary   string    `gorm:"not null"` // 概览，为了节省流量，在消息列表只显示这个字段，Info值为“”
	Info      string    `gorm:"not null"` // 消息内容，是个文本，详情显示
}

func StartInternalEventQueue() {
	go func(ctx context.Context, InternalEventBus *InternalEventBus) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				continue // 防止超时死锁
			case Event := <-InternalEventBus.Queue:
				// 将消息推给订阅者
				for _, Channel := range InternalEventBus.Subscribers {
					Channel <- Event
				}
				// glogger.GLogger.Debug("Internal Event:", Event)
				interdb.DB().Table("m_internal_notifies").Save(&MInternalNotify{
					UUID:    utils.MakeUUID("NOTIFY"),
					Type:    Event.Type, // INFO | ERROR | WARNING
					Status:  1,          // Default unread
					Event:   Event.Event,
					Ts:      Event.Ts,
					Summary: "RHILEX Internal Event: " + Event.Event,
					Info:    Event.String(),
				})
			}
		}
	}(typex.GCTX, __DefaultInternalEventBus)
}
