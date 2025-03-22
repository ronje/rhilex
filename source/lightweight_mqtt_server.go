package source

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type MqttServerConfig struct {
	ServerName string `json:"serverName" validate:"required" title:"服务名称"`
	ListenHost string `json:"host" validate:"required" title:"监听地址"`
	ListenPort int    `json:"port" validate:"required" title:"监听端口"`
	Anonymous  *bool  `json:"anonymous" validate:"required" title:"允许匿名连接"`
}
type MqttServer struct {
	typex.XStatus
	mainConfig MqttServerConfig
	status     typex.SourceState
	server     *mqtt.Server
	locker     sync.Mutex
}

/*
*
* 轻量级MQTT Server
*
 */
func NewMqttServer(e typex.Rhilex) typex.XSource {
	Anonymous := true
	h := MqttServer{
		mainConfig: MqttServerConfig{
			ServerName: "rhilex-mqtt-server",
			ListenHost: "127.0.0.1",
			ListenPort: 1883,
			Anonymous:  &Anonymous,
		},
		locker: sync.Mutex{},
	}
	h.RuleEngine = e
	return &h
}

func (ms *MqttServer) Init(inEndId string, configMap map[string]any) error {
	ms.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &ms.mainConfig); err != nil {
		return err
	}
	return nil
}

func (ms *MqttServer) Start(cctx typex.CCTX) error {
	ms.Ctx = cctx.Ctx
	ms.CancelCTX = cctx.CancelCTX
	Slog := slog.New(slog.NewTextHandler(glogger.Logrus.Out, &slog.HandlerOptions{
		Level: func() slog.Leveler {
			if core.GlobalConfig.LogLevel == "info" {
				return slog.LevelInfo
			}
			if core.GlobalConfig.LogLevel == "debug" {
				return slog.LevelDebug
			}
			if core.GlobalConfig.LogLevel == "error" {
				return slog.LevelError
			}
			if core.GlobalConfig.LogLevel == "warn" {
				return slog.LevelWarn
			}
			return slog.LevelInfo
		}(),
	}))
	ms.server = mqtt.New(&mqtt.Options{
		Logger: Slog,
	})
	ms.server.AddHook(&AuthHook{s: ms}, nil)
	if err := ms.server.AddListener(listeners.NewTCP(listeners.Config{
		ID:      ms.mainConfig.ServerName,
		Address: fmt.Sprintf("%s:%d", ms.mainConfig.ListenHost, ms.mainConfig.ListenPort),
	})); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if err := ms.server.Serve(); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	ms.status = typex.SOURCE_UP
	return nil
}

func (ms *MqttServer) Stop() {
	ms.status = typex.SOURCE_DOWN
	if ms.CancelCTX != nil {
		ms.CancelCTX()
	}
	if ms.server != nil {
		ms.server.Close()
		ms.server = nil
	}
}

func (ms *MqttServer) Status() typex.SourceState {
	if ms.server == nil {
		return typex.SOURCE_DOWN
	}
	return ms.status
}

func (ms *MqttServer) Details() *typex.InEnd {
	return ms.RuleEngine.GetInEnd(ms.PointId)
}

func (*MqttServer) DownStream([]byte) (int, error) {
	return 0, nil
}

func (*MqttServer) UpStream([]byte) (int, error) {
	return 0, nil
}

func (ms *MqttServer) Clients(page, size int) *mqtt.Clients {
	return ms.server.Clients
}
func (ms *MqttServer) FindClients(clientId string) (*mqtt.Client, bool) {
	return ms.server.Clients.Get(clientId)
}

// AuthHooks is an authentication hook which allows connection access
// for all users and read and write access to all topics.
type AuthHook struct {
	mqtt.HookBase
	s *MqttServer
}

// ID returns the ID of the hook.
func (h *AuthHook) ID() string {
	return "auth-hooks"
}

/*
*
* OnConnectAuthenticate
*
 */
func (h *AuthHook) OnConnectAuthenticate(C *mqtt.Client, pk packets.Packet) bool {
	if *h.s.mainConfig.Anonymous {
		return true
	}
	return h.check(C.ID, string(C.Properties.Username), "")
}

// TODO
func (h *AuthHook) check(string, string, string) bool {
	return true
}

// Provides indicates which hook methods this hook provides.
func (h *AuthHook) Provides(b byte) bool {
	return true
}

/*
*
* OnACLCheck
*
 */
func (h *AuthHook) OnACLCheck(client *mqtt.Client, topic string, write bool) bool {
	glogger.GLogger.Debugf("Mqtt Client ACLCheck, ClientId:(%s),Topic: (%v)", client.ID, topic)
	return true
}

/*
*
* OnSubscribe
*
 */
func (h *AuthHook) OnSubscribe(client *mqtt.Client, pk packets.Packet) packets.Packet {
	glogger.GLogger.Debugf("Mqtt Client Subscribe, ClientId:(%s),Topic: (%v)", client.ID, pk.TopicName)
	return pk
}

/*
*
* OnConnect
*
 */
func (h *AuthHook) OnConnect(client *mqtt.Client, pk packets.Packet) error {
	glogger.GLogger.Debugf("Mqtt Client Connected:(%s), Addr:(%s)", client.ID, client.Net.Conn.RemoteAddr())
	MqttEvent := MqttEvent{
		Action:    "connect",
		Clientid:  client.ID,
		Username:  string(client.Properties.Username),
		Ipaddress: client.Net.Remote,
		Ts:        time.Now().UnixMilli(),
		Topic:     pk.TopicName,
		Payload:   string(pk.Payload),
	}
	_, errWorkInEnd := h.s.RuleEngine.WorkInEnd(
		h.s.RuleEngine.GetInEnd(h.s.PointId),
		MqttEvent.JsonString(),
	)
	if errWorkInEnd != nil {
		glogger.GLogger.Error(errWorkInEnd)
	}
	return nil
}

/*
*
* OnDisconnect
*
 */
func (h *AuthHook) OnDisconnect(client *mqtt.Client, err error, expire bool) {
	glogger.GLogger.Debugf("Mqtt Client Disconnect:(%s), Addr:(%s)", client.ID, client.Net.Conn.RemoteAddr())
	MqttEvent := MqttEvent{
		Action:    "disconnect",
		Clientid:  client.ID,
		Username:  string(client.Properties.Username),
		Ipaddress: client.Net.Remote,
		Ts:        time.Now().UnixMilli(),
	}
	_, errWorkInEnd := h.s.RuleEngine.WorkInEnd(
		h.s.RuleEngine.GetInEnd(h.s.PointId),
		MqttEvent.JsonString(),
	)
	if errWorkInEnd != nil {
		glogger.GLogger.Error(errWorkInEnd)
	}
}

/*
*
* OnPublish
*
 */
func (h *AuthHook) OnPublish(client *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	glogger.GLogger.Debugf("Mqtt Client Publish:(%s), Addr:(%s), Topic:(%s), Payload(%s)",
		client.ID, client.Net.Conn.RemoteAddr(), pk.TopicName, string(pk.Payload))
	MqttEvent := MqttEvent{
		Action:    "publish",
		Clientid:  client.ID,
		Username:  string(client.Properties.Username),
		Ipaddress: client.Net.Remote,
		Ts:        time.Now().UnixMilli(),
		Topic:     pk.TopicName,
		Payload:   string(pk.Payload),
	}
	_, errWorkInEnd := h.s.RuleEngine.WorkInEnd(
		h.s.RuleEngine.GetInEnd(h.s.PointId),
		MqttEvent.JsonString(),
	)
	if errWorkInEnd != nil {
		glogger.GLogger.Error(errWorkInEnd)
	}
	return pk, nil
}

/*
*
* Push to rule
*
 */
type MqttEvent struct {
	Action    string `json:"action"`
	Clientid  string `json:"clientid"`
	Username  string `json:"username"`
	Ipaddress string `json:"ipaddress"`
	Ts        int64  `json:"ts"`
	Topic     string `json:"topic,omitempty"`
	Payload   string `json:"payload,omitempty"`
}

func (O MqttEvent) JsonString() string {
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}
