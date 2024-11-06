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
	"os"
	"os/exec"
	"time"
)

const __WINDOWS_SERVICENAME = "RhilexService"

// InstallService 安装服务
func InstallService() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sc.exe", "create", __WINDOWS_SERVICENAME, "binPath=", exePath)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error installing service: %v", err)
	}

	log.Println("Service installed. You can start it using 'sc start RhilexService'")
}

// UninstallService 卸载服务
func UninstallService() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sc.exe", "delete", __WINDOWS_SERVICENAME)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error uninstalling service: %v", err)
	}

	log.Println("Service uninstalled.")
}

// CheckStatus 查询服务状态
func CheckStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sc.exe", "query", __WINDOWS_SERVICENAME)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error checking service status: %v", err)
	}

	fmt.Println(string(output))
}

// RestartService 重启服务
func RestartService() {
	{
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		stopCmd := exec.CommandContext(ctx, "sc.exe", "stop", __WINDOWS_SERVICENAME)
		if err := stopCmd.Run(); err != nil {
			log.Fatalf("Error stopping service: %v", err)
		}

	}
	{
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		startCmd := exec.CommandContext(ctx, "sc.exe", "start", __WINDOWS_SERVICENAME)
		if err := startCmd.Run(); err != nil {
			log.Fatalf("Error starting service: %v", err)
		}

	}
	log.Println("Service restarted.")
}
