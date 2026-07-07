# Contract: workspace-io-v1

**Version:** v1.0  
**Date:** 2026-07-06  
**Parties:** capture CLI ↔ innate-works 目录  
**Issues:** cap-ws-config, cap-idea-write, cap-daily-read

## 配置（capture `config.yaml`）

```yaml
workspace:
  root: /absolute/path/to/innate-works
  ideas_inbox: ideas/inbox
  daily_today: daily/today.md
```

环境变量：`CAPTURE_WORKSPACE_ROOT` 覆盖 `workspace.root`。

提供方校验：`root` 下存在 `ideas/` 与 `daily/` 目录。

## Idea 写入（cap-idea-write）

**路径：** `{root}/{ideas_inbox}/YYYY-MM-DD-<slug>.md`

**正文模板：**

```markdown
# <Title>

**Captured:** YYYY-MM-DD
**Source:** capture-cli | capture-tui
**Stage:** inbox

## 一句话

<description>

## 原始上下文

<context or empty>

## 初步问题

- [ ]

## 下一步

- [ ] 留在 inbox 观察
- [ ] 晋升到 exploring/
```

**Slug 规则：** 小写；空格/下划线→`-`；仅 `[a-z0-9._-]`；冲突追加 `-2`, `-3`…

## Daily 读取（cap-daily-read）

**路径：** `{root}/{daily_today}`

- 不存在：从 `{root}/daily/_template/day.md` 复制，`__DATE__` → 今日
- 不写入 checkbox 状态（只读展示；`--open` 交用户编辑器）

## CLI 面（消费方期望）

```bash
capture idea add "<title>" [-d desc] [-c context]
capture idea list
capture daily [--open] [--section input|output|ideas]
```

## 变更记录

| 版本 | 日期 | 变更 |
|---|---|---|
| v1.0 | 2026-07-06 | 初稿，来自 exploring plan |
