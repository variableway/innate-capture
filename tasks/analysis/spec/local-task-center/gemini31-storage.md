# Storage Module Spec (Gemini-3.1 设计)

## 目标

持久化任务，包括新加入的阶段与分派信息，同时提供高效索引查询。

## 存储机制

双写策略：
1. **Markdown (Source of Truth)**: 所有数据、描述、YAML Frontmatter 均写入文件
2. **SQLite (Index Layer)**: 所有核心字段、阶段、状态均写入 DB 以便 `list` 与看板快速查询

## 数据兼容

- 确保 SQLite schema 变更（如增加 `stage`、`assigned_agent` 等）能自动迁移
- 确保解析现有 Markdown 文件时，新字段有合理默认值（如 `stage=inbox`）
- 保证重新构建索引时（`RebuildIndex`）不会丢失新字段数据

## 文件路径约定

- 继续使用 `~/.capture/tasks/YYYY/MM/TASK-XXXXX.md` 作为主存储路径
