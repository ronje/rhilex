package rhilexlib

import (
	"strings"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/sirupsen/logrus"
)

// | Ws log topic                     | 用途               |
// | -------------------------------- | ------------------ |
// | plugin/${name}/${uuid}           | 插件的日志 |
// | rule/log/${uuid}                 | 规则运行时的日志   |
// | app/console/${uuid}              | app运行地址         |
// | device/rule/test/${uuid}         | Test device        |
// | inend/rule/test/${uuid}          | Test inend         |

/*
*
* APP debug输出, Debug(".....")
*
 */
func DebugAPP(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		content := ""
		for i := 1; i <= top; i++ {
			content += L.ToStringMeta(L.Get(i)).String()
			if i != top {
				content += "  "
			}
		}
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "app/console/" + uuid,
		}).Info(content)
		return 0
	}
}

/*
*
* 辅助Debug使用, 用来向前端Dashboard打印日志的时候带上ID
*
 */
func DebugRule(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		content := ""
		for i := 1; i <= top; i++ {
			content += L.ToStringMeta(L.Get(i)).String()
			if i != top {
				content += "  "
			}
		}
		// ::::TEST_RULE:::: 用来标记是否是测试数据
		TestPrefix := "::::TEST_RULE::::"
		if strings.HasPrefix(content, TestPrefix) {
			if content[len(TestPrefix):] == "" {
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/test/" + uuid,
				}).Info("<Empty>")
			} else {
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/test/" + uuid,
				}).Info(content[len(TestPrefix):])
			}

		} else {
			glogger.GLogger.WithFields(logrus.Fields{
				"topic": "rule/log/" + uuid,
			}).Info(content)
		}

		return 0
	}
}

// 输出
func DebugCecolla(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		content := ""
		for i := 1; i <= top; i++ {
			content += L.ToStringMeta(L.Get(i)).String()
			if i != top {
				content += "  "
			}
		}
		glogger.GLogger.WithFields(logrus.Fields{
			"topic": "cecolla/console/" + uuid,
		}).Info(content)
		return 0
	}
}
