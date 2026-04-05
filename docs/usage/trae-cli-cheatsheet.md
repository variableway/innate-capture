# Trae 命令速查表

## 概述

本文基于本机 `Trae 3.5.43` 的 `trae -h` 输出整理，按日常开发场景归类常用命令，方便快速查询。

`trae` 的核心用途主要有三类：

- 打开项目、文件、diff、merge
- 管理扩展与 MCP
- 排错诊断、远程访问与命令行进入聊天

> 说明：当前 `trae -h` 顶层帮助中未看到公开的 `solo` 模式切换参数，CLI 更偏向“打开编辑器与执行辅助功能”，而不是完整映射界面中的所有模式开关。

## 1. 打开项目或文件

### 打开当前目录

```bash
trae .
```

适合进入当前项目目录开始工作。

### 打开指定目录

```bash
trae /path/to/project
```

### 打开单个文件

```bash
trae README.md
trae /path/to/file.go
```

### 同时打开多个路径

```bash
trae backend/ docs/ README.md
```

## 2. 精确定位到代码位置

### 跳到指定文件、行、列

```bash
trae -g main.go:120
trae -g internal/service/task_service.go:88:12
```

适合：

- 从终端报错直接跳到出错位置
- 从 grep、测试输出、编译日志里快速打开源码

## 3. 控制窗口行为

### 强制新开窗口

```bash
trae -n .
```

适合把不同项目隔离到不同窗口。

### 强制复用已有窗口

```bash
trae -r .
```

适合不想把桌面开很多窗口时使用。

### 添加目录到当前工作区

```bash
trae -a ../shared-lib
```

适合把多个仓库放到同一个工作区里联合查看。

### 从当前工作区移除目录

```bash
trae --remove ../shared-lib
```

## 4. 对比和合并文件

### 对比两个文件

```bash
trae -d old.txt new.txt
trae --diff before.go after.go
```

适合：

- 比较两个版本的配置文件
- 查看重构前后差异
- 作为图形化 diff 工具使用

### 三方合并文件

```bash
trae -m ours.go theirs.go base.go result.go
```

参数含义：

- `path1`：版本 1
- `path2`：版本 2
- `base`：共同基线版本
- `result`：合并结果输出文件

适合处理手动合并冲突。

## 5. 等待编辑完成再返回终端

### 配合脚本或 Git 流程使用

```bash
trae -w COMMIT_EDITMSG
trae --wait notes.md
```

作用：

- 命令会阻塞
- 直到对应文件或窗口关闭后才返回 shell

适合：

- 让 `Trae` 充当 Git 提交信息编辑器
- 在脚本中等待用户完成编辑

## 6. 扩展管理

### 查看已安装扩展

```bash
trae --list-extensions
trae --list-extensions --show-versions
```

### 按分类过滤扩展

```bash
trae --list-extensions --category themes
```

### 安装扩展

```bash
trae --install-extension ms-python.python
trae --install-extension golang.go
```

安装指定版本：

```bash
trae --install-extension vscode.csharp@1.2.3
```

安装本地 VSIX：

```bash
trae --install-extension ./my-extension.vsix
```

安装预发布版：

```bash
trae --install-extension ms-python.python --pre-release
```

### 卸载扩展

```bash
trae --uninstall-extension ms-python.python
```

### 更新全部扩展

```bash
trae --update-extensions
```

### 指定扩展目录

```bash
trae --extensions-dir ~/.trae-extensions
```

适合测试隔离环境，避免污染默认扩展目录。

## 7. Profile 与数据隔离

### 使用指定 Profile 打开项目

```bash
trae --profile work .
trae --profile clean-test .
```

适合：

- 工作项目和个人项目分离
- 为特定仓库准备独立配置
- 做演示或测试时使用临时 profile

### 指定用户数据目录

```bash
trae --user-data-dir /tmp/trae-user-data .
```

适合：

- 启动一个干净实例
- 排查是不是现有配置导致的问题
- 同时运行多个相互隔离的 `Trae` 实例

### 临时模式启动

```bash
trae --transient .
```

适合快速试验，不想复用原有缓存、状态和扩展数据时使用。

## 8. 调试界面问题或排错

### 查看版本

```bash
trae --version
```

### 查看状态与诊断信息

```bash
trae --status
```

适合遇到卡顿、扩展异常、进程占用异常时先查看总体状态。

### 输出更详细日志

```bash
trae --verbose
trae --log debug
trae --log trace
```

### 禁用所有扩展启动

```bash
trae --disable-extensions
```

适合确认问题是否由扩展引起。

### 禁用单个扩展

```bash
trae --disable-extension ms-python.python
```

### 关闭 GPU 加速

```bash
trae --disable-gpu
```

适合界面黑屏、闪烁、窗口渲染异常时尝试。

### 开关配置同步

```bash
trae --sync on
trae --sync off
```

## 9. 扩展开发与高级诊断

以下命令偏向扩展开发或底层排错，日常开发较少直接使用：

```bash
trae --enable-proposed-api publisher.extension
trae --inspect-extensions 9333
trae --inspect-brk-extensions 9333
trae --prof-startup
trae --cpuprofile <pid> <port> <duration> <type> <remote-debugging-port> <watch> <threshold> <interval>
trae --heapsnapshot <pid> <port> <type> <remote-debugging-port> <watch> <threshold> <interval>
```

用途概览：

- `--enable-proposed-api`：给扩展开启实验性 API
- `--inspect-extensions`：调试扩展宿主进程
- `--inspect-brk-extensions`：扩展宿主启动后先暂停，等调试器连接
- `--prof-startup`：分析启动性能
- `--cpuprofile`：抓 CPU profile
- `--heapsnapshot`：抓内存快照

## 10. Shell 集成与辅助信息

### 查询 shell 集成脚本路径

```bash
trae --locate-shell-integration-path zsh
trae --locate-shell-integration-path bash
```

可用取值：

- `bash`
- `zsh`
- `pwsh`
- `fish`

### 查看 telemetry 事件

```bash
trae --telemetry
```

主要用于开发和排查，不是日常开发高频命令。

## 11. MCP 管理

### 添加 MCP Server

```bash
trae --add-mcp '{"name":"server-name","command":"my-mcp-server"}'
```

作用：

- 向当前用户配置中注册一个 MCP Server 定义
- 供 Trae 内的 AI/Agent 后续调用外部工具或服务

适合：

- 接入自定义工具链
- 接入数据库、文档系统、浏览器自动化或内部平台能力

## 12. 子命令入口

### 命令行发起聊天

```bash
trae chat
```

从帮助文本看，`chat` 子命令用于“在当前工作目录中发起一个 chat session 并传入 prompt”。

如果你想看它的详细参数，可以继续执行：

```bash
trae chat -h
```

### 在浏览器中提供编辑器 UI

```bash
trae serve-web
```

作用：

- 启动一个服务
- 让编辑器界面可通过浏览器访问

### 通过安全隧道暴露当前机器

```bash
trae tunnel
```

作用：

- 让当前机器可被其他机器或 `vscode.dev` 访问
- 适合远程开发、临时协作或跨设备访问

## 13. 最常用命令清单

如果只记最常用的一组，优先记这些：

```bash
trae .
trae -g main.go:120
trae -d old.txt new.txt
trae -n .
trae -r .
trae --list-extensions
trae --install-extension golang.go
trae --disable-extensions
trae --status
trae chat
```

## 14. 推荐使用套路

### 场景 1：打开当前仓库开始写代码

```bash
trae .
```

### 场景 2：从报错日志跳到具体代码

```bash
trae -g internal/store/sqlite.go:42
```

### 场景 3：图形化比较两个文件

```bash
trae -d before.yaml after.yaml
```

### 场景 4：排查是不是扩展导致卡顿

```bash
trae --disable-extensions
```

### 场景 5：用干净环境复现问题

```bash
trae --transient .
```

### 场景 6：安装常用语言扩展

```bash
trae --install-extension golang.go
trae --install-extension ms-python.python
```

## 15. 补充说明

- `trae` 顶层帮助更偏“编辑器控制台入口”，不是所有 UI 能力都会暴露成 CLI 参数。
- 如果你想继续深入，下一步最值得看的通常是 `trae chat -h`，因为它更接近 AI 工作流。
- 如果你已经配置了 shell alias，也可以把常用命令再封装成更短的别名。
