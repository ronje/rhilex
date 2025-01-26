# Copyright (C) 2024 wwhai
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

# 服务相关的名称和路径定义
SERVICE_NAME="rhilex"
WORKING_DIRECTORY="/usr/local/rhilex"
EXECUTABLE_PATH="$WORKING_DIRECTORY/$SERVICE_NAME"
CONFIG_PATH="$WORKING_DIRECTORY/$SERVICE_NAME.ini"

SERVICE_FILE="/etc/init.d/$SERVICE_NAME"

STOP_SIGNAL="/var/run/rhilex-stop.signal"
UPGRADE_SIGNAL="/var/run/rhilex-upgrade.lock"

SOURCE_DIR="$PWD"

# 日志函数，用于记录不同级别的日志信息
log() {
    local level=$1
    shift
    echo "[$level] $(date +'%Y-%m-%d %H:%M:%S') - $@"
}

# 安装服务的函数
install() {
    cat > "$SERVICE_FILE" << EOL
#!/bin/sh

export PATH=\$PATH:/usr/local/bin:/usr/bin:/usr/sbin:/usr/local/sbin:/sbin

### BEGIN INIT INFO
# Provides:          rhilex
# Required-Start:    \$network \$local_fs \$remote_fs
# Required-Stop:     \$network \$local_fs \$remote_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: rhilex Service
# Description:       rhilex Service
### END INIT INFO

#
# Create Time: $(date +'%Y-%m-%d %H:%M:%S')
#

start() {
    rm -f $UPGRADE_SIGNAL
    rm -f $STOP_SIGNAL
    pid=\$(pgrep -x -n -f "$EXECUTABLE_PATH run -config=$CONFIG_PATH")
    if [ -n "\$pid" ]; then
        log INFO "rhilex is running with Pid:\${pid}"
        exit 0
    fi
    cd $WORKING_DIRECTORY
    daemon &
    log INFO "rhilex daemon started."
    exit 0
}

stop() {
    echo "1" > $STOP_SIGNAL
    if pgrep -x "rhilex" > /dev/null; then
        log INFO "rhilex process is running. Killing it..."
        pkill -x "rhilex"
        log INFO "rhilex process has been killed."
    else
        log WARNING "rhilex process is not running."
    fi
}

restart() {
    stop
    sleep 1
    start
}

status() {
    log INFO "Checking rhilex status."
    pid=\$(pgrep -x -n "rhilex")
    if [ -n "\$pid" ]; then
        log INFO "rhilex is running with Pid:\${pid}"
    else
        log INFO "rhilex is not running."
    fi
}

daemon() {
    while true; do
        if pgrep -x "rhilex" > /dev/null; then
            log INFO "rhilex process exists"
            sleep 3
            continue
        fi
        if ! pgrep -x "rhilex" > /dev/null; then
            if [ -e "$UPGRADE_SIGNAL" ]; then
                log INFO "File $UPGRADE_SIGNAL exists. May upgrade now."
                sleep 2
                continue
            elif [ -e "$STOP_SIGNAL" ]; then
                log INFO "$STOP_SIGNAL file found. Exiting."
                exit 0
            else
                log WARNING "Detected that rhilex process is interrupted. Restarting..."
                cd $WORKING_DIRECTORY
                $EXECUTABLE_PATH run -config=$CONFIG_PATH
                log WARNING "Detected that rhilex process has Restarted."
            fi
        fi
        sleep 4
    done
}

case "\$1" in
    start)
        start
    ;;
    restart)
        restart
    ;;
    stop)
        stop
    ;;
    status)
        status
    ;;
    *)
        log INFO "Usage: \$0 {start|restart|stop|status}"
        log INFO "You must specify one of the following options:"
        log INFO "    \$0 start    - Start the service"
        log INFO "    \$0 restart  - Restart the service"
        log INFO "    \$0 stop     - Stop the service"
        log INFO "    \$0 status   - Check the status of the service"
        exit 1
    ;;
esac

EOL

    mkdir -p $WORKING_DIRECTORY
    chmod +x "$SOURCE_DIR/$SERVICE_NAME"
    if cp -rfp "$SOURCE_DIR/$SERVICE_NAME" "$EXECUTABLE_PATH"; then
        log INFO "Copy rhilex to $WORKING_DIRECTORY"
    else
        log ERROR "Failed to copy rhilex to $WORKING_DIRECTORY"
        exit 1
    fi

    if cp -rfp "$SOURCE_DIR/$SERVICE_NAME.ini" "$CONFIG_PATH"; then
        log INFO "Copy rhilex.ini to $WORKING_DIRECTORY"
    else
        log ERROR "Failed to copy rhilex.ini to $WORKING_DIRECTORY"
        exit 1
    fi

    if cp -rfp "$SOURCE_DIR/license.lic" "$WORKING_DIRECTORY/"; then
        log INFO "Copy license.lic to $WORKING_DIRECTORY"
    else
        log ERROR "Failed to copy license.lic to $WORKING_DIRECTORY"
        exit 1
    fi

    if cp -rfp "$SOURCE_DIR/license.key" "$WORKING_DIRECTORY/"; then
        log INFO "Copy license.key to $WORKING_DIRECTORY"
    else
        log ERROR "Failed to copy license.key to $WORKING_DIRECTORY"
        exit 1
    fi

    chmod 755 $SERVICE_FILE
    if [ $? -eq 0 ]; then
        log INFO "rhilex service has been created and extracted."
    else
        log ERROR "Failed to create the rhilex service or extract files."
    fi
    exit 0
}

# 移除文件的函数，根据文件是否为目录选择不同的删除方式
__remove_files() {
    local file=$1
    log INFO "Removing $file."
    if [ -e "$file" ]; then
        if [ -d "$file" ]; then
            rm -rf "$file"
        else
            rm "$file"
        fi
        log INFO "$file removed."
    else
        log INFO "$file not found. No need to remove."
    fi
}

# 卸载服务的函数
uninstall() {
    if [ -e "$SERVICE_FILE" ]; then
        $SERVICE_FILE stop
    fi
    __remove_files "$SERVICE_FILE"
    __remove_files "$WORKING_DIRECTORY/license.*"
    __remove_files "$WORKING_DIRECTORY/$SERVICE_NAME*"
    __remove_files "$WORKING_DIRECTORY/md5.sum"
    __remove_files "$WORKING_DIRECTORY/*.txt"
    __remove_files "$WORKING_DIRECTORY/upload/"
    __remove_files "$WORKING_DIRECTORY/zold/"
    __remove_files "$WORKING_DIRECTORY/old/"
    __remove_files "$WORKING_DIRECTORY/zupgrade/"
    __remove_files "$STOP_SIGNAL"
    __remove_files "$UPGRADE_SIGNAL"
    log INFO "rhilex has been uninstalled."
}

# 调用服务文件的 start 命令启动服务
start() {
    $SERVICE_FILE start
}

# 调用服务文件的 restart 命令重启服务
restart() {
    $SERVICE_FILE restart
}

# 调用服务文件的 stop 命令停止服务
stop() {
    $SERVICE_FILE stop
}

# 调用服务文件的 status 命令查看服务状态
status() {
    $SERVICE_FILE status
}

# 根据输入的参数执行对应的操作
case "$1" in
    install)
        install
    ;;
    start)
        start
    ;;
    restart)
        restart
    ;;
    stop)
        stop
    ;;
    uninstall)
        uninstall
    ;;
    status)
        status
    ;;
    *)
        log INFO "Usage: $0 {start|restart|stop|status}"
        log INFO "You must specify one of the following options:"
        log INFO "    $0 start    - Start the service"
        log INFO "    $0 restart  - Restart the service"
        log INFO "    $0 stop     - Stop the service"
        log INFO "    $0 status   - Check the status of the service"
        exit 1
    ;;
esac

exit 0
