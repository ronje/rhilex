# Ithings网关
主要用来对接Ithings：https://doc.ithings.net.cn/pages/68b74b。

## 测试数据
```sh
curl --location '127.0.0.1:2580/api/v1/devices/create' \
--header 'Content-Type: application/json' \
--data '{
    "name": "ITHINGS_IOTHUB_CEC",
    "type": "ITHINGS_IOTHUB_CEC",
    "gid": "CEROOT",
    "config": {
        "serverEndpoint": "tcp://demo.ithings.net.cn:1883",
        "mode": "GATEWAY",
        "productId": "01D",
        "subProduct": "01J",
        "deviceName": "基站1",
        "devicePsk": "xSL17wB6+qtcj3n2Pqotfcr8WVE="
    },
    "description": "ITHINGS_IOTHUB_CEC"
}'
```