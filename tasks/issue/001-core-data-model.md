# Issue 001: Core Data Model Implementation

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement the core Issue data model based on the specification. This is the foundation of the entire system.

## Acceptance Criteria

- [ ] Issue struct with all fields defined
- [ ] Status enum with valid transitions
- [ ] Stage pipeline implemented
- [ ] Priority levels defined
- [ ] JSON/YAML serialization works
- [ ] Validation rules implemented

## Tasks

### 1. Define Issue Model
```go
type Issue struct {
    ID          string
    Number      int
    Title       string
    Description string
    Status      IssueStatus
    Stage       IssueStage
    Priority    Priority
    // ... more fields
}
```

### 2. Implement Status State Machine
- Define valid transitions
- Implement `CanTransition(from, to)` function
- Handle invalid transitions with proper errors

### 3. Implement Stage Pipeline
- Define stage sequence
- Implement stage advancement logic
- Support manual approval gates

### 4. Add Serialization
- YAML frontmatter support
- JSON marshaling/unmarshaling
- Custom field types (timestamps, enums)

### 5. Add Validation
- Required fields
- Field length limits
- Enum value validation

## Technical Notes

- Use string types for enums for readability
- Store timestamps as RFC3339
- Use pointers for optional fields
- Add JSON/YAML tags for all fields

## Dependencies

- None (foundation task)

## Related Specs

- [Issue Model Spec](../prd/spec/issue-model-spec.md)
