# Storage Specification

> **Version**: 1.0  
> **Status**: Draft

---

## 1. Overview

采用 Local First 的存储策略，以 Markdown 文件为数据源，SQLite 为索引和缓存，Git 为版本控制。

---

## 2. Storage Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Storage Layer                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Source of Truth                             │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │    │
│  │  │   Issues/    │  │    Wiki/     │  │   Config Files       │  │    │
│  │  │  Markdown    │  │  Markdown    │  │   (YAML)             │  │    │
│  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │    │
│  │                                                                  │    │
│  │  - Human readable                                                │    │
│  │  - Git version controlled                                        │    │
│  │  - Portable                                                      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                  │                                       │
│                                  ▼                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Index Layer                                 │    │
│  │                    ┌──────────────┐                              │    │
│  │                    │    SQLite    │                              │    │
│  │                    │   (Cache)    │                              │    │
│  │                    └──────────────┘                              │    │
│  │                                                                  │    │
│  │  - Fast queries                                                  │    │
│  │  - Full-text search                                              │    │
│  │  - Auto-generated from Markdown                                  │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                  │                                       │
│                                  ▼                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Cloud Sync (Optional)                       │    │
│  │                    ┌──────────────┐                              │    │
│  │                    │ Feishu Bitable│                              │    │
│  │                    └──────────────┘                              │    │
│  │                                                                  │    │
│  │  - Collaboration view                                            │    │
│  │  - Mobile access                                                 │    │
│  │  - Notification trigger                                          │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Directory Structure

```
~/.capture/
├── config.yaml              # Main configuration
├── agents.yaml             # Agent definitions
├── capture.db              # SQLite database (auto-generated)
│
├── issues/                 # Issues storage
│   └── 2026/
│       └── 04/
│           ├── TASK-00001.md
│           └── TASK-00002.md
│
├── wiki/                   # Wiki knowledge base
│   ├── index.md
│   ├── log.md
│   ├── entities/
│   ├── concepts/
│   ├── executions/
│   └── sources/
│
└── raw/                    # Raw external sources
    └── sources/
```

---

## 4. SQLite Schema

### 4.1 Issues Table

```sql
CREATE TABLE issues (
    id TEXT PRIMARY KEY,
    number INTEGER NOT NULL,
    identifier TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL CHECK(status IN ('todo', 'in_progress', 'done', 'cancelled', 'archived')),
    stage TEXT NOT NULL CHECK(stage IN ('inbox', 'analysis', 'planning', 'dispatch', 'execution', 'review', 'complete')),
    priority TEXT NOT NULL CHECK(priority IN ('urgent', 'high', 'medium', 'low', 'none')),
    assignee_type TEXT CHECK(assignee_type IN ('agent', 'user')),
    assignee_id TEXT,
    parent_id TEXT REFERENCES issues(id),
    source TEXT NOT NULL,
    creator_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    due_date DATETIME,
    file_path TEXT NOT NULL,
    wiki_entity_page TEXT,
    
    FOREIGN KEY (parent_id) REFERENCES issues(id)
);

CREATE INDEX idx_issues_status ON issues(status);
CREATE INDEX idx_issues_stage ON issues(stage);
CREATE INDEX idx_issues_assignee ON issues(assignee_id);
CREATE INDEX idx_issues_created ON issues(created_at);
```

### 4.2 Sub-tasks Table

```sql
CREATE TABLE sub_tasks (
    parent_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    child_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    
    PRIMARY KEY (parent_id, child_id)
);
```

### 4.3 Labels Table

```sql
CREATE TABLE labels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    color TEXT
);

CREATE TABLE issue_labels (
    issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    label_id TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    
    PRIMARY KEY (issue_id, label_id)
);
```

### 4.4 Tasks Table (Agent Tasks)

```sql
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    agent_id TEXT NOT NULL,
    runtime_id TEXT,
    status TEXT NOT NULL CHECK(status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    priority INTEGER DEFAULT 0,
    prompt TEXT,
    work_dir TEXT,
    session_id TEXT,
    result_summary TEXT,
    error_message TEXT,
    created_at DATETIME NOT NULL,
    started_at DATETIME,
    completed_at DATETIME,
    duration_seconds INTEGER
);

CREATE INDEX idx_tasks_issue ON tasks(issue_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_agent ON tasks(agent_id);
```

### 4.5 Wiki Index Table

```sql
CREATE TABLE wiki_pages (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK(type IN ('entity', 'concept', 'execution', 'source', 'comparison')),
    title TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    content_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    word_count INTEGER,
    
    -- For entities
    entity_type TEXT CHECK(entity_type IN ('issue', 'system', 'component', 'tech', 'person')),
    source_issue_id TEXT REFERENCES issues(id),
    
    -- Full-text search
    content_fts TEXT
);

-- Full-text search index
CREATE VIRTUAL TABLE wiki_fts USING fts5(
    title,
    content_fts,
    content='wiki_pages',
    content_rowid='rowid'
);

CREATE INDEX idx_wiki_type ON wiki_pages(type);
CREATE INDEX idx_wiki_entity_type ON wiki_pages(entity_type);
```

### 4.6 Wiki Links Table

```sql
CREATE TABLE wiki_links (
    from_page TEXT NOT NULL REFERENCES wiki_pages(id) ON DELETE CASCADE,
    to_page TEXT NOT NULL,
    link_type TEXT DEFAULT 'reference', -- reference, execution, concept
    
    PRIMARY KEY (from_page, to_page)
);

CREATE INDEX idx_wiki_links_to ON wiki_links(to_page);
```

### 4.7 Sync State Table

```sql
CREATE TABLE sync_state (
    issue_id TEXT PRIMARY KEY REFERENCES issues(id) ON DELETE CASCADE,
    feishu_record_id TEXT,
    last_synced_at DATETIME,
    sync_version INTEGER DEFAULT 1,
    local_hash TEXT NOT NULL,
    remote_hash TEXT
);
```

---

## 5. Store Interface

```go
// Store is the unified storage interface
type Store interface {
    // Issue operations
    IssueStore
    
    // Task operations
    TaskStore
    
    // Wiki operations
    WikiStore
    
    // Sync operations
    SyncStore
    
    // Lifecycle
    Init() error
    Close() error
    
    // Maintenance
    RebuildIndex() error
    Vacuum() error
}

// IssueStore defines issue storage operations
type IssueStore interface {
    CreateIssue(issue *Issue) error
    GetIssue(id string) (*Issue, error)
    UpdateIssue(issue *Issue) error
    DeleteIssue(id string) error
    ListIssues(filter IssueFilter) ([]*Issue, error)
    
    // Relations
    AddSubTask(parentID, childID string) error
    RemoveSubTask(parentID, childID string) error
    GetSubTasks(parentID string) ([]*Issue, error)
    
    // Labels
    AddLabel(issueID, labelID string) error
    RemoveLabel(issueID, labelID string) error
}

// TaskStore defines task storage operations
type TaskStore interface {
    CreateTask(task *Task) error
    GetTask(id string) (*Task, error)
    UpdateTask(task *Task) error
    DeleteTask(id string) error
    ListTasks(filter TaskFilter) ([]*Task, error)
    
    // Queue operations
    PollPendingTask(agentID string) (*Task, error)
    ClaimTask(taskID, agentID string) error
    CompleteTask(taskID string, result TaskResult) error
    FailTask(taskID string, err error) error
}

// WikiStore defines wiki storage operations
type WikiStore interface {
    // Pages
    CreatePage(page *WikiPage) error
    GetPage(id string) (*WikiPage, error)
    UpdatePage(page *WikiPage) error
    DeletePage(id string) error
    ListPages(filter WikiFilter) ([]*WikiPage, error)
    
    // Search
    Search(query string) ([]WikiSearchResult, error)
    SearchByTag(tag string) ([]*WikiPage, error)
    
    // Links
    AddLink(fromID, toID, linkType string) error
    GetBacklinks(pageID string) ([]string, error)
    GetOutgoingLinks(pageID string) ([]string, error)
}

// SyncStore defines sync state operations
type SyncStore interface {
    GetSyncState(issueID string) (*SyncState, error)
    UpdateSyncState(state *SyncState) error
    ListUnsynced() ([]*SyncState, error)
    MarkSynced(issueID string, remoteHash string) error
}
```

---

## 6. Dual-Write Strategy

### 6.1 Write Flow

```
┌─────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Write  │────►│  Markdown File   │────►│   SQLite Index  │
│ Request │     │  (Source of      │     │  (Auto-update)  │
└─────────┘     │   Truth)         │     └─────────────────┘
                └──────────────────┘
```

### 6.2 Implementation

```go
type DualWriteStore struct {
    fs     FileSystemStore  // Markdown files
    db     SQLiteStore      // SQLite index
}

func (s *DualWriteStore) CreateIssue(issue *Issue) error {
    // 1. Write to Markdown (source of truth)
    if err := s.fs.CreateIssue(issue); err != nil {
        return err
    }
    
    // 2. Update SQLite index
    if err := s.db.CreateIssue(issue); err != nil {
        // Log error but don't fail - can rebuild index later
        log.Printf("Failed to update index: %v", err)
    }
    
    return nil
}

func (s *DualWriteStore) GetIssue(id string) (*Issue, error) {
    // Read from SQLite (fast)
    issue, err := s.db.GetIssue(id)
    if err == nil {
        return issue, nil
    }
    
    // Fallback to file system
    return s.fs.GetIssue(id)
}

func (s *DualWriteStore) RebuildIndex() error {
    // Rebuild SQLite from Markdown files
    issues, err := s.fs.ListIssues(IssueFilter{})
    if err != nil {
        return err
    }
    
    if err := s.db.ClearIssues(); err != nil {
        return err
    }
    
    for _, issue := range issues {
        if err := s.db.CreateIssue(issue); err != nil {
            log.Printf("Failed to index issue %s: %v", issue.ID, err)
        }
    }
    
    return nil
}
```

---

## 7. File System Store

```go
// FileSystemStore implements Store using Markdown files
type FileSystemStore struct {
    root string // ~/.capture
}

func (s *FileSystemStore) CreateIssue(issue *Issue) error {
    // Determine file path: issues/YYYY/MM/TASK-XXXXX.md
    year := issue.CreatedAt.Format("2006")
    month := issue.CreatedAt.Format("01")
    dir := filepath.Join(s.root, "issues", year, month)
    
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    filePath := filepath.Join(dir, issue.ID + ".md")
    
    // Serialize to Markdown with YAML frontmatter
    content, err := s.serializeIssue(issue)
    if err != nil {
        return err
    }
    
    return os.WriteFile(filePath, []byte(content), 0644)
}

func (s *FileSystemStore) GetIssue(id string) (*Issue, error) {
    // Find file by glob pattern
    pattern := filepath.Join(s.root, "issues", "*", "*", id + ".md")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return nil, err
    }
    if len(matches) == 0 {
        return nil, ErrNotFound
    }
    
    // Parse Markdown
    data, err := os.ReadFile(matches[0])
    if err != nil {
        return nil, err
    }
    
    return s.parseIssue(string(data), matches[0])
}

func (s *FileSystemStore) ListIssues(filter IssueFilter) ([]*Issue, error) {
    var issues []*Issue
    
    err := filepath.Walk(filepath.Join(s.root, "issues"), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !strings.HasSuffix(path, ".md") {
            return nil
        }
        
        data, err := os.ReadFile(path)
        if err != nil {
            return nil // Skip unreadable files
        }
        
        issue, err := s.parseIssue(string(data), path)
        if err != nil {
            return nil // Skip unparseable files
        }
        
        if s.matchesFilter(issue, filter) {
            issues = append(issues, issue)
        }
        
        return nil
    })
    
    return issues, err
}
```

---

## 8. Related Specs

- [Issue Model Spec](./issue-model-spec.md)
- [Agent Runtime Spec](./agent-runtime-spec.md)
- [Wiki Integration Spec](./wiki-integration-spec.md)
