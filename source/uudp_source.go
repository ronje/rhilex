package source

import (
	"context"
	"net"

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
	uDPConn    *net.UDPConn
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
	if udps.uDPConn, err = net.ListenUDP("udp", addr); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	udps.status = typex.SOURCE_UP
	go func(ctx context.Context, u1 *udpSource) {
		data := make([]byte, udps.mainConfig.MaxDataLength)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, remoteAddr, err := u1.uDPConn.ReadFromUDP(data)
			if err != nil {
				glogger.GLogger.Error(err.Error())
				continue
			}
			glogger.GLogger.Debug("UDP Server Received:", data[:n])
			work, err := udps.RuleEngine.WorkInEnd(udps.RuleEngine.GetInEnd(udps.PointId), string(data[:n]))
			if !work {
				glogger.GLogger.Error(err)
				continue
			}
			// return ok
			_, err = u1.uDPConn.WriteToUDP([]byte("ok"), remoteAddr)
			if err != nil {
				glogger.GLogger.Error(err)
			}
		}
	}(udps.Ctx, udps)
	udps.status = typex.SOURCE_UP
	glogger.GLogger.Infof("UDP source started on [%v]:%v", udps.mainConfig.Host, udps.mainConfig.Port)
	return nil

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
	if udps.uDPConn != nil {
		err := udps.uDPConn.Close()
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
