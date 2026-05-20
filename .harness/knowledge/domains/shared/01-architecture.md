# harness — 架构

**harness** 是面向 Agent 工作流的**骨架仓库**（非单一业务应用）：提供 `skills` CLI、编辑器规则/技能规范源、`.agents` 下团队技能真源，以及 `.harness/knowledge` 项目知识库。

## 目录结构（顶层）

```
harness/
├── cli/                      # skills CLI（TypeScript，入口 cli/src/cli.ts）
├── .agents/skills/           # 团队维护技能真源（简体中文 SKILL.md）
├── skills/                   # 可供 skills add 安装的 skill 包副本
├── rules/                    # Cursor/Claude 规则规范源（cursor/*.mdc, claude/*.md）
├── .cursor/rules/            # 已 bootstrap 的 Cursor 规则
├── .claude/rules/            # Claude Code 规则镜像
├── .harness/
│   ├── knowledge/            # 项目知识库（index + domains + query）
│   ├── workflow/             # 功能工作流（feature-workflow.md）
│   └── session/              # sessionStart 注入（session-bootstrap.md）
├── docs/                     # agents 文档、cli-rules、ADR
└── .scratch/                 # 功能分拣 PRD/issue（可选）
```

## 入口与运行方式

| 组件 | 入口 | 说明 |
|------|------|------|
| CLI | `cli/src/cli.ts`（`pnpm --dir cli dev`） | 子命令：`add`、`rules add`、`list`、`experimental_install` 等 |
| 规则安装 | `cli/src/rules-add.ts` | 从 Git/本地 `rules/cursor\|claude` 复制到消费端 |
| 技能安装 | `cli/src/add.ts` + `installer.ts` | 安装到 `.agents/skills` 等 agent 目录 |
| 项目根解析 | `cli/src/project-root.ts` | `INIT_CWD` / `--cwd` 确定消费端根目录 |

## 技术栈

- **语言：** TypeScript（ESM，`cli/package.json` `type: module`）
- **运行时：** Node ≥ 22
- **包管理：** pnpm（`packageManager` 锁定版本）
- **测试：** vitest（`cli/tests/`）
- **UI 交互：** @clack/prompts

## 主要模块（cli/src）

| 模块 | 职责 |
|------|------|
| `add.ts` / `installer.ts` | 技能发现、安装、符号链接/复制 |
| `rules-add.ts` / `rules-lock.ts` | 规则包安装与锁文件 |
| `agents.ts` | 多 Agent 目录约定与检测 |
| `git.ts` / `blob.ts` | Git 克隆与 GitHub tree API |
| `skill-lock.ts` / `local-lock.ts` | 全局锁与项目 `skills-lock.json` |
| `sync.ts` | 从 node_modules 同步技能到 agent 目录 |
