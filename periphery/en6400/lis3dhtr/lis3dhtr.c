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

#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <linux/i2c-dev.h>
#include <sys/ioctl.h>
#include "lis3dhtr.h"

// 翻译加速度值
double _translateAccValue(uint8_t data)
{
    double resAcc = 0;
    if (data & 0x80)
    { // sign, < 0
        resAcc = ((0x7F - (data & 0x7F) + 1) / 64.0) * (-1);
    }
    else
    { // > 0
        resAcc = (double)data / 64.0;
    }
    return resAcc;
}

// 四舍五入到指定精度
double roundToPrecision(double value, int precisionFactor)
{
    return (double)((int)(value * precisionFactor + 0.5)) / precisionFactor;
}

// 读取加速度数据
void readAcceleration(double *xData, double *yData, double *zData)
{
    int file;
    char *filename = "/dev/i2c-1";
    if ((file = open(filename, O_RDWR)) < 0)
    {
        perror("Failed to open the i2c bus");
        return;
    }

    if (ioctl(file, I2C_SLAVE, LIS3DHTR_ADDR) < 0)
    {
        perror("Failed to acquire bus access and/or talk to slave");
        close(file);
        return;
    }

    // 0. test
    uint8_t rawData;
    if (read(file, &rawData, 1) != 1)
    {
        perror("Failed to read from the i2c bus");
        close(file);
        return;
    }

    // 1. set Data rate and low-power mode
    uint8_t ctrlData = 0x9F; // Low-power mode (5.376 kHz)
    if (write(file, &ctrlData, 1) != 1)
    {
        perror("Failed to write to the i2c bus");
        close(file);
        return;
    }

    // 2. set Full-scale and high-resolution disabled
    ctrlData = 0x00;
    if (write(file, &ctrlData, 1) != 1)
    {
        perror("Failed to write to the i2c bus");
        close(file);
        return;
    }

    // 3. read acceleration data
    uint8_t buffer[5];
    if (read(file, buffer, 5) != 5)
    {
        perror("Failed to read acceleration data");
        close(file);
        return;
    }

    close(file);

    const int precisionFactor = 10000;
    *xData = roundToPrecision(_translateAccValue(buffer[0]), precisionFactor);
    *yData = roundToPrecision(_translateAccValue(buffer[2]), precisionFactor);
    *zData = roundToPrecision(_translateAccValue(buffer[4]), precisionFactor);
}