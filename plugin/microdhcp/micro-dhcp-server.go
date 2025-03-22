package microdhcp

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/typex"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"gopkg.in/ini.v1"
)

// MicroDHCP 结构体封装了所有的私有参数
type MicroDHCP struct {
	server        *server4.Server
	leases        map[string]net.IP // MAC -> IP 的租约映射
	staticIPs     map[string]net.IP // MAC -> 静态IP 的映射
	blacklist     map[string]bool   // 黑名单MAC地址
	leaseDuration time.Duration     // 租约时长
	listenIP      net.IP
	listenPort    int
	mutex         sync.Mutex // 用于保护并发访问
}

// NewMicroDHCP 创建一个新的 MicroDHCP 实例
func NewMicroDHCP() *MicroDHCP {
	return &MicroDHCP{
		leases:        make(map[string]net.IP),
		staticIPs:     make(map[string]net.IP),
		blacklist:     make(map[string]bool),
		leaseDuration: 24 * time.Hour,
	}
}

// Init 初始化DHCP库，配置通过ini传入
// Init 初始化DHCP库，配置通过ini传入
func (md *MicroDHCP) Init(config *ini.Section) error {
	// 读取监听IP
	listenIPStr := config.Key("listen_ip").String()
	if listenIPStr != "" {
		listenIP := net.ParseIP(listenIPStr)
		if listenIP == nil {
			return errors.New("invalid listen IP address in config")
		}
		md.listenIP = listenIP
	}

	// 读取监听端口
	listenPort, err := config.Key("listen_port").Int()
	if err != nil {
		return fmt.Errorf("invalid listen port in config: %v", err)
	}
	md.listenPort = listenPort

	// 读取租约时长
	leaseDurationStr := config.Key("lease_duration").String()
	if leaseDurationStr != "" {
		leaseDuration, err := time.ParseDuration(leaseDurationStr)
		if err != nil {
			return fmt.Errorf("invalid lease duration in config: %v", err)
		}
		md.leaseDuration = leaseDuration
	}

	// 读取静态IP配置
	staticIPsStr := config.Key("static_ips").String()
	if staticIPsStr != "" {
		staticIPsPairs := strings.Split(staticIPsStr, ",")
		for _, pair := range staticIPsPairs {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				return errors.New("invalid static IP configuration")
			}
			mac := strings.TrimSpace(parts[0])
			ip := net.ParseIP(strings.TrimSpace(parts[1]))
			if ip == nil {
				return errors.New("invalid IP address in static IP configuration")
			}
			md.staticIPs[mac] = ip
		}
	}

	// 读取黑名单配置
	blacklistStr := config.Key("blacklist").String()
	if blacklistStr != "" {
		blacklistMACs := strings.Split(blacklistStr, ",")
		for _, mac := range blacklistMACs {
			md.blacklist[strings.TrimSpace(mac)] = true
		}
	}

	return nil
}

// Start 启动DHCP服务器
func (md *MicroDHCP) Start(rhilex typex.Rhilex) error {
	server, err := server4.NewServer("eth0", nil, md.handleDHCP)
	if err != nil {
		return err
	}
	md.server = server

	go func() {
		if err := md.server.Serve(); err != nil {
			fmt.Printf("DHCP server failed: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止DHCP服务器并释放资源
func (md *MicroDHCP) Stop() error {
	if md.server != nil {
		return md.server.Close()
	}
	return nil
}

// PluginMetaInfo 返回插件的元信息
func (md *MicroDHCP) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "MicroDHCP",
		Name:        "MicroDHCP",
		Version:     "v0.0.1",
		Description: "Micro DHCP Server Plugin",
	}
}

// Service 这个接口是无用的，请不要关注
func (md *MicroDHCP) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

// handleDHCP 处理DHCP请求
func (md *MicroDHCP) handleDHCP(conn net.PacketConn, peer net.Addr, req *dhcpv4.DHCPv4) {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	if md.blacklist[req.ClientHWAddr.String()] {
		fmt.Printf("MAC address %s is blacklisted\n", req.ClientHWAddr.String())
		return
	}

	switch req.MessageType() {
	case dhcpv4.MessageTypeDiscover:
		md.handleDiscover(conn, peer, req)
	case dhcpv4.MessageTypeRequest:
		md.handleRequest(conn, peer, req)
	default:
		fmt.Printf("Unhandled DHCP message type: %v\n", req.MessageType())
	}
}

// handleDiscover 处理DHCP Discover请求
func (md *MicroDHCP) handleDiscover(conn net.PacketConn, peer net.Addr, req *dhcpv4.DHCPv4) {
	ip := md.allocateIP(req.ClientHWAddr.String())
	if ip == nil {
		fmt.Printf("No available IP for MAC %s\n", req.ClientHWAddr.String())
		return
	}

	offer, err := dhcpv4.NewReplyFromRequest(req,
		dhcpv4.WithMessageType(dhcpv4.MessageTypeOffer),
		dhcpv4.WithYourIP(ip),
		dhcpv4.WithServerIP(net.IP{192, 168, 100, 1}),
		dhcpv4.WithLeaseTime(uint32(md.leaseDuration.Seconds())),
	)
	if err != nil {
		fmt.Printf("Failed to create DHCP Offer: %v\n", err)
		return
	}

	if _, err := conn.WriteTo(offer.ToBytes(), peer); err != nil {
		fmt.Printf("Failed to send DHCP Offer: %v\n", err)
	}
}

// handleRequest 处理DHCP Request请求
func (md *MicroDHCP) handleRequest(conn net.PacketConn, peer net.Addr, req *dhcpv4.DHCPv4) {
	ip := md.allocateIP(req.ClientHWAddr.String())
	if ip == nil {
		fmt.Printf("No available IP for MAC %s\n", req.ClientHWAddr.String())
		return
	}

	ack, err := dhcpv4.NewReplyFromRequest(req,
		dhcpv4.WithMessageType(dhcpv4.MessageTypeAck),
		dhcpv4.WithYourIP(ip),
		dhcpv4.WithServerIP(net.IP{192, 168, 100, 1}),
		dhcpv4.WithLeaseTime(uint32(md.leaseDuration.Seconds())),
	)
	if err != nil {
		fmt.Printf("Failed to create DHCP Ack: %v\n", err)
		return
	}

	if _, err := conn.WriteTo(ack.ToBytes(), peer); err != nil {
		fmt.Printf("Failed to send DHCP Ack: %v\n", err)
	}
}

// allocateIP 分配IP地址
func (md *MicroDHCP) allocateIP(mac string) net.IP {
	if ip, ok := md.staticIPs[mac]; ok {
		return ip
	}
	if ip, ok := md.leases[mac]; ok {
		return ip
	}

	for i := 2; i < 255; i++ {
		ip := net.IPv4(192, 168, 100, byte(i))
		if !md.isIPInUse(ip) {
			md.leases[mac] = ip
			return ip
		}
	}
	return nil
}

func (md *MicroDHCP) isIPInUse(ip net.IP) bool {
	for _, leasedIP := range md.leases {
		if leasedIP.Equal(ip) {
			return true
		}
	}
	return false
}

func (md *MicroDHCP) AddStaticIP(mac string, ip net.IP) {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	md.staticIPs[mac] = ip
}

func (md *MicroDHCP) BlacklistMAC(mac string) {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	md.blacklist[mac] = true
}

func (md *MicroDHCP) GetDHCPLeases() map[string]net.IP {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	return md.leases
}
