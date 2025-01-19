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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/eventbus"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/supervisor"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 设备监控器 5秒检查一下状态
*
 */
func StartMultimediaSupervisor(MultimediaCtx context.Context, Multimedia *MultimediaStream, ruleEngine typex.Rhilex) {
	UUID := Multimedia.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(Multimedia.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For Multimedia:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-context.Background().Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-MultimediaCtx.Done():
			{
				glogger.GLogger.Debugf("Multimedia Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
		}
		// 被删除后就直接退出监督进程
		currentMultimedia := GetMultimediaStream(UUID)
		if currentMultimedia == nil {
			glogger.GLogger.Debugf("Multimedia:%v Deleted, supervisor exit", UUID)
			return
		}

		// 资源可能不会及时DOWN
		currentMultimediaStatus := currentMultimedia.xMultimediaStream.Status()
		if currentMultimediaStatus == MEDIA_DOWN {
			ErrMsg := ""
			Slot := intercache.GetSlot("__DefaultRuleEngine")
			if Slot != nil {
				CacheValue, ok := Slot[currentMultimedia.UUID]
				if ok {
					ErrMsg = CacheValue.ErrMsg
				}
			}
			info := fmt.Sprintf("Multimedia:(%s,%s) DOWN, supervisor try to Restart, error message: %s",
				UUID, currentMultimedia.Name, ErrMsg)
			glogger.GLogger.Debug(info)
			lineS := "event.multimedia.down." + UUID
			eventbus.Publish(lineS, eventbus.EventMessage{
				Topic:   lineS,
				From:    "res-supervisor",
				Type:    "Multimedia",
				Event:   lineS,
				Ts:      uint64(time.Now().UnixMilli()),
				Payload: ErrMsg,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestMultimedia(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 云边协同器
*
 */
var loadMultimediaLocker = sync.Mutex{}

// LoadNewestMultimedia
func LoadNewestMultimedia(uuid string, ruleEngine typex.Rhilex) error {
	loadMultimediaLocker.Lock()
	defer loadMultimediaLocker.Unlock()
	mMultimedia, err := service.GetMultiMediaWithUUID(uuid)
	if err != nil {
		return err
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mMultimedia.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	old := GetMultimediaStream(uuid)
	if old != nil {
		old.xMultimediaStream.Stop()
	}
	RemoveMultimediaStream(uuid) // 删除内存里面的
	multimedia := MultimediaStream{}
	// Important !!!!!!!!
	multimedia.UUID = mMultimedia.UUID // 本质上是配置和内存的数据映射起来
	// 最新的配置
	multimedia.Config = mMultimedia.GetConfig()
	// 参数传给 --> startMultimedia()
	ctx, cancelCTX := typex.NewCCTX()
	err2 := LoadMultimediaStreamWithCtx(&multimedia, ctx, cancelCTX)
	if err2 != nil {
		glogger.GLogger.Error(err2)
		// return err2
	}
	go StartMultimediaSupervisor(ctx, &multimedia, ruleEngine)
	return nil

}
