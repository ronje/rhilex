// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ossupport

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

/*
*
* NTP 用于启用NTP时间同步
*
 */

func UpdateTimeByNtp() error {
	err2 := setNtp(false)
	if err2 != nil {
		return err2
	}
	err1 := setNtp(true)
	if err1 != nil {
		return err1
	}
	return nil
}

/*
*
* 验证时间格式 YYYY-MM-DD HH:MM:SS
*
 */
func isValidTimeFormat(input string) bool {
	expectedFormat := "2006-01-02 15:04:05"
	_, err := time.Parse(expectedFormat, input)
	return err == nil
}

/*
*
* 获取当前系统时间
*
 */
func GetSystemTime() (string, error) {
	cmd := exec.Command("date", "+%Y-%m-%d %H:%M:%S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}

/*
*
*
设置时间，格式为 "YYYY-MM-DD HH:MM:SS"
*
*/
func SetSystemTime(newTime string) error {
	if !isValidTimeFormat(newTime) {
		return fmt.Errorf("Invalid time format:%s, must be 'YYYY-MM-DD HH:MM:SS'", newTime)
	}
	// newTime := "2023-08-10 15:30:00"
	cmd := exec.Command("date", "-s", newTime)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

/*
*
* v: true|false
*
 */
func setNtp(v bool) error {
	cmd := exec.Command("timedatectl", "set-ntp", func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}(v))
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(err.Error() + ":" + string(bytes))
	}
	return nil
}

/*
*
* 时区
*
 */
type TimeZoneInfo struct {
	CurrentTimezone string `json:"currentTimezone"`
	NTPSynchronized string `json:"NTPSynchronized"`
}

func GetTimeZone() (TimeZoneInfo, error) {
	timezoneInfo, err := GetLinuxTimeZone()
	if err != nil {
		return TimeZoneInfo{}, err
	}
	return TimeZoneInfo{CurrentTimezone: timezoneInfo}, nil
}

// GetLinuxTimeZone 返回当前 Linux 系统的时区
func GetLinuxTimeZone() (string, error) {
	// Linux 系统的时区通常存储在 /etc/timezone 文件中
	timezoneFilePath := "/etc/timezone"
	// 读取时区文件内容
	content, err := os.ReadFile(timezoneFilePath)
	if err != nil {
		// 如果 /etc/timezone 文件不存在，尝试读取 /etc/localtime 文件
		localtimeFilePath := "/etc/localtime"
		_, err := os.Stat(localtimeFilePath)
		if err != nil {
			return "", err
		}
		realPath, err := os.Readlink(localtimeFilePath)
		if err != nil {
			return "", err
		}
		zoneinfoPath := "/usr/share/zoneinfo/"
		if strings.HasPrefix(realPath, zoneinfoPath) {
			return realPath[len(zoneinfoPath):], nil
		}
		return "", fmt.Errorf("cannot determine timezone from %s", localtimeFilePath)
	}

	// 去除内容中的换行符
	timezone := strings.TrimSpace(string(content))
	return timezone, nil
}

// SetTimeZone 设置系统时区
// timezone := "Asia/Shanghai"
func SetTimeZone(timezone string) error {
	cmd := exec.Command("timedatectl", "set-timezone", timezone)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(err.Error() + ":" + string(output))
	}
	return nil
}

/*
*
* 获取开机时间
*
 */

func GetUptime() (string, error) {
	var info unix.Sysinfo_t

	if err := unix.Sysinfo(&info); err != nil {
		return "0 Year 0 Month 0 Days 0 Hours 0 Minutes 0 Seconds", err
	}

	return formatUptime(int64(info.Uptime)), nil
}

func formatUptime(uptime int64) string {
	days := uptime / 86400
	hours := (uptime % 86400) / 3600
	minutes := (uptime % 3600) / 60
	seconds := uptime % 60
	return fmt.Sprintf("%d days %d Hours %02d Minutes %02d Seconds", days, hours, minutes, seconds)
}
