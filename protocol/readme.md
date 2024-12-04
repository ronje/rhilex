# 通用协议指南

## 1. 概述
本协议设计用于简化设备间的数据通信，确保数据传输的稳定性和可靠性。协议使用固定格式的帧进行数据封装，并包括了包头、有效负载和校验机制。通过此协议，设备可以在串口通信中有效地进行数据交换，确保数据的完整性和正确性。

### 主要特点：
- **包头**：固定格式，便于数据帧的解析。
- **数据长度**：每个数据包的长度是固定的，并包含了有效负载数据长度信息。
- **校验和**：通过 CRC 校验来保证数据的完整性。
- **错误处理**：在解析过程中，对于恶意数据包或数据长度不一致的情况，会采取丢弃或恢复的机制，避免对正常数据包的干扰。

## 2. 数据帧结构

### 2.1 包头结构 (`Header`)

包头结构包含以下字段：

- **Id**：2字节，标识帧的唯一ID。
- **Type**：2字节，表示数据帧的类型。
- **Length**：2字节，表示数据帧的有效负载（Payload）的长度。

```go
type Header struct {
    Id     [2]byte  // 包ID
    Type   [2]byte  // 包类型
    Length [2]byte  // 有效负载长度
}
```

### 2.2 数据帧结构 (`AppLayerFrame`)

数据帧包括以下几个部分：

- **Delimiter**：2字节的帧起始标识符，用于同步帧的开始。
- **Header**：包头，包含 `Id`、`Type`、`Length`。
- **Payload**：数据负载，实际传输的数据。
- **CrcSum**：4字节的 CRC 校验和，用于确保数据在传输过程中的完整性。
- **ReverseDelimiter**：2字节的帧结束标识符。

```go
type AppLayerFrame struct {
    Delimiter        [2]byte  // 起始标识符
    Header           Header   // 包头
    Payload          []byte   // 有效负载
    CrcSum           [4]byte  // CRC 校验和
    ReverseDelimiter [2]byte  // 结束标识符
}
```

### 2.3 帧的组成
每个数据帧的组成如下：
```
[Delimiter] [Header] [Payload] [CrcSum] [ReverseDelimiter]
```
- **Delimiter**：帧开始标识符，固定为 `0xAA 0xBB`。
- **Header**：包头，包含包的 ID、类型和有效负载长度。
- **Payload**：实际传输的数据。
- **CrcSum**：对整个数据帧（包括包头、有效负载）计算的 CRC 校验值，确保数据在传输过程中未被篡改。
- **ReverseDelimiter**：帧结束标识符，固定为 `0xBB 0xAA`。

### 2.4 帧示例
假设包头是 `0x01 0x02`，数据长度为 `0x05`，有效负载是 `0x11 0x22 0x33 0x44 0x55`，CRC 校验和为 `0xDEADBEEF`。完整的帧数据如下：

```
[0xAA 0xBB] [0x01 0x02 0x05] [0x11 0x22 0x33 0x44 0x55] [0xDE 0xAD 0xBE 0xEF] [0xBB 0xAA]
```

## 3. 编解码规则

### 3.1 编码规则

**AppLayerFrame 编码 (`Encode`)**：
1. 首先构造帧头 (`Header`)。
2. 将帧头的各个字段（`Id`、`Type`、`Length`）编码成字节序列。
3. 将有效负载（`Payload`）按照实际数据直接附加。
4. 计算整个帧数据（包括包头和负载）的 CRC 校验和。
5. 将 CRC 校验和附加到数据帧末尾。
6. 将帧开始标识符（`Delimiter`）和帧结束标识符（`ReverseDelimiter`）加入数据帧。

**编码示例**：

```go
func (frame *AppLayerFrame) Encode() ([]byte, error) {
    // 编码包头
    headerBytes, err := frame.Header.Encode()
    if err != nil {
        return nil, err
    }

    // 计算CRC
    crc := crc32.ChecksumIEEE(append(headerBytes, frame.Payload...))

    // 创建完整数据帧
    data := append(frame.Delimiter[:], headerBytes...)
    data = append(data, frame.Payload...)
    data = append(data, crc...)
    data = append(data, frame.ReverseDelimiter[:]...)

    return data, nil
}
```

### 3.2 解码规则

**AppLayerFrame 解码 (`Decode`)**：
1. 先检查帧的起始标识符（`Delimiter`）是否正确，确认数据帧的起始位置。
2. 解码包头并解析 `Id`、`Type` 和 `Length`。
3. 根据 `Length` 字段从数据中读取有效负载（`Payload`）。
4. 提取并验证 CRC 校验和，确保数据没有被篡改。
5. 检查帧的结束标识符（`ReverseDelimiter`）是否正确，确保数据帧结束。

**解码示例**：

```go
func DecodeAppLayerFrame(data []byte) (AppLayerFrame, error) {
    // 验证起始标识符
    if data[0] != 0xAA || data[1] != 0xBB {
        return AppLayerFrame{}, fmt.Errorf("invalid delimiter")
    }

    // 解码包头
    header, err := DecodeHeader(data[2:])
    if err != nil {
        return AppLayerFrame{}, err
    }

    // 获取有效负载
    payload := data[6 : 6+header.Length[0]]

    // 验证 CRC 校验和
    crc := crc32.ChecksumIEEE(data[:len(data)-4])
    if !bytes.Equal(crc, data[len(data)-4:]) {
        return AppLayerFrame{}, fmt.Errorf("CRC check failed")
    }

    // 验证结束标识符
    if data[len(data)-2] != 0xBB || data[len(data)-1] != 0xAA {
        return AppLayerFrame{}, fmt.Errorf("invalid reverse delimiter")
    }

    return AppLayerFrame{Header: header, Payload: payload}, nil
}
```

## 4. 错误处理与恢复机制

在数据解析过程中，可能会遇到以下情况：
1. **无效包头**：当包头的起始标识符不正确时，丢弃当前数据并重新同步。
2. **数据长度异常**：当 `Length` 字段的值与实际数据长度不符时，丢弃当前帧，恢复到下一个有效包。
3. **CRC 校验失败**：如果校验和失败，认为数据已损坏，丢弃该帧并恢复。

### 4.1 错误恢复机制
为了提高协议的鲁棒性，当遇到恶意数据包或帧解析失败时，协议会将错误的数据暂时存储到错误缓冲区（`errorBuffer`）。在下次数据接收时，协议会尝试恢复缓冲区中的错误数据并继续解析。此机制确保了恶意数据不会阻止正常数据包的接收。

## 5. 总结

本协议提供了一种简单而强大的数据帧格式，能够确保数据传输的完整性和可靠性。通过固定的包头结构、数据长度字段、CRC 校验以及错误恢复机制，设备可以高效地进行数据通信并避免因恶意数据包导致的通信中断。