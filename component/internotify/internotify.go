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
	"fmt"

	"github.com/hootrhino/rhilex/utils"
)

// ---------------------------------------------------------
// Type
// ---------------------------------------------------------
// - WARNING:  警告
// - ERROR:    错误
// - INFO:     信息
// - FATAL:    致命错误

type BaseEvent struct {
	Type    string
	Event   string
	Ts      uint64
	Summary string
	Info    interface{}
}

func (be BaseEvent) String() string {
	return fmt.Sprintf(`Event: [%s], [%s], %v`, be.Type, be.Event, be.Info)
}

/*
*
* Push
*
 */
func Insert(Event BaseEvent) error {
	// glogger.GLogger.Debug("Internal Event:", Event)
	InterNotifyDb().Table("m_internal_notifies").Save(&MInternalNotify{
		UUID:    utils.MakeUUID("NOTIFY"),
		Type:    Event.Type,  // INFO | ERROR | WARNING
		Status:  1,           // Default unread
		Event:   Event.Event, // 事件
		Ts:      Event.Ts,    // Unix毫秒 时间戳
		Summary: "Internal Event: " + Event.Event,
		Info:    Event.String(),
	})
	return nil
}
