package source

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net"
	"unicode/utf8"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type RHILEXUdpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址"`
	Port          int    `json:"port" validate:"required" title:"服务端口"`
	MaxDataLength int    `json:"maxDataLength" title:"最大数据包"`
}
type udpSource struct {
	typex.XStatus
	UdpConn    *net.UDPConn
	mainConfig RHILEXUdpConfig
	status     typex.SourceState
}

func NewUdpInEndSource(e typex.Rhilex) typex.XSource {
	udps := udpSource{}
	udps.mainConfig = RHILEXUdpConfig{
		Host:          "0.0.0.0",
		Port:          6200,
		MaxDataLength: 1024,
	}
	udps.RuleEngine = e
	return &udps
}

func (udps *udpSource) Start(cctx typex.CCTX) error {
	udps.Ctx = cctx.Ctx
	udps.CancelCTX = cctx.CancelCTX

	addr := &net.UDPAddr{IP: net.ParseIP(udps.mainConfig.Host), Port: udps.mainConfig.Port}
	var err error
	if udps.UdpConn, err = net.ListenUDP("udp", addr); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	go func(ctx context.Context, u1 *udpSource) {
		buffer := make([]byte, udps.mainConfig.MaxDataLength)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, remoteAddr, err := u1.UdpConn.ReadFromUDP(buffer)
			if err != nil {
				glogger.GLogger.Error(err.Error())
				continue
			}
			glogger.GLogger.Debug("UDP Server Received:", buffer[:n])
			go udps.handleClient(buffer[:n], remoteAddr)
		}
	}(udps.Ctx, udps)
	udps.status = typex.SOURCE_UP
	glogger.GLogger.Infof("UDP source started on [%v]:%v", udps.mainConfig.Host, udps.mainConfig.Port)
	return nil

}

/*
*
* 处理UDP客户端
*
 */
func (udps *udpSource) handleClient(data []byte, remoteAddr *net.UDPAddr) {
	ClientData := udp_client_data{
		ClientAddr: remoteAddr.String(),
	}
	if utf8.Valid(data) {
		ClientData.Data = string(data)
	} else {
		ClientData.Data = hex.EncodeToString(data)
	}
	ClientDataBytes, _ := json.Marshal(ClientData)
	work, err := udps.RuleEngine.WorkInEnd(udps.RuleEngine.GetInEnd(udps.PointId), string(ClientDataBytes))
	if !work {
		glogger.GLogger.Error(err)
	}
	// return ok
	_, err = udps.UdpConn.WriteToUDP([]byte("ok"), remoteAddr)
	if err != nil {
		glogger.GLogger.Error(err)
	}
}

/*
*
* 对外输出的数据格式
*
 */
type udp_client_data struct {
	ClientAddr string      `json:"clientAddr"`
	Data       interface{} `json:"data"`
}

func (O udp_client_data) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}
func (udps *udpSource) Details() *typex.InEnd {
	return udps.RuleEngine.GetInEnd(udps.PointId)
}

func (udps *udpSource) Test(inEndId string) bool {
	return true
}

func (udps *udpSource) Init(inEndId string, configMap map[string]interface{}) error {
	udps.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &udps.mainConfig); err != nil {
		return err
	}
	return nil
}

func (udps *udpSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (udps *udpSource) Stop() {
	udps.status = typex.SOURCE_DOWN
	if udps.CancelCTX != nil {
		udps.CancelCTX()
	}
	if udps.UdpConn != nil {
		err := udps.UdpConn.Close()
		if err != nil {
			glogger.GLogger.Error(err)
		}
	}

}

// 来自外面的数据
func (*udpSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*udpSource) UpStream([]byte) (int, error) {
	return 0, nil
}
