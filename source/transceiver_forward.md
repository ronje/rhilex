# USB、串口通信模块信号接受处理器
当通信模块向内部总线`transceiver.up.data.$ComName`发送事件时，在此可以收到，从而做业务逻辑处理。

## 模块列表
`transceiver.up.data.$ComName`表示来自某个设备的数据。下面是支持的设备列表.
- 妙想科技蓝牙模块：`transceiver.up.data.MX01`
- 原子科技LORA模块：`transceiver.up.data.ATK01`
- 亿百特科技E22：`transceiver.up.data.E22`
- 亿百特科技E34：`transceiver.up.data.E34`