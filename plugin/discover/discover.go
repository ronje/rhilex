package discover

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"gopkg.in/ini.v1"
)

// ProbeMessage 探针协议消息结构
type ProbeMessage struct {
	MessageType string `json:"message_type"`
	NodeID      string `json:"node_id"`
	Token       string `json:"token"`
}

func (pm *ProbeMessage) String() string {
	return fmt.Sprintf("ProbeMessage Type: %s, NodeID: %s, Token: %s", pm.MessageType, pm.NodeID, pm.Token)
}

// Node 节点信息
type Node struct {
	NodeID string
	Addr   *net.UDPAddr
}

// to string
func (n *Node) String() string {
	return fmt.Sprintf("NodeID: %s, Addr: %s", n.NodeID, n.Addr.String())
}

// NodeList 节点列表
type NodeList struct {
	nodes map[string]*Node
	mu    sync.RWMutex
}

// NewNodeList 初始化节点列表
func NewNodeList() *NodeList {
	return &NodeList{
		nodes: make(map[string]*Node),
	}
}

// AddNode 添加节点
func (nl *NodeList) AddNode(node *Node) {
	nl.mu.Lock()
	defer nl.mu.Unlock()
	nl.nodes[node.NodeID] = node
}

// RemoveNode 删除节点
func (nl *NodeList) RemoveNode(nodeID string) {
	nl.mu.Lock()
	defer nl.mu.Unlock()
	delete(nl.nodes, nodeID)
}

// GetNodes 获取节点列表
func (nl *NodeList) GetNodes() []*Node {
	nl.mu.RLock()
	defer nl.mu.RUnlock()
	var nodes []*Node
	for _, node := range nl.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// DiscoverPlugin 实现XPlugin接口的发现插件
type DiscoverPlugin struct {
	config            *ini.Section
	nodeList          *NodeList
	nodeName          string
	token             string
	broadcastInterval time.Duration
	udpPort           int
	conn              *net.UDPConn
	stopChan          chan struct{}
	metaInfo          typex.XPluginMetaInfo
	broadcast         bool
}

func NewDiscoverPlugin() *DiscoverPlugin {
	return &DiscoverPlugin{}
}

// Init 初始化插件
func (dp *DiscoverPlugin) Init(config *ini.Section) error {
	dp.config = config
	dp.nodeList = NewNodeList()

	// 解析配置信息
	dp.nodeName = config.Key("node_name").MustString("rhilex@local.node")
	dp.token = config.Key("token").MustString("rhilex_secret_token")
	dp.broadcastInterval = time.Duration(config.Key("broadcast_interval").MustInt(5)) * time.Second
	dp.udpPort = config.Key("udp_port").MustInt(2590)
	dp.stopChan = make(chan struct{})
	dp.metaInfo = typex.XPluginMetaInfo{
		UUID:        "discover",
		Name:        "Network Discover Plugin",
		Version:     "v1.0.0",
		Description: "Discovering nodes in the local network",
	}
	return nil
}

// Start 启动插件
func (dp *DiscoverPlugin) Start(rhilex typex.Rhilex) error {
	// 监听UDP端口
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", dp.udpPort))
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %v", err)
	}
	dp.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("error listening on UDP: %v", err)
	}

	// 启动广播
	dp.broadcast = true
	go dp.broadcastProbeMessages()

	// 接收消息
	go dp.receiveMessages()

	return nil
}

// Service 对外提供服务
func (dp *DiscoverPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	// 这里可以实现具体的服务逻辑
	return typex.ServiceResult{}
}

// Stop 停止插件
func (dp *DiscoverPlugin) Stop() error {
	dp.broadcast = false
	close(dp.stopChan)
	if dp.conn != nil {
		return dp.conn.Close()
	}
	return nil
}

// PluginMetaInfo 返回插件元信息
func (dp *DiscoverPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return dp.metaInfo
}

// broadcastProbeMessages 定期发送探针消息
func (dp *DiscoverPlugin) broadcastProbeMessages() {
	ticker := time.NewTicker(dp.broadcastInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !dp.broadcast {
				return
			}
			probeMsg := ProbeMessage{
				MessageType: "probe",
				NodeID:      dp.nodeName,
				Token:       dp.token,
			}
			glogger.GLogger.Debug("Broadcast Probe Messages:", probeMsg.String())
			probeData, err := json.Marshal(probeMsg)
			if err != nil {
				glogger.GLogger.Errorf("Error marshalling probe message: %v", err)
				continue
			}

			broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:8888")
			if err != nil {
				glogger.GLogger.Errorf("Error resolving broadcast address: %v", err)
				continue
			}

			_, err = dp.conn.WriteToUDP(probeData, broadcastAddr)
			if err != nil {
				glogger.GLogger.Errorf("Error sending probe message: %v", err)
			}
		case <-dp.stopChan:
			return
		}
	}
}

// receiveMessages 接收消息
func (dp *DiscoverPlugin) receiveMessages() {
	buf := make([]byte, 1024)
	for {
		n, addr, err := dp.conn.ReadFromUDP(buf)
		if err != nil {
			if dp.broadcast {
				glogger.GLogger.Errorf("Error reading from UDP: %v", err)
			}
			return
		}
		go dp.handleMessage(buf[:n], addr)
	}
}

// handleMessage 处理接收到的消息
func (dp *DiscoverPlugin) handleMessage(data []byte, addr *net.UDPAddr) {
	var msg ProbeMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		glogger.GLogger.Errorf("Error Unmarshal message from %s: %v", addr.String(), err)
		return
	}

	switch msg.MessageType {
	case "probe":
		// 验证Token
		if msg.Token != dp.token {
			glogger.GLogger.Errorf("Invalid token from %s", addr.String())
			return
		}
		// 回复确认加入
		responseMsg := ProbeMessage{
			MessageType: "response",
			NodeID:      dp.nodeName,
			Token:       dp.token,
		}
		responseData, err := json.Marshal(responseMsg)
		if err != nil {
			glogger.GLogger.Errorf("Error marshalling response message: %v", err)
			return
		}
		_, err = dp.conn.WriteToUDP(responseData, addr)
		if err != nil {
			glogger.GLogger.Errorf("Error sending response to %s: %v", addr.String(), err)
		}
	case "response":
		// 验证Token
		if msg.Token != dp.token {
			glogger.GLogger.Errorf("Invalid token from %s", addr.String())
			return
		}
		// 将对方加入节点列表
		dp.nodeList.AddNode(&Node{
			NodeID: msg.NodeID,
			Addr:   addr,
		})
		glogger.GLogger.Infof("Added node %s at %s to the node list", msg.NodeID, addr.String())
	}
}
