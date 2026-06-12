# Capture v2 核心结论（Core Conclusions）

> 整理时间：2026-06-04
> 性质：基于 Brainstorm 分析、可行性评估、架构设计和模块设计的最终结论汇总
> 用途：快速查阅所有关键决策，作为后续开发的第一参考文档

---

## 一、产品定位结论

### 1.1 目标用户

**核心用户是"需要轻量项目管理的开发者"**，不是专业产品经理，不是企业 PMO。

具体画像：
- 独立开发者：同时维护 2-5 个 side project
- 小团队 Tech Lead：2-8 人团队，不想用 Jira
- AI Agent 重度用户：需要追踪人机协作产生的工作项

### 1.2 一句话定位

> Capture v2 是「开发者在自己的机器上管理项目和任务的工具」。

**不是**：团队协作平台、在线 SaaS、Jira/Linear 的替代品。
**是**：本地优先的"个人工作空间"，与 GitHub/飞书形成互补。

### 1.3 核心价值主张

| 价值 | 说明 |
|------|------|
| **3秒记录** | 灵感来时，终端里一条命令完成记录，不打断心流 |
| **离线可用** | 所有数据在本地，无网络也能查看和编辑 |
| **全局视图** | 多个项目的任务在一个终端窗口里统一管理 |
| **数据主权** | 纯文本存储，用户完全拥有和控制自己的数据 |
| **GitHub原生** | 与开发者已有的工作流（GitHub Issues/Projects）无缝衔接 |

---

## 二、架构决策结论

### 2.1 核心概念层次

```
Workspace（工作空间）→ Project（项目）→ Task（任务）
         │
         └── Spec（规格文档，与 Project 关联）
```

- **Workspace 是物理目录**：`~/.capture/workspaces/<ws-id>/`
- **Project 是 Workspace 下的子目录**：`<ws>/projects/<proj-slug>/`
- **Task 是 Markdown 文件**：`<proj>/tasks/YYYY/MM/TASK-NNNNN.md`
- **Spec 是 Markdown 文件**：`<proj>/specs/SPEC-NNN-title.md`

### 2.2 存储架构（保持不变并扩展）

| 层级 | 格式 | 职责 |
|------|------|------|
| Source of Truth | Markdown + YAML Frontmatter | 人类可读、Git 友好、可手动修复 |
| Query Index | SQLite（per Workspace） | 快速过滤、分页、搜索 |
| 远端缓存 | SQLite + 内存缓存 | GitHub/飞书数据的本地副本 |

**关键决策**：每个 Workspace 独立 SQLite 数据库，避免全局锁，支持 Workspace 独立迁移。

### 2.3 远端集成策略

| 远端 | 优先级 | 方向 | 说明 |
|------|--------|------|------|
| GitHub Project | P0 | 只读（MVP）→ 双向（v2.3） | 核心差异化功能 |
| 飞书 Bitable | P1 | 保持 v1 行为 | 不破坏现有用户 |
| 飞书 Bot | P1 | 增加 Project 维度 |  Bot 命令扩展 |

**统一抽象**：所有远端通过 `internal/remote/Remote` 接口接入，可插拔扩展。

### 2.4 配置分层

优先级从高到低：环境变量 > Project 配置 > Workspace 配置 > 全局配置

```
~/.capture/config.yaml                 # 全局：编辑器、凭据、默认值
~/.capture/workspaces/<ws>/workspace.yaml   # Workspace：名称、默认 Project、远端覆盖
~/.capture/workspaces/<ws>/projects/<p>/project.yaml  # Project：GitHub repo、飞书表、自定义字段
```

---

## 三、技术选型结论

### 3.1 新增依赖

| 依赖 | 用途 | 版本 | 替代方案 | 选择理由 |
|------|------|------|----------|----------|
| `shurcooL/githubv4` | GitHub GraphQL API | latest | 手写 HTTP | 类型安全，官方 schema 对齐 |
| `golang.org/x/oauth2` | GitHub OAuth 认证 | latest | 无 | 标准库扩展，维护活跃 |

### 3.2 不引入的依赖（有意克制）

| 领域 | 不引入 | 原因 |
|------|--------|------|
| 模板引擎 | 不用 Helm/YAML 模板 | Go `text/template` 标准库足够 |
| ORM | 不用 GORM/sqlx | 现有手写 SQL 足够简单，避免抽象泄漏 |
| Web 框架 | 不用 Gin/Echo | v2 无 Web 界面计划 |
| 配置库 | 保持 Viper，不引入其他 | 已有技术债，够用 |

### 3.3 GitHub 认证方式优先级

1. **`gh auth status`**（首选）— 零配置，大多数开发者已安装
2. **`GITHUB_TOKEN` 环境变量** — 显式、安全、CI 友好
3. **`~/.config/gh/hosts.yml`** — gh CLI 的 token 存储位置
4. **手动配置** — `capture config set github.token xxx`

---

## 四、数据模型结论

### 4.1 v1 → v2 模型变更

**Task 模型新增字段**：
```go
WorkspaceID string  // 所属 Workspace
ProjectID   string  // 所属 Project
ExternalID  string  // GitHub Issue node_id / 飞书 record_id
ExternalURL string  // 远端链接
SyncVersion int     // 乐观锁/冲突检测
```

**v1 兼容性**：新增字段均为可选，旧文件读取时为空值。迁移时自动填充 `"default"` workspace 和 `"legacy"` project。

### 4.2 ID 格式保持不变

- Task：`TASK-NNNNN`（与 v1 一致）
- Spec：`SPEC-NNN`（新增）
- Project 通过 slug 标识（URL-friendly）

**不改为 `PROJ-TASK-NNNNN`**：避免破坏 v1 用户的肌肉记忆和既有文件引用。

### 4.3 新增核心模型

| 模型 | 关键字段 | 存储位置 |
|------|----------|----------|
| Workspace | id, name, path, default_project, remotes | `workspace.yaml` |
| Project | id, name, slug, github.repo, github.project_number, feishu.bitable_* | `project.yaml` |
| Spec | id, title, status, linked_tasks[], progress | `specs/SPEC-NNN-*.md` |

---

## 五、用户体验结论

### 5.1 CLI 命令设计原则

| 原则 | 说明 | 示例 |
|------|------|------|
| **常用操作保持极简** | 最常用的命令不超过 2 个单词 | `capture add "xxx"` |
| **高级操作显式声明** | 不常用的功能需要显式 flag | `capture add "xxx" --project p` |
| **交互式向导兜底** | 必填参数缺失时进入交互模式 | `capture workspace create` → 提示输入名称 |
| **上下文自动推断** | 减少用户的选择负担 | 自动从 git remote 推断 Project |

### 5.2 渐进式概念暴露

```
新用户视角：                    高级用户视角：
─────────────                  ─────────────
capture add                    capture add --project x --stage analysis
capture list                   capture list --workspace w --project x
capture kanban                 capture kanban --view stage
capture status                 capture project sync x --direction bidirectional

（只看到 Task）                （暴露 Workspace/Project/Sync/Spec 全部能力）
```

### 5.3 TUI 性能红线

| 指标 | 红线 | 测试场景 |
|------|------|----------|
| 启动时间 | < 500ms | 冷启动，100+ 任务 |
| 列表滚动 | > 30 FPS | 500 任务，快速上下滚动 |
| 搜索响应 | < 200ms | 500 任务，实时过滤 |
| 视图切换 | < 100ms | Kanban ↔ Stage ↔ List |

**不达标则削减功能**：宁可少一个视图，也不能让 TUI 卡顿。

---

## 六、实施策略结论

### 6.1 里程碑与版本映射

| 版本 | 里程碑 | 工期 | 核心交付 |
|------|--------|------|----------|
| v1.1 | Foundation | 1周 | 后台重构，无用户可见变化 |
| v1.2 | Workspace | 2周 | `workspace` + `project` 命令，v1 数据迁移 |
| **v2.0** | **GitHub Read** | **1.5周** | **GitHub Project 只读同步，大版本发布** |
| v2.1 | TUI + Spec | 1.5周 | 多视图 TUI、Spec 文档管理 |
| v2.2 | Feishu Upgrade | 0.5周 | Bot Project 支持、Bitable 多 Project |
| v2.3 | Bidirectional | 1周 | GitHub 双向同步、冲突解决 |

### 6.2 范围削减优先级（工期紧张时）

保留（不可削减）→ 依次削减：
1. ✅ **M1 Foundation + M2 Workspace**（v2 存在的基础）
2. ✅ **M3 GitHub Read**（核心差异化）
3. ⚠️ M4 TUI Enhancement（可削减为新视图，保留 Kanban 适配）
4. ⚠️ M5 Spec Management（可推迟到 v2.1）
5. ⚠️ M6 Feishu Upgrade（可保持 v1 行为）
6. ⚠️ M7 Bidirectional（可长期只读）

### 6.3 发布策略

- **渐进发布**：v1.1 → v1.2 → v2.0，不是一次性大爆炸升级
- **Feature Flag**：新功能默认关闭，用户显式开启
- **自动迁移**：首次运行 v2 时自动检测 v1 数据并迁移，无需手动命令
- **回滚通道**：迁移自动备份，提供 `capture migrate rollback`

---

## 七、风险结论

### 7.1 已识别的高风险项

| 风险 | 概率 | 影响 | 缓解措施 | 责任 Milestone |
|------|------|------|----------|---------------|
| GitHub GraphQL API 变更 | 中 | 高 | Remote 接口抽象 + 版本校验 + 缓存兜底 | M3 |
| v1 用户习惯断裂导致流失 | 中 | 高 | 默认行为保持 v1 一致 + 渐进暴露新功能 + 自动迁移 | M2 |
| TUI 性能不达标（>500任务） | 中 | 中 | SQLite 分页 + 虚拟滚动 + 压力测试 | M4 |
| 概念膨胀导致新用户困惑 | 中 | 高 | 渐进式概念暴露 + 默认 Workspace/Project 自动创建 | M2 |

### 7.2 关键假设

| 假设 | 验证方式 | 如果失败怎么办 |
|------|----------|---------------|
| 目标用户已安装 `gh` CLI | 新用户调研 | 提供详细的手动 token 获取指南 |
| GitHub Project v2 API 稳定 | 持续集成测试 | 降级为只读，禁用写入功能 |
| 开发者愿意管理本地目录数据 | 用户行为观察 | 提供自动备份到云端的可选插件 |
| v1 用户愿意升级 | 升级率统计 @ 2周 | 保留 v1 兼容层长期维护 |

---

## 八、关键成功因素（KSF）

### 8.1 必须成功的三件事

1. **第一次 GitHub 同步必须「哇」**
   - 用户运行 `capture sync github` 后，30 秒内看到自己项目的所有任务出现在终端
   - 零配置（对 gh CLI 用户）、自动发现 Project、清晰进度提示
   - **衡量**：新用户 7 日留存率 > 50%

2. **v1 用户升级必须「无感」**
   - `capture add` / `capture list` / `capture kanban` 行为与 v1 完全一致
   - 旧数据自动迁移，无需手动操作
   - **衡量**：v1 用户升级率 @ 2周 > 70%

3. **TUI 必须「丝滑」**
   - 500 任务下不卡顿、不闪烁、不崩溃
   - 视图切换流畅，搜索实时响应
   - **衡量**：TUI 使用时长 > 5 分钟/次

### 8.2 失败的三个信号

| 信号 | 含义 | 应对 |
|------|------|------|
| 新用户安装后 24 小时内未创建第一个 Task | 上手门槛过高 | 简化 init 流程，增加 onboarding |
| v1 用户升级后 3 天内回滚到 v1 | 迁移或体验有问题 | 紧急修复，提供 v1 兼容模式 |
| GitHub 同步成功率 < 90% | 集成不可靠 | 降级为只读，加强错误处理 |

---

## 九、与 v1 的关系结论

### 9.1 继承

| v1 能力 | v2 处理方式 |
|----------|------------|
| CLI 命令 | 完全保持，新增 `--workspace/--project` 可选参数 |
| Markdown 存储格式 | 完全保持，新增可选字段 |
| SQLite 索引 | 扩展 schema，兼容旧数据 |
| TUI Kanban | 保持三列看板，新增 Project 标识 |
| 飞书 Bot/Bitable | 保持 v1 行为，逐步扩展 |
| Task ID 格式（TASK-NNNNN）| 保持不变 |
| Stage 工作流 | 保持 9 阶段，Project 级可配置是后期扩展 |

### 9.2 废弃

| v1 功能 | 废弃原因 | 替代方案 |
|----------|----------|----------|
| 全局 `~/.capture/tasks/` 目录 | Workspace/Project 层级取代 | 自动迁移到 `workspaces/default/projects/legacy/` |
| `assign --terminal_session` | 使用频率极低，增加复杂度 | 从 CLI 中移除，保留字段兼容 |
| 飞书 Webhook 模式（考虑） | WebSocket 已足够，减少维护 | 标记为 deprecated，v2.3 移除 |

### 9.3 新增

| v2 新增能力 | 价值 |
|-------------|------|
| Workspace/Project 层级 | 多项目管理 |
| GitHub Project 同步 | 离线可用 + 统一视图 |
| Spec 文档管理 | 需求到代码的追踪 |
| Project 自动推断 | 零摩擦任务创建 |
| 多视图 TUI | Kanban / Stage / List / Spec |
| 同步冲突检测 | 多源数据一致性 |

---

## 十、一句话总结

> **Capture v2 的成败，取决于能否在「保持 v1 极简体验」的同时，让开发者第一次从终端看到自己的 GitHub Project 任务时，发出「哇」的感叹。**

其他一切都是锦上添花。
