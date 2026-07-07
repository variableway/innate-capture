# Issue: cap-docs-new

**Status:** open  
**Mode:** vertical  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** —  
**Assignee:** —  
**Branch:** `feature/cap-docs-new-<agent>`  
**Contract:** —

## 背景

用户向 CLI 文档：`docs/usage/workspace.md`；确保 `docs/index.md` 指向 project。

## 范围

**做：**

- 撰写 `docs/usage/workspace.md`（idea add、daily、config workspace.root）
- 核对 `docs/index.md`、`docs/project/*` 链接
- 根目录 `AGENTS.md` 指向 `docs/project/index.md`

**不做：**

- 删除旧 docs（见 cap-docs-purge）

## 验收标准

- [ ] 新用户可按 workspace.md 完成 E2E
- [ ] AGENTS.md 薄层 ≤15 行 + 链接

## 依赖

- 阻塞：cap-idea-write、cap-daily-read 行为稳定（可与实现末段并行）

## 过程日志

[`tasks/issues/cap-docs-new.md`](../tasks/issues/cap-docs-new.md)
