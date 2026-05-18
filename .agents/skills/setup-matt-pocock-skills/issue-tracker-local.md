# Issue 跟踪器：本地 Markdown

本仓库的 issue 与 PRD 以 markdown 文件存放在 `.scratch/`。

## 约定

- 每个特性一个目录：`.scratch/<feature-slug>/`
- PRD 为 `.scratch/<feature-slug>/PRD.md`
- 实现 issue 为 `.scratch/<feature-slug>/issues/<NN>-<slug>.md`，从 `01` 起编号
- 分拣状态记录在 issue 文件顶部附近的 `Status:` 行（角色字符串见 `triage-labels.md`）
- 评论与对话历史追加在文件末尾 `## Comments` 标题下

## 当技能说「发布到 issue 跟踪器」

在 `.scratch/<feature-slug>/` 下创建新文件（必要时创建目录）。

## 当技能说「获取相关 ticket」

读取所引路径的文件。用户通常会直接传入路径或 issue 编号。
