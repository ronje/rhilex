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
# rc.local 文件路径
RC_LOCAL="/etc/rc.local"

# 显示帮助信息
show_help() {
    echo "Usage: $0 [install|uninstall|restart|status|help]"
    echo "  install    : 安装到 rc.local 以实现开机自启"
    echo "  uninstall  : 从 rc.local 中移除该程序的开机自启配置"
    echo "  restart    : 重启系统以应用更改"
    echo "  status     : 检查程序是否已配置在 rc.local 中开机自启"
    echo "  help       : 显示此帮助信息"
}

# 安装到 rc.local
install_script() {
    if [ ! -f "$RC_LOCAL" ]; then
        sudo touch "$RC_LOCAL"
        sudo chmod +x "$RC_LOCAL"
    fi
    if ! grep -q "$EXEC_COMMAND" "$RC_LOCAL"; then
        sudo sed -i "/^exit 0/i $EXEC_COMMAND" "$RC_LOCAL"
        echo "程序已成功安装到 rc.local 实现开机自启。"
    else
        echo "程序已经配置在 rc.local 中开机自启。"
    fi
}

# 从 rc.local 卸载
uninstall_script() {
    if grep -q "$EXEC_COMMAND" "$RC_LOCAL"; then
        sudo sed -i "/$EXEC_COMMAND/d" "$RC_LOCAL"
        echo "程序已从 rc.local 中移除开机自启配置。"
    else
        echo "程序未配置在 rc.local 中开机自启。"
    fi
}

# 重启系统
restart_system() {
    sudo reboot
}

# 检查状态
check_status() {
    if grep -q "$EXEC_COMMAND" "$RC_LOCAL"; then
        echo "程序已配置在 rc.local 中开机自启。"
    else
        echo "程序未配置在 rc.local 中开机自启。"
    fi
}

case "$1" in
    install)
        install_script
    ;;
    uninstall)
        uninstall_script
    ;;
    restart)
        restart_system
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