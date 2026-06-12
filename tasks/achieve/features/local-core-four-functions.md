# 本地核心四大功能

> **加工、排序、委派、查看** - 本地系统的核心价值

---

## 1. 一句话总结

```
本地 = 加工(Processing) + 排序(Prioritization) + 委派(Delegation) + 查看(Tracking)
IM = 入口(Input) + 出口(Notification)
```

---

## 2. 四大功能详解

### 2.1 加工 (Processing)

**定义**: 将原始输入转化为结构化、可执行的知识

```
原始输入                    加工过程                    输出
────────                    ────────                    ────
"研究 Go vs Rust"    →    信息收集              →    研究报告
                           知识提取                    决策矩阵
                           对比分析                    执行计划
                           综合总结                    知识卡片
```

**为什么必须在本地**:
| 特性 | 说明 |
|------|------|
| **时间长** | 可能需要分钟级甚至小时级 |
| **资源多** | 需要访问大量参考资料 |
| **状态复杂** | 多阶段、可暂停恢复 |
| **输出结构化** | 不仅是文本，还有表格、代码、图表 |

**IM 不适合**: 不能让用户在聊天界面等 10 分钟看进度刷屏

---

### 2.2 排序 (Prioritization)

**定义**: 根据多维度因素决定任务执行顺序

```
排序维度:
├── 紧急度 (Deadline)
├── 重要性 (Impact)
├── 依赖关系 (Dependencies)
├── 资源可用性 (Resources)
├── 上下文连贯性 (Context)
└── 个人习惯 (Preference)

示例决策:
TASK-001: 紧急但不重要 → 插队但限时
TASK-002: 重要有依赖   → 等待依赖完成
TASK-003: 低优长耗时   → 后台批量处理
```

**为什么必须在本地**:
| 特性 | 说明 |
|------|------|
| **需要全局视图** | 要知道所有待办任务才能排序 |
| **复杂规则** | 不能简单 FIFO，要加权计算 |
| **动态调整** | 新任务加入可能需要重排序 |
| **个人化** | 每个人的优先级算法不同 |

**IM 不适合**: 聊天界面无法展示全局任务视图，做不了复杂排序决策

---

### 2.3 委派 (Delegation)

**定义**: 将任务分配到合适的执行环境

```
委派决策树:
                    任务到达
                       │
                       ▼
              ┌─────────────────┐
              │   分析任务特征   │
              └────────┬────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
    ┌─────────┐  ┌─────────┐  ┌─────────┐
    │ 敏感数据 │  │ 协作需求 │  │ 高算力  │
    │   Yes   │  │   Yes   │  │   Yes   │
    └────┬────┘  └────┬────┘  └────┬────┘
         │             │             │
         ▼             ▼             ▼
   ┌──────────┐  ┌──────────┐  ┌──────────┐
   │ 本地执行  │  │ 云端执行  │  │ 云端执行  │
   │ (Claude) │  │ (OpenClaw)│  │ (GPU实例) │
   └──────────┘  └──────────┘  └──────────┘
```

**为什么必须在本地**:
| 特性 | 说明 |
|------|------|
| **需要全局状态** | 要知道所有 runtime 的负载状况 |
| **数据位置感知** | 数据在哪，计算就应该在哪 |
| **成本优化** | 需要比较各选项的成本 |
| **故障转移** | 一个失败要能自动切换 |

**IM 不适合**: 聊天界面无法实时监控多个 runtime 状态

---

### 2.4 查看 (Tracking)

**定义**: 追踪任务执行状态和进度

```
查看维度:
├── 执行进度 (Progress %)
├── 当前阶段 (Current Stage)
├── 已用时间 (Elapsed)
├── 预估剩余 (ETA)
├── 中间产物 (Artifacts)
├── 日志输出 (Logs)
└── 异常信息 (Errors)

展示方式:
├── TUI Dashboard (实时)
├── Web Dashboard (详细)
├── CLI Status (快速)
└── IM Notification (轻量)
```

**为什么必须在本地**:
| 特性 | 说明 |
|------|------|
| **实时更新** | 秒级状态变化 |
| **详细信息** | 需要看日志、中间结果 |
| **历史追溯** | 需要查看过去执行记录 |
| **多任务视图** | 同时看多个任务状态 |

**IM 不适合**: 不能每分钟推送一条消息更新进度，太吵

---

## 3. 四大功能的关系

```
                    输入 (IM/Terminal)
                           │
                           ▼
                    ┌──────────────┐
                    │     排序     │ ◄── 全局任务队列
                    │ (Prioritize) │     决定执行顺序
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     委派     │ ◄── 选择执行环境
                    │   (Delegate) │     本地/云端/人工
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     加工     │ ◄── 实际执行
                    │  (Process)   │     深度处理
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     查看     │ ◄── 监控反馈
                    │   (Track)    │     进度状态
                    └──────┬───────┘
                           │
                           ▼
                    输出 (Notification)
```

**循环**: 查看过程中可能需要重新排序、重新委派、或者触发新的加工

---

## 4. 与 IM 的边界

### 4.1 IM 做这四件事的问题

| 功能 | IM 尝试做 | 问题 |
|------|----------|------|
| **加工** | 在聊天里等 AI 回复 | 时间太长，界面被占 |
| **排序** | 问用户"先做哪个" | 用户没有全局视图 |
| **委派** | 让用户选择 runtime | 用户不懂技术细节 |
| **查看** | 推送进度消息 | 消息刷屏，信息过载 |

### 4.2 正确的分工

```
IM 职责:
├── 输入: "帮我研究 xxx" (触发加工)
├── 查询: "现在在进行什么？" (快速查看摘要)
├── 决策: "A 和 B 选哪个？" (简单二选一)
└── 通知: "任务完成" (结果通知)

本地职责:
├── 加工: 实际的研究、分析、生成
├── 排序: 自动决定执行顺序
├── 委派: 自动选择最佳执行环境
└── 查看: 详细的实时仪表板
```

---

## 5. 技术实现要点

### 5.1 加工层 (Processing)

```go
type Processor interface {
    // 启动加工流程
    Start(ctx context.Context, task Task) (*Process, error)
    
    // 暂停/恢复
    Pause(processID string) error
    Resume(processID string) error
    
    // 获取状态
    Status(processID string) (*ProcessStatus, error)
    
    // 取消
    Cancel(processID string) error
}

// 加工流程内部
type Process struct {
    ID       string
    Task     Task
    Stages   []Stage      // 多阶段
    Current  int          // 当前阶段
    Progress float64      // 总体进度
    Logs     []LogEntry   // 执行日志
    Artifacts []Artifact  // 中间产物
}
```

### 5.2 排序层 (Prioritization)

```go
type Prioritizer interface {
    // 添加任务到队列
    Enqueue(task Task) error
    
    // 获取下一个要执行的任务
    Dequeue() (*Task, error)
    
    // 重新排序（当优先级变化时）
    Reorder() error
    
    // 查看队列
    Queue() ([]Task, error)
}

// 排序算法
type PriorityAlgorithm interface {
    Calculate(task Task, context Context) float64
}

// 示例：加权评分
func WeightedScore(task Task) float64 {
    urgency := 1.0 / (hoursUntilDeadline(task) + 1)
    importance := task.ImpactScore
    effort := 1.0 / (task.EstimatedHours + 1)
    
    return urgency*0.4 + importance*0.4 + effort*0.2
}
```

### 5.3 委派层 (Delegation)

```go
type Delegator interface {
    // 选择最佳执行环境
    SelectRuntime(task Task) (Runtime, error)
    
    // 执行任务
    Dispatch(task Task, runtime Runtime) (*Execution, error)
    
    // 监控运行中的任务
    Monitor(executionID string) (*ExecutionStatus, error)
}

type RuntimeSelector struct {
    runtimes []Runtime
    strategy SelectionStrategy
}

func (s *RuntimeSelector) Select(task Task) (Runtime, error) {
    // 1. 过滤不满足条件的 runtime
    candidates := s.filter(task)
    
    // 2. 按策略排序
    s.strategy.Rank(candidates, task)
    
    // 3. 返回最佳选项
    return candidates[0], nil
}
```

### 5.4 查看层 (Tracking)

```go
type Tracker interface {
    // 订阅任务更新
    Subscribe(taskID string, callback UpdateCallback) error
    
    // 获取当前状态
    Status(taskID string) (*TaskStatus, error)
    
    // 获取历史记录
    History(taskID string) ([]Event, error)
    
    // 获取实时日志
    Logs(taskID string, follow bool) (<-chan LogEntry, error)
}

// TUI Dashboard
type Dashboard struct {
    tasks []TaskStatus
    
    // 刷新显示
    Refresh()
    
    // 处理用户输入
    HandleInput(key string)
}
```

---

## 6. 用户交互示例

### 6.1 完整流程

```
[用户 Terminal]
$ capture add "研究 Go vs Rust"
✓ 创建任务 TASK-001

[本地 - 排序]
→ 评估优先级: 高 (用户直接创建)
→ 队列位置: #1

[本地 - 委派]
→ 分析: 需要大量搜索，无敏感数据
→ 选择: Cloud Runtime (OpenClaw)

[本地 - 加工]
→ 启动研究流程
→ Phase 1/4: 信息收集 (25%)
→ Phase 2/4: 知识提取 (50%)
→ Phase 3/4: 对比分析 (75%)
→ Phase 4/4: 生成报告 (100%)
→ 归档到 Wiki

[用户 TUI Dashboard]
┌─────────────────────────────────────┐
│ Task: TASK-001                      │
│ Status: ✅ Completed                │
│ Duration: 15m 32s                   │
│ Output: wiki/research/go-vs-rust    │
│                                     │
│ [View Report] [Create Issues]       │
└─────────────────────────────────────┘

[用户 IM]
💬 飞书通知: "TASK-001 完成，点击查看"

[用户 Terminal]
$ capture list
ID       TITLE               STATUS    PRIORITY
TASK-001 研究 Go vs Rust     Done      High
```

---

## 7. 总结

### 7.1 核心价值

| 功能 | 价值 | 为什么本地 |
|------|------|-----------|
| **加工** | 深度处理信息 | 时间长、资源多、状态复杂 |
| **排序** | 全局优化调度 | 需要全局视图、复杂算法 |
| **委派** | 最优资源利用 | 需要监控多 runtime |
| **查看** | 实时状态感知 | 需要详细信息、历史追溯 |

### 7.2 与 IM 的边界

```
IM 是 "Lightweight Interface"
- 输入触发
- 简单查询
- 结果通知

本地是 "Heavyweight Engine"
- 深度加工
- 智能决策
- 状态管理
```

### 7.3 一句话

> **本地做「重活」，IM 做「轻交互」。**

加工、排序、委派、查看是「重活」，必须在本地。

---

*文档完成: 2026-04-07*
