package rhilexlib

import (
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/interqueue"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

func DataToMqtt(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		err := handleDataFormat(rx, id, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}
func DataToMqttTopic(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		topic := l.ToString(3)
		data := l.ToString(4)
		err := handleMqttFormat(rx, id, topic, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

// 处理MQTT消息
// 支持自定义MQTT Topic, 需要在Target的to接口来实现这个
func handleMqttFormat(e typex.Rhilex, uuid string, topic string, incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		return interqueue.DefaultDataCacheQueue.PushOutQueue(outEnd, incoming)
	}
	msg := "target not found:" + uuid
	glogger.GLogger.Error(msg)
	return errors.New(msg)

}
