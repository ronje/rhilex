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

### 通用二进制数据解析器

#### 1. 概述

通用二进制数据解析器是一种灵活的工具，它可以根据用户定义的表达式，解析任意二进制数据。该解析器允许用户指定字段的名称、位长度、数据类型以及字节序，以实现对复杂数据结构的轻松解析。解析结果最终会封装为一个结构体，并可导出为 JSON 格式。

#### 2. 特性

- **自定义表达式**：支持通过表达式定义需要解析的数据结构。
- **位级控制**：字段长度以位为单位进行定义，支持位级精度解析。
- **数据类型支持**：支持 `int`、`float` 和 `string` 等常用数据类型。
- **字节序控制**：支持大端序（Big Endian，BE）和小端序（Little Endian，LE）。
- **解析结果封装**：解析后的数据以键值对的形式封装，并可转换为 JSON 格式输出。

#### 3. 表达式格式

解析器通过表达式来定义如何解析输入的二进制数据。表达式格式为：

```
Key:Length:Type:Endian; Key:Length:Type:Endian; ...
```

- **Key**：字段的名称，作为输出 JSON 中的键。
- **Length**：字段的长度，以位（bits）为单位。
- **Type**：字段的数据类型，支持 `int`（整数）、`float`（浮点数）、`string`（字符串）。
- **Endian**：字节序，支持 `BE`（大端序）和 `LE`（小端序）。

#### 4. 语法规则（BNF 范式）

以下是表达式的语法定义，采用 BNF 范式：

```
<expression>    ::= <field> ";" <expression> | <field>
<field>         ::= <key> ":" <length> ":" <type> ":" <endian>
<key>           ::= <identifier>
<length>        ::= <number>
<type>          ::= "int" | "float" | "string"
<endian>        ::= "BE" | "LE"
<identifier>    ::= 字母开头的字母数字序列
<number>        ::= 1-9 的数字序列，表示位的长度
```

#### 5. 示例

##### 5.1 表达式示例

```
ID:32:int:BE; Name:40:string:BE; Age:16:int:LE
```

该表达式表示：
- `ID`：一个 32 位（4 字节）的整数，使用大端序。
- `Name`：一个 40 位（5 字节）的字符串，使用大端序。
- `Age`：一个 16 位（2 字节）的整数，使用小端序。

##### 5.2 输入数据示例

二进制数据：
```
[0x00, 0x00, 0x00, 0x01, 'A', 'l', 'i', 'c', 'e', 0x20, 0x00]
```

##### 5.3 解析结果示例

解析器将会根据表达式对输入数据进行解析，得到如下结果：

```json
{
  "ID": 1,
  "Name": "Alice",
  "Age": 32
}
```

#### 6. 使用场景

- **网络协议解析**：在处理二进制协议数据时，协议报文通常以位为单位对字段进行定义，通用二进制数据解析器可以通过表达式轻松实现报文解析。
- **文件格式解析**：一些二进制文件格式（如图片、音频等）具有固定结构，通过表达式描述文件结构可以快速提取文件中的关键信息。
- **嵌入式设备数据处理**：在嵌入式系统中，传感器或控制器常常以二进制形式传输数据，该解析器能够方便地解析这些数据。

#### 7. 注意事项

- **字节序**：解析时请确保为每个字段选择正确的字节序，尤其是在跨平台传输数据时，字节序不一致可能导致解析错误。
- **字段对齐**：由于字段的长度是以位为单位，注意多个字段之间的位对齐问题，解析器会自动按字节进行对齐处理。

### 8. 实例展示：Modbus 03 报文解析示例

#### 8.1 Modbus 03 功能码报文简介

Modbus 协议广泛应用于工业控制和自动化系统中，其中 **03 功能码** 用于读取保持寄存器的数据。Modbus 03 报文通常包括以下字段：

- **设备地址（1 字节）**：标识从站设备的地址。
- **功能码（1 字节）**：表示操作的类型，对于读取保持寄存器，功能码为 `03`。
- **字节数（1 字节）**：返回的数据字节数。
- **寄存器数据（可变长度，n 字节）**：实际的寄存器值，通常为多个 16 位整数。
- **CRC 校验（2 字节）**：校验码，用于校验报文的完整性。

#### 8.2 表达式定义

为了解析 Modbus 03 报文，我们可以设计以下表达式来匹配报文中的字段。假设我们需要解析以下字段：
- **DeviceAddress**：设备地址，8 位（1 字节）。
- **FunctionCode**：功能码，8 位（1 字节）。
- **ByteCount**：字节数，8 位（1 字节）。
- **Register1**：第一个寄存器数据，16 位（2 字节），使用大端序。
- **Register2**：第二个寄存器数据，16 位（2 字节），使用大端序。
- **CRC**：CRC 校验，16 位（2 字节），使用小端序。

我们可以使用如下表达式：
```
DeviceAddress:8:int:BE; FunctionCode:8:int:BE; ByteCount:8:int:BE; Register1:16:int:BE; Register2:16:int:BE; CRC:16:int:LE
```

#### 8.3 示例数据

假设我们收到的 Modbus 03 功能码报文数据如下：

```
[0x01, 0x03, 0x04, 0x00, 0x0A, 0x00, 0x14, 0xA1, 0xB2]
```

报文结构解读如下：
- **DeviceAddress**：`0x01`，表示设备地址 1。
- **FunctionCode**：`0x03`，表示功能码 03，读取保持寄存器。
- **ByteCount**：`0x04`，表示返回的数据长度为 4 字节（两个寄存器）。
- **Register1**：`0x000A`，表示第一个寄存器的值为 10。
- **Register2**：`0x0014`，表示第二个寄存器的值为 20。
- **CRC**：`0xA1B2`，小端序校验码。

#### 8.4 解析结果

根据定义的表达式和示例数据，解析器将输出以下结果：

```json
{
  "DeviceAddress": 1,
  "FunctionCode": 3,
  "ByteCount": 4,
  "Register1": 10,
  "Register2": 20,
  "CRC": 45537
}
```

#### 8.5 解析步骤解析

1. **解析设备地址**：数据流的第一个字节 `0x01` 被解析为 `DeviceAddress`，其值为 `1`。
2. **解析功能码**：第二个字节 `0x03` 被解析为 `FunctionCode`，其值为 `3`，表示读取寄存器的操作。
3. **解析字节数**：第三个字节 `0x04` 被解析为 `ByteCount`，表明接下来返回的寄存器数据占 4 字节。
4. **解析寄存器数据**：接下来的 4 字节分别被解析为两个寄存器 `Register1` 和 `Register2`，其值分别为 `10` 和 `20`。
5. **解析 CRC 校验码**：最后 2 字节 `0xA1B2` 使用小端序解析为 `45537`。
