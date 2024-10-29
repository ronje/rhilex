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
	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

type Queue interface {
	PushInQueue(in *typex.InEnd, data string) error
	PushOutQueue(in *typex.OutEnd, data string) error
	PushDeviceQueue(in *typex.Device, data string) error
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
