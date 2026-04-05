# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make build              # Build binary (./capture)
make run                # Build and run
make test               # Run all tests (verbose)
make test-coverage      # Generate HTML coverage report
make fmt                # Format code (gofmt + goimports)
make lint               # Run golangci-lint
make clean              # Remove binary and coverage files

# Run a single test
go test ./internal/store/ -v -run TestMarkdownStore
go test ./internal/model/ -v -run TestTask
go test ./internal/bot/  -v -run TestMsgParser
```

## Architecture

Capture is a Go CLI/TUI tool for capturing ideas from Terminal and Feishu (飞书), storing them as Markdown files with YAML frontmatter and syncing to Feishu Bitable.

**Entry point**: `main.go` → `cmd.Execute()` (Cobra CLI)

**Layers**:
- `cmd/` — Cobra commands (add, list, show, edit, delete, status, stage, assign, kanban, bot, bot_serve, sync, config, init). Each command wires up `service.TaskService` with `store.MarkdownStore`.
- `internal/service/` — Business logic. `TaskService` uses functional options (`TaskOption`) for flexible mutations. Enforces status transition rules via `model.CanTransition()`.
- `internal/store/` — `Store` interface with `MarkdownStore` implementation. Files stored at `~/.capture/tasks/YYYY/MM/TASK-NNNNN.md` with YAML frontmatter. Markdown is the source of truth.
- `internal/model/` — Task model with status state machine (todo → in_progress → done, todo → cancelled → archived) and stage pipeline (inbox → mindstorm → analysis → planning → prd → tasks → dispatch → execution → review).
- `internal/bot/` — Feishu Bot handlers supporting both Webhook (HTTP) and WebSocket (long connection) modes. Parses Chinese commands: 记录/列出/删除/帮助.
- `internal/bitable/` — Feishu Bitable API client for syncing tasks to spreadsheets.
- `internal/feishu/` — Shared Feishu SDK wrapper.
- `internal/tui/` — bubbletea-based kanban board.
- `internal/config/` — Viper-based config from `~/.capture/config.yaml` with `CAPTURE_` env prefix.
- `pkg/idgen/` — Sequential TASK-NNNNN ID generator using file-based counter with mutex.
- `pkg/frontmatter/` — YAML frontmatter parser.

**Key patterns**:
- Storage interface (`store.Store`) allows swapping Markdown for SQLite or dual-write.
- Functional options pattern for task creation/update.
- Config: `~/.capture/config.yaml`, env vars with `CAPTURE_` prefix, or `--config`/`--data-dir` flags.

## Environment Variables

Required for Feishu Bot: `FEISHU_APP_ID`, `FEISHU_APP_SECRET`. See `.env.example` for full list.

## Module

`github.com/variableway/innate/capture` — Go 1.26.1, pure Go (no CGo, uses `modernc.org/sqlite`).
