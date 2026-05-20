# harness — 架构

**harness** 是面向 Agent 工作流的**骨架仓库**（非单一业务应用）：提供 `skills` CLI、编辑器规则/技能规范源、`.agents` 下团队技能真源，以及 `.harness/knowledge` 项目知识库。

## 目录结构（顶层）

```
harness/
├── cli/                 # skills CLI（TypeScript，入口 cli/src/cli.ts）
├── .agents/skills/      # 团队维护技能真源（简体中文 SKILL.md）
├── skills/              # 可供 skills add 安装的 skill 包副本
├── rules/               # Cursor/Claude 规则规范源（cursor/*.mdc, claude/*.md）
├── .cursor/rules/       # 已 bootstrap 的 Cursor 规则（与 rules/cursor 对齐）
├── .claude/rules/       # Claude Code 规则镜像
├── docs/                # agents 文档、cli-rules、ADR
├── .scratch/            # 功能分拣 PRD/issue（可选）
└── .harness/knowledge/  # 本知识库（index.yaml + domains + query）
```

## 入口与运行方式

| 组件 | 入口 | 说明 |
|------|------|------|
| CLI | `cli/src/cli.ts`（`pnpm --dir cli dev`） | 子命令：`add`、`rules add`、`list`、`experimental_install` 等 |
| 规则安装 | `cli/src/rules-add.ts` | 从 Git/本地 `rules/cursor|claude` 复制到消费端 |
| 技能安装 | `cli/src/add.ts` + `installer.ts` | 安装到 `.agents/skills` 等 agent 目录 |

## 技术栈

- **语言：** TypeScript（ESM，`cli/package.json` `type: module`）
- **运行时：** Node ≥ 22（见 `cli/package.json` engines）
- **包管理：** pnpm（`packageManager` 锁定版本）
- **测试：** vitest（`cli/tests/`）
- **UI 交互：** @clack/prompts

## 主要外部依赖（cli）

- Git 克隆：`cli/src/git.ts`
- 多 Agent 目录约定：`cli/src/agents.ts`
- 遥测：`cli/src/telemetry.ts`（可选）

## 与消费端项目的关系

其它仓库通过 `pnpm --dir <harness>/cli dev add|rules add <harness>` 安装技能/规则；安装目标项目根由 **`INIT_CWD` / `--cwd`** 解析（见 `cli/src/project-root.ts`），避免 `pnpm --dir cli` 时误写到 `cli/` 目录。
