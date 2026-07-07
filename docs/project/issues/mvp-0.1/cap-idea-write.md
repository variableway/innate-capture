# Issue: cap-idea-write

**Status:** open  
**Mode:** horizontal  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** [idea-capture](../features/idea-capture.md)  
**Assignee:** —  
**Branch:** `feature/cap-idea-write-<agent>`  
**Contract:** [workspace-io-v1](../spec/contracts/workspace-io-v1.md)

## 背景

实现 `capture idea add/list`，写入 innate-works inbox markdown。

## 范围

**做：**

- `internal/idea`：slug、模板渲染、Write、List
- `cmd/idea.go`：`add`, `list` 子命令
- flags：`-d` description, `-c` context

**不做：**

- SQLite 双写
- TUI（见 cap-tui-idea）

## 可改动的路径

- `internal/idea/`
- `cmd/idea.go`
- `cmd/root.go`（注册子命令）

## 验收标准

- [ ] `capture idea add "标题" -d "一句话"` 生成合法 inbox 文件
- [ ] slug 冲突自动 `-2`
- [ ] `capture idea list` 列出 inbox 下 md

## 依赖

- 阻塞：[cap-ws-config](./cap-ws-config.md)

## 过程日志

[`tasks/issues/cap-idea-write.md`](../tasks/issues/cap-idea-write.md)
