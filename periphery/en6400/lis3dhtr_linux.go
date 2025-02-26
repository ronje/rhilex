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

package en6400

/*
#cgo CFLAGS: -I./lis3dhtr -v
#cgo LDFLAGS: -lm
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "lis3dhtr.h"
#include "lis3dhtr.c"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// AccelerationData 存储加速度数据
type AccelerationData struct {
	X float64
	Y float64
	Z float64
}
// To String
func (a AccelerationData) String() string {
	return fmt.Sprintf("AccelerationData = X: %.2f, Y: %.2f, Z: %.2f", a.X, a.Y, a.Z)
}
// ReadAcceleration 读取加速度数据
func ReadAcceleration() (AccelerationData, error) {
	devicePath := C.CString("/dev/i2c-1")
	defer C.free(unsafe.Pointer(devicePath))

	file := C.open_i2c_device(devicePath)
	if int(file) < 0 {
		return AccelerationData{}, fmt.Errorf("Can not open I2C device")
	}
	defer C.close(file)
	// 设置 I2C 地址
	if int(C.set_i2c_address(file, C.int(C.LIS3DHTR_ADDR))) < 0 {
		return AccelerationData{}, fmt.Errorf("Can not set I2C address")
	}
	// 初始化 LIS3DHT 传感器
	if int(C.init_lis3dht(file)) < 0 {
		return AccelerationData{}, fmt.Errorf("Can not init LIS3DHT")
	}
	var xAcc, yAcc, zAcc C.float
	if int(C.read_acceleration_data(file, &xAcc, &yAcc, &zAcc)) != 0 {
		return AccelerationData{}, fmt.Errorf("Can not read acceleration data")
	}
	return AccelerationData{
		X: float64(xAcc),
		Y: float64(yAcc),
		Z: float64(zAcc),
	}, nil
}
