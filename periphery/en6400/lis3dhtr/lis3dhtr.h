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

#ifndef LIS3DHT_H
#define LIS3DHT_H

#include <stdint.h>

// I2C 地址和寄存器定义
#define LIS3DHTR_ADDR 0x18 // LIS3DHTR I2C 地址
#define DEVICE_ID_REG 0x0F // WHO_AM_I 寄存器
#define CTRL_REG1 0x20     // 控制寄存器 1
#define CTRL_REG4 0x23     // 控制寄存器 4
#define OUT_X_MSB_REG 0x29 // X 轴 MSB 寄存器
#define OUT_X_LSB_REG 0x28 // X 轴 LSB 寄存器
#define OUT_Y_MSB_REG 0x2B // Y 轴 MSB 寄存器
#define OUT_Y_LSB_REG 0x2A // Y 轴 LSB 寄存器
#define OUT_Z_MSB_REG 0x2D // Z 轴 MSB 寄存器
#define OUT_Z_LSB_REG 0x2C // Z 轴 LSB 寄存器

// 函数声明
int open_i2c_device(const char *device);
int set_i2c_address(int file, int addr);
int write_i2c_byte(int file, uint8_t reg, uint8_t value);
int read_i2c_6bytes(int file, uint8_t reg, uint8_t values[6]);
uint8_t read_i2c_byte(int file, uint8_t reg);
int init_lis3dht(int file);
int read_acceleration_data(int file, float *xAcc, float *yAcc, float *zAcc);
float translate_acc_value(uint8_t data);
float round_to_precision(float value, int precisionFactor);
#endif // LIS3DHT_H
