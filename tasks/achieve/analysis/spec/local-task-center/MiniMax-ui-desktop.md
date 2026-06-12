# Desktop UI Spec (MiniMax 设计)

## 角色定位

P2 级别本地增强界面，作为个人工作台的入口，提供原生系统整合。

## 模块要求

### 1. Quick Capture
- 全局快捷键唤起（如 `Cmd+Shift+Space`）的极简输入框
- 支持快速提交 idea 到 Inbox 阶段

### 2. Local Dashboard Shell
- 嵌套加载 Web 版本的 Dashboard
- 不需重复实现 Kanban 和筛选逻辑

### 3. Notification Center
- 集成系统原生通知（macOS/Windows）
- 提醒：Agent 执行完成、执行报错、长期阻塞任务

### 4. Embedded Terminal Entry
- 提供一键打开特定 repo + worktree 终端的快捷入口
- 帮助用户快速跳转到执行现场

## 实施前提

- 必须在 Local API 和 Web 版本稳定后实施
- 技术栈选型建议：Tauri 或 Electron
