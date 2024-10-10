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

package rhilexg1

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// isWirelessInterface checks if the given interface name corresponds to a wireless interface.
func isWirelessInterface(ifName string) bool {
	// On Linux, wireless interfaces typically have a directory under /sys/class/net/<iface>/wireless
	_, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/wireless", ifName))
	return !os.IsNotExist(err)
}

// getWlanList returns a list of wireless interfaces.
func getWlanList() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var wlanIfaces []net.Interface
	for _, iface := range ifaces {
		if isWirelessInterface(iface.Name) {
			wlanIfaces = append(wlanIfaces, iface)
		}
	}

	return wlanIfaces, nil
}

func SetWifi(iface, ssid, psk string, timeout time.Duration) error {
	if len(psk) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	// nmcli dev wifi connect SSID password pwd
	if WifiAlreadyConfig(ssid) {
		s := "nmcli connection up %s"
		cmd := exec.Command("sh", "-c", fmt.Sprintf(s, ssid))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Command error: %s,  %s", err, string(out))
		}
	} else {
		s := "nmcli dev wifi connect \"%s\" password \"%s\""
		cmd := exec.Command("sh", "-c", fmt.Sprintf(s, ssid, psk))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Command error: %s,  %s", err, string(out))
		}
	}
	return nil
}
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
* 升级版，带上了WIFI信号强度
*
 */
func ScanWlanList(WFace string) ([][2]string, error) {
	wifiList := [][2]string{}
	shell := `
iw dev %s scan | awk '
  /SSID/ { ssid=$2 } \
  /signal/ { signal=$2; if (!seen[ssid] || signal > seen[ssid]) { seen[ssid]=signal } } \
  END { for (s in seen) print s "," seen[s]}
' | sort
`
	cmd := exec.Command("sh", "-c", fmt.Sprintf(shell, WFace))
	fmt.Println("ScanWlanList == ", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	lines := bufio.NewScanner(strings.NewReader(string(output)))
	for lines.Scan() {
		line := lines.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			ssid := parts[0]
			signal := parts[1]
			if ssid != "" {
				wifiList = append(wifiList, [2]string{ssid, signal})
			}
		}
	}
	return wifiList, nil
}
