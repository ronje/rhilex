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

package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	intercache "github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/rhilexmanager"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 加载云边协同器
*
 */
func (e *RuleEngine) LoadCecollaWithCtx(cecollaInstance *typex.Cecolla,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	if config := rhilexmanager.DefaultCecollaTypeManager.Find(cecollaInstance.Type); config != nil {
		return e.loadCecollas(config.NewCecolla(e), cecollaInstance, ctx, cancelCTX)
	}
	return fmt.Errorf("unsupported Cecolla type:%s", cecollaInstance.Type)

}

/*
*
* 启动一个和RHILEX直连的外部云边协同器
*
 */

func (e *RuleEngine) loadCecollas(xCecolla typex.XCecolla, cecollaInstance *typex.Cecolla,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	// Bind
	// xCecolla: Interface
	// cecollaInstance: Real Worker, Running instance
	cecollaInstance.Cecolla = xCecolla
	e.SaveCecolla(cecollaInstance)
	// Load config
	// 要从数据库里面查Config
	config := e.GetCecolla(cecollaInstance.UUID).Config
	if config == nil {
		e.RemoveCecolla(cecollaInstance.UUID)
		err := fmt.Errorf("cecolla [%v] config is nil", cecollaInstance.Name)
		return err
	}
	if err := xCecolla.Init(cecollaInstance.UUID, config); err != nil {
		intercache.SetValue("__DefaultRuleEngine", cecollaInstance.UUID, intercache.CacheValue{
			UUID:          cecollaInstance.UUID,
			Status:        1,
			ErrMsg:        err.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
		configBytes, _ := json.Marshal(config)
		// 注册一个缓存器
		intercache.SetValue("__CecollaConfigMap", cecollaInstance.UUID, intercache.CacheValue{
			UUID:          cecollaInstance.UUID,
			Status:        1,
			ErrMsg:        err.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         string(configBytes),
		})
		return err
	}
	err2 := startCecolla(xCecolla, ctx, cancelCTX)
	if err2 != nil {
		glogger.GLogger.Error(err2)
		intercache.SetValue("__DefaultRuleEngine", cecollaInstance.UUID, intercache.CacheValue{
			UUID:          cecollaInstance.UUID,
			Status:        1,
			ErrMsg:        err2.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
	} else {
		intercache.DeleteValue("__DefaultRuleEngine", cecollaInstance.UUID) // 删除云边协同器缓存
		intercache.DeleteValue("__CecollaConfigMap", cecollaInstance.UUID)  // 删除配置缓存
	}
	glogger.GLogger.Infof("Cecolla [%v, %v] load successfully", cecollaInstance.Name, cecollaInstance.UUID)
	return nil
}

/*
*
* Start是异步进行的,当云边协同器的GetStatus返回状态UP时，正常运行，当Down时重启
*
 */
func startCecolla(xCecolla typex.XCecolla, ctx context.Context, cancelCTX context.CancelFunc) error {
	if err := xCecolla.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("Cecolla start error:", err)
		return err
	}
	xCecolla.SetState(typex.CEC_UP)
	return nil
}
