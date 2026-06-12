# 用例：CLI-Anything LibreOffice —— 自动生成演示文稿

本文档展示如何使用 CLI-Anything 的 LibreOffice 模块，将系统架构信息自动转换为演示文稿（PPT），并以 Capture 项目自身的架构图作为完整示例。

---

## 背景

Capture 项目的架构信息分散在以下位置：

- `AGENTS.md` —— 项目整体架构与模块说明
- `docs/features/input_ports_architecture.svg` —— 输入端口架构图（SVG）
- `docs/README.md` —— 文档目录结构说明

目标：将这些架构信息整理成一份结构化的演示文稿（PDF），用于团队沟通或方案汇报。

---

## 前置准备

### 1. 安装 CLI-Anything LibreOffice Harness

```bash
pip install cli-anything-libreoffice
```

### 2. 安装 LibreOffice（导出 PDF 必需）

**macOS:**
```bash
brew install libreoffice
```

**Ubuntu/Debian:**
```bash
sudo apt install libreoffice
```

**验证安装:**
```bash
cli-anything-libreoffice --help
```

---

## Skill 配置（让 AI Agent 自动调用）

### Claude Code

```bash
# 方式一：直接安装 CLI-Anything 官方 skill
npx skills add HKUDS/CLI-Anything --skill cli-anything-libreoffice -g -y

# 方式二：手动复制 skill 文件
mkdir -p ~/.claude/skills/cli-anything-libreoffice
cp $(python -c "import cli_anything.libreoffice; print(cli_anything.libreoffice.__path__[0])")/skills/SKILL.md \
   ~/.claude/skills/cli-anything-libreoffice/
```

### Codex

```bash
mkdir -p ~/.codex/skills/cli-anything-libreoffice
cp $(python -c "import cli_anything.libreoffice; print(cli_anything.libreoffice.__path__[0])")/skills/SKILL.md \
   ~/.codex/skills/cli-anything-libreoffice/
```

安装后重启 Agent，即可在对话中通过自然语言触发文档生成。

---

## 完整案例：Capture 架构图 → 演示文稿

### 架构信息梳理

Capture 项目的核心架构包含以下层次：

| 层级 | 组件 | 说明 |
|------|------|------|
| **输入层** | Terminal | markdown / text 格式输入 |
| | IM | 飞书/微信消息同步到本地或云端 |
| **统一入口** | Inbox | 所有任务和灵感的统一待处理队列 |
| **并行处理层** | 分析 Agent | 分类 · 关联 · 提炼 |
| | 调研 Agent | 联网 · 搜索 · 验证 |
| | 执行 Agent | Claude Code 实际编码 |
| **存储同步层** | 本地存储 | Markdown + Frontmatter + SQLite |
| | 飞书同步 | 多维表格双向同步 |

架构图原件：`docs/features/input_ports_architecture.svg`

---

### 步骤化命令执行

#### 步骤 1：创建演示文稿项目

```bash
cli-anything-libreoffice project new \
  -o capture-architecture.json \
  --type impress \
  --name "Capture 系统架构"
```

#### 步骤 2：封面页

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "Capture 系统架构设计" \
  --layout title
```

#### 步骤 3：输入层

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "输入层：多端采集" \
  --layout bullets

cli-anything-libreoffice --project capture-architecture.json \
  impress set-content \
  --slide 1 \
  --bullets '[
    "Terminal 输入：markdown / text 格式，适合开发者快速记录",
    "IM 输入：飞书/微信等即时通讯工具，消息自动同步到本地或云端",
    "统一格式转换：所有输入最终标准化为结构化任务数据"
  ]'
```

#### 步骤 4：Inbox 统一入口

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "Inbox：统一任务入口" \
  --layout title-content

cli-anything-libreoffice --project capture-architecture.json \
  impress set-content \
  --slide 2 \
  --content "所有灵感、任务、想法首先进入 Inbox，作为统一待处理队列，等待后续分析和分发。"
```

#### 步骤 5：并行处理层（三大 Agent）

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "并行处理层：三大 Agent" \
  --layout bullets

cli-anything-libreoffice --project capture-architecture.json \
  impress set-content \
  --slide 3 \
  --bullets '[
    "分析 Agent：分类 · 关联 · 提炼 —— 理解任务本质",
    "调研 Agent：联网 · 搜索 · 验证 —— 补充必要上下文",
    "执行 Agent：Claude Code —— 实际编码和工具调用"
  ]'
```

#### 步骤 6：架构全景图（嵌入 SVG）

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "架构全景图" \
  --layout title-only

cli-anything-libreoffice --project capture-architecture.json \
  impress add-image \
  --slide 4 \
  --image docs/features/input_ports_architecture.svg \
  --position center \
  --width 80%
```

#### 步骤 7：存储与同步层

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "存储与同步层" \
  --layout bullets

cli-anything-libreoffice --project capture-architecture.json \
  impress set-content \
  --slide 5 \
  --bullets '[
    "本地存储：Markdown + Frontmatter + SQLite 双写",
    "飞书同步：多维表格双向同步",
    "任务状态流：todo → in_progress → done / cancelled → archived"
  ]'
```

#### 步骤 8：设计要点总结

```bash
cli-anything-libreoffice --project capture-architecture.json \
  impress add-slide \
  --title "设计要点" \
  --layout bullets

cli-anything-libreoffice --project capture-architecture.json \
  impress set-content \
  --slide 6 \
  --bullets '[
    "纯 Go 实现，Cobra CLI + Bubble Tea TUI",
    "飞书 Bot 双模式：Webhook + WebSocket",
    "Agent 原生：支持 OpenCLI / CLI-Anything 扩展",
    "任务驱动：从灵感捕获到执行闭环"
  ]'
```

#### 步骤 9：导出 PDF

```bash
cli-anything-libreoffice --project capture-architecture.json \
  export render \
  capture-architecture.pdf \
  --preset pdf \
  --overwrite
```

---

## REPL 交互式方式

适合逐步调试和实时预览：

```bash
$ cli-anything-libreoffice

libreoffice> project new -o capture.json --type impress --name "Capture架构"
✓ Created Impress document: Capture架构

libreoffice[Capture架构]> impress add-slide --title "Capture 系统架构设计" --layout title
✓ Added slide 0: title layout

libreoffice[Capture架构]*> impress add-slide --title "输入层：多端采集" --layout bullets
✓ Added slide 1: bullets layout

libreoffice[Capture架构]*> impress set-content --slide 1 \
  --bullets '["Terminal: markdown/text", "IM: 飞书/微信同步", "统一格式转换"]'
✓ Updated slide 1

libreoffice[Capture架构]*> impress add-slide --title "并行处理层" --layout bullets
✓ Added slide 2: bullets layout

libreoffice[Capture架构]*> impress set-content --slide 2 \
  --bullets '["分析 Agent: 分类·关联·提炼", "调研 Agent: 联网·搜索·验证", "执行 Agent: Claude Code"]'
✓ Updated slide 2

libreoffice[Capture架构]*> impress add-slide --title "架构全景图" --layout title-only
✓ Added slide 3: title-only layout

libreoffice[Capture架构]*> impress add-image --slide 3 \
  --image docs/features/input_ports_architecture.svg --position center
✓ Added image to slide 3

libreoffice[Capture架构]*> export render capture-architecture.pdf --preset pdf --overwrite
✓ Rendered 4 slides to capture-architecture.pdf
```

---

## Agent 自动化工作流

配置好 Skill 后，只需对 Agent 说一句话：

> **"帮我把 capture 项目的架构图做成一份 PPT，包含输入层、Inbox、并行处理三大 Agent、存储同步这几页，最后导出 PDF"**

Agent 基于 Skill 知识自动执行全部命令，最终返回：

```
✓ 已生成 capture-architecture.pdf
  - P1 封面：Capture 系统架构设计
  - P2 输入层：Terminal + IM 多端采集
  - P3 Inbox：统一任务入口
  - P4 并行处理层：分析 / 调研 / 执行
  - P5 架构全景图（含 SVG 架构图）
  - P6 存储与同步
  - P7 设计要点总结
```

---

## 命令速查表

| 命令 | 说明 |
|------|------|
| `project new -o <file> --type impress` | 创建新的演示文稿项目 |
| `impress add-slide --title "..." --layout <type>` | 添加幻灯片 |
| `impress set-content --slide N --bullets '[...]'` | 设置列表内容 |
| `impress set-content --slide N --content "..."` | 设置正文内容 |
| `impress add-image --slide N --image <path>` | 插入图片 |
| `impress list-slides` | 列出所有幻灯片 |
| `export render <output> --preset pdf` | 导出为 PDF |
| `export render <output> --preset pptx` | 导出为 PPTX |
| `session undo` | 撤销上一步操作 |
| `session redo` | 重做 |

### 可用布局类型

- `title` —— 仅标题
- `title-content` —— 标题 + 正文
- `bullets` —— 标题 + 列表
- `two-column` —— 双栏
- `title-only` —— 仅标题（适合放全屏图片）

---

## 常见问题

| 问题 | 解决方案 |
|------|---------|
| SVG 图片无法插入？ | LibreOffice Impress 原生支持 SVG，确认文件路径正确 |
| 需要调整图片大小？ | 添加 `--width 80%` 或 `--height 60%` |
| 查看更多布局？ | 执行 `impress list-layouts` |
| 导出为 PPTX？ | `--preset pptx` |
| 命令执行失败？ | 先执行 `opencli doctor` 或检查 LibreOffice 是否安装 |

---

## 参考

- [CLI-Anything 官方仓库](https://github.com/HKUDS/CLI-Anything)
- [CLI-Anything LibreOffice Skill](https://github.com/HKUDS/CLI-Anything/blob/main/skills/cli-anything-libreoffice/SKILL.md)
- [HARNESS.md 方法论](https://github.com/HKUDS/CLI-Anything/blob/main/cli-anything-plugin/HARNESS.md)
