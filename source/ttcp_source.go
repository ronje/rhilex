package source

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TcpConnectionManager struct {
	connections map[string]net.Conn
	mu          sync.Mutex
}

func NewConnectionManager() *TcpConnectionManager {
	return &TcpConnectionManager{
		connections: make(map[string]net.Conn),
	}
}

func (cm *TcpConnectionManager) AddConnection(conn net.Conn) string {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	id := uuid.New().String()
	cm.connections[id] = conn
	glogger.GLogger.Info("TcpConnectionManager Add Connection:", id)
	return id
}

func (cm *TcpConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	glogger.GLogger.Info("TcpConnectionManager Remove Connection:", id)
	if conn, ok := cm.connections[id]; ok {
		conn.Close()
		delete(cm.connections, id)
	}
}

type RHILEXTcpConfig struct {
	Host          string `json:"host" validate:"required"`
	Port          int    `json:"port" validate:"required"`
	MaxDataLength int    `json:"maxDataLength"`
	KeepAlive     int    `json:"keepAlive"` // 客户端保活时间：ms
}
type TcpSource struct {
	typex.XStatus
	tCPListener       *net.TCPListener
	mainConfig        RHILEXTcpConfig
	connectionManager *TcpConnectionManager
	status            typex.SourceState
}

func NewTcpSource(e typex.Rhilex) typex.XSource {
	tcps := TcpSource{}
	tcps.mainConfig = RHILEXTcpConfig{
		Host:          "0.0.0.0",
		Port:          6201,
		MaxDataLength: 1024,
		KeepAlive:     5000,
	}
	tcps.connectionManager = NewConnectionManager()
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
			go tcps.handleClient(tcps.connectionManager.AddConnection(conn), conn)
		}
	}(tcps.Ctx, tcps)
	tcps.status = typex.SOURCE_UP
	glogger.GLogger.Infof("TCP source started on [%v]:%v", tcps.mainConfig.Host, tcps.mainConfig.Port)
	return nil

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
	headerSize = 2
	bufferSize = 1024
)

// 处理客户端连接
func (tcps *TcpSource) handleClient(id string, conn net.Conn) {
	defer tcps.connectionManager.RemoveConnection(id)
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
		ClientData := tcp_client_data{
			ClientAddr: conn.RemoteAddr().String(),
			Data:       hex.EncodeToString(data),
		}
		if utf8.Valid(data) {
			glogger.GLogger.Debug("TCP Server Received:", string(data))
		} else {
			glogger.GLogger.Debug("TCP Server Received:", hex.EncodeToString(data))
		}
		tcpClientDataBytes, _ := json.Marshal(ClientData)
		work, err := tcps.RuleEngine.WorkInEnd(tcps.RuleEngine.GetInEnd(tcps.PointId),
			string(tcpClientDataBytes))
		if !work {
			glogger.GLogger.Error(err)
		}
		if _, err = conn.Write([]byte("ok\r\n")); err != nil {
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
	ClientAddr string `json:"clientAddr"`
	Data       any    `json:"data"`
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

func (tcps *TcpSource) Init(inEndId string, configMap map[string]any) error {
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
