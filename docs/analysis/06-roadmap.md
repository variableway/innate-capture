# Capture v2 实施路线图

> 制定时间：2026-06-04
> 预估总工期：6-8 周（单开发者，兼职投入）
> 方法论：MVP 优先、快速迭代、风险前置

---

## 一、发布策略

采用 **Feature Flag + 渐进发布** 策略：

| 版本 | 代号 | 目标 | 用户影响 |
|------|------|------|----------|
| v1.1 | Foundation | 后台重构，无用户可见变化 | 零影响 |
| v1.2 | Workspace | 新增 Workspace/Project 管理，v1 功能保持 | 可选使用新功能 |
| v2.0 | GitHub | GitHub Project 集成 + 完整 Workspace 功能 | 推荐升级 |
| v2.1 | Spec | Spec 文档管理 + TUI 增强 | 功能增强 |
| v2.2 | Sync | 飞书多 Project 同步 + 双向 GitHub 同步 | 功能增强 |

---

## 二、里程碑详细计划

### Milestone 1: Foundation（基础重构）
**目标**：为 v2 做准备，不引入用户可见变化
**工期**：1 周
**风险等级**：低

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M1.1 提取 `pkg/frontmatter` 为独立可测试模块 | P0 | 4h | 新增 frontmatter_test.go，覆盖率 > 80% |
| M1.2 重构 `Store` interface，预留 Workspace/Project 扩展点 | P0 | 6h | TaskFilter 新增字段不破坏现有测试 |
| M1.3 创建 `internal/remote/` 抽象层框架 | P0 | 4h | Remote interface 定义完成，有 mock 实现 |
| M1.4 引入 `Task.WorkspaceID/ProjectID` 字段（nullable） | P0 | 4h | v1 任务读取正常，新增字段为空 |
| M1.5 SQLite schema migration：新增 projects 表（空表） | P0 | 2h | `capture list` 正常工作 |
| M1.6 集成测试：确保 v1 所有命令在重构后正常 | P0 | 4h | 全部现有测试通过 |

**交付物**：
- 代码重构完成，v1 功能 100% 保持
- `internal/remote/remote.go` 接口定义
- 扩展的 `TaskFilter` 结构

---

### Milestone 2: Workspace（工作空间）
**目标**：Workspace + Project 管理 + v1 数据迁移
**工期**：2 周
**风险等级**：中（数据迁移是关键）

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M2.1 实现 `internal/model/workspace.go` + `project.go` | P0 | 4h | 模型定义完整，有 JSON/YAML 序列化测试 |
| M2.2 实现 `internal/store/workspace_store.go` | P0 | 6h | 支持 CRUD，目录结构正确创建 |
| M2.3 实现 `internal/store/project_store.go` | P0 | 6h | 支持 CRUD，与 SQLite 同步 |
| M2.4 实现 `internal/service/workspace_service.go` | P0 | 4h | 支持 Create/List/Get/Delete |
| M2.5 实现 `internal/service/project_service.go` | P0 | 4h | 支持 Create/List/Get/Delete |
| M2.6 实现 `cmd/workspace.go` 命令族 | P0 | 6h | `capture workspace create/list/switch/show/delete` 可用 |
| M2.7 实现 `cmd/project.go` 命令族 | P0 | 6h | `capture project create/list/show/set/delete` 可用 |
| M2.8 修改 `cmd/init.go`：初始化 v2 目录结构 | P0 | 4h | `capture init` 创建 workspaces/ 目录 |
| M2.9 修改 `cmd/add.go`：支持 `--project` 参数 | P0 | 4h | 未指定 project 时使用 default |
| M2.10 修改 `cmd/list.go`：支持 `--workspace/--project` 过滤 | P0 | 4h | 过滤正确，无 project 任务也显示 |
| M2.11 实现 `cmd/migrate.go`：v1 → v2 自动迁移 | P0 | 8h | 迁移成功，数据完整，可回滚 |
| M2.12 TUI 最小适配：顶部显示当前 Project | P1 | 4h | TUI 正常启动，显示 project 名 |
| M2.13 集成测试：Workspace + Project 完整流程 | P0 | 6h | E2E 测试通过 |

**交付物**：
- `capture workspace *` 和 `capture project *` 命令可用
- v1 数据可自动迁移到 v2
- 所有 v1 命令在 v2 下正常工作（向后兼容）

**关键检查点**：
- [ ] 迁移 50+ 个 v1 任务，验证数据完整性
- [ ] 删除 workspace 后再创建同名 workspace，无残留数据
- [ ] TUI 在 default workspace + default project 下与 v1 行为一致

---

### Milestone 3: GitHub Integration（GitHub 集成）
**目标**：GitHub Project 只读同步
**工期**：1.5 周
**风险等级**：中（GraphQL API 是新领域）

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M3.1 调研 GitHub Project v2 GraphQL schema | P0 | 4h | 能手动执行 query 获取 project items |
| M3.2 实现 `internal/github/auth.go`：多来源 token 获取 | P0 | 4h | `gh auth`, `GITHUB_TOKEN`, 手动配置 三种方式 |
| M3.3 实现 `internal/github/client.go`：GraphQL 客户端 | P0 | 6h | 能成功查询 Project items |
| M3.4 实现 `internal/github/project.go`：Project/Item 查询 | P0 | 6h | 支持分页获取全部 items |
| M3.5 实现 `internal/github/mapper.go`：GitHub ↔ Task 映射 | P0 | 4h | Status/Priority/Title/Body 正确映射 |
| M3.6 实现 `internal/github/sync.go`：只读同步逻辑 | P0 | 6h | 执行同步后本地生成对应的 Task Markdown |
| M3.7 修改 `cmd/sync.go`：支持 `capture sync github --project` | P0 | 4h | 命令可用，有进度输出 |
| M3.8 实现缓存层：避免重复 API 调用 | P1 | 4h | 第二次同步使用缓存，速度提升 10x |
| M3.9 TUI Sync 状态面板（最小实现） | P1 | 4h | 显示最后同步时间和项目数 |
| M3.10 集成测试：GitHub 同步 E2E | P0 | 6h | 使用测试仓库验证同步正确性 |

**交付物**：
- `capture sync github` 可用
- GitHub Project items 可拉取为本地 Task
- Token 自动获取（gh CLI 优先）

**关键检查点**：
- [ ] 使用真实 GitHub Project 测试，100+ items 同步成功
- [ ] API 速率限制处理：超过限制时友好提示
- [ ] 网络中断恢复：中断后重新同步无重复数据

---

### Milestone 4: TUI Enhancement（TUI 增强）
**目标**：多视图、Project 导航、性能优化
**工期**：1 周
**风险等级**：低

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M4.1 实现 `tui/workspace_selector.go` | P1 | 4h | 多 workspace 时可选择 |
| M4.2 实现 `tui/project_selector.go` | P1 | 4h | 多 project 时可选择/过滤 |
| M4.3 重构 `tui/app.go`：支持多视图状态机 | P0 | 8h | Kanban/Stage/List/Spec 视图可切换 |
| M4.4 实现 `tui/views/stage_pipeline.go` | P1 | 6h | 按 Stage 分 9 列，可横向滚动 |
| M4.5 实现 `tui/views/task_list.go` | P1 | 6h | 列表视图，支持实时搜索 |
| M4.6 实现 `tui/views/spec_browser.go` | P1 | 6h | 显示 Project 下的 Spec 列表，可预览 |
| M4.7 TUI 性能优化：SQLite 分页 + 虚拟滚动 | P0 | 6h | 500+ 任务不卡顿 |
| M4.8 统一快捷键帮助面板 | P1 | 2h | `?` 键显示完整快捷键说明 |

**交付物**：
- TUI 支持 Workspace/Project 选择
- Kanban + Stage Pipeline + List + Spec 四个视图
- 500+ 任务流畅运行

**关键检查点**：
- [ ] 窗口 resize 时布局正确
- [ ] 快速切换视图无闪烁
- [ ] 搜索 500 条任务，响应时间 < 200ms

---

### Milestone 5: Spec Management（文档管理）
**目标**：Spec/PRD 文档的模板化创建和生命周期管理
**工期**：1 周
**风险等级**：低

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M5.1 设计 Spec 文档模板格式（Go text/template） | P0 | 4h | 有 default/feature/bugfix 三种模板 |
| M5.2 实现 `pkg/template/engine.go`：模板引擎 | P0 | 4h | 支持变量替换和条件渲染 |
| M5.3 实现 `internal/model/spec.go` | P0 | 2h | Spec 模型定义完成 |
| M5.4 实现 `internal/store/spec_store.go` | P0 | 4h | Markdown 存储 + SQLite 索引 |
| M5.5 实现 `internal/service/spec_service.go` | P0 | 4h | CRUD + 模板渲染 + 任务生成 |
| M5.6 实现 `cmd/spec.go` 命令族 | P0 | 6h | `capture spec create/list/show/edit/link/generate-tasks` |
| M5.7 实现 Spec → Task 自动生成功能 | P1 | 4h | Acceptance Criteria 每项生成一个 Task |
| M5.8 用户自定义模板支持 | P1 | 4h | `~/.capture/templates/` 下自定义模板优先 |
| M5.9 TUI Spec 视图集成 | P1 | 4h | TUI 中可浏览 Spec，按 Enter 查看详情 |

**交付物**：
- `capture spec create` 可用
- Spec 模板系统支持自定义
- Acceptance Criteria 可生成 Task

---

### Milestone 6: Feishu Upgrade（飞书升级）
**目标**：飞书 Bot 和 Bitable 支持多 Project
**工期**：0.5 周
**风险等级**：低

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M6.1 修改 `bot/msgparser.go`：识别 project 上下文 | P1 | 4h | Bot 消息可指定项目 |
| M6.2 修改 `bot/dispatcher.go`：路由到 ProjectService | P1 | 4h | 创建任务时正确归属 Project |
| M6.3 新增 Bot 命令：`项目 列出`、`项目 切换` | P1 | 4h | 飞书中可用 |
| M6.4 修改 `bitable/sync.go`：多 Project 同步 | P1 | 6h | 每个 Project 对应 Bitable 子表或视图 |
| M6.5 飞书消息中 GitHub URL 自动识别 | P2 | 4h | 发送 GitHub Issue URL 自动创建关联 Task |

**交付物**：
- 飞书 Bot 支持 Project 维度操作
- Bitable 同步支持多 Project

---

### Milestone 7: Bidirectional Sync（双向同步）
**目标**：GitHub Project 双向同步 + 冲突解决
**工期**：1 周
**风险等级**：高（冲突处理复杂）

| 任务 | 优先级 | 预估工时 | 验收标准 |
|------|--------|----------|----------|
| M7.1 实现 GitHub GraphQL mutations：创建/更新 item | P0 | 6h | 可通过 API 修改 GitHub Project item |
| M7.2 实现 `internal/github/sync.go` 双向同步逻辑 | P0 | 8h | 本地修改可写回 GitHub |
| M7.3 实现冲突检测机制 | P0 | 6h | 能检测本地和远端同时修改 |
| M7.4 实现冲突解决策略（远端优先/本地优先/最新优先） | P0 | 6h | `capture sync github --strategy local` 可用 |
| M7.5 交互式冲突解决（TUI） | P1 | 6h | TUI 中显示冲突项，用户选择保留哪个 |
| M7.6 同步日志和审计 | P1 | 4h | `capture sync status` 显示历史同步记录 |

**交付物**：
- `capture sync github --bidirectional` 可用
- 冲突检测和解决机制
- 同步日志

---

## 三、时间线总览

```
周次:  W1    W2    W3    W4    W5    W6    W7    W8
      ├─────┼─────┼─────┼─────┼─────┼─────┼─────┼─────┤
M1    ████  │     │     │     │     │     │     │     │ Foundation
M2          ████████████  │     │     │     │     │     │ Workspace
M3                      ███████ │     │     │     │     │ GitHub (read)
M4                            ██████│     │     │     │ TUI Enhancement
M5                                  ██████│     │     │ Spec
M6                                        ████│     │     │ Feishu Upgrade
M7                                            ████████│     │ Bidirectional
      ├─────┼─────┼─────┼─────┼─────┼─────┼─────┼─────┤
v1.1  ▲     │     │     │     │     │     │     │     │
v1.2        ▲─────┘     │     │     │     │     │     │
v2.0              ▲─────┘     │     │     │     │     │
v2.1                    ▲─────┘     │     │     │     │
v2.2                          ▲─────┘     │     │     │
v2.3                                ▲─────┘     │     │
```

**发布节奏**：
- **v1.1**（W1 末）：Foundation 完成后发布，零用户可见变化
- **v1.2**（W3 末）：Workspace 完成后发布，新增 workspace/project 命令
- **v2.0**（W4 末）：GitHub 只读同步完成，大版本发布
- **v2.1**（W6 初）：TUI + Spec 完成后发布
- **v2.2**（W7 初）：飞书升级完成后发布
- **v2.3**（W8 末）：双向同步完成后发布

---

## 四、风险缓解与应急预案

### 4.1 风险应对矩阵

| 风险 | 发生概率 | 影响 | 应对措施 | 触发条件 |
|------|----------|------|----------|----------|
| GitHub GraphQL API 不稳定 | 中 | 高 | 降级为只读同步，禁用写入 | API 返回 5xx > 10% |
| 数据迁移失败 | 低 | 高 | 自动回滚 + 手动恢复指引 | 迁移后任务数不匹配 |
| TUI 性能不达标 | 中 | 中 | 延迟加载 + 分页 | 100+ 任务时 FPS < 30 |
| 开发时间超支 | 中 | 中 | 削减 M5/M6 范围，保留核心 | 任何 Milestone 延迟 > 3 天 |
| v1 用户不愿意升级 | 中 | 中 | 保留 v1 CLI 兼容层，渐进引导 | 升级率 < 50% @ 2周 |

### 4.2 范围削减预案

如果工期紧张，按以下顺序削减功能：

1. **首先削减**：M6（飞书升级）→ 飞书保持 v1 行为，Bot 暂不支持 Project
2. **其次削减**：M5（Spec 管理）→ Spec 用纯 Markdown 手动管理，不建工具
3. **再次削减**：M7（双向同步）→ 只做 GitHub 只读同步，本地编辑不写回
4. **最后削减**：M4（TUI 增强）→ TUI 只做最小 Project 适配，不新增视图

**不可削减**：M1 + M2 + M3（Workspace + GitHub 只读同步是 v2 的核心价值）

---

## 五、成功指标（Definition of Done）

### 5.1 技术指标

| 指标 | 目标 | 验证方式 |
|------|------|----------|
| 代码测试覆盖率 | > 60% | `go test -cover` |
| v1 → v2 迁移成功率 | 100% | 自动化迁移测试（10 组不同规模数据） |
| GitHub 同步成功率 | > 98% | 连续 7 天同步日志分析 |
| TUI 帧率 | > 30 FPS | 500 任务下手动测试 |
| CLI 命令响应时间 | < 200ms | `capture list` 冷启动时间 |

### 5.2 用户体验指标

| 指标 | 目标 | 验证方式 |
|------|------|----------|
| 新用户上手时间 | < 5 分钟 | 观察测试用户从安装到创建第一个 Task |
| v1 用户升级摩擦 | 零数据丢失 | 迁移测试 |
| GitHub Project 首次同步 | < 30 秒（100 items） | 计时测试 |
| TUI 学习成本 | 常用操作无需文档 | 用户测试 |

---

## 六、后续展望（v2.x 及以后）

### v2.3+ 可能方向

| 方向 | 描述 | 优先级 |
|------|------|--------|
| GitLab / Gitee 集成 | 扩展 remote 支持其他代码托管平台 | 低 |
| AI 自动分类 | 根据 Task 标题/描述自动建议 Stage 和 Priority | 中 |
| Web Dashboard | 基于浏览器的数据可视化（可选模块） | 低 |
| 团队协作 | 多用户共享 Workspace（通过 Git） | 低 |
| 时间追踪 | Task 级别的时间记录和报告 | 中 |
| 自动化规则 | "当 GitHub PR 合并时自动将 Task 标记为 done" | 中 |
