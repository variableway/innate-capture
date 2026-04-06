# Capture System Specifications

> **Version**: 1.0  
> **Status**: Draft

---

## Specification Index

| Spec | Description | Status |
|------|-------------|--------|
| [Issue Model Spec](./issue-model-spec.md) | Issue 数据模型、状态流转、生命周期 | Draft |
| [Agent Runtime Spec](./agent-runtime-spec.md) | Agent 运行时、Daemon、任务执行 | Draft |
| [Wiki Integration Spec](./wiki-integration-spec.md) | Wiki 知识库集成 (Karpathy 模式) | Draft |
| [Storage Spec](./storage-spec.md) | 存储层、SQLite 双写策略 | Draft |
| [Sync Spec](./sync-spec.md) | Feishu Bitable 同步 | Draft |

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Capture System                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Input Layer                                                             │
│  ───────────                                                             │
│  • Terminal CLI (capture add/list/status)                               │
│  • TUI Kanban (bubbletea)                                               │
│  • Feishu Bot (webhook/websocket)                                       │
│  • GitHub Issues (webhook)                                              │
│                                                                          │
│  Core Layer                                                              │
│  ──────────                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │ Issue Service│  │  Task Queue  │  │   Wiki       │                  │
│  │              │  │              │  │  Service     │                  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                  │
│         │                 │                 │                            │
│         └────────┬────────┴────────┬────────┘                            │
│                  ▼                 ▼                                     │
│         ┌────────────────────────────────┐                               │
│         │         Storage Layer          │                               │
│         │  Markdown (Source) + SQLite    │                               │
│         └────────────────────────────────┘                               │
│                          │                                               │
│                          ▼                                               │
│         ┌────────────────────────────────┐                               │
│         │       Sync Layer               │                               │
│         │   Feishu Bitable Integration   │                               │
│         └────────────────────────────────┘                               │
│                                                                          │
│  Execution Layer                                                         │
│  ───────────────                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                                  │
│  │ Daemon  │  │ Daemon  │  │ Daemon  │  (Multi-Machine)                  │
│  │ (Claude)│  │ (Codex) │  │(OpenCode)│                                  │
│  └─────────┘  └─────────┘  └─────────┘                                  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Key Design Decisions

### 1. Local First

- Markdown files are the source of truth
- SQLite is auto-generated index
- Git provides version control
- Feishu is optional cloud enhancement

### 2. Issue + Wiki Dual Model

- **Issue**: Executable work unit with lifecycle
- **Wiki**: Accumulated knowledge base
- Issue execution results feed into Wiki
- Wiki knowledge informs new Issues

### 3. Multi-Agent Support

- Support Claude, Codex, OpenCode, and custom agents
- Each agent has its own workdir and config
- Tasks are dispatched based on agent capabilities
- Multi-machine deployment supported

### 4. Stage Pipeline

```
inbox → analysis → planning → dispatch → execution → review → complete
```

Each stage can have automated actions and manual approval gates.

---

## Data Flow

### Issue Creation Flow

```
Terminal Input
      │
      ▼
Create Issue (inbox stage)
      │
      ▼
Auto Ingest to Wiki
      │
      ├── Create entity page
      ├── Extract concepts
      └── Update index
      │
      ▼
Sync to Feishu
```

### Task Execution Flow

```
Issue in dispatch stage
      │
      ▼
Daemon claims task
      │
      ▼
Agent executes
      │
      ▼
Archive result to Wiki
      │
      ├── Create execution log
      ├── Extract components
      └── Update related entities
      │
      ▼
Update Issue status
      │
      ▼
Notify user (Feishu)
```

---

## Implementation Phases

### Phase 1: Core Framework (4-6 weeks)

- [ ] Issue data model implementation
- [ ] Markdown storage layer
- [ ] SQLite index
- [ ] Basic CLI commands

### Phase 2: Agent Integration (4-6 weeks)

- [ ] Daemon framework
- [ ] Agent backends (Claude/Codex)
- [ ] Task queue
- [ ] Execution flow

### Phase 3: Wiki & Sync (4 weeks)

- [ ] Wiki directory structure
- [ ] Issue Ingestion
- [ ] Execution archival
- [ ] Feishu sync

### Phase 4: Advanced Features (4 weeks)

- [ ] TUI Kanban
- [ ] Wiki query/lint
- [ ] Multi-machine sync
- [ ] Documentation

---

## Related Documents

- [Project Analysis Report](../../features/analyasis_report.md)
- [Karpathy Wiki Integration Analysis](../../features/karpathy_wiki_integration_analysis.md)
