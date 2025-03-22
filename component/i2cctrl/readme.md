<!--
 Copyright (C) 2025 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# I2C 操作库开发文档

## 一、概述
本项目实现了一套用于 I2C（Inter-Integrated Circuit）设备通信的操作库，包含 C 语言实现的底层操作和使用 Cgo 封装的 Go 语言接口。通过 `I2CController` 结构体对 I2C 设备进行管理，提供了初始化、读写单个或多个寄存器以及关闭设备等功能。

## 二、功能特性
1. **设备初始化**：支持指定 I2C 设备路径和从设备地址进行初始化。
2. **寄存器读写**：提供读写单个寄存器和多个寄存器的功能。
3. **资源管理**：包含关闭 I2C 设备的方法，确保资源正确释放。
4. **跨语言封装**：使用 Cgo 技术将 C 语言实现的 I2C 操作封装为 Go 语言接口，方便在 Go 项目中使用。

## 三、文件结构
### 1. C 代码部分
- `generic_i2c_ctrl.h`：定义了 I2C 操作的函数原型，供 C 代码调用。
- `generic_i2c_ctrl.c`：实现了 `generic_i2c_ctrl.h` 中声明的函数，通过 Linux 的 I2C 设备文件和 `ioctl` 系统调用进行实际的 I2C 通信。

### 2. Go 代码部分
- `i2c.go`：使用 Cgo 调用 C 代码，将 C 函数封装到 `I2CController` 结构体的方法中。
- `main.go`：提供了使用 `I2CController` 的示例代码。

## 四、C 代码接口说明

### 1. `generic_i2c_ctrl.h`
```c
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
```

### 2. 函数详细说明
#### `i2c_init`
- **功能**：初始化 I2C 设备，打开指定的 I2C 设备文件并设置从设备地址。
- **参数**：
  - `device_path`：I2C 设备文件的路径，如 `/dev/i2c-1`。
  - `slave_address`：I2C 从设备的地址。
- **返回值**：初始化成功返回 `true`，失败返回 `false`。

#### `i2c_read_single_register`
- **功能**：读取指定寄存器的值。
- **参数**：
  - `register_address`：要读取的寄存器地址。
  - `value`：用于存储读取到的值的指针。
- **返回值**：读取成功返回 `true`，失败返回 `false`。

#### `i2c_read_multiple_registers`
- **功能**：从指定起始地址开始，连续读取多个寄存器的值。
- **参数**：
  - `start_register_address`：起始寄存器地址。
  - `values`：用于存储读取到的值的数组指针。
  - `num_registers`：要读取的寄存器数量。
- **返回值**：读取成功返回 `true`，失败返回 `false`。

#### `i2c_write_single_register`
- **功能**：向指定寄存器写入一个值。
- **参数**：
  - `register_address`：要写入的寄存器地址。
  - `value`：要写入的值。
- **返回值**：写入成功返回 `true`，失败返回 `false`。

#### `i2c_write_multiple_registers`
- **功能**：从指定起始地址开始，连续写入多个寄存器的值。
- **参数**：
  - `start_register_address`：起始寄存器地址。
  - `values`：包含要写入的值的数组指针。
  - `num_registers`：要写入的寄存器数量。
- **返回值**：写入成功返回 `true`，失败返回 `false`。

#### `i2c_close`
- **功能**：关闭 I2C 设备文件，释放资源。
- **参数**：无
- **返回值**：无

## 五、Go 代码接口说明

### 1. `I2CController` 结构体
```go
type I2CController struct {
    devicePath    string
    slaveAddress  uint8
    isInitialized bool
}
```
- `devicePath`：I2C 设备文件的路径。
- `slaveAddress`：I2C 从设备的地址。
- `isInitialized`：标记设备是否已经初始化。

### 2. 方法详细说明
#### `NewI2CController`
- **功能**：创建一个新的 `I2CController` 实例。
- **参数**：
  - `devicePath`：I2C 设备文件的路径。
  - `slaveAddress`：I2C 从设备的地址。
- **返回值**：返回一个指向 `I2CController` 结构体的指针。

#### `Init`
- **功能**：初始化 I2C 设备，如果设备已经初始化则直接返回 `true`。
- **参数**：无
- **返回值**：初始化成功返回 `true`，失败返回 `false`。

#### `ReadSingleRegister`
- **功能**：读取指定寄存器的值。
- **参数**：
  - `registerAddress`：要读取的寄存器地址。
- **返回值**：返回读取到的值和读取是否成功的布尔值。

#### `ReadMultipleRegisters`
- **功能**：从指定起始地址开始，连续读取多个寄存器的值。
- **参数**：
  - `startRegisterAddress`：起始寄存器地址。
  - `numRegisters`：要读取的寄存器数量。
- **返回值**：返回包含读取到的值的切片和读取是否成功的布尔值。

#### `WriteSingleRegister`
- **功能**：向指定寄存器写入一个值。
- **参数**：
  - `registerAddress`：要写入的寄存器地址。
  - `value`：要写入的值。
- **返回值**：写入成功返回 `true`，失败返回 `false`。

#### `WriteMultipleRegisters`
- **功能**：从指定起始地址开始，连续写入多个寄存器的值。
- **参数**：
  - `startRegisterAddress`：起始寄存器地址。
  - `values`：包含要写入的值的切片。
- **返回值**：写入成功返回 `true`，失败返回 `false`。

#### `Close`
- **功能**：关闭 I2C 设备，释放资源，并将 `isInitialized` 标记设置为 `false`。
- **参数**：无
- **返回值**：无

## 六、使用示例
```go
package main

import (
    "fmt"
    "./i2c"
)

func main() {
    devicePath := "/dev/i2c-1"
    slaveAddress := uint8(0x50)

    // 创建 I2CController 实例
    controller := i2c.NewI2CController(devicePath, slaveAddress)

    // 初始化 I2C 设备
    if!controller.Init() {
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
```

## 七、编译和运行
确保 `generic_i2c_ctrl.c`、`generic_i2c_ctrl.h`、`i2c.go` 和 `main.go` 文件在同一目录下，然后使用以下命令编译和运行：
```sh
go build -o i2c_example main.go
sudo ./i2c_example  # 可能需要 sudo 权限来访问 I2C 设备文件
```

## 八、注意事项
1. **权限问题**：在 Linux 系统中，访问 I2C 设备文件通常需要 root 权限，因此运行程序时可能需要使用 `sudo`。
2. **内存管理**：C 代码中使用 `malloc` 分配的内存会在使用完后自动释放，Go 代码中使用 `C.CString` 分配的内存会在使用完后通过 `C.free` 释放，确保不会出现内存泄漏。
3. **错误处理**：在调用各个函数时，应检查返回值以确保操作成功。如果操作失败，可能需要根据具体的错误信息进行调试和处理。 