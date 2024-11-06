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
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Ifconfig 执行系统的 ifconfig 或 ipconfig 命令，并返回输出结果
func Ifconfig() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, "ipconfig", "/all")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, "ifconfig", "-a")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		// 将 GBK 转换为 UTF-8
		reader := transform.NewReader(
			bytes.NewReader(out.Bytes()), // 使用 bytes.NewReader
			simplifiedchinese.GBK.NewDecoder(),
		)
		decodedOutput, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return string(decodedOutput), nil

	}
	return out.String(), nil
}

// GetAllIps returns a slice of all non-loopback IPv4 addresses on the system.
func GetAllIps() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, fmt.Sprintf("[%10s] http://%s:2580", iface.Name, ip.String()))
		}
	}

	return ips, nil
}

// getEthList returns a list of Ethernet (wired) interfaces.
func GetEthList() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ethIfaces []net.Interface
	for _, iface := range ifaces {
		if iface.Name == "lo" {
			continue
		}
		ethIfaces = append(ethIfaces, iface)
	}

	return ethIfaces, nil
}

/**
 * print all Macs
 *
 */
func ShowMacAddress() ([]string, error) {
	var macs []string
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifa := range ifas {
		if ifa.Flags == net.FlagLoopback {
			continue
		}
		if slices.ContainsFunc(
			[]string{"loopback", "veth", "virbr", "tun", "tap", "docker"},
			func(s string) bool {
				return strings.Contains(strings.ToLower(ifa.Name), s)
			}) {
			continue
		}
		mac := ifa.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, fmt.Sprintf("[%10s]: %s", ifa.Name, strings.ToUpper(mac)))
		}
	}
	return macs, nil
}

// getCPUID 返回CPU的序列号。
func GetCPUID() (string, error) {
	var cpuID string
	var err error

	switch runtime.GOOS {
	case "windows":
		cpuID, err = getCPUIDWindows()
	case "linux":
		cpuID, err = getCPUIDLinux()
	default:
		err = fmt.Errorf("Unsupported System: %s", runtime.GOOS)
	}
	return strings.ToUpper(cpuID), err
}

// getCPUIDWindows 获取Windows系统的CPU序列号。
func getCPUIDWindows() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wmic", "cpu", "get", "ProcessorId")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(string(output), "\n")[1]), nil
}

// getCPUIDLinux 获取Linux系统的CPU序列号。
func getCPUIDLinux() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cat", "/proc/cpuinfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cpuinfo := string(output)
	start := strings.Index(cpuinfo, "Serial")
	if start == -1 {
		return "", fmt.Errorf("Can't Find CPU Info")
	}
	end := strings.Index(cpuinfo[start:], "\n")
	if end == -1 {
		return "", fmt.Errorf("Invalid CPU Info")
	}
	serialLine := strings.TrimSpace(cpuinfo[start : start+end])
	serial := strings.Split(serialLine, ":")[1]
	return strings.TrimSpace(serial), nil
}

/**
 * 启动网卡
 *
 */
func IfconfigSetIface(interfaceName, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ifconfig", interfaceName, status)
	log.Println("Ifconfig Set Iface = ", cmd.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start network interface %s: %w", interfaceName, err)
	}
	return nil
}
