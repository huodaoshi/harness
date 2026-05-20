---
description: 仓库文档与 Agent 产出使用简体中文
alwaysApply: true
---

# 文档语言：简体中文

与本仓库相关的**人类可读文档**一律使用简体中文撰写或改写。

## 适用路径

- `CLAUDE.md`
- `docs/**`（含 `docs/agents/`、`docs/adr/`）
- `CONTEXT.md`、`CONTEXT-MAP.md`（仓库根或子目录）
- `.scratch/**`（PRD、issue、评论、分拣备注）

## 技能文件（`SKILL.md`）

- **本仓维护**（`.agents/skills/` 中由团队编写或本地改写的技能）：使用简体中文
- **上游 / vendor**（自外部同步、第三方包或可预期会再同步的技能；含 `.claude/skills/**`）：保持英文，避免同步冲突
- 用户明确要求某一语言时，以用户要求为准

## 排除路径

- `skills-lock.json` 等工具/锁文件（除非用户明确要求翻译）

## 约定保留英文的部分

- 分拣 `Status:` 行的五个标准值：`needs-triage`、`needs-info`、`ready-for-agent`、`ready-for-human`、`wontfix`
- 目录 slug、issue 文件名中的英文片段（如 `.scratch/my-feature/issues/01-auth.md`）
- 代码标识符、CLI 命令、技术 API 名称

## 与用户交流

在 Claude Code 或 Cursor 中处理本仓库时，默认用简体中文回复用户。
