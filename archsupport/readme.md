# 跨平台

## 环境变量

如果要启用硬件特性，需要在启动的时候加入 `ARCHSUPPORT=$V` 环境变量来指定运行的版本, 例如要在 RHILEXG1-H3网关上运行：

```sh
export ARCHSUPPORT=RHILEXG1
rhilex run
```

## 支持硬件列表

| 硬件名              | 环境参数   | 示例                                |
| ------------------- | ---------- | ----------------------------------- |
| RHILEXPro1 版本网关 | RHILEXPRO1 | `ARCHSUPPORT=RHILEXPRO1 rhilex run` |
| RHILEXG1 版本网关   | RHILEXG1   | `ARCHSUPPORT=RHILEXG1 rhilex run`   |
| RHILEXG2 版本网关   | RHILEXG2   | `ARCHSUPPORT=RHILEXG2 rhilex run`   |
| 树莓派4             | RPI4       | `ARCHSUPPORT=RPI4B rhilex run`      |
| 玩客云S805          | WKYS805    | `ARCHSUPPORT=WKYS805 rhilex run`    |
| HAAS506             | HAAS506    | `ARCHSUPPORT=HAAS506 rhilex run`    |
