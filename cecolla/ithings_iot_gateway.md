# Ithings网关
主要用来对接Ithings：https://doc.ithings.net.cn/pages/68b74b。

## 测试数据
```sh
curl --location '127.0.0.1:2580/api/v1/devices/create' \
--header 'Content-Type: application/json' \
--data '{
    "name": "ITHINGS_IOTHUB_GATEWAY",
    "type": "ITHINGS_IOTHUB_GATEWAY",
    "gid": "DROOT",
    "config": {
        "ithingsConfig": {
            "serverEndpoint": "tcp://139.159.188.223:1883",
            "mode": "GATEWAY",
            "productId": "00c",
            "deviceName": "主机1",
            "devicePsk": "cJwhsaSdbecsdydkVr5bf1XS1p4="
        }
    },
    "description": "ITHINGS_IOTHUB_GATEWAY"
}'
```