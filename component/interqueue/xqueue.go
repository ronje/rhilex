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
	"math/rand"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/luaexecutor"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

var __DefaultXQueue *XQueue

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
	__DefaultXQueue = &XQueue{
		InQueue:      InQueue,
		OutQueue:     OutQueue,
		DeviceQueue:  DeviceQueue,
		rhilex:       rhilex,
		locker:       sync.Mutex{},
		maxQueueSize: maxQueueSize,
	}
	return __DefaultXQueue
}
func StartXQueue() {
	__DefaultXQueue.StartXQueue()
}
func (q *XQueue) startInQueue(ctx context.Context) {
	glogger.GLogger.Info("Start XQueue: InQueue")
	q.processQueue(ctx, q.InQueue, func(data QueueData) {
		if data.I == nil || data.E == nil {
			return
		}
		luaexecutor.RunSourceCallbacks(data.I, data.Data)
	})
}

func (q *XQueue) startDeviceQueue(ctx context.Context) {
	glogger.GLogger.Info("Start XQueue: DeviceQueue")
	q.processQueue(ctx, q.DeviceQueue, func(data QueueData) {
		if data.D == nil || data.E == nil {
			return
		}
		luaexecutor.RunDeviceCallbacks(data.D, data.Data)
	})
}

func (q *XQueue) startOutQueue(ctx context.Context) {
	glogger.GLogger.Info("Start XQueue: OutQueue")
	q.processQueue(ctx, q.OutQueue, func(data QueueData) {
		ProcessOutQueueData(data, data.E)
	})
}

func (q *XQueue) processQueue(ctx context.Context, queue []chan QueueData, callback func(QueueData)) {
	for {
		for _, qc := range queue {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-qc:
				if ok {
					callback(data)
				}
			case <-time.After(4 * time.Millisecond):
				continue
			}
		}
	}
}

func (q *XQueue) StartXQueue() {
	ctx := context.Background()
	go q.startInQueue(ctx)
	go q.startDeviceQueue(ctx)
	go q.startOutQueue(ctx)
}

func (q *XQueue) push(data QueueData, queue []chan QueueData) error {
	q.locker.Lock()
	defer q.locker.Unlock()
	var available []int
	for i := 0; i < len(queue); i++ {
		if len(queue[i])+1 <= q.maxQueueSize {
			available = append(available, i)
		}
	}
	if len(available) > 0 {
		randomIndex := available[rand.Intn(len(available))]
		queue[randomIndex] <- data
		return nil
	}
	return fmt.Errorf("Queue Send Failed: No available channels")
}

func (q *XQueue) PushInQueue(in *typex.InEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		I:    in,
		Data: data,
	}
	return q.push(qd, q.InQueue)
}
func PushInQueue(in *typex.InEnd, data string) error {
	return __DefaultXQueue.PushInQueue(in, data)
}
func (q *XQueue) PushDeviceQueue(device *typex.Device, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		D:    device,
		Data: data,
	}
	return q.push(qd, q.DeviceQueue)
}
func PushDeviceQueue(device *typex.Device, data string) error {
	return __DefaultXQueue.PushDeviceQueue(device, data)
}
func (q *XQueue) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		O:    out,
		Data: data,
	}
	return q.push(qd, q.OutQueue)
}
func PushOutQueue(out *typex.OutEnd, data string) error {
	return __DefaultXQueue.PushOutQueue(out, data)
}
