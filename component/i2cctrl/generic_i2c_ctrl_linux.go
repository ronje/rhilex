//go:build linux
// +build linux

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

/*
#cgo CFLAGS: -I./i2c
#include "generic_i2c_ctrl.h"
#include "generic_i2c_ctrl.c"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// I2CController 结构体用于封装 I2C 操作
type I2CController struct {
	devicePath    string
	slaveAddress  uint8
	isInitialized bool
}

// NewI2CController 创建一个新的 I2CController 实例
func NewI2CController(devicePath string, slaveAddress uint8) *I2CController {
	return &I2CController{
		devicePath:    devicePath,
		slaveAddress:  slaveAddress,
		isInitialized: false,
	}
}

// Init 初始化 I2C 设备
func (i *I2CController) Init() bool {
	if i.isInitialized {
		return true
	}
	cDevicePath := C.CString(i.devicePath)
	defer C.free(unsafe.Pointer(cDevicePath))
	result := bool(C.i2c_init(cDevicePath, C.uint8_t(i.slaveAddress)))
	if result {
		i.isInitialized = true
	}
	return result
}

// ReadSingleRegister 读取单个寄存器的值
func (i *I2CController) ReadSingleRegister(registerAddress uint8) (uint8, bool) {
	if !i.isInitialized {
		return 0, false
	}
	var value C.uint8_t
	success := bool(C.i2c_read_single_register(C.uint8_t(registerAddress), &value))
	return uint8(value), success
}

// ReadMultipleRegisters 读取多个寄存器的值
func (i *I2CController) ReadMultipleRegisters(startRegisterAddress uint8, numRegisters uint8) ([]uint8, bool) {
	if !i.isInitialized {
		return nil, false
	}
	values := make([]uint8, numRegisters)
	success := bool(C.i2c_read_multiple_registers(C.uint8_t(startRegisterAddress), (*C.uint8_t)(&values[0]), C.uint8_t(numRegisters)))
	return values, success
}

// WriteSingleRegister 写入单个寄存器的值
func (i *I2CController) WriteSingleRegister(registerAddress uint8, value uint8) bool {
	if !i.isInitialized {
		return false
	}
	return bool(C.i2c_write_single_register(C.uint8_t(registerAddress), C.uint8_t(value)))
}

// WriteMultipleRegisters 写入多个寄存器的值
func (i *I2CController) WriteMultipleRegisters(startRegisterAddress uint8, values []uint8) bool {
	if !i.isInitialized {
		return false
	}
	numRegisters := len(values)
	if numRegisters == 0 {
		return false
	}
	return bool(C.i2c_write_multiple_registers(C.uint8_t(startRegisterAddress), (*C.uint8_t)(&values[0]), C.uint8_t(numRegisters)))
}

// Close 关闭 I2C 设备
func (i *I2CController) Close() {
	if i.isInitialized {
		C.i2c_close()
		i.isInitialized = false
	}
}
