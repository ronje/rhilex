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
// GNU Affero General Public License for more regarding.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package multimedia

import (
	"context"
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

func LoadMultimediaStreamWithCtx(multimedia *MultimediaStream,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	__DefaultMultimediaRuntime.locker.Lock()
	defer __DefaultMultimediaRuntime.locker.Unlock()
	// 这里添加实际的加载逻辑，例如加载资源等
	// 例如：
	// loadResource(multimedia)
	return nil
}

func StartMultimediaStream(uuid string, Env *lua.LTable) error {
	__DefaultMultimediaRuntime.locker.RLock()
	multimedia := GetMultimediaStream(uuid)
	__DefaultMultimediaRuntime.locker.RUnlock()
	if multimedia == nil {
		return fmt.Errorf("multimedia not exists:%s", uuid)
	}
	// 这里可以添加启动多媒体流的逻辑
	// 例如：
	// multimedia.Start()
	return nil
}

/*
*
* 从内存里面删除多媒体流
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
* 停止应用并不删除应用, 将其进程结束，状态置 0
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
	glogger.GLogger.Info("MultimediaStream stopped:", uuid)
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

	// 停止旧的多媒体流
	if err := stopMultimediaStreamHelper(oldMultimedia); err != nil {
		return err
	}

	// 删除旧的多媒体流
	if err := removeMultimediaStreamHelper(oldMultimedia.UUID); err != nil {
		return err
	}

	// 加载新的多媒体流
	if err := loadNewestMultimediaHelper(newMultimedia.UUID, __DefaultMultimediaRuntime.RuleEngine); err != nil {
		return err
	}

	return nil
}

func stopMultimediaStreamHelper(multimedia *MultimediaStream) error {
	if multimedia.xMultimediaStream != nil {
		multimedia.xMultimediaStream.Stop()
	}
	glogger.GLogger.Info("MultimediaStream stopped:", multimedia.UUID)
	return nil
}

func removeMultimediaStreamHelper(uuid string) error {
	delete(__DefaultMultimediaRuntime.MultimediaStreams, uuid)
	glogger.GLogger.Info("MultimediaStream removed:", uuid)
	return nil
}

func loadNewestMultimediaHelper(uuid string, ruleEngine typex.Rhilex) error {
	// 这里添加实际的加载最新多媒体流的逻辑
	// 例如：
	// loadResource(newMultimedia)
	return nil
}

func GetMultimediaStream(uuid string) *MultimediaStream {
	__DefaultMultimediaRuntime.locker.RLock()
	defer __DefaultMultimediaRuntime.locker.RUnlock()
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
	__DefaultMultimediaRuntime.locker.RLock()
	defer __DefaultMultimediaRuntime.locker.RUnlock()
	return len(__DefaultMultimediaRuntime.MultimediaStreams)
}

func ListMultimediaStream() []*MultimediaStream {
	__DefaultMultimediaRuntime.locker.RLock()
	defer __DefaultMultimediaRuntime.locker.RUnlock()
	cecollalets := []*MultimediaStream{}
	for _, v := range __DefaultMultimediaRuntime.MultimediaStreams {
		cecollalets = append(cecollalets, v)
	}
	return cecollalets
}
