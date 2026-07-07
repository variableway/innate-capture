# Conventions: innate-capture

## 代码

- Go 模块：`projects/capture/`
- 新包：`internal/workspace`、`internal/idea`、`internal/daily`
- 测试：新包需 table-driven 单测；路径解析覆盖 edge case

## Git 与多 Agent

- 分支：`feature/<issue-id>-<agent>`
- 不并行改同一 issue
- PR 合并后更新 `features/` 状态与 task log

## Agent 约定

1. 入口：[`../index.md`](../index.md)  
2. 合约：[`contracts/workspace-io-v1.md`](./contracts/workspace-io-v1.md)  
3. 做完追加 [`../tasks/issues/`](../tasks/issues/)  

## 与 innate-works

- inbox 格式对齐 `innate-works/ideas/_template/inbox-entry.md`
- daily 不修改 checkbox 语义，只读或 `--open` 交编辑器
