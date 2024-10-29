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
# M-Bus

M-Bus（Meter-Bus）是一种用于远程读取计量设备（如水表、电表、气表等）的欧洲标准（EN 13757-1）通信协议。M-Bus协议定义了物理层、数据链路层和网络层的通信规范。以下是一个简化的M-Bus报文规范概要：
## 报文结构
M-Bus报文由以下几个部分组成：
1. **起始字符（Start character）**:
   - 始终为 `0x68`。
2. **长度域（Length field）**:
   - 包含控制域和数据域的总长度，不包括起始字符、长度域和校验和。
3. **控制域（Control field）**:
   - 包含地址信息和控制字节。
4. **地址域（Address field）**:
   - 标识发送或接收设备的地址。
5. **数据域（Data field）**:
   - 包含实际的数据，如计量读数。
6. **校验和（Checksum）**:
   - 是整个报文（从长度域开始到数据域结束）的8位累加和的反码。
## 报文格式
```
+--------+--------+--------+--------+--------+--------+--------+--------+
| Start  | Length | Control| Address| Data   | Data   | ...    | Checksum|
| 0x68   |        | Field  | Field  | Field  | Field  |        |         |
+--------+--------+--------+--------+--------+--------+--------+--------+
```
## 报文示例
以下是M-Bus报文的一个简单示例：
```
68 08 08 01 72 02 2F 16
```
- `68`：起始字符
- `08`：长度域（控制域和数据域的总长度为8字节）
- `08`：控制域（通常包含地址信息和控制字节）
- `01`：地址域（设备地址）
- `72 02`：数据域（示例数据）
- `2F`：校验和（累加和的反码）
## 控制域
控制域的格式如下：
```
+--------+--------+
| C_field| C_field|
| (MSB)  | (LSB)  |
+--------+--------+
```
- **C_field (MSB)**: 通常包含地址信息。
- **C_field (LSB)**: 包含帧类型（如短帧、长帧）、传输模式（如确认模式、无确认模式）等信息。
## 地址域
地址域可以是单字节或多字节，取决于实际设备地址的长度。
## 数据域
数据域包含实际的数据，例如：
- 计量值
- 状态信息
- 事件记录
## 校验和
校验和计算方法如下：
1. 从长度域开始到数据域的最后一个字节，计算所有字节的累加和。
2. 取累加和的低8位。
3. 计算累加和的反码（即将每个位取反）。
## 注意
以上仅是一个简化的M-Bus报文规范概要。完整的M-Bus协议规范包含了更多的细节，如不同类型的帧（如SND_NKE、REQ_UD2等），错误处理，以及各种特殊情况的通信流程。如需详细规范，请参考EN 13757-1等相关的欧洲标准文档。
