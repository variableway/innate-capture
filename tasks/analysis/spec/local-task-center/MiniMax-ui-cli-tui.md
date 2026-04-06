# CLI/TUI UI Spec (MiniMax 设计)

## 角色定位

P0 级别主控台，个人本地优先工作流的核心承载界面。

## CLI 模块要求

### 1. Capture Entry
- 必须支持快速单行录入 `capture add "idea"`
- 支持带参数指定优先级、标签

### 2. Workflow Control
- `capture stage <id> <stage>` 推进阶段
- `capture status <id> <status>` 改变状态

### 3. Dispatch Control
- `capture assign <id> --agent codex --repo path/to/repo` 分派执行环境

### 4. Info Inspection
- `capture show <id>` 展示完整元数据，包含阶段与执行分派信息

## TUI 模块要求

### 1. Workflow Board
- 扩展现有的 status 三列看板
- 新增 stage 视图，按 `inbox / mindstorm / analysis / planning / prd / tasks / dispatch / execution / review` 展示任务

### 2. Task Card
- 必须展示：ID、标题、优先级、标签
- 新增摘要展示：当前阶段 (stage) 和 分派的 Agent (agent)

### 3. Detail View
- 支持展开展示执行状态和历史

### 4. Execution Monitor (Future)
- TUI 内部支持展示关联终端会话日志的 tail 摘要
