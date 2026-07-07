# Overview: innate-capture

**Slug:** `innate-capture`  
**Repo:** `innate-works/innate-capture/`  
**Updated:** 2026-07-06

## 一句话

Go CLI/TUI，作为 **innate-works 的日常入口**：向 workspace 写 idea、读每日清单。

## 当前阶段

| 项 | 说明 |
|---|---|
| 阶段 | building |
| 下一步 | 见 [index.md](./index.md) 当前 Sprint / Issue |

> Sprint 状态、分工、里程碑在 **issue** 维护（如 [cap-sprint-01](./issues/cap-sprint-01-workspace-mvp.md)），不在本文件堆叠。

## 关键约束

- `workspace.root` 指向 innate-works 根（含 `ideas/`、`daily/`）
- inbox / daily 事实来源在 innate-works，capture 不写第二套
- 飞书代码保留，启用与否由当时 Sprint issue 说明

## 禁止事项

- 不要新建 `.agents/` 或第二套 project 文档
- 不要未经 issue 扩大 spec
- 不要把 Sprint 进度写进本 overview（写 issue + task log）

## 技术栈（摘要）

| 层 | 技术 |
|---|---|
| CLI | Go, cobra, viper |
| TUI | kanban（bubbletea） |
| 对接 | innate-works 文件系统 |

## 里程碑（项目级，少改）

| 日期 | 事件 |
|---|---|
| 2026-07-06 | 初始化 `docs/project/` |
