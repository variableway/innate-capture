# 任务分发架构设计

## 架构目标

构建一套基于飞书 Bot + 本地 Daemon 的多机任务分发体系，使 Capture 能够将任务下发到多台 worker 机器上执行，并收集执行结果。

## 总体架构

```
┌─────────────────────────────────────────────────────────────┐
│                     Master (Capture)                        │
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌─────────────┐ │
│  │ Task Center  │───▶│ Dispatcher   │───▶│ Feishu Bot  │ │
│  │   Core       │    │ (Assign)     │    │ (Outgoing)  │ │
│  └──────────────┘    └──────────────┘    └──────┬──────┘ │
│                                                    │         │
└────────────────────────────────────────────────────┼─────────┘
                                                     │
                    ┌────────────────────────────────┼────────┐
                    │                                │         │
                    ▼                                ▼         ▼
             ┌────────────┐               ┌────────────┐ ┌────────────┐
             │ Machine A  │               │ Machine B  │ │ Machine C  │
             │            │               │            │ │            │
             │ Feishu Bot │◀──消息───────│Feishu Bot │◀┘Feishu Bot │
             │            │               │            │ │            │
             │ Local      │               │ Local      │ │ Local      │
             │ Daemon     │               │ Daemon     │ │ Daemon     │
             │            │               │            │ │            │
             │ Executor   │               │ Executor   │ │ Executor   │
             └─────┬──────┘               └─────┬──────┘ └─────┬──────┘
                   │                              │              │
                   ▼                              ▼              ▼
             执行并回调                       执行并回调      执行并回调
```

## 核心组件

### 1. Master Side

#### Dispatcher

负责：
- 管理任务 → 机器的映射关系
- 将任务打包成消息，通过对应机器的 Bot 发送
- 维护各机器的在线状态

#### Feishu Bot (Master)

负责：
- 接收来自 Worker 机器的执行结果回调
- 将结果写入 Task Center，更新 stage 和 execution 信息

### 2. Worker Side

#### Feishu Bot (per machine)

每台机器一个独立的飞书应用 + Bot，作为该机器的唯一消息入口：
- 接收来自 Master 的任务指令
- 将指令转发给本地 Daemon

#### Local Daemon

常驻进程，负责：
- 监听飞书 Bot 收到的消息
- 解析任务指令（Task ID、命令、仓库路径）
- 调用 Executor 执行
- 将执行结果通过 Bot 回调给 Master

#### Executor

执行具体任务：
- 在指定 git repo 中执行 shell 命令
- 支持工作目录（worktree）切换
- 返回执行结果（stdout/stderr/exit code）

## 消息协议

### 任务下发消息

```json
{
  "type": "task_dispatch",
  "task_id": "TASK-00001",
  "machine_id": "machine-A",
  "command": "bash /workspace/repo/run.sh",
  "repo": "/workspace/repo",
  "worktree": "feature-xyz",
  "callback_url": "https://master/callback/execute"
}
```

### 结果回调消息

```json
{
  "type": "task_result",
  "task_id": "TASK-00001",
  "machine_id": "machine-A",
  "status": "success",
  "exit_code": 0,
  "stdout": "...",
  "stderr": "...",
  "duration_ms": 12345
}
```

## 数据模型扩展

在 Capture Task 模型基础上新增：

| 字段 | 说明 |
|------|------|
| dispatch.machine_id | 目标机器 ID |
| dispatch.bot_token | 该机器的 Bot 访问令牌 |
| execution.callback_url | 结果回调地址 |

## 部署架构

### 每台 Worker 需要

1. 一个独立的飞书应用（创建独立 Bot）
2. 部署 `capture-worker` Daemon
3. 配置机器 ID 和 Bot 凭证

### Master 需要

1. Capture 主程序
2. 每个 Worker Bot 的访问凭证
3. 回调接收服务（可以是 Webhook 或 WebSocket）

## 消息可靠性

- **至少一次**：Daemon 收到任务后应持久化，重复消息幂等处理
- **回调确认**：Master 下发任务后，如果在超时时间内未收到回调，触发重发
- **Daemon 心跳**：定期向 Master 报告机器在线状态
