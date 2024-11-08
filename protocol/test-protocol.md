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

以下是几个测试示例，用于验证协议的编码、解码、以及错误处理和恢复机制。

### 测试示例 1：正常数据帧

此测试示例验证正常情况下数据帧的编码和解码流程。

**输入数据：**

```go
header := Header{
    Id:     [2]byte{0x01, 0x02},
    Type:   [2]byte{0x03, 0x04},
    Length: [2]byte{0x00, 0x05}, // Payload 长度为 5
}

payload := []byte{0x11, 0x22, 0x33, 0x44, 0x55}

frame := AppLayerFrame{
    Delimiter:        [2]byte{0xAA, 0xBB},
    Header:           header,
    Payload:          payload,
    CrcSum:           calculateCRC(header, payload), // 假设 calculateCRC 是计算 CRC 的函数
    ReverseDelimiter: [2]byte{0xBB, 0xAA},
}

// 编码数据帧
encodedData, err := frame.Encode()
if err != nil {
    fmt.Println("编码失败：", err)
    return
}

// 解码数据帧
decodedFrame, err := DecodeAppLayerFrame(encodedData)
if err != nil {
    fmt.Println("解码失败：", err)
    return
}

// 验证解码后的数据
fmt.Printf("解码后的 Header: %+v\n", decodedFrame.Header)
fmt.Printf("解码后的 Payload: %+v\n", decodedFrame.Payload)
```

**期望输出：**
- 解码后的 `Header` 与原始 `Header` 相同。
- 解码后的 `Payload` 与原始 `Payload` 相同。

### 测试示例 2：错误的起始标识符

此测试验证当起始标识符不正确时，解码器会丢弃该帧并返回错误。

**输入数据：**

```go
data := []byte{
    0xAB, 0xCD, // 错误的起始标识符
    0x01, 0x02, 0x03, 0x04, 0x00, 0x05,
    0x11, 0x22, 0x33, 0x44, 0x55,
    0xDE, 0xAD, 0xBE, 0xEF,
    0xBB, 0xAA,
}

// 解码数据帧
_, err := DecodeAppLayerFrame(data)
if err != nil {
    fmt.Println("解码失败：", err) // 应输出错误信息，例如“invalid delimiter”
}
```

**期望输出：**
- 解码失败，输出 `invalid delimiter` 错误信息。

### 测试示例 3：数据长度异常

此测试验证当 `Length` 字段的值与实际 `Payload` 长度不符时，解码器会丢弃该帧并返回错误。

**输入数据：**

```go
header := Header{
    Id:     [2]byte{0x01, 0x02},
    Type:   [2]byte{0x03, 0x04},
    Length: [2]byte{0x00, 0x07}, // 声明 Payload 长度为 7
}

payload := []byte{0x11, 0x22, 0x33, 0x44, 0x55} // 实际 Payload 长度为 5

frame := AppLayerFrame{
    Delimiter:        [2]byte{0xAA, 0xBB},
    Header:           header,
    Payload:          payload,
    CrcSum:           calculateCRC(header, payload),
    ReverseDelimiter: [2]byte{0xBB, 0xAA},
}

// 编码数据帧
encodedData, _ := frame.Encode()

// 解码数据帧
_, err := DecodeAppLayerFrame(encodedData)
if err != nil {
    fmt.Println("解码失败：", err) // 应输出错误信息，例如“data length mismatch”
}
```

**期望输出：**
- 解码失败，输出 `data length mismatch` 错误信息。

### 测试示例 4：CRC 校验失败

此测试验证当 CRC 校验失败时，解码器会丢弃该帧并返回错误。

**输入数据：**

```go
header := Header{
    Id:     [2]byte{0x01, 0x02},
    Type:   [2]byte{0x03, 0x04},
    Length: [2]byte{0x00, 0x05}, // Payload 长度为 5
}

payload := []byte{0x11, 0x22, 0x33, 0x44, 0x55}

frame := AppLayerFrame{
    Delimiter:        [2]byte{0xAA, 0xBB},
    Header:           header,
    Payload:          payload,
    CrcSum:           [4]byte{0xFF, 0xFF, 0xFF, 0xFF}, // 错误的 CRC 校验和
    ReverseDelimiter: [2]byte{0xBB, 0xAA},
}

// 编码数据帧
encodedData, _ := frame.Encode()

// 解码数据帧
_, err := DecodeAppLayerFrame(encodedData)
if err != nil {
    fmt.Println("解码失败：", err) // 应输出错误信息，例如“CRC check failed”
}
```

**期望输出：**
- 解码失败，输出 `CRC check failed` 错误信息。

### 测试示例 5：误判恢复

此测试验证当数据越界并导致错误读取时，解析器会在错误处理后恢复到下一个有效数据包。

**输入数据：**

```go
// 正常的第一帧
frame1 := AppLayerFrame{
    Delimiter:        [2]byte{0xAA, 0xBB},
    Header:           Header{Id: [2]byte{0x01, 0x02}, Type: [2]byte{0x03, 0x04}, Length: [2]byte{0x00, 0x05}},
    Payload:          []byte{0x11, 0x22, 0x33, 0x44, 0x55},
    CrcSum:           calculateCRC(header, payload),
    ReverseDelimiter: [2]byte{0xBB, 0xAA},
}

// 不完整的恶意帧
malformedData := []byte{0xAA, 0xBB, 0x01, 0x02, 0x03}

// 正常的第二帧
frame2 := AppLayerFrame{
    Delimiter:        [2]byte{0xAA, 0xBB},
    Header:           Header{Id: [2]byte{0x02, 0x03}, Type: [2]byte{0x04, 0x05}, Length: [2]byte{0x00, 0x03}},
    Payload:          []byte{0x66, 0x77, 0x88},
    CrcSum:           calculateCRC(header, payload),
    ReverseDelimiter: [2]byte{0xBB, 0xAA},
}

// 模拟读取并合并数据
data := append(frame1.Encode(), malformedData...)
data = append(data, frame2.Encode()...)

// 解码数据帧，检查恢复机制
frames, err := DecodeMultipleFrames(data)
if err != nil {
    fmt.Println("解码中出现错误：", err)
} else {
    fmt.Println("成功解码帧数量：", len(frames))
    for _, f := range frames {
        fmt.Printf("解码后的帧：%+v\n", f)
    }
}
```

**期望输出：**
- 解码器成功跳过不完整的恶意帧，正确解析 `frame1` 和 `frame2`