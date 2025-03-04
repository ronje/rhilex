// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ossupport

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// UDHCPClient 封装 udhcpc 功能的结构体
type UDHCPClient struct {
	Interface string
}

// NewUDHCPClient 创建一个新的 UDHCPClient 实例
func NewUDHCPClient(iface string) *UDHCPClient {
	return &UDHCPClient{
		Interface: iface,
	}
}

// RunCommand 执行 udhcpc 命令的辅助函数
func (u *UDHCPClient) RunCommand(args ...string) (string, error) {
	cmdArgs := append([]string{"-i", u.Interface}, args...)
	cmd := exec.Command("udhcpc", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run udhcpc command: %v, output: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// RequestIP 向 DHCP 服务器请求 IP 地址
func (u *UDHCPClient) RequestIP() (string, error) {
	return u.RunCommand()
}

// RequestIPInBackground 以守护进程模式请求 IP 地址
func (u *UDHCPClient) RequestIPInBackground() (string, error) {
	return u.RunCommand("-b")
}

// RequestIPWithTimeout 以指定超时时间请求 IP 地址
func (u *UDHCPClient) RequestIPWithTimeout(timeout int) (string, error) {
	return u.RunCommand("-t", fmt.Sprintf("%d", timeout))
}

// RequestIPQuietly 以安静模式请求 IP 地址
func (u *UDHCPClient) RequestIPQuietly() (string, error) {
	return u.RunCommand("-q")
}

// RequestIPWithScript 以指定脚本请求 IP 地址
func (u *UDHCPClient) RequestIPWithScript(scriptPath string) (string, error) {
	return u.RunCommand("-s", scriptPath)
}

// AllocateIPAddressUsingUdhcpc 使用 udhcpc 命令为指定的网络接口分配 IP 地址。
// interfaceName 是要分配 IP 地址的网络接口的名称。
// 超时时间设置为5秒，如果命令执行超时，返回错误。
func AllocateIPAddressUsingUdhcpc(interfaceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "udhcpc", "-i", interfaceName)
	log.Println("AllocateIPAddressUsingUdhcpc = ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("udhcpc command failed for interface %s: %w, output: %s", interfaceName, err, string(output))
	}
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("udhcpc command for interface %s timed out", interfaceName)
	}
	fmt.Printf("udhcpc output for interface %s: %s\n", interfaceName, string(output))
	return nil
}
