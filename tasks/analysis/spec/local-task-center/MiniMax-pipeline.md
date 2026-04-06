# Pipeline Module Spec (MiniMax 设计)

## 目标

维护任务中心的工作流状态与阶段流转。

## 阶段定义 (Stage)

1. `inbox`: 刚刚收集进来的 raw idea
2. `mindstorm`: 正在进行发散性思考
3. `analysis`: 正在进行可行性与结构分析
4. `planning`: 正在排期与出架构
5. `prd`: 正在出详细需求
6. `tasks`: 正在做任务拆解
7. `dispatch`: 已分配给具体的 agent/repo
8. `execution`: agent 正在执行
9. `review`: 执行完毕，等待人工验收

## 状态定义 (Status)

维持现有：
- `todo`: 未开始
- `in_progress`: 进行中
- `done`: 已完成
- `cancelled`: 已取消
- `archived`: 已归档

## 状态与阶段的正交性

- Stage 代表任务在"研发管线"里的位置
- Status 代表在当前 Stage 里的工作状态
