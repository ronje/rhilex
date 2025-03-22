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

	"github.com/hootrhino/rhilex/component/intermetric"
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

// 初始化队列
func InitXQueue(rhilex typex.Rhilex, maxQueueSize int) *XQueue {
	queueCount := 10
	inQueue := make([]chan QueueData, queueCount)
	outQueue := make([]chan QueueData, queueCount)
	deviceQueue := make([]chan QueueData, queueCount)

	// 初始化队列通道
	initQueues := func(queues []chan QueueData) {
		for i := 0; i < queueCount; i++ {
			queues[i] = make(chan QueueData, maxQueueSize)
		}
	}
	initQueues(inQueue)
	initQueues(outQueue)
	initQueues(deviceQueue)

	__DefaultXQueue = &XQueue{
		InQueue:      inQueue,
		OutQueue:     outQueue,
		DeviceQueue:  deviceQueue,
		rhilex:       rhilex,
		locker:       sync.Mutex{},
		maxQueueSize: maxQueueSize,
	}
	return __DefaultXQueue
}

// 启动队列
func StartXQueue() {
	__DefaultXQueue.StartXQueue()
}

// 通用的推送函数
func pushData(q *XQueue, data QueueData, queue []chan QueueData) error {
	// 减少锁的使用范围
	var available []int
	for i := 0; i < len(queue); i++ {
		if len(queue[i])+1 <= q.maxQueueSize {
			available = append(available, i)
		}
	}
	if len(available) == 0 {
		return fmt.Errorf("Queue Send Failed: No available channels")
	}
	// 更高效的随机数生成
	randomIndex := available[rand.Intn(len(available))]
	queue[randomIndex] <- data
	return nil
}

// 通用的启动队列处理函数
func startQueue(ctx context.Context, q *XQueue, queue []chan QueueData, logMsg string, callback func(QueueData)) {
	glogger.GLogger.Info("Start XQueue: " + logMsg)
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

// 启动输入队列
func (q *XQueue) startInQueue(ctx context.Context) {
	startQueue(ctx, q, q.InQueue, "InQueue", func(data QueueData) {
		if data.I == nil || data.E == nil {
			return
		}
		luaexecutor.RunSourceCallbacks(data.I, data.Data)
	})
}

// 启动设备队列
func (q *XQueue) startDeviceQueue(ctx context.Context) {
	startQueue(ctx, q, q.DeviceQueue, "DeviceQueue", func(data QueueData) {
		if data.D == nil || data.E == nil {
			return
		}
		luaexecutor.RunDeviceCallbacks(data.D, data.Data)
	})
}

// 启动输出队列
func (q *XQueue) startOutQueue(ctx context.Context) {
	startQueue(ctx, q, q.OutQueue, "OutQueue", func(data QueueData) {
		ProcessOutQueueData(data, data.E)
	})
}

// 启动所有队列
func (q *XQueue) StartXQueue() {
	ctx := context.Background()
	go q.startInQueue(ctx)
	go q.startDeviceQueue(ctx)
	go q.startOutQueue(ctx)
}

// 通用的推送包装函数
func pushWrapper[QueueType any](q *XQueue, pushFunc func(*XQueue, QueueType, string) error, target QueueType, data string) error {
	return pushFunc(q, target, data)
}

// 推送数据到输入队列
func (q *XQueue) PushInQueue(in *typex.InEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		I:    in,
		Data: data,
	}
	return pushData(q, qd, q.InQueue)
}
func PushInQueue(in *typex.InEnd, data string) error {
	return pushWrapper(__DefaultXQueue, (*XQueue).PushInQueue, in, data)
}

// 推送数据到设备队列
func (q *XQueue) PushDeviceQueue(device *typex.Device, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		D:    device,
		Data: data,
	}
	return pushData(q, qd, q.DeviceQueue)
}
func PushDeviceQueue(device *typex.Device, data string) error {
	return pushWrapper(__DefaultXQueue, (*XQueue).PushDeviceQueue, device, data)
}

// 推送数据到输出队列
func (q *XQueue) PushOutQueue(out *typex.OutEnd, data string) error {
	qd := QueueData{
		E:    q.rhilex,
		O:    out,
		Data: data,
	}
	return pushData(q, qd, q.OutQueue)
}
func PushOutQueue(out *typex.OutEnd, data string) error {
	return pushWrapper(__DefaultXQueue, (*XQueue).PushOutQueue, out, data)
}

type QueueData struct {
	Debug bool // 是否是Debug消息
	I     *typex.InEnd
	O     *typex.OutEnd
	D     *typex.Device
	E     typex.Rhilex
	Data  string
}

func (qd QueueData) String() string {
	return "QueueData@In:" + qd.I.UUID + ", Data:" + qd.Data
}

func ProcessOutQueueData(qd QueueData, e typex.Rhilex) {
	if qd.O != nil {
		target := e.GetOutEnd(qd.O.UUID)
		if target != nil {
			if _, err := target.Target.To(qd.Data); err != nil {
				glogger.GLogger.Error(err)
				intermetric.IncOutFailed()
			} else {
				intermetric.IncOut()
			}
		}
	}
}
