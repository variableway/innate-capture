# Issue 002: Markdown File Storage

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement file-based storage for Issues using Markdown with YAML frontmatter.

## Acceptance Criteria

- [ ] Issues saved to `~/.capture/issues/YYYY/MM/TASK-XXXXX.md`
- [ ] YAML frontmatter with all metadata
- [ ] Markdown body for description
- [ ] File creation works
- [ ] File reading works
- [ ] File updating works
- [ ] File deletion works
- [ ] List with filters works

## Tasks

### 1. Create Directory Structure
```
~/.capture/issues/2026/04/TASK-00001.md
```

### 2. Implement Issue Store Interface
```go
type IssueStore interface {
    Create(issue *Issue) error
    Get(id string) (*Issue, error)
    Update(issue *Issue) error
    Delete(id string) error
    List(filter IssueFilter) ([]*Issue, error)
}
```

### 3. Implement File Operations
- Create: Write YAML + Markdown to file
- Get: Parse file and return Issue
- Update: Rewrite file with new content
- Delete: Remove file
- List: Walk directory, parse matching files

### 4. Handle File Formats
```markdown
---
id: TASK-00001
title: "Example Issue"
status: todo
---

Description here...
```

### 5. Error Handling
- Handle missing files
- Handle parse errors gracefully
- Handle directory creation

## Technical Notes

- Use `os.MkdirAll` for directory creation
- Use file locking for concurrent access
- Store file path in Issue struct for quick access

## Dependencies

- Issue 001: Core Data Model

## Related Specs

- [Storage Spec](../prd/spec/storage-spec.md)
