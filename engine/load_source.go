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

package engine

import (
	"context"
	"fmt"
	"time"

	intercache "github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/registry"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 加载输入资源
*
 */
func (e *RuleEngine) LoadInEndWithCtx(in *typex.InEnd,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	if config := registry.DefaultSourceRegistry.Find(in.Type); config != nil {
		return e.loadSource(config.NewSource(e), in, ctx, cancelCTX)
	}
	return fmt.Errorf("unsupported InEnd type:%s", in.Type)
}

//
// start Sources
//
/*
* Life cycle
+------------------+       +------------------+   +---------------+        +---------------+
|     Init         |------>|   Start          |-->|     Test      |--+ --->|  Stop         |
+------------------+  ^    +------------------+   +---------------+  |     +---------------+
                      |                                              |
                      |                                              |
                      +-------------------Error ---------------------+
*/
func (e *RuleEngine) loadSource(source typex.XSource, in *typex.InEnd,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	in.Source = source
	e.SaveInEnd(in)
	// Load config
	config := e.GetInEnd(in.UUID).Config
	if config == nil {
		e.RemoveInEnd(in.UUID)
		err := fmt.Errorf("source [%v, %v] config is nil", in.UUID, in.Name)
		return err
	}
	if err := source.Init(in.UUID, config); err != nil {
		glogger.GLogger.Error(err)
		intercache.SetValue("__DefaultRuleEngine", in.UUID, intercache.CacheValue{
			UUID:          in.UUID,
			Status:        1,
			ErrMsg:        err.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
		return err
	}

	err2 := e.startSource(source, ctx, cancelCTX)
	if err2 != nil {
		glogger.GLogger.Error(err2)
		intercache.SetValue("__DefaultRuleEngine", in.UUID, intercache.CacheValue{
			UUID:          in.UUID,
			Status:        1,
			ErrMsg:        err2.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
	} else {
		intercache.DeleteValue("__DefaultRuleEngine", in.UUID)
	}
	glogger.GLogger.Infof("InEnd [%v, %v] load successfully", in.UUID, in.Name)
	return nil
}

func (e *RuleEngine) startSource(source typex.XSource,
	ctx context.Context, cancelCTX context.CancelFunc) error {

	if err := source.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("Source start error:", err)
		return err
	}
	// LoadNewestSource
	// 2023-06-14新增： 重启成功后数据会丢失,还得加载最新的Rule到设备中
	Source := source.Details()
	if Source != nil {
		for _, rule := range Source.BindRules {
			RuleInstance := typex.NewLuaRule(e,
				rule.UUID,
				rule.Name,
				rule.Description,
				rule.FromSource,
				rule.FromDevice,
				rule.Success,
				rule.Actions,
				rule.Failed)
			if err1 := e.LoadRule(RuleInstance); err1 != nil {
				return err1
			}
		}
	}
	return nil
}
