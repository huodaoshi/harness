# Issue tracker：本地 Markdown

本仓库的 issue 与 PRD 以 markdown 文件形式存放在 `.scratch/`。**正文使用简体中文**；下文中的路径、`Status:` 角色名等约定保持英文 slug/标识符不变。

## 约定

- 每个功能一个目录：`.scratch/<feature-slug>/`
- PRD：`.scratch/<feature-slug>/PRD.md`
- 实现类 issue：`.scratch/<feature-slug>/issues/<NN>-<slug>.md`，编号从 `01` 起
- 分拣状态写在每个 issue 文件顶部附近的 `Status:` 行（角色字符串见 `triage-labels.md`）
- 评论与对话历史追加在文件末尾的 `## 评论` 标题下

## 当技能要求「发布到 issue tracker」

在 `.scratch/<feature-slug>/` 下创建新文件（目录不存在则创建）。文件内容（标题、正文、评论）使用简体中文。

## 当技能要求「获取相关工单」

读取用户给出的路径；用户通常会直接提供路径或 issue 编号。
