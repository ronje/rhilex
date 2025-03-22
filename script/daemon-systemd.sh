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

# 可执行文件路径和执行指令
EXEC_FILE="/usr/local/rhilex/rhilex"
EXEC_COMMAND="/usr/local/rhilex/rhilex run -config= /usr/local/rhilex/rhilex.ini"
# systemd 服务文件名称
SERVICE_NAME="rhilex.service"
# systemd 服务文件路径
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME"

# 显示帮助信息
show_help() {
    echo "Usage: $0 [install|uninstall|restart|status|help]"
    echo "  install    : 安装为 systemd 服务以实现开机自启"
    echo "  uninstall  : 卸载 systemd 服务"
    echo "  restart    : 重启 systemd 服务"
    echo "  status     : 查看 systemd 服务的状态"
    echo "  help       : 显示此帮助信息"
}

# 安装为 systemd 服务
install_service() {
    cat << EOF | sudo tee "$SERVICE_FILE"
[Unit]
Description=Rhilex Service
After=network.target

[Service]
ExecStart=$EXEC_COMMAND
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE_NAME"
    sudo systemctl start "$SERVICE_NAME"
    echo "程序已成功安装为 systemd 服务并启动。"
}

# 卸载 systemd 服务
uninstall_service() {
    sudo systemctl stop "$SERVICE_NAME"
    sudo systemctl disable "$SERVICE_NAME"
    sudo rm -f "$SERVICE_FILE"
    sudo systemctl daemon-reload
    echo "systemd 服务已卸载。"
}

# 重启 systemd 服务
restart_service() {
    sudo systemctl restart "$SERVICE_NAME"
    echo "systemd 服务已重启。"
}

# 查看 systemd 服务状态
check_status() {
    sudo systemctl status "$SERVICE_NAME"
}

case "$1" in
    install)
        install_service
    ;;
    uninstall)
        uninstall_service
    ;;
    restart)
        restart_service
    ;;
    status)
        check_status
    ;;
    help)
        show_help
    ;;
    *)
        show_help
        exit 1
    ;;
esac