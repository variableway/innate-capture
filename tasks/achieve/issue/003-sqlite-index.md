# Issue 003: SQLite Index Layer

> **Type**: Feature  
> **Priority**: P0  
> **Status**: Planned  
> **Assignee**: TBD  
> **Estimate**: 1 week

---

## Description

Implement SQLite database as fast index layer for queries, auto-generated from Markdown files.

## Acceptance Criteria

- [ ] SQLite schema defined
- [ ] Tables: issues, tasks, wiki_pages
- [ ] Full-text search enabled
- [ ] Auto-rebuild from files
- [ ] Dual-write strategy works
- [ ] Query performance < 100ms for 10k issues

## Tasks

### 1. Design Schema
```sql
CREATE TABLE issues (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT,
    stage TEXT,
    -- ...
);

CREATE VIRTUAL TABLE wiki_fts USING fts5(...);
```

### 2. Implement SQLite Store
```go
type SQLiteStore struct {
    db *sql.DB
}
```

### 3. Implement Dual-Write Store
```go
type DualWriteStore struct {
    fs FileSystemStore
    db SQLiteStore
}
```

### 4. Add Index Rebuild
- Walk all Markdown files
- Parse and insert into SQLite
- Handle errors gracefully

### 5. Add Full-Text Search
- FTS5 virtual table
- Search across titles and content
- Ranking by relevance

## Technical Notes

- Use `modernc.org/sqlite` for pure Go
- Schema migrations with version table
- WAL mode for better concurrency

## Dependencies

- Issue 001: Core Data Model
- Issue 002: Markdown Storage

## Related Specs

- [Storage Spec](../prd/spec/storage-spec.md)
