# Issue 006: Agent Backend Implementations

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 2 weeks

---

## Description

Implement unified backend interface for Claude, Codex, and OpenCode agents.

## Acceptance Criteria

- [ ] Common Backend interface
- [ ] Claude backend works
- [ ] Codex backend works
- [ ] OpenCode backend works
- [ ] Streaming output supported
- [ ] Session resume supported
- [ ] Error handling works

## Tasks

### 1. Define Backend Interface
```go
type Backend interface {
    Execute(ctx context.Context, prompt string, opts ExecOptions) (*Session, error)
    Health() error
    Version() (string, error)
}
```

### 2. Implement Claude Backend
```go
type ClaudeBackend struct {
    executable string
    model      string
}

func (b *ClaudeBackend) Execute(...) (*Session, error)
```

### 3. Implement Codex Backend
```go
type CodexBackend struct {
    executable string
    model      string
}
```

### 4. Implement OpenCode Backend
```go
type OpenCodeBackend struct {
    executable string
    model      string
}
```

### 5. Implement Session Streaming
```go
type Session struct {
    Messages <-chan Message
    Result   <-chan Result
}
```

### 6. Add Backend Factory
```go
func NewBackend(agentType string, config Config) (Backend, error)
```

## Technical Notes

- Parse agent output for structured data
- Handle agent process lifecycle
- Support environment variable injection

## Dependencies

- Issue 005: Agent Daemon Framework

## Related Specs

- [Agent Runtime Spec](../prd/spec/agent-runtime-spec.md)
