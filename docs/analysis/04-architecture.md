# Capture v2 架构设计

> 设计时间：2026-06-04
> 设计原则：本地优先、终端原生、Markdown 原生、Git 友好、渐进复杂度、远端可插拔

---

## 一、系统架构总览

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                              交互层 (Presentation)                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   CLI 命令   │  │  TUI 看板   │  │  飞书 Bot   │  │   Spec/PRD 文档      │  │
│  │  (Cobra)    │  │ (bubbletea) │  │ (WS/HTTP)   │  │   (Markdown)        │  │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────────────────────┘  │
└─────────┼────────────────┼────────────────┼──────────────────────────────────┘
          │                │                │
          └────────────────┴────────────────┘
                              │
┌─────────────────────────────┼─────────────────────────────────────────────────┐
│                         服务层 (Service)                                       │
│  ┌──────────────────────────┴─────────────────────────────────────────────┐   │
│  │                        TaskService                                        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │   │
│  │  │  Workspace  │  │   Project   │  │    Task     │  │    Spec     │   │   │
│  │  │   Service   │  │   Service   │  │   Service   │  │   Service   │   │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                      SyncService (Remote Orchestrator)                  │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐  │   │
│  │  │  GitHubRemote   │  │  FeishuRemote   │  │    Remote Interface     │  │   │
│  │  │   (GraphQL)     │  │   (REST/WS)     │  │    (可插拔扩展)          │  │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────────┘
                                    │
┌───────────────────────────────────┼───────────────────────────────────────────┐
│                              存储层 (Storage)                                  │
│  ┌────────────────────────────────┴────────────────────────────────────────┐  │
│  │                         Store Interface                                  │  │
│  │  ┌──────────────────────────────────────────────────────────────────┐   │  │
│  │  │                     DualStore (per Workspace)                     │   │  │
│  │  │  ┌──────────────────────────┐  ┌──────────────────────────────┐  │   │  │
│  │  │  │   MarkdownStore          │  │      SQLiteStore             │  │   │  │
│  │  │  │ (Source of Truth)        │  │   (Query Index)              │  │   │  │
│  │  │  │                          │  │                              │  │   │  │
│  │  │  │  ~/.capture/             │  │  capture.db                  │  │   │  │
│  │  │  │  └── workspaces/         │  │  ├── tasks                   │  │   │  │
│  │  │  │      └── <ws>/           │  │  ├── tags                    │  │   │  │
│  │  │  │          └── projects/   │  │  ├── sync_logs               │  │   │  │
│  │  │  │              └── <proj>/ │  │  └── project_configs         │  │   │  │
│  │  │  │                  ├──tasks│  │                              │  │   │  │
│  │  │  │                  └──specs│  │                              │  │   │  │
│  │  │  └──────────────────────────┘  └──────────────────────────────┘  │   │  │
│  │  └──────────────────────────────────────────────────────────────────┘   │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────────────────────────┘
```

---

## 二、核心概念模型

### 2.1 概念层次

```
Workspace（工作空间）
    ├── 配置：workspace.yaml（名称、默认 Project、远端配置）
    ├── SQLite：capture.db（索引）
    ├── Projects（项目）
    │   ├── 配置：project.yaml（GitHub repo、飞书表、自定义字段）
    │   ├── Specs（规格文档）
    │   │   └── SPEC-001-feature-name.md
    │   └── Tasks（任务）
    │       └── 2026/06/TASK-00001.md
    └── Templates（模板）
        └── spec-default.md
```

### 2.2 实体关系

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│  Workspace  │ 1───N │   Project   │ 1───N │    Task     │
│             │       │             │       │             │
│  - id       │       │  - id       │       │  - id       │
│  - name     │       │  - name     │       │  - title    │
│  - path     │       │  - repo_url │       │  - status   │
│  - config   │       │  - github   │       │  - stage    │
└─────────────┘       │  - feishu   │       │  - project  │
                      └─────────────┘       │  - external │
                             │              └─────────────┘
                             │ 1───N
                      ┌─────────────┐
                      │    Spec     │
                      │             │
                      │  - id       │
                      │  - title    │
                      │  - project  │
                      │  - tasks[]  │
                      └─────────────┘
```

### 2.3 关键设计决策

| 决策 | 选择 | 理由 |
|------|------|------|
| Workspace 是物理目录还是逻辑概念？ | **物理目录** | 支持 Git 版本化、支持跨设备同步、数据透明可审计 |
| Project 是否必须关联 GitHub repo？ | **不必须，但推荐** | 支持非代码项目（如产品规划、学习计划） |
| Task 是否必须属于 Project？ | **v2 是，v1 兼容否** | v1 遗留任务归入 default Project |
| Spec 文档存储在哪里？ | **Project 目录下的 specs/** | 与 Task 同级，便于关联和导航 |
| 远端配置在哪一级？ | **Workspace 级统一配置，Project 级可覆盖** | 灵活且不过度分散 |

---

## 三、存储层架构

### 3.1 目录结构（v2）

```
~/.capture/                                 # 全局根目录
├── config.yaml                             # 全局配置：默认编辑器、远端凭据等
│                                           #
├── workspaces/                             # 所有工作空间
│   └── default/                            # 默认工作空间（自动创建）
│       ├── workspace.yaml                  # Workspace 配置
│       │   ├── name: "Personal"
│       │   ├── default_project: "capture"
│       │   └── remotes:
│       │       ├── github:
│       │       │   ├── token_source: "gh_cli"
│       │       │   └── default_org: "variableway"
│       │       └── feishu:
│       │           ├── app_id: "${FEISHU_APP_ID}"
│       │           └── app_secret: "${FEISHU_APP_SECRET}"
│       │
│       ├── capture.db                      # Workspace 级 SQLite 数据库
│       │
│       ├── projects/                       # 项目列表
│       │   └── innate-capture/             # 项目目录（slug 命名）
│       │       ├── project.yaml            # 项目配置
│       │       │   ├── name: "Capture"
│       │       │   ├── description: "..."
│       │       │   ├── github:
│       │       │   │   ├── repo: "variableway/innate-capture"
│       │       │   │   └── project_number: 1
│       │       │   ├── feishu:
│       │       │   │   ├── bitable_app_token: "..."
│       │       │   │   └── bitable_table_id: "..."
│       │       │   └── custom_fields: [...]
│       │       │
│       │       ├── specs/                  # 规格文档
│       │       │   └── SPEC-001-workspace.md
│       │       │
│       │       └── tasks/                  # 任务文件
│       │           └── 2026/
│       │               └── 06/
│       │                   └── TASK-00001.md
│       │
│       └── templates/                      # 工作空间级模板
│           └── spec-default.md
│
└── .backup-v1-20260604/                    # v1 迁移备份（如需要）
```

### 3.2 DualStore 演进

v1 的 `DualStore` 在 v2 中变为 **Workspace 级** 的存储管理器：

```go
// WorkspaceStore 管理一个 Workspace 下的所有存储
type WorkspaceStore struct {
    workspaceID string
    workspacePath string
    markdown    *MarkdownStore      // per-project
    sqlite      *SQLiteStore        // per-workspace
    projects    map[string]*ProjectStore
}

type ProjectStore struct {
    projectID   string
    projectPath string
    markdown    *MarkdownStore      // 本项目 Task 存储
    specs       *SpecStore          // 本项目 Spec 文档存储
}
```

**关键设计**：
- 每个 Workspace 一个 SQLite 数据库（避免全局锁、支持 Workspace 独立迁移）
- Markdown 文件按 Project 分目录（支持 Git 子模块式的项目管理）
- 全局 `StoreManager` 负责 Workspace 的创建、切换、枚举

### 3.3 SQLite Schema 演进

```sql
-- workspaces 表（在全局 SQLite 中，或每个 workspace 独立 DB）
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- projects 表
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    github_repo TEXT,
    github_project_number INTEGER,
    feishu_bitable_app_token TEXT,
    feishu_bitable_table_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(workspace_id, slug)
);

-- tasks 表（v1 兼容 + 新增字段）
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    title TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'todo',
    stage TEXT NOT NULL DEFAULT 'inbox',
    priority TEXT DEFAULT 'medium',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    source TEXT DEFAULT 'cli',
    file_path TEXT NOT NULL,
    feishu_record_id TEXT DEFAULT '',
    external_id TEXT DEFAULT '',       -- GitHub Issue/PR ID
    external_url TEXT DEFAULT '',      -- GitHub Issue/PR URL
    assigned_agent TEXT DEFAULT '',
    assigned_repository TEXT DEFAULT '',
    assigned_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- specs 表
CREATE TABLE specs (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    file_path TEXT NOT NULL,
    status TEXT DEFAULT 'draft',       -- draft, reviewing, approved, implemented
    linked_tasks TEXT,                 -- JSON array of task IDs
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- sync_logs 表（扩展）
CREATE TABLE sync_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT,
    sync_type TEXT NOT NULL,           -- github, feishu
    direction TEXT NOT NULL,           -- push, pull, bidirectional
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    status TEXT,
    records_count INTEGER DEFAULT 0,
    error_message TEXT
);
```

---

## 四、远端同步架构

### 4.1 Remote 接口设计

```go
// Remote 定义所有远端同步源的统一接口
type Remote interface {
    Name() string                      // "github", "feishu"
    
    // Project 级操作
    ListProjects(ctx context.Context) ([]RemoteProject, error)
    GetProject(ctx context.Context, id string) (*RemoteProject, error)
    
    // Task/Item 级操作
    ListItems(ctx context.Context, projectID string) ([]RemoteItem, error)
    GetItem(ctx context.Context, projectID, itemID string) (*RemoteItem, error)
    CreateItem(ctx context.Context, projectID string, item *RemoteItem) (*RemoteItem, error)
    UpdateItem(ctx context.Context, projectID string, item *RemoteItem) (*RemoteItem, error)
    DeleteItem(ctx context.Context, projectID, itemID string) error
    
    // 同步状态
    SyncStatus(ctx context.Context, projectID string) (*SyncStatus, error)
}

type RemoteProject struct {
    ID          string
    Name        string
    Description string
    URL         string
    ExternalID  string          // 远端原始 ID
}

type RemoteItem struct {
    ID         string
    ExternalID string          // GitHub Issue node_id / Feishu record_id
    Title      string
    Status     string
    Priority   string
    URL        string
    UpdatedAt  time.Time
    Raw        map[string]interface{}  // 远端原始数据（用于冲突检测）
}
```

### 4.2 GitHub Remote 架构

```
┌─────────────────────────────────────────────┐
│           GitHubRemote                       │
│  ┌───────────────────────────────────────┐  │
│  │  Authentication                       │  │
│  │  ├── gh auth status (优先)            │  │
│  │  ├── GITHUB_TOKEN 环境变量            │  │
│  │  └── ~/.config/gh/hosts.yml           │  │
│  └───────────────────────────────────────┘  │
│  ┌───────────────────────────────────────┐  │
│  │  GraphQL Client (githubv4)            │  │
│  │  ├── query ProjectItems               │  │
│  │  ├── query ProjectFields              │  │
│  │  └── mutation UpdateItem              │  │
│  └───────────────────────────────────────┘  │
│  ┌───────────────────────────────────────┐  │
│  │  Field Mapper                         │  │
│  │  ├── GitHub Status → TaskStatus       │  │
│  │  ├── GitHub Priority → TaskPriority   │  │
│  │  └── Custom Fields → Task.Tags/...    │  │
│  └───────────────────────────────────────┘  │
│  ┌───────────────────────────────────────┐  │
│  │  Cache Layer                          │  │
│  │  └── ~/.capture/cache/github/...      │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

### 4.3 同步策略

**Phase 1（MVP）：只读拉取（Pull Only）**
```
GitHub Project ──Pull──► Local SQLite/Markdown
    │                          │
    └──── 单向同步，本地只读缓存 ────┘
```

**Phase 2：双向同步（Bidirectional）**
```
GitHub Project ◄──Sync──► Local SQLite/Markdown
    │                          │
    ├── 冲突检测：对比 UpdatedAt 和 Raw hash
    ├── 冲突策略：
    │   ├── 默认：远端优先（GitHub wins）
    │   ├── 可选：本地优先（Local wins）
    │   └── 手动：冲突时提示用户选择
    └── 批量同步：支持按 Project 全量/增量同步
```

**Phase 3：飞书同步升级**
```
GitHub Project ◄──Sync──► Local ◄──Sync──► 飞书 Bitable
    │                         │
    └──── Capture 作为中间协调者 ─────┘
```

### 4.4 冲突解决机制

```go
type ConflictResolver interface {
    // 检测本地和远端是否有冲突
    Detect(local, remote *model.Task) (*Conflict, error)
    
    // 解决冲突
    Resolve(conflict *Conflict, strategy ResolutionStrategy) (*model.Task, error)
}

type Conflict struct {
    TaskID      string
    LocalSHA    string      // 本地内容 hash
    RemoteSHA   string      // 远端内容 hash
    LocalTime   time.Time
    RemoteTime  time.Time
    Fields      []string    // 冲突的字段列表
}

type ResolutionStrategy int
const (
    StrategyRemoteWins ResolutionStrategy = iota   // 远端优先
    StrategyLocalWins                               // 本地优先
    StrategyNewestWins                              // 最新修改优先
    StrategyManual                                  // 交互式选择
)
```

---

## 五、配置架构

### 5.1 配置分层

```
┌─────────────────────────────────────────────┐
│  全局配置：~/.capture/config.yaml            │
│  ├── app                                    │
│  ├── defaults                               │
│  ├── editor                                 │
│  └── remote_credentials (GitHub, Feishu)    │
├─────────────────────────────────────────────┤
│  Workspace 配置：<ws>/workspace.yaml         │
│  ├── name, description                      │
│  ├── default_project                        │
│  └── remote_overrides                       │
├─────────────────────────────────────────────┤
│  Project 配置：<proj>/project.yaml           │
│  ├── name, description                      │
│  ├── github: repo, project_number           │
│  ├── feishu: bitable config                 │
│  ├── custom_fields                          │
│  └── templates                              │
└─────────────────────────────────────────────┘
```

### 5.2 配置合并规则

```go
// 优先级：Project > Workspace > Global（高优先级覆盖低优先级）
// 示例：GitHub token
// 1. 先读取 ~/.capture/config.yaml 中的 github.token
// 2. 如果 Workspace 有覆盖，使用 Workspace 的值
// 3. 如果 Project 有覆盖，使用 Project 的值
// 4. 如果环境变量 GITHUB_TOKEN 存在，最高优先级
```

---

## 六、TUI 架构演进

### 6.1 视图层次

```
App (bubbletea Model)
├── WorkspaceSelector        # 选择当前 Workspace（如果只有一个则跳过）
├── ProjectSelector          # 选择当前 Project
├── MainView
│   ├── KanbanView           # 按 Status 分三列（与 v1 兼容）
│   ├── StageView            # 按 Stage 分多列（inbox→mindstorm→...→review）
│   ├── ListView             # 列表视图，支持搜索/过滤
│   └── SpecView             # Spec 文档浏览器
├── DetailView               # Task/Spec 详情
└── SyncStatusView           # 远端同步状态面板
```

### 6.2 状态管理

```go
type App struct {
    // 导航状态
    currentView    ViewState
    currentWorkspace string
    currentProject   string
    
    // 服务依赖
    workspaceSvc   *service.WorkspaceService
    projectSvc     *service.ProjectService
    taskSvc        *service.TaskService
    specSvc        *service.SpecService
    syncSvc        *service.SyncService
    
    // 视图组件
    workspaceSelector *components.Selector
    projectSelector   *components.Selector
    kanban            *views.Kanban
    taskList          *views.TaskList
    specBrowser       *views.SpecBrowser
    detail            *views.Detail
    
    // 全局状态
    width, height  int
    err            error
}
```

---

## 七、飞书 Bot 架构演进

### 7.1 命令扩展

```
# v1 命令（保持兼容）
记录 <内容>           → 在当前 Project 创建任务
列出                  → 列出当前 Project 的任务
删除 <TASK-ID>        → 删除任务
帮助                  → 显示帮助

# v2 新增命令
项目 列出             → 列出所有 Project
项目 切换 <project>   → 切换当前默认 Project
项目 任务 <project>   → 列出指定 Project 的任务
项目 同步 <project>   → 手动触发指定 Project 的同步
```

### 7.2 消息解析增强

```go
// 识别消息中的 GitHub URL，自动提取信息
type MessageContext struct {
    RawText      string
    ProjectHint  string        // 用户显式指定的 Project
    GitHubURLs   []string      // 消息中包含的 GitHub Issue/PR URL
    Command      string
    Args         []string
}

// 示例：
// 用户发送："记录 修复登录bug https://github.com/org/repo/issues/123"
// Bot 解析：
//   - 创建 Task："修复登录bug"
//   - 关联 external_id 到 GitHub Issue #123
//   - 自动推断 Project 为 "repo"
```

---

## 八、Spec 文档架构

### 8.1 文档模板

```markdown
---
id: SPEC-001
title: "Workspace 管理功能"
project: innate-capture
status: draft          # draft, reviewing, approved, implemented
author: "patrick"
created_at: 2026-06-04
updated_at: 2026-06-04
linked_tasks:
  - TASK-00001
  - TASK-00002
tags: [feature, architecture]
---

## Background

## Goals

## Non-Goals

## Design

## Acceptance Criteria

- [ ] 可以创建 Workspace
- [ ] 可以切换 Workspace
- [ ] ...

## References
```

### 8.2 文档生命周期

```
draft → reviewing → approved → implemented → archived
  │         │           │            │
  │         │           │            └── 所有 linked_tasks 完成
  │         │           └── PRD 评审通过
  │         └── 发起评审（可生成飞书文档或 GitHub Discussion）
  └── `capture spec create` 创建
```

### 8.3 与 Task 的关联

```go
// Spec 文档中的 Acceptance Criteria 可一键生成 Task
type SpecService interface {
    CreateSpec(projectID string, template string) (*Spec, error)
    ParseAcceptanceCriteria(spec *Spec) ([]TaskDraft, error)
    GenerateTasks(spec *Spec) ([]*model.Task, error)
    LinkTask(specID, taskID string) error
}
```

---

## 九、关键技术决策汇总

| 决策项 | 方案 | 备选方案 | 选择理由 |
|--------|------|----------|----------|
| Workspace 存储 | 物理文件目录 | 单 SQLite 数据库 | Git 友好、数据透明 |
| GitHub API | GraphQL v4 | REST v3 | Project v2 只有 GraphQL |
| 认证方式 | gh CLI 优先 | 手动 token | 降低用户配置成本 |
| Task ID 格式 | `TASK-NNNNN`（保持） | `PROJ-TASK-NNNNN` | v1 兼容，Project 通过字段区分 |
| 同步方向 | 只读优先（MVP） | 双向同步 | 降低 MVP 风险 |
| 冲突解决 | 远端优先默认 | 本地优先 | GitHub 通常是团队协作源 |
| Spec 模板 | Go text/template | Helm/YAML 模板 | 标准库、零依赖 |
| TUI 框架 | bubbletea（保持） | tview/fyne | 已有技术债，生态更好 |
