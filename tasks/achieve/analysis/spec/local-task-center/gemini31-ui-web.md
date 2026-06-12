# Web UI Spec (Gemini-3.1 设计)

## 角色定位

P1 级别观察与管理界面，负责 dashboard、聚合筛选、历史回顾，需要依赖稳定的 Local API。

## 模块要求

### 1. Dashboard Overview
- 活跃任务摘要：今日完成、阻塞中、执行中任务
- Stage 漏斗统计：展示有多少 idea 停留在 mindstorm，有多少在 execution

### 2. Stage Board
- 基于 Kanban 的多阶段可视化流转界面
- 支持拖拽变更 stage 或 status

### 3. Agent View
- 聚合视图：当前每个 agent 在处理什么任务
- 展示 agent 的运行状态与历史表现

### 4. Repository View
- 按工作区或代码仓库维度的任务视图
- 清晰展示该 repo 还有哪些关联的 TODO 未执行

### 5. Review Center
- 专属页面用于审批 `stage=review` 的任务
- 展示 Agent 执行结果与输出日志摘要

## 实施前提

- 必须先完成 Local API 层抽象
- 建议基于轻量级前端框架（如 React/Vue）配合 Tailwind 构建
