# Epic: cap-sprint-01 — workspace MVP

**Status:** open  
**Mode:** horizontal（Sprint 父任务）  
**Feature:** —  
**Assignee:** human（协调）  
**Branch:** —

## 背景

第一个 Sprint：让 capture 成为 innate-works 前端（idea inbox + daily）。对应 exploring plan P0–P3。

## Sprint 目标

- [ ] workspace 配置可用（`cap-ws-config`）
- [ ] `capture idea add/list` 写 inbox（`cap-idea-write`）
- [ ] `capture daily` 读 today（`cap-daily-read`）
- [ ] （可选）TUI 提交 idea（`cap-tui-idea`）
- [ ] 用户文档（`cap-docs-new`）
- [ ] E2E 通过后删旧 docs（`cap-docs-purge`，**最后**）

## 子 Issue

| # | Issue | Assignee | 状态 | 依赖 |
|---|---|---|---|---|
| 1 | [cap-ws-config](./cap-ws-config.md) | A | open | — |
| 2 | [cap-idea-write](./cap-idea-write.md) | B | open | ws-config |
| 3 | [cap-daily-read](./cap-daily-read.md) | C | open | ws-config |
| 4 | [cap-tui-idea](./cap-tui-idea.md) | D | open | idea-write |
| 5 | [cap-docs-new](./cap-docs-new.md) | E | open | idea + daily |
| 6 | [cap-docs-purge](./cap-docs-purge.md) | human | open | E2E |

```text
cap-ws-config ──┬──► cap-idea-write ──► cap-tui-idea
                └──► cap-daily-read ──► cap-docs-new ──► cap-docs-purge
```

## 本 Sprint 约束

- 飞书 **不启用**（代码保留）
- 新 idea **不写** SQLite
- 旧 `docs/analysis` 等 **等 purge 再删**

## 里程碑（Sprint 内）

| 日期 | 事件 |
|---|---|
| 2026-07-06 | plan → docs/project scaffold |
| — | P1 E2E：idea add + daily |
| — | cap-docs-purge 完成 → Sprint done |

## 过程日志

子 issue 记施工细节；Sprint 级摘要见 [`tasks/issues/cap-sprint-01-workspace-mvp.md`](../tasks/issues/cap-sprint-01-workspace-mvp.md)
