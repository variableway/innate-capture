# AGENTS.md

> **项目事实来源：** [`docs/project/index.md`](docs/project/index.md) — spec、issues、多 Agent 日志均在该目录。

## 入口

1. 读 [`docs/project/index.md`](docs/project/index.md)
2. 读当前 [`docs/project/issues/<id>.md`](docs/project/issues/)
3. 读 [`docs/project/tasks/issues/<id>.md`](docs/project/tasks/issues/) **最后一条**

## 事实来源

- 范围与验收：[`docs/project/spec/README.md`](docs/project/spec/README.md)
- 合约：[`docs/project/spec/contracts/workspace-io-v1.md`](docs/project/spec/contracts/workspace-io-v1.md)

## 交付

- 分支 `feature/<issue>-<agent>`；合并后追加 task log

## 禁止

- 不要创建 `.agents/` 或第二套 spec
- 不要未经 issue 扩大范围
- E2E 前不要执行 `cap-docs-purge`

## 代码根

Go 项目：`projects/capture/`（`go build` 在此目录）

用户 CLI 手册：[`docs/usage/workspace.md`](docs/usage/workspace.md)
