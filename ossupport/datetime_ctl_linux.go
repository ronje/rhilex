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
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"golang.org/x/sys/unix"
)

/*
*
* NTP 用于启用NTP时间同步
*
 */

// UpdateTimeByNtp 通过NTP服务器更新系统时间
func UpdateTimeByNtp() error {
	// 先关闭NTP（这里实际未做操作，只是调用函数）
	if err := setNtp(false); err != nil {
		return err
	}
	// 再开启NTP
	if err := setNtp(true); err != nil {
		return err
	}
	return nil
}

/*
*
* 验证时间格式 YYYY-MM-DD HH:MM:SS
*
 */
// isValidTimeFormat 验证输入的时间字符串是否符合 "YYYY-MM-DD HH:MM:SS" 格式
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
// GetSystemTime 获取当前系统时间
func GetSystemTime() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行 date 命令获取系统时间
	cmd := exec.CommandContext(ctx, "date", "+%Y-%m-%d %H:%M:%S")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get system time: %w", err)
	}

	// 去除输出中的换行符
	return strings.TrimSpace(string(output)), nil
}

/*
*
* 设置时间，格式为 "YYYY-MM-DD HH:MM:SS"
*
 */
// SetSystemTime 设置系统时间，需要输入符合 "YYYY-MM-DD HH:MM:SS" 格式的时间字符串
func SetSystemTime(newTime string) error {
	if !isValidTimeFormat(newTime) {
		return fmt.Errorf("invalid time format: %s, must be 'YYYY-MM-DD HH:MM:SS'", newTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行 date 命令设置系统时间
	cmd := exec.CommandContext(ctx, "date", "-s", newTime)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set system time: %w", err)
	}

	return nil
}

/*
*
* v: true|false
*
 */

// setNtp 尝试从提供的NTP服务器列表中获取时间，并设置本地系统时间。
// 如果成功获取时间并设置系统时间，则返回nil；否则返回错误。
func setNtp(v bool) error {
	if !v {
		return nil
	}

	ntpServers := []string{
		"ntp.sjtu.edu.cn",
		"ntp.neu.edu.cn",
		"ntp.bupt.edu.cn",
		"ntp.shu.edu.cn",
		"ntp.tuna.tsinghua.edu.cn",
		"ntp1.aliyun.com",
		"ntp2.aliyun.com",
		"ntp3.aliyun.com",
		"ntp4.aliyun.com",
		"ntp5.aliyun.com",
		"ntp6.aliyun.com",
		"ntp7.aliyun.com",
		"0.cn.pool.ntp.org",
		"1.cn.pool.ntp.org",
		"2.cn.pool.ntp.org",
		"3.cn.pool.ntp.org",
		"time1.cloud.tencent.com",
		"time2.cloud.tencent.com",
		"time3.cloud.tencent.com",
		"time4.cloud.tencent.com",
		"time5.cloud.tencent.com",
	}

	var ntpTime time.Time
	var err error

	// 遍历NTP服务器列表，尝试获取时间
	for _, server := range ntpServers {
		ntpTime, err = queryNTPTime(server)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("all NTP servers failed to respond: %w", err)
	}

	ntpTimeStr := ntpTime.Format("2006-01-02 15:04:05")

	// 执行 date 命令设置系统时间
	cmd := exec.Command("date", "-s", ntpTimeStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set system time: %w, output: %s", err, string(output))
	}

	log.Printf("System time set to: %s\n", ntpTimeStr)
	log.Printf("Command output: %s\n", string(output))

	return nil
}

// queryNTPTime 尝试从指定的NTP服务器获取时间
func queryNTPTime(server string) (time.Time, error) {
	// 向NTP服务器查询时间
	response, err := ntp.Query(server)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to query NTP server %s: %w", server, err)
	}
	return response.Time, nil
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

// GetTimeZone 获取当前系统的时区信息
func GetTimeZone() (TimeZoneInfo, error) {
	timezone, err := GetLinuxTimeZone()
	if err != nil {
		return TimeZoneInfo{}, err
	}
	return TimeZoneInfo{CurrentTimezone: timezone}, nil
}

// GetLinuxTimeZone 读取 /etc/timezone 文件获取Linux系统的时区
func GetLinuxTimeZone() (string, error) {
	timezoneFilePath := "/etc/timezone"
	content, err := os.ReadFile(timezoneFilePath)
	if err != nil {
		// 读取失败时默认返回 "Asia/Shanghai"
		return "Asia/Shanghai", nil
	}
	// 去除前后空格
	return strings.TrimSpace(string(content)), nil
}

// SetTimeZone 设置系统时区
func SetTimeZone(timezone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行 timedatectl 命令设置时区
	cmd := exec.CommandContext(ctx, "timedatectl", "set-timezone", timezone)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set timezone: %w, output: %s", err, string(output))
	}

	return nil
}

/*
*
* 获取开机时间
*
 */

// GetUptime 获取系统的开机时间
func GetUptime() (string, error) {
	var info unix.Sysinfo_t

	if err := unix.Sysinfo(&info); err != nil {
		return "0 Year 0 Month 0 Days 0 Hours 0 Minutes 0 Seconds", fmt.Errorf("failed to get system uptime: %w", err)
	}

	return formatUptime(int64(info.Uptime)), nil
}

// formatUptime 将开机时间的秒数格式化为人类可读的时间格式
func formatUptime(uptime int64) string {
	days := uptime / 86400
	hours := (uptime % 86400) / 3600
	minutes := (uptime % 3600) / 60
	seconds := uptime % 60
	return fmt.Sprintf("%d days %d Hours %02d Minutes %02d Seconds", days, hours, minutes, seconds)
}
