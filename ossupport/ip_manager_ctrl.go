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
	"os/exec"
	"strings"
	"time"
)

// IPManager provides methods to interact with network interfaces using the 'ip' command
type IPManager struct {
}

// NewIPManager creates a new IPManager instance with a context
func NewIPManager(ctx context.Context) *IPManager {
	return &IPManager{}
}

// Up brings up a network interface
func (ipm *IPManager) Up(interfaceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "link", "set", interfaceName, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error bringing up interface %s: %s, output: %s", interfaceName, err, output)
	}
	return nil
}

// Down brings down a network interface
func (ipm *IPManager) Down(interfaceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "link", "set", interfaceName, "down")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error bringing down interface %s: %s, output: %s", interfaceName, err, output)
	}
	return nil
}

// AddAddress adds an IP address to a network interface
func (ipm *IPManager) AddAddress(interfaceName, ipCIDR string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "addr", "add", ipCIDR, "dev", interfaceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error adding address %s to interface %s: %s, output: %s", ipCIDR, interfaceName, err, output)
	}
	return nil
}

// DeleteAddress removes an IP address from a network interface
func (ipm *IPManager) DeleteAddress(interfaceName, ipCIDR string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "addr", "del", ipCIDR, "dev", interfaceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error deleting address %s from interface %s: %s, output: %s", ipCIDR, interfaceName, err, output)
	}
	return nil
}

// ListInterfaces lists all network interfaces
func (ipm *IPManager) ListInterfaces() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "link", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error listing interfaces: %s, output: %s", err, output)
	}

	// Parse the output to get the list of interfaces
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var interfaces []string
	for _, line := range lines {
		if strings.HasPrefix(line, "2:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				interfaces = append(interfaces, parts[1])
			}
		}
	}
	return interfaces, nil
}
