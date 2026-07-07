# Issue: cap-docs-purge

**Status:** open  
**Mode:** pipeline  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** —  
**Assignee:** human  
**Branch:** —  
**Contract:** —

## 背景

**Task 1 最后一项。** innate-capture 旧文档整删，避免 Agent 误读。

## 范围

**删除（不迁移）：**

- `docs/analysis/`
- `docs/use-cases/`
- `docs/features/`（旧，含 svg）
- `docs/usage/idea-capture.md`, `trae-cli-cheatsheet.md`, `feishu-setup.md`
- `docs/example-feature.md`
- 旧 `docs/README.md`（已由新 index 替代）

**保留：**

- `docs/index.md`, `docs/usage/workspace.md`, `docs/project/**`

**代码保留：** `internal/feishu/`, `bot`, `sync`

## 验收标准

- [ ] E2E（§spec README）已通过
- [ ] `docs/` 下无 analysis、use-cases、旧 usage
- [ ] git 提交仅含删除 + 已存在新文档

## 依赖

- 阻塞：**全部 P1 issue done** + E2E 通过

## 过程日志

[`tasks/issues/cap-docs-purge.md`](../tasks/issues/cap-docs-purge.md)
