# goweb-core

`goweb-core` 只承载 Go Web 运行时内核，不放业务模块，不继承旧仓库后端架构。

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
