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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package archsupport

import (
	"fmt"
	"log"
	"os"
)

// HaaS506-LD1是基于匠芯创D213（平头哥玄铁C906 RISC-V内核）为主控的工业级Linux可编程网关RTU
// https://www.yuque.com/haas506/wiki/agyla4pphb7fxpd1
const (
	__HAAS506_LED2_SYSDEV_PATH = "/sys/devices/platform/leds/leds/led2/brightness"
	__HAAS506_LED3_SYSDEV_PATH = "/sys/devices/platform/leds/leds/led3/brightness"
	__HAAS506_LED4_SYSDEV_PATH = "/sys/devices/platform/leds/leds/led4/brightness"
	__HAAS506_LED5_SYSDEV_PATH = "/sys/devices/platform/leds/leds/led5/brightness"
)

// echo "0" > /sys/devices/platform/leds/leds/led2/brightness
// echo "0" > /sys/devices/platform/leds/leds/led3/brightness
// echo "0" > /sys/devices/platform/leds/leds/led4/brightness
// echo "0" > /sys/devices/platform/leds/leds/led5/brightness

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506" {
		_HAAS506_LEDAllInit()
	}
}
func _HAAS506_LEDAllInit() {
	_HAAS506_LedInit(__HAAS506_LED2_SYSDEV_PATH, "out")
	_HAAS506_LedInit(__HAAS506_LED3_SYSDEV_PATH, "out")
	_HAAS506_LedInit(__HAAS506_LED4_SYSDEV_PATH, "out")
	_HAAS506_LedInit(__HAAS506_LED5_SYSDEV_PATH, "out")
}

func _HAAS506_LedInit(Pin string, direction string) {
	log.Printf("[HAAS506_LEDInit] LED(%s, %s) Init...", Pin, direction)
	log.Printf("[HAAS506_LEDInit] LED(%s, %s) Init Ok.", Pin, direction)
}

func HAAS506_LEDSet(pin, value int) error {
	return _HAAS506_LEDSet(pin, value)
}

func HAAS506_LEDGet(pin int) (int, error) {
	if pin == 2 {
		return _HAAS506_LED_get(__HAAS506_LED2_SYSDEV_PATH)
	}
	if pin == 3 {
		return _HAAS506_LED_get(__HAAS506_LED3_SYSDEV_PATH)
	}
	if pin == 4 {
		return _HAAS506_LED_get(__HAAS506_LED4_SYSDEV_PATH)
	}
	if pin == 5 {
		return _HAAS506_LED_get(__HAAS506_LED5_SYSDEV_PATH)
	}
	return 0, fmt.Errorf("read LED value failed: %d", pin)
}

func _HAAS506_LEDSet(pin, value int) error {
	if pin == 2 {
		return _HAAS506_LED_set(__HAAS506_LED2_SYSDEV_PATH, value)
	}
	if pin == 3 {
		return _HAAS506_LED_set(__HAAS506_LED3_SYSDEV_PATH, value)
	}
	if pin == 4 {
		return _HAAS506_LED_set(__HAAS506_LED4_SYSDEV_PATH, value)
	}
	if pin == 5 {
		return _HAAS506_LED_set(__HAAS506_LED5_SYSDEV_PATH, value)
	}
	return fmt.Errorf("invalid led number:%d, value: %d", pin, value)
}
func _HAAS506_LED_get(LEDPath string) (int, error) {
	bites, err := os.ReadFile(LEDPath)
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
	return 0, fmt.Errorf("read LED value failed: %s, value: %v", LEDPath, bites)
}
func _HAAS506_LED_set(LEDPath string, value int) error {
	if value == 1 {
		err := os.WriteFile(LEDPath, []byte{'1'}, 0644)
		if err != nil {
			return err
		}
	}
	if value == 0 {
		err := os.WriteFile(LEDPath, []byte{'0'}, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
