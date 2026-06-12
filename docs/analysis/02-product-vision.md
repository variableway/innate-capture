# Capture v2 产品愿景与 Brainstorm 分析

> 分析时间：2026-06-04
> 分析方法：Starbursting + First Principles + Reverse Brainstorming + SCAMPER
> 分析对象：Capture 从 v1（个人任务捕捉工具）到 v2（项目/产品工作空间管理平台）的演进

---

## 一、产品定位回顾

### v1 现状

Capture v1 是一个**个人任务捕捉与轻量管理工具**：

- **入口**：Terminal CLI + TUI 看板 + 飞书 Bot
- **存储**：本地 Markdown + YAML Frontmatter（Source of Truth）+ SQLite（索引查询）
- **核心对象**：Task（单一层级，无项目归属）
- **工作流**：inbox → mindstorm → analysis → planning → prd → tasks → dispatch → execution → review
- **同步**：单向/双向同步到飞书多维表格
- **用户场景**：开发者随手记录灵感、跟踪个人待办、Agent 任务分派

### v2 目标

将 Capture 升级为**面向开发者和技术团队的「产品工作空间」**：

- **新增核心概念**：Workspace（工作空间）→ Project（项目）→ Task（任务）三级结构
- **新增集成**：GitHub Project 深度集成（读取进度、状态同步）
- **新增能力**：Spec/PRD 文档的模板化创建、版本管理与跨项目分发
- **目标用户**：独立开发者、小团队 Tech Lead、需要轻量项目管理的工程师

---

## 二、Starbursting 分析：充分理解问题空间

### 2.1 Who — 谁会使用/维护/受影响？

| 用户角色 | 使用场景 | 痛点 |
|----------|----------|------|
| **独立开发者** | 同时维护多个 side project，需要统一入口查看所有项目的任务和进度 | 在多个 GitHub repo 之间切换查看 Issues/Projects 太分散 |
| **小团队 Tech Lead** | 需要跟踪团队任务进度，但不想用 Jira 这样重的工具 | Linear 不错但需要联网，想有一个本地优先 + 可选同步的方案 |
| **AI Agent 使用者** | 用 Claude/Codex 等 Agent 写代码，需要记录 Agent 产生的任务和上下文 | Agent 产生的 TODO 散落在对话中，无法系统化跟踪 |
| **产品经理（兼职）** | 在开发团队中兼任 PM，需要写 PRD、Spec，并跟踪实现进度 | Notion 太重，飞书文档不好和代码仓库关联 |
| **飞书用户** | 已经习惯用飞书 Bot 记录灵感 | 希望在飞书里也能看到项目维度的任务列表 |

**关键洞察**：核心用户是"需要轻量项目管理的开发者"，而非专业 PM。工具必须保持"开发者友好"（CLI 优先、Git 友好、Markdown 优先）。

### 2.2 What — 要解决什么核心问题？

v1 解决的问题：**"灵感来了，如何 3 秒内记录下来？"**

v2 要解决的问题：**"我同时推进 3 个项目，每个项目有需求文档、开发任务、GitHub Issues，如何在一个地方看到全貌？"**

具体子问题：
1. 如何在本地统一管理多个项目的任务，而不依赖网络？
2. 如何让项目的需求文档（Spec/PRD）与代码仓库、任务列表保持关联？
3. 如何让 GitHub Project 的进度在本地可见，甚至可编辑后同步回去？
4. 如何让飞书 Bot 也支持按项目查询和创建任务？

### 2.3 Why — 为什么现在做？

- **v1 的瓶颈**：单任务列表在任务量超过 50 后变得难以管理，没有项目维度的筛选和分组
- **AI 编程时代**：开发者与 AI Agent 协作的频率越来越高，需要一个工具来管理"人机协作产生的工作项"
- **本地优先趋势**：Notion/Linear 都是 SaaS，数据不在本地。开发者越来越重视数据主权
- **GitHub Project 不足**：GitHub Project 功能越来越强，但在终端中无法快速查看和操作

### 2.4 Where — 在哪里使用/部署？

| 场景 | 数据存储 | 交互方式 |
|------|----------|----------|
| 本地开发 | `~/.capture/` 目录下，按 Workspace 组织 | CLI + TUI |
| 与 GitHub 协作 | GitHub Project 数据作为远端源，本地缓存 | `capture sync github` |
| 与飞书协作 | 飞书多维表格作为远端源，本地缓存 | Bot 消息 + `capture sync feishu` |
| 跨设备同步 | Git 仓库（Workspace 配置可提交到 Git） | `git push/pull` |

**关键洞察**：Workspace 的配置和任务数据应该可以作为一个 Git 仓库管理，实现"配置即代码"。

### 2.5 When — 什么时候发布/迭代？

| 里程碑 | 时间 | 交付物 |
|--------|------|--------|
| MVP | 2-3 周 | Workspace + Project 基础模型 + CLI 命令 + TUI 适配 |
| v2.1 | 1-2 周 | GitHub Project 只读同步 |
| v2.2 | 1-2 周 | Spec 文档模板 + 分发机制 |
| v2.3 | 2-3 周 | GitHub Project 双向同步 + 冲突解决 |
| v2.4 | 1-2 周 | 飞书 Bot Project 支持 + 飞书多维表格多 Project 同步 |

### 2.6 How — 如何衡量成功？

| 指标 | 目标 | 测量方式 |
|------|------|----------|
| 任务完成率 | 周任务完成率 > 60% | 本地 SQLite 统计 |
| 项目活跃度 | 每个 Workspace 平均管理 2+ Project | 数据分析 |
| GitHub 同步成功率 | > 95% | 同步日志 |
| 用户留存 | 连续 4 周每周至少使用 1 次 | 本地匿名埋点（可选） |

---

## 三、First Principles 分析：从根本约束重新构建

### 3.1 拆解到基本事实

**事实 1**：开发者的时间极度碎片化，任何需要超过 10 秒的操作都会被放弃。
**事实 2**：开发者对工具的首要偏好是"不离开终端"。
**事实 3**：Markdown 是开发者最熟悉的文档格式，Git 是最熟悉的数据管理工具。
**事实 4**：GitHub 是代码托管的事实标准，GitHub Project 是多数开发者的项目管理工具。
**事实 5**：飞书是中国技术团队的主要协作工具，Bot 是飞书内最高效的交互方式。
**事实 6**：本地数据必须能被用户完全理解和控制（纯文本、可手动编辑）。

### 3.2 从事实推导设计原则

| 事实 | 推导原则 | 设计决策 |
|------|----------|----------|
| 事实 1 | 操作必须极简 | CLI 命令保持短命名，常用操作 1-2 个参数 |
| 事实 2 | 终端优先，GUI 辅助 | CLI 和 TUI 是第一公民，Web 界面暂不考虑 |
| 事实 3 | Markdown + Git 原生支持 | 所有文档和配置都是 Markdown/YAML，Workspace 可 git init |
| 事实 4 | GitHub 是默认远端 | GitHub Project 集成优先级最高 |
| 事实 5 | 飞书是重要补充 | 飞书 Bot 和 Bitable 同步继续维护，但不作为核心 |
| 事实 6 | 数据透明可审计 | 所有数据存储为人类可读的文本文件，无二进制封闭格式 |

### 3.3 差异化定位

```
┌─────────────────────────────────────────────────────────────┐
│                    项目管理工具谱系                           │
├─────────────────────────────────────────────────────────────┤
│  重 ↔─────────────────────────────────────────────→ 轻      │
│                                                              │
│  Jira    Notion    Linear    GitHub Project    Capture v2    │
│   │        │         │            │              │          │
│   │        │         │            │              │          │
│  企业级   文档库    云端优先    代码仓库自带     本地优先    │
│  流程重   功能全    体验好      功能中等         终端原生    │
│                                                              │
│  Capture v2 的差异化 = 本地优先 + 终端原生 + Markdown 原生    │
│                       + GitHub 深度集成 + AI Agent 友好       │
└─────────────────────────────────────────────────────────────┘
```

**一句话定位**：Capture 是"开发者在自己的机器上管理项目和任务的工具"，不是"团队在线协作平台"。

---

## 四、Reverse Brainstorming：识别失败因素

### 4.1 如何确保 Capture v2 彻底失败？

1. **概念膨胀**：引入 Workspace/Project/Repo/Task/Issue/Spec/PRD 等 10+ 概念，新用户完全无法理解
2. **数据迁移灾难**：v1 用户的任务数据无法平滑迁移到 v2，导致数据丢失或重复
3. **GitHub 集成不可靠**：同步经常失败，错误信息晦涩，用户不敢依赖
4. **TUI 性能崩溃**：任务量超过 100 后 TUI 卡顿、闪烁、崩溃
5. **存储格式锁定**：引入二进制格式或复杂的数据库 schema，用户无法手动修复数据
6. **飞书同步断裂**：v2 改动导致飞书 Bot 和 Bitable 同步失效
7. **CLI 命令过长**：创建项目任务的命令需要 5+ 个 flag，用户记不住
8. **文档模板僵化**：Spec 模板只有一种格式，无法适配不同团队习惯

### 4.2 反转后的防护策略

| 失败因素 | 防护策略 | 具体措施 |
|----------|----------|----------|
| 概念膨胀 | 渐进式概念暴露 | 默认只显示 Task 列表，Project/Workspace 筛选器默认折叠 |
| 数据迁移灾难 | 自动迁移 + 回滚 | `capture migrate` 命令，自动备份 v1 数据，失败可恢复 |
| GitHub 集成不可靠 | 只读先行 + 缓存 | v2.1 只做只读同步，所有数据本地缓存，失败不影响本地使用 |
| TUI 性能崩溃 | 分页 + 虚拟化 | SQLite LIMIT/OFFSET 分页，TUI 只渲染可视区域 |
| 存储格式锁定 | 纯文本承诺 | 所有数据继续用 Markdown + YAML，新增 JSON 仅用于 API 缓存 |
| 飞书同步断裂 | 兼容层保留 | v1 的 Feishu/Bot 模块作为 legacy 兼容层保留，逐步迁移 |
| CLI 命令过长 | 交互式向导 | `capture project task add` 支持交互式 prompts，也支持 flags |
| 文档模板僵化 | 模板可扩展 | 模板用 Go text/template，用户可在 `~/.capture/templates/` 自定义 |

---

## 五、SCAMPER 分析：系统化功能创新

### 5.1 Substitute（替代）

| 现有方案 | 替代方案 | 价值 |
|----------|----------|------|
| 飞书多维表格作为远端 | GitHub Project 作为首要远端 | 开发者更常用 GitHub |
| SQLite 单一数据库 | 每个 Workspace 独立 SQLite | 支持 Workspace 级别的隔离和迁移 |
| YAML frontmatter | YAML frontmatter（保持） | 用户已熟悉，无需改变 |
| 手动 stage 切换 | 基于规则自动 stage 推荐 | 根据 GitHub PR 状态自动建议 stage |

### 5.2 Combine（组合）

| 组合 A + B | 组合后功能 |
|------------|------------|
| Task + GitHub Issue | 本地 Task 可与 GitHub Issue 双向关联，创建 Task 时可选择同步创建 Issue |
| Spec 文档 + Task | Spec 文档中的 "Acceptance Criteria" 自动生成为子任务 |
| Workspace + Git Repo | Workspace 目录可以直接 `git init`，配置和任务数据版本化 |
| TUI + Project 视图 | TUI 新增 Project 选择器，看板可按 Project 过滤 |

### 5.3 Adapt（适应/借鉴）

| 借鉴来源 | 借鉴内容 | 适配到 Capture |
|----------|----------|---------------|
| Linear | Issue 的状态和优先级设计 | 保留现有 TaskStatus/TaskPriority，新增 Project-level 自定义状态 |
| Obsidian | 本地优先 + 链接式知识管理 | Spec 文档支持 WikiLinks `[[TASK-00001]]` 自动关联任务 |
| GitHub Projects | 自定义字段（Custom Fields） | Project 可配置自定义字段，映射到 Task frontmatter |
| make | Makefile 式的任务依赖 | Task 支持 `depends_on` 字段，展示依赖关系 |

### 5.4 Modify（修改）

| 现有功能 | 修改方向 |
|----------|----------|
| 固定 9 阶段工作流 | Project 可配置自定义 stage 序列 |
| 单任务列表 | 支持 Project 内 Task 分组（如 Sprint、Epic） |
| 纯本地存储 | 支持 Git 仓库作为 Workspace 的远端备份 |
| TUI 三列看板 | TUI 支持多视图：Kanban（按 status）、Pipeline（按 stage）、List（按 project） |

### 5.5 Put to another use（另作他用）

| 现有能力 | 新用途 |
|----------|--------|
| Markdown 存储引擎 | 不仅存 Task，也存 Spec/PRD/Meeting Notes |
| 飞书 Bot 消息解析 | 解析飞书消息中的 URL，自动提取 GitHub Issue/PR 信息创建任务 |
| YAML frontmatter | 作为 Project 配置文件的格式（替代 JSON） |
| idgen（TASK-NNNNN） | 扩展为 `PROJ-TASK-NNNNN` 格式，支持跨 Project 唯一 ID |

### 5.6 Eliminate（消除）

| 可消除的 | 理由 |
|----------|------|
| v1 的 `assign` 命令中的 terminal_session | 使用频率低，增加复杂度 |
| 飞书 Webhook 模式（考虑） | WebSocket 模式已足够，减少维护负担 |
| Bitable 的 "双向写回" 默认开启 | 默认只读同步，双向需显式配置 |
| 全局任务列表（默认） | 默认显示当前 Project 的任务，全局列表作为可选视图 |

### 5.7 Reverse（反转）

| 反转前 | 反转后 |
|--------|--------|
| 用户手动创建 Task | Task 从 GitHub Issues/PRs/飞书消息自动导入 |
| 本地数据同步到远端 | 远端数据拉取到本地，本地编辑后选择性推送 |
| Task 驱动文档 | 文档（Spec）驱动 Task 生成 |
| 用户选择 Project | 根据当前 Git 仓库自动推断 Project |

---

## 六、总结：v2 核心设计信条

```
┌────────────────────────────────────────────────────────────┐
│  1. 本地优先（Local First）                                 │
│     所有数据首先存储在本地，远端同步是可选增强                │
├────────────────────────────────────────────────────────────┤
│  2. 终端原生（Terminal Native）                             │
│     CLI 和 TUI 是第一公民，所有功能都必须在终端可用           │
├────────────────────────────────────────────────────────────┤
│  3. Markdown 原生（Markdown Native）                        │
│     所有文档、配置、任务存储都是人类可读的纯文本               │
├────────────────────────────────────────────────────────────┤
│  4. Git 友好（Git Friendly）                                │
│     Workspace 可作为一个 Git 仓库管理，数据可版本化            │
├────────────────────────────────────────────────────────────┤
│  5. 渐进复杂度（Progressive Disclosure）                    │
│     新用户只看到 Task，高级用户才暴露 Project/Workspace        │
├────────────────────────────────────────────────────────────┤
│  6. 远端可插拔（Pluggable Remotes）                         │
│     GitHub/飞书/其他远端通过统一接口接入，可独立开关            │
└────────────────────────────────────────────────────────────┘
```
