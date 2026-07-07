# Feature: daily-read

**Status:** planned  
**Issues:** [cap-daily-read](../issues/cap-daily-read.md)

## 说明

`capture daily` 展示 innate-works `daily/today.md`；缺失时从 `_template/day.md` 引导创建。

## 相关代码

- `internal/daily/`
- `cmd/daily.go`

## 验收

- [ ] 终端输出完整 today 或按 `--section` 切片
- [ ] `--open` 使用 `defaults.editor`
- [ ] 不在 capture 内维护第二份 daily 状态
