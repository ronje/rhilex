// Copyright (C) 2025 wwhai
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

package multimedia

import (
	"context"
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/glogger"
)

func LoadMultimediaStreamWithCtx(multimedia *MultimediaStream,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()

	return nil
}

func StartMultimediaStream(uuid string, Env *lua.LTable) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	// 先从缓存里面获取
	multimedia := GetMultimediaStream(uuid)
	if multimedia != nil {
		return fmt.Errorf("multimedia already exists:%s", uuid)
	}
	return nil
}

/*
*
* 从内存里面删除cecollalet
*
 */
func RemoveMultimediaStream(uuid string) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	delete(__DefaultMultimediaRuntime.MultimediaStreams, uuid)
	glogger.GLogger.Info("MultimediaStream removed:", uuid)
	return nil
}

/*
*
* 停止应用并不删除应用, 将其进程结束，状态置0
*
 */
func StopMultimediaStream(uuid string) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	// 先从缓存里面获取
	multimedia := GetMultimediaStream(uuid)
	if multimedia == nil {
		return fmt.Errorf("multimedia not exists:%s", uuid)
	}
	// 停止
	if multimedia.xMultimediaStream != nil {
		multimedia.xMultimediaStream.Stop()
	}
	glogger.GLogger.Info("MultimediaStream removed:", uuid)
	return nil
}

/*
*
* 更新应用信息
*
 */
func UpdateMultimediaStream(newMultimedia MultimediaStream) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	oldMultimedia := GetMultimediaStream(newMultimedia.UUID)
	if oldMultimedia == nil {
		return fmt.Errorf("update failed, multimedia not exists:%s", newMultimedia.UUID)
	}
	if err := StopMultimediaStream(oldMultimedia.UUID); err != nil {
		return err
	}
	if err := RemoveMultimediaStream(oldMultimedia.UUID); err != nil {
		return err
	}
	if err := LoadNewestMultimedia(newMultimedia.UUID, __DefaultMultimediaRuntime.RuleEngine); err != nil {
		return err
	}
	return nil
}
func GetMultimediaStream(uuid string) *MultimediaStream {
	if multimedia, ok := __DefaultMultimediaRuntime.MultimediaStreams[uuid]; ok {
		return multimedia
	}
	return nil
}

/*
*
* 获取列表
*
 */
func MultimediaStreamCount() int {
	return len(__DefaultMultimediaRuntime.MultimediaStreams)
}
func AllMultimediaStream() []*MultimediaStream {
	return ListMultimediaStream()
}
func ListMultimediaStream() []*MultimediaStream {
	cecollalets := []*MultimediaStream{}
	for _, v := range __DefaultMultimediaRuntime.MultimediaStreams {
		cecollalets = append(cecollalets, v)
	}
	return cecollalets
}
