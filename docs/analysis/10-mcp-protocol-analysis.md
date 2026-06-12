# MCP 协议深度分析：Capture 作为 MCP Server 的对接方案

> 分析时间：2026-06-04
> 分析范围：MCP 传输协议、Claude/OpenAI/Kimi 的 MCP 接入方式、Capture v2 的实现路径
> 协议版本：MCP Specification 2025-03-26 / 2025-11-25 / 2025-06-18 (current stable)

---

## 一、MCP 协议核心概览

### 1.1 什么是 MCP

MCP（Model Context Protocol）是 Anthropic 于 2024-11 开源的协议，2025-12 捐赠给 Linux Foundation 下的 Agentic AI Foundation。它定义了 AI Agent 与外部工具/数据源之间的标准通信方式，被业界称为 **"AI 的 USB-C"**。

**核心设计目标**：
- 将 N×M 的集成问题（N 个 AI 平台 × M 个工具）简化为 N+M
- 一次构建 MCP Server，任意 MCP Client 都能使用

### 1.2 协议栈分层

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  Tools / Resources / Prompts / Sampling / Elicitation        │
├─────────────────────────────────────────────────────────────┤
│                    Protocol Layer                            │
│  JSON-RPC 2.0 (Request / Response / Notification)            │
├─────────────────────────────────────────────────────────────┤
│                    Transport Layer                           │
│  stdio  │  Streamable HTTP  │  HTTP+SSE (legacy)            │
├─────────────────────────────────────────────────────────────┤
│                    Connection Lifecycle                      │
│  Initialize → Active Session → Termination                   │
└─────────────────────────────────────────────────────────────┘
```

### 1.3 JSON-RPC 2.0 消息格式

MCP 的所有通信都基于 JSON-RPC 2.0，UTF-8 编码。

**Request（客户端 → 服务器）**：
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "capture_list_tasks",
    "arguments": {
      "workspace_id": "default",
      "project_id": "innate-capture",
      "status": "todo"
    }
  }
}
```

**Response（服务器 → 客户端）**：
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "[TASK-001] 设计 Workspace API\n[TASK-002] 实现 WorkspaceStore"
      }
    ],
    "isError": false
  }
}
```

**Notification（单向，无 id）**：
```json
{
  "jsonrpc": "2.0",
  "method": "notifications/tools/list_changed",
  "params": {}
}
```

**Error Response**：
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "workspace_id 'default' not found"
    }
  }
}
```

---

## 二、传输层详解

### 2.1 三种传输方式对比

| 特性 | stdio | Streamable HTTP | HTTP+SSE (Legacy) |
|------|-------|-----------------|-------------------|
| **方向性** | 双向（stdin/stdout） | 双向（GET/POST 同端点） | 双向（POST + 独立 SSE） |
| **端点** | 无（子进程） | 单端点 `POST /mcp` | 双端点（`POST /messages` + `GET /sse`） |
| **Session ID** | 不适用（单进程对） | `Mcp-Session-Id` HTTP Header | 查询参数 `?sessionId=xxx` |
| **并发客户端** | 1（单租户） | 多（多租户） | 多 |
| **防火墙友好** | N/A（本地） | 是（标准 HTTP） | 是 |
| **状态** | Active | Active（2025-03-26 起推荐） | Deprecated（保留向后兼容） |
| **适用场景** | 本地 CLI 工具 | 远程/云端服务 | 旧版客户端兼容 |

### 2.2 stdio 传输（推荐用于本地）

**工作原理**：
```
Client Process                    Server Process (Capture)
    │                                 │
    ├── spawn: capture mcp serve ──>│
    │                                 │
    ├── stdin: {"jsonrpc":"2.0",...}→│
    │                                 │
    │←─ stdout: {"jsonrpc":"2.0",...}┤
    │                                 │
    ├── stdin: {"jsonrpc":"2.0",...}→│
    │←─ stdout: {"jsonrpc":"2.0",...}┤
    │                                 │
    ├── close stdin ───────────────>│ (server exits)
```

**关键规则**：
- 消息以 newline（`\n`）分隔
- 消息内不能包含 embedded newlines
- Server 只能写 JSON-RPC 到 stdout，`stderr` 用于日志
- Client 只能写 JSON-RPC 到 stdin

**Capture v2 的 stdio 实现（Go 伪代码）**：
```go
package mcp

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

type StdioTransport struct {
    scanner *bufio.Scanner
    writer  *bufio.Writer
}

func NewStdioTransport() *StdioTransport {
    return &StdioTransport{
        scanner: bufio.NewScanner(os.Stdin),
        writer:  bufio.NewWriter(os.Stdout),
    }
}

func (t *StdioTransport) ReadMessage() (*JSONRPCMessage, error) {
    if !t.scanner.Scan() {
        return nil, t.scanner.Err()
    }
    var msg JSONRPCMessage
    if err := json.Unmarshal(t.scanner.Bytes(), &msg); err != nil {
        return nil, err
    }
    return &msg, nil
}

func (t *StdioTransport) WriteMessage(msg *JSONRPCMessage) error {
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    // 关键：一条消息一行，不能包含换行
    if _, err := t.writer.Write(data); err != nil {
        return err
    }
    if _, err := t.writer.WriteString("\n"); err != nil {
        return err
    }
    return t.writer.Flush()
}
```

### 2.3 Streamable HTTP 传输（推荐用于远程）

**工作原理**（2025-03-26 引入，替代 SSE）：
```
Client                                    Capture MCP Server
  │                                              │
  ├── POST /mcp ────────────────────────────────>│
  │   {"jsonrpc":"2.0","id":1,"method":"initialize"}│
  │                                              │
  │<─ 200 OK ────────────────────────────────────┤
  │   Mcp-Session-Id: abc-123                    │
  │   {"jsonrpc":"2.0","id":1,"result":{...}}    │
  │                                              │
  ├── POST /mcp ────────────────────────────────>│
  │   Mcp-Session-Id: abc-123                    │
  │   {"jsonrpc":"2.0","id":2,"method":"tools/list"}│
  │                                              │
  │<─ 200 OK ────────────────────────────────────┤
  │   Content-Type: application/json             │
  │   {"jsonrpc":"2.0","id":2,"result":{...}}    │
  │                                              │
  │   [或流式响应]                                │
  │<─ 200 OK ────────────────────────────────────┤
  │   Content-Type: text/event-stream            │
  │   data: {"jsonrpc":"2.0",...}                │
  │   data: {"jsonrpc":"2.0",...}                │
```

**关键特性**：
- **单端点**：`POST /mcp` 处理所有请求
- **Session 管理**：Server 通过 `Mcp-Session-Id` Header 设置 session，Client 在后续请求中回传
- **响应模式可选**：
  - 非流式：`Content-Type: application/json`
  - 流式：`Content-Type: text/event-stream`（SSE chunk 格式）
- **安全要求**：必须验证 `Origin` Header，默认绑定 `127.0.0.1`

**SSE 事件格式**（流式响应时）：
```
HTTP/1.1 200 OK
Content-Type: text/event-stream
Mcp-Session-Id: abc-123

data: {"jsonrpc":"2.0","id":2,"result":{...}}

data: {"jsonrpc":"2.0","id":3,"result":{...}}

```

### 2.4 Legacy HTTP+SSE 传输（已废弃）

**工作原理**：
```
Client                                    Server
  │                                              │
  ├── GET /sse ─────────────────────────────────>│
  │<─ event: endpoint ───────────────────────────┤
  │   data: /messages?sessionId=abc              │
  │                                              │
  ├── POST /messages?sessionId=abc ─────────────>│
  │   {"jsonrpc":"2.0","id":1,"method":"tools/list"}│
  │                                              │
  │<─ SSE event: {"jsonrpc":"2.0","id":1,...} ───┤
```

**状态**：2025-03-26 起正式废弃，但 SDK 保留向后兼容。

---

## 三、MCP 连接生命周期

### 3.1 初始化握手

```
Client                                              Server
  │                                                   │
  ├── Request: initialize ──────────────────────────> │
  │   {                                               │
  │     "jsonrpc": "2.0",                             │
  │     "id": 1,                                      │
  │     "method": "initialize",                       │
  │     "params": {                                   │
  │       "protocolVersion": "2025-03-26",            │
  │       "capabilities": {                           │
  │         "tools": { "listChanged": true },         │
  │         "resources": { "subscribe": true }        │
  │       },                                          │
  │       "clientInfo": {                             │
  │         "name": "claude-code",                    │
  │         "version": "1.0.0"                        │
  │       }                                           │
  │     }                                             │
  │   }                                               │
  │                                                   │
  │<─ Response: initialize ────────────────────────── │
  │   {                                               │
  │     "jsonrpc": "2.0",                             │
  │     "id": 1,                                      │
  │     "result": {                                   │
  │       "protocolVersion": "2025-03-26",            │
  │       "capabilities": {                           │
  │         "tools": { "listChanged": true },         │
  │         "resources": {}                           │
  │       },                                          │
  │       "serverInfo": {                             │
  │         "name": "capture-mcp",                    │
  │         "version": "2.0.0"                        │
  │       }                                           │
  │     }                                             │
  │   }                                               │
  │                                                   │
  ├── Notification: notifications/initialized ──────> │
  │   {                                               │
  │     "jsonrpc": "2.0",                             │
  │     "method": "notifications/initialized"         │
  │   }                                               │
```

### 3.2 活跃会话阶段

双方可以：
- 发送 Request（期望 Response）
- 发送 Notification（不期望响应）
- 取消在途请求（`notifications/cancelled`）
- 发送进度更新（`notifications/progress`）

### 3.3 终止连接

- 关闭传输（stdio: 关闭 stdin；HTTP: 关闭连接）
- 发送协议级终止消息
- 或静默断开（不推荐）

---

## 四、MCP 核心能力（Primitives）

### 4.1 Tools（工具调用）

Server 暴露可调用的工具，Client（AI）根据用户请求决定调用哪个。

**Tools/List**：
```json
// Request
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}

// Response
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "capture_list_tasks",
        "description": "列出指定项目和状态的任务",
        "inputSchema": {
          "type": "object",
          "properties": {
            "workspace_id": { "type": "string", "description": "Workspace ID" },
            "project_id": { "type": "string", "description": "Project ID" },
            "status": { "type": "string", "enum": ["todo", "in_progress", "done"] }
          },
          "required": ["workspace_id"]
        }
      },
      {
        "name": "capture_create_task",
        "description": "在指定项目下创建新任务",
        "inputSchema": {
          "type": "object",
          "properties": {
            "project_id": { "type": "string" },
            "title": { "type": "string" },
            "description": { "type": "string" },
            "priority": { "type": "string", "enum": ["high", "medium", "low"] }
          },
          "required": ["project_id", "title"]
        }
      }
    ]
  }
}
```

**Tools/Call**：
```json
// Request
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "capture_list_tasks",
    "arguments": {
      "workspace_id": "default",
      "project_id": "innate-capture",
      "status": "todo"
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "找到 3 个任务:\n1. [TASK-001] 设计 Workspace API (high)\n2. [TASK-002] 实现 WorkspaceStore (medium)\n3. [TASK-003] 编写 CLI 命令 (medium)"
      }
    ],
    "isError": false
  }
}
```

### 4.2 Resources（资源读取）

Server 暴露可读的资源，通过 URI 标识。

**Resources/List**：
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list",
  "params": {}
}
```

**Resources/Read**：
```json
// Request
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/read",
  "params": {
    "uri": "capture://task/TASK-001"
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "contents": [
      {
        "uri": "capture://task/TASK-001",
        "mimeType": "text/markdown",
        "text": "---\nid: TASK-001\ntitle: ...\n---"
      }
    ]
  }
}
```

### 4.3 Prompts（提示模板）

Server 暴露可复用的 Prompt 模板。

**Prompts/Get**：
```json
// Request
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "prompts/get",
  "params": {
    "name": "brainstorm-scamper",
    "arguments": {
      "task_id": "TASK-001"
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": 6,
  "result": {
    "description": "SCAMPER 分析",
    "messages": [
      {
        "role": "user",
        "content": {
          "type": "text",
          "text": "请对功能 '设计 Workspace API' 进行 SCAMPER 分析..."
        }
      }
    ]
  }
}
```

---

## 五、对接 Claude MCP

### 5.1 Claude Desktop

**配置方式**：`~/Library/Application Support/Claude/claude_desktop_config.json`（macOS）

```json
{
  "mcpServers": {
    "capture": {
      "command": "capture",
      "args": ["mcp", "serve", "--transport", "stdio"]
    }
  }
}
```

**Claude Code 终端**：
```bash
# 添加 stdio 模式
claude mcp add --transport stdio capture -- capture mcp serve

# 添加 HTTP 模式
claude mcp add --transport http capture http://localhost:8080/mcp

# 列出已配置
claude mcp list

# 交互式管理
/mcp
```

### 5.2 Claude 的 MCP 消息流程

```
Claude Desktop/Claude Code
        │
        ├── 1. spawn: capture mcp serve
        │
        ├── 2. stdin → {"method":"initialize",...}
        │
        ├── 3. stdout ← {"result":{"protocolVersion":"2025-03-26",...}}
        │
        ├── 4. stdin → {"method":"notifications/initialized"}
        │
        ├── 5. stdin → {"method":"tools/list"}
        │
        ├── 6. stdout ← {"result":{"tools":[...]}}
        │
        ├── [用户说："列出 innate-capture 项目的待办任务"]
        │
        ├── 7. Claude 内部决策：调用 capture_list_tasks
        │
        ├── 8. stdin → {"method":"tools/call","params":{"name":"capture_list_tasks",...}}
        │
        └── 9. stdout ← {"result":{"content":[{"text":"..."}]}}
```

### 5.3 Claude 的 Messages API 与 MCP 的关系

Claude 有两个不同的 API/协议：

| 协议 | 用途 | 通信方向 |
|------|------|----------|
| **Messages API** | 直接与 Claude LLM 对话（HTTP REST API） | 用户/应用 → Anthropic API |
| **MCP Protocol** | Claude 作为 Client，连接外部 Tools/Resources | Claude ↔ 本地/远程 MCP Server |

**关键理解**：
- Messages API 是你"调用 Claude"的方式
- MCP 是 Claude"调用你的工具"的方式
- 两者互补，不是替代关系

**实际交互流程**：
```
User: "列出我的待办任务"
  │
  ▼
Claude Desktop (Host)
  │── 通过 Messages API 发送用户消息到 Claude LLM
  │
Claude LLM 分析：用户要查看任务 → 需要调用 capture_list_tasks
  │
  ▼
Claude Desktop (MCP Client)
  │── 通过 MCP Protocol 发送 tools/call 到 Capture MCP Server
  │
Capture MCP Server
  │── 查询 SQLite / Markdown
  │── 返回任务列表
  │
  ▼
Claude Desktop
  │── 将工具结果通过 Messages API 送回 Claude LLM
  │
Claude LLM 生成自然语言回复
  │
  ▼
User: " innate-capture 项目有 3 个待办：..."
```

---

## 六、对接 OpenAI MCP

### 6.1 OpenAI 的 MCP 支持现状

- **2025-03**：OpenAI 宣布全面支持 MCP
- **Responses API**：原生连接远程 MCP 服务器
- **Agents SDK**：支持 MCP 工具调用
- **ChatGPT Desktop**：2025-09 添加自定义 MCP 支持
- **严格认证要求**：OAuth 2.1 + Dynamic Client Registration 强制

### 6.2 OpenAI Responses API 调用 MCP

```python
from openai import OpenAI

client = OpenAI()

resp = client.responses.create(
    model="gpt-5",
    tools=[
        {
            "type": "mcp",
            "server_label": "capture",
            "server_url": "https://capture.example.com/mcp",
            # OAuth 2.1 认证（OpenAI 强制要求）
            "require_approval": "never"  # 或指定白名单工具
        }
    ],
    input="列出 innate-capture 项目中所有高优先级的待办任务"
)

print(resp.output_text)
```

### 6.3 OpenAI 与 Claude 的 MCP 差异

| 维度 | Claude (Anthropic) | OpenAI |
|------|-------------------|--------|
| **本地 stdio** | ✅ 原生支持 | ❌ 仅远程 HTTP |
| **远程 HTTP** | ✅ Streamable HTTP + SSE | ✅ Streamable HTTP |
| **认证方式** | 可选（本地无需认证） | OAuth 2.1 强制 |
| **Dynamic Client Registration** | 可选 | 强制 |
| **Bearer Token** | 支持 | 不支持 |
| **工具审批** | 首次使用时提示 | 可配置 `require_approval` |
| **配置方式** | JSON 文件 / CLI | 代码中声明 / ChatGPT Settings |

### 6.4 对 Capture 的启示

如果 Capture 要同时支持 Claude 和 OpenAI：
- **stdio 模式**：仅 Claude/Claude Code 支持，OpenAI 不支持
- **HTTP 模式**：两者都支持，但 OpenAI 强制 OAuth 2.1
- **推荐策略**：
  - 本地开发/Claude 用户：stdio（零配置）
  - OpenAI/远程访问：Streamable HTTP + OAuth 2.1

---

## 七、对接 Kimi MCP

### 7.1 Kimi 的 MCP 支持现状

截至 2026-06，公开信息中 **Kimi 官方尚未宣布原生 MCP 客户端支持**。但 MCP 是开放标准，理论上任何实现 JSON-RPC 2.0 + 标准传输的客户端都可以连接 MCP Server。

### 7.2 Kimi 的潜在对接方式

**方式 A：通过 Kimi API + 自定义 MCP Client（如果 Kimi 支持 Function Calling）**

```python
# 假设 Kimi API 支持 function calling（类似 OpenAI）
# Capture 需要提供 OpenAPI 风格的 function schema

functions = [
    {
        "name": "capture_list_tasks",
        "description": "列出指定项目的任务",
        "parameters": {
            "type": "object",
            "properties": {
                "workspace_id": {"type": "string"},
                "project_id": {"type": "string"},
                "status": {"type": "string", "enum": ["todo", "in_progress", "done"]}
            },
            "required": ["workspace_id"]
        }
    }
]

# Kimi 决定在合适的时候调用 function
# 应用层需要实现 MCP Client 来连接 Capture Server
```

**方式 B：Kimi 的 Kimi Code / 桌面客户端未来支持 MCP**

如果 Kimi 未来推出类似 Claude Code 的终端工具或 MCP 客户端支持，配置方式预计类似：

```json
{
  "mcpServers": {
    "capture": {
      "command": "capture",
      "args": ["mcp", "serve"]
    }
  }
}
```

**方式 C：通过第三方 MCP Bridge**

使用社区工具将 MCP Server 桥接到 Kimi API：

```
Kimi API ←→ MCP Bridge (开源工具) ←→ Capture MCP Server (stdio/HTTP)
```

### 7.3 对 Capture 的建议

| 策略 | 说明 |
|------|------|
| **当前** | 专注 Claude 生态（stdio + HTTP），Claude 的 MCP 支持最成熟 |
| **近期** | 同时支持 OpenAI（HTTP + OAuth 2.1） |
| **远期** | 监控 Kimi/国内大模型的 MCP 动态，保持协议兼容即可对接 |
| **通用** | Capture 作为 MCP Server 保持"供应商中立"，任何实现标准协议的客户端都能连接 |

---

## 八、Capture v2 作为 MCP Server 的实现方案

### 8.1 总体架构

```
┌──────────────────────────────────────────────────────────────────────┐
│                        Capture v2 MCP Server                          │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │                    MCP Protocol Handler                        │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │   │
│  │  │ Initialize│  │ Tools    │  │ Resources│  │ Prompts      │  │   │
│  │  │ Handler   │  │ Handler  │  │ Handler  │  │ Handler      │  │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────┘  │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                          │                                            │
│  ┌───────────────────────┼───────────────────────────────────────┐   │
│  │                   Transport Layer                                │   │
│  │  ┌──────────────┐    ┌──────────────────┐    ┌──────────────┐  │   │
│  │  │ StdioTransport│    │ HTTPTransport    │    │ (extensible) │  │   │
│  │  │              │    │ (Streamable HTTP)│    │              │  │   │
│  │  │ stdin/stdout │    │ POST /mcp        │    │ WebSocket?   │  │   │
│  │  └──────────────┘    └──────────────────┘    └──────────────┘  │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                          │                                            │
│  ┌───────────────────────┼───────────────────────────────────────┐   │
│  │                   Service Layer (复用现有)                        │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │   │
│  │  │Workspace │  │ Project  │  │ Task     │  │ Spec         │  │   │
│  │  │Service   │  │ Service  │  │ Service  │  │ Service      │  │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────┘  │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                          │                                            │
│  ┌───────────────────────┼───────────────────────────────────────┐   │
│  │                   Store Layer (复用现有)                          │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │   │
│  │  │Markdown  │  │ SQLite   │  │ Workspace│  │ Project      │  │   │
│  │  │Store     │  │ Store    │  │ Store    │  │ Store        │  │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────┘  │   │
│  └───────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘
```

### 8.2 暴露的 Tools 设计

```go
package mcp

var captureTools = []Tool{
    {
        Name:        "capture_list_tasks",
        Description: "列出指定工作空间和项目的任务",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "workspace_id": map[string]string{
                    "type":        "string",
                    "description": "Workspace ID，如 'default'",
                },
                "project_id": map[string]string{
                    "type":        "string",
                    "description": "Project ID/slug，如 'innate-capture'",
                },
                "status": map[string]interface{}{
                    "type":        "string",
                    "description": "按状态筛选",
                    "enum":        []string{"todo", "in_progress", "done", "cancelled"},
                },
                "stage": map[string]interface{}{
                    "type":        "string",
                    "description": "按阶段筛选",
                    "enum":        []string{"inbox", "mindstorm", "analysis", "planning", "prd", "tasks", "dispatch", "execution", "review"},
                },
                "priority": map[string]interface{}{
                    "type":        "string",
                    "description": "按优先级筛选",
                    "enum":        []string{"high", "medium", "low"},
                },
            },
            "required": []string{"workspace_id"},
        },
        Handler: handleListTasks,
    },
    {
        Name:        "capture_create_task",
        Description: "在指定项目下创建新任务",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "project_id":  map[string]string{"type": "string"},
                "title":       map[string]string{"type": "string"},
                "description": map[string]string{"type": "string"},
                "priority": map[string]interface{}{
                    "type": "string",
                    "enum": []string{"high", "medium", "low"},
                },
                "stage": map[string]interface{}{
                    "type": "string",
                    "enum": []string{"inbox", "mindstorm", "analysis", "planning", "prd", "tasks", "dispatch", "execution", "review"},
                },
                "tags": map[string]interface{}{
                    "type": "array",
                    "items": map[string]string{"type": "string"},
                },
            },
            "required": []string{"project_id", "title"},
        },
        Handler: handleCreateTask,
    },
    {
        Name:        "capture_update_task_status",
        Description: "更新任务状态",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "task_id": map[string]string{"type": "string"},
                "status": map[string]interface{}{
                    "type": "string",
                    "enum": []string{"todo", "in_progress", "done", "cancelled", "archived"},
                },
            },
            "required": []string{"task_id", "status"},
        },
        Handler: handleUpdateTaskStatus,
    },
    {
        Name:        "capture_get_task_context",
        Description: "获取任务完整上下文（包括关联的 Spec、相关任务、GitHub 同步状态）",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "task_id": map[string]string{"type": "string"},
            },
            "required": []string{"task_id"},
        },
        Handler: handleGetTaskContext,
    },
    {
        Name:        "capture_list_projects",
        Description: "列出工作空间下的所有项目",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "workspace_id": map[string]string{"type": "string"},
            },
            "required": []string{"workspace_id"},
        },
        Handler: handleListProjects,
    },
    {
        Name:        "capture_sync_github",
        Description: "触发指定项目的 GitHub Project 同步",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "project_id": map[string]string{"type": "string"},
                "direction": map[string]interface{}{
                    "type":    "string",
                    "enum":    []string{"pull", "push", "bidirectional"},
                    "default": "pull",
                },
            },
            "required": []string{"project_id"},
        },
        Handler: handleSyncGitHub,
    },
}
```

### 8.3 暴露的 Resources 设计

```go
var captureResources = []Resource{
    {
        URI:         "capture://workspaces",
        Name:        "Workspaces",
        Description: "所有工作空间列表",
        MimeType:    "application/json",
    },
    {
        URI:         "capture://workspace/{workspace_id}/projects",
        Name:        "Projects",
        Description: "指定工作空间下的项目",
        MimeType:    "application/json",
    },
    {
        URI:         "capture://workspace/{workspace_id}/tasks",
        Name:        "Tasks",
        Description: "指定工作空间下的任务",
        MimeType:    "application/json",
    },
    {
        URI:         "capture://task/{task_id}",
        Name:        "Task Detail",
        Description: "单个任务的完整 Markdown 内容",
        MimeType:    "text/markdown",
    },
    {
        URI:         "capture://spec/{spec_id}",
        Name:        "Spec Detail",
        Description: "Spec 文档的完整 Markdown 内容",
        MimeType:    "text/markdown",
    },
}
```

### 8.4 暴露的 Prompts 设计

```go
var capturePrompts = []Prompt{
    {
        Name:        "brainstorm-scamper",
        Description: "使用 SCAMPER 方法分析任务或功能",
        Arguments: []PromptArgument{
            {Name: "task_id", Description: "要分析的任务 ID", Required: true},
        },
    },
    {
        Name:        "spec-grill-me",
        Description: "对 Spec 进行压力测试式审问",
        Arguments: []PromptArgument{
            {Name: "spec_id", Description: "要审问的 Spec ID", Required: true},
        },
    },
    {
        Name:        "task-first-principles",
        Description: "用第一性原理拆解任务",
        Arguments: []PromptArgument{
            {Name: "task_id", Description: "要拆解的任务 ID", Required: true},
        },
    },
    {
        Name:        "multi-role-review",
        Description: "多角色评审（PM/Engineer/Designer）",
        Arguments: []PromptArgument{
            {Name: "spec_id", Description: "要评审的 Spec ID", Required: true},
        },
    },
}
```

---

## 九、命令行接口设计

```bash
# 启动 MCP Server（stdio 模式，默认）
capture mcp serve

# 启动 MCP Server（stdio 模式，显式）
capture mcp serve --transport stdio

# 启动 MCP Server（Streamable HTTP 模式）
capture mcp serve --transport http --port 8080

# 启动 MCP Server（绑定到特定地址）
capture mcp serve --transport http --host 127.0.0.1 --port 8080

# 检查 MCP Server 状态
capture mcp status

# 测试 MCP Server（本地自检）
capture mcp test
```

---

## 十、实现优先级建议

| 阶段 | 内容 | 工期 | 理由 |
|------|------|------|------|
| **Phase 1** | stdio transport + Tools/List + Tools/Call | 3-5 天 | Claude Desktop/Claude Code 原生支持，零配置 |
| **Phase 2** | Resources + Prompts | 2-3 天 | 增强 AI 对 Capture 数据的访问能力 |
| **Phase 3** | Streamable HTTP transport | 3-4 天 | 支持远程/多租户场景，OpenAI 兼容 |
| **Phase 4** | OAuth 2.1 + Session 管理 | 4-5 天 | OpenAI 强制要求，生产环境必需 |
| **Phase 5** | Notifications + Progress | 2-3 天 | 长任务（如同步）的进度反馈 |

---

## 十一、参考资源

| 资源 | URL |
|------|-----|
| MCP 官方规范 | `modelcontextprotocol.io/specification` |
| MCP TypeScript SDK | `github.com/modelcontextprotocol/typescript-sdk` |
| MCP Python SDK | `github.com/modelcontextprotocol/python-sdk` |
| Claude MCP 配置指南 | `systemprompt.io/guides/claude-code-mcp-servers-extensions` |
| OpenAI MCP 指南 | `developers.openai.com/api/docs/guides/tools-connectors-mcp` |
| MCP 传输层详解 | `learn.engineering.vips.edu/mcp/mcp-transports-stdio-vs-sse` |
| MCP Inspector（调试工具） | `github.com/modelcontextprotocol/inspector` |
