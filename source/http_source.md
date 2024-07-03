# Http Server
提供HTTP数据推入接口。

## 数据格式
### 请求数据
```json
{
	"data": {
		"k":"v"
	}
}
```
### 成功返回
```json
{
	"message": "success",
	"code":    200
}
```
### 失败返回
```json
{
	"message": "error-msg",
	"code":    500
}
```
## 示例
```py
import requests
url = 'http://example.com/api'
data = {
    "data": {
        "k": "v"
    }
}
response = requests.post(url, json=data)
print(response.status_code)
print(response.text)
```