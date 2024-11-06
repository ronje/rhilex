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

package ossupport

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// StartNetworkInterface 启动指定的网络接口。
// interfaceName 是要启动的网络接口的名称，例如 "eth0" 或 "wlan0"。
// 如果启动成功，返回 nil；如果发生错误，返回错误。
func IpLinkSetIface(interfaceName, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ip", "link", "set", interfaceName, status)
	log.Println("IpLinkSetIface = ", cmd.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start network interface %s: %w", interfaceName, err)
	}
	return nil
}
