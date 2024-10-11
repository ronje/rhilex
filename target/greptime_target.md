# GrepTimeDB 北向资源
## 配置
定义:
```go
type GrepTimeDbTargetConfig struct {
	GwSn     string `json:"gwsn" validate:"required"`       // 服务地址
	Host     string `json:"host" validate:"required"`       // 服务地址
	Port     int    `json:"port" validate:"required"`       // 服务端口
	Username string `json:"username" validate:"required"`   // 用户
	Password string `json:"password" validate:"required"`   // 密码
	DataBase string `json:"database" validate:"required"`   // 数据库
	Table    string `json:"table" validate:"required"`      // 数据表名
}
```

JSON:
```json
{
  "gwsn": "rhilex",
  "host": "127.0.0.1",
  "port": 4001,
  "username": "rhilex",
  "password": "rhilex",
  "database": "rhilex",
  "table": "rhilex"
}

```

## 示例
```lua
local err = data:ToGreptimeDB('UUID', {"v1","v2","v3","v4"})
```