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

#ifndef LIS3DHTR_H
#define LIS3DHTR_H

#include <stdint.h>

#define LIS3DHTR_ADDR 24
#define DEVICE_ID_REG 0x0F
#define CTRL_REG0 0x1E
#define OUT_X_MSB_REG 0x29
#define OUT_X_LSB_REG 0x28
#define OUT_Y_MSB_REG 0x2B
#define OUT_Y_LSB_REG 0x2A
#define OUT_Z_MSB_REG 0x2D
#define OUT_Z_LSB_REG 0x2C
#define CTRL_REG1_REG 0x20
#define CTRL_REG4_REG 0x23

// 翻译加速度值
double _translateAccValue(uint8_t data);

// 四舍五入到指定精度
double roundToPrecision(double value, int precisionFactor);

// 读取加速度数据
void readAcceleration(double *xData, double *yData, double *zData);

#endif