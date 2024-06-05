# 通用读写串口
主要用来操作普通串口，只有读操作，不支持指令控制.

## 配置
```json
{
    "name": "GENERIC_UART_RW",
    "type": "GENERIC_UART_RW",
    "gid": "DROOT",
    "config": {
        "commonConfig": {
            "tag": "uart",
            "timeSlice": 50,
            "autoRequest": true,
            "readFormat": "HEX"
        },
        "portUuid": "COM3"
    },
    "description": "GENERIC_UART_RW"
}
```