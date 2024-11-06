package ossupport

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// GetWindowsMACAddress 获取网卡MAC
func GetWindowsFirstMacAddress() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe", "-Command",
		`wmic nicconfig where "Index=1" get MACAddress`)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	macAddress := strings.TrimSpace(string(output))
	if strings.Contains(macAddress, "MACAddress") {
		macAddress = strings.Split(macAddress, "\n")[1]
	}
	return strings.ToUpper(macAddress), nil
}
