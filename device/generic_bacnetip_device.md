## BACnet设备

### 设备type值
`GENERIC_BACNET_IP`
### 配置示例
#### 广播配置模式
```json
{
    "name": "GENERIC_BACNET_IP",
    "type": "GENERIC_BACNET_IP",
    "gid": "DROOT",
    "config": {
        "commonConfig": {
            "frequency": 1000
        },
        "bacnetConfig": {
            "interface": "ETH1",
            "mode": "BROADCAST",
            "localIp": "192.168.10.163",
            "subnetCidr": 24,
            "localPort": 47808
        }
    },
    "description": "GENERIC_BACNET_IP"
}
```


> 字段含义

- type: bacnet采集驱动运行模式（可选SINGLE、BROADCAST）
- ip:  bacnet设备ip（仅type=SINGLE时生效）
- port:  bacnet端口，通常是47808（仅type=SINGLE时生效）
- isMstp: 是否为mstp over ip设备，若是则子网号必须填写（仅type=SINGLE时生效）
- subnet: 虚拟子网号，一般是mstp转Ip网关用于给mstp组分配的一个网络号
- localIp: 本地地址（仅type=BROADCAST时生效）
- subnetCidr: 子网掩码长度（localIp和subnetCidr共同计算出广播地址）
- localPort:  本地监听端口，填0表示默认47808
- frequency:  采集间隔，单位毫秒

### 点位配置
- tag: 点位名称
- alias: 别名
- bacnetDeviceId: 设备id，整数类型（若isMstp=1，则deviceId应该必填；若是bacnetip设备，则填1即可）
- objectType: 点位类型，必填（下拉框，枚举值）
- objectId: 对象id，必填（范围0-4194303)

`objectType`可选类型
* AI
* AO
* AV
* BI
* BO
* BV
* MI
* MO
* MV

### 采集后输出的数据格式
```json
{
    "gas": {
        "tag": "gas",
        "deviceId": 1,
        "propertyType": "AnalogInput",
        "propertyInstance": 2,
        "value": 77.89
    },
    "humi": {
        "tag": "humi",
        "deviceId": 1,
        "propertyType": "AnalogInput",
        "propertyInstance": 1,
        "value": 56.78
    },
    "temp": {
        "tag": "temp",
        "deviceId": 1,
        "propertyType": "AnalogInput",
        "propertyInstance": 0,
        "value": 12.34
    }
}
```

### 注意事项
* 由于bacnet需要本地监听UDP端口用于收发信令，因此`localPort`不能配置重复
* 若使用广播模式`BROADCAST`，则建议`localPort`设置为47808
