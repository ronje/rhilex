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
	"fmt"
	"log"
	"os"
	"os/exec"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506" {
		_HAAS506_DI_Init()
		_HAAS506_DO_Init()
		_HAAS506_LED_Init()
	}
}

func _HAAS506_GPIOInit(Pin string, direction string) {
	//gpio export
	cmd := fmt.Sprintf("echo %s > /sys/class/gpio/export", Pin)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Println("[HAAS506_GPIOInit] error", err, string(output))
		return
	}
	//gpio set direction
	cmd = fmt.Sprintf("echo %s > /sys/class/gpio/gpio%s/direction", direction, Pin)
	output, err = exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Println("[HAAS506_GPIOInit] error", err, string(output))
	}
}

func HAAS506_GPIOGetByFile(pin byte) (int, error) {
	return _HAAS506_GPIO_Get(fmt.Sprintf(HAAS506_DI_SYSDEV_PATH, pin))
}

func _HAAS506_GPIO_Get(gpioPath string) (int, error) {
	bites, err := os.ReadFile(gpioPath)
	if err != nil {
		return 0, err
	}
	if len(bites) > 0 {
		if bites[0] == '0' || bites[0] == 48 {
			return 0, nil
		}
		if bites[1] == '1' || bites[0] == 49 {
			return 1, nil
		}
	}
	return 0, fmt.Errorf("read gpio value failed: %s, value: %v", gpioPath, bites)
}

func HAAS506_GPIOSetByFile(pin, value int) error {
	return _HAAS506_GPIO_Set(fmt.Sprintf(HAAS506_DI_SYSDEV_PATH, pin), value)
}

func _HAAS506_GPIO_Set(gpioPath string, value int) error {
	if value == 1 {
		err := os.WriteFile(gpioPath, []byte{'1'}, 0644)
		if err != nil {
			return err
		}
	}
	if value == 0 {
		err := os.WriteFile(gpioPath, []byte{'0'}, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
