<!--
 Copyright (C) 2024 wwhai

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
# ATK-LORA-01

ATK-LORA-01 一款体积小、微功率、低功耗、高性能远距离 LORA 无线串口模块。模块设计是采用高效的 ISM 频段射频 SX1278 扩频芯片，模块的工作频率在 410Mhz~441Mhz，以 1Mhz 频率为步进信道，共 32 个信道。可通过 AT 指令在线修改串口速率，发射功率，空中速率，工作模式等各种参数，并且支持固件升级功能。

## Topic
| Topic               | Args | Result(HTTP返回JSON里的data字段) |
| ------------------- | ---- | -------------------------------- |
| lora.atk01.cmd.send | ""   | String                           |

## 数据协议
### 协议结构
```c
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

```

## Arduino示例
要在Arduino中实现一个基于上述结构体定义的简单协议，用于上报A0和A1模拟引脚的数据，你可以参考以下代码实例。这里假设你使用的是标准的Arduino Uno板，它具有A0和A1模拟输入引脚。

```cpp
#include <Arduino.h>
#include <Wire.h> // For I2C communication

#define HEADER_SIZE 4
#define PAYLOAD_MAX_SIZE 256

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

// Function to calculate the CRC8 of a buffer
uint8_t crc8(const uint8_t* data, size_t len);

void setup() {
    Serial.begin(9600); // Initialize serial for debugging
    Wire.begin(); // Initialize I2C for communication
}

void loop() {
    static int counter = 0; // Counter to simulate data change
    AppLayerFrame frame;

    // Prepare frame header
    frame.header.Type[0] = 'D'; // Data type identifier
    frame.header.Type[1] = 'T';
    frame.header.Length[0] = 2; // Payload length
    frame.header.Length[1] = 0;

    // Prepare payload
    frame.payload[0] = analogRead(A0); // Read A0
    frame.payload[1] = analogRead(A1); // Read A1
    frame.crc = crc8((uint8_t*)&frame, HEADER_SIZE + 2); // Calculate CRC8 over header and payload

    // Send frame via I2C
    Wire.write((uint8_t*)&frame, HEADER_SIZE + PAYLOAD_MAX_SIZE + 1); // Send full frame

    delay(1000); // Delay to avoid flooding the bus
    counter++; // Increment counter
}

uint8_t crc8(const uint8_t* data, size_t len) {
    const uint8_t POLYNOMIAL = 0x07; // CRC-8 polynomial
    uint8_t crc = 0;

    for (size_t pos = 0; pos < len; pos++) {
        crc ^= *data++;
        for (uint8_t bit = 0x80; bit > 0; bit >>= 1) {
            if (crc & bit)
                crc = (crc << 1) ^ POLYNOMIAL;
            else
                crc <<= 1;
        }
    }

    return crc;
}
```

在这个示例中，我们定义了结构体`Header`和`AppLayerFrame`。在`loop`函数内，我们读取A0和A1模拟引脚的数据，并将其放入`payload`字段。然后计算CRC-8校验值并存储在`crc`字段中。最后通过I2C发送整个帧。
