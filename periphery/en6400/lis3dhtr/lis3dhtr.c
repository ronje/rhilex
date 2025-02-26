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

#include "lis3dhtr.h"
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>
#include <math.h>

// 打开 I2C 设备
int open_i2c_device(const char *device)
{
    int file = open(device, O_RDWR);
    if (file < 0)
    {
        perror("Unable to open I2C device");
        return -1;
    }
    return file;
}

// 设置 I2C 地址
int set_i2c_address(int file, int addr)
{
    if (ioctl(file, I2C_SLAVE, addr) < 0)
    {
        perror("Unable to set I2C address");
        return -1;
    }
    return 0;
}

// 写寄存器值
int write_i2c_byte(int file, uint8_t reg, uint8_t value)
{
    uint8_t data[2] = {reg, value};
    if (write(file, data, 2) != 2)
    {
        perror("Failed to write to the I2C bus");
        return -1;
    }
    return 0;
}

// 读取连续 6 个字节
int read_i2c_6bytes(int file, uint8_t reg, uint8_t values[6])
{
    uint8_t reg1 = reg | 0x80; // 启用寄存器自动递增模式
    // 设置寄存器地址
    if (write(file, &reg1, 1) != 1)
    {
        perror("Failed to set register address");
        return -1;
    }

    // 批量读取 6 个字节
    if (read(file, values, 6) != 6)
    {
        perror("Failed to read from the I2C bus");
        return -1;
    }
    return 0;
}
// 读取寄存器值
uint8_t read_i2c_byte(int file, uint8_t reg)
{
    uint8_t value;
    if (write(file, &reg, 1) != 1)
    {
        perror("Failed to set register address");
        return -1;
    }
    if (read(file, &value, 1) != 1)
    {
        perror("Failed to read from the I2C bus");
        return -1;
    }
    return value;
}

// 初始化 LIS3DHT 传感器
int init_lis3dht(int file)
{
    uint8_t rawData;

    // 设备检测（读取 WHO_AM_I 寄存器）
    rawData = read_i2c_byte(file, DEVICE_ID_REG);
    if (rawData == 0x33)
    {
        printf("LIS3DHTR detected.\n");
    }
    else
    {
        printf("Unknown chip detected: 0x%02X\n", rawData);
        return -1;
    }

    // 配置控制寄存器 1 (设置数据速率和低功耗模式)
    // 设置为低功耗模式 5.376 kHz
    if (write_i2c_byte(file, CTRL_REG1, 0x9F) < 0)
    {
        return -1;
    }
    // 再读一下CTRL_REG1
    uint8_t v1 = read_i2c_byte(file, CTRL_REG1);
    if (v1 != 0x9F)
    {
        return -1;
    }
    // 配置控制寄存器 4 (设置满量程)
    if (write_i2c_byte(file, CTRL_REG4, 0x00) < 0)
    { // 设置为 ±2g
        return -1;
    }
    uint8_t v2 = read_i2c_byte(file, CTRL_REG4);
    if (v2 != 0)
    {
        return -1;
    }
    return 0;
}

// 读取加速度数据
int read_acceleration_data(int file, float *xAcc, float *yAcc, float *zAcc)
{
    uint8_t buffer[6];
    float xData, yData, zData;
    int ret = read_i2c_6bytes(file, OUT_X_MSB_REG, buffer);
    if (ret != 0)
    {
        printf("read acceleration error\n");
        return ret;
    }
    xData = translate_acc_value(buffer[0]);
    yData = translate_acc_value(buffer[2]);
    zData = translate_acc_value(buffer[4]);
    int precisionFactor = 10000; // 控制精度的因子
    float xDataRounded = round_to_precision(xData, precisionFactor);
    float yDataRounded = round_to_precision(yData, precisionFactor);
    float zDataRounded = round_to_precision(zData, precisionFactor);
    *xAcc = xDataRounded;
    *yAcc = yDataRounded;
    *zAcc = zDataRounded;
    return 0;
}
// 转换加速度值的函数
float translate_acc_value(uint8_t data)
{
    float resAcc = 0;
    // 检查第 7 位（符号位）是否为 1
    if (data & 0x80)
    {
        // 负数情况
        // 取反加 1 得到绝对值
        uint8_t absValue = (~data & 0x7F) + 1;
        resAcc = -(float)absValue / 64.0;
    }
    else
    {
        // 正数情况
        resAcc = (float)data / 64.0;
    }
    return resAcc;
}

// 精度控制函数，四舍五入
float round_to_precision(float value, int precisionFactor)
{
    return round(value * precisionFactor) / precisionFactor;
}