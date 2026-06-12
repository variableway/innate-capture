# Issue 008: Issue to Wiki Ingestion

> **Type**: Feature  
> **Priority**: P1  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement automatic ingestion of Issues into Wiki when they are created or analyzed.

## Acceptance Criteria

- [ ] Issue creates entity page in wiki
- [ ] Concepts are extracted from Issue
- [ ] Related entities are linked
- [ ] Index is updated
- [ ] Log entry is added
- [ ] Issue stores wiki references

## Tasks

### 1. Implement Ingest Service
```go
type IngestService interface {
    IngestIssue(issue *Issue) error
}
```

### 2. Create Entity Page
```markdown
---
type: entity
entity_type: issue
source: TASK-00001
---
# TASK-00001: Title

## Description
...

## Related
- [[entities/related]]
```

### 3. Extract Concepts (LLM)
```go
func extractConcepts(issue *Issue) []Concept {
    // Use LLM to extract key concepts
}
```

### 4. Update Index
```go
func (s *IngestService) updateIndex()
```

### 5. Append to Log
```go
## [2026-04-06] ingest | TASK-00001
- Title: ...
- Entities: ...
- Concepts: ...
```

## Technical Notes

- Trigger on Issue creation and stage advancement
- Async processing to not block user
- Handle LLM errors gracefully

## Dependencies

- Issue 007: Wiki Structure

## Related Specs

- [Wiki Integration Spec](../prd/spec/wiki-integration-spec.md)
