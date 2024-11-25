// 抽象设备：
// 1.0 以后的大功能：支持抽象设备，抽象设备就是外挂的设备，Rhilex本来是个规则引擎，但是1.0之前的版本没有对硬件设备进行抽象支持
// 因此，1.0以后增加对硬件的抽象
// Target Source 描述了数据的流向，抽象设备描述了数据的载体。
// 举例：外挂一个设备，这个设备具备双工控制功能，例如电磁开关等，此时它强调的是设备的物理功能，而数据则不是主体。
// 因此需要抽象出来一个层专门来描述这些设备
package typex

type DeviceState int

const (
	// 设备故障
	DEV_DOWN DeviceState = 0
	// 设备启用
	DEV_UP DeviceState = 1
	// 暂停，这是个占位值，只为了和其他地方统一值,但是没用
	_ DeviceState = 2
	// 外部停止
	DEV_STOP DeviceState = 3
	// 准备态
	DEV_PENDING DeviceState = 4
	// 被禁用
	DEV_DISABLE DeviceState = 5
)

func (s DeviceState) String() string {
	if s == 0 {
		return "DOWN"
	}
	if s == 1 {
		return "UP"
	}
	if s == 2 {
		return "PAUSE"
	}
	if s == 3 {
		return "STOP"
	}
	if s == 4 {
		return "PENDING"
	}
	if s == 5 {
		return "DISABLE"
	}
	return "ERROR"
}

type DeviceType string

func (d DeviceType) String() string {
	return string(d)

}

/**
 * 免费版只支持这三个
 *
 */
const (
	GENERIC_UART_RW       DeviceType = "GENERIC_UART_RW"       // 通用读写串口
	GENERIC_MODBUS_MASTER DeviceType = "GENERIC_MODBUS_MASTER" // 通用 GENERIC_MODBUS_MASTER
	GENERIC_MODBUS_SLAVER DeviceType = "GENERIC_MODBUS_SLAVER" // 通用 GENERIC_MODBUS_SLAVER
)

/**
 * 企业版
 *
 */
const (
	SIEMENS_PLC                 DeviceType = "SIEMENS_PLC"                 // SIEMENS-S71200
	GENERIC_SNMP                DeviceType = "GENERIC_SNMP"                // SNMP 协议支持
	GENERIC_CAMERA              DeviceType = "GENERIC_CAMERA"              // 通用摄像头
	GENERIC_BACNET_IP           DeviceType = "GENERIC_BACNET_IP"           // 通用Bacnet IP模式
	BACNET_ROUTER_GW            DeviceType = "BACNET_ROUTER_GW"            // 通用BACNET 路由模式
	GENERIC_HTTP_DEVICE         DeviceType = "GENERIC_HTTP_DEVICE"         // HTTP采集器
	TENCENT_IOTHUB_GATEWAY      DeviceType = "TENCENT_IOTHUB_GATEWAY"      // 腾讯云物联网平台
	ITHINGS_IOTHUB_GATEWAY      DeviceType = "ITHINGS_IOTHUB_GATEWAY"      // ITHINGS物联网平台
	LORA_WAN_GATEWAY            DeviceType = "LORA_WAN_GATEWAY"            // LoraWan
	KNX_GATEWAY                 DeviceType = "KNX_GATEWAY"                 // KNX 网关
	GENERIC_MBUS_EN13433_MASTER DeviceType = "GENERIC_MBUS_EN13433_MASTER" // 通用 Mbus
	DLT6452007_MASTER           DeviceType = "DLT6452007_MASTER"           // DLT6452004
	CJT1882004_MASTER           DeviceType = "CJT1882004_MASTER"           // CJT1882004
	SZY2062016_MASTER           DeviceType = "SZY2062016_MASTER"           // SZY2062016
	GENERIC_USER_PROTOCOL       DeviceType = "GENERIC_USER_PROTOCOL"       // 自定义协议
	GENERIC_AIS_RECEIVER        DeviceType = "GENERIC_AIS_RECEIVER"        // 通用AIS
	GENERIC_NEMA_GNS_PROTOCOL   DeviceType = "GENERIC_NEMA_GNS_PROTOCOL"   // GPS采集器
)

type DCAModel struct {
	UUID    string      `json:"uuid"`
	Command string      `json:"command"`
	Args    interface{} `json:"args"`
}
type DCAResult struct {
	Error error
	Data  string
}

// 真实工作设备,即具体实现
type XDevice interface {
	// 初始化 通常用来获取设备的配置
	Init(devId string, configMap map[string]interface{}) error
	// 启动, 设备的工作进程
	Start(CCTX) error
	// 新特性, 适用于自定义协议读写
	OnCtrl(cmd []byte, args []byte) ([]byte, error)
	// 设备当前状态
	Status() DeviceState
	// 停止设备, 在这里释放资源,一般是先置状态为STOP,然后CancelContext()
	Stop()

	Details() *Device
	// 状态
	SetState(DeviceState)
	// 外部调用, 该接口是个高级功能, 准备为了设计分布式部署设备的时候用, 但是相当长时间内都不会开启
	// 默认情况下该接口没有用
	OnDCACall(UUID string, Command string, Args interface{}) DCAResult
}

/*
*
* 子设备网络拓扑[2023-04-17新增]
*
 */
type DeviceTopology struct {
	Id       string                 // 子设备的ID
	Name     string                 // 子设备名
	LinkType int                    // 物理连接方式: 0-ETH 1-WIFI 3-BLE 4 LORA 5 OTHER
	State    int                    // 状态: 0-Down 1-Working
	Info     map[string]interface{} // 子设备的一些额外信息
}
