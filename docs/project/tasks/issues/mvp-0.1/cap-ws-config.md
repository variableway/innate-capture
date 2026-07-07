# Task log: cap-ws-config

Issue：[`issues/cap-ws-config.md`](../../issues/cap-ws-config.md)

---

## 2026-07-06 — human — opened

**做了：** promote scaffold；issue 自 exploring plan  
**证据：** `docs/project/issues/cap-ws-config.md`  
**未做：** 代码  
**给下一 Agent：** 实现 `internal/workspace` + config  
**状态：** open

---

## 2026-07-07 — agent — domain migration

**做了：**
- `internal/workspace` 迁移为 domain 结构：`model.go` + `service.go` + `workspace_service.go` + `repo/fs_repo.go`
- 保留兼容 API：`workspace.Root/InboxDir/DailyPath/Validate` 继续可用（内部委托 `DefaultService`）
- 保持原行为不变（env 覆盖、路径拼接、校验错误语义）

**验证：**
- `go test ./internal/workspace ./internal/idea ./cmd` 通过

**状态：** open

---
