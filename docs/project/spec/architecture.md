# Architecture: innate-capture workspace 集成

## 系统概览

```text
capture CLI/TUI
    │
    ├── internal/workspace   ← 读 config.workspace.root，解析路径
    ├── internal/idea        ← 渲染 inbox md，写 innate-works/ideas/inbox/
    └── internal/daily       ← 读/引导 innate-works/daily/today.md
              │
              ▼
    <innate-works>/ideas/inbox/ · daily/today.md
```

## 关键决策

| 日期 | 决策 | 原因 |
|---|---|---|
| 2026-07-06 | inbox 只写文件，不写 SQLite | 与 innate-works 单一事实来源一致 |
| 2026-07-06 | 飞书代码保留不调用 | P4 再启用 |
| 2026-07-06 | 旧 docs 整删不迁移 | 避免 Agent 误读 |

## 模块与路径

| 模块 | 路径 | 说明 |
|---|---|---|
| workspace | `internal/workspace/` | Root、InboxDir、DailyPath |
| idea | `internal/idea/` | Slug、模板、Write、List |
| daily | `internal/daily/` | Show、Open、Bootstrap |
| cmd | `cmd/idea.go`, `cmd/daily.go` | 子命令 |
| config | `internal/model/config.go` | `WorkspaceConfig` |
| 保留 | `internal/feishu/`, `bot/`, `bitable/` | 本阶段不动 |

## 外部依赖

| 依赖 | 用途 |
|---|---|
| innate-works FS | ideas/inbox、daily 模板与 today |
| viper | `workspace.*` 配置 |
