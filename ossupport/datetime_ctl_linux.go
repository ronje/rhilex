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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "date", "+%Y-%m-%d %H:%M:%S")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "date", "-s", newTime)
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
	for _, server := range ntpServers {
		ntpTime, err = queryNTPTime(server)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("all NTP servers failed to respond")
	}
	ntpTimeStr := ntpTime.Format("2006-01-02 15:04:05")
	cmd := exec.Command("date", "-s", ntpTimeStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set system time: %v", err)
	}
	log.Printf("System time set to: %s\n", ntpTimeStr)
	log.Printf("Command output: %s\n", output)
	return nil
}

// queryNTPTime 尝试从指定的NTP服务器获取时间
func queryNTPTime(server string) (time.Time, error) {
	response, err := ntp.Query(server)
	if err != nil {
		return time.Time{}, err
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

func GetTimeZone() (TimeZoneInfo, error) {
	timezoneInfo, err := GetLinuxTimeZone()
	if err != nil {
		return TimeZoneInfo{}, err
	}
	return TimeZoneInfo{CurrentTimezone: timezoneInfo}, nil
}

func GetLinuxTimeZone() (string, error) {
	timezoneFilePath := "/etc/timezone"
	content, err := os.ReadFile(timezoneFilePath)
	if err != nil {
		return "Asia/Shanghai", nil
	}
	timezone := strings.TrimSpace(string(content))
	return timezone, nil
}

// SetTimeZone 设置系统时区
// timezone := "Asia/Shanghai"
func SetTimeZone(timezone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "timedatectl", "set-timezone", timezone)
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
