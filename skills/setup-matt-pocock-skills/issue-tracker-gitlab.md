# Issue 跟踪器：GitLab

本仓库的 issue 与 PRD 存放在 GitLab Issues。所有操作使用 [`glab`](https://gitlab.com/gitlab-org/cli) CLI。

## 约定

- **创建 issue**：`glab issue create --title "..." --description "..."`。多行描述用 heredoc。`--description -` 可打开编辑器。
- **读取 issue**：`glab issue view <number> --comments`。机器可读输出用 `-F json`。
- **列出 issue**：`glab issue list -F json`，配合适当的 `--label` 过滤。
- **评论**：`glab issue note <number> --message "..."`。GitLab 称评论为 note。
- **添加/移除标签**：`glab issue update <number> --label "..."` / `--unlabel "..."`。多标签可用逗号分隔或重复该 flag。
- **关闭**：`glab issue close <number>`。`glab issue close` 不接受关闭说明，须先用 `glab issue note <number> --message "..."` 发布说明，再关闭。
- **合并请求**：GitLab 称 PR 为 merge request。使用 `glab mr create`、`glab mr view`、`glab mr note` 等——与 `gh pr ...` 形态相同，以 `mr` 代 `pr`，以 `note`/`--message` 代 `comment`/`--body`。

从 `git remote -v` 推断仓库——在 clone 内运行 `glab` 时会自动识别。

## 当技能说「发布到 issue 跟踪器」

创建 GitLab issue。

## 当技能说「获取相关 ticket」

运行 `glab issue view <number> --comments`。
