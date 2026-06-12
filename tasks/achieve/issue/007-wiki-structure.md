# Issue 007: Wiki Directory Structure

> **Type**: Feature  
> **Priority**: P1  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 3 days

---

## Description

Create the Wiki directory structure and implement basic page management for knowledge accumulation.

## Acceptance Criteria

- [ ] Directory structure created
- [ ] `wiki/index.md` auto-generated
- [ ] `wiki/log.md` maintained
- [ ] Entity pages can be created
- [ ] Concept pages can be created
- [ ] Pages can be linked with wiki-links

## Tasks

### 1. Create Directory Structure
```
wiki/
├── index.md
├── log.md
├── entities/
├── concepts/
├── executions/
└── sources/
```

### 2. Implement Page Interface
```go
type WikiPage struct {
    ID      string
    Type    string // entity|concept|execution|source
    Title   string
    Content string
    Links   []string
}
```

### 3. Implement Index Generation
```go
func (w *WikiService) UpdateIndex() error
```

### 4. Implement Log Maintenance
```go
func (w *WikiService) AppendLog(entry LogEntry) error
```

### 5. Implement Wiki-Links
- Parse `[[entity/name]]` syntax
- Generate backlinks
- Validate links

## Technical Notes

- Use Obsidian-compatible wiki-link syntax
- Auto-generate index on every page update
- Log entries use consistent format for grep

## Dependencies

- Issue 002: Markdown Storage

## Related Specs

- [Wiki Integration Spec](../prd/spec/wiki-integration-spec.md)
