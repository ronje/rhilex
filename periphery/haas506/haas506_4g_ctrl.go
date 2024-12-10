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

package haas506

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/ossupport"
)

const _ML307R_4G_PATH = "/dev/ttyUSB2"

type ModuleInfo struct {
	Up    bool
	ICCID string
	IMEL  string
	CSQ   int32
	COPS  string
}

/*
*
* 初始化4G模组
*
 */
func InitML307R4G(path string) {
	if err := turnOffEcho(path); err != nil {
		log.Println("ML307RInit4G turnOffEcho error:", err)
		return
	}
	log.Println("ML307RInit4G resetCard ok.")

}

// /dev/ttyUSB2

/*
*
* 获取4G数据
*
 */
func Get4GBaseInfo() ModuleInfo {
	info := ModuleInfo{
		ICCID: "0",
		CSQ:   0,
		COPS:  "UNKNOWN",
	}

	csq := ML307RGet4G_CSQ(_ML307R_4G_PATH)
	if csq == 0 {
		time.Sleep(100 * time.Millisecond)
		csq = ML307RGet4G_CSQ(_ML307R_4G_PATH)
	}
	cops, err1 := ML307RGetCOPS(_ML307R_4G_PATH)
	if err1 != nil {
		return info
	}
	cm := "UNKNOWN"
	if strings.Contains(cops, "CMCC") {
		cm = "CHINA CMCC"
	}
	if strings.Contains(cops, "MOBILE") {
		cm = "CHINA MOBILE"
	}
	if strings.Contains(cops, "UNICOM") {
		cm = "CHINA UNICOM"
	}
	iccid, err2 := ML307RGetICCID(_ML307R_4G_PATH)
	if err2 != nil {
		return info
	}
	imel, err3 := ML307RGetIMEL(_ML307R_4G_PATH)
	if err3 != nil {
		return info
	}
	Up, _ := ossupport.IsInterfaceUp("eth0")
	info.Up = Up
	info.IMEL = imel
	info.COPS = cm
	info.CSQ = csq
	info.ICCID = iccid
	return info
}

/*
*
  - 初始化4G模组
    echo -e "AT+QCFG=\"usbnet\",1\r\n" >/dev/ttyUSB2  //驱动模式
    echo -e "AT+QNETDEVCTL=3,1,1\r\n" >/dev/ttyUSB2   //自动拨号
    echo -e "AT+QCFG=\"nat\",1 \r\n" >/dev/ttyUSB2    //网卡模式
    echo -e "AT+CFUN=1,1\r\n" >/dev/ttyUSB2           //重启
*/

const (
	__AT_TIMEOUT = 300 * time.Millisecond // timeout ms
)

const (
	__TURN_OFF_ECHO = "ATE0\r\n"           // ECHO Mode
	__GET_CSQ_CMD   = "AT+CSQ\r\n"         // CSQ  信号
	__GET_COPS_CMD  = "AT+COPS?\r\n"       // COPS 运营商
	__DAIL_CMD      = "AT+MDIALUP=1,1\r\n" // 拨号
	__GET_IMEL_CMD  = "AT+CGSN=1\r\n"      // Get IMEL
	__GET_INFO_CMD  = "ATI\r\n"            // Get INFO
	__SAVE_CONFIG   = "AT&W\r\n"           // SaveConfig
)

/**
 * 开启4G
 *
 */
func ML307RTurnOn4G() error {
	__ML307R_AT(_ML307R_4G_PATH, __DAIL_CMD, __AT_TIMEOUT)
	{
		ctx, Cancel := context.WithDeadline(context.Background(), time.Now().Add(3000*time.Millisecond))
		defer Cancel()
		output, err := exec.CommandContext(ctx, "sh", "-c", `ifconfig eth0 up`).CombinedOutput()
		log.Println("[ML307RTurnOn4G] ifconfig eth0 up:", ", Output=", string(output))
		if err != nil {
			return err
		}
	}
	{
		ctx, Cancel := context.WithDeadline(context.Background(), time.Now().Add(3000*time.Millisecond))
		defer Cancel()
		output, err := exec.CommandContext(ctx, "sh", "-c", `udhcpc -i eth0`).CombinedOutput()
		log.Println("[ML307RTurnOn4G] udhcpc -i eth0:", ", Output=", string(output))
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * 关闭4g
 *
 */
func ML307RTurnOff4G() error {
	ctx, Cancel := context.WithDeadline(context.Background(), time.Now().Add(3000*time.Millisecond))
	defer Cancel()
	output, err := exec.CommandContext(ctx, "sh", "-c", `ifconfig eth0 down`).CombinedOutput()
	log.Println("[ML307RTurnOff4G] ifconfig eth0 down:", ", Output=", string(output))
	if err != nil {
		return err
	}
	return nil
}

/*
*
* APN 配置
 */
func ML307RGetAPN(path string) (string, error) {
	return "", nil
}

// 场景恒等于1
func ML307RSetAPN(path string, ptype int, apn, username, password string, auth, cdmaPwd int) (string, error) {
	return "", nil
}

/*
*
  - 获取信号: +CSQ: 39,99
  - 0：没有信号。
  - 1-9：非常弱的信号，可能无法建立连接。
  - 10-14：较弱的信号，但可能可以建立连接。
  - 15-19：中等强度的信号。
  - 20-31：非常强的信号，信号质量非常好。
    ML307RGet4G_CSQ: 返回值代表信号格
*/
func ML307RGet4G_CSQ(path string) int32 {
	return __Get4G_CSQ(path)
}

func ML307RGetINFO(path string) (string, error) {
	return __ML307R_AT(path, __GET_INFO_CMD, __AT_TIMEOUT)
}

/*
*
* 获取运营商
+COPS:
(2,"CHINA MOBILE","CMCC","46000",7),
(1,"CHINA MOBILE","CMCC","46000",0),
(3,"CHN-UNICOM","UNICOM","46001",7),
(3,"CHN-CT","CT","46011",7),
(1,"460 15","460 15","46015",7),,(0-4),(0-2)
*/
func ML307RGetCOPS(path string) (string, error) {
	// +COPS: 0,0,"CHINA MOBILE",7
	// +COPS: 0,0,"CHIN-UNICOM",7
	return __ML307R_AT(path, __GET_COPS_CMD, __AT_TIMEOUT)
}

/*
*
* 获取ICCID, 用户查询电话卡号
* +QCCID: 89860025128306012474
 */
func ML307RGetICCID(path string) (string, error) {
	return "UNKNOWN", nil

}

/**
 * 获取IMEL
 *
 */
func ML307RGetIMEL(path string) (string, error) {
	return __ML307R_AT(path, __GET_IMEL_CMD, __AT_TIMEOUT)

}
func __Get4G_CSQ(path string) int32 {
	csq := int32(0)
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return csq
	}
	defer file.Close()
	_, err = file.WriteString(__GET_CSQ_CMD)
	if err != nil {
		return csq
	}
	// 4G 模组的绝大多数指令都是100毫秒
	timeout := 200 * time.Millisecond
	buffer := [1]byte{}
	var responseData []byte
	b1 := 0
	for {
		if b1 == 4 {
			break
		}
		deadline := time.Now().Add(timeout)
		file.SetReadDeadline(deadline)
		n, err := file.Read(buffer[:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				break
			}
		}
		if n > 0 {
			if buffer[0] == 10 {
				b1++
			}
			if buffer[0] != 10 {
				responseData = append(responseData, buffer[0])
			}
		}
	}
	log.Println("[ML307R __Get4G_CSQ]:", __GET_CSQ_CMD, ", Output:", string(responseData))
	if len(responseData) > 6 {
		// +CSQ: 30,99
		response := string(responseData[6:])
		parts := strings.Split(response, ",")
		if len(parts) == 2 {
			v, err := strconv.Atoi(parts[0])
			if err == nil {
				csq = int32(v)
			}
		}
	}

	return csq
}

func turnOffEcho(path string) error {
	return __ExecuteAT(path, __TURN_OFF_ECHO)
}

func __ExecuteAT(path string, cmd string) error {
	buffer, err0 := __ML307R_AT(path, cmd, __AT_TIMEOUT)
	log.Println("[ML307R __ExecuteAT]:", cmd, ", Output:", (buffer))
	if err0 != nil {
		return err0
	}
	buffer, err1 := __ML307R_AT(path, __SAVE_CONFIG, __AT_TIMEOUT)
	log.Println("[ML307R __ExecuteAT]:", cmd, ", Output:", (buffer))
	if err1 != nil {
		return err1
	}
	return nil
}

/*
*
  - ML307R 系列AT指令封装
    指令格式：AT+<CMD>\r\n
    指令返回值：\r\nCMD\r\n\r\nOK\r\n

解析结果

	[

		"",
		"CMD",
		"",
		"OK",
		""

	]

*
*/
func __ML307R_AT(path string, command string, timeout time.Duration) (string, error) {
	// 打开设备文件以供读写
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 写入AT指令
	_, err = file.WriteString(command)
	if err != nil {
		return "", err
	}
	buffer := [1]byte{}
	var responseData []byte
	b1 := 0
	for {
		if b1 == 4 {
			break
		}
		deadline := time.Now().Add(timeout)
		file.SetReadDeadline(deadline)
		n, err := file.Read(buffer[:])
		if err != nil {
			return "", err
		}
		if n > 0 {
			if buffer[0] == 10 {
				b1++
			}
			if buffer[0] != 10 {
				responseData = append(responseData, buffer[0])
			}
		}
	}
	log.Println("[ML307R __ML307R_AT]:", command, ", Output:", string(responseData))
	return string(responseData), nil
}
