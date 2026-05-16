## 文档语言

本仓库中由人与 Agent 维护的可读文档一律使用**简体中文**，包括但不限于：

- `CLAUDE.md`、`docs/agents/`、`CONTEXT.md`、`docs/adr/`
- `.scratch/` 下的 PRD、issue 与评论
- Agent 新建或大幅改写的 markdown 文档

`.claude/`、`.agents/` 下的 `SKILL.md` 等上游技能文件保持英文，不要翻译。

与用户对话时使用简体中文。

## Agent 技能

### Issue tracker（问题跟踪）

Issue 以 markdown 形式存放在 `.scratch/<feature>/`。详见 `docs/agents/issue-tracker.md`。

### Triage labels（分拣状态）

每个 issue 文件顶部的 `Status:` 行表示分拣状态；五个标准角色名 1:1 映射。详见 `docs/agents/triage-labels.md`。

### Domain docs（领域文档）

单上下文布局：仓库根目录的 `CONTEXT.md` 与 `docs/adr/`。详见 `docs/agents/domain.md`。
