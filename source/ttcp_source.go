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
	Host          string `json:"host" validate:"required"`
	Port          int    `json:"port" validate:"required"`
	MaxDataLength int    `json:"maxDataLength"`
	KeepAlive     int    `json:"keepAlive"` // 客户端保活时间：ms
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
		KeepAlive:     5000,
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
	glogger.GLogger.Infof("TCP source started on [%v]:%v", tcps.mainConfig.Host, tcps.mainConfig.Port)
	return nil

}

// calculateCRC 计算给定数据的 CRC32 校验值并返回 4 字节的结果
// 00 01 A9 C7 E8 B8 01
var __crcIEEETable = crc32.MakeTable(crc32.Koopman)

func calculateCRC(data []byte) uint32 {
	crcValue := crc32.Checksum(data, __crcIEEETable)
	return ByteToUint32([]byte{
		byte(crcValue >> 24),
		byte(crcValue >> 16),
		byte(crcValue >> 8),
		byte(crcValue),
	})
}

/*
*
* 大端转换法
*
 */
func ByteToUint32(b []byte) uint32 {
	var result uint32
	for i := 0; i < 4; i++ {
		result |= uint32(b[i]) << (8 * (3 - i))
	}
	return result
}
func ByteToUint16(b []byte) uint16 {
	var result uint16
	result |= uint16(b[0]) << 8
	result |= uint16(b[1])
	return result
}

const (
	headerSize = 6    // LENGTH|CRC4 CRC3 CRC2 CRC1
	bufferSize = 1024 // 设置缓冲区大小
)

// 处理客户端连接
func (tcps *TcpSource) handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, headerSize)
	buffer := make([]byte, bufferSize)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(tcps.mainConfig.KeepAlive) * time.Millisecond))
		if _, err := conn.Read(header); err != nil {
			glogger.GLogger.Error(err)
			return
		}
		conn.SetReadDeadline(time.Time{})
		length := ByteToUint16([]byte{header[0], header[1]})
		if length > bufferSize {
			glogger.GLogger.Error("exceed max buffer size")
			return
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		N, errRead := conn.Read(buffer[:length])
		if errRead != nil {
			glogger.GLogger.Error(errRead)
			return
		}
		conn.SetReadDeadline(time.Time{})
		data := buffer[:N]
		calculatedCRC := calculateCRC(data)
		dataCrc := ByteToUint32([]byte{header[2], header[3], header[4], header[5]})
		if calculatedCRC != dataCrc {
			glogger.GLogger.Error("CRC check failed")
			return
		}
		ClientData := tcp_client_data{
			ClientAddr: conn.RemoteAddr().String(),
		}
		if utf8.Valid(data) {
			ClientData.Data = string(data)
		} else {
			ClientData.Data = hex.EncodeToString(data)
		}
		glogger.GLogger.Debug("TCP Server Received:", ClientData)
		tcpClientDataBytes, _ := json.Marshal(ClientData)
		work, err := tcps.RuleEngine.WorkInEnd(tcps.RuleEngine.GetInEnd(tcps.PointId),
			string(tcpClientDataBytes))
		if !work {
			glogger.GLogger.Error(err)
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

func (O tcp_client_data) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}
func (tcps *TcpSource) Details() *typex.InEnd {
	return tcps.RuleEngine.GetInEnd(tcps.PointId)
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
