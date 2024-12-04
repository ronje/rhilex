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

// AllocateIPAddressUsingUdhcpc 使用 udhcpc 命令为指定的网络接口分配 IP 地址。
// interfaceName 是要分配 IP 地址的网络接口的名称。
// 超时时间设置为5秒，如果命令执行超时，返回错误。
func AllocateIPAddressUsingUdhcpc(interfaceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "udhcpc", "-i", interfaceName)
	log.Println("AllocateIPAddressUsingUdhcpc = ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("udhcpc command failed for interface %s: %w, output: %s", interfaceName, err, string(output))
	}
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("udhcpc command for interface %s timed out", interfaceName)
	}
	fmt.Printf("udhcpc output for interface %s: %s\n", interfaceName, string(output))
	return nil
}
