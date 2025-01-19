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

package multimedia

import (
	"sync"

	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

var __DefaultMultimediaRuntime *MultimediaRuntime

/*
*
* 管理器
*
 */
type MultimediaRuntime struct {
	locker            sync.Mutex
	RuleEngine        typex.Rhilex
	MultimediaStreams map[string]*MultimediaStream
}

func InitMultimediaRuntime(re typex.Rhilex) *MultimediaRuntime {
	__DefaultMultimediaRuntime = &MultimediaRuntime{
		RuleEngine:        re,
		locker:            sync.Mutex{},
		MultimediaStreams: make(map[string]*MultimediaStream),
	}
	// Cecolla Config
	intercache.RegisterSlot("__MultimediaBinding")
	return __DefaultMultimediaRuntime
}

func Stop() {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	for _, v := range __DefaultMultimediaRuntime.MultimediaStreams {
		StopMultimediaStream(v.UUID)
	}
	intercache.UnRegisterSlot("__MultimediaBinding")
	glogger.GLogger.Info("multimedia stopped")
}
