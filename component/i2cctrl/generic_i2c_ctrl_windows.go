//go:build windows
// +build windows

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

// I2CController 结构体用于封装 I2C 操作
type I2CController struct {
}

// NewI2CController 创建一个新的 I2CController 实例
func NewI2CController(devicePath string, slaveAddress uint8) *I2CController {
	return &I2CController{}
}

// Init 初始化 I2C 设备
func (i *I2CController) Init() bool {
	return true
}

// ReadSingleRegister 读取单个寄存器的值
func (i *I2CController) ReadSingleRegister(registerAddress uint8) (uint8, bool) {
	return 0, true
}

// ReadMultipleRegisters 读取多个寄存器的值
func (i *I2CController) ReadMultipleRegisters(startRegisterAddress uint8, numRegisters uint8) ([]uint8, bool) {
	return []uint8{}, true

}

// WriteSingleRegister 写入单个寄存器的值
func (i *I2CController) WriteSingleRegister(registerAddress uint8, value uint8) bool {
	return true
}

// WriteMultipleRegisters 写入多个寄存器的值
func (i *I2CController) WriteMultipleRegisters(startRegisterAddress uint8, values []uint8) bool {
	return true
}

// Close 关闭 I2C 设备
func (i *I2CController) Close() {
}
