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

package ossupport

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// GPIO 结构体表示一个 GPIO 引脚
type GPIO struct {
	Number int
}

// NewGPIO 创建一个新的 GPIO 实例
func NewGPIO(number int) *GPIO {
	return &GPIO{
		Number: number,
	}
}

// Init 初始化 GPIO，如果未导出则导出
func (g *GPIO) Init() error {
	// 检查 GPIO 是否已经导出
	exported, err := g.isExported()
	if err != nil {
		return err
	}

	if !exported {
		// 未导出则进行导出操作
		err := g.export()
		if err != nil {
			return err
		}
		// 等待一段时间，确保导出操作完成
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// isExported 检查 GPIO 是否已经导出
func (g *GPIO) isExported() (bool, error) {
	gpioPath := filepath.Join("/sys/class/gpio", fmt.Sprintf("gpio%d", g.Number))
	_, err := os.Stat(gpioPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// export 导出 GPIO
func (g *GPIO) export() error {
	exportFile := "/sys/class/gpio/export"
	err := os.WriteFile(exportFile, []byte(strconv.Itoa(g.Number)), 0644)
	if err != nil {
		return fmt.Errorf("failed to export GPIO %d: %v", g.Number, err)
	}
	return nil
}

// SetDirection 设置 GPIO 的方向（in 或 out）
func (g *GPIO) SetDirection(direction string) error {
	directionPath := filepath.Join("/sys/class/gpio", fmt.Sprintf("gpio%d", g.Number), "direction")
	err := os.WriteFile(directionPath, []byte(direction), 0644)
	if err != nil {
		return fmt.Errorf("failed to set direction for GPIO %d: %v", g.Number, err)
	}
	return nil
}

// SetValue 设置 GPIO 的值（0 或 1）
func (g *GPIO) SetValue(value int) error {
	valuePath := filepath.Join("/sys/class/gpio", fmt.Sprintf("gpio%d", g.Number), "value")
	err := os.WriteFile(valuePath, []byte(strconv.Itoa(value)), 0644)
	if err != nil {
		return fmt.Errorf("failed to set value for GPIO %d: %v", g.Number, err)
	}
	return nil
}

// Unexport 取消导出 GPIO
func (g *GPIO) Unexport() error {
	unexportFile := "/sys/class/gpio/unexport"
	err := os.WriteFile(unexportFile, []byte(strconv.Itoa(g.Number)), 0644)
	if err != nil {
		return fmt.Errorf("failed to unexport GPIO %d: %v", g.Number, err)
	}
	return nil
}

func TestGpio() {
	// 创建一个 GPIO 实例，假设使用 GPIO 17
	gpio := NewGPIO(17)

	// 初始化 GPIO
	err := gpio.Init()
	if err != nil {
		fmt.Printf("Failed to initialize GPIO: %v\n", err)
		return
	}

	// 设置 GPIO 方向为输出
	err = gpio.SetDirection("out")
	if err != nil {
		fmt.Printf("Failed to set GPIO direction: %v\n", err)
		return
	}

	// 设置 GPIO 值为高电平
	err = gpio.SetValue(1)
	if err != nil {
		fmt.Printf("Failed to set GPIO value: %v\n", err)
		return
	}

	fmt.Println("GPIO set to high for 5 seconds...")
	time.Sleep(5 * time.Second)

	// 设置 GPIO 值为低电平
	err = gpio.SetValue(0)
	if err != nil {
		fmt.Printf("Failed to set GPIO value: %v\n", err)
		return
	}

	// 取消导出 GPIO
	err = gpio.Unexport()
	if err != nil {
		fmt.Printf("Failed to unexport GPIO: %v\n", err)
		return
	}

	fmt.Println("GPIO operation completed.")
}
