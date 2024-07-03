package mqttserver

import (
	"fmt"
	"sync"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"gopkg.in/ini.v1"
)

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"host"`
	Port   int    `ini:"port"`
}

type MqttServer struct {
	Enable     bool
	Host       string
	Port       int
	mqttServer *mqtt.Server
	topics     map[string][]_topic // Topic 订阅表
	ruleEngine typex.Rhilex
	uuid       string
	locker     sync.Mutex
}

func NewMqttServer() typex.XPlugin {
	return &MqttServer{
		Host:   "127.0.0.1",
		Port:   1883,
		topics: map[string][]_topic{},
		uuid:   "RHILEX-MqttServer",
		locker: sync.Mutex{},
	}
}

func (s *MqttServer) Init(config *ini.Section) error {
	var mainConfig _serverConfig
	if err := utils.InIMapToStruct(config, &mainConfig); err != nil {
		return err
	}
	s.Host = mainConfig.Host
	s.Port = mainConfig.Port
	return nil
}

func (s *MqttServer) Start(r typex.Rhilex) error {
	s.ruleEngine = r

	server := mqtt.New(&mqtt.Options{})
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "node1",
		Address: fmt.Sprintf("%v:%v", s.Host, s.Port),
	})
	if err := server.AddListener(tcp); err != nil {
		return err
	}
	if err := server.Serve(); err != nil {
		return err
	}
	//
	// 本地服务器
	//
	s.mqttServer = server
	server.AddHook(&AuthHook{s: s}, nil)
	glogger.GLogger.Infof("MqttServer start at [%s:%v] successfully", s.Host, s.Port)
	return nil
}

func (s *MqttServer) Stop() error {
	if s.mqttServer != nil {
		return s.mqttServer.Close()
	} else {
		return nil
	}

}

func (s *MqttServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        s.uuid,
		Name:        "MqttServer",
		Version:     "v2.0.0",
		Description: "Simple Light Weight MqttServer",
	}
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
func (h *AuthHook) OnConnectAuthenticate(C *mqtt.Client, pk packets.Packet) bool {
	return true
}

// Provides indicates which hook methods this hook provides.
func (h *AuthHook) Provides(b byte) bool {
	return true
}
func (h *AuthHook) OnConnect(C *mqtt.Client, pk packets.Packet) error {
	glogger.GLogger.Infof("Mqtt Client Connected:(%s), Addr:(%s)", C.ID, C.Net.Conn.RemoteAddr())
	return nil
}
func (h *AuthHook) OnDisconnect(C *mqtt.Client, err error, expire bool) {
	h.s.locker.Lock()
	defer h.s.locker.Unlock()
	delete(h.s.topics, C.ID)
	glogger.GLogger.Infof("Mqtt Client Disconnect:(%s), Addr:(%s)", C.ID, C.Net.Conn.RemoteAddr())
}

// OnACLCheck returns true/allowed for all checks.
func (h *AuthHook) OnACLCheck(client *mqtt.Client, topic string, write bool) bool {
	glogger.GLogger.Infof("Mqtt Client ACLCheck, ClientId:(%s),Topic: (%v)", client.ID, topic)
	h.s.locker.Lock()
	defer h.s.locker.Unlock()
	_, ok := h.s.topics[client.ID]
	if !ok {
		h.s.topics[client.ID] = []_topic{{Topic: topic}}
	} else {
		h.s.topics[client.ID] = append(h.s.topics[client.ID], _topic{Topic: topic})
	}
	return true
}
