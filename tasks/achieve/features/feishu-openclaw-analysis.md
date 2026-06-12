# Feishu OpenClaw 与本地 Agent Runtime 混合架构分析

> **分析日期**: 2026-04-07  
> **核心洞察**: Feishu OpenClaw 可作为远程 Agent Runtime，本地专注于任务编排

---

## 1. 现状分析

### 1.1 Feishu OpenClaw 能力

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Feishu OpenClaw Platform                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      AI Capabilities                             │    │
│  │                                                                  │    │
│  │  • 代码生成 (Code Generation)                                    │    │
│  │  • 代码审查 (Code Review)                                        │    │
│  │  • 文档理解 (Document Understanding)                             │    │
│  │  • 知识问答 (Knowledge Q&A)                                      │    │
│  │  • 多轮对话 (Multi-turn Conversation)                            │    │
│  │  • 工具调用 (Tool Calling)                                       │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Data Access                                 │    │
│  │                                                                  │    │
│  │  • 飞书文档 (Native Access)     ✅ 直接读取                      │    │
│  │  • 飞书表格 (Bitable)           ✅ 直接读写                      │    │
│  │  • 飞书群聊 (Group Chat)        ✅ 上下文感知                    │    │
│  │  • 外部 API (External APIs)     ✅ 通过插件                      │    │
│  │  • GitHub/GitLab                ✅ 通过集成                      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Execution Environment                       │    │
│  │                                                                  │    │
│  │  • 云端沙箱 (Sandbox)           ✅ 隔离执行                      │    │
│  │  • 预装环境 (Pre-installed)     ✅ Python, Node, Go              │    │
│  │  • 文件存储 (File Storage)      ✅ 临时/持久                     │    │
│  │  • 网络访问 (Network)           ✅ 受限但可用                    │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 架构对比

| 维度 | 纯本地方案 (Original) | 纯云端方案 (OpenClaw) | 混合方案 (Hybrid) |
|------|---------------------|---------------------|------------------|
| **Agent 位置** | 本地机器 | 飞书云端 | 本地 + 云端 |
| **代码位置** | 本地文件系统 | 云端沙箱/飞书文档 | 按需选择 |
| **执行环境** | 用户机器 | 飞书服务器 | 灵活切换 |
| **数据隐私** | 高（不上云） | 中（飞书云端） | 可控 |
| **网络依赖** | 低 | 高 | 中 |
| **成本** | API Token 成本 | 飞书订阅成本 | 优化选择 |
| **可访问性** | 单机器 | 多设备 | 多设备 |
| **复杂度** | 高（Daemon管理） | 低（托管） | 中（抽象层） |

---

## 2. 新的架构设计

### 2.1 混合 Runtime 架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Hybrid Agent Runtime Architecture                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   Local Layer (本地层)                                                   │
│   ─────────────────                                                     │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                    Task Orchestrator                             │   │
│   │   - Issue 管理                                                    │   │
│   │   - 任务路由 (Router)                                            │   │
│   │   - 状态同步                                                      │   │
│   │   - Wiki 维护                                                     │   │
│   └───────────────────────────┬─────────────────────────────────────┘   │
│                               │                                          │
│           ┌───────────────────┼───────────────────┐                      │
│           │                   │                   │                      │
│           ▼                   ▼                   ▼                      │
│   ┌───────────────┐   ┌───────────────┐   ┌───────────────┐            │
│   │ Local Runtime │   │ Local Runtime │   │ Feishu        │            │
│   │ (Claude)      │   │ (Codex)       │   │ OpenClaw      │            │
│   │               │   │               │   │ (Remote)      │            │
│   │ • 敏感代码     │   │ • 快速原型     │   │               │            │
│   │ • 本地测试     │   │ • 日常开发     │   │ • 文档生成     │            │
│   │ • 私有项目     │   │ • 标准任务     │   │ • 数据分析     │            │
│   │               │   │               │   │ • 协作场景     │            │
│   └───────┬───────┘   └───────┬───────┘   └───────┬───────┘            │
│           │                   │                   │                      │
│           └───────────────────┼───────────────────┘                      │
│                               │                                          │
│                               ▼                                          │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                    Unified Output                                │   │
│   │   - 结果归档到 Wiki                                              │   │
│   │   - 同步到 Feishu Bitable                                        │   │
│   │   - 本地 Markdown 存储                                           │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Runtime 选择策略

```go
type RuntimeType string

const (
    RuntimeLocalClaude   RuntimeType = "local_claude"
    RuntimeLocalCodex    RuntimeType = "local_codex"
    RuntimeLocalOpenCode RuntimeType = "local_opencode"
    RuntimeRemoteFeishu  RuntimeType = "remote_feishu"  // OpenClaw
)

type RuntimeSelector interface {
    // 根据任务特征选择最佳 Runtime
    Select(issue *Issue) (RuntimeType, error)
}

// 默认选择策略
type DefaultRuntimeSelector struct {
    rules []SelectionRule
}

type SelectionRule struct {
    Name        string
    Condition   func(issue *Issue) bool
    Runtime     RuntimeType
    Priority    int
    Description string
}

// 预定义规则
var DefaultRules = []SelectionRule{
    {
        Name: "敏感代码本地执行",
        Condition: func(i *Issue) bool {
            return i.Tags.Contains("sensitive") || 
                   i.Tags.Contains("private-repo") ||
                   i.Context.Location.Contains("internal")
        },
        Runtime:     RuntimeLocalClaude,
        Priority:    100,
        Description: "包含敏感信息的任务必须在本地执行",
    },
    {
        Name: "飞书文档操作使用云端",
        Condition: func(i *Issue) bool {
            return i.Title.Contains("文档") ||
                   i.Title.Contains("表格") ||
                   i.Context.Trigger == "feishu_bot"
        },
        Runtime:     RuntimeRemoteFeishu,
        Priority:    90,
        Description: "飞书相关操作使用 OpenClaw 更高效",
    },
    {
        Name: "代码生成任务偏好本地",
        Condition: func(i *Issue) bool {
            return i.Tags.Contains("coding") ||
                   i.Tags.Contains("implementation")
        },
        Runtime:     RuntimeLocalClaude,
        Priority:    50,
        Description: "代码生成默认使用本地 Claude",
    },
    {
        Name: "分析任务可使用云端",
        Condition: func(i *Issue) bool {
            return i.Stage == StageAnalysis ||
                   i.Tags.Contains("research") ||
                   i.Tags.Contains("analysis")
        },
        Runtime:     RuntimeRemoteFeishu,
        Priority:    40,
        Description: "分析类任务可使用云端资源",
    },
}
```

---

## 3. Feishu OpenClaw 适配器设计

### 3.1 OpenClaw Runtime 实现

```go
package runtime

// OpenClawRuntime implements AgentRuntime for Feishu OpenClaw
type OpenClawRuntime struct {
    client      *feishu.Client
    appID       string
    appSecret   string
    
    // OpenClaw specific
    webhookURL  string
    callbackURL string
}

func NewOpenClawRuntime(config OpenClawConfig) (*OpenClawRuntime, error) {
    return &OpenClawRuntime{
        client:     feishu.NewClient(config.AppID, config.AppSecret),
        appID:      config.AppID,
        appSecret:  config.AppSecret,
    }, nil
}

// Execute implements AgentRuntime
func (r *OpenClawRuntime) Execute(ctx context.Context, task *Task) (*ExecutionResult, error) {
    // 1. 准备上下文
    context := r.prepareContext(task)
    
    // 2. 构建 OpenClaw 请求
    req := &OpenClawRequest{
        Prompt:      task.Prompt,
        Context:     context,
        Tools:       r.getAvailableTools(),
        CallbackURL: r.callbackURL,
    }
    
    // 3. 发送到 OpenClaw
    resp, err := r.client.OpenClaw().Execute(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("openclaw execute failed: %w", err)
    }
    
    // 4. 等待结果（同步或异步）
    result, err := r.waitForResult(ctx, resp.ExecutionID)
    if err != nil {
        return nil, err
    }
    
    // 5. 转换结果格式
    return r.convertResult(result), nil
}

// prepareContext 准备 OpenClaw 上下文
func (r *OpenClawRuntime) prepareContext(task *Task) OpenClawContext {
    context := OpenClawContext{
        // 关联的飞书文档
        Documents: r.getRelatedDocuments(task),
        
        // 历史对话
        History: r.getConversationHistory(task.IssueID),
        
        // 可用工具
        Tools: []string{
            "read_document",
            "write_document", 
            "search_knowledge",
            "code_interpreter",
            "web_search",
        },
    }
    
    return context
}

// getRelatedDocuments 获取关联的飞书文档
func (r *OpenClawRuntime) getRelatedDocuments(task *Task) []DocumentRef {
    var docs []DocumentRef
    
    // 从 task 上下文中提取文档链接
    for _, ref := range task.Context.WikiRefs {
        if strings.HasPrefix(ref, "https://open.feishu.cn/docx/") {
            docs = append(docs, DocumentRef{
                URL:  ref,
                Type: "docx",
            })
        }
    }
    
    return docs
}
```

### 3.2 OpenClaw 任务格式

```go
// OpenClawRequest represents a request to Feishu OpenClaw
type OpenClawRequest struct {
    // 执行 ID（用于追踪）
    ExecutionID string `json:"execution_id"`
    
    // 主要提示词
    Prompt string `json:"prompt"`
    
    // 系统提示词（可选）
    SystemPrompt string `json:"system_prompt,omitempty"`
    
    // 上下文信息
    Context OpenClawContext `json:"context"`
    
    // 工具配置
    Tools []string `json:"tools"`
    
    // 回调 URL（异步模式）
    CallbackURL string `json:"callback_url,omitempty"`
    
    // 超时设置
    TimeoutSeconds int `json:"timeout_seconds,omitempty"`
}

type OpenClawContext struct {
    // 关联文档
    Documents []DocumentRef `json:"documents,omitempty"`
    
    // 历史对话
    History []Message `json:"history,omitempty"`
    
    // 知识库
    KnowledgeBase string `json:"knowledge_base,omitempty"`
    
    // 用户信息
    User UserInfo `json:"user"`
}

type DocumentRef struct {
    URL        string `json:"url"`
    Type       string `json:"type"` // docx, sheet, bitable
    AccessToken string `json:"access_token,omitempty"`
}

type OpenClawResponse struct {
    ExecutionID string `json:"execution_id"`
    Status      string `json:"status"` // running, completed, failed
    
    // 执行结果
    Result *OpenClawResult `json:"result,omitempty"`
    
    // 错误信息
    Error *OpenClawError `json:"error,omitempty"`
}

type OpenClawResult struct {
    // 文本输出
    TextOutput string `json:"text_output"`
    
    // 生成的文档
    GeneratedDocs []DocumentRef `json:"generated_docs,omitempty"`
    
    // 代码输出
    CodeOutputs []CodeOutput `json:"code_outputs,omitempty"`
    
    // 执行日志
    ExecutionLog []LogEntry `json:"execution_log,omitempty"`
    
    // 使用的工具
    ToolsUsed []ToolUsage `json:"tools_used,omitempty"`
}
```

---

## 4. 任务路由与编排

### 4.1 智能任务路由

```go
// TaskRouter routes tasks to appropriate runtimes
type TaskRouter struct {
    runtimes map[RuntimeType]AgentRuntime
    selector RuntimeSelector
}

func (r *TaskRouter) Route(ctx context.Context, task *Task) (*ExecutionResult, error) {
    // 1. 选择 Runtime
    runtimeType, err := r.selector.Select(task.Issue)
    if err != nil {
        return nil, err
    }
    
    log.Printf("Routing task %s to runtime: %s", task.ID, runtimeType)
    
    // 2. 获取 Runtime 实例
    runtime, ok := r.runtimes[runtimeType]
    if !ok {
        return nil, fmt.Errorf("runtime %s not available", runtimeType)
    }
    
    // 3. 准备任务（根据 Runtime 类型调整）
    preparedTask := r.prepareTaskForRuntime(task, runtimeType)
    
    // 4. 执行
    result, err := runtime.Execute(ctx, preparedTask)
    if err != nil {
        // 失败时尝试 fallback
        return r.fallbackExecute(ctx, task, runtimeType, err)
    }
    
    return result, nil
}

// prepareTaskForRuntime 根据 Runtime 类型准备任务
func (r *TaskRouter) prepareTaskForRuntime(task *Task, rt RuntimeType) *Task {
    switch rt {
    case RuntimeRemoteFeishu:
        // 为 OpenClaw 优化提示词
        return r.optimizeForOpenClaw(task)
    default:
        return task
    }
}

// optimizeForOpenClaw 优化任务以适应 OpenClaw
func (r *TaskRouter) optimizeForOpenClaw(task *Task) *Task {
    optimized := *task
    
    // 添加飞书文档引用格式
    if len(task.Context.WikiRefs) > 0 {
        optimized.Prompt = fmt.Sprintf(
            "参考文档：\n%s\n\n任务：\n%s",
            formatDocRefs(task.Context.WikiRefs),
            task.Prompt,
        )
    }
    
    // OpenClaw 更适合结构化输出
    optimized.Prompt += "\n\n请以结构化格式输出结果。"
    
    return &optimized
}

// fallbackExecute 失败时的备用执行
func (r *TaskRouter) fallbackExecute(
    ctx context.Context, 
    task *Task, 
    failedRuntime RuntimeType,
    originalErr error,
) (*ExecutionResult, error) {
    log.Printf("Runtime %s failed: %v, trying fallback", failedRuntime, originalErr)
    
    // 云端失败，尝试本地
    if failedRuntime == RuntimeRemoteFeishu {
        if localRuntime, ok := r.runtimes[RuntimeLocalClaude]; ok {
            return localRuntime.Execute(ctx, task)
        }
    }
    
    // 本地失败，尝试云端
    if strings.HasPrefix(string(failedRuntime), "local_") {
        if cloudRuntime, ok := r.runtimes[RuntimeRemoteFeishu]; ok {
            return cloudRuntime.Execute(ctx, task)
        }
    }
    
    return nil, originalErr
}
```

### 4.2 统一输出处理

无论使用哪个 Runtime，输出格式统一：

```go
// ExecutionResult 统一的执行结果
type ExecutionResult struct {
    // 基本信息
    TaskID      string        `json:"task_id"`
    RuntimeType RuntimeType   `json:"runtime_type"`
    Status      string        `json:"status"` // success, failed, partial
    
    // 输出内容
    Summary     string        `json:"summary"`
    Artifacts   []Artifact    `json:"artifacts"`
    
    // 代码相关
    CodeChanges []CodeChange  `json:"code_changes,omitempty"`
    GitCommit   string        `json:"git_commit,omitempty"`
    
    // 文档相关
    GeneratedDocs []DocumentRef `json:"generated_docs,omitempty"`
    
    // 元数据
    Duration    int           `json:"duration_seconds"`
    TokensUsed  int           `json:"tokens_used,omitempty"`
    Cost        float64       `json:"cost,omitempty"`
}

// Artifact 统一的产物格式
type Artifact struct {
    Type        string `json:"type"` // code, doc, image, data
    Name        string `json:"name"`
    Content     string `json:"content,omitempty"`
    URL         string `json:"url,omitempty"`
    Size        int    `json:"size,omitempty"`
    MIMEType    string `json:"mime_type,omitempty"`
}
```

---

## 5. 配置更新

### 5.1 新的 agents.yaml 配置

```yaml
version: "2.0"

# Runtime 配置
runtimes:
  # 本地 Runtime
  local:
    enabled: true
    agents:
      claude:
        type: "claude"
        executable: "/usr/local/bin/claude"
        model: "claude-sonnet-4"
        max_concurrent: 3
        priority: 10  # 优先级，数字越小越优先
        
      codex:
        type: "codex"
        executable: "codex"
        model: "gpt-4o"
        max_concurrent: 2
        priority: 20
  
  # 远程 Runtime (Feishu OpenClaw)
  remote:
    feishu_openclaw:
      enabled: true
      app_id: "${FEISHU_APP_ID}"
      app_secret: "${FEISHU_APP_SECRET}"
      webhook_url: "https://open.feishu.cn/open-apis/ai/v1/execute"
      callback_url: "https://your-server.com/callbacks/openclaw"
      max_concurrent: 5
      priority: 30
      
      # OpenClaw 特定配置
      features:
        - "document_read"
        - "document_write"
        - "bitable_access"
        - "code_interpreter"
        - "web_search"
      
      # 自动路由规则
      auto_route:
        prefer_for:
          - "document_analysis"
          - "data_processing"
          - "collaborative_tasks"
        avoid_for:
          - "sensitive_code"
          - "private_repositories"

# 路由策略
routing:
  default_runtime: "local_claude"
  
  # 故障转移
  failover:
    enabled: true
    attempts: 2
    fallback_order:
      - "local_claude"
      - "remote_feishu"
      - "local_codex"
  
  # 用户覆盖
  allow_user_override: true
```

---

## 6. 使用场景示例

### 6.1 场景 1: 文档分析任务 → OpenClaw

```
用户: capture add "分析产品需求文档并提取技术要点"

系统自动路由:
- 检测到关键词 "文档"
- 选择 Runtime: remote_feishu (OpenClaw)
- OpenClaw 直接读取飞书文档
- 分析结果归档到本地 Wiki
```

### 6.2 场景 2: 敏感代码任务 → Local

```
用户: capture add "实现内部加密算法" --tag=sensitive

系统自动路由:
- 检测到标签 "sensitive"
- 选择 Runtime: local_claude
- 本地执行，代码不上云
- 结果本地归档
```

### 6.3 场景 3: 故障转移

```
任务: 生成 API 文档

尝试 1: local_claude → 失败（网络问题）
自动故障转移 ↓
尝试 2: remote_feishu → 成功
结果: 使用 OpenClaw 生成文档
```

---

## 7. 优势与挑战

### 7.1 优势

| 优势 | 说明 |
|------|------|
| **简化部署** | 不需要本地安装 Claude/Codex |
| **跨设备** | 飞书账号即服务，任意设备访问 |
| **原生集成** | 与飞书文档/表格无缝协作 |
| **成本优化** | 按需选择 Runtime，优化成本 |
| **弹性伸缩** | 云端自动扩缩容 |

### 7.2 挑战与解决方案

| 挑战 | 解决方案 |
|------|----------|
| **数据隐私** | 敏感任务强制本地执行，路由规则可配置 |
| **网络依赖** | 本地 Runtime 作为离线备份 |
| **Vendor Lock-in** | 抽象 Runtime 接口，易于切换 |
| **调试困难** | 云端执行日志同步到本地 |

---

## 8. 实施建议

### 8.1 实施顺序

1. **Phase 1**: 抽象 Runtime 接口（1 周）
2. **Phase 2**: 实现 OpenClaw Runtime（1 周）
3. **Phase 3**: 实现智能路由（1 周）
4. **Phase 4**: 故障转移与监控（1 周）

### 8.2 推荐配置

```yaml
# 起步配置（只有 OpenClaw）
runtimes:
  remote:
    feishu_openclaw:
      enabled: true
  local:
    enabled: false

# 生产配置（混合）
runtimes:
  remote:
    feishu_openclaw:
      enabled: true
  local:
    enabled: true
    agents:
      claude:
        enabled: true
      codex:
        enabled: true

routing:
  default_runtime: "local_claude"
  failover:
    enabled: true
```

---

## 9. 结论

### 9.1 核心观点

1. **Feishu OpenClaw 是可行的 Remote Runtime**：功能完善，与飞书生态深度集成
2. **混合架构是最佳选择**：兼顾灵活性、成本和隐私
3. **任务编排是核心价值**：本地专注于路由、编排、知识管理
4. **抽象层是关键**：Runtime 接口屏蔽底层差异

### 9.2 下一步行动

1. **立即**：定义 `AgentRuntime` 接口
2. **本周**：实现 `OpenClawRuntime` 适配器
3. **下周**：实现智能路由和故障转移
4. **验证**：选择 3-5 个任务测试混合执行

---

*分析完成: 2026-04-07*
