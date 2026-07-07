# Feature: tui-capture

**Status:** planned  
**Issues:** [cap-tui-idea](../issues/cap-tui-idea.md)

## 说明

在现有 kanban TUI 增加：快捷键 `i` 提交 idea；可选「今日」只读 daily 视图。

## 依赖

- [cap-idea-write](../issues/cap-idea-write.md) 的 `internal/idea` 包

## 验收

- [ ] TUI 提交生成与 CLI 相同 inbox 文件
- [ ] 不破坏现有 kanban 主流程（P2 可渐进）
