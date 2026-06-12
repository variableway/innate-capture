# Issue 011: TUI Kanban Board

> **Type**: Feature  
> **Priority**: P1  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement interactive TUI Kanban board using Bubble Tea for visual issue management.

## Acceptance Criteria

- [ ] Kanban board displays issues in columns
- [ ] Columns: Todo, In Progress, Done
- [ ] Keyboard navigation works
- [ ] Issue details view works
- [ ] Status can be changed
- [ ] Real-time updates

## Tasks

### 1. Setup Bubble Tea
```go
import tea "github.com/charmbracelet/bubbletea"
```

### 2. Implement Board Model
```go
type Board struct {
    columns []Column
    cursor  Position
}

type Column struct {
    title string
    issues []Issue
}
```

### 3. Implement Views
- Board view with columns
- Issue detail view
- Help view

### 4. Implement Keybindings
```
h/l - Move between columns
j/k - Move between issues
Enter - View details
s - Change status
q - Quit
```

### 5. Add Real-time Updates
```go
func (b *Board) ListenForUpdates()
```

## Technical Notes

- Use lipgloss for styling
- Support mouse events
- Refresh on window resize

## Dependencies

- Issue 004: Basic CLI

## Related Specs

- None
