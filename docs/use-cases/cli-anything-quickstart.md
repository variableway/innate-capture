# CLI-Anything 快速上手教程

用 Claude Code + cli-anything-libreoffice 生成演示文稿。

---

## 1. 安装

```bash
pip install cli-anything-libreoffice
brew install libreoffice        # macOS
# sudo apt install libreoffice  # Ubuntu
cli-anything-libreoffice --help
```

## 2. 重要：单进程限制

cli-anything-libreoffice 的每次 CLI 调用是**独立进程**，修改只在内存中，不会自动持久化。

**错误用法（多次 CLI 调用）：**
```bash
# ✗ 每次调用都是新进程，add-slide 的修改会丢失
cli-anything-libreoffice --project demo.json impress add-slide -t "标题" -c "内容"
cli-anything-libreoffice --project demo.json document save  # 保存的是原始文件
```

**正确用法一：REPL 模式（单进程）**
```bash
cli-anything-libreoffice
◆ libreoffice ❯ document new --type impress -n "demo" -o demo.json
◆ libreoffice ["demo"] ❯ impress add-slide -t "标题" -c "内容"
◆ libreoffice ["demo"] ❯ impress add-slide -t "第二页" -c "内容"
◆ libreoffice ["demo"] ❯ document save
◆ libreoffice ["demo"] ❯ exit
```

**正确用法二：Python API（推荐，Claude Code 可直接调用）**
```python
from cli_anything.libreoffice.core.session import Session
from cli_anything.libreoffice.core import document as doc_mod
from cli_anything.libreoffice.core import impress as impress_mod
from cli_anything.libreoffice.core import export as export_mod

sess = Session()
proj = doc_mod.create_document(doc_type='impress', name='我的PPT')
sess.set_project(proj, 'output.json')

impress_mod.add_slide(sess.get_project(), title='标题', content='内容')
impress_mod.add_slide(sess.get_project(), title='第二页', content='更多内容')

sess.save_session()
export_mod.export(sess.get_project(), 'output.pdf', preset='pdf', overwrite=True)
```

## 3. Claude Code 自动化

安装 Skill 后，对 Claude Code 说：

> **"帮我做一份项目架构 PPT，导出 PDF"**

Agent 会生成 Python 脚本调用 cli-anything-libreoffice API，自动完成：
1. 创建 Impress 项目（`doc_mod.create_document`）
2. 添加幻灯片（`impress_mod.add_slide`）
3. 保存项目（`sess.save_session`）
4. 导出 PDF（`export_mod.export`）

## 4. 命令速查

### 文档操作
| 命令 | 说明 |
|------|------|
| `document new --type impress -n "名" -o file.json` | 创建演示文稿（自动保存） |
| `document save` | 保存当前项目（仅 REPL 内有效） |
| `document info` | 查看文档信息 |
| `document json` | 输出原始 JSON |

### 幻灯片操作（REPL 或 Python API）
| 命令 | 说明 |
|------|------|
| `impress add-slide -t "标题" -c "内容"` | 添加幻灯片 |
| `impress set-content <index> -t "标题" -c "内容"` | 修改幻灯片 |
| `impress list-slides` | 列出所有幻灯片 |
| `impress remove-slide <index>` | 删除幻灯片 |
| `impress add-element <slide> --type text_box -t "文字"` | 添加文本框 |

### 导出
| 命令 | 说明 |
|------|------|
| `export render out.pdf --preset pdf` | 导出 PDF |
| `export render out.pptx --preset pptx` | 导出 PPTX |
| `export render out.odp --preset odp` | 导出 ODP |
| `export presets` | 列出所有导出格式 |

### 可用导出预设
`odt` `ods` `odp` `html` `text` `pdf` `docx` `xlsx` `pptx` `csv`

## 5. 完整 Python 脚本示例

```python
#!/usr/bin/env python3
"""用 cli-anything-libreoffice 生成项目架构 PPT"""
from cli_anything.libreoffice.core.session import Session
from cli_anything.libreoffice.core import document as doc_mod
from cli_anything.libreoffice.core import impress as impress_mod
from cli_anything.libreoffice.core import export as export_mod

sess = Session()
proj = doc_mod.create_document(
    doc_type='impress',
    name='项目架构',
    profile='presentation_16_9'
)
sess.set_project(proj, 'architecture.json')

slides = [
    ("系统架构设计", "团队沟通 · 方案汇报"),
    ("输入层", "• Terminal：命令行输入\n• IM：飞书/微信消息同步"),
    ("处理层", "• 分析 Agent：分类·关联·提炼\n• 执行 Agent：Claude Code"),
    ("存储与同步", "• 本地：Markdown + SQLite\n• 云端：飞书多维表格"),
]

for title, content in slides:
    impress_mod.add_slide(sess.get_project(), title=title, content=content)

sess.save_session()
export_mod.export(sess.get_project(), 'architecture.pdf', preset='pdf', overwrite=True)
print("✓ 已生成 architecture.pdf")
```

## 6. 常见问题

**Q: `--project` 模式下 add-slide 不生效？**
A: 这是已知限制。每次 CLI 调用是独立进程，内存修改不会持久化。使用 REPL 模式或 Python API。

**Q: 导出失败？**
A: 确认 LibreOffice 已安装：`soffice --version`

**Q: 如何使用 Claude Code 自动化？**
A: 安装 Skill：`npx skills add HKUDS/CLI-Anything --skill cli-anything-libreoffice -g -y`

---

参考：[CLI-Anything 官方仓库](https://github.com/HKUDS/CLI-Anything)
