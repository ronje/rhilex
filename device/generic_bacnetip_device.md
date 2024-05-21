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


字段含义

根据您的要求，这是优化后的Markdown表格：
| 参数名     | 描述                      | 类型/要求                                                                             |
| ---------- | ------------------------- | ------------------------------------------------------------------------------------- |
| type       | bacnet采集驱动运行模式    | 可选：`SINGLE`、`BROADCAST`                                                           |
| ip         | bacnet设备ip              | 仅当`type=SINGLE`时生效                                                               |
| port       | bacnet端口，默认为47808   | 整数类型<br>仅当`type=SINGLE`时生效                                                   |
| isMstp     | 是否为mstp over ip设备    | 布尔类型<br>若是，则`subnet`必须填写<br>仅当`type=SINGLE`时生效                       |
| subnet     | 虚拟子网号                | 整数类型<br>一般是mstp转Ip网关用于给mstp组分配的一个网络号<br>仅当`type=SINGLE`时生效 |
| localIp    | 本地地址                  | 仅当`type=BROADCAST`时生效                                                            |
| subnetCidr | 子网掩码长度              | 整数类型<br>`localIp`和`subnetCidr`共同计算出广播地址<br>仅当`type=BROADCAST`时生效   |
| localPort  | 本地监听端口，默认为47808 | 整数类型<br>填0表示默认47808<br>仅当`type=BROADCAST`时生效                            |
| frequency  | 采集间隔，单位毫秒        | 整数类型                                                                              |

### 点位配置
根据您的要求，这是优化后的Markdown表格：
| 参数名         | 描述     | 类型/要求                                                         |
| -------------- | -------- | ----------------------------------------------------------------- |
| tag            | 点位名称 |                                                                   |
| alias          | 别名     |                                                                   |
| bacnetDeviceId | 设备id   | 整数类型<br>若`isMstp=1`，则必填；<br>若是bacnetip设备，则填1即可 |
| objectType     | 点位类型 | 必填<br>下拉框，枚举值                                            |
| objectId       | 对象id   | 必填<br>范围0-4194303                                             |

> 请注意，`bacnetDeviceId` 字段的内容取决于 `isMstp` 的值，这可能在表格中不易表达。在实际应用中，可能需要额外的说明或逻辑来处理这个字段。

`objectType`可选类型

| ObjectType | 英文名称         | 中文名称   |
| ---------- | ---------------- | ---------- |
| AI         | AnalogInput      | 模拟输入   |
| AO         | AnalogOutput     | 模拟输出   |
| AV         | AnalogValue      | 模拟值     |
| BI         | BinaryInput      | 二进制输入 |
| BO         | BinaryOutput     | 二进制输出 |
| BV         | BinaryValue      | 二进制值   |
| MI         | MultiStateInput  | 多状态输入 |
| MO         | MultiStateOutput | 多状态输出 |
| MV         | MultiStateValue  | 多状态值   |

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
