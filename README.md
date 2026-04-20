# goweb-core

`goweb-core` 只承载 Go Web 运行时内核，不放业务模块，不继承旧仓库后端架构。

## 当前版本

- repo version: `v0.9.0-rc.1`
- 发布状态：release candidate，不是 final stable
- 兼容 scaffold：`goweb-scaffold v0.9.0-rc.1` / `v0.9.x`

## 文档入口

- [Versioning](docs/versioning.md)
- [Compatibility](docs/compatibility.md)
- [Release Policy](docs/releases/RELEASE_POLICY.md)
- [Release Checklist](docs/releases/RELEASE_CHECKLIST.md)
- [Release Notes](docs/releases/v0.9.0-rc.1.md)

## 包结构

- `app`：应用启动与关闭封装
- `config`：YAML 配置加载
- `db`：PostgreSQL 与 GORM 初始化
- `httpx`：基于 `net/http` 的轻量路由、上下文、中间件、统一响应
- `auth`：JWT 签发与解析
- `rbac`：RBAC 接口约束与按 `permission_code` 校验的中间件
- `files`：本地文件存储帮助函数
- `validate`：`validator/v10` 封装
- `errorsx`：统一业务错误定义
- `logx`：基于 `slog` 的日志初始化

## 设计边界

- 只使用 `net/http`
- 数据库只针对 PostgreSQL
- ORM 只使用 GORM
- 不依赖 Gin、Chi、Casbin、RabbitMQ
- RBAC 只提供抽象，不提供具体业务表实现

## 边界说明

- core 只负责运行时内核，不负责业务模块
- codegen、lock/export/source migration 由 `goweb-scaffold` 管理
- 当前仓库不包含 admin/user 页面，不包含任何业务表实现

## 快速示例

```go
logger := logx.New(logx.Config{Level: "info"})
engine := httpx.NewEngine(logger)
engine.Use(
	httpx.RequestID(),
	httpx.Logger(logger),
	httpx.Recover(logger),
)

engine.GET("/ping", func(c *httpx.Context) {
	c.Success(map[string]any{"pong": true})
})

application := app.New("demo", ":8080", engine, logger)
if err := application.Run(); err != nil {
	panic(err)
}
```
