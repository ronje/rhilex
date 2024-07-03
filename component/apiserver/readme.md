# Http Server
## 简介
HTTP Server 是 RHILEX 的 WEB API 提供者，主要用来支持 Dashboard 以及部分性能监控。
<img src="./structure.png"/>

## 注意
gin框架会自带默认值，如果接口对默认值有要求，建议将字段设置为指针形式.
```go
type S1 struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

type S2 struct {
	ID     *int  `json:"id"`
	Status *bool `json:"status"`
}

```
其中S1和S2的两个结构体的字段默认值是不同的，例如上面的例子：S2.ID如果不传会是0，但是当你传0的时候他还是0，这样就有歧义，为了避免问题，可以设计为指针形式。