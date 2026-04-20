# Release Policy

## 定位

`goweb-core` 当前采用 release candidate 方式发布，目标是稳定运行时边界，而不是提前声明 `1.0`。

## 当前规则

- 当前版本：`v0.9.0-rc.1`
- 未经过完整业务试跑前，不发布 `v1.0.0`
- `rc` 版本可以继续修正文档、边界说明、兼容性说明和轻量缺陷
- 运行时公共 API 的破坏性变更必须显式记录到 changelog 和 release notes

## 与 scaffold 的关系

- core 是运行时内核
- scaffold 是可直接运行的业务脚手架
- scaffold 依赖 core 的 semver 版本
- template version、lock/export/source migration 由 scaffold 单独管理
