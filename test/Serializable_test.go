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

package test

import (
	"testing"

	"github.com/hootrhino/rhilex/protocol"
)

// go test -timeout 30s -run ^TestGenerateGenericFrameCode github.com/hootrhino/rhilex/test -v -count=1
func TestGenerateGenericFrameCode(t *testing.T) {
	type CJT188Frame0x01 struct {
		Start        byte    // 帧起始符
		MeterType    byte    // 仪表类型
		Address      [7]byte // 地址域
		CtrlCode     byte    // 控制码
		DataLength   byte    // 数据长度域
		DataType     [2]byte // 数据长度域
		DataArea     []byte  // 数据域
		SerialNumber byte
		CheckSum     byte // 校验码
		End          byte // 结束符
	}
	f := CJT188Frame0x01{}
	t.Log("============\n", protocol.GenerateCode(f))
}

// 示例结构体
type Person struct {
	Name string
	Age  int
}

// 实现 Validatable 接口
func (p *Person) Validate() error {
	return nil
}

// 设置和获取操作
func (p *Person) SetName(name string) {
	p.Name = name
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) SetAge(age int) {
	p.Age = age
}

func (p *Person) GetAge() int {
	return p.Age
}

// go test -timeout 30s -run ^TestSerializable github.com/hootrhino/rhilex/test -v -count=1

func TestSerializable(t *testing.T) {
	// 创建一个通用结构体
	person := &Person{Name: "Alice", Age: 30}
	t.Log(protocol.GenerateCode(*person))
	genericStruct := protocol.GenericFrame[*Person]{Data: person}

	// 序列化结构体
	data, err := genericStruct.Serialize()
	if err != nil {
		t.Fatalf("Error serializing: %s", err)
	}
	t.Log("Serialized data:", data)

	// 反序列化字节为结构体
	newGenericStruct := protocol.GenericFrame[*Person]{}
	err = newGenericStruct.Deserialize(data)
	if err != nil {
		t.Fatalf("Error deserializing: %s", err)
	}
	newPerson := newGenericStruct.Data
	t.Log("Deserialized Person:", newPerson)

	// 测试 set 和 get 操作
	newPerson.SetName("Bob")
	newPerson.SetAge(25)
	t.Log("Updated Person:", newPerson.GetName(), newPerson.GetAge())

	// 测试CRC不匹配
	invalidData := append(data[:len(data)-2], 0, 1) // 修改CRC
	err = newGenericStruct.Deserialize(invalidData)
	if err != nil {
		t.Log("Validation error during deserialization:", err)
	}
}
