# USB、串口通信模块信号接受处理器
当通信模块向内部总线`transceiver.upstream.data.$ComName`发送事件时，在此可以收到，从而做业务逻辑处理。
```go
internotify.Push(internotify.BaseEvent{
	Type:    "transceiver.upstream.data",
	Event:   "transceiver.upstream.data.MX01-BLE-Module",
	Ts:      uint64(time.Now().UnixMilli()),
	Summary: "transceiver.upstream.data",
	Info:    []byte("HELLO WORLD\r\n"),
})
```
其中 `Event`里面带上通信模块的名称。方便过滤数据。