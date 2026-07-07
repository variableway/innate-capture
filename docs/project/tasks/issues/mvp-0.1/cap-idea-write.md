# Task log: cap-idea-write

Issue：[`issues/cap-idea-write.md`](../../issues/cap-idea-write.md)

---

## 2026-07-06 — human — opened

**做了：** issue 契约就绪  
**未做：** 实现  
**给下一 Agent：** 等 cap-ws-config 合并；读 workspace-io-v1 idea 节  
**状态：** open

---

## 2026-07-06 — agent — implemented

**做了：**
- `internal/idea`：`Slug`、`Write`、`List`、模板渲染（workspace-io-v1）
- `cmd/idea.go`：`capture idea add`（`-d`/`-c`）、`capture idea list`
- slug 冲突自动 `-2`、`-3`…
- 单测：`internal/idea/idea_test.go`、`slug_test` 逻辑在 `idea_test.go`

**验证：**
- `go test ./...` 通过
- E2E：`CAPTURE_WORKSPACE_ROOT=innate-works capture idea add/list` 成功

**给下一 Agent：** `cap-tui-idea` 可复用 `idea.Write(..., idea.SourceTUI)`  
**状态：** done（待合并）

---

## 2026-07-07 — agent — domain migration

**做了：**
- `internal/idea` 迁移为 domain 结构：`model.go` + `service.go` + `idea_service.go` + `repo/fs_repo.go`
- `Write/List` 改为调用 `DefaultService`，对外函数签名保持兼容
- `idea_service` 显式依赖 `workspace.Service`，便于后续扩展和测试替身注入

**验证：**
- `go test ./internal/workspace ./internal/idea ./cmd` 通过

**状态：** done（待合并）

---
