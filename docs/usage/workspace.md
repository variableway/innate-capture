# Capture × innate-works

用户向 CLI 手册。开发文档见 [`project/`](../project/index.md)。

## 配置

编辑 `~/.capture/config.yaml`（或 `capture init` 后）：

```yaml
workspace:
  root: D:/innate-works/innate-works   # innate-works 根目录
  ideas_inbox: ideas/inbox
  daily_today: daily/today.md
```

或环境变量：`CAPTURE_WORKSPACE_ROOT`

## 提交 idea

```bash
capture idea add "我的新想法"
capture idea add "标题" -d "一句话描述"
capture idea add "标题" -d "描述" -c "触发场景、链接"
capture idea list
```

写入：`{workspace.root}/ideas/inbox/YYYY-MM-DD-<slug>.md`

## 查看每日清单

```bash
capture daily
capture daily --section output
capture daily --open    # 用编辑器打开 today.md
```

读取：`{workspace.root}/daily/today.md`

## 飞书

本阶段未启用。代码保留，配置见历史文档（purge 前）或 P4。

## 开发

Agent 入口：[`docs/project/index.md`](../project/index.md)
