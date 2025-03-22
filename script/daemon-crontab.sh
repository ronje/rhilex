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

# 显示帮助信息
show_help() {
    echo "Usage: $0 [install|uninstall|restart|status|help]"
    echo "  install    : 安装到 crontab 以实现开机自启"
    echo "  uninstall  : 从 crontab 中移除该程序的开机自启配置"
    echo "  restart    : 此功能对 crontab 不适用，请手动重启系统"
    echo "  status     : 检查程序是否已配置在 crontab 中开机自启"
    echo "  help       : 显示此帮助信息"
}

# 安装到 crontab
install_script() {
    TEMP_CRON=$(mktemp)
    crontab -l > "$TEMP_CRON" 2>/dev/null
    if ! grep -q "@reboot $EXEC_COMMAND" "$TEMP_CRON"; then
        echo "@reboot $EXEC_COMMAND" >> "$TEMP_CRON"
        crontab "$TEMP_CRON"
        echo "程序已成功安装到 crontab 实现开机自启。"
    else
        echo "程序已经配置在 crontab 中开机自启。"
    fi
    rm -f "$TEMP_CRON"
}

# 从 crontab 卸载
uninstall_script() {
    TEMP_CRON=$(mktemp)
    crontab -l > "$TEMP_CRON" 2>/dev/null
    if grep -q "@reboot $EXEC_COMMAND" "$TEMP_CRON"; then
        sed -i "/@reboot $EXEC_COMMAND/d" "$TEMP_CRON"
        crontab "$TEMP_CRON"
        echo "程序已从 crontab 中移除开机自启配置。"
    else
        echo "程序未配置在 crontab 中开机自启。"
    fi
    rm -f "$TEMP_CRON"
}

# 重启系统（crontab 场景下提示手动重启）
restart_system() {
    echo "此功能对 crontab 不适用，请手动重启系统。"
}

# 检查状态
check_status() {
    TEMP_CRON=$(mktemp)
    crontab -l > "$TEMP_CRON" 2>/dev/null
    if grep -q "@reboot $EXEC_COMMAND" "$TEMP_CRON"; then
        echo "程序已配置在 crontab 中开机自启。"
    else
        echo "程序未配置在 crontab 中开机自启。"
    fi
    rm -f "$TEMP_CRON"
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