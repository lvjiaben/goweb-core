# Versioning

## 当前版本

- repo version: `v0.9.0-rc.1`
- 版本状态：release candidate

## 版本策略

`goweb-core` 只管理运行时内核版本，不管理业务模块版本，也不管理 codegen template version。

遵循 semver 的内容：

- `goweb-core` 仓库版本
- 对外公开的 Go 包边界

不在本仓库内管理的内容：

- scaffold 业务模块生命周期
- codegen lock/export/source version
- template version

## 兼容原则

- `v0.9.x` 内优先保持对 `goweb-scaffold v0.9.x` 的运行时兼容
- 破坏性变更需要进入下一个 semver 版本，而不是在 `rc` 内静默落地

## 当前边界

本仓库只包含：

- `app`
- `config`
- `db`
- `httpx`
- `auth`
- `rbac`
- `files`
- `validate`
- `errorsx`
- `logx`

不包含业务模块、codegen、菜单模型、角色模型和任何前端资源。
