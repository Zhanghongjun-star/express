# 顺风送运小程序（后端服务）

基于 [Kratos v3](https://github.com/go-kratos/kratos) 构建的微服务初始框架，包含 HTTP / gRPC 双协议、Protobuf 优先的 API、Wire 依赖注入、OpenAPI 生成以及 CRUD 示例。

## 技术栈

- Go 1.26+
- Kratos v3 微服务框架
- gRPC + HTTP 双协议
- Protobuf + buf 代码生成
- Wire 依赖注入
- OpenAPI (Swagger) 文档生成
- Docker 部署支持

## 项目结构

```text
api/                  Protobuf API 定义与生成代码
cmd/                  应用入口
configs/              配置文件
internal/server/      HTTP 和 gRPC 服务构造
internal/service/     面向传输层的服务方法
internal/biz/         业务用例、实体、错误、仓库接口
internal/data/        仓库实现
third_party/          Protobuf 依赖
openapi.yaml          生成的 OpenAPI 文档
Dockerfile
Makefile
```

## 开发命令

安装代码生成器：

```bash
make init
```

重新生成 API 和 OpenAPI：

```bash
make api
```

重新生成配置 protobuf：

```bash
make config
```

执行所有生成步骤、Wire 和模块清理：

```bash
make all
```

构建：

```bash
make build
```

测试：

```bash
go test ./...
```

## 本地运行

```bash
go run ./cmd/server -conf ./configs
```

默认端口：

- HTTP: `0.0.0.0:8000`
- gRPC: `0.0.0.0:9000`

## 团队使用说明

本仓库作为团队后端项目的初始框架，已经配置好本地 Kratos 框架依赖：

```go.mod
replace (
	github.com/go-kratos/kratos/v3 => C:/Users/31487/Desktop/express
	github.com/go-kratos/kratos/contrib/otel/v3 => C:/Users/31487/Desktop/express/contrib/otel
)
```

团队成员克隆后，需要保证本地 `C:\Users\31487\Desktop\express` 路径存在，或使用 `go env` 配置自己的本地框架路径。

## Docker 运行

```bash
docker build -t shunfeng-miniprogram .
docker run --rm -p 8000:8000 -p 9000:9000 -v ./configs:/data/conf shunfeng-miniprogram
```
