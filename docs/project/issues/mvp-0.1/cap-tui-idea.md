# Issue: cap-tui-idea

**Status:** open  
**Mode:** horizontal  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** [tui-capture](../features/tui-capture.md)  
**Assignee:** —  
**Branch:** `feature/cap-tui-idea-<agent>`  
**Contract:** [workspace-io-v1](../spec/contracts/workspace-io-v1.md)

## 背景

kanban TUI 增加快速提交 idea（`i` 键）与可选今日视图。

## 范围

**做：**

- 复用 `internal/idea.Write`
- `internal/tui` 增加提交对话框
- Source 标 `capture-tui`

**不做：**

- 重写 kanban
- daily 编辑

## 可改动的路径

- `internal/tui/`
- `cmd/kanban.go`（若需挂接）

## 验收标准

- [ ] TUI 提交与 CLI 生成相同格式文件
- [ ] 无 workspace 配置时友好报错

## 依赖

- 阻塞：[cap-idea-write](./cap-idea-write.md)

## 过程日志

[`tasks/issues/cap-tui-idea.md`](../tasks/issues/cap-tui-idea.md)
