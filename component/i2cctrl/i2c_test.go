// Copyright (C) 2025 wwhai
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
package i2cctrl

import (
	"fmt"
	"testing"
)

func Test_I2C_RW(t *testing.T) {
	devicePath := "/dev/i2c-1"
	slaveAddress := uint8(0x50)

	// 创建 I2CController 实例
	controller := NewI2CController(devicePath, slaveAddress)

	// 初始化 I2C 设备
	if !controller.Init() {
		fmt.Println("Failed to initialize I2C device")
		return
	}

	registerAddress := uint8(0x01)
	valueToWrite := uint8(0xAA)

	// 写入单个寄存器的值
	if controller.WriteSingleRegister(registerAddress, valueToWrite) {
		fmt.Printf("Successfully wrote 0x%02X to register 0x%02X\n", valueToWrite, registerAddress)
	}

	// 读取单个寄存器的值
	value, success := controller.ReadSingleRegister(registerAddress)
	if success {
		fmt.Printf("Read value 0x%02X from register 0x%02X\n", value, registerAddress)
	}

	// 关闭 I2C 设备
	controller.Close()
}
