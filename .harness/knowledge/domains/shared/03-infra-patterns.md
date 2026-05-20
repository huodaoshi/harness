# harness — 基础设施与横切模式

## 锁文件与安装目标

| 文件 | 位置 | 用途 |
|------|------|------|
| `skills-lock.json` | 消费端项目根（全局 lock 在 XDG/`~/.agents`） | 技能版本与来源追踪 |
| `rules-lock.json` | 消费端项目根 | 规则包恢复（`skills rules experimental_install`） |
| `skills-lock.json`（本仓根） | harness 根 | 本仓 CLI 开发时产生的 lock（若存在） |

## 项目根解析（CLI）

`cli/src/project-root.ts`：`resolveCliProjectRoot()`  
顺序：**`--cwd`** → **`INIT_CWD`** → **`process.cwd()`**。  
用于 `skills add`、`skills rules add`、`experimental_install`、`experimental_sync`。

## Agent 目录约定

- 通用技能目录：`.agents/skills/`（多数 universal agent）
- Cursor：`.cursor/rules/`（`*.mdc`）
- Claude Code：`.claude/rules/`（`*.md`）
- 定义表：`cli/src/agents.ts`（50+ agent 的 skillsDir / globalSkillsDir）

## 规则双写与规范源

- **编辑源头：** `rules/cursor/`、`rules/claude/`
- **IDE 使用：** `.cursor/rules/`、`.claude/rules/`（bootstrap 后可用 `rules add .` 再同步）
- **语言规则：** `_lang/*.mdc` / `_lang/*.md`，frontmatter `paths` 控制生效范围

## 文档与语言

- 人类文档默认**简体中文**（`.cursor/rules/docs-zh.mdc`）
- Windows 终端：PowerShell 5.x 避免 `&&`（`windows-shell.mdc`）
- Git：默认分支 `master`，feature 分支命名见 `git.mdc`

## 测试与质量

- 测试目录：`cli/tests/`（vitest）
- 类型检查：`pnpm --dir cli run type-check`
- 格式化：Prettier（`cli/.prettierrc`，tabWidth 4）

## 无服务端组件

本仓**无**数据库、消息队列、HTTP API 服务；CLI 为本地 Node 进程，Git 与网络仅用于 clone/遥测。
