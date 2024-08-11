package interqueue

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultYQueue *YQueue

/*
*
* YQueue
*
 */
type YQueue struct {
	Queue        *list.List
	OutQueue     *list.List
	InQueue      *list.List
	DeviceQueue  *list.List
	rhilex       typex.Rhilex
	inLocker     sync.RWMutex
	deviceLocker sync.RWMutex
	outLocker    sync.RWMutex
	queueSize    int
}

func InitYQueue(rhilex typex.Rhilex, queueSize int) *YQueue {
	DefaultYQueue = &YQueue{
		Queue:        list.New(),
		OutQueue:     list.New(),
		InQueue:      list.New(),
		DeviceQueue:  list.New(),
		rhilex:       rhilex,
		inLocker:     sync.RWMutex{},
		deviceLocker: sync.RWMutex{},
		outLocker:    sync.RWMutex{},
		queueSize:    queueSize,
	}
	return DefaultYQueue
}

/*
*
*GetDeviceQueue
*
 */
func (q *YQueue) GetDeviceQueue() *list.List {
	return q.DeviceQueue
}

func StartYQueue() {
	ctx := typex.GCTX
	go func(ctx context.Context, queue *YQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			queue.inLocker.Lock()
			for listE := queue.InQueue.Back(); listE != nil; listE = queue.Queue.Back() {
				switch T := listE.Value.(type) {
				case QueueData:
					if T.I != nil {
						T.E.RunSourceCallbacks(T.I, T.Data)
					}
					queue.InQueue.Remove(listE)
				}
			}
			queue.inLocker.Unlock()
			time.Sleep(8 * time.Millisecond)
		}

	}(ctx, DefaultYQueue)
	go func(ctx context.Context, queue *YQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			queue.deviceLocker.Lock()
			for listE := queue.DeviceQueue.Back(); listE != nil; listE = queue.Queue.Back() {
				switch T := listE.Value.(type) {
				case QueueData:
					if T.D != nil {
						T.E.RunDeviceCallbacks(T.D, T.Data)
					}
					queue.DeviceQueue.Remove(listE)
				}
			}
			queue.deviceLocker.Unlock()
			time.Sleep(8 * time.Millisecond)
		}

	}(ctx, DefaultYQueue)
	go func(ctx context.Context, queue *YQueue) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			queue.outLocker.Lock()
			for listE := queue.OutQueue.Back(); listE != nil; listE = queue.Queue.Back() {
				switch QueueData := listE.Value.(type) {
				case QueueData:
					if QueueData.O != nil {
						ProcessOutQueueData(QueueData, QueueData.E)
					}
					queue.OutQueue.Remove(listE)
				}
			}
			queue.outLocker.Unlock()
			time.Sleep(8 * time.Millisecond)
		}

	}(ctx, DefaultYQueue)
}

/*
*
*PushInQueue
*
 */
func (q *YQueue) PushInQueue(in *typex.InEnd, data string) error {

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
func (q *YQueue) PushDeviceQueue(Device *typex.Device, data string) error {
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
func (q *YQueue) PushOutQueue(out *typex.OutEnd, data string) error {
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
func (q *YQueue) pushIn(d QueueData) error {
	if q.InQueue.Len()+1 > q.queueSize {
		msg := fmt.Sprintf("Exceed max YQueue size:%v", q.queueSize)
		return errors.New(msg)
	}
	q.inLocker.Lock()
	defer q.inLocker.Unlock()
	q.InQueue.PushBack(d)
	return nil
}

/*
*
* Push
*
 */
func (q *YQueue) pushOut(d QueueData) error {
	if q.OutQueue.Len()+1 > q.queueSize {
		msg := fmt.Sprintf("Exceed max YQueue size:%v", q.queueSize)
		return errors.New(msg)
	}
	q.outLocker.Lock()
	defer q.outLocker.Unlock()
	q.OutQueue.PushBack(d)
	return nil
}

/*
*
* Push
*
 */
func (q *YQueue) pushDevice(d QueueData) error {
	if q.DeviceQueue.Len()+1 > q.queueSize {
		msg := fmt.Sprintf("Exceed max YQueue size:%v", q.queueSize)
		return errors.New(msg)
	}
	q.deviceLocker.Lock()
	defer q.deviceLocker.Unlock()
	q.DeviceQueue.PushBack(d)
	return nil
}
