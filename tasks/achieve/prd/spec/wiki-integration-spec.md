# Wiki Integration Specification

> **Version**: 1.0  
> **Status**: Draft  
> **Based on**: Karpathy Wiki Pattern + Issue System

---

## 1. Overview

将 Karpathy 的 LLM Knowledge Base 模式与 Issue 系统整合，实现知识的自动累积和复利。Wiki 作为持久化知识层，Issue 作为执行层。

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Wiki-Issue Integration Architecture                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   Input Layer                    Processing Layer         Storage Layer  │
│   ───────────                    ────────────────         ─────────────  │
│                                                                          │
│   ┌──────────┐                  ┌──────────────┐         ┌────────────┐ │
│   │ Terminal │─────────────────►│   Ingest     │────────►│ Wiki/      │ │
│   └──────────┘                  │   Pipeline   │         │ entities/  │ │
│                                  └──────────────┘         └────────────┘ │
│   ┌──────────┐                         │                 ┌────────────┐ │
│   │ Feishu   │─────────────────────────┤────────────────►│ issues/    │ │
│   └──────────┘                         │                 └────────────┘ │
│                                  ┌─────┴──────┐         ┌────────────┐ │
│   ┌──────────┐                  │   Agent    │────────►│ wiki/      │ │
│   │  GitHub  │─────────────────►│ Execution  │         │ executions/│ │
│   └──────────┘                  └─────┬──────┘         └────────────┘ │
│                                        │                                │
│                                        ▼                                │
│                                  ┌──────────────┐                       │
│                                  │   Archive    │                       │
│                                  └──────────────┘                       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Directory Structure

```
~/.capture/
├── issues/                     # Issues (可执行任务)
│   └── 2026/
│       └── 04/
│           ├── TASK-00001.md
│           └── TASK-00002.md
│
├── wiki/                       # Wiki (知识库)
│   ├── index.md               # 内容索引
│   ├── log.md                 # 时间线日志
│   ├── README.md              # Wiki 说明
│   │
│   ├── entities/              # 实体页面
│   │   ├── TASK-00001.md     # Issue 实体
│   │   ├── auth-system.md    # 系统实体
│   │   ├── jwt-middleware.md # 组件实体
│   │   └── ...
│   │
│   ├── concepts/              # 概念页面
│   │   ├── jwt.md
│   │   ├── authentication.md
│   │   ├── microservices.md
│   │   └── ...
│   │
│   ├── executions/            # 执行归档
│   │   ├── TASK-00001-2026-04-06.md
│   │   └── ...
│   │
│   ├── sources/               # 源文档摘要
│   │   └── ...
│   │
│   └── comparisons/           # 对比分析
│       └── ...
│
└── raw/                       # 原始资源 (Karpathy pattern)
    └── sources/              # 外部文档、论文等
        └── ...
```

---

## 4. Page Types

### 4.1 Entity Page (`wiki/entities/*.md`)

```markdown
---
type: entity
entity_type: issue|system|component|tech|person
created: 2026-04-06
updated: 2026-04-06
tags: [auth, security]
source: TASK-00001
---

# Auth System

## Description
用户认证系统，基于 JWT 实现...

## Related Entities
- [[entities/jwt-middleware]]
- [[entities/user-service]]

## Related Concepts
- [[concepts/jwt]]
- [[concepts/stateless-auth]]

## History
- [2026-04-06] Created from [[entities/TASK-00001]]
- [2026-04-06] Updated after [[executions/TASK-00001-2026-04-06]]
```

### 4.2 Concept Page (`wiki/concepts/*.md`)

```markdown
---
type: concept
domain: architecture|security|performance
created: 2026-04-06
updated: 2026-04-06
---

# JWT (JSON Web Token)

## Definition
JWT is an open standard for securely transmitting information...

## Applications
- [[entities/auth-system]]
- [[entities/TASK-00001]]

## Best Practices
- Keep tokens small
- Use HTTPS
- Set appropriate expiration

## See Also
- [[concepts/oauth2]]
- [[concepts/session-management]]
```

### 4.3 Execution Log (`wiki/executions/*.md`)

```markdown
---
type: execution
issue: TASK-00001
agent: claude@macbook-pro-001
status: completed
duration: 2700  # seconds
started: 2026-04-06T10:00:00Z
completed: 2026-04-06T10:45:00Z
---

# TASK-00001 Execution - 2026-04-06

## Summary
Successfully implemented JWT authentication middleware.

## Agent
- **Name**: Claude
- **Machine**: macbook-pro-001
- **Model**: claude-sonnet-4

## Generated Artifacts
- `auth/jwt.go` - JWT middleware implementation
- `auth/password.go` - bcrypt password hashing
- `docker-compose.yml` - Redis configuration

## Technical Decisions
1. Used `github.com/golang-jwt/jwt/v5` for JWT handling
2. Access token TTL: 15 minutes
3. Refresh token TTL: 7 days

## Output
```
[Agent output here...]
```

## Related
- [[entities/TASK-00001]]
- [[entities/jwt-middleware]]
- [[concepts/jwt]]
```

### 4.4 Index Page (`wiki/index.md`)

```markdown
# Wiki Index

Auto-generated index of all wiki pages.

## Entities

### Issues
- [TASK-00001: 实现用户认证系统](./entities/TASK-00001.md) - Status: completed
- [TASK-00002: 集成 OAuth2](./entities/TASK-00002.md) - Status: in_progress

### Systems
- [Auth System](./entities/auth-system.md) - Authentication and authorization
- [User Service](./entities/user-service.md) - User management

### Components
- [JWT Middleware](./entities/jwt-middleware.md) - JWT validation middleware

## Concepts

### Architecture
- [JWT](./concepts/jwt.md) - JSON Web Token
- [OAuth2](./concepts/oauth2.md) - OAuth 2.0 flow

### Security
- [Password Hashing](./concepts/password-hashing.md) - bcrypt best practices

## Recent Executions
- [TASK-00001 - 2026-04-06](./executions/TASK-00001-2026-04-06.md)

---
*Last updated: 2026-04-06T10:45:00Z*
```

### 4.5 Log Page (`wiki/log.md`)

```markdown
# Wiki Log

Chronological log of all wiki activities.

## [2026-04-06 10:45:00] execute | TASK-00001
- Agent: claude@macbook-pro-001
- Duration: 45min
- Result: success
- Output: [[executions/TASK-00001-2026-04-06]]
- Entities updated: [[entities/jwt-middleware]], [[entities/auth-system]]

## [2026-04-06 10:00:00] dispatch | TASK-00001
- Agent: claude@macbook-pro-001
- Runtime: rt-001

## [2026-04-06 09:30:00] ingest | TASK-00001
- Source: terminal
- Title: 实现用户认证系统
- Entities created: [[entities/TASK-00001]], [[entities/auth-system]]
- Concepts extracted: [[concepts/jwt]], [[concepts/authentication]]
```

---

## 5. Integration Workflows

### 5.1 Issue Ingestion Flow

```go
// When Issue is created or enters Analysis stage
func (w *WikiService) IngestIssue(issue *Issue) error {
    // 1. Create entity page for issue
    entityPage := &EntityPage{
        ID: issue.ID,
        Type: "issue",
        Title: issue.Title,
        Content: issue.Description,
        Source: issue.ID,
    }
    w.createEntityPage(entityPage)
    
    // 2. Extract concepts using LLM
    concepts := w.extractConcepts(issue)
    for _, concept := range concepts {
        w.upsertConceptPage(concept)
    }
    
    // 3. Create/update related entities
    entities := w.extractEntities(issue)
    for _, entity := range entities {
        w.upsertEntityPage(entity)
    }
    
    // 4. Update index
    w.updateIndex()
    
    // 5. Append to log
    w.appendLog(LogEntry{
        Time: time.Now(),
        Action: "ingest",
        Target: issue.ID,
        Details: map[string]any{
            "entities_created": len(entities),
            "concepts_extracted": len(concepts),
        },
    })
    
    // 6. Update issue with wiki refs
    issue.WikiRefs = WikiReferences{
        EntityPage: fmt.Sprintf("entities/%s", issue.ID),
        ConceptPages: conceptNames(concepts),
    }
    
    return nil
}
```

### 5.2 Execution Archive Flow

```go
// When Agent completes task
func (w *WikiService) ArchiveExecution(task *Task, result *TaskResult) error {
    // 1. Create execution log page
    execPage := &ExecutionPage{
        IssueID: task.IssueID,
        TaskID: task.ID,
        Agent: task.AgentID,
        Status: result.Status,
        Duration: result.Duration,
        Output: result.Output,
        Artifacts: result.Artifacts,
    }
    w.createExecutionPage(execPage)
    
    // 2. Extract generated components
    components := w.extractComponents(result)
    for _, comp := range components {
        w.upsertEntityPage(&EntityPage{
            Type: "component",
            Title: comp.Name,
            Source: task.IssueID,
        })
    }
    
    // 3. Update related entity pages
    w.updateEntityWithExecution(task.IssueID, execPage)
    
    // 4. Update index
    w.updateIndex()
    
    // 5. Append to log
    w.appendLog(LogEntry{
        Time: time.Now(),
        Action: "execute",
        Target: task.IssueID,
        Details: map[string]any{
            "agent": task.AgentID,
            "duration": result.Duration,
            "components": len(components),
        },
    })
    
    return nil
}
```

### 5.3 Knowledge Query Flow

```go
// Query wiki for knowledge
func (w *WikiService) Query(query string) (*QueryResult, error) {
    // 1. Search index for relevant pages
    pages := w.searchIndex(query)
    
    // 2. Read relevant entity/concept pages
    context := w.gatherContext(pages)
    
    // 3. Use LLM to synthesize answer
    answer := w.synthesize(query, context)
    
    // 4. Optionally save answer as new wiki page
    if answer.ShouldSave {
        w.createComparisonPage(answer)
    }
    
    return answer, nil
}
```

---

## 6. Wiki Service Interface

```go
type WikiService interface {
    // Ingestion
    IngestIssue(ctx context.Context, issue *Issue) error
    IngestSource(ctx context.Context, sourcePath string) error
    
    // Archival
    ArchiveExecution(ctx context.Context, task *Task, result *TaskResult) error
    
    // Query
    Query(ctx context.Context, query string) (*QueryResult, error)
    Search(ctx context.Context, keywords []string) ([]SearchResult, error)
    GetEntity(ctx context.Context, entityID string) (*EntityPage, error)
    GetConcept(ctx context.Context, conceptID string) (*ConceptPage, error)
    
    // Maintenance
    UpdateIndex(ctx context.Context) error
    Lint(ctx context.Context) (*LintReport, error)
    
    // Cross-references
    GetBacklinks(ctx context.Context, pageID string) ([]string, error)
    GetRelated(ctx context.Context, pageID string) ([]string, error)
}
```

---

## 7. Lint Checks

```go
type LintReport struct {
    Orphans        []string  // Pages with no inbound links
    BrokenLinks    []string  // Links to non-existent pages
    StaleClaims    []StaleClaim  // Claims contradicted by newer info
    MissingConcepts []string // Concepts mentioned but have no page
    EmptyPages     []string  // Pages with no content
}

// Automated checks
func (w *WikiService) Lint() (*LintReport, error) {
    // 1. Find orphan pages
    // 2. Find broken wiki-links
    // 3. Detect contradictions between pages
    // 4. Find mentioned but missing concepts
    // 5. Check for empty pages
}
```

---

## 8. Related Specs

- [Issue Model Spec](./issue-model-spec.md)
- [Agent Runtime Spec](./agent-runtime-spec.md)
- [Storage Spec](./storage-spec.md)
