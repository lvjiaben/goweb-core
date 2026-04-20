# Release Checklist

## 文档与版本

- [ ] `VERSION` 已更新
- [ ] `CHANGELOG.md` 已写明当前 release candidate 状态
- [ ] `README.md` 已包含版本与文档入口
- [ ] `docs/versioning.md` 已更新
- [ ] `docs/compatibility.md` 已更新
- [ ] `docs/releases/v0.9.0-rc.1.md` 已存在

## 代码与验证

- [ ] `go test ./...`
- [ ] `go build ./...`
- [ ] 运行时边界未引入业务模块依赖

## 发布动作

- [ ] annotated tag 已尝试创建
- [ ] tag push 结果已记录
- [ ] release candidate 状态已明确说明为非 final stable
