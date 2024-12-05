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

package alarmcenter

import (
	"errors"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/typex"
)

var __DefaultAlarmCenter *AlarmCenter

type AlarmCenter struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[string, *vm.Program]
}

func InitAlarmCenter(e typex.Rhilex) {
	__DefaultAlarmCenter = &AlarmCenter{
		e:        e,
		registry: orderedmap.NewOrderedMap[string, *vm.Program](),
	}
}

// Load Expr
func LoadExpr(uuid, exprs string) (*vm.Program, error) {
	Program, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return nil, err
	}
	__DefaultAlarmCenter.registry.Set(uuid, Program)
	return Program, nil
}

// 测试
func TestRunExpr(exprs string, in map[string]any) (bool, error) {
	Program, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return false, err
	}
	output, err := expr.Run(Program, in)
	return output.(bool), err
}

// Run Expr
func RunExpr(uuid string, in map[string]any) (bool, error) {
	Program, ok := __DefaultAlarmCenter.registry.Get(uuid)
	if ok {
		// TODO: 这里有危险，应该增加资源限制控制.
		output, err := expr.Run(Program, in)
		return output.(bool), err
	}
	return false, errors.New("Invalid Expr vm")
}

// Remove Expr
func RemoveExpr(uuid string) {
	__DefaultAlarmCenter.registry.Delete(uuid)
}

// Flush
func FlushAlarmCenter() {
	for _, v := range __DefaultAlarmCenter.registry.Keys() {
		__DefaultAlarmCenter.registry.Delete(v)
	}
}
