# Issue: cap-daily-read

**Status:** implemented (pending PR)
**Mode:** horizontal  
**Sprint:** [cap-sprint-01-workspace-mvp](./cap-sprint-01-workspace-mvp.md)  
**Feature:** [daily-read](../features/daily-read.md)  
**Assignee:** —  
**Branch:** `feature/cap-daily-read-<agent>`  
**Contract:** [workspace-io-v1](../spec/contracts/workspace-io-v1.md)

## 背景

实现 `capture daily` 读取 innate-works 今日清单。

## 范围

**做：**

- `internal/daily`：Read、BootstrapFromTemplate、PrintSection
- `cmd/daily.go`：`daily`, `--open`, `--section`
- today 缺失时从 `daily/_template/day.md` 复制

**不做：**

- 修改 today 内 checkbox（MVP 只读；`--open` 除外）

## 可改动的路径

- `internal/daily/`
- `cmd/daily.go`

## 验收标准

- [x] `capture daily` 打印 today 全文
- [x] `--section ideas` 只打 Ideas 焦点节
- [x] 无 today 时 bootstrap 成功

## 依赖

- 阻塞：[cap-ws-config](./cap-ws-config.md)

## 过程日志

[`tasks/issues/cap-daily-read.md`](../tasks/issues/cap-daily-read.md)
