<!--
 Copyright (C) 2025 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# `rhilex-daemon-systemctl.sh` 指南

## 1. 脚本概述
`rhilex-daemon-systemctl.sh` 是一个用于管理 `Rhilex` 服务的脚本，它可以帮助你方便地安装、卸载、重启 `Rhilex` 服务，获取服务状态，还提供了帮助信息。该脚本基于 `systemd` 系统服务管理，因此只适用于支持 `systemd` 的系统。

## 2. 脚本的准备
### 2.1 下载脚本
将上述脚本内容保存为一个文件，例如 `rhilex-daemon-systemctl.sh`。可以使用以下命令创建并编辑该文件：
```bash
nano rhilex-daemon-systemctl.sh
```
将脚本内容复制到编辑器中，然后保存并退出（在 `nano` 编辑器中，按 `Ctrl + X`，然后按 `Y` 确认保存，最后按 `Enter` 退出）。

### 2.2 添加执行权限
给脚本添加执行权限，以便可以直接运行该脚本：
```bash
chmod +x rhilex-daemon-systemctl.sh
```

## 3. 检查系统是否支持 `systemd`
在运行脚本之前，需要确保你的系统支持 `systemd` 服务管理。脚本在运行时会自动检查，如果系统不支持 `systemd`，脚本将输出错误信息并退出。你也可以手动检查 `systemctl` 命令是否存在：
```bash
command -v systemctl
```
如果输出 `systemctl` 的路径，则表示系统支持 `systemd`。

## 4. 脚本的使用
### 4.1 查看帮助信息
如果你不确定如何使用该脚本，可以使用 `help` 选项来查看帮助信息：
```bash
./rhilex-daemon-systemctl.sh help
```
脚本将输出使用方法和各个选项的说明，例如：
```sh
Usage: ./rhilex-daemon-systemctl.sh [install|uninstall|restart|status|help]
Options:
  install: Install and start the Rhilex service.
  uninstall: Stop and uninstall the Rhilex service.
  restart: Restart the Rhilex service.
  status: Get the status of the Rhilex service.
  help: Show this help message.
```

### 4.2 安装服务
使用 `install` 选项来安装并启动 `Rhilex` 服务：
```bash
./rhilex-daemon-systemctl.sh install
```
脚本将创建服务文件，重新加载 `systemd` 管理器的配置，启用并启动服务。如果安装成功，将输出相应的成功信息。

### 4.3 卸载服务
使用 `uninstall` 选项来停止并卸载 `Rhilex` 服务：
```bash
./rhilex-daemon-systemctl.sh uninstall
```
脚本将停止服务，禁用服务，删除服务文件，并重新加载 `systemd` 管理器的配置。如果卸载成功，将输出相应的成功信息。

### 4.4 重启服务
使用 `restart` 选项来重启 `Rhilex` 服务：
```bash
./rhilex-daemon-systemctl.sh restart
```
脚本将重启服务，并输出重启成功的信息。

### 4.5 获取服务状态
使用 `status` 选项来获取 `Rhilex` 服务的状态信息：
```bash
./rhilex-daemon-systemctl.sh status
```
脚本将调用 `systemctl status` 命令，输出服务的详细状态信息，包括服务是否正在运行、最后一次启动时间等。

## 5. 注意事项
- 脚本需要以具有足够权限的用户（如 `root` 用户或使用 `sudo`）来运行，因为涉及到系统服务的管理和文件的读写操作。
- 确保 `/usr/local/rhilex/` 目录存在，并且 `/usr/local/rhilex/rhilex` 可执行文件存在且具有执行权限，否则服务可能无法正常启动。
- 如果在使用过程中遇到问题，可以先查看脚本输出的错误信息，或者使用 `help` 选项查看帮助信息。如果问题仍然存在，可以检查系统日志或寻求相关技术支持。