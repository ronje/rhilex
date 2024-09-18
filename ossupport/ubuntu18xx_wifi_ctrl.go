package ossupport

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/glogger"
)

type WLANInterface struct {
	Interface string `yaml:"-" json:"interface"`
	SSID      string `yaml:"-" json:"ssid"`
	Password  string `yaml:"-" json:"password"`
	Security  string `yaml:"-" json:"security"`
}

type WlanConfig struct {
	Wlan0 WLANInterface `yaml:"-" json:"wlan0"`
}

/*
*
* WIFI 控制
nmcli device wifi rescan
nmcli device wifi list

获取WIFI网卡: iw dev | awk '$1=="Interface"{print $2}'
扫描WIFI列表: iwlist wlx0cc6551c5026 scan | grep ESSID | awk -F: '{print $2}' | sed 's/"//g'
*
*/
func ScanWIFIWithNmcli() ([]string, error) {

	wifiListReturn := []string{}
	var errReturn error
	finished := make(chan bool)

	go func() {
		{
			// 第一遍先扫描手边的 WIFI SSID
			cmd := exec.Command("sh", "-c", "nmcli device wifi rescan")
			output, err := cmd.CombinedOutput()
			if err != nil {
				errReturn = fmt.Errorf("scan WIFI error:%s,%s", string(output), err)
				return
			}
			glogger.GLogger.Debug("ScanWIFIWithNmcli:", cmd.String(), " OutPut:", string(output))

		}
		WFace := ""
		{
			cmd := exec.Command("sh", "-c", `iw dev | awk '$1=="Interface"{print $2}'`)
			output, err := cmd.CombinedOutput()
			if err != nil {
				stringWithoutNewlines := strings.Replace(string(output), "\n", "", -1)
				errReturn = fmt.Errorf("get WLAN Interface error:%s,%s", stringWithoutNewlines, err)
				return
			}
			if len(output) > 0 {
				WFace = string(output)
			} else {
				errReturn = fmt.Errorf("get WLAN Interface error:%s,%s", string(output), err)
				return
			}
			glogger.GLogger.Debug("ScanWIFIWithNmcli:", cmd.String(), " OutPut:", WFace)

		}

		{
			shell := `iwlist %s scan | grep ESSID | awk -F: '{print $2}' | sed 's/"//g'`
			stringWithoutNewlines := strings.Replace(WFace, "\n", "", -1)
			cmd := exec.Command("sh", "-c", fmt.Sprintf(shell, stringWithoutNewlines))
			output, err := cmd.CombinedOutput()
			if err != nil {
				errReturn = fmt.Errorf("scan WIFI error:%s,%s", string(output), err)
				return
			}
			for _, v := range strings.Split(string(output), "\n") {
				// AAA\nBBB\nCCC\n
				if v != "" {
					wifiListReturn = append(wifiListReturn, v)
				}
			}
			glogger.GLogger.Debug("Scan WIFI With Nmcli:", cmd.String(), " OutPut :", wifiListReturn)
		}
		finished <- true
	}()
	select {
	case <-time.After(10 * time.Second): // 超时时间6秒
		errReturn = fmt.Errorf("scan WIFI timeout")
		return wifiListReturn, errReturn
	case <-finished:
		return wifiListReturn, errReturn
	}
}

/*
*
  - 初始化
    // 删除之前的连接
    // if exists ${name} -> nmcli connection delete ${name}
    // 重新连接
    // sudo nmcli dev wifi connect "ssid" password "password"
*/
func WifiAlreadyConfig(wifiSSIDName string) bool {
	connectionsDir := "/etc/NetworkManager/system-connections/"
	files, err := os.ReadDir(connectionsDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return false
	}
	for _, file := range files {
		if wifiSSIDName == file.Name() {
			return true
		}
	}
	return false
}

/*
 *
 * 扫描
 *
 */
func ScanWlanList() error {
	cmd := exec.Command("nmcli", "device", "wifi", "rescan")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	return nil
}

/*
*
* 获取WIFI列表
*
 */
func GetWlanList() ([]string, error) {
	// 执行 nmcli 命令来获取WIFI列表
	cmd := exec.Command("nmcli", "--fields", "SSID,MODE,FREQ,SIGNAL,BARS,SECURITY", "device", "wifi", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	lines := strings.Split(string(output), "\n")
	var wifiList []string
	wifiList = append(wifiList, lines...)
	return wifiList, nil
}

/*
*
* 升级版，带上了WIFI信号强度
*
 */
func GetWlanListSignal() ([][2]string, error) {
	wifiList := [][2]string{}
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,signal", "device", "wifi", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	lines := bufio.NewScanner(strings.NewReader(string(output)))
	for lines.Scan() {
		line := lines.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			ssid := parts[0]
			signal := parts[1]
			wifiList = append(wifiList, [2]string{ssid, signal})
		}
	}
	return wifiList, nil
}

/*
*
* 关闭WIFI开关
*
 */
func DisableWifi() error {
	cmd := exec.Command("nmcli", "radio", "wifi", "off")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	return nil
}

/*
*
* 打开WIFI开关
*
 */
func EnableWifi() error {
	cmd := exec.Command("nmcli", "radio", "wifi", "on")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	return nil
}
