package source

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"unicode/utf8"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// RHILEXUdpConfig 定义UDP服务的配置信息
type RHILEXUdpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址"`
	Port          int    `json:"port" validate:"required" title:"服务端口"`
	MaxDataLength int    `json:"maxDataLength" title:"最大数据包"`
}

// udpSource 表示一个UDP数据源
type udpSource struct {
	typex.XStatus
	UdpConn    *net.UDPConn
	mainConfig RHILEXUdpConfig
	status     typex.SourceState
}

// NewUdpInEndSource 创建一个新的UDP数据源实例
func NewUdpInEndSource(e typex.Rhilex) typex.XSource {
	udps := udpSource{
		mainConfig: RHILEXUdpConfig{
			Host:          "0.0.0.0",
			Port:          6200,
			MaxDataLength: 1024,
		},
	}
	udps.RuleEngine = e
	return &udps
}

// Init 初始化UDP数据源，绑定配置信息
func (udps *udpSource) Init(inEndId string, configMap map[string]interface{}) error {
	udps.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &udps.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind UDP source config: %v", err)
		return err
	}
	// 验证配置
	if err := udps.validateConfig(); err != nil {
		glogger.GLogger.Errorf("Invalid UDP source config: %v", err)
		return err
	}
	return nil
}

// validateConfig 验证UDP配置信息的有效性
func (udps *udpSource) validateConfig() error {
	if udps.mainConfig.Port <= 0 || udps.mainConfig.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", udps.mainConfig.Port)
	}
	if udps.mainConfig.MaxDataLength <= 0 {
		return fmt.Errorf("invalid max data length: %d", udps.mainConfig.MaxDataLength)
	}
	if net.ParseIP(udps.mainConfig.Host) == nil {
		return fmt.Errorf("invalid host address: %s", udps.mainConfig.Host)
	}
	return nil
}

// Start 启动UDP数据源，开始监听UDP端口
func (udps *udpSource) Start(cctx typex.CCTX) error {
	udps.Ctx = cctx.Ctx
	udps.CancelCTX = cctx.CancelCTX

	addr := &net.UDPAddr{IP: net.ParseIP(udps.mainConfig.Host), Port: udps.mainConfig.Port}
	var err error
	udps.UdpConn, err = net.ListenUDP("udp", addr)
	if err != nil {
		glogger.GLogger.Errorf("Failed to listen on UDP address: %v", err)
		return err
	}

	go udps.listenForUDP()

	udps.status = typex.SOURCE_UP
	glogger.GLogger.Infof("UDP source started on [%v]:%v", udps.mainConfig.Host, udps.mainConfig.Port)
	return nil
}

// listenForUDP 监听UDP端口，处理接收到的数据
func (udps *udpSource) listenForUDP() {
	buffer := make([]byte, udps.mainConfig.MaxDataLength)
	for {
		select {
		case <-udps.Ctx.Done():
			return
		default:
		}

		n, remoteAddr, err := udps.UdpConn.ReadFromUDP(buffer)
		if err != nil {
			glogger.GLogger.Errorf("Error reading UDP data: %v", err)
			continue
		}

		glogger.GLogger.Debugf("UDP Server Received from %s: %s", remoteAddr.String(), hex.EncodeToString(buffer[:n]))

		go udps.handleClient(buffer[:n], remoteAddr)
	}
}

// handleClient 处理接收到的UDP客户端数据
func (udps *udpSource) handleClient(data []byte, remoteAddr *net.UDPAddr) {
	clientData := udp_client_data{
		ClientAddr: remoteAddr.String(),
		Data:       hex.EncodeToString(data),
	}

	if utf8.Valid(data) {
		glogger.GLogger.Debugf("UDP Server Received valid UTF-8 data from %s: %s", remoteAddr.String(), string(data))
	} else {
		glogger.GLogger.Debugf("UDP Server Received non-UTF-8 data from %s: %s", remoteAddr.String(), hex.EncodeToString(data))
	}

	clientDataBytes, err := json.Marshal(clientData)
	if err != nil {
		glogger.GLogger.Errorf("Failed to marshal client data: %v", err)
		return
	}

	work, err := udps.RuleEngine.WorkInEnd(udps.RuleEngine.GetInEnd(udps.PointId), string(clientDataBytes))
	if err != nil || !work {
		glogger.GLogger.Errorf("Failed to process client data: %v", err)
	}

	// Send response
	response := []byte("ok\r\n")
	if _, err := udps.UdpConn.WriteToUDP(response, remoteAddr); err != nil {
		glogger.GLogger.Errorf("Failed to send response to client: %v", err)
	}
}

// udp_client_data 定义对外输出的数据格式
type udp_client_data struct {
	ClientAddr string      `json:"clientAddr"`
	Data       interface{} `json:"data"`
}

// String 将udp_client_data转换为JSON字符串
func (o udp_client_data) String() string {
	if bytes, err := json.Marshal(o); err != nil {
		glogger.GLogger.Errorf("Failed to marshal udp_client_data: %v", err)
		return ""
	} else {
		return string(bytes)
	}
}

// Details 获取UDP数据源的详细信息
func (udps *udpSource) Details() *typex.InEnd {
	return udps.RuleEngine.GetInEnd(udps.PointId)
}

// Status 获取UDP数据源的当前状态
func (udps *udpSource) Status() typex.SourceState {
	return udps.status
}

// Stop 停止UDP数据源，释放资源
func (udps *udpSource) Stop() {
	udps.status = typex.SOURCE_DOWN
	if udps.CancelCTX != nil {
		udps.CancelCTX()
	}
	if udps.UdpConn != nil {
		if err := udps.UdpConn.Close(); err != nil {
			glogger.GLogger.Errorf("Failed to close UDP connection: %v", err)
		}
	}
}
