# Task log: cap-daily-read

Issue：[`issues/cap-daily-read.md`](../../issues/cap-daily-read.md)

---

## 2026-07-06 — human — opened

**做了：** issue 契约就绪  
**未做：** 实现  
**给下一 Agent：** 等 cap-ws-config；读 workspace-io-v1 daily 节  
**状态：** open

---

## 2026-07-06 — claude — implemented

**做了：** `internal/daily`（Read / BootstrapFromTemplate / PrintSection）+ `cmd/daily.go`（`daily` / `--open` / `--section`）；table-driven 单测 8 个全过；真实 workspace 冒烟通过（`daily` 全文 / `--section ideas` / `--section output` / 非法 section 报错退出 1）。  
**未做：** 未提交、未开 PR；`--open` 仅做参数与配置校验，未端到端（依赖用户 `defaults.editor` 在 PATH）；未改 `internal/workspace`——template 路径在 daily 包内拼接，符合本 issue「可改路径」。  
**给下一 Agent：** cap-ws-config 实际已就绪（`internal/workspace` 存在且被 `internal/idea` 使用），原「阻塞」状态过时；可开 PR，合并后翻 issue 与 `features/daily-read.md` 状态。  
**状态：** implemented (uncommitted)

---
