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

package interqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/luaexecutor"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultXQueue *XQueue

/*
*
* NewXQueue
*
 */

/*
*
* XQueue
*
 */
type XQueue struct {
	Queue        []chan QueueData
	OutQueue     []chan QueueData
	InQueue      []chan QueueData
	DeviceQueue  []chan QueueData
	rhilex       typex.Rhilex
	locker       sync.Mutex
	maxQueueSize int
}

func InitXQueue(rhilex typex.Rhilex, maxQueueSize int) *XQueue {
	InQueue := make([]chan QueueData, 10)
	OutQueue := make([]chan QueueData, 10)
	DeviceQueue := make([]chan QueueData, 10)
	for i := 0; i < 10; i++ {
		InQueue[i] = make(chan QueueData, maxQueueSize)
		OutQueue[i] = make(chan QueueData, maxQueueSize)
		DeviceQueue[i] = make(chan QueueData, maxQueueSize)
	}
	DefaultXQueue = &XQueue{
		InQueue:      InQueue,
		OutQueue:     OutQueue,
		DeviceQueue:  DeviceQueue,
		rhilex:       rhilex,
		locker:       sync.Mutex{},
		maxQueueSize: maxQueueSize,
	}
	return DefaultXQueue
}

/*
*
* 内部队列
*
 */
func StartXQueue() {
	// InQueue
	go func(ctx context.Context, DefaultXQueue *XQueue) {
		glogger.GLogger.Info("Start XQueue: InQueue")
		for {
			for _, Queue := range DefaultXQueue.InQueue {
				select {
				case <-ctx.Done():
					return
				case Data := <-Queue:
					if Data.I == nil {
						continue
					}
					if Data.E == nil {
						continue
					}
					luaexecutor.RunSourceCallbacks(Data.I, Data.Data)
				case <-time.After(4 * time.Millisecond):
					continue
				}
			}
		}
	}(typex.GCTX, DefaultXQueue)
	// DeviceQueue
	go func(ctx context.Context, DefaultXQueue *XQueue) {
		glogger.GLogger.Info("Start XQueue: DeviceQueue")
		for {
			for _, Queue := range DefaultXQueue.DeviceQueue {
				select {
				case <-ctx.Done():
					return
				case Data := <-Queue:
					if Data.D == nil {
						continue
					}
					if Data.E == nil {
						continue
					}
					luaexecutor.RunDeviceCallbacks(Data.D, Data.Data)
				case <-time.After(4 * time.Millisecond):
					continue
				}
			}
		}
	}(typex.GCTX, DefaultXQueue)
	// OutQueue
	go func(ctx context.Context, DefaultXQueue *XQueue) {
		glogger.GLogger.Info("Start XQueue: OutQueue")
		for {
			for _, Queue := range DefaultXQueue.OutQueue {
				select {
				case <-ctx.Done():
					return
				case Data := <-Queue:
					ProcessOutQueueData(Data, Data.E)
				case <-time.After(4 * time.Millisecond):
					continue
				}
			}
		}
	}(typex.GCTX, DefaultXQueue)
}

/*
*
*PushInQueue
*
 */
func (q *XQueue) PushInQueue(in *typex.InEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		I:    in,
		O:    nil,
		Data: data,
	}
	return q.pushIn(qd)
}

/*
*
* PushDeviceQueue
*
 */
func (q *XQueue) PushDeviceQueue(Device *typex.Device, data string) error {
	qd := QueueData{
		D:    Device,
		E:    q.rhilex,
		I:    nil,
		O:    nil,
		Data: data,
	}
	return q.pushDevice(qd)
}

/*
*
* PushOutQueue
*
 */
func (q *XQueue) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		D:    nil,
		I:    nil,
		O:    out,
		Data: data,
	}
	return q.pushOut(qd)
}

/*
*
* Push
*
 */
func (q *XQueue) pushIn(d QueueData) error {
	return q.handleQueue(d, q.InQueue)
}

/*
*
* Push
*
 */
func (q *XQueue) pushOut(d QueueData) error {
	return q.handleQueue(d, q.OutQueue)
}

/*
*
* Push
*
 */
func (q *XQueue) pushDevice(d QueueData) error {
	return q.handleQueue(d, q.DeviceQueue)
}

/*
*
* handle
*
 */
func (q *XQueue) handleQueue(qData QueueData, Queue []chan QueueData) error {
	send := false
	for i := 0; i < 10; i++ {
		if len((Queue[i]))+1 > q.maxQueueSize {
			continue
		}
		Queue[i] <- qData
		send = true
		break
	}
	if send {
		return nil
	}
	return fmt.Errorf("Queue Send Failed")
}
