# Top 10 Trending Brainstorm Skills（GitHub & Skill Marketplace 2026）

> 调研时间：2026-06-04
> 覆盖范围：GitHub Trending Repositories、Claude Code Skills Marketplace、GPT Store、开源 Agent 框架
> 筛选标准：Stars/Installs 数量、社区活跃度、与产品构思/创意发散的直接相关性

---

## 目录

| 排名 | 名称 | 类型 | Stars/Installs | 核心能力 |
|:---:|:---|:---|:---:|:---|
| 1 | **Superpowers** (`obra/superpowers`) | Claude Skill | 40.9K ⭐ | 多 Agent 开发工作流，含 `/brainstorm` |
| 2 | **STORM** (Stanford OVAL Lab) | GitHub 开源 | 28.1K ⭐ | AI 协作知识探索，生成深度研究报告 |
| 3 | **PM Skills Marketplace** (`phuryn/pm-skills`) | Claude Skill | 65 Skills 套件 | 产品经理专用 brainstorm-ideas 系列 |
| 4 | **grill-me** (`mattpocock/skills`) | Claude Skill | 156.2K Installs | 对计划和架构进行压力测试式提问 |
| 5 | **Creative Brainstorming Skill** | Claude/GPT Skill | 21 ⭐（新兴） | SCAMPER/TRIZ/Six Thinking Hats 系统化 |
| 6 | **andrej-karpathy-skills** (`forrestchang`) | Claude Skill | 144K ⭐ | 编码前的假设澄清与简化思维 |
| 7 | **OpenHands** (`All-Hands-AI`) | GitHub 开源 | 72.1K ⭐ | 自主软件开发 Agent，从需求到代码 |
| 8 | **MetaGPT** | GitHub 开源 | 50K+ ⭐ | 多角色协作（PM+Architect+Engineer） |
| 9 | **Co-STORM** (Stanford OVAL Lab) | GitHub 开源 | 绑定 STORM | 人机协作式头脑风暴，实时讨论 |
| 10 | **best-aeo-skill** (`metawhisp`) | Claude Skill | 热门付费 | AI 搜索优化视角的内容构思 |

---

## #1 Superpowers (`obra/superpowers`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code Skills Marketplace |
| **仓库** | `github.com/obra/superpowers` |
| **Stars** | 40.9K ⭐ / 3.1K Forks |
| **安装** | `npx skills add obra/superpowers` |
| **类型** | 多 Agent 开发工作流框架（含 brainstorm 能力） |

### 核心能力

Superpowers 是 Claude Skill 生态中**最完整的多 Agent 开发工作流**，核心技能链：

```
/brainstorm → /write-plan → /execute-plan → code-review → merge
```

- **`/brainstorm`**：通过结构化提问精炼想法，保存设计文档
- **`/write-plan`**：将设计拆分为 2-5 分钟可执行的任务
- **`/execute-plan`**：为每个任务派遣独立 Subagent，带规格合规检查
- **`using-git-worktrees`**：隔离分支，确保干净基线
- **`test-driven-development`**：强制 RED-GREEN-REFACTOR 纪律

### 为什么值得参考

> **关键洞察**：Superpowers 的 `/brainstorm` 不是自由发散，而是"结构化收敛"——通过提问迫使使用者澄清假设，再生成可执行计划。

对 Capture v2 的启示：
- Task 创建前可以增加"结构化提问"环节
- `capture add` 可以支持 `--brainstorm` 模式，通过 3-5 个问题帮助用户理清任务本质

---

## #2 STORM (Stanford OVAL Lab)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | GitHub 开源 |
| **仓库** | `github.com/stanford-oval/storm` |
| **Stars** | 28,130 ⭐ (2026.04 统计) |
| **机构** | Stanford OVAL Lab |
| **论文** | EMNLP 2024 |

### 核心能力

STORM = **Synthesis of Topic Outlines through Retrieval and Multi-perspective Question Asking**

工作流程：
```
用户输入主题
    ↓
AI 模拟"专家对话"（多视角提问）
    ↓
生成研究大纲（Outline）
    ↓
基于搜索引擎检索资料
    ↓
生成带引用的 Wikipedia 风格长文
```

### 扩展：Co-STORM

- **人机协作模式**：用户可以作为参与者加入 AI Agent 之间的讨论
- **思维导图**：动态将讨论组织为思维导图，降低认知负荷
- **用户评估**：78% 的人类评估者认为优于传统 RAG Chatbot

### 为什么值得参考

> **关键洞察**：STORM 的核心创新是"模拟多专家对话"来激发深度思考，而非单轮问答。

对 Capture v2 的启示：
- Spec 文档编写可以借鉴"多视角提问"机制
- `capture spec create` 后自动触发"PM / Engineer / Designer"三个视角的追问

---

## #3 PM Skills Marketplace (`phuryn/pm-skills`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code / Claude Cowork |
| **仓库** | `github.com/phuryn/pm-skills` |
| **规模** | 65 Skills + 36 Commands + 8 Plugins |
| **安装** | `claude plugin marketplace add phuryn/pm-skills` |

### Brainstorm 相关 Skills

| Skill | 用途 |
|-------|------|
| `brainstorm-ideas-existing` | 现有产品的多视角构思（PM/Designer/Engineer） |
| `brainstorm-ideas-new` | 新产品初期的发散构思 |
| `brainstorm-experiments-existing` | 为现有产品设计实验验证假设 |
| `brainstorm-experiments-new` | 新产品的 Lean Startup Pretotype 设计 |
| `identify-assumptions-*` | 识别 Value/Usability/Viability/Feasibility 风险假设 |
| `prioritize-assumptions` | Impact × Risk 矩阵优先级排序 |
| `opportunity-solution-tree` | Teresa Torres 的 OST 方法 |

### Commands（工作流）

```bash
/discover     # 完整发现周期：构思 → 假设 → 优先级 → 实验
/brainstorm   # 多视角构思（ideas|experiments × existing|new）
```

### 为什么值得参考

> **关键洞察**：PM Skills Marketplace 将散落的 PM 方法论（Teresa Torres、Marty Cagan、Alberto Savoia）编码为可复用的 Skill。

对 Capture v2 的启示：
- 可以将 Task Stage 工作流（inbox→mindstorm→analysis...）编码为 Skill
- 每个 Stage 对应一个引导式提问模板

---

## #4 grill-me (`mattpocock/skills`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code / Cursor / Codex / 35+ Agents |
| **作者** | Matt Pocock（TypeScript 领域知名专家） |
| **安装** | `npx skills add mattpocock/skills --skill grill-me` |
| **安装量** | 156.2K Installs |
| **合集 Stars** | 87.3K ⭐ |

### 核心能力

**grill-me** 是一个"压力测试式提问" Skill：

```
输入：一个功能规格或架构计划
输出：一系列尖锐问题，暴露假设和依赖链
```

**典型对话**：
```
User: "我要添加实时协作功能到笔记应用"
grill-me:
  → 你打算用 WebSocket 还是 SSE？为什么？
  → 冲突解决策略是什么？Operational Transform 还是 CRDT？
  → 如果用户同时编辑同一段文本，期望行为是什么？
  → 离线编辑如何与在线同步协调？
  → 这个功能的用户数量预期是多少？
```

### 为什么值得参考

> **关键洞察**：grill-me 的核心价值是"在编码前暴露错误假设"，避免"基于错误前提的优雅实现"。

对 Capture v2 的启示：
- `capture spec create` 后可以自动触发 grill-me 式审问
- `capture stage <task> analysis` 时增加"假设检查清单"

---

## #5 Creative Brainstorming Skill (MCP Market)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code / MCP Market |
| **来源** | `mcpmarket.com/tools/skills/creative-brainstorming-ideation` |
| **Stars** | 21 ⭐（新兴，但方法论完整） |
| **类型** | 系统化创意构思 Skill |

### 核心方法论

该 Skill 将传统创新方法论编码为 AI 工作流：

| 方法 | 用途 |
|------|------|
| **SCAMPER** | 产品改进的 7 个维度（Substitute/Combine/Adapt/Modify/Put to another use/Eliminate/Reverse） |
| **TRIZ** | 发明问题解决理论，打破技术矛盾 |
| **Six Thinking Hats** | 六顶思考帽，平衡决策视角 |
| **Morphological Matrix** | 形态分析，系统性组合解 |
| **反向思维** | 从"如何失败"反推成功策略 |
| **角色扮演** | 强制关联和随机刺激 |

### 工作流

```
Divergent（发散）→ Convergence（收敛）→ Evaluation（评估）
    ↓                    ↓                  ↓
  大量想法生成          聚类和筛选         Impact vs Feasibility 矩阵
```

### 为什么值得参考

> **关键洞察**：这是少数将经典创新方法论（TRIZ、Six Thinking Hats）系统性集成到 AI Agent 中的 Skill。

对 Capture v2 的启示：
- `capture mindstorm` 命令可以内置 SCAMPER 模板
- TUI 的 Mindstorm 视图可以提供"方法论选择器"

---

## #6 andrej-karpathy-skills (`forrestchang`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code Skills |
| **作者** | Forrest Chang（基于 Andrej Karpathy 的观察） |
| **Stars** | 144K ⭐（增长最快之一的 AI Workflow 仓库） |
| **安装** | `/plugin marketplace add forrestchang/andrej-karpathy-skills` |
| **类型** | 编码行为准则（Encoded Preference Skill） |

### 核心原则

Karpathy 在 2026.01 发表了关于 AI Coding Agent 痛点的 viral 观察，被编码为 4 条规则：

1. **Think Before Coding**：明确陈述假设，多解释存在时列出而非默默选一个
2. **Simplicity First**：最小代码解决问题，50 行能搞定不要 500 行
3. **Surgical Changes**：只碰必须改的地方，不改相邻代码
4. **Goal-Driven Execution**：定义成功标准并循环验证

### 为什么值得参考

> **关键洞察**：这 4 条规则虽然面向编码，但本质是"任何创造性工作前的思维纪律"。

对 Capture v2 的启示：
- 在 `capture stage <task> analysis` 阶段嵌入"假设澄清检查清单"
- Task 描述模板可以增加"成功标准"字段

---

## #7 OpenHands (`All-Hands-AI`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | GitHub 开源 |
| **仓库** | `github.com/All-Hands-AI/OpenHands` |
| **Stars** | 72,100 ⭐ |
| **定位** | 开源版 "Devin" |

### 核心能力

OpenHands 提供四种使用模式：

| 模式 | 说明 |
|------|------|
| **Python SDK** | 编程定义 Agent |
| **CLI** | 终端使用（对标 Claude Code / Codex） |
| **Desktop GUI** | React 前端界面 |
| **Cloud Platform** | 云端 + Slack/Jira/Linear 集成 |

### Brainstorm 相关能力

- **自主软件开发**：从自然语言描述到完整代码实现
- **多步工作流编排**：可以执行需要多个步骤的复杂任务
- **上下文保持**：长任务中的状态管理

### 为什么值得参考

> **关键洞察**：OpenHands 代表了"AI Agent 从辅助到执行"的趋势，其多模式架构（SDK/CLI/GUI/Cloud）值得参考。

对 Capture v2 的启示：
- Capture 的 Bot 模式（Webhook/WebSocket）可以借鉴 OpenHands 的多模式设计
- `capture dispatch` 功能可以向 OpenHands 的自主执行方向演进

---

## #8 MetaGPT

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | GitHub 开源 |
| **Stars** | 50K+ ⭐ |
| **定位** | 多角色协作 Agent 框架 |

### 核心能力

MetaGPT 模拟软件公司的组织架构：

```
产品经理（PM）→ 架构师（Architect）→ 项目经理（Project Manager）→ 工程师（Engineer）
     ↓                ↓                        ↓                      ↓
  写 PRD          写技术设计              拆 Task                 写代码
```

### 为什么值得参考

> **关键洞察**：MetaGPT 的核心创新是"角色分离"——不同 Agent 扮演不同角色，各司其职。

对 Capture v2 的启示：
- Spec 文档可以内置"角色视角"标签（PM view / Engineer view / Designer view）
- `capture spec review` 可以模拟多角色评审

---

## #9 Co-STORM (Stanford OVAL Lab)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | GitHub 开源（STORM 的扩展） |
| **机构** | Stanford OVAL Lab |
| **论文** | EMNLP 2024 |
| **用户** | 70,000+ |

### 核心能力

Co-STORM 是 STORM 的**人机协作扩展**：

- **AI Agent 圆桌讨论**：多个 AI Agent 模拟不同领域专家进行讨论
- **人类参与者**：用户可以观察或主动参与讨论，引导方向
- **实时思维导图**：讨论内容动态组织为可视化思维导图
- **未知发现**：帮助用户发现"自己不知道自己不知道"的知识盲区

### 为什么值得参考

> **关键洞察**：Co-STORM 解决了传统头脑风暴的"群体思维"问题——通过 AI 扮演"异见者"角色，强制多元化视角。

对 Capture v2 的启示：
- `capture mindstorm` 可以内置"Devil's Advocate"（魔鬼代言人）模式
- 为 Task 生成"反对理由"，帮助用户评估想法的可行性

---

## #10 best-aeo-skill (`metawhisp`)

### 基本信息

| 属性 | 详情 |
|------|------|
| **平台** | Claude Code / Cursor / Codex / 35+ Agents |
| **安装** | `npx skills add metawhisp/best-aeo-skill` |
| **类型** | AI 搜索优化（GEO/AEO）Skill |
| **热门度** | Claude Skill Marketplace 热门付费 Skill |

### 核心能力

best-aeo-skill 从"AI 如何引用内容"的视角帮助构思和优化：

- **ChatGPT 优化**：权威长文、共识性内容
- **Claude 优化**：精确归因、准确性优先
- **Perplexity 优化**：引用密集、学术和新闻源
- **Google Gemini**：社区聚焦、SEO + AI 信号混合
- **Google AI Overviews**：直接答案格式

### 为什么值得参考

> **关键洞察**：best-aeo-skill 代表了一个新兴趋势——"为 AI 消费而设计内容"，这与传统 SEO 截然不同。

对 Capture v2 的启示：
- Spec 文档模板可以考虑"AI 可读性"——清晰的标题层级、明确的成功标准、可验证的假设
- Capture 的 Markdown 输出本身就是"AI 友好"的（纯文本、结构化）

---

## 汇总对比表

| 排名 | 名称 | 平台 | Stars/Installs | 类型 | 核心方法论 |
|:---:|:---|:---|:---:|:---|:---|
| 1 | **Superpowers** | Claude Skill | 40.9K ⭐ | 工作流框架 | 结构化 brainstorm → 计划 → 执行 |
| 2 | **STORM** | GitHub | 28.1K ⭐ | 研究生成 | 多视角提问 + 检索增强 |
| 3 | **PM Skills Marketplace** | Claude Skill | 65 Skills | PM 方法论 | Teresa Torres / Cagan / Savoia |
| 4 | **grill-me** | Claude Skill | 156.2K Installs | 压力测试 | 假设暴露 + 依赖链澄清 |
| 5 | **Creative Brainstorming** | Claude/GPT | 21 ⭐ | 系统化创新 | SCAMPER / TRIZ / Six Hats |
| 6 | **karpathy-skills** | Claude Skill | 144K ⭐ | 行为准则 | 编码前思维纪律 |
| 7 | **OpenHands** | GitHub | 72.1K ⭐ | 自主 Agent | 自然语言 → 代码 |
| 8 | **MetaGPT** | GitHub | 50K+ ⭐ | 多角色协作 | PM/Architect/Engineer 角色分离 |
| 9 | **Co-STORM** | GitHub | 绑定 STORM | 人机协作 | AI 圆桌 + 实时思维导图 |
| 10 | **best-aeo-skill** | Claude Skill | 热门付费 | AI 搜索优化 | 多平台 AI 引用优化 |

---

## 趋势洞察

### Trend 1: 从"自由发散"到"结构化收敛"

2023-2024 的 AI Brainstorm 工具强调"生成更多想法"，2025-2026 的趋势是"帮助用户收敛到可执行方案"。

代表性：Superpowers `/brainstorm` → `/write-plan` 的流水线设计。

### Trend 2: 从"单 Agent"到"多 Agent 协作"

一个 AI 独自 brainstorm → 多个 AI Agent 扮演不同角色互相讨论。

代表性：MetaGPT 的角色分离、Co-STORM 的多 Agent 圆桌。

### Trend 3: 从"通用 Chat"到"编码方法论"

通用 AI 对话 → 将特定领域的方法论（PM 框架、TRIZ、Design Thinking）编码为 Skill。

代表性：PM Skills Marketplace、Creative Brainstorming Skill。

### Trend 4: 从"生成内容"到"暴露盲区"

AI 帮助用户生成想法 → AI 帮助用户发现自己没考虑到的东西。

代表性：grill-me 的假设暴露、Co-STORM 的"未知发现"。

---

## 附录：Capture v2 内置 Skill 的集成架构思路

> 关键问题：如果 Capture 要内置上述 Brainstorm Skill 能力，是否需要调用 Codex、Claude Code 等外部 CLI？
> 答案：**不直接调用，而是提供标准化的上下文输出，让用户/Agent 消费。**

---

### 问题本质

直接调用 `claude` 或 `codex` CLI 存在三个致命问题：

| 问题 | 说明 |
|------|------|
| **绑定性** | Capture 变成 Claude/Codex 的附属工具，失去独立性 |
| **脆弱性** | 外部 CLI 的版本变化、API 变更、输出格式变化都会破坏 Capture |
| **排他性** | 使用 ChatGPT、Gemini、本地 Ollama 的用户被排除在外 |

**正确的思路**：Capture 做好"数据的组织者和上下文的打包者"，AI Agent 做好"推理和生成"。两者通过标准接口连接，而非硬编码依赖。

---

### 核心设计原则：Capture 是上下文源，不是 Agent 控制器

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              设计原则                                         │
│                                                                             │
│   Capture 不调用 AI → Capture 输出上下文 → 任意 AI 工具消费上下文               │
│                                                                             │
│   capture context --task TASK-001        ──►   claude -p "分析这个任务"       │
│   capture context --spec SPEC-001        ──►   chatgpt "评审这个 Spec"        │
│   capture context --project innate-capture ──►   codex "生成实现计划"          │
│   capture context --workspace default    ──►   ollama "总结本周进度"           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 四种集成模式（从轻到重）

#### 模式一：管道式输出（Pipe Mode）—— 推荐默认

Capture 将项目/任务/Spec 的上下文输出为结构化文本，用户通过 Unix 管道传给任意 AI 工具。

```bash
# 基础用法：把 Task 上下文输出到终端
capture context --task TASK-001

# 管道传给 Claude Code
capture context --task TASK-001 | claude -p "分析这个任务的依赖关系和风险点"

# 管道传给 Codex
capture context --spec SPEC-001 | codex "基于这个 Spec 生成实现计划"

# 管道传给 ChatGPT（通过 openai CLI）
capture context --project innate-capture | openai chat "总结这个项目的当前状态"

# 管道传给本地 Ollama（完全离线）
capture context --workspace default | ollama run llama3.2 "本周有哪些高优先级任务？"
```

**`capture context` 输出格式**：

```markdown
# Capture Context Export
## Metadata
- Export Time: 2026-06-04T11:30:00Z
- Source: task TASK-001
- Workspace: default
- Project: innate-capture

## Task
- ID: TASK-001
- Title: "设计 Workspace API"
- Status: in_progress
- Stage: analysis
- Priority: high
- Description: |
    定义 Workspace 的数据结构和 CRUD 接口...

## Linked Spec
- SPEC-001: "Workspace 管理功能"
  - Status: approved
  - Acceptance Criteria:
    - [ ] 可以创建 Workspace
    - [ ] 可以切换 Workspace

## Related Tasks
- TASK-002: "实现 WorkspaceStore" (todo)
- TASK-003: "编写 workspace CLI 命令" (todo)

## GitHub Project Sync Status
- Last Sync: 2026-06-04T10:00:00Z
- Synced Items: 0/0
```

**优点**：
- ✅ 零绑定：不依赖任何特定 AI 工具
- ✅ Unix 哲学：符合开发者习惯
- ✅ 可组合：可以链式处理 `capture context | grep "high" | claude ...`
- ✅ 可审计：输出就是纯文本，用户完全知道给了 AI 什么上下文

---

#### 模式二：Prompt 模板库（Template Mode）—— 内置方法论

Capture 将 Top 10 Skill 中的方法论编码为本地 Prompt 模板，用户一键复制后贴到任意 AI 聊天窗口。

```bash
# 列出可用的 brainstorm 模板
capture template list
# →
#   brainstorm-scamper      SCAMPER 产品创新七维法
#   brainstorm-six-hats     六顶思考帽决策法
#   spec-grill-me           规格压力测试提问
#   task-first-principles   第一性原理分析
#   review-multi-role       多角色评审（PM/Engineer/Designer）

# 生成针对特定任务的 Prompt
capture template apply brainstorm-scamper --task TASK-001
# → 输出完整 Prompt + Task 上下文，用户直接复制到 ChatGPT/Claude

# 或者直接输出到剪贴板
capture template apply spec-grill-me --spec SPEC-001 | pbcopy
# → "已复制到剪贴板，直接粘贴到 AI 对话框即可"
```

**模板本质**：Go `text/template` 文件，存储在 `~/.capture/templates/`。

```
~/.capture/templates/
├── brainstorm-scamper.md       # SCAMPER 模板
├── brainstorm-six-hats.md      # 六顶思考帽模板
├── spec-grill-me.md            # 压力测试模板
├── review-multi-role.md        # 多角色评审模板
└── custom/                     # 用户自定义模板
    └── my-template.md
```

**模板示例**（`brainstorm-scamper.md`）：

```markdown
# SCAMPER 分析：{{.Task.Title}}

请对以下产品/功能进行 SCAMPER 分析，从 7 个维度提出改进思路：

## 目标功能
{{.Task.Title}}
{{.Task.Description}}

## 7 维分析

### S - Substitute（替代）
这个功能中，什么可以被替换？（技术、流程、人员、工具）

### C - Combine（组合）
这个功能可以和什么合并？（其他功能、产品、流程）

### A - Adapt（适应）
从哪里可以借鉴类似方案？（竞品、其他行业、开源项目）

### M - Modify（修改）
如果改变某个属性（大小、频率、顺序），会发生什么？

### P - Put to another use（另作他用）
这个功能还有什么其他用途？

### E - Eliminate（消除）
如果去掉某个部分，会变得更好吗？

### R - Reverse（反转）
如果反转流程或角色，会产生什么新思路？

## 输出格式
对每个维度，请给出 2-3 个具体可执行的想法，按可行性排序。
```

**优点**：
- ✅ 零外部依赖：不需要安装任何 AI CLI
- ✅ 方法论可沉淀：团队可以共享和积累 Prompt 模板
- ✅ 完全透明：用户知道 AI 收到了什么 Prompt

---

#### 模式三：Agent Dispatch 扩展（Dispatch Mode）—— v1 概念延伸

Capture v1 已有 `capture assign` 命令记录 Agent 执行上下文。v2 可以扩展为"向任意 Agent 分发任务"。

```bash
# v1 的 assign 命令（记录上下文）
capture assign TASK-001 --agent codex --model gpt-5 --repo ~/workspace/project

# v2 的扩展：dispatch 命令（分发任务到 Agent）
capture dispatch TASK-001 --agent claude
# → 将 Task 上下文打包，打开 Claude Code 并自动加载上下文

capture dispatch SPEC-001 --agent codex --prompt "生成实现计划"
# → 将 Spec 上下文传给 Codex，生成代码实现

# dispatch 到本地脚本（通用接口）
capture dispatch TASK-001 --agent ~/scripts/my-ai-script.sh
# → 将上下文作为环境变量和 stdin 传给脚本
```

**`dispatch` 的本质**：

```go
func Dispatch(task *model.Task, agentConfig AgentConfig) error {
    // 1. 打包上下文
    context := buildContext(task)
    
    // 2. 根据 agent 类型选择分发方式
    switch agentConfig.Type {
    case "claude":
        // 写入临时文件，调用 claude -f
        return dispatchToClaude(context, agentConfig)
    case "codex":
        // 通过 stdin 管道
        return dispatchToCodex(context, agentConfig)
    case "script":
        // 调用用户自定义脚本
        return dispatchToScript(context, agentConfig)
    case "mcp":
        // 通过 MCP 协议
        return dispatchToMCP(context, agentConfig)
    }
}
```

**优点**：
- ✅ 自动化：一键将任务上下文传给 Agent
- ✅ 可扩展：支持任意 Agent（通过脚本或 MCP）
- ✅ 可追踪：Capture 记录每次 dispatch 的历史

---

#### 模式四：MCP Server 模式（MCP Mode）—— 未来方向

Capture 自身可以作为 MCP Server，向外部 Agent 暴露项目/任务/Spec 数据。

```
┌─────────────────────────────────────────────────────────────┐
│                        MCP 架构                              │
│                                                             │
│   ┌──────────────┐         MCP Protocol          ┌────────┐ │
│   │  Claude Code │  ◄──────────────────────────►  │ Capture│ │
│   │  / Codex     │    resources/tools/prompts    │ Server │ │
│   │  / Cursor    │                               │        │ │
│   └──────────────┘                               └────────┘ │
│                                                          │  │
│   Capture 暴露的 MCP Resources:                           │  │
│   ├── capture://workspace/{id}/tasks                     │  │
│   ├── capture://project/{id}/specs                       │  │
│   ├── capture://task/{id}                                │  │
│   └── capture://context/current                          │  │
│                                                          │  │
│   Capture 暴露的 MCP Tools:                              │  │
│   ├── capture_create_task                                │  │
│   ├── capture_update_task_status                         │  │
│   ├── capture_list_tasks                                 │  │
│   └── capture_export_context                             │  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**用户使用体验**：

```bash
# 1. 启动 Capture MCP Server
capture mcp serve
# → MCP Server running on stdio/sse

# 2. 在 Claude Code 中配置 MCP
# {
#   "mcpServers": {
#     "capture": {
#       "command": "capture",
#       "args": ["mcp", "serve"]
#     }
#   }
# }

# 3. 在 Claude Code 中自然语言使用
# User: "帮我列出 innate-capture 项目中所有高优先级的待办任务"
# Claude: [通过 MCP 调用 capture_list_tasks] → 返回结果 → 展示给用户
```

**优点**：
- ✅ 双向交互：Agent 不仅可以读 Capture 数据，还可以创建/更新 Task
- ✅ 标准协议：MCP 是 Anthropic 推动的开放标准，生态正在快速增长
- ✅ 实时同步：Agent 操作即时反映到 Capture 的 Markdown/SQLite 中

---

### 四种模式的选择策略

| 模式 | 复杂度 | 外部依赖 | 适用场景 | 推荐优先级 |
|------|--------|----------|----------|-----------|
| **模式一：管道式输出** | 低 | 无 | 所有用户，尤其是多 AI 工具用户 | P0（MVP 必须有） |
| **模式二：Prompt 模板** | 低 | 无 | 方法论沉淀，团队共享 | P0（MVP 必须有） |
| **模式三：Agent Dispatch** | 中 | 可选 | 自动化工作流，重度 AI 用户 | P1（v2.1） |
| **模式四：MCP Server** | 高 | MCP Client | 深度 AI 集成，双向同步 | P2（v2.3+） |

---

### 结论：Capture 的 Skill 集成哲学

```
┌────────────────────────────────────────────────────────────────────┐
│                                                                    │
│   ❌ 不要做：Capture 内部调用 `claude` 或 `codex` CLI               │
│                                                                    │
│   ✅ 要做：                                                        │
│      1. Capture 输出标准化的、人类可读的上下文                      │
│      2. 用户选择用哪个 AI 工具消费这个上下文                        │
│      3. 可选地，Capture 提供一键分发（dispatch）降低操作成本        │
│      4. 未来通过 MCP 协议实现深度双向集成                           │
│                                                                    │
│   Capture 的定位：「开发者工作流的上下文中枢」                       │
│   不是「某个 AI 工具的包装器」                                      │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

---

## 对 Capture v2 的 Skill 设计建议

基于以上 Top 10 的观察和集成架构思路，Capture v2 的 Skill 能力可以这样设计：

### Phase 1（MVP）：上下文输出 + 模板库

```bash
# 1. 上下文输出（模式一）
capture context --task TASK-001

# 2. Prompt 模板（模式二）
capture template list
capture template apply brainstorm-scamper --task TASK-001 | pbcopy
```

### Phase 2（v2.1）：Agent Dispatch

```bash
# 3. 一键分发（模式三）
capture dispatch TASK-001 --agent claude --prompt "分析依赖关系"
```

### Phase 3（v2.3+）：MCP Server

```bash
# 4. MCP 服务（模式四）
capture mcp serve
```

### 内置模板清单（参考 Top 10 Skill）

```
~/.capture/templates/
├── brainstorm/
│   ├── scamper.md              # #1 参考 Superpowers / Creative Brainstorming
│   ├── six-thinking-hats.md    # 参考 Creative Brainstorming Skill
│   ├── first-principles.md     # 参考 First Principles Thinking
│   └── reverse-brainstorm.md   # 参考 Reverse Brainstorming
├── review/
│   ├── grill-me.md             # #2 参考 grill-me
│   ├── karpathy-discipline.md  # 参考 karpathy-skills
│   └── multi-role.md           # #3 参考 MetaGPT / PM Skills
├── spec/
│   ├── storm-outline.md        # 参考 STORM 的大纲生成
│   └── acceptance-criteria.md  # 参考 PM Skills Marketplace
└── custom/                     # 用户自定义模板
```

### 命令映射

```bash
# 头脑风暴（替代自由发散）
capture mindstorm "新功能想法" --template scamper

# 规格评审（替代随意看看）
capture spec review SPEC-001 --template grill-me

# 任务分析（替代直接开干）
capture task analyze TASK-001 --template first-principles

# 导出上下文（给任意 AI 使用）
capture context --task TASK-001 --format markdown | claude -p "帮我分析"
```
