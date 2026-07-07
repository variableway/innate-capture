# Tasks: innate-capture

## 分工

| 文件 | 内容 |
|---|---|
| `issues/<id>.md` | 任务契约 |
| `tasks/issues/<id>.md` | 施工日志 |

## 推荐顺序

```text
cap-ws-config
    ├── cap-idea-write ──► cap-tui-idea
    └── cap-daily-read
            └── cap-docs-new
                    └── cap-docs-purge（最后，human）
```

## 四件套 Prompt

```text
项目：innate-capture
1. 读 docs/project/index.md
2. 读 docs/project/issues/<id>.md
3. 读 docs/project/tasks/issues/<id>.md（最后一条）
4. 本轮目标：（≤3 条子任务）

分支：feature/<id>-<agent>
合约：docs/project/spec/contracts/workspace-io-v1.md
```

## 索引

见 [`issues/README.md`](../issues/README.md)
