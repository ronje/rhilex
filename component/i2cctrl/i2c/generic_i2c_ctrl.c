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

#include "generic_i2c_ctrl.h"
#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>

// 全局变量，用于存储 I2C 设备文件描述符
static int i2c_fd = -1;
// 全局变量，用于存储从设备地址
static uint8_t i2c_slave_address;

// 初始化 I2C 设备
bool i2c_init(const char *device_path, uint8_t slave_address)
{
    // 打开 I2C 设备文件
    i2c_fd = open(device_path, O_RDWR);
    if (i2c_fd < 0)
    {
        perror("Failed to open I2C device");
        return false;
    }

    // 设置从设备地址
    if (ioctl(i2c_fd, I2C_SLAVE, slave_address) < 0)
    {
        perror("Failed to set I2C slave address");
        close(i2c_fd);
        i2c_fd = -1;
        return false;
    }

    i2c_slave_address = slave_address;
    return true;
}

// 读取单个寄存器的值
bool i2c_read_single_register(uint8_t register_address, uint8_t *value)
{
    if (i2c_fd < 0)
    {
        fprintf(stderr, "I2C device is not initialized\n");
        return false;
    }

    // 先写入要读取的寄存器地址
    if (write(i2c_fd, &register_address, 1) != 1)
    {
        perror("Failed to write register address");
        return false;
    }

    // 读取寄存器的值
    if (read(i2c_fd, value, 1) != 1)
    {
        perror("Failed to read register value");
        return false;
    }

    return true;
}

// 读取多个寄存器的值
bool i2c_read_multiple_registers(uint8_t start_register_address, uint8_t *values, uint8_t num_registers)
{
    if (i2c_fd < 0)
    {
        fprintf(stderr, "I2C device is not initialized\n");
        return false;
    }

    // 先写入要读取的起始寄存器地址
    if (write(i2c_fd, &start_register_address, 1) != 1)
    {
        perror("Failed to write start register address");
        return false;
    }

    // 读取多个寄存器的值
    if (read(i2c_fd, values, num_registers) != num_registers)
    {
        perror("Failed to read multiple register values");
        return false;
    }

    return true;
}

// 写入单个寄存器的值
bool i2c_write_single_register(uint8_t register_address, uint8_t value)
{
    if (i2c_fd < 0)
    {
        fprintf(stderr, "I2C device is not initialized\n");
        return false;
    }

    uint8_t buffer[2] = {register_address, value};
    // 写入寄存器地址和值
    if (write(i2c_fd, buffer, 2) != 2)
    {
        perror("Failed to write register value");
        return false;
    }

    return true;
}

// 写入多个寄存器的值
bool i2c_write_multiple_registers(uint8_t start_register_address, const uint8_t *values, uint8_t num_registers)
{
    if (i2c_fd < 0)
    {
        fprintf(stderr, "I2C device is not initialized\n");
        return false;
    }

    uint8_t *buffer = (uint8_t *)malloc(num_registers + 1);
    if (buffer == NULL)
    {
        perror("Failed to allocate memory");
        return false;
    }

    buffer[0] = start_register_address;
    for (uint8_t i = 0; i < num_registers; i++)
    {
        buffer[i + 1] = values[i];
    }

    // 写入起始寄存器地址和多个寄存器的值
    if (write(i2c_fd, buffer, num_registers + 1) != num_registers + 1)
    {
        perror("Failed to write multiple register values");
        free(buffer);
        return false;
    }

    free(buffer);
    return true;
}

// 关闭 I2C 设备
void i2c_close()
{
    if (i2c_fd >= 0)
    {
        close(i2c_fd);
        i2c_fd = -1;
    }
}