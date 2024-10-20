# 自定义交互协议指南

## 1. 概述

本指南提供了一个自定义交互协议的结构和实现，适用于各种数据传输场景。该协议由 `Header` 和 `AppLayerFrame` 两部分组成，并包含 CRC 校验功能以确保数据传输的完整性。

## 2. 数据结构

### 2.1 Header

`Header` 结构体用于描述数据的元信息，具体字段如下：

- **Type**: `[2]byte`
  表示数据的类型，使用 2 字节表示。

- **Length**: `[2]byte`
  表示有效载荷的长度，使用 2 字节表示。

### 2.2 AppLayerFrame

`AppLayerFrame` 结构体用于封装完整的数据帧，具体字段如下：

- **Header**: `Header`
  包含数据的头部信息。

- **Payload**: `[]byte`
  表示实际的数据内容，长度可变。

## 3. CRC 校验码

在数据传输中，CRC（循环冗余校验）是一种常用的错误检测机制。本协议使用 CRC-8 校验码以确保数据的完整性。CRC 的计算方法在实现中提供，用户可以直接使用该方法进行数据校验。

## 4. 使用场景

本协议适用于各种需要数据传输的场景，例如：

- 网络通信
- 文件传输
- 设备间数据交互

通过使用 `Header`，接收方可以快速了解接收到的数据类型及其长度，从而进行相应处理。

## 5. C语言实现示例
```c
#include <stdio.h>
#include <stdint.h>
#include <string.h>

#define HEADER_SIZE 4 // 2 bytes for Type + 2 bytes for Length
#define PAYLOAD_MAX_SIZE 256 // Maximum payload size

// Header structure
typedef struct {
    uint8_t Type[2];
    uint8_t Length[2];
} Header;

// AppLayerFrame structure
typedef struct {
    Header header;
    uint8_t payload[PAYLOAD_MAX_SIZE];
    uint8_t crc; // CRC-8 checksum
} AppLayerFrame;

// Function to create an AppLayerFrame
void createFrame(AppLayerFrame *frame, uint8_t type1, uint8_t type2, uint8_t *data, uint16_t length) {
    frame->header.Type[0] = type1;
    frame->header.Type[1] = type2;
    frame->header.Length[0] = (length >> 8) & 0xFF; // High byte
    frame->header.Length[1] = length & 0xFF;       // Low byte
    memcpy(frame->payload, data, length);
    frame->crc = crc8(frame); // Calculate CRC for the entire frame
}

// Function to calculate CRC-8 checksum
uint8_t crc8(const AppLayerFrame *frame) {
    uint8_t crc = 0x00; // Initial value
    for (size_t i = 0; i < HEADER_SIZE + (frame->header.Length[0] << 8 | frame->header.Length[1]); i++) {
        uint8_t byte = ((uint8_t*)frame)[i]; // Access bytes of the frame
        crc ^= byte; // XOR with the current byte
        for (int j = 0; j < 8; j++) {
            if (crc & 0x80) {
                crc = (crc << 1) ^ 0x1D; // Polynomial for CRC-8
            } else {
                crc <<= 1;
            }
        }
    }
    return crc; // Final CRC value
}

// Function to display the frame
void displayFrame(const AppLayerFrame *frame) {
    printf("Header Type: %02X %02X\n", frame->header.Type[0], frame->header.Type[1]);
    printf("Payload Length: %d\n", (frame->header.Length[0] << 8) | frame->header.Length[1]);
    printf("Payload: ");
    for (int i = 0; i < (frame->header.Length[0] << 8 | frame->header.Length[1]); i++) {
        printf("%02X ", frame->payload[i]);
    }
    printf("\n");
    printf("CRC-8 Checksum: %02X\n", frame->crc);
}

int main() {
    AppLayerFrame frame;
    uint8_t data[] = {0xDE, 0xAD, 0xBE, 0xEF}; // Example payload

    // Create a frame with type 0x01 and payload data
    createFrame(&frame, 0x01, 0x02, data, sizeof(data));

    // Display the constructed frame
    displayFrame(&frame);

    return 0;
}
```

### 说明

1. **Header 结构体**：定义数据头部信息。
2. **AppLayerFrame 结构体**：包含头部、有效载荷和 CRC-8 校验码。
3. **createFrame 函数**：在创建帧时计算并存储 CRC-8 校验码。
4. **crc8 函数**：计算给定 `AppLayerFrame` 的 CRC-8 校验码，使用多项式 `0x1D`。
5. **displayFrame 函数**：显示帧的信息，包括类型、有效载荷和 CRC 校验码。
