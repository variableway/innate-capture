# Capture v2 可行性分析

> 分析时间：2026-06-04
> 分析维度：技术可行性、资源可行性、风险与依赖、数据兼容性

---

## 一、需求范围确认

### v2 新增核心需求

| 编号 | 需求 | 优先级 | 复杂度 |
|------|------|--------|--------|
| R1 | Workspace 管理：创建、切换、配置多个工作空间 | P0 | 中 |
| R2 | Project 管理：在 Workspace 下创建和管理项目 | P0 | 中 |
| R3 | Task 归属 Project：任务必须属于某个 Project，支持跨 Project 查询 | P0 | 中 |
| R4 | GitHub Project 集成：读取 GitHub Project 的 items、views、fields | P0 | 中 |
| R5 | Spec/PRD 文档管理：模板创建、编辑、关联 Project/Task | P1 | 中 |
| R6 | GitHub Project 双向同步：本地修改可写回 GitHub | P2 | 高 |
| R7 | 飞书 Bot Project 支持：Bot 可按项目查询和创建任务 | P1 | 低 |
| R8 | v1 数据自动迁移：无缝升级，不丢失历史数据 | P0 | 低 |

### 保持不变的 v1 能力

- CLI 命令行交互方式
- Markdown + YAML Frontmatter 存储格式
- SQLite 索引层
- 双写存储架构（Markdown source of truth + SQLite index）
- 飞书 Bitable 同步（需适配多 Project）
- TUI 看板（需适配多 Project）
- 飞书 Bot WebSocket/Webhook 模式

---

## 二、技术可行性分析

### 2.1 技术栈评估

| 技术/依赖 | 当前状态 | v2 需求 | 可行性 | 备注 |
|-----------|----------|---------|--------|------|
| Go 1.21+ | ✅ 已使用 | 继续使用 | ✅ 完全可行 | 无需升级 |
| Cobra | ✅ CLI 框架 | 新增 workspace/project 子命令 | ✅ 完全可行 | 原生支持多级命令 |
| bubbletea | ✅ TUI 框架 | 新增 Project 选择器、多视图 | ✅ 可行 | 已有 list/table 组件 |
| Viper | ✅ 配置管理 | 新增 Workspace/Project 配置 | ✅ 完全可行 | 支持多配置文件 |
| modernc.org/sqlite | ✅ SQLite 驱动 | 每个 Workspace 独立 DB | ✅ 可行 | 支持多数据库连接 |
| yaml.v3 | ✅ YAML 解析 | 继续用于 frontmatter | ✅ 完全可行 | 无需变更 |
| oapi-sdk-go/v3 | ✅ 飞书 SDK | 继续用于 Bot/Bitable | ✅ 完全可行 | 无需变更 |
| GitHub GraphQL API | ❌ 未使用 | 读取 GitHub Project v2 数据 | ✅ 可行 | GitHub Project v2 使用 GraphQL，Go 有 shurcooL/githubv4 |
| GitHub REST API | ❌ 未使用 | 写入 GitHub Project items | ✅ 可行 | 标准 REST API，go-github 库成熟 |
| Go text/template | ❌ 未使用 | Spec 文档模板引擎 | ✅ 完全可行 | 标准库，零依赖 |

### 2.2 GitHub Project 集成技术评估

GitHub Project v2（新版）使用 GraphQL API，这是最大的技术新领域。

**读取可行性**：
- GitHub GraphQL API 支持查询 Project 的 `items`（Issue/PR/Draft）
- 支持自定义字段查询（Status、Priority、Size 等）
- 支持 Views（Board/Table/Roadmap）的结构查询
- **限制**：GitHub Project v2 API 仍在演进中，部分字段可能变化

**写入可行性**：
- GraphQL mutations 支持：创建 item、更新字段值、添加 item 到 Project
- 需要 `project` scope 的 Personal Access Token (classic) 或 fine-grained token
- **限制**：API 速率限制（每小时 5000 点 for GraphQL），个人使用完全足够

**认证方式**：
- 推荐：`gh auth` CLI 的 token（用户已安装 GitHub CLI 的概率高）
- 备选：用户手动配置 `GITHUB_TOKEN` 环境变量
- 备选：读取 `~/.config/gh/hosts.yml` 中的 OAuth token

### 2.3 存储层演进可行性

当前 v1 存储结构：
```
~/.capture/
├── config.yaml
├── capture.db          (SQLite)
└── tasks/
    └── 2026/
        └── 04/
            └── TASK-00001.md
```

v2 目标存储结构：
```
~/.capture/                          # 全局配置
├── config.yaml
├── workspaces/
│   └── default/                     # 默认 Workspace
│       ├── workspace.yaml           # Workspace 配置
│       ├── capture.db               # Workspace 级 SQLite
│       ├── projects/
│       │   └── innate-capture/      # Project 目录
│       │       ├── project.yaml     # Project 配置（GitHub repo 关联等）
│       │       ├── specs/           # Spec/PRD 文档
│       │       │   └── 001-spec.md
│       │       └── tasks/           # Task Markdown 文件
│       │           └── 2026/04/TASK-00001.md
│       └── templates/               # 项目级模板
└── tasks/                           # v1 兼容：无 Project 归属的遗留任务
    └── ...
```

**可行性结论**：
- ✅ 纯文件系统操作，无技术壁垒
- ✅ Markdown 文件路径变更只需修改 `MarkdownStore.taskPath()` 方法
- ✅ SQLite 数据库可按 Workspace 隔离，每个 Workspace 独立 `NewSQLiteStore()`
- ⚠️ v1 任务需要迁移到 `default` Workspace 的某个 Project 下（或保留为 unassigned）

### 2.4 数据模型变更可行性

v1 `Task` 模型需新增字段：
```go
type Task struct {
    // ... 现有字段 ...
    WorkspaceID string `yaml:"workspace_id" json:"workspace_id"`  // 新增
    ProjectID   string `yaml:"project_id" json:"project_id"`      // 新增
    ExternalID  string `yaml:"external_id" json:"external_id"`    // 新增：GitHub Issue ID 等
    // ...
}
```

**YAML 兼容性**：yaml.v3 对未知字段默认忽略，新增字段不会破坏旧文件读取。但旧文件读取后 `WorkspaceID`/`ProjectID` 为空，需要在迁移时填充默认值。

**SQLite 兼容性**：使用 `ALTER TABLE ADD COLUMN`，支持 nullable，无迁移风险。

---

## 三、资源可行性分析

### 3.1 开发工作量估算

| 模块 | 功能点 | 预估人天 | 难度 |
|------|--------|----------|------|
| **Workspace 模型** | Workspace/Project 数据结构、验证、CRUD | 2 | 低 |
| **CLI 扩展** | `capture workspace *`、`capture project *` 命令族 | 2 | 低 |
| **存储层改造** | 多 Workspace 路径、多 Project 路径、迁移逻辑 | 3 | 中 |
| **Task 归属改造** | Task 新增 project_id、CLI/TUI 适配 | 2 | 中 |
| **GitHub 集成** | GraphQL 客户端、Project 读取、缓存 | 3 | 中 |
| **TUI 适配** | Project 选择器、多视图切换、分组展示 | 2 | 中 |
| **Spec 文档** | 模板引擎、创建命令、关联 Task | 2 | 低 |
| **飞书适配** | Bot 命令支持 Project 参数、Bitable 多表同步 | 2 | 低 |
| **数据迁移** | v1→v2 自动迁移、回滚、备份 | 1 | 低 |
| **测试** | 单元测试、集成测试、E2E 测试 | 2 | 中 |
| **文档** | 使用文档、开发文档 | 1 | 低 |
| **合计** | | **22 人天** | |

### 3.2 维护成本评估

| 维度 | v1 | v2 | 变化 |
|------|-----|-----|------|
| 代码量 | ~3500 行 Go | ~6000-7000 行 Go | +80% |
| 依赖数量 | ~15 个外部包 | +2 个（githubv4, go-github） | +13% |
| 测试覆盖 | 部分 | 需补充到 60%+ | 增加 |
| 配置复杂度 | 1 个 config.yaml | 1 个 global + N 个 workspace.yaml | 增加 |
| 远端集成点 | 1 个（飞书） | 2 个（飞书 + GitHub） | 翻倍 |

---

## 四、风险与依赖分析

### 4.1 高风险项

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| GitHub Project v2 GraphQL API 变更 | 中 | 高 | 抽象 `GitHubRemote` 接口，API 变更只影响适配层；加版本校验 |
| v1 数据迁移失败 | 低 | 高 | 迁移前自动备份；迁移失败自动回滚；提供 `capture migrate --dry-run` |
| TUI 性能退化（任务量 > 500） | 中 | 中 | SQLite 分页查询；TUI 虚拟滚动；压力测试 |
| 概念过多导致用户体验下降 | 中 | 高 | 渐进式概念暴露；默认 Workspace 自动创建；Project 默认为 default |

### 4.2 中风险项

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| GitHub Token 获取困难（非 gh CLI 用户） | 中 | 低 | 提供详细的 token 获取指引；支持多种 token 来源 |
| 飞书 Bitable 多 Project 同步复杂 | 中 | 中 | 每个 Project 对应 Bitable 一个子表；同步失败不影响本地 |
| Workspace 目录权限问题 | 低 | 中 | 详细的错误提示；自动 `chmod`；权限检查命令 |

### 4.3 外部依赖

| 依赖 | 状态 | 风险等级 |
|------|------|----------|
| GitHub GraphQL API 稳定性 | GitHub 官方 API，稳定性高 | 低 |
| 飞书 oapi-sdk-go/v3 | 官方 SDK，持续维护 | 低 |
| modernc.org/sqlite | 活跃维护，纯 Go | 低 |
| bubbletea/lipgloss | Charm 公司维护，非常活跃 | 低 |

---

## 五、数据兼容性矩阵

### 5.1 v1 → v2 数据迁移策略

```
v1 数据                  v2 数据
─────────               ─────────
config.yaml      →      ~/.capture/config.yaml（保留，添加 workspace 配置）
capture.db       →      ~/.capture/workspaces/default/capture.db
~/.capture/tasks/ →     ~/.capture/workspaces/default/projects/<default-project>/tasks/
                        （遗留任务归入 default workspace 的 default project）
```

### 5.2 迁移流程

```bash
# v2 首次运行时自动检测 v1 数据
capture init                    # 初始化 v2 目录结构
# 自动执行：
# 1. 备份 ~/.capture/ → ~/.capture/.backup-v1-<timestamp>/
# 2. 创建 default workspace
# 3. 创建 default project（命名为 "legacy" 或 "default"）
# 4. 移动所有 tasks 到 default project
# 5. 重建 SQLite 索引
# 6. 生成迁移报告
```

### 5.3 回滚策略

```bash
# 如果迁移后出现问题
capture migrate rollback        # 从备份恢复 v1 数据
# 或直接手动恢复：
rm -rf ~/.capture/
cp -r ~/.capture/.backup-v1-<timestamp> ~/.capture/
```

---

## 六、可行性结论

### 总体评估：✅ 高度可行

| 维度 | 评分 | 说明 |
|------|------|------|
| 技术可行性 | ⭐⭐⭐⭐⭐ | 全部技术方案成熟，无未知技术领域 |
| 资源可行性 | ⭐⭐⭐⭐☆ | ~22 人天工作量，合理可控 |
| 数据兼容性 | ⭐⭐⭐⭐⭐ | 纯文本存储，迁移路径清晰，可回滚 |
| 风险评估 | ⭐⭐⭐⭐☆ | 主要风险可控，有明确缓解措施 |
| 战略价值 | ⭐⭐⭐⭐⭐ | 显著扩展产品边界，提升用户粘性 |

### 关键成功因素

1. **保持 v1 的简洁性**：v2 新增概念不能破坏 v1 用户的既有使用习惯
2. **GitHub 集成的可靠性**：这是 v2 最大的差异化，第一次体验必须成功
3. **数据迁移的无感知**：用户升级后不应感受到任何摩擦
4. **渐进式发布**：MVP 只包含 Workspace + Project + GitHub 只读同步，后续迭代逐步添加功能
