# 轻量级MqttServer
## 配置
```json
{
    "type": "GENERIC_MQTT_SERVER",
    "name": "GENERIC_MQTT_SERVER Server",
    "description": "GENERIC_MQTT_SERVER Server",
    "config": {
        "serverName": "GENERIC_MQTT_SERVER",
        "host": "0.0.0.0",
        "port": 1883,
        "anonymous": true
    }
}
```
## 消息
```go
type MqttEvent struct {
	Action    string `json:"action"`
	Clientid  string `json:"clientid"`
	Username  string `json:"username"`
	Ipaddress string `json:"ipaddress"`
	Ts        int64  `json:"ts"`
	Topic     string `json:"topic,omitempty"`
	Payload   string `json:"payload,omitempty"`
}
```
