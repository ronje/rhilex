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

package typex

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
	TAOJINGCHI_UARTHMI_MASTER   DeviceType = "TAOJINGCHI_UARTHMI_MASTER"   // 陶晶池串口屏
)

type DCAModel struct {
	UUID    string `json:"uuid"`
	Command string `json:"command"`
	Args    any    `json:"args"`
}
type DCAResult struct {
	Error error
	Data  string
}

// 真实工作设备,即具体实现
type XDevice interface {
	// 初始化 通常用来获取设备的配置
	Init(devId string, configMap map[string]any) error
	// 启动, 设备的工作进程
	Start(CCTX) error
	// 新特性, 适用于自定义协议读写
	OnCtrl(cmd []byte, args []byte) ([]byte, error)
	// 设备当前状态
	Status() SourceState
	// 停止设备, 在这里释放资源,一般是先置状态为STOP,然后CancelContext()
	Stop()

	Details() *Device
	// 状态
	SetState(SourceState)
	// 外部调用, 该接口是个高级功能, 准备为了设计分布式部署设备的时候用, 但是相当长时间内都不会开启
	// 默认情况下该接口没有用
	OnDCACall(UUID string, Command string, Args any) DCAResult
}

/*
*
* 子设备网络拓扑[2023-04-17新增]
*
 */
type DeviceTopology struct {
	Id       string         // 子设备的ID
	Name     string         // 子设备名
	LinkType int            // 物理连接方式: 0-ETH 1-WIFI 3-BLE 4 LORA 5 OTHER
	State    int            // 状态: 0-Down 1-Working
	Info     map[string]any // 子设备的一些额外信息
}
