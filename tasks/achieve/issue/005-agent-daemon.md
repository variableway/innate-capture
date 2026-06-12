# Issue 005: Agent Daemon Framework

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 2 weeks

---

## Description

Implement the local Agent Daemon that polls for tasks and executes them using configured agents.

## Acceptance Criteria

- [ ] Daemon starts and runs in background
- [ ] Polls task queue periodically
- [ ] Claims tasks atomically
- [ ] Executes tasks with configured agent
- [ ] Reports results back
- [ ] Health endpoint works
- [ ] Graceful shutdown

## Tasks

### 1. Daemon Lifecycle
```go
type Daemon struct {
    config Config
    agents map[string]Agent
    tasks  TaskQueue
}

func (d *Daemon) Run(ctx context.Context) error
func (d *Daemon) Stop() error
```

### 2. Implement Task Polling
```go
func (d *Daemon) pollLoop(ctx context.Context) {
    for {
        task, err := d.claimTask()
        if err != nil || task == nil {
            time.Sleep(d.config.PollInterval)
            continue
        }
        d.handleTask(task)
    }
}
```

### 3. Implement Task Execution
```go
func (d *Daemon) handleTask(task *Task) {
    agent := d.agents[task.AgentID]
    result := agent.Execute(task.Prompt)
    d.reportResult(task.ID, result)
}
```

### 4. Health Endpoint
```http
GET http://localhost:19514/health
{
    "status": "healthy",
    "agents": [...],
    "active_tasks": 0
}
```

### 5. CLI Commands
```bash
capture daemon start [--foreground]
capture daemon stop
capture daemon status
capture daemon logs
```

## Technical Notes

- Use context for graceful shutdown
- Implement semaphore for max concurrent tasks
- Log to file in background mode

## Dependencies

- Issue 003: SQLite Index (for task queue)

## Related Specs

- [Agent Runtime Spec](../prd/spec/agent-runtime-spec.md)
