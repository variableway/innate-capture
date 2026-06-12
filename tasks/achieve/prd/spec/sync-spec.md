# Sync Specification

> **Version**: 1.0  
> **Status**: Draft

---

## 1. Overview

定义 Issue 和 Feishu Bitable 之间的同步机制，实现本地优先的数据存储与云端协作视图的双向同步。

---

## 2. Sync Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Sync Architecture                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Local Side                          Cloud Side                          │
│  ─────────                           ─────────                           │
│                                                                          │
│  ┌──────────────┐                   ┌─────────────────┐                  │
│  │   Markdown   │                   │  Feishu Bitable │                  │
│  │   Issues     │◄────────────────►│  (多维表格)      │                  │
│  └──────┬───────┘    Sync Engine   └────────┬────────┘                  │
│         │                                   │                            │
│  ┌──────▼───────┐                   ┌───────▼─────────┐                  │
│  │    SQLite    │                   │   Feishu Bot    │                  │
│  │    Index     │                   │   (通知/交互)    │                  │
│  └──────────────┘                   └─────────────────┘                  │
│                                                                          │
│  Sync Strategy: Local First                                              │
│  - Local Markdown is source of truth                                     │
│  - Feishu is read-only view + notification trigger                       │
│  - Bidirectional sync for status updates                                 │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Feishu Bitable Schema

### 3.1 Table Fields

| Field | Type | Local Field | Description |
|-------|------|-------------|-------------|
| Issue ID | Text | `id` | TASK-XXXXX format |
| Title | Text | `title` | Issue title |
| Description | Text | `description` | Full description |
| Status | Single Select | `status` | todo, in_progress, done, cancelled, archived |
| Stage | Single Select | `stage` | inbox, analysis, planning, dispatch, execution, review, complete |
| Priority | Single Select | `priority` | urgent, high, medium, low, none |
| Assignee | Text | `assignee_id` | Agent or user ID |
| Source | Single Select | `source` | cli, tui, feishu_bot, github |
| Created At | DateTime | `created_at` | Creation timestamp |
| Updated At | DateTime | `updated_at` | Last update |
| Due Date | DateTime | `due_date` | Optional deadline |
| Wiki Link | URL | `wiki_refs.entity_page` | Link to wiki entity |
| Local Path | Text | `file_path` | Local file path |

### 3.2 Field Mapping

```go
var FeishuFieldMapping = map[string]string{
    "Issue ID":     "id",
    "Title":        "title",
    "Description":  "description",
    "Status":       "status",
    "Stage":        "stage",
    "Priority":     "priority",
    "Assignee":     "assignee_id",
    "Source":       "source",
    "Created At":   "created_at",
    "Updated At":   "updated_at",
    "Due Date":     "due_date",
    "Wiki Link":    "wiki_entity_page",
    "Local Path":   "file_path",
}

var FeishuStatusOptions = map[string]string{
    "todo":        "待办",
    "in_progress": "进行中",
    "done":        "已完成",
    "cancelled":   "已取消",
    "archived":    "已归档",
}

var FeishuStageOptions = map[string]string{
    "inbox":      "收件箱",
    "analysis":   "分析中",
    "planning":   "规划中",
    "dispatch":   "待分派",
    "execution":  "执行中",
    "review":     "审核中",
    "complete":   "已完成",
}

var FeishuPriorityOptions = map[string]string{
    "urgent": "紧急",
    "high":   "高",
    "medium": "中",
    "low":    "低",
    "none":   "无",
}
```

---

## 4. Sync Service Interface

```go
// SyncService defines sync operations
type SyncService interface {
    // One-time sync
    SyncIssue(ctx context.Context, issueID string) error
    SyncAll(ctx context.Context, filter SyncFilter) error
    
    // Continuous sync
    StartAutoSync(ctx context.Context, interval time.Duration) error
    StopAutoSync()
    
    // Bidirectional sync
    PullFromFeishu(ctx context.Context) error  // Get updates from Feishu
    PushToFeishu(ctx context.Context) error    // Send updates to Feishu
    
    // Conflict resolution
    ResolveConflict(ctx context.Context, issueID string, strategy ConflictStrategy) error
    
    // Status
    GetSyncStatus(ctx context.Context, issueID string) (*SyncStatus, error)
}

type SyncFilter struct {
    Since      *time.Time
    Status     []string
    Stage      []string
    UnsyncedOnly bool
}

type SyncStatus struct {
    IssueID        string
    FeishuRecordID string
    LastSyncedAt   *time.Time
    LocalHash      string
    RemoteHash     string
    InSync         bool
}

type ConflictStrategy string

const (
    ConflictLocalWins     ConflictStrategy = "local_wins"
    ConflictRemoteWins    ConflictStrategy = "remote_wins"
    ConflictTimestampWins ConflictStrategy = "timestamp_wins"
    ConflictManual        ConflictStrategy = "manual"
)
```

---

## 5. Sync Flows

### 5.1 Push to Feishu (Local → Cloud)

```go
func (s *SyncService) PushToFeishu(ctx context.Context) error {
    // 1. Find unsynced issues
    states, err := s.store.ListUnsynced()
    if err != nil {
        return err
    }
    
    for _, state := range states {
        issue, err := s.store.GetIssue(state.IssueID)
        if err != nil {
            log.Printf("Failed to get issue %s: %v", state.IssueID, err)
            continue
        }
        
        // 2. Calculate local hash
        localHash := calculateHash(issue)
        
        if state.FeishuRecordID == "" {
            // 3a. Create new record
            recordID, err := s.feishu.CreateRecord(issue)
            if err != nil {
                log.Printf("Failed to create Feishu record: %v", err)
                continue
            }
            
            state.FeishuRecordID = recordID
            state.LastSyncedAt = time.Now()
            state.LocalHash = localHash
            
        } else if state.LocalHash != localHash {
            // 3b. Update existing record
            if err := s.feishu.UpdateRecord(state.FeishuRecordID, issue); err != nil {
                log.Printf("Failed to update Feishu record: %v", err)
                continue
            }
            
            state.LastSyncedAt = time.Now()
            state.LocalHash = localHash
        }
        
        // 4. Update sync state
        if err := s.store.UpdateSyncState(state); err != nil {
            log.Printf("Failed to update sync state: %v", err)
        }
    }
    
    return nil
}
```

### 5.2 Pull from Feishu (Cloud → Local)

```go
func (s *SyncService) PullFromFeishu(ctx context.Context) error {
    // 1. Fetch all records from Feishu
    records, err := s.feishu.ListRecords()
    if err != nil {
        return err
    }
    
    for _, record := range records {
        issueID := record.Fields["Issue ID"].(string)
        
        // 2. Get local issue
        localIssue, err := s.store.GetIssue(issueID)
        if err != nil {
            if err == ErrNotFound {
                // Create new local issue from Feishu
                issue := s.convertFromFeishu(record)
                if err := s.store.CreateIssue(issue); err != nil {
                    log.Printf("Failed to create local issue: %v", err)
                }
            }
            continue
        }
        
        // 3. Check for conflicts
        syncState, _ := s.store.GetSyncState(issueID)
        remoteHash := calculateHashFromFeishu(record)
        
        if syncState != nil && syncState.LocalHash != remoteHash {
            // Conflict detected
            switch s.conflictStrategy {
            case ConflictRemoteWins:
                issue := s.convertFromFeishu(record)
                s.store.UpdateIssue(issue)
            case ConflictLocalWins:
                // Do nothing, keep local
            case ConflictTimestampWins:
                remoteUpdated := parseTime(record.Fields["Updated At"])
                if remoteUpdated.After(localIssue.UpdatedAt) {
                    issue := s.convertFromFeishu(record)
                    s.store.UpdateIssue(issue)
                }
            }
        }
    }
    
    return nil
}
```

---

## 6. Auto Sync

### 6.1 Configuration

```yaml
# config.yaml
sync:
  enabled: true
  mode: "push_only"  # push_only | pull_only | bidirectional
  interval: "5m"     # Auto sync interval
  
  feishu:
    app_token: "YOUR_BITABLE_APP_TOKEN"
    table_id: "YOUR_TABLE_ID"
    
  conflict_strategy: "timestamp_wins"  # local_wins | remote_wins | timestamp_wins
  
  # What to sync
  filters:
    statuses: ["todo", "in_progress", "done"]
    exclude_archived: true
```

### 6.2 Auto Sync Loop

```go
func (s *SyncService) StartAutoSync(ctx context.Context, interval time.Duration) error {
    ticker := time.NewTicker(interval)
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                ticker.Stop()
                return
            case <-ticker.C:
                switch s.config.Mode {
                case "push_only":
                    s.PushToFeishu(ctx)
                case "pull_only":
                    s.PullFromFeishu(ctx)
                case "bidirectional":
                    s.PullFromFeishu(ctx)
                    s.PushToFeishu(ctx)
                }
            }
        }
    }()
    
    return nil
}
```

---

## 7. Event-Based Sync

### 7.1 Triggers

```go
type SyncTrigger int

const (
    TriggerManual SyncTrigger = iota
    TriggerIssueCreate
    TriggerIssueUpdate
    TriggerIssueStatusChange
    TriggerStageAdvance
    TriggerAgentComplete
    TriggerSchedule
)

// Register event handlers
func (s *SyncService) RegisterTriggers(eventBus *EventBus) {
    eventBus.Subscribe(EventIssueCreated, func(e Event) {
        s.SyncIssue(context.Background(), e.IssueID)
    })
    
    eventBus.Subscribe(EventIssueUpdated, func(e Event) {
        s.SyncIssue(context.Background(), e.IssueID)
    })
    
    eventBus.Subscribe(EventTaskCompleted, func(e Event) {
        s.SyncIssue(context.Background(), e.IssueID)
    })
}
```

---

## 8. Feishu Bot Integration

### 8.1 Notification Flow

```
Local Issue Update
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Sync Service │────►│ Feishu Bitable│────►│ Feishu Bot   │
│              │     │ (Update row)  │     │ (Notify user)│
└──────────────┘     └──────────────┘     └──────────────┘
```

### 8.2 Notification Types

```go
type NotificationType string

const (
    NotifyTaskAssigned   NotificationType = "task_assigned"
    NotifyTaskStarted    NotificationType = "task_started"
    NotifyTaskCompleted  NotificationType = "task_completed"
    NotifyTaskFailed     NotificationType = "task_failed"
)

func (s *SyncService) sendNotification(issue *Issue, notifType NotificationType) {
    card := s.buildNotificationCard(issue, notifType)
    s.feishuBot.SendCard(issue.Context.ChatID, card)
}
```

---

## 9. Related Specs

- [Issue Model Spec](./issue-model-spec.md)
- [Storage Spec](./storage-spec.md)
