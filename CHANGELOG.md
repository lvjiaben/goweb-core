# Changelog

## v0.9.0-rc.1

Release candidate，仍处于 `1.0` 前的收口阶段，不是 final stable。

已完成：

- `app` 应用启动与关闭封装
- `config` YAML 配置读取
- `db` PostgreSQL + GORM 初始化
- `httpx` 基于 `net/http` 的轻量路由、上下文、中间件、统一响应
- `auth` JWT 基础能力
- `rbac` 按 `permission_code` 的抽象与中间件接口
- `files` 本地文件存储辅助
- `validate` validator 封装
- `errorsx` 统一业务错误
- `logx` 日志初始化

不包含：

- 业务模块
- 代码生成器
- 队列、定时任务、WebSocket
- 用户端或后台页面

说明：

- 当前版本用于给 `goweb-scaffold` 提供稳定运行时内核
- 发布策略与兼容边界见 `docs/versioning.md`、`docs/compatibility.md`
