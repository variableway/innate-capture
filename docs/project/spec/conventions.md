# Conventions: innate-capture

## 目标

本文件定义 `innate-capture` 的代码组织、`cmd` 编写方式和多 Agent 协作约定。默认目标是：

- 先满足 `docs/project/issues/` 的 issue 边界，不扩 scope
- 业务域代码清晰分层，可替换实现、可测试
- `cmd` 只做编排，不承载业务逻辑

## 代码组织

- Go 模块根目录：项目根（`go build .` 在根目录执行）
- 业务 domain：优先放 `internal/<domain>/`
- 公共能力（非当前 sprint 业务域，或可复用 SDK/适配器）：放 `pkg/<name>/`
  - 当前约定：飞书相关迁移到 `pkg/feishu`、`pkg/bitable`、`pkg/bot`
- 公共基础设施适配层：放 `internal/platform/<capability>/`
  - 当前约定：共享文件系统仓储接口放在 `internal/platform/fsrepo`
- 新包需有单测，优先 table-driven；路径解析需覆盖 edge case

## Domain 包结构约定

以 `internal/workspace` 为例，一个 domain 推荐结构：

- `model.go`：domain 输入/输出模型
- `service.go`：对外接口（API 契约）
- `<domain>_service.go`：`service` 的默认实现
- `repo/`：数据访问或外部依赖适配（文件系统、DB、HTTP SDK）
- 可选：
  - `entity/`：仅在存在数据库实体映射时引入
  - `tui.go` 或 `tui/`：仅在该 domain 需要 TUI 交互时引入

约束：

- `model.go`、`service.go` 放包根目录
- 旧函数接口可保留为 facade（兼容层），内部委托默认 service
- 当类或方法规模明显增长时再拆子目录，避免过度设计

## Cmd 编写约定

`cmd` 层定位：**CLI 编排层**，不写业务细节。

### 1) 命令注册

- 每个命令文件定义一个主命令变量：`xxxCmd`
- 在 `init()` 中完成 flags 绑定和子命令挂载
- 根命令是否注册必须显式，和当前 sprint 范围保持一致

### 2) RunE 职责边界

`RunE` 只做这些事：

- 读取参数/flags
- 加载配置与依赖（store/service）
- 调用 domain service
- 统一输出与错误返回

`RunE` 不做这些事：

- 不直接写文件系统/数据库细节（下沉到 repo/service）
- 不内联复杂业务分支（下沉到 domain）

### 3) 错误与输出

- 一律返回 `error`，不在 `RunE` 里 `os.Exit`
- 错误信息可读、可定位（包含关键路径或参数）
- 输出保持稳定，便于后续自动化与测试

### 4) 配置与环境变量

- 统一通过 `config.Load(...)` + `config.Config`
- 环境变量覆盖规则在 domain/service 内统一处理（如 workspace root）
- 不在多个命令重复实现同一配置解析逻辑

## 测试约定

- 新增 domain 至少覆盖：
  - 正常路径
  - 边界条件
  - 错误路径（缺文件、非法参数、空配置等）
- 对 repo 依赖可抽象接口，优先通过注入提升可测性

## Git 与多 Agent

- 分支：`feature/<issue-id>-<agent>`
- 不并行改同一 issue
- PR 合并后更新 `features/` 状态与 task log

## Agent 约定

1. 入口：[`../index.md`](../index.md)
2. 合约：[`contracts/workspace-io-v1.md`](./contracts/workspace-io-v1.md)
3. 完成后追加 [`../tasks/issues/`](../tasks/issues/) 对应日志

## 与 innate-works 对齐

- inbox 格式对齐 `innate-works/ideas/_template/inbox-entry.md`
- daily 不修改 checkbox 语义，只读或 `--open` 交编辑器

