# RHILEX
RHILEX 商业版本仓库。

<div style="background-color: #fcf8e3; border: 1px solid #faebcc; padding: 15px; color: #8a6d3b;">
  <strong>警告！</strong> 本软件源代码收中华人民共和国著作权法保护，未经许可获取即违反中华人民共和国著作权法，我们有权利追究责任。
</div>

## 构建
安装go1.18以后的版本。推荐使用1.22。

```sh
go get
go build
```

## 跨平台
```sh
# arm32
make arm32
# arm64
make arm64
```


## 条件编译
条件编译使用tag来控制，例如：
```sh
go build -tags gocv
```

## 发布
```sh
bash ./release_pkg.sh
```