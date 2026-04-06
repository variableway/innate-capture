# Issue 012: Feishu Bot Integration

> **Type**: Feature  
> **Priority**: P1  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement Feishu Bot for creating and managing Issues via chat messages.

## Acceptance Criteria

- [ ] Bot receives messages via webhook
- [ ] `记录 <content>` creates issue
- [ ] `列出` shows issue list
- [ ] `删除 <id>` deletes issue
- [ ] `状态 <id>` shows status
- [ ] Notifications sent on task complete

## Tasks

### 1. Setup Webhook Handler
```go
func (h *BotHandler) HandleWebhook(w http.ResponseWriter, r *http.Request)
```

### 2. Implement Commands
```go
func (h *BotHandler) handleCreate(content string, user User) error
func (h *BotHandler) handleList(user User) error
func (h *BotHandler) handleDelete(id string, user User) error
func (h *BotHandler) handleStatus(id string, user User) error
```

### 3. Parse Commands
```
记录 实现用户认证系统
列出
删除 TASK-00001
状态 TASK-00001
```

### 4. Send Notifications
```go
func (h *BotHandler) NotifyTaskComplete(issue Issue)
```

### 5. WebSocket Support (optional)
```go
func (h *BotHandler) HandleWebSocket(conn *websocket.Conn)
```

## Technical Notes

- Verify Feishu signatures
- Support both webhook and WebSocket
- Store chat context in Issue

## Dependencies

- Issue 010: Feishu Sync

## Related Specs

- None
