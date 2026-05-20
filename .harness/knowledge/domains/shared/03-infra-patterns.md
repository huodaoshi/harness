# harness — 基础设施模式

## 路径与安装约定

| 模式 | 说明 |
|------|------|
| 消费端项目根 | `rules add` / `add` 用 `INIT_CWD` 或 `--cwd`，避免装到 `cli/` 目录 |
| Agent 目录 | `agents.ts` 定义 Cursor、Claude、Codex 等目标路径 |
| 锁文件 | 项目根 `skills-lock.json`、`rules-lock.json`；全局 `~/.agents/.skill-lock.json` |
| 规则双写 | 规范源 `rules/`，bootstrap 到 `.cursor/rules`、`.claude/rules` |

## sessionStart hook（阶段 D）

| 路径 | 说明 |
|------|------|
| `.harness/session/session-bootstrap.md` | 新会话注入文案 |
| `.cursor/hooks.json` + `.cursor/hooks/*` | Cursor sessionStart |
| `.claude/hooks/*` + `settings.json` 的 `hooks` | Claude Code SessionStart |
| 模板 | `.agents/skills/init-knowledge/resources/session/`、`resources/hooks/` |

Windows 依赖 **Git Bash**；无 bash 时 hook fail-open。

## 外部依赖

- **Git**：克隆技能包（`git.ts`）
- **GitHub API**：`blob.ts` 拉取 tree（可选 `GITHUB_TOKEN`）
- **网络**：安装远程包时需要；init-knowledge/learn 本身不联网

## 文档与 ADR

- `docs/cli-rules.md` — CLI 规则安装约定
- `docs/adr/0001-skills-cli-rules-install.md` — 结构决策
- `docs/agents/knowledge.md` — 知识库维护说明
