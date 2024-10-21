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

# 外部通信模块管理器

本项目旨在设计RHILEX平台的外设接入框架，以解决标准及非标准硬件（如网卡、蓝牙、WIFI、Lora、Zigbee、433M、2.4G等）在RHILEX内部集成的问题。该框架专注于克服标准软件与各类硬件之间交互的难题。框架内的外设多数为板载固定配置，不支持动态扩展。

## 环境参数

以下列出了支持的外部通信模块及其对应的标识：

```sh
[
    # WIFI_SUPPORT: 支持无线局域网通信，用于高速数据传输和互联网接入。
    'WIFI_SUPPORT',
    # BLC_SUPPORT: 支持低功耗蓝牙通信，适用于近距离低速率数据交换和设备配对。
    'BLC_SUPPORT',
    # BLE_SUPPORT: 支持蓝牙低能耗技术，适合电池供电设备的低功耗数据传输。
    'BLE_SUPPORT',
    # ZIGBEE_SUPPORT: 支持Zigbee协议，一种低速短距离传输的无线网上协议，适用于智能家居和工业自动化。
    'ZIGBEE_SUPPORT',
    # RF24g_SUPPORT: 支持2.4GHz射频通信，常用于无线传感器网络和小型设备间的数据传输。
    'RF24g_SUPPORT',
    # RF433M_SUPPORT: 支持433MHz射频通信，具有较好的穿透能力和较远的传输距离，适用于远距离无线控制。
    'RF433M_SUPPORT',
    # MN4G_SUPPORT: 支持4G移动通信网络，提供高速移动宽带连接。
    'MN4G_SUPPORT',
    # MN5G_SUPPORT: 支持5G移动通信网络，提供更高的数据传输速度和更低的延迟。
    'MN5G_SUPPORT',
    # NBIoT_SUPPORT: 支持窄带物联网通信，专为大规模物联网设备设计的低功耗广域网络技术。
    'NBIoT_SUPPORT',
    # LORA_SUPPORT: 支持LoRa长距离低功耗无线通信技术，适用于广域物联网应用。
    'LORA_SUPPORT',
    # LORA_WAN_SUPPORT: 支持LoRaWAN，一种基于LoRa技术的低功耗广域网络协议，用于物联网设备的长距离通信。
    'LORA_WAN_SUPPORT',
    # IR_SUPPORT: 支持红外线通信，常用于遥控器和短距离数据传输。
    'IR_SUPPORT',
    # BEEP_SUPPORT: 支持蜂鸣器信号，通常用于简单的音频提示或警报。
    'BEEP_SUPPORT'
]
```

## 交互流程

### 输入请求

客户端发送一个JSON格式的请求，包含设备名称、主题和参数：

```json
{
    "name": "EC200A-4G-DTU",
    "topic": "mn4g.ec200a.info.csq",
    "args": ""
}
```

### 返回响应

服务器处理请求后，返回一个JSON格式的响应，包含状态码、消息和数据：

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "name": "EC200A-4G-DTU",
        "topic": "mn4g.ec200a.info.csq",
        "args": "",
        "result": "{\"cops\":\"CMCC\",\"csq\":30,\"iccid\":\"11223344556677\"}"
    }
}
```

响应中的`result`字段包含了具体的设备信息，如运营商、信号强度和SIM卡ID。

## 命令行参数
启动的时候需要带上参数：
```ini
KEY=value
```
其中`KEY`是模块的名称，`value`是ini里面的`transceiver.*`配置，比如下面这个：
```ini
[transceiver.atk01]
# Address: Device is on COM3 serial port for communication
address = /dev/ttyUSB0
# io_timeout: Timeout for I/O ops (30 sec), prevents indefinite waiting
io_timeout = 30
# at_timeout: Timeout for AT cmds (200 ms), adjusts responsiveness
at_timeout = 200
# BaudRate: Data transfer speed set to standard 9600 baud
baudrate = 9600
# DataBits: Each character uses 8 bits for transmission
data_bits = 8
# Parity: parity ('N' 'O' 'D') check, if additional bits for error detection
parity = N
# Stopbits: Single stop bit (1) marks end of each character transmission
stop_bits = 1
# Transport Protocol: 1|2|3, goto homepage for detail
transport_protocol = 1
```
启动指令 `LORA_SUPPORT=atk01` 表示启用atk01这个模块，设备位于`/dev/ttyUSB0`下，使用的是固定报文协议格式。更多请参考文档。