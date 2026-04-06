# Intake Module Spec (Gemini-3.1 设计)

## 目标

负责收集所有来源的输入，并将其转换为统一的 Task 对象。

## 输入源

1. CLI (`capture add`)
2. TUI (新建面板)
3. 飞书 Bot Webhook / WebSocket

## 数据映射

不管输入源是什么，最终必须包含：
- Title
- Description (可选)
- Tags (可选)
- Priority (默认 medium)

## 默认状态

所有新进入的任务：
- status 默认为 `todo`
- stage 默认为 `inbox`
