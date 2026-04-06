# Issue 004: Basic CLI Commands

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement basic CLI commands for Issue management using Cobra framework.

## Acceptance Criteria

- [ ] `capture add <title>` - Create new issue
- [ ] `capture list` - List all issues
- [ ] `capture show <id>` - Show issue details
- [ ] `capture edit <id>` - Edit issue
- [ ] `capture delete <id>` - Delete issue
- [ ] `capture status <id> <status>` - Update status
- [ ] Help text and flags work

## Tasks

### 1. Setup Cobra Framework
```go
var rootCmd = &cobra.Command{
    Use:   "capture",
    Short: "Capture ideas and manage tasks",
}
```

### 2. Implement Add Command
```bash
capture add "Implement user authentication"
  --priority=high
  --stage=inbox
```

### 3. Implement List Command
```bash
capture list
  --status=todo
  --stage=analysis
  --format=table|json
```

### 4. Implement Show Command
```bash
capture show TASK-00001
  --format=json
```

### 5. Implement Edit Command
```bash
capture edit TASK-00001
  --title="New title"
  --priority=low
```

### 6. Implement Delete Command
```bash
capture delete TASK-00001
  --force  # Skip confirmation
```

### 7. Implement Status Command
```bash
capture status TASK-00001 in_progress
```

## Technical Notes

- Use Viper for config management
- Support `--config` flag for custom config path
- Support `--data-dir` flag for custom data directory

## Dependencies

- Issue 001: Core Data Model
- Issue 002: Markdown Storage

## Related Specs

- None
