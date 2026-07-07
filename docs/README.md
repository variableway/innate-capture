# Capture 文档（当前实现）

本目录记录 `innate-capture` 当前已经实现并可用的能力。

## 当前可用 CLI

- `capture config`
  - `config show`：显示当前生效完整配置（文件不存在时显示默认值）
  - `config get <key>` / `config set <key> <value>`
- `capture doctor`
  - 校验 workspace 目录结构（`ideas/`、`daily/`）
- `capture idea`
  - `idea add <title> -d -c`：写入 `ideas/inbox/*.md`
  - `idea list`：列出 inbox 条目
- `capture daily`
  - 默认读取当天 `daily/YYYY-MM-DD.md`
  - `--date YYYY-MM-DD` 读取指定日期
  - `--reset` 从 `daily/_template/day.md` 重建目标日期文件
  - `--section input|output|ideas` 输出指定区块
  - `--open` 用 `defaults.editor` 打开目标日期文件

## 当前存储与约束

- idea：文件存储（`ideas/inbox/*.md`）
- daily：文件存储（`daily/YYYY-MM-DD.md`，保留历史，不覆盖旧日期）
- workspace root：
  - 配置项 `workspace.root`
  - 环境变量 `CAPTURE_WORKSPACE_ROOT` 可覆盖

## 运行与验证

- 构建：`go build .`
- 测试：`go test ./...`
- 本地最小验证：
  - `capture doctor --data-dir .`
  - `capture idea add "demo" --data-dir .`
  - `capture daily --reset --data-dir .`
