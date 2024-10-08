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
	"fmt"
	"os/exec"
	"time"
)

// nmcli device connect "MyWiFiNetwork" password "MySecurePassword" ifname "wlan0"
func SetWifi(iface, ssid, psk string, timeout time.Duration) error {
	s := "nmcli dev wifi connect \"%s\" password \"%s\" ifname \"%s\""
	{
		cmd := exec.Command("sh", "-c", fmt.Sprintf(s, ssid, ssid, iface))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("SetWifi failed, error:(%s), output:(%s)", err, string(out))
		}
	}
	{
		cmd := exec.Command("sh", "-c", `service networking restart`)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(err.Error() + ":" + string(output))
		}
	}
	return nil
}
