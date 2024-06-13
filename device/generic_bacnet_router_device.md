# Bacnet 路由器模式
单独提供一个专门用来映射其他类型设备的点位都Bacnet的接入网关。比如将Modbus采集到的数据映射成bacnet点位。

## 配置
```json
{
    "gid": "DROOT",
    "name": "BACNET_ROUTER_GW",
    "type": "BACNET_ROUTER_GW",
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
    "description": "BACNET_ROUTER_GW"
}
```

## 读写
对虚拟点位进行操作：
```lua
Actions = {
    function(args)
        Debug("args=", args)
        local DeviceId = "DEVICEIBNVQPJ4"
        device:CtrlDevice(DeviceId, "setValue", json:T2J({ tag = "Temperature", value = 13.14 }))
        device:CtrlDevice(DeviceId, "setValue", json:T2J({ tag = "Humidity", value = 52.11 }))
        return true, args
    end
}
```