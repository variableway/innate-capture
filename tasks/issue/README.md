# Issue Tracking

> **Project**: Capture System  
> **Status**: Planning Phase

---

## Issue Index

| ID | Title | Priority | Status | Estimate |
|----|-------|----------|--------|----------|
| [001](./001-core-data-model.md) | Core Data Model | P0 | Planned | 1w |
| [002](./002-markdown-storage.md) | Markdown Storage | P0 | Planned | 1w |
| [003](./003-sqlite-index.md) | SQLite Index | P0 | Planned | 1w |
| [004](./004-basic-cli.md) | Basic CLI | P0 | Planned | 1w |
| [005](./005-agent-daemon.md) | Agent Daemon | P0 | Planned | 2w |
| [006](./006-agent-backends.md) | Agent Backends | P0 | Planned | 2w |
| [007](./007-wiki-structure.md) | Wiki Structure | P1 | Planned | 3d |
| [008](./008-issue-ingestion.md) | Issue Ingestion | P1 | Planned | 1w |
| [009](./009-execution-archival.md) | Execution Archival | P1 | Planned | 1w |
| [010](./010-feishu-sync.md) | Feishu Sync | P1 | Planned | 1w |
| [011](./011-tui-kanban.md) | TUI Kanban | P1 | Planned | 1w |
| [012](./012-feishu-bot.md) | Feishu Bot | P1 | Planned | 1w |

---

## Priority Legend

- **P0**: Must have - Core functionality
- **P1**: Should have - Important features
- **P2**: Nice to have - Enhancements

---

## Implementation Phases

### Phase 1: Core Framework (Weeks 1-4)

P0 Issues:
1. Core Data Model
2. Markdown Storage
3. SQLite Index
4. Basic CLI

**Goal**: Working CLI for Issue CRUD

### Phase 2: Agent Integration (Weeks 5-8)

P0 Issues:
5. Agent Daemon
6. Agent Backends

**Goal**: Agent can execute tasks

### Phase 3: Wiki & Cloud (Weeks 9-12)

P1 Issues:
7. Wiki Structure
8. Issue Ingestion
9. Execution Archival
10. Feishu Sync

**Goal**: Knowledge accumulation and cloud sync

### Phase 4: UI & Polish (Weeks 13-16)

P1 Issues:
11. TUI Kanban
12. Feishu Bot

**Goal**: Complete user experience

---

## Issue Template

```markdown
# Issue XXX: Title

> **Type**: Feature|Bug|Docs  
> **Priority**: P0|P1|P2  
> **Status**: Planned|In Progress|Done  
> **Assignee**: TBD  
> **Estimate**: X weeks/days

---

## Description

Brief description of what needs to be done.

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2

## Tasks

### 1. Task Name
Description and code example.

### 2. Task Name
Description and code example.

## Technical Notes

Important implementation details.

## Dependencies

- Issue XXX: Dependency

## Related Specs

- [Spec Name](../prd/spec/spec-file.md)
```

---

## Related Documents

- [PRD Specifications](../prd/spec/)
- [Project Analysis](../../features/analyasis_report.md)
- [Karpathy Wiki Integration](../../features/karpathy_wiki_integration_analysis.md)
