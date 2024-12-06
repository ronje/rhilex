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

package cecollalet

import (
	"context"
	"fmt"
	"sync"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/luaruntime"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/sirupsen/logrus"
)

var __DefaultCecollaletRuntime *CecollaletRuntime

func InitCecollaletRuntime(re typex.Rhilex) *CecollaletRuntime {
	__DefaultCecollaletRuntime = &CecollaletRuntime{
		RuleEngine:  re,
		locker:      sync.Mutex{},
		Cecollalets: make(map[string]*Cecollalet),
	}
	// Cecolla Config
	intercache.RegisterSlot("__CecollaBinding")
	return __DefaultCecollaletRuntime
}

/*
*
* 加载本地文件到lua虚拟机, 但是并不执行
*
 */
func LoadCecollalet(cecollalet *Cecollalet, luaSource string) error {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	// 重新读
	cecollalet.VM().DoString(string(luaSource))
	// 检查函数入口
	CecollaletMainVM := cecollalet.VM().GetGlobal("Main")
	if CecollaletMainVM == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if CecollaletMainVM.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 抽取main
	fMain := *CecollaletMainVM.(*lua.LFunction)
	cecollalet.SetMainFunc(&fMain)
	// 加载库
	// LoadCecollaletLib(cecollalet, __DefaultCecollaletRuntime.RuleEngine)
	luaruntime.LoadRuleLibGroup(__DefaultCecollaletRuntime.RuleEngine, "CECOLLA", cecollalet.UUID, cecollalet.VM())
	// 加载到内存里
	__DefaultCecollaletRuntime.Cecollalets[cecollalet.UUID] = cecollalet
	return nil
}

/*
* 此时才是真正的启动入口:
* 启动 function Main(args) --do-some-thing-- return 0 end
*
 */
func StartCecollalet(uuid string, Env *lua.LTable) error {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	cecollalet, ok := __DefaultCecollaletRuntime.Cecollalets[uuid]
	if !ok {
		return fmt.Errorf("Cecollalet not exists:%s", uuid)
	}
	if cecollalet.CecollaletState == 1 {
		return fmt.Errorf("Cecollalet already started:%s", uuid)
	}
	ctx, cancel := context.WithCancel(typex.GCTX)
	cecollalet.SetCnC(ctx, cancel)
	go func(cecollalet *Cecollalet) {
		defer func() {
			glogger.GLogger.Debug("Cecollalet exit:", cecollalet.UUID)
			cecollalet.VM().Pop(1) // 防止registry溢出
			cecollalet.CecollaletState = 0
		}()
		glogger.GLogger.Debugf("Ready to run Cecollalet:%s", cecollalet.UUID)
		cecollalet.CecollaletState = 1
		err := cecollalet.VM().CallByParam(lua.P{
			Fn:      cecollalet.GetMainFunc(),
			NRet:    1,
			Protect: true,
			Handler: &lua.LFunction{
				GFunction: func(*lua.LState) int {
					glogger.GLogger.Debug("Protect Mode Call")
					return 0
				},
			},
		}, lua.LString(uuid), Env)
		if err == nil {
			if cecollalet.KilledBy == "RHILEX" {
				glogger.GLogger.Infof("Cecollalet %s Killed By RHILEX", cecollalet.UUID)
			}
			if cecollalet.KilledBy == "NORMAL" || cecollalet.KilledBy == "" {
				glogger.GLogger.Infof("Cecollalet %s NORMAL Exited", cecollalet.UUID)
			}
			return
		}
		Debugger, Ok := cecollalet.vm.GetStack(1)
		if Ok {
			LValue, _ := cecollalet.vm.GetInfo("f", Debugger, lua.LNil)
			cecollalet.vm.GetInfo("l", Debugger, lua.LNil)
			cecollalet.vm.GetInfo("S", Debugger, lua.LNil)
			cecollalet.vm.GetInfo("u", Debugger, lua.LNil)
			cecollalet.vm.GetInfo("n", Debugger, lua.LNil)
			LFunction := LValue.(*lua.LFunction)
			LastCall := lua.DbgCall{
				Name: "_main",
			}
			if len(LFunction.Proto.DbgCalls) > 0 {
				LastCall = LFunction.Proto.DbgCalls[0]
			}
			glogger.GLogger.WithFields(logrus.Fields{
				"topic": "cecollalet/console/" + uuid,
			}).Errorf("Function Name: [%s],"+
				"What: [%s], Source Line: [%d],"+
				" Last Call: [%s], Error message: %s",
				Debugger.Name, Debugger.What, Debugger.CurrentLine,
				LastCall.Name, err.Error(),
			)
		}
		//
		// 检查是自己死的还是被RHILEX杀死
		// 1 正常结束
		// 2 被rhilex删除
		// 3 跑飞了
		// 中间出现异常挂了，此时要根据: auto start 来判断是否抢救
		time.Sleep(5 * time.Second)
		if cecollalet.KilledBy == "RHILEX" {
			glogger.GLogger.Infof("Cecollalet %s Killed By RHILEX, No need to rescue", cecollalet.UUID)
			return
		}
		if cecollalet.KilledBy == "NORMAL" {
			glogger.GLogger.Infof("Cecollalet %s NORMAL Exited, No need to rescue", cecollalet.UUID)
			return
		}
		glogger.GLogger.Warnf("Cecollalet %s Exited With error: %s, Maybe accident, Try to survive",
			cecollalet.UUID, err.Error())
	}(cecollalet)
	glogger.GLogger.Info("Cecollalet started:", cecollalet.UUID)
	return nil
}

/*
*
* 从内存里面删除cecollalet
*
 */
func RemoveCecollalet(uuid string) error {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	if cecollalet, ok := __DefaultCecollaletRuntime.Cecollalets[uuid]; ok {
		cecollalet.Remove()
		delete(__DefaultCecollaletRuntime.Cecollalets, uuid)
	}
	glogger.GLogger.Info("Cecollalet removed:", uuid)
	return nil
}

/*
*
* 停止应用并不删除应用, 将其进程结束，状态置0
*
 */
func StopCecollalet(uuid string) error {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	if cecollalet, ok := __DefaultCecollaletRuntime.Cecollalets[uuid]; ok {
		cecollalet.Stop()
	}
	glogger.GLogger.Info("Cecollalet removed:", uuid)
	return nil
}

/*
*
* 更新应用信息
*
 */
func UpdateCecollalet(cecollalet Cecollalet) error {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	if oldCecollalet, ok := __DefaultCecollaletRuntime.Cecollalets[cecollalet.UUID]; ok {
		oldCecollalet.Name = cecollalet.Name
		oldCecollalet.Version = cecollalet.Version
		glogger.GLogger.Info("Cecollalet updated:", cecollalet.UUID)
		return nil
	}
	return fmt.Errorf("update failed, cecollalet not exists:%s", cecollalet.UUID)

}
func GetCecollalet(uuid string) *Cecollalet {
	if cecollalet, ok := __DefaultCecollaletRuntime.Cecollalets[uuid]; ok {
		return cecollalet
	}
	return nil
}

/*
*
* 获取列表
*
 */
func CecollaletCount() int {
	return len(__DefaultCecollaletRuntime.Cecollalets)
}
func AllCecollalet() []*Cecollalet {
	return ListCecollalet()
}
func ListCecollalet() []*Cecollalet {
	cecollalets := []*Cecollalet{}
	for _, v := range __DefaultCecollaletRuntime.Cecollalets {
		cecollalets = append(cecollalets, v)
	}
	return cecollalets
}

func Stop() {
	__DefaultCecollaletRuntime.locker.Lock()
	defer __DefaultCecollaletRuntime.locker.Unlock()
	for _, cecollalet := range __DefaultCecollaletRuntime.Cecollalets {
		glogger.GLogger.Info("Stop Cecollalet:", cecollalet.UUID)
		cecollalet.Stop()
		glogger.GLogger.Info("Stop Cecollalet:", cecollalet.UUID, " Successfully")
	}
	intercache.UnRegisterSlot("__CecollaBinding")
	glogger.GLogger.Info("cecollalet stopped")

}
func GetRhilex() typex.Rhilex {
	return __DefaultCecollaletRuntime.RuleEngine
}
