# Copyright (C) 2025 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

#!/bin/bash

# 检查系统是否支持 systemd
check_systemd_support() {
    if ! command -v systemctl >/dev/null 2>&1; then
        echo "This system does not support systemd services."
        exit 1
    fi
}

# 服务名称
SERVICE_NAME="rhilex.service"
# 服务文件内容
SERVICE_CONTENT="[Unit]
Description=Rhilex Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/usr/local/rhilex/
ExecStart=/usr/local/rhilex/rhilex run
Environment=\"PATH=/usr/local/rhilex/bin:/usr/local/bin:/usr/bin:/bin\"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target"

# 安装服务的函数
install_service() {
    echo "Creating service file..."
    echo "$SERVICE_CONTENT" | sudo tee /etc/systemd/system/$SERVICE_NAME > /dev/null
    sudo systemctl daemon-reload
    sudo systemctl enable $SERVICE_NAME
    sudo systemctl start $SERVICE_NAME
    echo "Service installed and started successfully."
}

# 卸载服务的函数
uninstall_service() {
    echo "Stopping service..."
    sudo systemctl stop $SERVICE_NAME
    echo "Disabling service..."
    sudo systemctl disable $SERVICE_NAME
    echo "Removing service file..."
    sudo rm /etc/systemd/system/$SERVICE_NAME
    sudo systemctl daemon-reload
    echo "Service uninstalled successfully."
}

# 重启服务的函数
restart_service() {
    echo "Restarting service..."
    sudo systemctl restart $SERVICE_NAME
    echo "Service restarted successfully."
}

# 获取服务状态的函数
get_service_status() {
    echo "Getting service status..."
    sudo systemctl status $SERVICE_NAME
}

# 显示帮助信息的函数
show_help() {
    echo "Usage: $0 [install|uninstall|restart|status|help]"
    echo "Options:"
    echo "  install: Install and start the Rhilex service."
    echo "  uninstall: Stop and uninstall the Rhilex service."
    echo "  restart: Restart the Rhilex service."
    echo "  status: Get the status of the Rhilex service."
    echo "  help: Show this help message."
}

# 检查系统是否支持 systemd
check_systemd_support

# 检查脚本参数
if [ "$1" == "install" ]; then
    install_service
    elif [ "$1" == "uninstall" ]; then
    uninstall_service
    elif [ "$1" == "restart" ]; then
    restart_service
    elif [ "$1" == "status" ]; then
    get_service_status
    elif [ "$1" == "help" ]; then
    show_help
else
    echo "Invalid option. Use '$0 help' to see the available options."
    exit 1
fi