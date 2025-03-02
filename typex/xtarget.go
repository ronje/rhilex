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

// TargetType
type TargetType string

func (i TargetType) String() string {
	return string(i)
}

/*
*
* 输出资源类型
*
 */
const (
	MONGO_SINGLE          TargetType = "MONGO_SINGLE"          // To MongoDB
	MQTT_TARGET           TargetType = "MQTT"                  // To Mqtt Server
	HTTP_TARGET           TargetType = "HTTP"                  // To Http Target
	TDENGINE_TARGET       TargetType = "TDENGINE"              // To TDENGINE
	GRPC_CODEC_TARGET     TargetType = "GRPC_CODEC_TARGET"     // To GRPC Target
	RHILEX_GRPC_TARGET    TargetType = "RHILEX_GRPC_TARGET"    // To GRPC Target
	UDP_TARGET            TargetType = "UDP_TARGET"            // To UDP Server
	GENERIC_UART_TARGET   TargetType = "GENERIC_UART_TARGET"   // To GENERIC_UART_TARGET DTU
	TCP_TRANSPORT         TargetType = "TCP_TRANSPORT"         // To TCP Transport
	SEMTECH_UDP_FORWARDER TargetType = "SEMTECH_UDP_FORWARDER" // To Chirp stack UDP
	GREPTIME_DATABASE     TargetType = "GREPTIME_DATABASE"     // To GREPTIME DATABASE
)

// Stream from source and to target
type XTarget interface {
	//
	// 用来初始化传递资源配置
	//
	Init(outEndId string, configMap map[string]any) error
	//
	// 启动资源
	//
	Start(CCTX) error
	//
	// 获取资源状态
	//
	Status() SourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *OutEnd
	//
	// 数据出口
	//
	To(data any) (any, error)
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}
