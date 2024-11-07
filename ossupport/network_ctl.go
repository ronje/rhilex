package ossupport

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

/*
*
* 获取网卡的MAC地址
*
 */
func GetLinuxMacAddr(ifaceName string) (string, error) {
	filePath := filepath.Join("/sys/class/net", ifaceName, "address")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read MAC address file for %s: %w", ifaceName, err)
	}
	macAddr := strings.TrimSpace(string(content))
	// A standard MAC address is 17 characters long (6 groups of 2 hexadecimal digits + 5 colons).
	if len(macAddr) < 17 {
		return "", fmt.Errorf("invalid MAC address length for %s: %s", ifaceName, macAddr)
	}
	return strings.ToUpper(macAddr), nil
}

type NetInterfaceInfo struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
	Addr string `json:"addr"`
}

/*
*
* 获取网卡
*
 */
func GetAvailableInterfaces() ([]NetInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netInterfaces := make([]NetInterfaceInfo, 0, len(interfaces))
	for _, inter := range interfaces {
		info := NetInterfaceInfo{
			Name: inter.Name,
			Mac:  inter.HardwareAddr.String(),
		}
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}
		for i := range addrs {
			addr := addrs[i].String()
			cidr, _, _ := net.ParseCIDR(addr)
			if cidr == nil {
				continue
			}
			if cidr.To4() != nil {
				info.Addr = addr
				break
			}
		}
		netInterfaces = append(netInterfaces, info)
	}
	return netInterfaces, nil
}
