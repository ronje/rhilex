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

package haas506

import (
	"bytes"
	"os/exec"
	"strings"
)

// isInterfaceUp checks if a network interface is up.
// It returns true if the interface is up, false otherwise, and an error if the check fails.
func isInterfaceUp(interfaceName string) (bool, error) {
	// Execute the 'ip link show' command for the specified interface
	cmd := exec.Command("ip", "link", "show", interfaceName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	// Parse the output to determine if the interface is up
	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "state UP") {
			return true, nil
		}
	}

	return false, nil
}
