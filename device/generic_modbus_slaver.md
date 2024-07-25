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
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# Modbus Slaver
Modbus从机模式。
## 配置
TCP模式
```json
{
    "name": "GENERIC_MODBUS_SLAVER",
    "type": "GENERIC_MODBUS_SLAVER",
    "gid": "DROOT",
    "config": {
        "commonConfig": {
            "mode": "TCP"
        },
        "hostConfig": {
            "host": "0.0.0.0",
            "port": 6005
        }
    }
}
```
串口模式

```json
{
    "name": "GENERIC_MODBUS_SLAVER",
    "type": "GENERIC_MODBUS_SLAVER",
    "gid": "DROOT",
    "config": {
        "commonConfig": {
            "mode": "UART"
        },
        "portUuid": "/dev/ttyS1"
    }
}
```