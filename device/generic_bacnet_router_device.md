# Bacnet 路由器模式
单独提供一个专门用来映射其他类型设备的点位都Bacnet的接入网关。比如将Modbus采集到的数据映射成bacnet点位。

## 配置
```json
{
    "gid": "DROOT",
    "name": "GENERIC_BACNET_ROUTER",
    "type": "GENERIC_BACNET_ROUTER",
    "config": {
        "bacnetRouterConfig": {
            "deviceId": 123,
            "deviceName": "rhilex",
            "localPort": 47808,
            "mode": "BROADCAST",
            "netWorkId": 1,
            "networkCidr": "192.168.10.163/24",
            "vendorId": 123
        }
    },
    "description": "GENERIC_BACNET_ROUTER"
}
```