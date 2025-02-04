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

func (e *RuleEngine) LoadOutEndWithCtx(in *typex.OutEnd, ctx context.Context,
	cancelCTX context.CancelFunc) error {
	if config := registry.DefaultTargetRegistry.Find(in.Type); config != nil {
		return e.loadTarget(config.NewTarget(e), in, ctx, cancelCTX)
	}
	return fmt.Errorf("unsupported Target type:%s", in.Type)
}

// Start output target
//
// Target life cycle:
//
//	Register -> Start -> running/restart cycle
func (e *RuleEngine) loadTarget(target typex.XTarget, out *typex.OutEnd,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	// Set sources to inend
	out.Target = target
	e.SaveOutEnd(out)
	// Load config
	config := e.GetOutEnd(out.UUID).Config
	if config == nil {
		e.RemoveOutEnd(out.UUID)
		err := fmt.Errorf("target [%v] config is nil", out.Name)
		return err
	}
	if err := target.Init(out.UUID, config); err != nil {
		glogger.GLogger.Error(err)
		intercache.SetValue("__DefaultRuleEngine", out.UUID, intercache.CacheValue{
			UUID:          out.UUID,
			Status:        1,
			ErrMsg:        err.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
		return err
	}
	err2 := startTarget(target, ctx, cancelCTX)
	if err2 != nil {
		glogger.GLogger.Error(err2)
		intercache.SetValue("__DefaultRuleEngine", out.UUID, intercache.CacheValue{
			UUID:          out.UUID,
			Status:        1,
			ErrMsg:        err2.Error(),
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
		})
	} else {
		intercache.DeleteValue("__DefaultRuleEngine", out.UUID)
	}
	glogger.GLogger.Infof("Target [%v, %v] load successfully", out.Name, out.UUID)
	return nil
}

func startTarget(target typex.XTarget, ctx context.Context, cancelCTX context.CancelFunc) error {
	if err := target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("abstractDevice start error:", err)
		return err
	}
	return nil
}
