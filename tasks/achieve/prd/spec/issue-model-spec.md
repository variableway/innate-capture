# Issue Data Model Specification

> **Version**: 1.0  
> **Status**: Draft  
> **Based on**: Linear SDK + Multica Analysis

---

## 1. Overview

Issue 是系统的核心工作单元，代表一个可执行的任务或需求。本规范定义 Issue 的数据结构、状态流转和生命周期。

---

## 2. Data Structure

### 2.1 Core Issue Model

```go
type Issue struct {
    // Identity
    ID          string    `yaml:"id" json:"id"`           // TASK-XXXXX format
    Number      int       `yaml:"number" json:"number"`   // Auto-increment per workspace
    Identifier  string    `yaml:"identifier" json:"identifier"` // PREFIX-NUMBER
    
    // Content
    Title       string    `yaml:"title" json:"title"`
    Description string    `yaml:"description" json:"description"`
    
    // Status & Stage
    Status      IssueStatus `yaml:"status" json:"status"`     // todo, in_progress, done, cancelled, archived
    Stage       IssueStage  `yaml:"stage" json:"stage"`       // Pipeline stage
    Priority    Priority    `yaml:"priority" json:"priority"` // urgent, high, medium, low, none
    
    // Assignment
    AssigneeType string    `yaml:"assignee_type,omitempty" json:"assignee_type"` // "agent" | "user"
    AssigneeID   string    `yaml:"assignee_id,omitempty" json:"assignee_id"`
    
    // Relations
    ParentID    string    `yaml:"parent_id,omitempty" json:"parent_id"`
    SubTasks    []string  `yaml:"sub_tasks,omitempty" json:"sub_tasks"`
    Labels      []string  `yaml:"labels,omitempty" json:"labels"`
    
    // Metadata
    Source      string    `yaml:"source" json:"source"`           // cli, tui, feishu_bot, github
    CreatorID   string    `yaml:"creator_id" json:"creator_id"`
    CreatedAt   time.Time `yaml:"created_at" json:"created_at"`
    UpdatedAt   time.Time `yaml:"updated_at" json:"updated_at"`
    DueDate     *time.Time `yaml:"due_date,omitempty" json:"due_date"`
    
    // Context
    Context     IssueContext `yaml:"context" json:"context"`
    
    // Execution tracking
    Dispatch    IssueDispatch    `yaml:"dispatch,omitempty" json:"dispatch"`
    Execution   IssueExecution   `yaml:"execution,omitempty" json:"execution"`
    Sync        IssueSync        `yaml:"sync,omitempty" json:"sync"`
    
    // Wiki integration (Karpathy pattern)
    WikiRefs    WikiReferences   `yaml:"wiki_refs,omitempty" json:"wiki_refs"`
}
```

### 2.2 Status Enum

```go
type IssueStatus string

const (
    StatusTodo       IssueStatus = "todo"
    StatusInProgress IssueStatus = "in_progress"
    StatusDone       IssueStatus = "done"
    StatusCancelled  IssueStatus = "cancelled"
    StatusArchived   IssueStatus = "archived"
)

// Valid transitions
var validTransitions = map[IssueStatus][]IssueStatus{
    StatusTodo:       {StatusInProgress, StatusDone, StatusCancelled},
    StatusInProgress: {StatusDone, StatusCancelled, StatusTodo},
    StatusDone:       {StatusArchived},
    StatusCancelled:  {StatusTodo, StatusArchived},
    StatusArchived:   {},
}
```

### 2.3 Stage Pipeline

```go
type IssueStage string

const (
    StageInbox     IssueStage = "inbox"      // Initial capture
    StageAnalysis  IssueStage = "analysis"   // AI analysis & decomposition
    StagePlanning  IssueStage = "planning"   // Sub-task planning
    StageDispatch  IssueStage = "dispatch"   // Agent assignment
    StageExecution IssueStage = "execution"  // In progress
    StageReview    IssueStage = "review"     // Result review
    StageComplete  IssueStage = "complete"   // Done
)

// Stage transitions (forward only with human approval)
var stagePipeline = []IssueStage{
    StageInbox, StageAnalysis, StagePlanning, StageDispatch, 
    StageExecution, StageReview, StageComplete,
}
```

### 2.4 Context & Tracking

```go
type IssueContext struct {
    Trigger   string `yaml:"trigger" json:"trigger"`       // What triggered creation
    Location  string `yaml:"location" json:"location"`     // Working directory
    RelatedTo string `yaml:"related_to" json:"related_to"` // Related issue/wiki entity
    ChatID    string `yaml:"chat_id,omitempty" json:"chat_id"`    // Feishu chat
    MessageID string `yaml:"message_id,omitempty" json:"message_id"` // Source message
}

type IssueDispatch struct {
    AgentID       string     `yaml:"agent_id,omitempty" json:"agent_id"`
    AgentName     string     `yaml:"agent_name,omitempty" json:"agent_name"`
    RuntimeID     string     `yaml:"runtime_id,omitempty" json:"runtime_id"`
    MachineID     string     `yaml:"machine_id,omitempty" json:"machine_id"`
    AssignedAt    *time.Time `yaml:"assigned_at,omitempty" json:"assigned_at"`
    Repository    string     `yaml:"repository,omitempty" json:"repository"`
    Worktree      string     `yaml:"worktree,omitempty" json:"worktree"`
}

type IssueExecution struct {
    SessionID     string     `yaml:"session_id,omitempty" json:"session_id"`
    WorkDir       string     `yaml:"work_dir,omitempty" json:"work_dir"`
    ResultSummary string     `yaml:"result_summary,omitempty" json:"result_summary"`
    BranchName    string     `yaml:"branch_name,omitempty" json:"branch_name"`
    CommitSHA     string     `yaml:"commit_sha,omitempty" json:"commit_sha"`
    StartedAt     *time.Time `yaml:"started_at,omitempty" json:"started_at"`
    CompletedAt   *time.Time `yaml:"completed_at,omitempty" json:"completed_at"`
    Duration      int        `yaml:"duration_seconds,omitempty" json:"duration_seconds"`
}

type IssueSync struct {
    FeishuRecordID string     `yaml:"feishu_record_id,omitempty" json:"feishu_record_id"`
    LastSyncedAt   *time.Time `yaml:"last_synced_at,omitempty" json:"last_synced_at"`
    SyncVersion    int        `yaml:"sync_version,omitempty" json:"sync_version"`
}
```

### 2.5 Wiki Integration

```go
type WikiReferences struct {
    EntityPage   string   `yaml:"entity_page,omitempty" json:"entity_page"`     // wiki/entities/TASK-XXX
    ConceptPages []string `yaml:"concept_pages,omitempty" json:"concept_pages"` // Related concepts
    ExecutionLog string   `yaml:"execution_log,omitempty" json:"execution_log"` // wiki/executions/...
}
```

---

## 3. File Format

### 3.1 Markdown Storage

```markdown
---
id: TASK-00001
number: 1
identifier: CAP-1
title: "实现用户认证系统"
status: todo
stage: analysis
priority: high
source: cli
creator_id: user-001
created_at: 2026-04-06T10:00:00Z
updated_at: 2026-04-06T10:05:00Z
context:
  trigger: "terminal_input"
  location: "/Users/patrick/projects/myapp"
  related_to: ""
assignee_type: "agent"
assignee_id: "agent-001"
wiki_refs:
  entity_page: "entities/TASK-00001"
  concept_pages: ["concepts/jwt", "concepts/authentication"]
---

# 实现用户认证系统

## 描述
实现基于 JWT 的用户认证系统...

## 分析结果

### 技术方案
- 使用 JWT for stateless auth
- bcrypt for password hashing
- Redis for token blacklist

### 子任务
- [ ] SUB-001: 设计数据库 Schema
- [ ] SUB-002: 实现 JWT 中间件
- [ ] SUB-003: 实现登录/注册 API

## 执行记录
- [Execution Log](../executions/TASK-00001-2026-04-06.md)
```

### 3.2 Storage Path

```
~/.capture/
└── issues/
    └── 2026/
        └── 04/
            ├── TASK-00001.md
            ├── TASK-00002.md
            └── ...
```

---

## 4. API Interface

```go
// IssueService defines issue operations
type IssueService interface {
    // CRUD
    Create(ctx context.Context, title string, opts CreateOptions) (*Issue, error)
    Get(ctx context.Context, id string) (*Issue, error)
    Update(ctx context.Context, id string, updates UpdateOptions) (*Issue, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter IssueFilter) ([]*Issue, error)
    
    // Status & Stage
    TransitionStatus(ctx context.Context, id string, to IssueStatus) error
    AdvanceStage(ctx context.Context, id string) error
    
    // Assignment
    Assign(ctx context.Context, id string, assigneeType, assigneeID string) error
    Unassign(ctx context.Context, id string) error
    
    // Relations
    AddSubTask(ctx context.Context, parentID, subTaskID string) error
    SetParent(ctx context.Context, childID, parentID string) error
    
    // Wiki integration
    LinkWikiEntity(ctx context.Context, id string, entityPath string) error
    
    // Sync
    SyncToFeishu(ctx context.Context, id string) error
}
```

---

## 5. Validation Rules

| Field | Rule |
|-------|------|
| ID | Format: `TASK-\d{5}`, unique |
| Title | Required, max 200 chars |
| Status | Must use validTransitions |
| Stage | Must use stagePipeline |
| Assignee | If assignee_type="agent", agent must exist |
| DueDate | If set, must be in future |

---

## 6. Related Specs

- [Agent Runtime Spec](./agent-runtime-spec.md)
- [Wiki Integration Spec](./wiki-integration-spec.md)
- [Storage Spec](./storage-spec.md)
