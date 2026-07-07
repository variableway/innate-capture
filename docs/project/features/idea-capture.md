# Feature: idea-capture

**Status:** planned  
**Issues:** [cap-idea-write](../issues/cap-idea-write.md), [cap-tui-idea](../issues/cap-tui-idea.md)

## 说明

`capture idea add/list` 将想法写入 innate-works `ideas/inbox/`，格式见 [workspace-io-v1](../spec/contracts/workspace-io-v1.md)。

## 相关代码

- `internal/idea/`
- `cmd/idea.go`

## 验收

- [ ] 文件名 `YYYY-MM-DD-<slug>.md`
- [ ] frontmatter 等价字段：Captured, Source, Stage=inbox
- [ ] `idea list` 按日期倒序列出 inbox
