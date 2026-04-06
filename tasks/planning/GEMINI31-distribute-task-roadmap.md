# 任务分发实施计划 (GEMINI31 设计)

## 目标

在 Capture 中实现基于飞书 Bot + 本地 Daemon 的多机任务分发能力。

## 实施原则

- **轻量优先**：先验证核心链路，不做过度设计
- **复用优先**：尽量复用 Capture 现有飞书 Bot 积累
- **逐步演进**：先单机验证，再多机扩展

## 实施阶段

### Phase 1：单机验证（MVP）

**目标**：验证 Master → Worker Bot → Daemon → 执行 → 回调这条链路的可行性。

#### 产出

- Worker Daemon 程序（Go）
- 单机消息收发验证

#### 任务

| ID | 任务 | 优先级 |
|----|------|--------|
| P1-01 | 设计消息协议（task_dispatch / task_result） | P0 |
| P1-02 | 实现 Worker Daemon：监听 Bot 消息 | P0 |
| P1-03 | 实现 Worker Daemon：本地 shell 执行器 | P0 |
| P1-04 | 实现 Worker Daemon：结果回调 Master | P0 |
| P1-05 | 在 Capture 中新增 `dispatch send <task_id> --machine <machine_id>` 命令 | P0 |
| P1-06 | 单机 E2E 验证：Capture 下发 → Worker 执行 → 结果回流 | P0 |

#### 交付标准

- 能在单机上通过 Bot 收到任务并执行 shell 命令
- 执行结果能正确回调到 Capture 并更新任务状态

---

### Phase 2：多机支持

**目标**：支持多台 worker 机器注册与统一调度。

#### 产出

- 多机注册机制
- 机器状态追踪
- CLI 多机调度命令

#### 任务

| ID | 任务 | 优先级 |
|----|------|--------|
| P2-01 | 设计机器注册表（machine_id → bot_credentials） | P0 |
| P2-02 | 实现 `capture worker register <machine_id> --bot-token <token>` | P0 |
| P2-03 | 实现 `capture worker list` 查看已注册机器 | P0 |
| P2-04 | 实现 `capture worker status <machine_id>` 查看机器在线状态 | P1 |
| P2-05 | Daemon 心跳机制：定期向 Master 报告在线状态 | P1 |
| P2-06 | `capture dispatch --machine <machine_id>` 发送任务到指定机器 | P0 |
| P2-07 | `capture dispatch broadcast` 向所有在线机器广播任务 | P1 |

#### 交付标准

- 可以注册多台 worker 机器
- 可以向指定机器发送任务
- 能追踪各机器在线状态

---

### Phase 3：执行增强

**目标**：增强执行层面的能力，支持更复杂的任务场景。

#### 任务

| ID | 任务 | 优先级 |
|----|------|--------|
| P3-01 | 支持 worktree 切换执行 | P1 |
| P3-02 | 支持带超时的命令执行 | P1 |
| P3-03 | 支持 stdout/stderr 实时流式回调（进度可见） | P1 |
| P3-04 | 实现任务取消机制（发送 cancel → Worker abort） | P2 |
| P3-05 | 支持 interactive session（通过 Bot 会话转发 terminal 输入） | P2 |

---

### Phase 4：调度增强

**目标**：在 CLI/TUI 中实现更智能的调度体验。

#### 任务

| ID | 任务 | 优先级 |
|----|------|--------|
| P4-01 | CLI 显示各机器当前执行中的任务 | P1 |
| P4-02 | TUI 新增 Machine View：按机器维度查看任务分布 | P1 |
| P4-03 | `capture dispatch --idle-machine` 自动选择空闲机器 | P2 |
| P4-04 | 支持任务优先级和机器亲和性（某些任务只在特定机器执行） | P2 |

---

## 技术选型

### Worker Daemon 实现

- **语言**：Go（与 Capture 主程序一致，便于复用）
- **飞书 SDK**：复用 Capture 现有 `github.com/larksuite/oapi-sdk-go/v3`
- **消息监听模式**：WebSocket 长连接（与 Capture Bot serve 一致）

### 消息协议

- **格式**：JSON
- **传输**：飞书 Bot 消息（文本消息内嵌 JSON）或 飞书 Card 消息
- **回调**：Worker → Master 通过飞书 Bot 消息或 HTTP POST

### 消息协议权衡

考虑到飞书 Bot 消息有长度限制（最大 4000 字符），建议：
- 小任务直接通过 Bot 文本消息下发
- 大任务或复杂结果通过 HTTP 回调

---

## 文件结构

```
capture/
├── cmd/
│   ├── worker.go          # Worker Daemon 入口
│   └── dispatch.go        # Master 下发命令
├── internal/
│   ├── dispatch/
│   │   ├── sender.go     # 向 Worker Bot 发送消息
│   │   └── callback.go   # 接收 Worker 回调
│   ├── worker/
│   │   ├── daemon.go     # Worker Daemon 主程序
│   │   ├── executor.go   # Shell 执行器
│   │   ├── heartbeat.go  # 心跳上报
│   │   └── registry.go   # 机器注册表
│   └── bot/
│       ├── worker_bot.go # Worker 侧 Bot 处理
│       └── master_bot.go  # Master 侧 Bot 处理
```

---

## 测试策略

| 阶段 | 测试方式 |
|------|----------|
| Phase 1 | 手动 E2E：单机组装验证 |
| Phase 2 | 单元测试：registry、dispatch、executor |
| Phase 3 | 集成测试：多机器并发下发 |
| Phase 4 | TUI 交互测试 |

---

## 风险与对策

| 风险 | 影响 | 对策 |
|------|------|------|
| 飞书 Bot 消息延迟 | 任务下发慢 | 增加 HTTP 回调作为备用通道 |
| Worker 网络中断 | 任务卡住 | 设置执行超时 + 超时重试 |
| Bot token 泄露 | 安全风险 | Bot token 只存 Worker 本地，不走 Master 存储 |
| 多机器并发回调 | Master 处理不过来 | 回调队列 + 限流 |
