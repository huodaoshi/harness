# Issue 跟踪器：GitHub

本仓库的 issue 与 PRD 存放在 GitHub Issues。所有操作使用 `gh` CLI。

## 约定

- **创建 issue**：`gh issue create --title "..." --body "..."`。多行正文用 heredoc。
- **读取 issue**：`gh issue view <number> --comments`，用 `jq` 过滤评论，并获取标签。
- **列出 issue**：`gh issue list --state open --json number,title,body,labels,comments --jq '[.[] | {number, title, body, labels: [.labels[].name], comments: [.comments[].body]}]'`，配合适当的 `--label`、`--state` 过滤。
- **评论**：`gh issue comment <number> --body "..."`
- **添加/移除标签**：`gh issue edit <number> --add-label "..."` / `--remove-label "..."`
- **关闭**：`gh issue close <number> --comment "..."`

从 `git remote -v` 推断仓库——在 clone 内运行 `gh` 时会自动识别。

## 当技能说「发布到 issue 跟踪器」

创建 GitHub issue。

## 当技能说「获取相关 ticket」

运行 `gh issue view <number> --comments`。
