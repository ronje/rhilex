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

package en6400

import (
	"context"
	"fmt"
	"os"
	"time"
)

// FeedWatchdog 向 /dev/watchdog0 写入 0 进行喂狗操作
func FeedWatchdog() error {
	file, err := os.OpenFile("/dev/watchdog0", os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("can not open /dev/watchdog0: %w", err)
	}
	defer file.Close()

	// 写入 0 进行喂狗
	_, err = file.Write([]byte{0})
	if err != nil {
		return fmt.Errorf("write /dev/watchdog0 failed: %w", err)
	}

	return nil
}

// CancelWatchdog 向 /dev/watchdog0 写入 V 取消看门狗功能
func CancelWatchdog() error {
	file, err := os.OpenFile("/dev/watchdog0", os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("can not open /dev/watchdog0: %w", err)
	}
	defer file.Close()

	// 写入 V 取消看门狗
	_, err = file.Write([]byte("V"))
	if err != nil {
		return fmt.Errorf("write /dev/watchdog0 失败: %w", err)
	}

	return nil
}

// StartWatchdog 启动看门狗
func StartWatchdog() error {
	go func() {
		CancelWatchdog()
		for {
			select {
			case <-context.Background().Done():
				return
			default:
			}
			FeedWatchdog()
			time.Sleep(5 * time.Second)
		}
	}()
	return nil
}
