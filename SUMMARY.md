# Capture — Project Summary

> **One-liner**: A Go-based CLI/TUI tool for capturing ideas from Terminal and Feishu (飞书), saving them as Markdown files and syncing to Feishu Bitable.

---

## Purpose

Capture was born from the need to quickly record random ideas and tasks that pop up during work or daily life — whether at the terminal or via a Feishu message — and eventually turn them into executable, trackable work items. It bridges the gap between "fleeting thoughts" and "structured tasks" with minimal friction.

---

## Architecture

```
innate-capture/
├── projects/capture/          # Main Go application
│   ├── cmd/                   # Cobra CLI commands (add, list, show, edit, delete, status, stage, assign, kanban, bot, sync, init, config)
│   ├── internal/
│   │   ├── model/             # Task model, config model, status state machine
│   │   ├── store/             # Dual-write storage: Markdown (source of truth) + SQLite (fast index)
│   │   ├── service/           # TaskService business logic with functional options
│   │   ├── bot/               # Feishu Bot handlers (Webhook + WebSocket modes)
│   │   ├── bitable/           # Feishu Bitable API client and sync engine
│   │   ├── feishu/            # Shared Feishu SDK wrapper
│   │   ├── tui/               # bubbletea-based interactive kanban board
│   │   └── config/            # Viper-based configuration management
│   ├── pkg/
│   │   ├── idgen/             # TASK-NNNNN sequential ID generator
│   │   └── frontmatter/       # YAML frontmatter parser
│   └── main.go                # Entry point
├── linear/                    # Git submodule — Linear SDK reference
├── multica/                   # Git submodule — Multica reference implementation
├── tasks/                     # Development tracking & specifications
│   ├── prd/                   # Product requirements and specs
│   ├── issue/                 # Implementation issues (001–012)
│   ├── planning/              # AI-generated planning documents
│   ├── analysis/              # Feasibility and architecture analysis
│   └── features/              # Feature exploration documents
└── docs/                      # User-facing and internal documentation
```

---

## Tech Stack

| Layer        | Technology                                      |
|--------------|-------------------------------------------------|
| Language     | Go 1.26.1                                       |
| CLI Framework| Cobra                                           |
| TUI Framework| bubbletea + lipgloss                            |
| Config       | Viper                                           |
| Storage      | Markdown files + SQLite (`modernc.org/sqlite`)  |
| Feishu SDK   | `github.com/larksuite/oapi-sdk-go/v3`           |
| Testing      | testify                                         |

> **Note**: Pure Go — no CGo required.

---

## Key Features

### Implemented
- **CLI Task Management** — Create, list, show, edit, delete, and change task status/stage from the terminal
- **Markdown + YAML Frontmatter Storage** — Tasks stored as human-readable files at `~/.capture/tasks/YYYY/MM/TASK-NNNNN.md`
- **SQLite Index** — Fast querying and filtering alongside file-based source of truth
- **Status & Stage Workflows** — Status: `todo → in_progress → done` (or `todo → cancelled → archived`); Stage pipeline: `inbox → mindstorm → analysis → planning → prd → tasks → dispatch → execution → review`
- **AI Agent Assignment** — Record agent context (agent name, model, repo, worktree, terminal session) for dispatched tasks
- **TUI Kanban** — Interactive terminal board for visualizing and managing tasks
- **Feishu Bot** — Webhook and WebSocket modes; supports Chinese commands (`记录`, `列出`, `删除`, `帮助`)
- **Feishu Bitable Sync** — Push and bidirectional sync with Feishu spreadsheets

### Planned / In Exploration
- **Agent Execution Engine** — Direct integration with Claude Code, Codex, or other AI agents to execute tasks
- **Notification System** — Task status and execution result notifications via Feishu/WeChat
- **Wiki Integration** — Knowledge accumulation and RAG-style local research wiki

---

## Key Design Decisions

1. **Dual-Write Storage** — Markdown files are the source of truth; SQLite provides fast querying. This keeps data portable and human-readable while enabling performant filtering.
2. **Pure Go SQLite** — Uses `modernc.org/sqlite` to avoid CGo, simplifying cross-compilation and deployment.
3. **Two Bot Modes** — WebSocket for local development (no public URL needed); Webhook for production deployment.
4. **Functional Options Pattern** — `TaskService` uses `TaskOption` for flexible, backward-compatible mutations.
5. **Sequential ID Generator** — File-based counter with mutex produces human-friendly `TASK-NNNNN` IDs.

---

## Task Model

```yaml
---
id: TASK-00001
title: "优化项目构建脚本"
status: todo
stage: analysis
priority: high
tags: [优化, 构建]
created_at: 2026-04-03T10:30:00+08:00
updated_at: 2026-04-03T10:30:00+08:00
source: cli
dispatch:
  agent: "codex"
  model: "gpt-5"
  repository: "/Users/demo/workspace/project"
  worktree: "/Users/demo/workspace/project"
  terminal_session: "term-1"
  assigned_at: 2026-04-03T11:00:00+08:00
sync:
  feishu_record_id: ""
  last_synced_at: null
---
```

---

## Development Status

| Phase | Focus | Issues | Status |
|-------|-------|--------|--------|
| Phase 1 | Core Framework (CLI, Storage, Data Model) | 001–004 | ✅ Implemented |
| Phase 2 | Agent Integration (Execution Engine) | 005–006 | 🔄 Planned |
| Phase 3 | Wiki & Cloud Sync | 007–010 | 🔄 Partially Implemented (Sync) |
| Phase 4 | UI & Bot Polish | 011–012 | ✅ Implemented |

---

## Quick Start

```bash
cd projects/capture

# Build
go build -o capture .

# Initialize
capture init

# Add a task
capture add "优化项目构建脚本" -d "减少构建时间" -t "优化,构建" -p high

# List tasks
capture list

# Start TUI kanban
capture kanban

# Run tests
go test ./...
```

---

## Environment Variables

```bash
FEISHU_APP_ID               # Feishu app ID (for bot)
FEISHU_APP_SECRET           # Feishu app secret (for bot)
FEISHU_VERIFICATION_TOKEN   # Webhook verification token
FEISHU_ENCRYPT_KEY          # Webhook encrypt key
FEISHU_BITABLE_APP_TOKEN    # Bitable app token (for sync)
FEISHU_BITABLE_TABLE_ID     # Bitable table ID (for sync)
```

---

## Related Resources

- [AGENTS.md](./AGENTS.md) — Agent-specific coding guidance
- [CLAUDE.md](./CLAUDE.md) — Claude Code development instructions
- [tasks/issue/](./tasks/issue/) — Implementation issues and roadmap
- [tasks/prd/spec/](./tasks/prd/spec/) — Product specifications
- [docs/usage/](./docs/usage/) — User-facing documentation

---

**License**: MIT
