# 通用控制串口
主要用来操作普通串口,支持控制操作。

## 配置
```json
{
    "name": "GENERIC_USER_PROTOCOL",
    "type": "GENERIC_USER_PROTOCOL",
    "description": "GENERIC_USER_PROTOCOL",
    "gid": "DROOT",
    "config": {
        "commonConfig": {
            "mode": "UART",
            "autoRequest": true,
            "batchRequest": false
        },
        "hostConfig": {
            "host": "127.0.0.1",
            "port": 6001,
            "timeout": 5000
        },
        "uartConfig": {
            "uart": "COM1",
            "baudRate": 9600,
            "dataBits": 8,
            "stopBits": 1,
            "parity": "N"
        }
    }
}
```