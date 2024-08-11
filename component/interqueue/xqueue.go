package interqueue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultXQueue Queue

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
	Queue       chan QueueData
	OutQueue    chan QueueData
	InQueue     chan QueueData
	DeviceQueue chan QueueData
	rhilex      typex.Rhilex
}

func InitXQueue(rhilex typex.Rhilex, maxQueueSize int) Queue {
	DefaultXQueue = &XQueue{
		Queue:       make(chan QueueData, maxQueueSize),
		OutQueue:    make(chan QueueData, maxQueueSize),
		InQueue:     make(chan QueueData, maxQueueSize),
		DeviceQueue: make(chan QueueData, maxQueueSize),
		rhilex:      rhilex,
	}
	return DefaultXQueue
}
func (q *XQueue) GetSize() int {
	return cap(q.Queue)
}

/*
*
* Push
*
 */
func (q *XQueue) Push(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.Queue)+1 > q.GetSize() {
		msg := fmt.Sprintf("exceed max queue size:%v", q.GetSize())
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.Queue <- d
		return nil
	}
}

/*
*
* GetQueue
*
 */
func (q *XQueue) GetQueue() chan QueueData {
	return q.Queue
}

/*
*
* GetQueue
*
 */
func (q *XQueue) GetInQueue() chan QueueData {
	return q.InQueue
}

/*
*
* GetQueue
*
 */
func (q *XQueue) GetOutQueue() chan QueueData {
	return q.OutQueue
}

/*
*
*GetDeviceQueue
*
 */
func (q *XQueue) GetDeviceQueue() chan QueueData {
	return q.DeviceQueue
}

// TODO: 下个版本更换为可扩容的Chan
func StartXQueue() {
	ctx := typex.GCTX
	go func(ctx context.Context, queue Queue) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				// 优雅地处理上下文取消
				return
			case qd := <-queue.GetInQueue():
				if qd.I != nil {
					qd.E.RunSourceCallbacks(qd.I, qd.Data)
				}
			case qd := <-queue.GetDeviceQueue():
				if qd.D != nil {
					qd.E.RunDeviceCallbacks(qd.D, qd.Data)
				}
			case qd := <-queue.GetOutQueue():
				ProcessOutQueueData(qd, qd.E)
			case <-ticker.C:
				continue
			}
		}
	}(ctx, DefaultXQueue)
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
	err := q.pushIn(qd)
	if err != nil {
		glogger.GLogger.Error("Push InQueue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
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
	err := q.pushDevice(qd)
	if err != nil {
		glogger.GLogger.Error("Push Device Queue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
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
	err := q.pushOut(qd)
	if err != nil {
		glogger.GLogger.Error("Push OutQueue error:", err)
		intermetric.IncInFailed()
	} else {
		intermetric.IncIn()
	}
	return err
}

/*
*
* Push
*
 */
func (q *XQueue) pushIn(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.InQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.InQueue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.InQueue <- d
		return nil
	}
}

/*
*
* Push
*
 */
func (q *XQueue) pushOut(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.OutQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.OutQueue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.OutQueue <- d
		return nil
	}
}

/*
*
* Push
*
 */
func (q *XQueue) pushDevice(d QueueData) error {
	// 动态扩容
	// if len(q.Queue)+1 > q.GetSize() {
	// }
	if len(q.DeviceQueue)+1 > q.GetSize() {
		msg := fmt.Sprintf("attached max queue size, max size is:%v, current size is: %v",
			q.GetSize(), len(q.DeviceQueue)+1)
		glogger.GLogger.Error(msg)
		return errors.New(msg)
	} else {
		q.DeviceQueue <- d
		return nil
	}
}
