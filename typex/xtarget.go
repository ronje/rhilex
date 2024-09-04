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
	NATS_TARGET           TargetType = "NATS"                  // To Nats.io
	HTTP_TARGET           TargetType = "HTTP"                  // To Http Target
	TDENGINE_TARGET       TargetType = "TDENGINE"              // To TDENGINE
	GRPC_CODEC_TARGET     TargetType = "GRPC_CODEC_TARGET"     // To GRPC Target
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
	Init(outEndId string, configMap map[string]interface{}) error
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
	To(data interface{}) (interface{}, error)
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}
