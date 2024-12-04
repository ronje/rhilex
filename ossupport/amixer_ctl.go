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
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Amixer provides methods to interact with ALSA mixer
type Amixer struct {
	cardName    string
	controlName string
}

// NewAmixer creates a new Amixer instance
func NewAmixer(cardName, controlName string) *Amixer {
	return &Amixer{
		cardName:    cardName,
		controlName: controlName,
	}
}

// GetVolume retrieves the current volume level
func (a *Amixer) GetVolume() (int, error) {
	cmd := exec.Command("amixer", "-c", a.cardName, "get", a.controlName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "[%") {
			re := regexp.MustCompile(`\[([0-9]+)%\]`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				volume, err := strconv.Atoi(matches[1])
				if err != nil {
					return 0, err
				}
				return volume, nil
			}
		}
	}

	return 0, fmt.Errorf("volume not found")
}

// SetVolume sets the volume to the specified level
func (a *Amixer) SetVolume(volume int) error {
	cmd := exec.Command("amixer", "-c", a.cardName, "set", a.controlName, strconv.Itoa(volume)+"%")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error setting volume: %s, output: %s", err, output)
	}
	return nil
}

// Mute toggles the mute state
func (a *Amixer) Mute() error {
	cmd := exec.Command("amixer", "-c", a.cardName, "set", a.controlName, "mute")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error muting: %s, output: %s", err, output)
	}
	return nil
}

// Unmute toggles the unmute state
func (a *Amixer) Unmute() error {
	cmd := exec.Command("amixer", "-c", a.cardName, "set", a.controlName, "unmute")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error unmuting: %s, output: %s", err, output)
	}
	return nil
}
