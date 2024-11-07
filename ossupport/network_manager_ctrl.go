// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ossupport

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// NetworkManager provides methods to interact with NetworkManager
type NetworkManager struct {
}

// NewNetworkManager creates a new NetworkManager instance with a context
func NewNetworkManager() *NetworkManager {
	return &NetworkManager{}
}

// Connect connects to a network connection profile by name
func (nm *NetworkManager) Connect(profileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nmcli", "connection", "up", profileName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error connecting to profile %s: %s, output: %s", profileName, err, output)
	}
	return nil
}

// Disconnect disconnects from a network connection profile by name
func (nm *NetworkManager) Disconnect(profileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nmcli", "connection", "down", profileName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error disconnecting from profile %s: %s, output: %s", profileName, err, output)
	}
	return nil
}

// CheckConnectionStatus checks the status of a network connection profile
func (nm *NetworkManager) CheckConnectionStatus(profileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nmcli", "connection", "show", profileName, "--active")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error checking connection status for profile %s: %s, output: %s", profileName, err, output)
	}

	// Parse the output to determine the status
	outputStr := string(output)
	if strings.Contains(outputStr, "connected") {
		return "connected", nil
	} else if strings.Contains(outputStr, "disconnected") {
		return "disconnected", nil
	}
	return "unknown", nil
}

// ListConnections lists all network connections
func (nm *NetworkManager) ListConnections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nmcli", "connection", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error listing connections: %s, output: %s", err, output)
	}

	// Parse the output to get the list of connections
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var connections []string
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME") {
			continue // Skip header line
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			connections = append(connections, parts[0])
		}
	}
	return connections, nil
}
