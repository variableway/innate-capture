# Catpure Ideas

> 捕捉灵感，管理任务 — 通过 Terminal 或飞书 Bot 快速记录想法，保存为 Markdown 文件并同步到飞书多维表格。

## 功能特性

- **CLI 命令行** — 快速创建、查看、编辑、删除任务
- **Task Center 工作流** — 支持从 inbox、mindstorm、analysis 到 dispatch、execution、review 的阶段化流转
- **TUI 看板** — 终端交互式看板界面，可视化任务状态
- **飞书 Bot** — 通过飞书机器人接收消息，自动创建任务
- **飞书多维表格同步** — 双向同步任务到飞书 Bitable
- **Markdown 存储** — 任务以 Markdown + YAML frontmatter 格式存储
- **Agent 分派记录** — 为任务记录 agent、model、repo、worktree、terminal 等执行上下文 - 暂时不是重点
- **双模式 Bot** — 支持 Webhook 和 WebSocket 两种连接模式

这个是最初的项目，当时现在应该不仅仅是这个了，还需要 ：
1.  Task Center 需要Project 开发的详细进度，这部分可以和GITHUB Project 直接相关，同时主要包括了： 
   - 项目的任务
   - 项目任务的需求，spec等文档保存和分发到不同的项目
   - Workspace的管理，也就是需要一个Workspace 来管理不同的产品项目
[orig-ideas](../../capture-ori.md)，这个是一开始的想法和一些实现，现在需要在细化一下