# innate-capture

`innate-capture` 是 `innate-works` 的终端前端，用来管理 workspace 下的 idea inbox 和 daily 文件。

## 当前已实现能力

- `capture config`
  - `config show`：显示完整生效配置（无配置文件时显示默认值）
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
  - `--open` 使用 `defaults.editor` 打开目标日期文件

## 安装与运行

在项目根目录：

```bash
go build .
./capture --help
```

如果你用 Task：

```bash
task build
task install
```

## 配置说明

默认会从 `~/.capture/config.yaml` 读取引导配置。

当前推荐只使用一个工作目录概念：`workspace.root`。

- 规范键：`workspace.root`
- 兼容键：`data_dir` / `app.data_dir`（内部会映射到 `workspace.root`）

示例：

```bash
capture config set workspace.root .
capture config show
```

Windows/Nushell 路径建议用单引号或正斜杠：

```bash
capture config set workspace.root 'D:\innate-works\innate-works\innate-capture'
capture config set workspace.root D:/innate-works/innate-works/innate-capture
```

## 快速验证

```bash
capture doctor --data-dir .
capture idea add "demo idea" -d "one line" --data-dir .
capture idea list --data-dir .
capture daily --reset --data-dir .
capture daily --date 2026-07-08 --section ideas --data-dir .
```

## 文档索引

- 项目规范与 issue：`docs/project/index.md`
- 合约：`docs/project/spec/contracts/workspace-io-v1.md`
- 约定：`docs/project/spec/conventions.md`
