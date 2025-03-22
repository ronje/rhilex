<!--
 Copyright (C) 2025 wwhai

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

# Txt文本KV数据库

## 一、概述
本项目实现了一个简单的 Golang 文本数据库，主要用于本地存储 `Map` 形式的数据，以 `Key` 作为主键，支持基本的增删改查操作。数据以纯文本（`.txt`）形式保存，适用于一些对数据存储要求不高、数据量较小的本地应用场景。

## 二、功能特性
1. **数据存储**：将键值对数据以文本形式存储在本地文件中。
2. **主键支持**：使用 `Key` 作为主键，确保数据的唯一性。
3. **增删改查**：支持对数据进行添加、删除、更新和查询操作。

## 三、代码结构
### 1. 主要结构体
```go
type TextDB struct {
    filePath string
}
```
- `TextDB` 结构体表示文本数据库，其中 `filePath` 字段用于指定存储数据的文本文件路径。

### 2. 主要函数和方法
#### 2.1 `NewTextDB` 函数
```go
func NewTextDB(filePath string) *TextDB
```
- **功能**：创建一个新的 `TextDB` 实例。
- **参数**：
  - `filePath`：存储数据的文本文件路径。
- **返回值**：返回一个指向 `TextDB` 结构体的指针。

#### 2.2 `Add` 方法
```go
func (db *TextDB) Add(key, value string) error
```
- **功能**：将键值对添加到数据库中。
- **参数**：
  - `key`：要添加的键。
  - `value`：要添加的值。
- **返回值**：如果添加成功，返回 `nil`；如果键已存在，返回错误信息。

#### 2.3 `Delete` 方法
```go
func (db *TextDB) Delete(key string) error
```
- **功能**：从数据库中删除指定键的键值对。
- **参数**：
  - `key`：要删除的键。
- **返回值**：如果删除成功，返回 `nil`；如果键不存在，返回错误信息。

#### 2.4 `Update` 方法
```go
func (db *TextDB) Update(key, value string) error
```
- **功能**：更新数据库中指定键的值。
- **参数**：
  - `key`：要更新的键。
  - `value`：更新后的值。
- **返回值**：如果更新成功，返回 `nil`；如果键不存在，返回错误信息。

#### 2.5 `Get` 方法
```go
func (db *TextDB) Get(key string) (string, error)
```
- **功能**：从数据库中获取指定键的值。
- **参数**：
  - `key`：要查询的键。
- **返回值**：如果查询成功，返回对应的值和 `nil`；如果键不存在，返回空字符串和错误信息。

#### 2.6 `Exists` 方法
```go
func (db *TextDB) Exists(key string) (bool, error)
```
- **功能**：检查指定键是否存在于数据库中。
- **参数**：
  - `key`：要检查的键。
- **返回值**：如果键存在，返回 `true` 和 `nil`；如果键不存在，返回 `false` 和 `nil`；如果出现错误，返回 `false` 和错误信息。

#### 2.7 `readAllLines` 方法
```go
func (db *TextDB) readAllLines() ([]string, error)
```
- **功能**：读取文件的所有行。
- **返回值**：返回包含文件所有行的字符串切片和可能的错误信息。

#### 2.8 `writeAllLines` 方法
```go
func (db *TextDB) writeAllLines(lines []string) error
```
- **功能**：将所有行写入文件。
- **参数**：
  - `lines`：要写入文件的字符串切片。
- **返回值**：如果写入成功，返回 `nil`；如果出现错误，返回错误信息。

## 四、使用示例
```go
package main

import (
    "fmt"
)

func main() {
    // 创建一个新的文本数据库实例
    db := NewTextDB("data.txt")

    // 添加数据
    err := db.Add("key1", "value1")
    if err != nil {
        fmt.Println("Error adding data:", err)
    }

    // 获取数据
    value, err := db.Get("key1")
    if err != nil {
        fmt.Println("Error getting data:", err)
    } else {
        fmt.Println("Value for key1:", value)
    }

    // 更新数据
    err = db.Update("key1", "newvalue1")
    if err != nil {
        fmt.Println("Error updating data:", err)
    }

    // 再次获取数据
    value, err = db.Get("key1")
    if err != nil {
        fmt.Println("Error getting data:", err)
    } else {
        fmt.Println("Updated value for key1:", value)
    }

    // 删除数据
    err = db.Delete("key1")
    if err != nil {
        fmt.Println("Error deleting data:", err)
    }

    // 尝试获取已删除的数据
    value, err = db.Get("key1")
    if err != nil {
        fmt.Println("Error getting data:", err)
    } else {
        fmt.Println("Value for key1:", value)
    }
}
```

## 五、注意事项
1. **数据格式**：数据以 `key:value` 的形式存储在文本文件中，每行一个键值对。因此，键和值中不能包含冒号 `:`，否则会影响数据的解析。
2. **性能限制**：此实现是一个简单的文本数据库，不适合处理大量数据或高并发场景。在处理大量数据时，读写文件的操作可能会成为性能瓶颈。
3. **文件路径**：确保指定的文件路径具有读写权限，否则可能会导致文件操作失败。
