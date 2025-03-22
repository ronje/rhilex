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
package usbmonitor

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// USBMonitor 结构体
type USBMonitor struct {
	watcher  *fsnotify.Watcher
	callback func(event, device string)
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewUSBMonitor 创建 USBMonitor 实例
func NewUSBMonitor() (*USBMonitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %v", err)
	}

	return &USBMonitor{
		watcher:  watcher,
		stopChan: make(chan struct{}),
	}, nil
}

// Start 监听 USB 串口设备变化
func (u *USBMonitor) Start() error {
	usbPath := "/dev" // 串口设备通常位于 /dev 目录下，如 /dev/ttyUSB0, /dev/ttyACM0

	// 监听 /dev 目录
	err := u.watcher.Add(usbPath)
	if err != nil {
		return fmt.Errorf("failed to watch %s: %v", usbPath, err)
	}

	u.wg.Add(1)
	go u.monitor()

	return nil
}

// Stop 停止监听
func (u *USBMonitor) Stop() {
	close(u.stopChan)
	u.watcher.Close()
	u.wg.Wait()
}

// Callback 设置回调函数
func (u *USBMonitor) Callback(callback func(event, device string)) {
	u.callback = callback
}

// 监听 USB 串口设备的插拔事件
func (u *USBMonitor) monitor() {
	defer u.wg.Done()
	for {
		select {
		case <-u.stopChan:
			return
		case event, ok := <-u.watcher.Events:
			if !ok {
				return
			}

			// 仅关注 USB 串口设备（如 ttyUSB* 和 ttyACM*）
			if !isUSBSerialDevice(event.Name) {
				continue
			}

			var eventType string
			if event.Op&fsnotify.Create != 0 {
				eventType = "ADDED"
			} else if event.Op&fsnotify.Remove != 0 {
				eventType = "REMOVED"
			}

			if eventType != "" && u.callback != nil {
				u.callback(eventType, event.Name)
			}
		case err, ok := <-u.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

// 判断是否为 USB 串口设备
func isUSBSerialDevice(device string) bool {
	base := filepath.Base(device)
	return strings.HasPrefix(base, "ttyUSB") || strings.HasPrefix(base, "ttyACM")
}
