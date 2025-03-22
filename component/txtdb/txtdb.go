// Copyright (C) 2025 wwhai
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

package txtdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TextDB 表示文本数据库结构体
type TextDB struct {
	filePath string
}

// NewTextDB 创建一个新的文本数据库实例
func NewTextDB(filePath string) *TextDB {
	return &TextDB{
		filePath: filePath,
	}
}

// Add 添加数据到数据库
func (db *TextDB) Add(key, value string) error {
	// 检查键是否已存在
	exists, err := db.Exists(key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("key %s already exists", key)
	}

	// 以追加模式打开文件
	file, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入键值对
	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	return err
}

// Delete 从数据库中删除数据
func (db *TextDB) Delete(key string) error {
	// 读取文件内容
	lines, err := db.readAllLines()
	if err != nil {
		return err
	}

	// 查找并删除指定键的行
	newLines := make([]string, 0, len(lines))
	deleted := false
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			deleted = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !deleted {
		return fmt.Errorf("key %s not found", key)
	}

	// 将更新后的内容写回文件
	return db.writeAllLines(newLines)
}

// Update 更新数据库中的数据
func (db *TextDB) Update(key, value string) error {
	// 读取文件内容
	lines, err := db.readAllLines()
	if err != nil {
		return err
	}

	// 查找并更新指定键的行
	updated := false
	for i, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			lines[i] = fmt.Sprintf("%s:%s", key, value)
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("key %s not found", key)
	}

	// 将更新后的内容写回文件
	return db.writeAllLines(lines)
}

// Get 从数据库中获取数据
func (db *TextDB) Get(key string) (string, error) {
	// 读取文件内容
	lines, err := db.readAllLines()
	if err != nil {
		return "", err
	}

	// 查找指定键的值
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("key %s not found", key)
}

// Exists 检查键是否存在于数据库中
func (db *TextDB) Exists(key string) (bool, error) {
	// 读取文件内容
	lines, err := db.readAllLines()
	if err != nil {
		return false, err
	}

	// 查找指定键
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return true, nil
		}
	}

	return false, nil
}

// readAllLines 读取文件的所有行
func (db *TextDB) readAllLines() ([]string, error) {
	file, err := os.Open(db.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeAllLines 将所有行写入文件
func (db *TextDB) writeAllLines(lines []string) error {
	file, err := os.OpenFile(db.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
