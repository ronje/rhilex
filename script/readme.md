<!--
 Copyright (C) 2024 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

# Rhilex 守护进程
## 概述
这个脚本是一个用于管理 "rhilex" 服务的守护进程管理器。它提供了安装、启动、停止、重启和卸载服务以及检查服务状态的功能。脚本包含一个日志记录函数来记录事件，并使用信号来停止和升级服务。
## 安装
要安装 rhilex 服务，请使用 `install` 命令运行脚本：
```bash
sudo ./rhilex-daemon.sh install
```

这将创建必要的服务文件，并将可执行文件和配置文件复制到相应的目录。
## 配置
rhilex 服务的配置存储在 `rhilex.ini` 文件中，默认位于 `/usr/local/rhilex` 目录。您可以修改这个文件来调整服务的设置。
## 使用方法
### 启动服务
要启动 rhilex 服务，请使用 `start` 命令：
```bash
sudo ./rhilex-daemon.sh start
```

### 停止服务
要停止 rhilex 服务，请使用 `stop` 命令：
```bash
sudo ./rhilex-daemon.sh stop
```

### 重启服务
要重启 rhilex 服务，请使用 `restart` 命令：
```bash
sudo ./rhilex-daemon.sh restart
```

### 检查服务状态
要检查 rhilex 服务的状态，请使用 `status` 命令：
```bash
sudo ./rhilex-daemon.sh status
```

### 卸载服务
要卸载 rhilex 服务，请使用 `uninstall` 命令：
```bash
sudo ./rhilex-daemon.sh uninstall
```

这将删除服务文件、可执行文件和配置文件。
