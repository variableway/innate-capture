# Local Task Center 模块 Spec 索引

## 模块列表

- [MiniMax-intake.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-intake.md)
- [MiniMax-pipeline.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-pipeline.md)
- [MiniMax-assignment.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-assignment.md)
- [MiniMax-execution.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-execution.md)
- [MiniMax-dashboard.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-dashboard.md)
- [MiniMax-storage.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-storage.md)
- [MiniMax-ui-overview.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-ui-overview.md)
- [MiniMax-ui-cli-tui.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-ui-cli-tui.md)
- [MiniMax-ui-web.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-ui-web.md)
- [MiniMax-ui-desktop.md](file:///Users/patrick/workspace/variableway/innate/capture/docs/spec/local-task-center/MiniMax-ui-desktop.md)

## 模块关系

1. intake 接收输入并创建任务
2. pipeline 维护流程阶段与状态
3. assignment 负责下发到 agent 和 repo
4. execution 跟踪执行过程
5. storage 负责持久化与索引
6. dashboard 消费前面模块的数据做展示
7. ui-overview 统一三类界面的职责与阶段
8. ui-cli-tui 负责主控台能力
9. ui-web 负责 dashboard 与聚合观察
10. ui-desktop 负责本地增强壳层
