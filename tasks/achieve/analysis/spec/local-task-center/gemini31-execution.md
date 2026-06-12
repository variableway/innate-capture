# Execution Tracking Module Spec (Gemini-3.1 设计)

## 目标

跟踪已分配任务的执行状态，并收集执行结果。

## 记录维度

- 状态：agent 执行状态（`running`, `success`, `failed`, `timeout`）
- 开始时间：执行会话开始的时间戳
- 结束时间：执行会话结束的时间戳
- 结果摘要：成功信息、失败原因、或输出文件的路径

## 会话绑定 (Future)

在后续的实现中，Execution 模块需要能绑定到一个真实的终端会话（例如 `tmux` 或 `zellij` 的 session ID），以便用户可以随时进入“案发现场”。

## 业务流转

- Agent 开始工作时，标记 stage = `execution`
- Agent 完成工作后，标记 stage = `review`，等待人类确认
