# Spec: innate-capture workspace 集成

## 问题

innate-works 需要快速捕获 idea 与查看每日清单；手工编辑 markdown 摩擦大。capture 改为 workspace 前端。

## MVP 范围

- [ ] `workspace.root` 配置与路径解析（`cap-ws-config`）
- [ ] `capture idea add/list` → `ideas/inbox/YYYY-MM-DD-<slug>.md`
- [ ] `capture daily` → 读/引导 `daily/today.md`
- [ ] 用户文档 `docs/usage/workspace.md`（`cap-docs-new`）
- [ ] E2E 通过后删除旧 `docs/`（`cap-docs-purge`）

## 非目标

- 飞书 sync / bot 启用（P4 另立项）
- 新 idea 写入 SQLite / `~/.capture/tasks/`
- 扩展旧 `capture add` 文档或行为
- TUI 为 P2，可延后

## 验收标准（项目级）

- [ ] `capture idea add "测试" -d "描述"` 在 innate-works `ideas/inbox/` 生成合法 md
- [ ] `capture daily` 输出 `daily/today.md` 内容
- [ ] `workspace.root` 错误时 `capture doctor`（或等价）给出明确提示
- [ ] 旧 `docs/analysis` 等已删除，仅留 `docs/project/` + `usage/workspace.md`

## 相关文档

- [architecture.md](./architecture.md)
- [contracts/workspace-io-v1.md](./contracts/workspace-io-v1.md)
- [conventions.md](./conventions.md)

## 变更记录

| 日期 | 变更 |
|---|---|
| 2026-07-06 | 自 exploring plan promote scaffold |
