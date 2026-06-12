# Capture v2 核心洞察（Core Insights）

> 整理时间：2026-06-04
> 来源：Brainstorm 分析、架构设计、可行性评估过程中的深层发现
> 性质：这些洞察未必显而易见，但对产品成败具有决定性影响

---

## Insight 1: Capture 的真正竞品不是 Linear/Jira，而是 "开发者自己的记忆"

### 表面认知
Capture 是一个任务管理工具，竞品是 Linear、Jira、Notion、GitHub Issues。

### 深层洞察
**Capture 的核心竞品是"开发者的脑内记忆 + 散落在各处的 TODO 注释"。**

大多数开发者并不需要一个"完美的项目管理工具"，他们需要的是：
- 灵感闪过时，3 秒内完成记录（不打开浏览器、不登录、不选模板）
- 离线时也能查看和编辑自己的任务
- 第二天打开电脑时，昨天记的东西还在，且能找到

这就是为什么 v1 的 CLI `capture add "xxx"` 比任何 Web 工具都更有价值 —— 它消除了"记录成本"。

### 对 v2 的启示
- **Workspace/Project 不能增加记录成本**：`capture add` 必须保持 3 秒完成，Project 自动推断
- **本地优先不是功能，是核心价值主张**：用户选择 Capture 是因为数据在自己电脑上
- **不要试图"替代 Linear"**：Capture 是 Linear 的补充，而非替代。用户在 Linear 上协作，在 Capture 上思考和记录。

### 设计体现
```bash
# 坏的：增加记录摩擦
capture add "修复bug" --workspace personal --project my-app --stage inbox

# 好的：保持极简，Project 自动推断（从当前 git repo 或默认值）
capture add "修复bug"
# → 自动归属到当前 Workspace 的当前 Project
```

---

## Insight 2: GitHub Project 集成的真正价值不是"同步"，而是"离线缓存 + 统一视图"

### 表面认知
GitHub Project 集成 = 把 GitHub 上的 tasks 拉下来，本地编辑后再推回去。

### 深层洞察
**GitHub Project 集成的真正价值是：让开发者在不打开浏览器的情况下，拥有所有项目任务的全局视图。**

开发者每天面临的真实场景：
- 打开了 5 个 GitHub repos，每个都有 Project board
- 需要逐个点击才能知道有哪些待办
- 网络慢时 GitHub 页面加载 3-5 秒
- 想快速看一下"这周所有项目有哪些高优先级任务"几乎不可能

Capture v2 解决的是：**一次 `capture sync github`，之后所有任务在终端本地可查，0 延迟。**

### 对 v2 的启示
- **只读同步（MVP）比双向同步更有价值**：用户最痛的是"看不到"，不是"改不了"
- **缓存层是关键基础设施**：不是优化，是核心功能。用户期望同步一次后，断网也能用
- **统一过滤查询是杀手级体验**：`capture list --status todo --priority high` 跨所有 Project 查询

### 设计体现
```
GitHub Project ──sync──► Local Cache (SQLite) ──query──► CLI/TUI
    │                                              │
    │  用户大部分时间在这里操作                     │  0 延迟，离线可用
    │  （浏览、筛选、搜索）                         │
    │                                              │
    └── 只在需要同步时联网 ─────────────────────────┘
```

---

## Insight 3: Workspace 的物理目录设计，本质是在卖"数据主权"

### 表面认知
Workspace 用物理目录存储是为了方便和清晰。

### 深层洞察
**让 Workspace 成为物理目录 + Git 仓库，是在满足开发者对"数据完全可控"的心理需求。**

开发者群体的一个深层心理特征：
- 不信任 SaaS 产品的"数据锁定"（Notion 导出麻烦、Linear 数据拿不出来）
- 喜欢能 `cat`、`grep`、`git diff` 的数据格式
- 对"我的数据在我电脑上"有强烈的安全感

Capture v2 的目录结构：`~/.capture/workspaces/<ws>/projects/<proj>/tasks/2026/06/TASK-00001.md`

这不是过度设计，而是**产品承诺的可视化** —— 用户随时可以用任何文本编辑器打开、修改、备份自己的全部数据。

### 对 v2 的启示
- **所有数据必须是人眼可读的**：禁止任何二进制或需要专用工具解析的格式
- **Workspace 目录结构本身就是文档**：新用户打开目录就能理解数据组织方式
- **Git 友好是长期留存的关键**：用户可以 `git init` 自己的工作空间，实现跨设备同步

### 设计体现
```bash
# 用户可以随时这么做：
cd ~/.capture/workspaces/default/projects/innate-capture
git init
git add .
git commit -m "backup my project data"
git push origin main

# 数据完全透明：
cat tasks/2026/06/TASK-00001.md
# 直接看到完整的 YAML frontmatter + Markdown body
```

---

## Insight 4: v1 → v2 的最大风险不是技术，是"习惯断裂"

### 表面认知
数据迁移是技术问题：把文件从 A 目录移动到 B 目录，更新数据库 schema。

### 深层洞察
**v1 → v2 的最大风险是：v1 用户升级后发现"熟悉的东西变了"，产生认知失调，进而放弃使用。**

具体场景：
- v1 用户习惯 `capture list` 看到所有任务
- v2 升级后，默认只显示当前 Project 的任务，用户以为数据丢了
- v1 用户习惯 `~/.capture/tasks/` 目录结构
- v2 升级后找不到原来的文件，恐慌

### 对 v2 的启示
- **默认行为必须与 v1 完全一致**：`capture list` 默认显示所有任务（跨 Project），而非仅当前 Project
- **文件路径变更必须可逆**：v1 的 `~/.capture/tasks/` 可以作为符号链接指向迁移后的位置
- **升级过程必须"无感知"**：自动迁移，不需要用户手动执行命令（除非出错）
- **保留逃生通道**：`CAPTURE_V1_MODE=1` 环境变量可以临时回到 v1 行为模式

### 设计体现
```bash
# 升级后，v1 用户的日常操作完全不变：
capture add "新想法"           # ✓ 和 v1 一样
capture list                  # ✓ 默认显示所有任务（包含 legacy）
capture kanban                # ✓ 看板显示所有任务

# 高级功能需要显式选择：
capture list --project my-app  # 新增：按项目筛选
capture project list           # 新增：管理项目
```

---

## Insight 5: 飞书和 GitHub 不是并列关系，而是"主从关系"

### 表面认知
Capture 同时支持飞书和 GitHub，两个都是重要的远端集成。

### 深层洞察
**对于目标用户（开发者），GitHub 是"工作系统"，飞书是"通知系统"。两者的战略权重完全不同。**

| 维度 | GitHub | 飞书 |
|------|--------|------|
| 用户停留时间 | 每天 2-4 小时 | 每天 30 分钟（查看消息） |
| 数据权威来源 | 是（Issue/PR 是真实工作） | 否（Bitable 是副本） |
| 用户操作深度 | 深（创建、编辑、关闭 Issue） | 浅（查看、简单回复） |
| 网络依赖 | 可接受（本身就是 Web 产品） | 必须在线 |
| 用户对集成的期待 | 高（希望终端能操作 GitHub） | 中（Bot 能记录就行） |

### 对 v2 的启示
- **GitHub 集成必须做到"专业级"**：字段映射准确、同步可靠、冲突处理好
- **飞书集成做到"够用就行"**：Bot 能按项目创建/查询任务即可，Bitable 同步保持现状
- **资源分配**：GitHub 模块投入 60% 的远端集成精力，飞书 40%
- **飞书的新增功能（Bot 支持 Project）是 P1，不是 P0**

### 设计体现
```
同步优先级：
  P0: GitHub Project 只读同步（MVP 必须有）
  P0: 飞书 Bitable 保持 v1 行为（不破坏现有功能）
  P1: GitHub 双向同步
  P1: 飞书 Bot 支持 Project 维度
  P2: 飞书 Bitable 多 Project 同步
```

---

## Insight 6: Spec 文档管理的真正痛点不是"写文档"，而是"需求到代码的追踪断裂"

### 表面认知
Spec 管理 = 提供 PRD 模板，让用户写需求文档。

### 深层洞察
**开发者的真正痛点是：写了 PRD 后，不知道哪些需求已经实现、哪些还在开发、哪些被放弃了。**

典型场景：
1. 产品经理写了 PRD，里面有 20 条 Acceptance Criteria
2. 开发过程中，有些 AC 被拆成了 GitHub Issues
3. 有些 AC 直接被代码实现了，没有关联 Issue
4. 2 个月后，没人知道 PRD 里的第 7 条 AC 到底做了没有

Capture v2 的 Spec 模块解决的是：**一条 AC 对应一个 Task，Task 的状态变化自动反映到 Spec 的完成度上。**

### 对 v2 的启示
- **Spec 不是文档工具，是追踪工具**：核心价值是"需求 → 任务 → 代码"的可追溯性
- **自动生成 Task 是关键功能**：`capture spec generate-tasks SPEC-001` 必须好用
- **Spec 视图需要显示完成度**：`4/10 AC completed`，而不是只显示文档内容
- **关联关系是双向的**：Task 页面能看到所属的 Spec，Spec 页面能看到所有关联 Task

### 设计体现
```markdown
---
id: SPEC-001
title: "用户认证模块"
status: in_progress          # 自动计算：如果所有 linked_tasks 都 done → implemented
progress: 3/5                # 自动计算
linked_tasks:
  - TASK-00001  # done
  - TASK-00002  # in_progress
  - TASK-00003  # todo
  - TASK-00004  # done
  - TASK-00005  # done
---

## Acceptance Criteria

- [x] 用户可以通过邮箱注册                    → TASK-00001 ✅
- [ ] 用户可以通过手机号注册                  → TASK-00002 🔄
- [ ] 注册后发送验证邮件                      → TASK-00003 ⬜
- [x] 用户可以通过邮箱+密码登录               → TASK-00004 ✅
- [x] 登录失败 3 次后锁定账户 15 分钟         → TASK-00005 ✅
```

---

## Insight 7: TUI 的"流畅感"比"功能多"更重要

### 表面认知
TUI 需要支持很多视图（Kanban、Stage、List、Spec、Sync Status...）。

### 深层洞察
**开发者对 TUI 的容忍度极低：一旦感觉到"卡"或"闪烁"，会立刻退出并再也不打开。**

原因：
- 开发者每天使用高度优化的终端工具（vim、tmux、htop）
- 他们对响应延迟极其敏感（200ms 就能感觉到）
- TUI 卡顿会被解读为"这个工具不成熟"

### 对 v2 的启示
- **性能优化不是"优化"，是"准入门槛"**：TUI 在 500 任务下必须流畅，否则功能再多也没用
- **先做一个流畅的 Kanban，再做五个卡顿的视图**：宁可少做视图，也要保证体验
- **虚拟滚动是必选项，不是可选项**：从第一天就要考虑，不能事后补救
- **SQLite 分页查询是 TUI 的生死线**：`ListTasks` 必须支持 LIMIT/OFFSET

### 设计体现
```go
// TUI 加载任务时强制分页
func (a *App) loadTasks() tea.Cmd {
    return func() tea.Msg {
        // 不是一次性加载所有任务！
        filter := model.TaskFilter{
            Limit:  100,  // 每页 100 条
            Offset: a.pageOffset,
        }
        tasks, err := a.taskSvc.List(context.Background(), filter)
        // ...
    }
}

// 虚拟滚动：只渲染可视区域
func (a *App) viewKanban() string {
    visibleTasks := a.getVisibleRange(a.allTasks)  // 只取当前窗口能显示的任务
    // 渲染 visibleTasks，而非 allTasks
}
```

---

## Insight 8: "当前 Project 自动推断"是降低认知负担的核武器

### 表面认知
用户创建 Task 时需要指定 Project，不然系统不知道放到哪里。

### 深层洞察
**如果用户每次创建 Task 都要显式选择 Project，Task 创建量会下降 50% 以上。**

人类行为的深层规律：
- 每次需要"做选择"时，都有认知成本
- 选择越多，放弃概率越高
- "默认正确"比"可以配置"重要 10 倍

### 对 v2 的启示
- **Project 推断必须智能且可覆盖**：有默认，但用户可以改
- **推断优先级**：
  1. `--project` flag（用户显式指定）
  2. 当前 Git 仓库的 remote origin → 匹配 Project 的 `github.repo`
  3. 环境变量 `CAPTURE_PROJECT`
  4. Workspace 的 `default_project`
  5. "legacy" project（兜底）
- **TUI 中当前 Project 要醒目显示**：让用户知道自己在哪个上下文

### 设计体现
```bash
# 场景：用户在 ~/workspace/innate-capture 目录下
# 该目录 git remote origin = git@github.com:variableway/innate-capture.git
# Project "innate-capture" 配置了 github.repo = "variableway/innate-capture"

capture add "修复登录bug"
# → 自动推断 Project = "innate-capture"
# → 不需要用户指定！

# 如果用户想在另一个项目下创建：
capture add "学习 Rust" --project side-projects
# → 显式覆盖推断结果
```

---

## Insight 9: 远端同步的"最终一致性"设计，比"强一致性"更实际

### 表面认知
同步系统应该保证本地和远端数据完全一致，任何差异都是 bug。

### 深层洞察
**GitHub Project、飞书 Bitable、本地 Markdown 三个系统是独立演化的，追求强一致性会引入无限复杂度和脆弱性。**

现实中的同步场景：
- GitHub 上有人关闭了 Issue，本地还没同步 → 不一致
- 本地修改了 Task 标题，但网络断了无法推送到 GitHub → 不一致
- 飞书 Bitable 被同事手动改了字段 → 与本地不一致
- 用户在 GitHub 网页上直接改了 Project 字段 → 与本地不一致

**同步的正确心智模型不是"镜像"，而是"缓存 + 定期刷新"。**

### 对 v2 的启示
- **接受不一致是常态，不是异常**：设计系统时假设三个数据源永远不完全一致
- **明确标注"同步状态"**：每个 Task 显示 `synced_at` 和 `sync_status`
- **用户手动触发同步，而非自动实时同步**：避免冲突和 API 限流
- **冲突解决策略默认"远端优先"**：因为 GitHub 是团队协作源，本地是个人工作区

### 设计体现
```yaml
---
id: TASK-00001
title: "修复登录bug"
sync:
  github:
    item_id: "PVTI_lADOABCD123"
    last_synced_at: 2026-06-04T10:00:00Z
    sync_status: synced          # synced / pending / conflict / error
  feishu:
    record_id: "recABC123"
    last_synced_at: 2026-06-04T09:00:00Z
    sync_status: synced
---
```

---

## Insight 10: Capture v2 的成功，取决于"第一次 GitHub 同步"的体验

### 表面认知
产品成功取决于功能完整度、文档质量、社区运营。

### 深层洞察
**Capture v2 的成功，80% 取决于用户第一次运行 `capture sync github` 时的体验。**

这是一个"Aha Moment"（顿悟时刻）：
- 如果用户运行后，30 秒内看到自己 GitHub Project 的所有任务出现在终端里 → "哇，这太棒了"
- 如果用户运行后，看到"token not found"、"project not found"、"GraphQL error" → "这玩意不行"

开发者对新工具的耐心窗口极短：
- 第一次使用：愿意给 5 分钟
- 如果第一次失败：大概率放弃
- 如果第一次成功：会成为忠实用户并推荐给同事

### 对 v2 的启示
- **token 获取必须零配置（对 gh CLI 用户）**：`gh auth status` 有有效 token → 直接可用
- **Project 配置必须自动发现**：用户输入 `github.com/owner/repo`，自动查询可用的 Project numbers
- **首次同步必须显示进度**："Fetching 127 items from GitHub... Done"，不能黑屏等待
- **错误信息必须可执行**：不是"GraphQL error: Node not found"，而是"找不到 GitHub Project，请确认仓库开启了 Project 功能：`capture project config --github-project 1`"

### 设计体现
```bash
$ capture project sync innate-capture
✓ GitHub token found via `gh auth status`
✓ Found GitHub Project #1: "Innate Capture Development"
⏳ Fetching items from GitHub (127 total)...
  ████████████████████████████████████████  100%  (127/127)
✓ Sync complete: 127 items fetched
  - 42 new tasks created locally
  - 85 existing tasks updated
  - 0 conflicts detected

$ capture list --project innate-capture
ID          TITLE                           STATUS    STAGE
TASK-G001   设计 Workspace API              todo      analysis
TASK-G002   实现 GitHub GraphQL 客户端      todo      inbox
...
```

---

## 洞察汇总矩阵

| # | 洞察 | 影响面 | 优先级 | 验证方式 |
|:---|:---|:---|:---:|:---|
| 1 | 竞品是"开发者记忆"，不是 Linear | 产品定位 | P0 | 用户访谈：为什么选择 Capture？ |
| 2 | GitHub 集成的价值是"离线缓存+统一视图" | 架构设计 | P0 | MVP 发布后用户留存率 |
| 3 | Workspace 物理目录 = 数据主权承诺 | 存储设计 | P0 | 用户是否 git init workspace |
| 4 | v1→v2 最大风险是习惯断裂 | 迁移设计 | P0 | v1 用户升级率 @ 2 周 |
| 5 | GitHub 是主，飞书是从 | 资源分配 | P0 | GitHub vs 飞书功能使用频率 |
| 6 | Spec 是追踪工具，不是文档工具 | 功能设计 | P1 | Spec 到 Task 的生成使用率 |
| 7 | TUI 流畅感 > 功能数量 | 体验设计 | P0 | TUI 使用时长和留存 |
| 8 | Project 自动推断降低 50% 放弃率 | CLI 设计 | P0 | Task 创建量对比 v1 |
| 9 | 同步是缓存模型，不是镜像模型 | 同步架构 | P1 | 同步失败率和用户投诉 |
| 10 | 第一次 GitHub 同步决定产品生死 | 上线策略 | P0 | 新用户 7 日留存率 |

---

## 这些洞察如何指导日常决策

当面临以下选择时，参考上述洞察：

**"要不要加一个配置项？"**
→ 参考 Insight 8：每增加一个配置项，就增加一分认知负担。先做好自动推断，配置项作为覆盖手段。

**"GitHub 同步先做只读还是直接上双向？"**
→ 参考 Insight 2：只读同步解决 80% 的痛点，且风险低。双向同步是锦上添花。

**"TUI 要不要加更多视图？"**
→ 参考 Insight 7：先确保 Kanban 在 500 任务下 60 FPS，再考虑加视图。

**"飞书 Bot 的新功能优先级怎么排？"**
→ 参考 Insight 5：飞书保持可用即可，不要投入过多。GitHub 集成才是护城河。

**"v1 用户升级后抱怨找不到原来的任务了"**
→ 参考 Insight 4：这是致命问题。`capture list` 默认必须显示所有任务，Project 筛选是可选的。
