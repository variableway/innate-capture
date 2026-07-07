# Issue: cap-ws-config

**Status:** open  
**Mode:** horizontal  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** —
**Assignee:** —  
**Branch:** `feature/cap-ws-config-<agent>`  
**Contract:** [workspace-io-v1](../spec/contracts/workspace-io-v1.md)

## 背景

capture 需读取 `workspace.root` 并解析 `ideas/inbox`、`daily/today.md` 绝对路径。

## 范围

**做：**

- `WorkspaceConfig` 并入 `model.Config` / viper
- `internal/workspace`：`Root()`, `InboxDir()`, `DailyPath()`, `Validate()`
- `CAPTURE_WORKSPACE_ROOT` 环境变量
- `capture config` 可 get/set workspace 项（或 `config set workspace.root`）
- 可选：`capture doctor` 检查 innate-works 目录结构

**不做：**

- idea/daily 业务逻辑
- 飞书配置变更

## 可改动的路径

- `internal/config/model.go`
- `internal/workspace/`
- `cmd/config.go`
- `cmd/doctor.go`（新建可选）

**禁止改动：** `internal/feishu/`, `internal/bot/`

## 验收标准

- [ ] 合法 root 解析出 inbox / daily 路径
- [ ] 缺 `ideas/` 或 `daily/` 时 Validate 失败并提示
- [ ] 单测覆盖路径拼接与环境变量覆盖

## 依赖

- 阻塞：无  
- 流水线上游：无

## 过程日志

[`tasks/issues/cap-ws-config.md`](../tasks/issues/cap-ws-config.md)
