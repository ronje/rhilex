# 腾讯云物联网网关
主要用来对接腾讯云，参考文档：https://cloud.tencent.com/document/product/1081/34916。

## 测试配置
```sh
curl --location '127.0.0.1:2580/api/v1/devices/create' \
--header 'Content-Type: application/json' \
--data '{
    "name": "TENCENT_IOTHUB_GATEWAY",
    "type": "TENCENT_IOTHUB_GATEWAY",
    "gid": "DROOT",
    "description": "测试",
    "config": {
        "tencentConfig": {
            "mode": "DEVICE",
            "productId": "AEU6EFTRU6",
            "deviceName": "rhilex001",
            "devicePsk": "ysOqfTtUCaLBj5UeEOfCtQ=="
        }
    },
    "uuid": "TENCENT_IOTHUB_GATEWAY"
}'
```