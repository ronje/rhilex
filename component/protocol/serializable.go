// Copyright (C) 2024 wwhai
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

package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"go/format"
	"hash/crc32"
	"log"
	"reflect"
	"strings"
)

// 定义一个验证接口
type Validatable interface {
	Validate() error // 验证数据有效性
}

// 定义一个序列化接口
type Serializable[T any] interface {
	Serialize() ([]byte, error) // 返回序列化的字节
	Deserialize(data []byte) error
}

// 通用结构体
type GenericFrame[T Validatable] struct {
	Data T
}

// 实现 Serializable 接口
func (g *GenericFrame[T]) Serialize() ([]byte, error) {
	if err := g.Data.Validate(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(g.Data)
	if err != nil {
		return nil, err
	}

	dataBytes := buf.Bytes()
	var finalBuf bytes.Buffer
	length := uint32(len(dataBytes))
	if err := binary.Write(&finalBuf, binary.LittleEndian, length); err != nil {
		return nil, err
	}
	finalBuf.Write(dataBytes)
	crc := crc32.ChecksumIEEE(finalBuf.Bytes())
	if err := binary.Write(&finalBuf, binary.LittleEndian, crc); err != nil {
		return nil, err
	}
	return finalBuf.Bytes(), nil
}

func (g *GenericFrame[T]) Deserialize(data []byte) error {
	if len(data) < 8 { // 4字节长度 + 8字节CRC
		return fmt.Errorf("data is too short to contain length and CRC")
	}

	// 读取长度
	var length uint32
	if err := binary.Read(bytes.NewReader(data[:4]), binary.LittleEndian, &length); err != nil {
		return err
	}

	crcBytes := [4]byte{}
	copy(crcBytes[:], data[len(data)-4:])
	receivedCRC := binary.LittleEndian.Uint32(crcBytes[:])

	expectedCRC := crc32.ChecksumIEEE(data[:len(data)-4])
	if receivedCRC != expectedCRC {
		return fmt.Errorf("CRC does not match, data may be corrupted")
	}
	buf := bytes.NewBuffer(data[4 : len(data)-4])
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&g.Data)
	if err != nil {
		return err
	}

	// 验证反序列化后的数据
	return g.Data.Validate()
}

/**
 * 生成通用序列化代码
 *
 */
func generateCode(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return "Error: provided value is not a struct"
	}

	structName := typ.Name()
	var sb strings.Builder

	// 结构体定义
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		sb.WriteString(fmt.Sprintf("    %s %s\n", field.Name, Uint8ToByte(field.Type)))
	}
	sb.WriteString("}\n\n")

	// Get 和 Set 方法
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		// Get 方法
		sb.WriteString(fmt.Sprintf("func (f *%s) Get%s() %s {\n", structName, field.Name, Uint8ToByte(field.Type)))
		sb.WriteString(fmt.Sprintf("    return f.%s\n", field.Name))
		sb.WriteString("}\n\n")

		// Set 方法
		sb.WriteString(fmt.Sprintf("func (f *%s) Set%s(value %s) {\n", structName, field.Name, Uint8ToByte(field.Type)))
		sb.WriteString(fmt.Sprintf("    f.%s = value\n", field.Name))
		sb.WriteString("}\n\n")
	}

	// String 方法
	sb.WriteString(fmt.Sprintf("func (f %s) String() string {\n", structName))
	sb.WriteString("    bytes, _ := json.Marshal(f)\n")
	sb.WriteString("    return string(bytes)\n")
	sb.WriteString("}\n\n")

	// Validate 方法
	sb.WriteString(fmt.Sprintf("func (f %s) Validate() error {\n", structName))
	sb.WriteString("    return nil\n")
	sb.WriteString("}\n\n")

	// Serialize 方法
	sb.WriteString(fmt.Sprintf("func (f %s) Serialize() ([]byte, error) {\n", structName))
	sb.WriteString(fmt.Sprintf("    genericStruct := protocol.GenericFrame[%s]{Data: f}\n", structName))
	sb.WriteString("    return genericStruct.Serialize()\n")
	sb.WriteString("}\n\n")

	// DeSerialize 函数
	sb.WriteString(fmt.Sprintf("func DeSerialize%s(b []byte) (%s, error) {\n", structName, structName))
	sb.WriteString(fmt.Sprintf("    f := %s{}\n", structName))
	sb.WriteString(fmt.Sprintf("    genericStruct := protocol.GenericFrame[%s]{Data: f}\n", structName))
	sb.WriteString("    return f, genericStruct.Deserialize(b)\n")
	sb.WriteString("}\n")

	return sb.String()
}

/**
 * uint8 -> byte
 *
 */
func Uint8ToByte(Type reflect.Type) string {
	if Type.Kind() == reflect.Uint8 {
		return "byte"
	}
	return Type.Name()
}

/**
 * 输出格式化后的代码
 *
 */
func GenerateCode(v interface{}) string {
	code := generateCode(v)
	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		log.Fatalf("Failed to format code: %v", err)
	}
	return string(formattedCode)
}
