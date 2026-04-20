# Compatibility

## 当前兼容矩阵

| repo | version | compatible with | notes |
| --- | --- | --- | --- |
| `goweb-core` | `v0.9.0-rc.1` | `goweb-scaffold v0.9.0-rc.1` / `v0.9.x` | 作为运行时内核提供 `app/config/db/httpx/auth/rbac/files/validate/errorsx/logx` |

## 边界说明

- `goweb-core` 不感知 scaffold 的业务模块
- `goweb-core` 不直接管理 codegen template version
- scaffold 当前 template version 为 `v7`，但这属于 scaffold 侧兼容面，不属于 core 的版本面

## 发布前提

只有在 `goweb-scaffold` 的第十一阶段真实业务试跑通过后，才考虑把 `rc` 升为稳定版。
