# Assignment Module Spec (MiniMax 设计)

## 目标

把经过分析和拆解的任务，分配给对应的执行环境和 AI Agent。

## 记录维度

- Agent：指定执行当前任务的代理程序（如 `claude-code`, `codex`, `kimi`）
- Model：指定执行使用的模型（如 `gpt-4o`, `claude-3-opus`）
- Repository：指定任务应该在哪一个代码仓库中执行
- Worktree：指定代码仓库中的哪个工作区（可选）
- TerminalSession：分配的终端会话 ID（可选，用于重连）
- AssignedAt：分配时间戳

## 业务流转

- 一旦任务分配成功，该任务的 stage 应自动变更为 `dispatch`
- 下一个阶段是 `execution`，通常由 Agent 开始执行时主动触发
