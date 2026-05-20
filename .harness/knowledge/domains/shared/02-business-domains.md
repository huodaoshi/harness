# harness — 业务域与能力模块

本仓无传统「业务服务」；以下为**能力域**划分，便于 Agent 定位代码。

## CLI 命令域（`cli/src/cli.ts`）

| 命令组 | 模块 | 职责 |
|--------|------|------|
| 技能 | `add.ts`, `install.ts`, `installer.ts`, `skills.ts` | 发现、安装、更新技能；`skills-lock` / 本地 lock |
| 规则 | `rules-add.ts`, `rules-lock.ts` | `skills rules add`、从 `rules-lock.json` 恢复 |
| 发现 | `find.ts` | 搜索技能 |
| 列表/移除 | `list.ts`, `remove.ts` | 已安装技能管理 |
| 同步 | `sync.ts` | `experimental_sync`（node_modules 技能） |
| 规则包解析 | `source-parser.ts`, `git.ts`, `blob.ts` | GitHub/本地/Well-known 来源 |

## 规范源域

| 路径 | 内容 |
|------|------|
| `rules/cursor/`、`rules/claude/` | 编辑器规则规范源（`_lang/`、git、docs-zh、windows-shell 等） |
| `skills/*/` | 每个子目录一个可安装 skill（`SKILL.md`） |
| `.agents/skills/*/` | 团队技能真源（init-knowledge、learn、triage、to-prd 等） |

## 文档与分拣

| 路径 | 内容 |
|------|------|
| `docs/agents/` | issue-tracker、triage-labels、domain、knowledge |
| `docs/adr/` | 架构决策（如 skills/rules CLI） |
| `docs/cli-rules.md` | `skills rules` 约定 |
| `.scratch/<feature>/` | PRD、issue markdown 分拣 |

## Agent 技能清单（`.agents/skills/`，18 个）

caveman、diagnose、grill-me、grill-with-docs、handoff、improve-codebase-architecture、init-knowledge、landscape、learn、product-council、prototype、setup-matt-pocock-skills、tdd、to-issues、to-prd、triage、write-a-skill、zoom-out。

## 知识库域

| 路径 | 内容 |
|------|------|
| `.harness/knowledge/index.yaml` | domain 索引、validate_commands |
| `.harness/knowledge/domains/shared/` | 本四件套 |
| `.harness/knowledge/query/` | 按域查询说明 |
| `.agents/skills/init-knowledge/resources/learner-workflow.md` | 学习规程（随技能分发） |
