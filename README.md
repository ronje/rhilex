# RHILEX
RHILEX 商业版本仓库。

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

## 工作流程
- [指南](./contribute.md "指南")