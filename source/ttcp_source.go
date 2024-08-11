package source

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net"
	"time"
	"unicode/utf8"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type RHILEXTcpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址"`
	Port          int    `json:"port" validate:"required" title:"服务端口"`
	MaxDataLength int    `json:"maxDataLength" title:"最大数据包"`
}
type TcpSource struct {
	typex.XStatus
	tCPListener *net.TCPListener
	mainConfig  RHILEXTcpConfig
	status      typex.SourceState
}

func NewTcpSource(e typex.Rhilex) typex.XSource {
	tcps := TcpSource{}
	tcps.mainConfig = RHILEXTcpConfig{
		Host:          "0.0.0.0",
		Port:          6201,
		MaxDataLength: 1024,
	}
	tcps.RuleEngine = e
	return &tcps
}

func (tcps *TcpSource) Start(cctx typex.CCTX) error {
	tcps.Ctx = cctx.Ctx
	tcps.CancelCTX = cctx.CancelCTX
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d",
		tcps.mainConfig.Host, tcps.mainConfig.Port))
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if tcps.tCPListener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	tcps.status = typex.SOURCE_UP
	go func(ctx context.Context, tcps *TcpSource) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			conn, err := tcps.tCPListener.Accept()
			if err != nil {
				glogger.GLogger.Error(err)
				continue
			}
			glogger.GLogger.Debug("Tcp Client Connected:", conn.RemoteAddr())
			go tcps.handleClient(conn)
		}
	}(tcps.Ctx, tcps)
	tcps.status = typex.SOURCE_UP
	glogger.GLogger.Infof("UDP source started on [%v]:%v", tcps.mainConfig.Host, tcps.mainConfig.Port)
	return nil

}

// calculateCRC 计算给定数据的 CRC32 校验值并返回 4 字节的结果
func calculateCRC(data []byte) [4]byte {
	crcTable := crc32.MakeTable(crc32.IEEE)
	crcValue := crc32.Checksum(data, crcTable)
	return [4]byte{
		byte(crcValue >> 24),
		byte(crcValue >> 16),
		byte(crcValue >> 8),
		byte(crcValue),
	}
}

const (
	headerSize = 4
	bufferSize = 1024 // 设置缓冲区大小
)

// 处理客户端连接
func (tcps *TcpSource) handleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, bufferSize)
	header := make([]byte, headerSize)
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, err := conn.Read(header); err != nil {
			glogger.GLogger.Error(err)
			return
		}
		conn.SetReadDeadline(time.Time{})
		length := int(header[0])
		if length > bufferSize {
			glogger.GLogger.Error("exceed max buffer size")
			return
		}
		data := buffer[:length]
		if _, err := conn.Read(data); err != nil {
			glogger.GLogger.Error(err)
			return
		}
		// crc := header[1:4]
		// calculatedCRC := calculateCRC(data)
		// if string(calculatedCRC[:]) != string(crc) {
		// 	glogger.GLogger.Error("CRC check failed")
		// 	return
		// }
		glogger.GLogger.Debug("TCP Server Received:", data)
		tcpClientData := udp_client_data{
			ClientAddr: conn.RemoteAddr().String(),
		}
		if utf8.Valid(data) {
			tcpClientData.Data = string(data)
		} else {
			tcpClientData.Data = hex.EncodeToString(data)
		}

		tcpClientDataBytes, _ := json.Marshal(tcpClientData)
		work, err := tcps.RuleEngine.WorkInEnd(tcps.RuleEngine.GetInEnd(tcps.PointId),
			string(tcpClientDataBytes))
		if !work {
			glogger.GLogger.Error(err)
			continue
		}
		// return ok
		_, err = conn.Write([]byte("ok"))
		if err != nil {
			glogger.GLogger.Error(err)
		}
	}
}

/*
*
* 对外输出的数据格式
*
 */
type tcp_client_data struct {
	ClientAddr string      `json:"clientAddr"`
	Data       interface{} `json:"data"`
}

func (tcps *TcpSource) Details() *typex.InEnd {
	return tcps.RuleEngine.GetInEnd(tcps.PointId)
}

func (tcps *TcpSource) Test(inEndId string) bool {
	return true
}

func (tcps *TcpSource) Init(inEndId string, configMap map[string]interface{}) error {
	tcps.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &tcps.mainConfig); err != nil {
		return err
	}
	return nil
}

func (tcps *TcpSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (tcps *TcpSource) Stop() {
	tcps.status = typex.SOURCE_DOWN
	if tcps.CancelCTX != nil {
		tcps.CancelCTX()
	}
	if tcps.tCPListener != nil {
		err := tcps.tCPListener.Close()
		if err != nil {
			glogger.GLogger.Error(err)
		}
	}
}

// 来自外面的数据
func (*TcpSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*TcpSource) UpStream([]byte) (int, error) {
	return 0, nil
}
