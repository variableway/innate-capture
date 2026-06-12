# Agent Runtime Specification

> **Version**: 1.0  
> **Status**: Draft  
> **Based on**: Multica Daemon + Agent System Analysis

---

## 1. Overview

Agent Runtime 是本地 AI Agent 的执行环境，负责任务领取、执行和结果上报。支持多机器部署，实现任务的分布式执行。

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Agent Runtime Architecture                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                      Central Server (Optional)                   │    │
│  │  - Task Queue Management                                        │    │
│  │  - Runtime Registration                                         │    │
│  │  - Result Aggregation                                           │    │
│  └───────────────────────────┬─────────────────────────────────────┘    │
│                              │                                           │
│         ┌────────────────────┼────────────────────┐                      │
│         │                    │                    │                      │
│         ▼                    ▼                    ▼                      │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐              │
│  │   Daemon    │      │   Daemon    │      │   Daemon    │              │
│  │  (MacBook)  │      │  (Windows)  │      │  (Linux)    │              │
│  │             │      │             │      │             │              │
│  │ ┌─────────┐ │      │ ┌─────────┐ │      │ ┌─────────┐ │              │
│  │ │ Agent   │ │      │ │ Agent   │ │      │ │ Agent   │ │              │
│  │ │(claude) │ │      │ │(codex)  │ │      │ │(opencode│ │              │
│  │ └────┬────┘ │      │ └────┬────┘ │      │ └────┬────┘ │              │
│  │ ┌─────────┐ │      │ ┌─────────┐ │      │ ┌─────────┐ │              │
│  │ │ Agent   │ │      │ │ Agent   │ │      │ │ Agent   │ │              │
│  │ │(codex)  │ │      │ │(custom) │ │      │ │(claude) │ │              │
│  │ └─────────┘ │      │ └─────────┘ │      │ └─────────┘ │              │
│  └─────────────┘      └─────────────┘      └─────────────┘              │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Configuration

### 3.1 Agent Configuration

```yaml
# ~/.capture/agents.yaml
version: "1.0"

# Machine identification
machine:
  id: "macbook-pro-001"      # Unique machine ID
  name: "MacBook Pro"         # Human readable name
  platform: "darwin"          # darwin | windows | linux

# Agent definitions
agents:
  claude-primary:
    type: "claude"            # claude | codex | opencode
    name: "Claude"
    executable: "/usr/local/bin/claude"
    model: "claude-sonnet-4"
    max_concurrent: 3
    work_dir: "~/.capture/workspaces/{workspace}/claude"
    env:
      CLAUDE_CODE_DEBUG: "1"
    
  codex-secondary:
    type: "codex"
    name: "Codex"
    executable: "codex"
    model: "gpt-4o"
    max_concurrent: 2
    work_dir: "~/.capture/workspaces/{workspace}/codex"
    
  opencode-research:
    type: "opencode"
    name: "OpenCode"
    executable: "opencode"
    model: "gemini-2.5-pro"
    max_concurrent: 2
    work_dir: "~/.capture/workspaces/{workspace}/opencode"

# Runtime settings
runtime:
  poll_interval: 30s
  heartbeat_interval: 30s
  task_timeout: 30m
  health_check_port: 19514
```

### 3.2 Runtime Configuration

```go
type RuntimeConfig struct {
    // Identity
    MachineID    string `yaml:"machine_id"`
    MachineName  string `yaml:"machine_name"`
    Platform     string `yaml:"platform"`
    
    // Server connection (optional, for distributed mode)
    ServerURL    string `yaml:"server_url,omitempty"`
    AuthToken    string `yaml:"auth_token,omitempty"`
    
    // Runtime settings
    PollInterval      time.Duration `yaml:"poll_interval"`
    HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
    TaskTimeout       time.Duration `yaml:"task_timeout"`
    HealthPort        int           `yaml:"health_check_port"`
    
    // Agent definitions
    Agents map[string]AgentConfig `yaml:"agents"`
}

type AgentConfig struct {
    Type           string            `yaml:"type"`
    Name           string            `yaml:"name"`
    Executable     string            `yaml:"executable"`
    Model          string            `yaml:"model,omitempty"`
    MaxConcurrent  int               `yaml:"max_concurrent"`
    WorkDir        string            `yaml:"work_dir"`
    Env            map[string]string `yaml:"env,omitempty"`
    Skills         []string          `yaml:"skills,omitempty"` // Skill file paths
}
```

---

## 4. Daemon Service

### 4.1 Daemon Interface

```go
// Daemon is the local agent runtime
type Daemon interface {
    // Lifecycle
    Start(ctx context.Context) error
    Stop() error
    Restart() error
    Status() DaemonStatus
    
    // Agent management
    RegisterAgent(agent AgentConfig) error
    UnregisterAgent(agentID string) error
    ListAgents() []AgentStatus
    
    // Task management
    ClaimTask(ctx context.Context) (*Task, error)
    StartTask(ctx context.Context, taskID string) error
    CompleteTask(ctx context.Context, taskID string, result TaskResult) error
    FailTask(ctx context.Context, taskID string, err error) error
    CancelTask(ctx context.Context, taskID string) error
    
    // Health & monitoring
    Health() HealthStatus
    GetMetrics() Metrics
}

type DaemonStatus struct {
    State       string    `json:"state"` // running, stopped, error
    StartedAt   time.Time `json:"started_at"`
    Agents      []AgentStatus `json:"agents"`
    ActiveTasks []TaskStatus  `json:"active_tasks"`
}

type AgentStatus struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Type          string    `json:"type"`
    State         string    `json:"state"` // idle, working, error
    RunningTasks  int       `json:"running_tasks"`
    MaxConcurrent int       `json:"max_concurrent"`
    LastHeartbeat time.Time `json:"last_heartbeat"`
}
```

### 4.2 CLI Commands

```bash
# Daemon management
capture daemon start [--foreground]
capture daemon stop
capture daemon status
capture daemon logs [-f] [-n 100]
capture daemon restart

# Agent management
capture agent list
capture agent status <agent-id>
capture agent ping <agent-id>
```

---

## 5. Task Execution

### 5.1 Task Model

```go
type Task struct {
    ID          string    `json:"id"`
    IssueID     string    `json:"issue_id"`      // Linked issue
    AgentID     string    `json:"agent_id"`      // Assigned agent
    RuntimeID   string    `json:"runtime_id"`    // Machine runtime
    
    // Content
    Title       string    `json:"title"`
    Prompt      string    `json:"prompt"`        // Full prompt for agent
    Context     TaskContext `json:"context"`
    
    // Status
    Status      TaskStatus  `json:"status"`      // pending, running, completed, failed, cancelled
    Priority    int         `json:"priority"`
    
    // Timestamps
    CreatedAt   time.Time   `json:"created_at"`
    StartedAt   *time.Time  `json:"started_at,omitempty"`
    CompletedAt *time.Time  `json:"completed_at,omitempty"`
    
    // Execution
    WorkDir     string      `json:"work_dir,omitempty"`
    SessionID   string      `json:"session_id,omitempty"`  // Agent session for resume
    Result      *TaskResult `json:"result,omitempty"`
}

type TaskContext struct {
    Issue       IssueSummary    `json:"issue"`
    Skills      []SkillData     `json:"skills,omitempty"`
    Repos       []RepoInfo      `json:"repos,omitempty"`
    PriorContext *PriorContext  `json:"prior_context,omitempty"` // For resume
}

type TaskResult struct {
    Status     string    `json:"status"`      // success, failed, blocked
    Output     string    `json:"output"`      // Text summary
    Artifacts  []Artifact `json:"artifacts"`  // Generated files
    GitCommit  string    `json:"git_commit,omitempty"`
    GitBranch  string    `json:"git_branch,omitempty"`
    Duration   int       `json:"duration_seconds"`
}
```

### 5.2 Execution Flow

```
┌─────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Claim  │────►│  Prepare │────►│  Execute │────►│  Report  │
└─────────┘     └──────────┘     └──────────┘     └──────────┘
     │               │               │               │
     ▼               ▼               ▼               ▼
Poll task      Setup workdir    Run agent        Update issue
from queue     Checkout repo    Stream output    Archive to wiki
               Load skills      Capture result   Notify user
```

### 5.3 Agent Backend Interface

```go
// Backend is the unified interface for AI agents
type Backend interface {
    // Execute runs a prompt and returns a session for streaming results
    Execute(ctx context.Context, prompt string, opts ExecOptions) (*Session, error)
    
    // Health check
    Health() error
    
    // Version info
    Version() (string, error)
}

type ExecOptions struct {
    Cwd             string
    Model           string
    SystemPrompt    string
    MaxTurns        int
    Timeout         time.Duration
    ResumeSessionID string
    Env             map[string]string
}

type Session struct {
    Messages <-chan Message  // Stream events
    Result   <-chan Result   // Final outcome
}

type Message struct {
    Type    MessageType      // text, thinking, tool_use, tool_result, error
    Content string
    Tool    string           // For tool_use/tool_result
    Input   map[string]any   // For tool_use
    Output  string           // For tool_result
}

type Result struct {
    Status     string // completed, failed, timeout, aborted
    Output     string
    Error      string
    DurationMs int64
    SessionID  string
}
```

---

## 6. Communication Protocol

### 6.1 Daemon ↔ Server (Optional)

```go
// For distributed mode with central server

// Runtime registration
type RegisterRequest struct {
    MachineID   string         `json:"machine_id"`
    MachineName string         `json:"machine_name"`
    Agents      []AgentInfo    `json:"agents"`
}

type RegisterResponse struct {
    RuntimeID   string `json:"runtime_id"`
    AuthToken   string `json:"auth_token"`
}

// Task claiming
type ClaimTaskRequest struct {
    RuntimeID string `json:"runtime_id"`
    AgentID   string `json:"agent_id"`
}

type ClaimTaskResponse struct {
    Task    *Task  `json:"task,omitempty"`
    HasTask bool   `json:"has_task"`
}

// Heartbeat
type HeartbeatRequest struct {
    RuntimeID     string        `json:"runtime_id"`
    AgentStatuses []AgentStatus `json:"agent_statuses"`
    ActiveTasks   []TaskStatus  `json:"active_tasks"`
}

type HeartbeatResponse struct {
    Commands []ServerCommand `json:"commands,omitempty"` // Cancel, config update, etc.
}
```

### 6.2 Local Mode (No Server)

In local-only mode, the daemon reads tasks directly from local SQLite queue.

```go
type LocalTaskQueue interface {
    // Poll for pending tasks
    Poll(ctx context.Context, agentID string) (*Task, error)
    
    // Update task status
    UpdateStatus(taskID string, status TaskStatus) error
    
    // Store result
    StoreResult(taskID string, result TaskResult) error
}
```

---

## 7. Health & Monitoring

### 7.1 Health Endpoint

```http
GET http://localhost:19514/health

Response:
{
    "status": "healthy",
    "version": "1.0.0",
    "machine_id": "macbook-pro-001",
    "started_at": "2026-04-06T10:00:00Z",
    "agents": [
        {
            "id": "claude-primary",
            "name": "Claude",
            "state": "idle",
            "running_tasks": 0,
            "max_concurrent": 3
        }
    ],
    "active_tasks": 0
}
```

### 7.2 Metrics

```go
type Metrics struct {
    TasksTotal     int64     `json:"tasks_total"`
    TasksSucceeded int64     `json:"tasks_succeeded"`
    TasksFailed    int64     `json:"tasks_failed"`
    AverageDuration float64  `json:"average_duration_seconds"`
    Uptime         int64     `json:"uptime_seconds"`
}
```

---

## 8. Security Considerations

1. **Workdir Isolation**: Each task gets isolated working directory
2. **Env Sanitization**: Sensitive env vars are filtered before passing to agents
3. **Token Management**: Auth tokens stored in OS keychain, not plain files
4. **Network**: Agents run with restricted network access (optional)

---

## 9. Related Specs

- [Issue Model Spec](./issue-model-spec.md)
- [Wiki Integration Spec](./wiki-integration-spec.md)
