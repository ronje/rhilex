package ossupport

import (
	"fmt"
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
