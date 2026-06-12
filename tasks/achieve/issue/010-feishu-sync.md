# Issue 010: Feishu Bitable Sync

> **Type**: Feature  
> **Priority**: P1  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement bidirectional sync between local Issues and Feishu Bitable for cloud collaboration.

## Acceptance Criteria

- [ ] Issues sync to Feishu Bitable
- [ ] Status changes sync back
- [ ] Auto-sync on changes
- [ ] Conflict resolution works
- [ ] Sync status tracked

## Tasks

### 1. Setup Feishu API Client
```go
type FeishuClient struct {
    appToken string
    tableID  string
}
```

### 2. Implement Push
```go
func (s *SyncService) PushToFeishu() error
```

### 3. Implement Pull
```go
func (s *SyncService) PullFromFeishu() error
```

### 4. Implement Auto-Sync
```go
func (s *SyncService) StartAutoSync(interval time.Duration)
```

### 5. Track Sync State
```sql
CREATE TABLE sync_state (
    issue_id TEXT PRIMARY KEY,
    feishu_record_id TEXT,
    last_synced_at DATETIME,
    local_hash TEXT
);
```

### 6. Conflict Resolution
```go
func (s *SyncService) ResolveConflict(issueID string, strategy ConflictStrategy)
```

## Technical Notes

- Local First: Markdown is source of truth
- Bitable is read-only view by default
- Bidirectional for status updates only

## Dependencies

- Issue 004: Basic CLI

## Related Specs

- [Sync Spec](../prd/spec/sync-spec.md)
