# Capture v2 模块设计

> 设计时间：2026-06-04
> 基于现有 v1 代码库 (`projects/capture/`) 进行增量设计

---

## 一、现有模块结构（v1）

```
projects/capture/
├── main.go
├── cmd/                        # CLI 命令（Cobra）
│   ├── root.go                 # 根命令、配置初始化
│   ├── init.go                 # `capture init`
│   ├── add.go                  # `capture add`
│   ├── list.go                 # `capture list`
│   ├── show.go                 # `capture show`
│   ├── edit.go                 # `capture edit`
│   ├── delete.go               # `capture delete`
│   ├── status.go               # `capture status`
│   ├── stage.go                # `capture stage`
│   ├── assign.go               # `capture assign`
│   ├── kanban.go               # `capture kanban`
│   ├── sync.go                 # `capture sync`
│   ├── bot.go / bot_serve.go   # `capture bot`
│   └── config.go               # `capture config`
├── internal/
│   ├── model/                  # 数据模型
│   │   ├── task.go
│   │   └── config.go
│   ├── store/                  # 存储层
│   │   ├── store.go            # Store interface
│   │   ├── markdown.go         # Markdown 文件存储
│   │   ├── sqlite.go           # SQLite 索引存储
│   │   └── dual.go             # 双写协调
│   ├── service/                # 业务逻辑层
│   │   └── task_service.go
│   ├── config/                 # 配置管理
│   │   └── config.go
│   ├── tui/                    # TUI 界面
│   │   ├── app.go
│   │   ├── keys.go
│   │   └── styles.go
│   ├── bot/                    # 飞书 Bot
│   │   ├── msgparser.go
│   │   ├── dispatcher.go
│   │   ├── webhook.go
│   │   └── websocket.go
│   ├── bitable/                # 飞书多维表格
│   │   ├── client.go
│   │   ├── mapper.go
│   │   └── sync.go
│   └── feishu/                 # 飞书 SDK 封装
│       ├── client.go
│       └── connection.go
└── pkg/
    ├── idgen/                  # ID 生成器
    │   └── idgen.go
    └── frontmatter/            # YAML frontmatter 解析
        └── frontmatter.go
```

---

## 二、v2 模块增量设计

### 2.1 总体变化概览

```diff
 projects/capture/
 ├── main.go
 ├── cmd/
 │   ├── root.go
+│   ├── workspace.go            # `capture workspace *` 命令族
+│   ├── project.go              # `capture project *` 命令族
 │   ├── init.go
 │   ├── add.go                  # 修改：支持 --project 参数
 │   ├── list.go                 # 修改：支持按 workspace/project 过滤
 │   ├── show.go
 │   ├── edit.go
 │   ├── delete.go
 │   ├── status.go
 │   ├── stage.go
 │   ├── assign.go
 │   ├── kanban.go               # 修改：增加 Project 选择器
 │   ├── sync.go                 # 修改：支持 github/feishu 多远端
 │   ├── bot.go / bot_serve.go   # 修改：支持 project 上下文
 │   ├── config.go
+│   ├── spec.go                 # `capture spec *` 命令族
+│   └── migrate.go              # `capture migrate` v1→v2 迁移
 ├── internal/
 │   ├── model/
 │   │   ├── task.go             # 修改：新增 WorkspaceID, ProjectID, ExternalID 等字段
 │   │   ├── config.go           # 修改：新增 Workspace/Project 配置结构
+│   │   ├── workspace.go        # 新增：Workspace 模型
+│   │   ├── project.go          # 新增：Project 模型
+│   │   ├── spec.go             # 新增：Spec 文档模型
+│   │   └── remote.go           # 新增：远端同步通用模型
 │   ├── store/
 │   │   ├── store.go            # 修改：扩展 Store interface
 │   │   ├── markdown.go         # 修改：支持 Project 子目录
 │   │   ├── sqlite.go           # 修改：新增 projects/specs 表
 │   │   ├── dual.go             # 修改：支持 Workspace 级初始化
+│   │   ├── workspace_store.go  # 新增：Workspace 元数据管理
+│   │   ├── project_store.go    # 新增：Project 元数据管理
+│   │   └── spec_store.go       # 新增：Spec 文档存储
 │   ├── service/
 │   │   ├── task_service.go     # 修改：新增 Project 上下文
+│   │   ├── workspace_service.go # 新增
+│   │   ├── project_service.go  # 新增
+│   │   ├── spec_service.go     # 新增
+│   │   └── sync_service.go     # 新增：统一同步编排
 │   ├── config/
 │   │   └── config.go           # 修改：支持分层配置合并
 │   ├── tui/
 │   │   ├── app.go              # 大幅修改：多视图 + Project 导航
 │   │   ├── keys.go
 │   │   ├── styles.go
+│   │   ├── workspace_selector.go # 新增
+│   │   ├── project_selector.go   # 新增
+│   │   ├── views/
+│   │   │   ├── kanban.go       # 从 app.go 提取
+│   │   │   ├── task_list.go    # 新增：列表视图
+│   │   │   ├── spec_browser.go # 新增：Spec 文档浏览器
+│   │   │   └── sync_status.go  # 新增：同步状态面板
+│   │   └── components/
+│   │       ├── selector.go     # 新增：通用选择器组件
+│   │       └── search.go       # 新增：搜索框组件
 │   ├── bot/
 │   │   ├── msgparser.go        # 修改：识别 project 上下文和 GitHub URL
 │   │   ├── dispatcher.go       # 修改：路由到 ProjectService
 │   │   ├── webhook.go
 │   │   └── websocket.go
 │   ├── bitable/
 │   │   ├── client.go
 │   │   ├── mapper.go           # 修改：支持 Project 字段映射
 │   │   └── sync.go             # 修改：多 Project 同步
 │   ├── feishu/
 │   │   └── client.go
+│   ├── github/                 # 新增：GitHub 集成模块
+│   │   ├── client.go           # GraphQL 客户端封装
+│   │   ├── project.go          # Project/Item 查询
+│   │   ├── mapper.go           # GitHub ↔ Task 字段映射
+│   │   ├── sync.go             # GitHub 同步逻辑
+│   │   └── auth.go             # Token 获取与验证
+│   ├── remote/                 # 新增：远端抽象层
+│   │   ├── remote.go           # Remote interface
+│   │   └── registry.go         # Remote 注册中心
 │   └── config/
 │       └── config.go
 └── pkg/
     ├── idgen/
     │   └── idgen.go
     ├── frontmatter/
     │   └── frontmatter.go
+    └── template/               # 新增：文档模板引擎
+        └── engine.go
```

---

## 三、新增模块详细设计

### 3.1 `internal/model/workspace.go`

```go
package model

import "time"

// Workspace 代表一个独立的工作空间，包含一组相关的 Project
type Workspace struct {
    ID          string            `yaml:"id" json:"id"`
    Name        string            `yaml:"name" json:"name"`
    Description string            `yaml:"description" json:"description"`
    Path        string            `yaml:"path" json:"path"`              // 本地目录路径
    DefaultProject string         `yaml:"default_project" json:"default_project"`
    Remotes     WorkspaceRemotes  `yaml:"remotes" json:"remotes"`
    CreatedAt   time.Time         `yaml:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `yaml:"updated_at" json:"updated_at"`
}

type WorkspaceRemotes struct {
    GitHub *GitHubRemoteConfig `yaml:"github,omitempty" json:"github,omitempty"`
    Feishu *FeishuRemoteConfig `yaml:"feishu,omitempty" json:"feishu,omitempty"`
}

type GitHubRemoteConfig struct {
    TokenSource string `yaml:"token_source" json:"token_source"` // "gh_cli", "env", "file"
    Token       string `yaml:"token,omitempty" json:"token,omitempty"`
    DefaultOrg  string `yaml:"default_org,omitempty" json:"default_org,omitempty"`
}

type FeishuRemoteConfig struct {
    AppID             string `yaml:"app_id" json:"app_id"`
    AppSecret         string `yaml:"app_secret" json:"app_secret"`
    VerificationToken string `yaml:"verification_token,omitempty" json:"verification_token,omitempty"`
    EncryptKey        string `yaml:"encrypt_key,omitempty" json:"encrypt_key,omitempty"`
}

func NewWorkspace(name string) *Workspace {
    now := time.Now()
    id := generateWorkspaceID(name) // slug 化，如 "personal"
    return &Workspace{
        ID:        id,
        Name:      name,
        Path:      id,
        CreatedAt: now,
        UpdatedAt: now,
    }
}
```

### 3.2 `internal/model/project.go`

```go
package model

import "time"

// Project 代表 Workspace 下的一个项目，可以关联 GitHub repo、飞书表格等
type Project struct {
    ID          string            `yaml:"id" json:"id"`
    WorkspaceID string            `yaml:"workspace_id" json:"workspace_id"`
    Name        string            `yaml:"name" json:"name"`
    Slug        string            `yaml:"slug" json:"slug"`              // URL-friendly ID
    Description string            `yaml:"description" json:"description"`
    GitHub      *ProjectGitHub    `yaml:"github,omitempty" json:"github,omitempty"`
    Feishu      *ProjectFeishu    `yaml:"feishu,omitempty" json:"feishu,omitempty"`
    CustomFields []CustomField    `yaml:"custom_fields,omitempty" json:"custom_fields,omitempty"`
    CreatedAt   time.Time         `yaml:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `yaml:"updated_at" json:"updated_at"`
}

type ProjectGitHub struct {
    Repo           string `yaml:"repo" json:"repo"`                      // "owner/repo"
    ProjectNumber  int    `yaml:"project_number,omitempty" json:"project_number,omitempty"`
    DefaultViewID  string `yaml:"default_view_id,omitempty" json:"default_view_id,omitempty"`
}

type ProjectFeishu struct {
    BitableAppToken string `yaml:"bitable_app_token,omitempty" json:"bitable_app_token,omitempty"`
    BitableTableID  string `yaml:"bitable_table_id,omitempty" json:"bitable_table_id,omitempty"`
}

type CustomField struct {
    Name  string `yaml:"name" json:"name"`
    Type  string `yaml:"type" json:"type"`   // "text", "select", "number", "date"
    Value string `yaml:"value,omitempty" json:"value,omitempty"`
}

func NewProject(workspaceID, name string) *Project {
    now := time.Now()
    slug := slugify(name)
    return &Project{
        ID:          generateProjectID(workspaceID, slug),
        WorkspaceID: workspaceID,
        Name:        name,
        Slug:        slug,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}
```

### 3.3 `internal/model/spec.go`

```go
package model

import "time"

// SpecStatus 代表规格文档的生命周期
type SpecStatus string

const (
    SpecStatusDraft      SpecStatus = "draft"
    SpecStatusReviewing  SpecStatus = "reviewing"
    SpecStatusApproved   SpecStatus = "approved"
    SpecStatusImplemented SpecStatus = "implemented"
    SpecStatusArchived   SpecStatus = "archived"
)

// Spec 代表一个规格/PRD 文档
type Spec struct {
    ID          string        `yaml:"id" json:"id"`
    ProjectID   string        `yaml:"project_id" json:"project_id"`
    Title       string        `yaml:"title" json:"title"`
    Status      SpecStatus    `yaml:"status" json:"status"`
    Author      string        `yaml:"author" json:"author"`
    Tags        []string      `yaml:"tags" json:"tags"`
    LinkedTasks []string      `yaml:"linked_tasks" json:"linked_tasks"`
    CreatedAt   time.Time     `yaml:"created_at" json:"created_at"`
    UpdatedAt   time.Time     `yaml:"updated_at" json:"updated_at"`
    FilePath    string        `yaml:"-" json:"-"`
    Body        string        `yaml:"-" json:"-"`    // Markdown body
}

func NewSpec(projectID, title string) *Spec {
    now := time.Now()
    return &Spec{
        ID:        generateSpecID(projectID),
        ProjectID: projectID,
        Title:     title,
        Status:    SpecStatusDraft,
        Tags:      []string{},
        LinkedTasks: []string{},
        CreatedAt: now,
        UpdatedAt: now,
    }
}
```

### 3.4 `internal/model/task.go`（增量变更）

```go
// v2 在 v1 Task 基础上新增以下字段：

type Task struct {
    // ... v1 现有字段 ...
    
    // v2 新增
    WorkspaceID string `yaml:"workspace_id" json:"workspace_id"`
    ProjectID   string `yaml:"project_id" json:"project_id"`
    
    // 远端关联
    ExternalID  string `yaml:"external_id" json:"external_id"`    // GitHub Issue node_id / Feishu record_id
    ExternalURL string `yaml:"external_url" json:"external_url"`  // 远端链接
    
    // 同步元数据
    SyncVersion int    `yaml:"sync_version" json:"sync_version"`  // 用于乐观锁/冲突检测
    
    // ...
}

// v1 兼容：如果 WorkspaceID/ProjectID 为空，表示这是 v1 遗留任务
// 迁移时会自动填充为 "default" workspace 和 "legacy" project
```

### 3.5 `internal/store/workspace_store.go`

```go
package store

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
    "github.com/variableway/innate/capture/internal/model"
)

// WorkspaceStore 管理工作空间的元数据和目录
type WorkspaceStore struct {
    rootPath string  // ~/.capture/
}

func NewWorkspaceStore(rootPath string) *WorkspaceStore {
    return &WorkspaceStore{rootPath: rootPath}
}

func (s *WorkspaceStore) workspacesDir() string {
    return filepath.Join(s.rootPath, "workspaces")
}

func (s *WorkspaceStore) workspacePath(id string) string {
    return filepath.Join(s.workspacesDir(), id)
}

func (s *WorkspaceStore) CreateWorkspace(ctx context.Context, ws *model.Workspace) error {
    path := s.workspacePath(ws.ID)
    if err := os.MkdirAll(path, 0755); err != nil {
        return fmt.Errorf("failed to create workspace dir: %w", err)
    }
    
    // 创建子目录结构
    dirs := []string{
        filepath.Join(path, "projects"),
        filepath.Join(path, "templates"),
    }
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            return err
        }
    }
    
    // 保存 workspace.yaml
    ws.Path = path
    data, err := yaml.Marshal(ws)
    if err != nil {
        return err
    }
    configPath := filepath.Join(path, "workspace.yaml")
    return os.WriteFile(configPath, data, 0644)
}

func (s *WorkspaceStore) GetWorkspace(ctx context.Context, id string) (*model.Workspace, error) {
    path := s.workspacePath(id)
    configPath := filepath.Join(path, "workspace.yaml")
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("workspace %s not found: %w", id, err)
    }
    
    var ws model.Workspace
    if err := yaml.Unmarshal(data, &ws); err != nil {
        return nil, err
    }
    ws.Path = path
    return &ws, nil
}

func (s *WorkspaceStore) ListWorkspaces(ctx context.Context) ([]*model.Workspace, error) {
    entries, err := os.ReadDir(s.workspacesDir())
    if err != nil {
        return nil, err
    }
    
    var workspaces []*model.Workspace
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }
        ws, err := s.GetWorkspace(ctx, entry.Name())
        if err == nil {
            workspaces = append(workspaces, ws)
        }
    }
    return workspaces, nil
}

func (s *WorkspaceStore) DeleteWorkspace(ctx context.Context, id string) error {
    path := s.workspacePath(id)
    return os.RemoveAll(path)
}
```

### 3.6 `internal/store/project_store.go`

```go
package store

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
    "github.com/variableway/innate/capture/internal/model"
)

// ProjectStore 管理单个 Workspace 下的 Project
type ProjectStore struct {
    workspacePath string
    sqlite        *SQLiteStore  // 复用现有 SQLite 存储项目元数据
}

func NewProjectStore(workspacePath string, sqlite *SQLiteStore) *ProjectStore {
    return &ProjectStore{
        workspacePath: workspacePath,
        sqlite:        sqlite,
    }
}

func (s *ProjectStore) projectPath(slug string) string {
    return filepath.Join(s.workspacePath, "projects", slug)
}

func (s *ProjectStore) CreateProject(ctx context.Context, project *model.Project) error {
    path := s.projectPath(project.Slug)
    if err := os.MkdirAll(path, 0755); err != nil {
        return err
    }
    
    // 创建子目录
    for _, d := range []string{"tasks", "specs"} {
        if err := os.MkdirAll(filepath.Join(path, d), 0755); err != nil {
            return err
        }
    }
    
    // 保存 project.yaml
    data, err := yaml.Marshal(project)
    if err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(path, "project.yaml"), data, 0644); err != nil {
        return err
    }
    
    // 写入 SQLite
    _, err = s.sqlite.db.ExecContext(ctx,
        `INSERT INTO projects (id, workspace_id, name, slug, description, github_repo, github_project_number, feishu_bitable_app_token, feishu_bitable_table_id, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        project.ID, project.WorkspaceID, project.Name, project.Slug, project.Description,
        project.GitHub.Repo, project.GitHub.ProjectNumber,
        project.Feishu.BitableAppToken, project.Feishu.BitableTableID,
        project.CreatedAt.Format(timeFormat), project.UpdatedAt.Format(timeFormat),
    )
    return err
}

// GetProject, ListProjects, UpdateProject, DeleteProject ...
```

### 3.7 `internal/service/workspace_service.go`

```go
package service

import (
    "context"
    "fmt"
    
    "github.com/variableway/innate/capture/internal/model"
    "github.com/variableway/innate/capture/internal/store"
)

type WorkspaceService struct {
    wsStore    *store.WorkspaceStore
    dataDir    string
}

func NewWorkspaceService(wsStore *store.WorkspaceStore, dataDir string) *WorkspaceService {
    return &WorkspaceService{wsStore: wsStore, dataDir: dataDir}
}

func (s *WorkspaceService) Create(ctx context.Context, name string) (*model.Workspace, error) {
    ws := model.NewWorkspace(name)
    if err := s.wsStore.CreateWorkspace(ctx, ws); err != nil {
        return nil, err
    }
    
    // 自动创建 SQLite 数据库
    dbPath := filepath.Join(ws.Path, "capture.db")
    _, err := store.NewSQLiteStore(filepath.Dir(dbPath))
    if err != nil {
        return nil, fmt.Errorf("failed to init workspace database: %w", err)
    }
    
    return ws, nil
}

func (s *WorkspaceService) Get(ctx context.Context, id string) (*model.Workspace, error) {
    return s.wsStore.GetWorkspace(ctx, id)
}

func (s *WorkspaceService) List(ctx context.Context) ([]*model.Workspace, error) {
    return s.wsStore.ListWorkspaces(ctx)
}

func (s *WorkspaceService) Delete(ctx context.Context, id string) error {
    return s.wsStore.DeleteWorkspace(ctx, id)
}

func (s *WorkspaceService) GetDefault(ctx context.Context) (*model.Workspace, error) {
    workspaces, err := s.List(ctx)
    if err != nil {
        return nil, err
    }
    if len(workspaces) == 0 {
        return nil, fmt.Errorf("no workspace found, run `capture workspace create` first")
    }
    return workspaces[0], nil
}
```

### 3.8 `internal/service/project_service.go`

```go
package service

import (
    "context"
    "fmt"
    "strings"
    
    "github.com/variableway/innate/capture/internal/model"
    "github.com/variableway/innate/capture/internal/store"
)

type ProjectService struct {
    projStore  *store.ProjectStore
}

func NewProjectService(projStore *store.ProjectStore) *ProjectService {
    return &ProjectService{projStore: projStore}
}

func (s *ProjectService) Create(ctx context.Context, workspaceID, name string, opts ...ProjectOption) (*model.Project, error) {
    if strings.TrimSpace(name) == "" {
        return nil, fmt.Errorf("project name cannot be empty")
    }
    
    project := model.NewProject(workspaceID, name)
    for _, opt := range opts {
        opt(project)
    }
    
    if err := s.projStore.CreateProject(ctx, project); err != nil {
        return nil, fmt.Errorf("failed to create project: %w", err)
    }
    
    return project, nil
}

func (s *ProjectService) Get(ctx context.Context, id string) (*model.Project, error) {
    return s.projStore.GetProject(ctx, id)
}

func (s *ProjectService) ListByWorkspace(ctx context.Context, workspaceID string) ([]*model.Project, error) {
    return s.projStore.ListProjects(ctx, workspaceID)
}

// ProjectOption 功能选项模式
type ProjectOption func(*model.Project)

func WithGitHubRepo(repo string) ProjectOption {
    return func(p *model.Project) {
        p.GitHub = &model.ProjectGitHub{Repo: repo}
    }
}

func WithDescription(desc string) ProjectOption {
    return func(p *model.Project) {
        p.Description = desc
    }
}
```

### 3.9 `internal/remote/remote.go`

```go
package remote

import (
    "context"
    "time"
)

// Remote 定义所有远端同步源的统一接口
type Remote interface {
    Name() string
    IsConfigured() bool
    
    // 项目级
    ListProjects(ctx context.Context) ([]RemoteProject, error)
    
    // 任务级（Item）
    ListItems(ctx context.Context, projectID string) ([]RemoteItem, error)
    GetItem(ctx context.Context, projectID, itemID string) (*RemoteItem, error)
    CreateItem(ctx context.Context, projectID string, item *RemoteItem) (*RemoteItem, error)
    UpdateItem(ctx context.Context, projectID string, item *RemoteItem) (*RemoteItem, error)
    DeleteItem(ctx context.Context, projectID, itemID string) error
    
    // 同步
    LastSyncTime(ctx context.Context, projectID string) (time.Time, error)
}

type RemoteProject struct {
    ID          string
    Name        string
    Description string
    URL         string
    ExternalID  string
}

type RemoteItem struct {
    ID         string
    ExternalID string
    Title      string
    Body       string
    Status     string
    Priority   string
    URL        string
    UpdatedAt  time.Time
    Raw        map[string]interface{}
}

// Registry 管理所有可用的 Remote 实现
type Registry struct {
    remotes map[string]Remote
}

func NewRegistry() *Registry {
    return &Registry{remotes: make(map[string]Remote)}
}

func (r *Registry) Register(remote Remote) {
    r.remotes[remote.Name()] = remote
}

func (r *Registry) Get(name string) (Remote, bool) {
    remote, ok := r.remotes[name]
    return remote, ok
}

func (r *Registry) List() []Remote {
    var list []Remote
    for _, remote := range r.remotes {
        list = append(list, remote)
    }
    return list
}
```

### 3.10 `internal/github/client.go`

```go
package github

import (
    "context"
    "fmt"
    
    "github.com/shurcooL/githubv4"
    "golang.org/x/oauth2"
)

// Client 封装 GitHub GraphQL API 调用
type Client struct {
    v4    *githubv4.Client
    token string
}

func NewClient(ctx context.Context, token string) *Client {
    src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    httpClient := oauth2.NewClient(ctx, src)
    return &Client{
        v4:    githubv4.NewClient(httpClient),
        token: token,
    }
}

// NewClientFromGHCLI 尝试从 gh CLI 获取 token
func NewClientFromGHCLI(ctx context.Context) (*Client, error) {
    token, err := getTokenFromGHCLI()
    if err != nil {
        return nil, err
    }
    return NewClient(ctx, token), nil
}

// QueryProjectItems 查询 GitHub Project 的所有 items
func (c *Client) QueryProjectItems(ctx context.Context, org, projectNumber int) (*ProjectItemsQuery, error) {
    var query ProjectItemsQuery
    variables := map[string]interface{}{
        "login": githubv4.String(org),
        "number": githubv4.Int(projectNumber),
        "first": githubv4.Int(100),
    }
    err := c.v4.Query(ctx, &query, variables)
    return &query, err
}

// ProjectItemsQuery GraphQL 查询结构
// 注意：实际需要根据 GitHub GraphQL schema 定义完整的结构
```

### 3.11 `cmd/workspace.go`

```go
package cmd

import (
    "context"
    "fmt"
    
    "github.com/spf13/cobra"
    "github.com/variableway/innate/capture/internal/service"
    "github.com/variableway/innate/capture/internal/store"
)

var workspaceCmd = &cobra.Command{
    Use:   "workspace",
    Short: "Manage workspaces",
    Long:  "Create, list, switch and manage capture workspaces.",
}

var workspaceCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Create a new workspace",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := getDataDir()
        wsStore := store.NewWorkspaceStore(dir)
        svc := service.NewWorkspaceService(wsStore, dir)
        
        ws, err := svc.Create(context.Background(), args[0])
        if err != nil {
            return err
        }
        
        fmt.Printf("Created workspace: %s (%s)\n", ws.Name, ws.ID)
        fmt.Printf("Path: %s\n", ws.Path)
        return nil
    },
}

var workspaceListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all workspaces",
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := getDataDir()
        wsStore := store.NewWorkspaceStore(dir)
        svc := service.NewWorkspaceService(wsStore, dir)
        
        workspaces, err := svc.List(context.Background())
        if err != nil {
            return err
        }
        
        if len(workspaces) == 0 {
            fmt.Println("No workspaces found. Run `capture workspace create <name>` to create one.")
            return nil
        }
        
        fmt.Println("WORKSPACE          PATH")
        for _, ws := range workspaces {
            fmt.Printf("%-18s %s\n", ws.Name, ws.Path)
        }
        return nil
    },
}

var workspaceSwitchCmd = &cobra.Command{
    Use:   "switch <workspace-id>",
    Short: "Switch to a workspace",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // 将当前 workspace 写入全局配置
        viper.Set("current_workspace", args[0])
        return viper.WriteConfig()
    },
}

func init() {
    workspaceCmd.AddCommand(workspaceCreateCmd)
    workspaceCmd.AddCommand(workspaceListCmd)
    workspaceCmd.AddCommand(workspaceSwitchCmd)
    rootCmd.AddCommand(workspaceCmd)
}
```

### 3.12 `cmd/project.go`

```go
package cmd

import (
    "context"
    "fmt"
    
    "github.com/spf13/cobra"
    "github.com/variableway/innate/capture/internal/service"
    "github.com/variableway/innate/capture/internal/store"
)

var projectCmd = &cobra.Command{
    Use:   "project",
    Short: "Manage projects within a workspace",
}

var projectCreateCmd = &cobra.Command{
    Use:   "create <name>",
    Short: "Create a new project",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        dir := getDataDir()
        workspaceID := viper.GetString("current_workspace")
        if workspaceID == "" {
            workspaceID = "default"
        }
        
        // 获取 workspace 路径，创建 project
        wsStore := store.NewWorkspaceStore(dir)
        ws, err := wsStore.GetWorkspace(context.Background(), workspaceID)
        if err != nil {
            return fmt.Errorf("workspace not found: %w", err)
        }
        
        sqlite, err := store.NewSQLiteStore(ws.Path)
        if err != nil {
            return err
        }
        defer sqlite.Close()
        
        projStore := store.NewProjectStore(ws.Path, sqlite)
        svc := service.NewProjectService(projStore)
        
        repo, _ := cmd.Flags().GetString("github-repo")
        desc, _ := cmd.Flags().GetString("description")
        
        opts := []service.ProjectOption{}
        if desc != "" {
            opts = append(opts, service.WithDescription(desc))
        }
        if repo != "" {
            opts = append(opts, service.WithGitHubRepo(repo))
        }
        
        project, err := svc.Create(context.Background(), workspaceID, args[0], opts...)
        if err != nil {
            return err
        }
        
        fmt.Printf("Created project: %s (%s)\n", project.Name, project.Slug)
        return nil
    },
}

var projectListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all projects in current workspace",
    RunE: func(cmd *cobra.Command, args []string) error {
        // ...
        fmt.Println("PROJECT            DESCRIPTION")
        // ...
        return nil
    },
}

func init() {
    projectCreateCmd.Flags().String("github-repo", "", "Associated GitHub repo (owner/repo)")
    projectCreateCmd.Flags().String("description", "", "Project description")
    
    projectCmd.AddCommand(projectCreateCmd)
    projectCmd.AddCommand(projectListCmd)
    rootCmd.AddCommand(projectCmd)
}
```

---

## 四、CLI 命令完整映射

### 4.1 v2 CLI 命令树

```
capture
├── workspace                    # 工作空间管理
│   ├── create <name>            # 创建工作空间
│   ├── list                     # 列出工作空间
│   ├── switch <id>              # 切换工作空间
│   ├── show [<id>]              # 显示工作空间详情
│   └── delete <id>              # 删除工作空间
│
├── project                      # 项目管理（在当前工作空间）
│   ├── create <name>            # 创建项目 [--github-repo] [--description]
│   ├── list                     # 列出项目
│   ├── show <slug>              # 显示项目详情
│   ├── set <slug>               # 设置当前默认项目
│   ├── delete <slug>            # 删除项目
│   └── sync <slug>              # 手动同步项目（GitHub/飞书）
│
├── add "title"                  # 创建任务 [--project] [--stage] [--priority] [--tag]
├── list                         # 列出任务 [--workspace] [--project] [--status] [--stage]
├── show <TASK-ID>
├── edit <TASK-ID>
├── delete <TASK-ID>
├── status <TASK-ID> <status>
├── stage <TASK-ID> <stage>
├── assign <TASK-ID> [opts]
│
├── spec                         # 规格文档管理
│   ├── create <title>           # 创建 Spec [--project] [--template]
│   ├── list                     # 列出 Spec [--project]
│   ├── show <SPEC-ID>
│   ├── edit <SPEC-ID>
│   ├── link <SPEC-ID> <TASK-ID> # 关联任务
│   └── generate-tasks <SPEC-ID> # 从 Acceptance Criteria 生成任务
│
├── kanban                       # TUI 看板 [--project]
├── sync                         # 同步
│   ├── github                   # 同步 GitHub Project
│   ├── feishu                   # 同步飞书 Bitable
│   └── status                   # 查看同步状态
│
├── bot                          # 飞书 Bot
│   └── serve [--mode]
│
├── migrate                      # v1 → v2 数据迁移
│   ├── run                      # 执行迁移
│   ├── status                   # 检查迁移状态
│   └── rollback                 # 回滚迁移
│
├── init                         # 初始化（v2 目录结构）
└── config                       # 配置管理
    ├── get <key>
    └── set <key> <value>
```

### 4.2 常用快捷命令设计

```bash
# 最常用的场景保持极简

# 场景 1：快速记录一个想法（保持 v1 习惯）
capture add "优化构建脚本"
# → 自动使用当前 workspace + 当前 project

# 场景 2：给指定项目添加任务
capture add "修复登录 bug" --project innate-capture

# 场景 3：查看当前项目的所有待办
capture list --status todo

# 场景 4：启动看板（自动进入当前 project）
capture kanban

# 场景 5：拉取 GitHub Project 最新状态
capture project sync innate-capture

# 场景 6：创建 Spec 并生成任务
capture spec create "用户认证模块" --project innate-capture
capture spec generate-tasks SPEC-001
```

---

## 五、TUI 界面流程

### 5.1 启动流程

```
启动 capture kanban
    │
    ├── 检查是否有多个 Workspace
    │   ├── 是 → 显示 Workspace 选择器
    │   └── 否 → 直接进入
    │
    ├── 检查是否有多个 Project
    │   ├── 是 → 显示 Project 选择器（可跳过，默认显示全部）
    │   └── 否 → 直接进入
    │
    └── 加载看板视图（Kanban by Status）
```

### 5.2 视图切换快捷键

```
1 / k          → Kanban 视图（按 Status）
2 / s          → Stage Pipeline 视图（按 Stage 阶段）
3 / l          → List 视图（可搜索/过滤）
4 / p          → Project 概览视图
5 / d          → Sync 状态面板
Enter          → 查看详情
n              → 新建 Task
N              → 新建 Spec
/              → 搜索
q / Ctrl+C     → 退出
```

---

## 六、关键接口契约

### 6.1 Store Interface（v2 扩展）

```go
package store

import (
    "context"
    "github.com/variableway/innate/capture/internal/model"
)

// Store 保持 v1 契约，新增 Project 过滤能力
type Store interface {
    CreateTask(ctx context.Context, task *model.Task) error
    GetTask(ctx context.Context, id string) (*model.Task, error)
    UpdateTask(ctx context.Context, task *model.Task) error
    DeleteTask(ctx context.Context, id string) error
    ListTasks(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error)
}

// TaskFilter 扩展
type TaskFilter struct {
    WorkspaceID *string
    ProjectID   *string
    Status      *model.TaskStatus
    Stage       *model.TaskStage
    Priority    *model.TaskPriority
    Tags        []string
    Source      *string
    ExternalID  *string
}

// 新增接口

type WorkspaceStore interface {
    CreateWorkspace(ctx context.Context, ws *model.Workspace) error
    GetWorkspace(ctx context.Context, id string) (*model.Workspace, error)
    ListWorkspaces(ctx context.Context) ([]*model.Workspace, error)
    UpdateWorkspace(ctx context.Context, ws *model.Workspace) error
    DeleteWorkspace(ctx context.Context, id string) error
}

type ProjectStore interface {
    CreateProject(ctx context.Context, project *model.Project) error
    GetProject(ctx context.Context, id string) (*model.Project, error)
    ListProjects(ctx context.Context, workspaceID string) ([]*model.Project, error)
    UpdateProject(ctx context.Context, project *model.Project) error
    DeleteProject(ctx context.Context, id string) error
}

type SpecStore interface {
    CreateSpec(ctx context.Context, spec *model.Spec) error
    GetSpec(ctx context.Context, id string) (*model.Spec, error)
    ListSpecs(ctx context.Context, projectID string) ([]*model.Spec, error)
    UpdateSpec(ctx context.Context, spec *model.Spec) error
    DeleteSpec(ctx context.Context, id string) error
}
```

---

## 七、模块依赖关系图

```
                    ┌─────────────┐
                    │   cmd/*.go  │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
      ┌─────────────┐ ┌────────┐ ┌─────────────┐
      │  service/*  │ │ tui/*  │ │   bot/*     │
      └──────┬──────┘ └───┬────┘ └──────┬──────┘
             │            │             │
             └────────────┼─────────────┘
                          │
                    ┌─────┴─────┐
                    │ store/*   │
                    │ model/*   │
                    └─────┬─────┘
                          │
             ┌────────────┼────────────┐
             │            │            │
             ▼            ▼            ▼
      ┌─────────────┐ ┌────────┐ ┌─────────────┐
      │  github/*   │ │feishu/*│ │  bitable/*  │
      │  (remote)   │ │(remote)│ │  (remote)   │
      └─────────────┘ └────────┘ └─────────────┘
```

**依赖原则**：
- `cmd` 依赖 `service` 和 `tui`
- `service` 依赖 `store` 和 `model`
- `store` 依赖 `model`
- `remote`（github/feishu/bitable）依赖 `model`，被 `service` 使用
- **禁止循环依赖**：`remote` 不依赖 `service`
