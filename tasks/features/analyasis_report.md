# 集成需求分析报告

## 1. Linear SDK Issue 模型分析

### 1.1 Issue 数据结构

Linear SDK (TypeScript) 中的 Issue 是一个核心实体，具有以下主要字段：

```typescript
// 核心字段
- id: string                    // 唯一标识符
- identifier: string             // 人类可读的标识符 (如 TEAM-123)
- title: string                  // 标题
- description?: string           // 描述（Markdown）
- state: WorkflowState          // 工作流状态
- priority: number               // 优先级 (0-4)
- assignee?: User               // 指派用户
- creator: User                 // 创建者
- labels: IssueLabel[]          // 标签
- createdAt: DateTime           // 创建时间
- updatedAt: DateTime           // 更新时间
- dueDate?: DateTime            // 截止日期
- estimate?: number             // 预估点数
- parent?: Issue                // 父Issue
- children: Issue[]             // 子Issues
- comments: Comment[]           // 评论
- cycle?: Cycle                 // 所属周期
- project?: Project             // 所属项目
- subscribers: User[]           // 订阅者
- url: string                   // URL链接
- branchName: string            // 建议分支名
```

### 1.2 Issue 状态流转

Linear 使用 **WorkflowState** 而非简单枚举：

```
backlog → unstarted → started → completed → canceled
```

每个 WorkflowState 包含：
- `name`: 状态名称
- `color`: 状态颜色
- `type`: 状态类型 (backlog, unstarted, started, completed, canceled)
- `position`: 排序位置

### 1.3 关键操作

```typescript
// Issue 创建
issueCreate(input: IssueCreateInput): Issue

// Issue 更新
issueUpdate(id: string, input: IssueUpdateInput): Issue

// Issue 删除
issueDelete(id: string): boolean

// 状态变更
issueUpdate(id: string, { stateId: newStateId }): Issue

// 分配变更
issueUpdate(id: string, { assigneeId: userId }): Issue
```

---

## 2. Multica Issue 管理与 Agent 授权分析

### 2.1 Issue 数据模型 (Go + PostgreSQL)

```go
type Issue struct {
    ID                 pgtype.UUID        // UUID 主键
    WorkspaceID        pgtype.UUID        // 工作空间ID
    Number             int32              // 递增编号
    Title              string             // 标题
    Description        pgtype.Text        // 描述
    Status             string             // 状态: backlog, todo, in_progress, done, cancelled
    Priority           string             // 优先级: urgent, high, medium, low, none
    AssigneeType       pgtype.Text        // 指派类型: "agent" | "user"
    AssigneeID         pgtype.UUID        // 指派ID
    CreatorType        string             // 创建者类型
    CreatorID          pgtype.UUID        // 创建者ID
    ParentIssueID      pgtype.UUID        // 父Issue
    AcceptanceCriteria []byte             // 验收标准 (JSON)
    ContextRefs        []byte             // 上下文引用 (JSON)
    Position           float64            // 排序位置
    DueDate            pgtype.Timestamptz // 截止日期
    CreatedAt          pgtype.Timestamptz // 创建时间
    UpdatedAt          pgtype.Timestamptz // 更新时间
}
```

### 2.2 Agent 模型

```go
type Agent struct {
    ID                 pgtype.UUID        // UUID
    WorkspaceID        pgtype.UUID        // 工作空间
    Name               string             // Agent名称
    RuntimeMode        string             // 运行模式: "local" | "cloud"
    RuntimeConfig      []byte             // 运行配置 (JSON)
    RuntimeID          pgtype.UUID        // 关联的运行时
    Visibility         string             // 可见性: "public" | "private"
    Status             string             // 状态: "idle" | "working"
    MaxConcurrentTasks int32              // 最大并发任务数
    OwnerID            pgtype.UUID        // 所有者
    Instructions       string             // Agent指令
    Tools              []byte             // 工具配置 (JSON)
    Triggers           []byte             // 触发器配置 (JSON)
}
```

### 2.3 任务队列模型 (AgentTaskQueue)

```go
type AgentTaskQueue struct {
    ID               pgtype.UUID        // 任务ID
    AgentID          pgtype.UUID        // 执行Agent
    IssueID          pgtype.UUID        // 关联Issue
    RuntimeID        pgtype.UUID        // 运行时ID
    Status           string             // 状态: pending, dispatched, running, completed, failed, cancelled
    Priority         int32              // 优先级
    Context          []byte             // 任务上下文 (JSON)
    SessionID        pgtype.Text        // Agent会话ID
    WorkDir          pgtype.Text        // 工作目录
    TriggerCommentID pgtype.UUID        // 触发评论
    DispatchedAt     pgtype.Timestamptz // 派发时间
    StartedAt        pgtype.Timestamptz // 开始时间
    CompletedAt      pgtype.Timestamptz // 完成时间
    Result           []byte             // 执行结果
    Error            pgtype.Text        // 错误信息
}
```

### 2.4 Issue → Agent 授权流程

```
┌─────────────┐     创建      ┌─────────────┐
│   User      │──────────────>│   Issue     │
└─────────────┘               └──────┬──────┘
                                     │
                                     │ 指派给Agent
                                     ▼
                              ┌─────────────┐
                              │    Agent    │
                              └──────┬──────┘
                                     │
                                     │ 触发条件检查
                              ┌──────┴──────┐
                              │  on_assign  │────┐
                              │  on_comment │────┤
                              │  on_mention │────┤
                              └─────────────┘    │
                                                 ▼
                                        ┌────────────────┐
                                        │ AgentTaskQueue │
                                        │   (pending)    │
                                        └───────┬────────┘
                                                │
                    ┌───────────────────────────┼───────────────────────────┐
                    │                           │                           │
                    ▼                           ▼                           ▼
            ┌──────────────┐          ┌──────────────┐          ┌──────────────┐
            │   Daemon     │          │   Daemon     │          │   Daemon     │
            │  (Machine A) │          │  (Machine B) │          │  (Machine C) │
            └──────────────┘          └──────────────┘          └──────────────┘
```

### 2.5 任务派发机制

```go
// 1. ClaimTask - 原子性领取任务
func (s *TaskService) ClaimTask(ctx context.Context, agentID pgtype.UUID) (*db.AgentTaskQueue, error) {
    // 检查并发限制
    running, _ := s.Queries.CountRunningTasks(ctx, agentID)
    if running >= int64(agent.MaxConcurrentTasks) {
        return nil, nil // 无容量
    }
    // 原子性领取
    task, _ := s.Queries.ClaimAgentTask(ctx, agentID)
}

// 2. 状态流转
pending → dispatched → running → completed/failed/cancelled
```

### 2.6 是否需要将分析任务拆解成 Task

**是的，需要多层抽象：**

```
┌─────────────────────────────────────────────────────────────┐
│                        Input (Raw)                          │
│         (Terminal 输入 / 飞书消息 / GitHub Issue)             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                        Issue (高层抽象)                       │
│  - 业务目标描述                                              │
│  - 验收标准                                                  │
│  - 优先级/截止日期                                           │
└───────────────────────┬─────────────────────────────────────┘
                        │
            ┌───────────┴───────────┐
            │                       │
            ▼                       ▼
┌──────────────────┐    ┌──────────────────┐
│   Analysis Task   │    │   Analysis Task   │
│  (问题分析拆解)    │    │  (方案设计规划)    │
└────────┬─────────┘    └────────┬─────────┘
         │                       │
         ▼                       ▼
┌──────────────────┐    ┌──────────────────┐
│   Sub-Task 1     │    │   Sub-Task 2     │
│  (具体实现A)      │    │  (具体实现B)      │
└────────┬─────────┘    └────────┬─────────┘
         │                       │
         └───────────┬───────────┘
                     │
                     ▼
┌──────────────────────────────────────────┐
│           Agent Task Queue               │
│     (可分配给不同机器并行执行)              │
└──────────────────────────────────────────┘
```

---

## 3. Terminal 输入管理方案

### 3.1 输入 → Issue 转化流程

```
Terminal Input
     │
     ▼
┌─────────────────────┐
│   capture add       │  ← CLI 命令
│   "实现用户认证系统"   │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│   Task Creation     │  ← 初始 Stage = inbox
│   - ID: TASK-00001  │
│   - Title: ...      │
│   - Source: cli     │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│   Auto Analysis     │  ← 本地 AI 快速分析
│   - 提取关键词        │
│   - 建议 Stage        │
│   - 预估复杂度        │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│   Issue Storage     │
│   - Markdown File   │
│   - SQLite Index    │
│   - Feishu Bitable  │
└─────────────────────┘
```

### 3.2 Issue 分析后本地保留内容

```yaml
# ~/.capture/tasks/2026/04/TASK-00001.md
---
id: TASK-00001
title: "实现用户认证系统"
stage: analysis
status: todo
priority: high
source: cli
created_at: 2026-04-06T10:00:00Z
updated_at: 2026-04-06T10:05:00Z
context:
  trigger: "terminal_input"
  location: "/Users/patrick/projects/myapp"
  related_to: ""
analysis:
  complexity: "medium"
  estimated_hours: 8
  suggested_agents: ["claude", "codex"]
  subtasks:
    - id: SUB-001
      title: "设计数据库Schema"
      agent: "claude"
    - id: SUB-002
      title: "实现JWT中间件"
      agent: "codex"
---

# 详细描述
实现基于JWT的用户认证系统...

## 分析结果

### 技术方案
- 使用 JWT for stateless auth
- bcrypt for password hashing
- Redis for token blacklist

### 子任务分解
...
```

---

## 4. 飞书消息输入管理方案

### 4.1 飞书消息 → Issue 转化流程

```
Feishu IM Message
       │
       ▼
┌─────────────────────────┐
│   Feishu Bot Webhook    │
│   / WebSocket           │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│   Message Parser        │
│   - 记录 <content>      │
│   - 列出                │
│   - 删除 <TASK-ID>      │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│   NLP Analysis          │  ← 可选：意图识别
│   - 是否为任务创建？      │
│   - 紧急程度？           │
│   - 是否需要回复？        │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│   Issue Creation        │
│   - Title: message      │
│   - Source: feishu_bot  │
│   - Reporter: user_id   │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│   Feishu Bitable Sync   │
│   + Local Markdown      │
└─────────────────────────┘
```

### 4.2 Issue 分析后处理

飞书消息创建的 Issue 需要额外处理：

```go
type FeishuIssueContext struct {
    ChatID      string    // 群聊ID
    ChatName    string    // 群聊名称
    SenderID    string    // 发送者ID
    SenderName  string    // 发送者名称
    MessageID   string    // 消息ID
    MessageType string    // 消息类型 (text/image/file)
    IsAtBot     bool      // 是否@了Bot
    ReplyTo     string    // 回复哪条消息
}
```

### 4.3 飞书特有的触发器

```go
// 飞书 Bot 触发器配置
type FeishuTriggers struct {
    OnDirectMessage bool // 私聊触发
    OnGroupAt       bool // 群内@触发
    OnKeyword       []string // 关键词触发
    AutoReply       bool // 是否自动回复
}
```

---

## 5. 场景打通方案

### 5.1 统一数据模型

```
┌─────────────────────────────────────────────────────────────────┐
│                        Unified Issue Model                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐    ┌─────────────────┐    ┌──────────────┐ │
│  │   Terminal      │    │   Feishu Bot    │    │  GitHub      │ │
│  │     Input       │    │    Message      │    │   Issues     │ │
│  └────────┬────────┘    └────────┬────────┘    └──────┬───────┘ │
│           │                      │                     │         │
│           └──────────────────────┼─────────────────────┘         │
│                                  ▼                               │
│                    ┌─────────────────────────┐                   │
│                    │      Issue Adapter      │                   │
│                    │  - normalizeInput()     │                   │
│                    │  - extractMetadata()    │                   │
│                    │  - createUnifiedIssue() │                   │
│                    └───────────┬─────────────┘                   │
│                                │                                 │
│                                ▼                                 │
│                    ┌─────────────────────────┐                   │
│                    │    Unified Issue        │                   │
│                    │  - id: ISS-XXXXX        │                   │
│                    │  - title                │                   │
│                    │  - description          │                   │
│                    │  - source_type          │                   │
│                    │  - source_metadata      │                   │
│                    │  - status               │                   │
│                    │  - stage                │                   │
│                    └───────────┬─────────────┘                   │
│                                │                                 │
└────────────────────────────────┼─────────────────────────────────┘
                                 │
                                 ▼
```

### 5.2 Issue → Task 表转化

```go
// 统一 Stage Pipeline
type StagePipeline struct {
    Inbox     string // 收件箱 - 原始输入
    Analysis  string // 分析 - AI分析拆解
    Planning  string // 规划 - 子任务规划
    Dispatch  string // 派发 - 分配给Agent
    Execution string // 执行 - Agent运行
    Review    string // 审核 - 结果检查
    Complete  string // 完成
}

// Issue 到 Task 的转化
type IssueToTaskConverter struct {
    // 当 Issue 进入 Dispatch Stage 时
    func Convert(issue Issue) []Task {
        // 1. 读取 issue.analysis.subtasks
        // 2. 为每个 subtask 创建 Task
        // 3. 设置依赖关系
        // 4. 返回 Task 列表
    }
}
```

### 5.3 飞书多维表格与本地 Multica 看板同步

```
┌──────────────────┐         ┌──────────────────┐         ┌──────────────────┐
│  Local Markdown  │◄───────►│   Sync Engine    │◄───────►│  Feishu Bitable  │
│   (Source of     │         │                  │         │   (云端备份/     │
│    Truth)        │         │ - 双向同步        │         │    协作视图)     │
└──────────────────┘         │ - 冲突解决        │         └──────────────────┘
                             │ - 增量更新        │                  ▲
┌──────────────────┐         └──────────────────┘                  │
│   SQLite Cache   │◄───────────────────────────────────────────────┘
│   (快速查询)      │         Webhook / Polling
└──────────────────┘
       ▲
       │
┌──────────────────┐
│   Multica TUI    │
│   (本地看板)      │
│  - Kanban View   │
│  - Agent Status  │
│  - Task Monitor  │
└──────────────────┘
```

**同步策略：**

```go
type SyncStrategy struct {
    // 本地优先，飞书为备份和协作视图
    LocalFirst bool
    
    // 同步触发时机
    Triggers struct {
        OnTaskCreate    bool
        OnTaskUpdate    bool
        OnStatusChange  bool
        OnAgentComplete bool
    }
    
    // 冲突解决策略
    ConflictResolution string // "local_wins" | "timestamp_wins" | "manual"
}
```

### 5.4 Agent 运行时配置与并行执行

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Agent Runtime Configuration                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    User Configuration                        │    │
│  │  ~/.capture/agents.yaml                                      │    │
│  │  ─────────────────────                                       │    │
│  │  agents:                                                     │    │
│  │    claude-macbook:                                           │    │
│  │      machine: macbook-pro                                    │    │
│  │      type: claude                                            │    │
│  │      model: claude-sonnet-4                                  │    │
│  │      max_concurrent: 3                                       │    │
│  │      work_dir: ~/.capture/workspaces/{workspace}             │    │
│  │                                                              │    │
│  │    codex-desktop:                                            │    │
│  │      machine: windows-desktop                                │    │
│  │      type: codex                                             │    │
│  │      model: gpt-4o                                           │    │
    │      max_concurrent: 2                                       │    │
│  │                                                              │    │
│  │    opencode-server:                                          │    │
│  │      machine: linux-server                                   │    │
│  │      type: opencode                                          │    │
│  │      model: gemini-2.5-pro                                   │    │
│  │      max_concurrent: 5                                       │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│                              │                                       │
│                              ▼                                       │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                 Task Dispatch Engine                         │    │
│  │                                                              │    │
│  │  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐        │    │
│  │  │  Scheduler  │   │  Scheduler  │   │  Scheduler  │        │    │
│  │  │  (MacBook)  │   │  (Windows)  │   │  (Linux)    │        │    │
│  │  └──────┬──────┘   └──────┬──────┘   └──────┬──────┘        │    │
│  │         │                 │                 │               │    │
│  │    ┌────┴────┐       ┌────┴────┐       ┌────┴────┐          │    │
│  │    │ Daemon  │       │ Daemon  │       │ Daemon  │          │    │
│  │    │ Process │       │ Process │       │ Process │          │    │
│  │    └────┬────┘       └────┬────┘       └────┬────┘          │    │
│  │         │                 │                 │               │    │
│  │         ▼                 ▼                 ▼               │    │
│  │    ┌─────────┐       ┌─────────┐       ┌─────────┐          │    │
│  │    │ Agent   │       │ Agent   │       │ Agent   │          │    │
│  │    │ Pool    │       │ Pool    │       │ Pool    │          │    │
│  │    │ [ ] [ ] │       │ [ ] [ ] │       │ [ ] [ ] │          │    │
│  │    │ [ ]     │       │ [ ]     │       │ [ ] [ ] │          │    │
│  │    └─────────┘       └─────────┘       └─────────┘          │    │
│  │                                                              │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.5 Agent 运行结果体现

```
┌────────────────────────────────────────────────────────────────────┐
│                     Agent 运行结果反馈                               │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. 本地文件系统                                                     │
│     ├── ~/.capture/tasks/YYYY/MM/TASK-XXXXX.md                      │
│     │   └── 更新 execution 字段                                      │
│     ├── ~/.capture/workspaces/{ws}/issues/{id}/                     │
│     │   ├── output.log                                              │
│     │   ├── git_changes.patch                                       │
│     │   └── artifacts/                                              │
│     └── ~/.capture/agents/{agent}/sessions/{session}/               │
│         └── transcript.json                                         │
│                                                                     │
│  2. Git 提交                                                         │
│     └── 在指定 repo 中创建 commit/branch                            │
│                                                                     │
│  3. 飞书通知                                                         │
│     ├── 任务完成卡片                                                 │
│     │   ├── 执行摘要                                                 │
│     │   ├── 关键输出                                                 │
│     │   └── 点击查看详情按钮                                          │
│     └── 异常告警                                                     │
│                                                                     │
│  4. 实时状态 (WebSocket)                                             │
│     ├── task:dispatch                                               │
│     ├── task:progress                                               │
│     ├── task:completed                                              │
│     └── task:failed                                                 │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

### 5.6 飞书客户端获取 Agent Daemon 状态

```
┌──────────────────────────────────────────────────────────────────────┐
│               Feishu Client Agent Status Monitor                      │
├──────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  方案1: 飞书 Bot 状态查询                                              │
│  ─────────────────────────                                            │
│  User ──► Feishu App ──► Bot Message "状态"                          │
│                            │                                          │
│                            ▼                                          │
│                    ┌───────────────┐                                  │
│                    │ Status Aggregator│                               │
│                    │ - 查询所有注册的 Daemon                          │
│                    │ - 汇总任务状态                                   │
│                    └───────┬───────┘                                  │
│                            │                                          │
│                            ▼                                          │
│                    ┌───────────────┐                                  │
│                    │  Feishu Card   │                                  │
│                    │  ┌─────────┐  │                                  │
│                    │  │机器状态 │  │                                  │
│                    │  │🟢 MacBook│  │                                  │
│                    │  │🟡 WinPC  │  │                                  │
│                    │  │🔴 Server │  │                                  │
│                    │  │         │  │                                  │
│                    │  │任务队列 │  │                                  │
│                    │  │运行: 3  │  │                                  │
│                    │  │等待: 5  │  │                                  │
│                    │  └─────────┘  │                                  │
│                    └───────────────┘                                  │
│                                                                       │
│  方案2: 心跳检测 (Daemon Health Check)                                │
│  ─────────────────────────────────────                                │
│                                                                       │
│  每个 Daemon 定期发送心跳到中央服务：                                   │
│  ┌─────────────┐      ┌──────────────┐      ┌──────────────────┐     │
│  │   Daemon    │─────►│  Server/API  │◄─────│  Feishu Bot      │     │
│  │ (各机器)     │ 心跳  │  (状态聚合)   │ 查询 │  (用户查询入口)   │     │
│  └─────────────┘      └──────────────┘      └──────────────────┘     │
│       every 30s                                                        │
│                                                                       │
│  心跳数据：                                                             │
│  {                                                                     │
│    daemon_id: "macbook-pro-001",                                       │
│    runtime_id: "rt-xxx",                                               │
│    status: "online",    // online | busy | offline                     │
│    tasks: {                                                            │
│      running: 2,                                                       │
│      max_concurrent: 3                                                 │
│    },                                                                  │
│    system: {                                                           │
│      cpu: "45%",                                                       │
│      memory: "60%"                                                     │
│    },                                                                  │
│    last_seen: "2026-04-06T10:30:00Z"                                   │
│  }                                                                     │
│                                                                       │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 6. 可行性分析

### 6.1 技术可行性

| 组件 | 技术选型 | 可行性 | 说明 |
|------|---------|--------|------|
| Issue 存储 | Markdown + SQLite + Feishu Bitable | ✅ 高 | 已有成熟方案 |
| Agent 执行 | Claude/Codex/OpenCode CLI | ✅ 高 | multica 已验证 |
| 任务队列 | 本地 SQLite / PostgreSQL | ✅ 高 | 单机可用SQLite |
| 多端同步 | Feishu Bitable API | ✅ 高 | 飞书提供完善API |
| Daemon 管理 | Go HTTP + WebSocket | ✅ 高 | 已有参考实现 |
| 本地看板 | Bubble Tea (Go TUI) | ✅ 高 | 已有基础实现 |

### 6.2 架构设计

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              Capture System                                │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                          Input Layer                                 │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐     │  │
│  │  │   CLI      │  │   TUI      │  │ Feishu Bot │  │  GitHub    │     │  │
│  │  │  (capture) │  │  (kanban)  │  │  (webhook) │  │  (webhook) │     │  │
│  │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘     │  │
│  │        └───────────────┴───────────────┴───────────────┘             │  │
│  │                              │                                        │  │
│  │                              ▼                                        │  │
│  │                   ┌─────────────────────┐                             │  │
│  │                   │   Issue Adapter     │                             │  │
│  │                   │  (normalize input)  │                             │  │
│  │                   └──────────┬──────────┘                             │  │
│  └──────────────────────────────┼───────────────────────────────────────┘  │
│                                 │                                          │
│  ┌──────────────────────────────┼───────────────────────────────────────┐  │
│  │                         Core Layer                                     │  │
│  │                              │                                        │  │
│  │  ┌───────────────────────────┴────────────────────────────────────┐  │  │
│  │  │                    Issue Service                                │  │  │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌───────────┐ │  │  │
│  │  │  │  Create    │  │  Analyze   │  │   Plan     │  │  Dispatch │ │  │  │
│  │  │  └────────────┘  └────────────┘  └────────────┘  └───────────┘ │  │  │
│  │  └───────────────────────────┬────────────────────────────────────┘  │  │
│  │                              │                                        │  │
│  │  ┌───────────────────────────┴────────────────────────────────────┐  │  │
│  │  │                    Storage Layer                                │  │  │
│  │  │  ┌─────────────────┐  ┌───────────────┐  ┌──────────────────┐ │  │  │
│  │  │  │ Markdown Store  │  │ SQLite Cache  │  │ Feishu Bitable   │ │  │  │
│  │  │  │ (Source of Truth)│  │ (Fast Query)  │  │ (Cloud Sync)     │ │  │  │
│  │  │  └─────────────────┘  └───────────────┘  └──────────────────┘ │  │  │
│  │  └────────────────────────────────────────────────────────────────┘  │  │
│  │                              │                                        │  │
│  └──────────────────────────────┼───────────────────────────────────────┘  │
│                                 │                                          │
│  ┌──────────────────────────────┼───────────────────────────────────────┐  │
│  │                      Execution Layer                                   │  │
│  │                              │                                        │  │
│  │  ┌───────────────────────────┴────────────────────────────────────┐  │  │
│  │  │                  Task Queue Service                             │  │  │
│  │  │                      │                                         │  │  │
│  │  │          ┌───────────┼───────────┐                             │  │  │
│  │  │          ▼           ▼           ▼                             │  │  │
│  │  │    ┌─────────┐  ┌─────────┐  ┌─────────┐                       │  │  │
│  │  │    │ Daemon  │  │ Daemon  │  │ Daemon  │    (Multi-Machine)    │  │  │
│  │  │    │(MacBook)│  │(Windows)│  │ (Linux) │                       │  │  │
│  │  │    └────┬────┘  └────┬────┘  └────┬────┘                       │  │  │
│  │  │         │            │            │                            │  │  │
│  │  │    ┌────┴────┐  ┌────┴────┐  ┌────┴────┐                       │  │  │
│  │  │    │  Agent  │  │  Agent  │  │  Agent  │                       │  │  │
│  │  │    │ (claude)│  │ (codex) │  │(opencode)│                      │  │  │
│  │  │    └─────────┘  └─────────┘  └─────────┘                       │  │  │
│  │  └────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                     │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                      Notification Layer                              │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐     │  │
│  │  │   TUI      │  │  Feishu    │  │  GitHub    │  │  Webhook   │     │  │
│  │  │  (本地看板) │  │  (消息推送) │  │  (PR/Comment)│  │  (外部集成) │     │  │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘     │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                            │
└──────────────────────────────────────────────────────────────────────────┘
```

### 6.3 实现计划

#### Phase 1: 基础框架 (4-6 周)

| 周次 | 任务 | 产出 |
|------|------|------|
| 1 | 重构 Issue 模型，统一数据层 | Issue CRUD API |
| 2 | 实现 Stage Pipeline | Stage 流转逻辑 |
| 3 | 集成飞书 Bitable 同步 | 双向同步功能 |
| 4 | 完善 TUI Kanban 看板 | 可视化看板 |
| 5-6 | 测试与优化 | Alpha 版本 |

#### Phase 2: Agent 集成 (4-6 周)

| 周次 | 任务 | 产出 |
|------|------|------|
| 1 | Agent 配置管理 | agents.yaml 配置 |
| 2 | 本地 Daemon 框架 | daemon 启动/管理 |
| 3 | 任务队列实现 | Task Queue Service |
| 4 | Agent 执行器 | Claude/Codex 集成 |
| 5 | 多机器支持 | 跨机器任务分发 |
| 6 | 测试与优化 | Beta 版本 |

#### Phase 3: 场景完善 (2-4 周)

| 周次 | 任务 | 产出 |
|------|------|------|
| 1 | 飞书 Bot 增强 | 状态查询、通知 |
| 2 | 结果反馈机制 | 多通道结果推送 |
| 3 | 监控与日志 | 运行监控 |
| 4 | 文档与发布 | v1.0 发布 |

### 6.4 时间表

```
2026 Q2 (4-6月)
├── 4月: Phase 1 - 基础框架
│   ├── Week 1-2: Issue 模型 + Stage Pipeline
│   ├── Week 3-4: 飞书同步 + TUI 看板
│   └── Week 5-6: 测试优化
│
├── 5月: Phase 2 - Agent 集成
│   ├── Week 1-2: Agent 配置 + Daemon 框架
│   ├── Week 3-4: 任务队列 + Agent 执行器
│   └── Week 5-6: 多机器支持
│
└── 6月: Phase 3 - 场景完善 + 发布
    ├── Week 1-2: 飞书 Bot 增强 + 结果反馈
    └── Week 3-4: 监控 + 文档 + v1.0
```

### 6.5 技术栈总结

| 层级 | 技术 | 说明 |
|------|------|------|
| 语言 | Go 1.26+ | 统一技术栈，无 CGo |
| CLI 框架 | Cobra | 命令行接口 |
| TUI 框架 | Bubble Tea | 交互式看板 |
| 存储 | Markdown + SQLite | 本地优先 |
| 云端同步 | Feishu Bitable API | 飞书多维表格 |
| Agent | Claude/Codex/OpenCode CLI | 多厂商支持 |
| 通讯 | HTTP/REST + WebSocket | Daemon 通信 |
| 配置 | YAML + Viper | 配置管理 |

### 6.6 Local First 策略

```
┌─────────────────────────────────────────────────────────────┐
│                     Local First Architecture                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 数据所有权                                                │
│     - 所有数据首先存储在本地 Markdown 文件                    │
│     - 用户完全拥有数据，可随时导出/迁移                       │
│                                                              │
│  2. 离线优先                                                  │
│     - 所有核心功能无需网络即可工作                            │
│     - 飞书同步为可选功能                                      │
│                                                              │
│  3. 云增强 (Cloud Enhancement)                                │
│     - 飞书: 协作视图 + 移动端访问 + 通知推送                  │
│     - Agent: 使用云端 AI 服务 (OpenClaw API)                 │
│                                                              │
│  4. 局域网优先                                                │
│     - Daemon 首先通过局域网发现彼此                          │
│     - 可选中心服务器用于公网穿透                              │
│                                                              │
│  5. 隐私保护                                                  │
│     - 敏感代码仅在本地处理                                    │
│     - Agent 执行环境隔离                                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. 任务分解

### 7.1 核心任务列表

```markdown
## P0 (Must Have)

- [ ] 重构 Issue 数据模型，兼容 Linear 风格
- [ ] 实现 Stage Pipeline (inbox → analysis → planning → dispatch → execution → review)
- [ ] Markdown + SQLite 双写存储
- [ ] Feishu Bitable 双向同步
- [ ] 增强 TUI Kanban 看板
- [ ] Agent 配置管理 (agents.yaml)
- [ ] 本地 Daemon 基础框架
- [ ] 任务队列实现
- [ ] Claude/Codex Agent 执行器

## P1 (Should Have)

- [ ] 多机器 Daemon 支持
- [ ] 飞书 Bot 状态查询
- [ ] Agent 执行结果反馈 (多渠道)
- [ ] 任务依赖管理
- [ ] Git 集成 (自动 commit/branch)
- [ ] 运行时监控面板

## P2 (Nice to Have)

- [ ] GitHub Issue 集成
- [ ] 自定义 Agent Skill
- [ ] 任务模板
- [ ] 执行报告生成
- [ ] 团队协作功能
```

---

## 8. 总结

### 8.1 关键设计决策

1. **Issue 作为高层抽象**: 从任何输入(Terminal/飞书/GitHub)都先转化为 Issue，再拆解为 Task

2. **Stage Pipeline**: 明确的阶段流转，支持人机协作（人工确认后进入下一阶段）

3. **Local First**: 本地 Markdown 为数据源头，飞书为增强视图

4. **多 Agent 并行**: 支持不同厂商 AI Agent 在多机器上并行执行任务

5. **统一配置**: 通过 agents.yaml 管理所有 Agent 配置，支持机器粒度和 Agent 粒度设置

### 8.2 预期实现场景

**场景1: 早晨任务确认**
```
1. 用户在 TUI 看板上查看今日任务
2. 确认后，任务自动分派到各机器的 Agent
3. 用户关闭电脑出门，Agent 在后台执行
4. 飞书推送任务完成通知
5. 晚上回家查看执行结果和生成的代码
```

**场景2: 飞书远程管理**
```
1. 用户在手机上通过飞书 Bot 创建任务
2. 任务同步到本地系统
3. 家中电脑的 Agent 自动领取并执行
4. 执行结果推送到飞书
5. 用户在手机上查看摘要，回家查看详情
```

**场景3: 多 Agent 协作**
```
1. 复杂 Issue 被拆分为多个 Sub-Task
2. Sub-Task 1 分配给 MacBook 的 Claude
3. Sub-Task 2 分配给 Windows 的 Codex
4. Sub-Task 3 分配给 Linux Server 的 OpenCode
5. 所有结果汇总后统一提交
```

### 8.3 下一步行动

1. **立即开始**: 重构 Issue 模型，参考 Linear 和 Multica 的设计
2. **本周完成**: Stage Pipeline 的状态机实现
3. **下周开始**: Feishu Bitable 同步模块
4. **并行进行**: TUI Kanban 增强

---

*报告生成时间: 2026-04-06*
*基于: Linear SDK + Multica 代码分析*
