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

// 生成 Getter、Setter 和 New 方法
func generateGetSetMethods(v interface{}) string {
	t := reflect.TypeOf(v)
	var buf bytes.Buffer

	// 生成 New 函数
	buf.WriteString(fmt.Sprintf("func New%s(", t.Name()))
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldType := field.Type

		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%s %s", strings.ToLower(string(fieldName[0])), fieldType))
	}
	buf.WriteString(fmt.Sprintf(") *%s {\n", t.Name()))
	buf.WriteString(fmt.Sprintf("    return &%s{\n", t.Name()))
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		buf.WriteString(fmt.Sprintf("        %s: %s,\n", fieldName, strings.ToLower(string(fieldName[0]))))
	}
	buf.WriteString("    }\n")
	buf.WriteString("}\n\n")

	// 生成 Getter 和 Setter
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldType := field.Type

		// 生成 Getter
		getterName := "Get" + fieldName
		buf.WriteString(fmt.Sprintf("func (p *%s) %s() %s {\n", t.Name(), getterName, fieldType))
		buf.WriteString(fmt.Sprintf("    return p.%s\n", fieldName))
		buf.WriteString("}\n\n")

		// 生成 Setter
		setterName := "Set" + fieldName
		buf.WriteString(fmt.Sprintf("func (p *%s) %s(value %s) {\n", t.Name(), setterName, fieldType))
		buf.WriteString(fmt.Sprintf("    p.%s = value\n", fieldName))
		buf.WriteString("}\n\n")
	}

	return buf.String()
}

func GenerateCode(v interface{}) string {
	code := generateGetSetMethods(v)
	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		log.Fatalf("Failed to format code: %v", err)
	}
	return string(formattedCode)
}
