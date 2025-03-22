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

#ifndef GENERIC_I2C_CTRL_H
#define GENERIC_I2C_CTRL_H

#include <stdint.h>
#include <stdbool.h>

// 初始化 I2C 设备
bool i2c_init(const char *device_path, uint8_t slave_address);

// 读取单个寄存器的值
bool i2c_read_single_register(uint8_t register_address, uint8_t *value);

// 读取多个寄存器的值
bool i2c_read_multiple_registers(uint8_t start_register_address, uint8_t *values, uint8_t num_registers);

// 写入单个寄存器的值
bool i2c_write_single_register(uint8_t register_address, uint8_t value);

// 写入多个寄存器的值
bool i2c_write_multiple_registers(uint8_t start_register_address, const uint8_t *values, uint8_t num_registers);

// 关闭 I2C 设备
void i2c_close();

#endif