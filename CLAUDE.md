# harness

> Agent 工作流骨架：skills CLI、规则/技能规范源、团队技能与项目知识库。

## 仓库结构

```
harness/
├── cli/                 # skills CLI（TypeScript）
├── .agents/skills/      # 团队技能真源
├── skills/              # skills add 可安装的 skill 副本
├── rules/               # Cursor/Claude 规则规范源
├── .cursor/rules/       # Cursor 规则（bootstrap）
├── .claude/rules/       # Claude Code 规则镜像
├── docs/agents/         # Agent 协作文档
├── docs/adr/            # 架构决策
├── .scratch/            # 功能分拣（可选）
├── .harness/
│   ├── knowledge/       # 项目知识库
│   ├── workflow/        # 功能工作流（feature-workflow.md）
│   └── session/         # 会话引导（session-bootstrap.md）
```

## 技术栈

| 区域 | 技术 |
|------|------|
| CLI | TypeScript (ESM)、Node ≥ 22、pnpm、vitest、@clack/prompts |
| 规则 | Cursor `.mdc`、Claude `.md` |
| 文档 | Markdown（默认简体中文） |

## 域清单

- **harness**（shared, typescript）— 仓库根 `./`；知识见 `.harness/knowledge/domains/shared/`

## 常用命令

```bash
pnpm --dir cli exec vitest run
pnpm --dir cli run type-check
pnpm --dir cli dev add . --list
pnpm --dir cli dev rules add . -a cursor
```

## 全局约定

- **禁止自动提交**：除非用户明确要求，不要执行 `git commit`
- **提交格式**：`<type>(<scope>): <subject>`，详见 `.cursor/rules/git.mdc`
- **文档语言**：人类可读文档默认**简体中文**（见 `.cursor/rules/docs-zh.mdc`）
- **Windows 终端**：勿用 `&&` 连接命令，见 `windows-shell.mdc`
- 其他约定见 **`.cursor/rules/`**、**`.claude/rules/`**、**`.harness/knowledge/`** 与 **`.harness/workflow/`**

## Agent 工作流（harness）

### Issue 与分拣（`.scratch/`）

- Issue 存放在 **`.scratch/<feature>/`**，详见 `docs/agents/issue-tracker.md`
- 文件顶部 **`Status:`** 与分拣角色见 `docs/agents/triage-labels.md`
- 新建或分拣 issue 时使用 **triage** 技能

### 领域文档（可选）

- 根目录 **`CONTEXT.md`**、**`docs/adr/`** 按需建立，见 `docs/agents/domain.md`
- 探索代码前若存在上述文件应先读；不存在则静默继续

### 功能工作流

做**新功能、重构、修 bug、领 issue** 前，先 **Read** **`.harness/workflow/feature-workflow.md`**：

1. 按文首「先选哪一种工作类型」进入对应章节（**只跟一条流程，不要混用**）
2. 遵守文内「全局硬规则」（开工前有 scope、收工前跑 `validate_commands`）
3. 未在该章节列出的步骤**默认不做**

### 实现阶段（写代码时）

1. 读取 **`.harness/knowledge/index.yaml`** 中目标 domain 的 `rules` 与 `skills[]`
2. 先读 **`.harness/knowledge/query/<domain>.md`**（若存在）
3. 在 `domain.path` 下找相似实现
4. 按分层实现并运行 `validate_commands`

## 知识体系（渐进式披露）

### 第一层：Rules

**`.cursor/rules/`**（`*.mdc`）与 **`.claude/rules/`**（`*.md`，若存在）按各文件 frontmatter 的 `paths` 匹配。

### 第二层：按域查询说明

**`.harness/knowledge/query/<domain>.md`** — 由 **init-knowledge** 生成，说明该 domain 应先读哪些知识文件。

### 第三层：Domain 知识文档

**`.harness/knowledge/domains/<type>/`** 下四份文档：

- `01-architecture.md`
- `02-business-domains.md`
- `03-infra-patterns.md`
- `04-dev-guide.md`

## 自学习与索引

- **增量学习**：**learn** 技能；先 **Read** **`.harness/knowledge/learner-workflow.md`**，并依据 **`.harness/knowledge/index.yaml`**
- **全量 / 首次**：在 **harness** 仓库根对目标仓运行 **init-knowledge**（可传目标仓绝对路径）
- **项目索引**：**`.harness/knowledge/index.yaml`**；domain 划分变更后重跑 **init-knowledge**（如 `--domain <name>`）
- **规则包同步**（可选）：`skills rules add <harness-repo> -a cursor` 等，见 `rules/README.md`
