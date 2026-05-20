# harness

Agent 工作流骨架仓库：在团队与业务项目之间，统一 **技能（Skills）**、**编辑器规则（Rules）**、**项目知识库（Knowledge）** 与 **会话/功能工作流** 的落地方式。

适合作为：

- **harness 自身** — 维护 CLI、团队技能、规则规范源与文档真源；
- **技能/规则包来源** — 通过 `skills add`、`skills rules add` 安装到其它 Cursor / Claude Code 项目；
- **消费端模板参考** — 用 `init-knowledge` 在目标仓库生成 `.harness/`、hook 与 `CLAUDE.md`。

人类可读文档默认使用**简体中文**（见 `.cursor/rules/docs-zh.mdc`）。面向 Agent 的入口另有 **`CLAUDE.md`**（项目总览与实现阶段约定）。

---

## 目录

- [核心能力](#核心能力)
- [仓库结构](#仓库结构)
- [快速开始](#快速开始)
- [在其它项目中使用](#在其它项目中使用)
- [Agent 工作流](#agent-工作流)
- [团队技能一览](#团队技能一览)
- [知识体系（三层）](#知识体系三层)
- [CLI 开发](#cli-开发)
- [文档索引](#文档索引)
- [约定与贡献](#约定与贡献)

---

## 核心能力

| 能力 | 说明 |
|------|------|
| **skills CLI** | 从 Git / 本地路径安装 Agent 技能；支持 `rules add` 安装 Cursor/Claude 规则包（见 [`cli/README.md`](cli/README.md)） |
| **团队技能** | `.agents/skills/` 为真源；`skills/` 为可 `skills add` 的副本 |
| **规则规范源** | `rules/cursor`、`rules/claude` → 安装到 `.cursor/rules`、`.claude/rules` |
| **项目知识库** | `.harness/knowledge/`：索引、按域四件套、query、学习规程 |
| **功能工作流** | `.harness/workflow/feature-workflow.md`：新功能 / 重构 / 修 bug / 领 issue 的完整流程 |
| **会话引导** | `sessionStart` hook + `.harness/session/session-bootstrap.md` 新会话注入（**默认 caveman**） |
| **Issue 分拣（可选）** | `.scratch/<feature>/` 本地 markdown issue，见 [`docs/agents/issue-tracker.md`](docs/agents/issue-tracker.md) |

---

## 仓库结构

```
harness/
├── cli/                      # skills CLI（TypeScript，Node ≥ 22）
├── .agents/skills/           # 团队技能真源（维护时改这里）
├── skills/                   # 供 skills add 安装的技能副本
├── rules/                    # Cursor/Claude 规则规范源
│   ├── cursor/               # *.mdc → .cursor/rules/
│   └── claude/               # *.md → .claude/rules/
├── .cursor/                  # Cursor：已 bootstrap 的规则 + hooks
├── .claude/                  # Claude Code：规则镜像 + hooks + settings.json
├── .harness/
│   ├── knowledge/            # 知识库（index.yaml、domains、query、learner-workflow）
│   ├── workflow/             # feature-workflow.md
│   └── session/              # session-bootstrap.md
├── docs/
│   ├── agents/               # issue 跟踪、分拣、知识库说明等
│   ├── adr/                  # 架构决策记录
│   └── cli-rules.md          # rules add 约定
├── .scratch/                 # 功能分拣 PRD/issue（可选）
├── CLAUDE.md                 # Agent 项目总览（init-knowledge 可生成）
└── README.md                 # 本文件（面向人类）
```

**维护约定：**

- 改技能 → 先改 **`.agents/skills/<name>/`**，再同步到 **`skills/<name>/`**
- 改规则 → 只改 **`rules/`**，需要 IDE 生效时执行 `rules add`（见 [`rules/README.md`](rules/README.md)）
- 改知识库流程 → **`.harness/knowledge/learner-workflow.md`**，并复制到 `init-knowledge/resources/`
- 改功能工作流 → **`.harness/workflow/feature-workflow.md`**，并复制到 `init-knowledge/resources/workflow/`

---

## 快速开始

### 环境要求

- **Node.js** ≥ 22（CLI 见 `cli/package.json`）
- **pnpm**（CLI 目录已锁定 `packageManager` 版本）
- Windows 上跑 hook 脚本需 **Git Bash**（或 PATH 上的 `bash`）

### 克隆与 CLI

```powershell
git clone <你的 harness 仓库 URL>
cd harness\cli
pnpm install
```

在 **harness 仓库根** 使用 CLI（开发模式）：

```powershell
cd D:\harness

# 列出本仓库可作为技能包提供的技能
pnpm --dir cli dev add . --list

# 安装指定技能到当前目录（示例）
pnpm --dir cli dev add . --skill init-knowledge --skill learn -y

# 将本仓库规则安装到当前项目的 .cursor/rules
pnpm --dir cli dev rules add . -a cursor
```

发布版也可通过 `npx skills` 使用（见 [`cli/README.md`](cli/README.md)）。

### 本仓自检

```powershell
pnpm --dir cli exec vitest run
pnpm --dir cli run type-check
```

---

## 在其它项目中使用

在**目标项目根**打开终端（`cd` 到业务仓库），用本机 harness 路径调用 CLI。

```powershell
cd D:\your-app

# 安装团队技能（示例）
pnpm --dir D:\harness\cli dev add D:\harness --skill init-knowledge --skill learn -y

# 安装编辑器规则
pnpm --dir D:\harness\cli dev rules add D:\harness -a cursor
pnpm --dir D:\harness\cli dev rules add D:\harness -a claude-code
```

**重要：** 使用 `pnpm --dir …\cli dev …` 时，CLI 会优先用 **`INIT_CWD`**（你发起命令时的目录）作为项目根，避免装到 `cli/` 下。仍不对时加 **`--cwd D:\your-app`**。详见 [`docs/cli-rules.md`](docs/cli-rules.md)。

### 初始化目标项目的知识库与 hook

在 Cursor 中加载 **init-knowledge** 技能，对目标仓库执行（或在 harness 根传入目标绝对路径）：

```text
/init-knowledge
/init-knowledge D:\your-app
```

默认会生成：

- `.harness/knowledge/`（索引 + domain 文档 + query）
- `.harness/workflow/feature-workflow.md`
- `.harness/session/session-bootstrap.md`
- `.cursor/hooks*`、`.claude/hooks*`（阶段 D，可用 `--no-hooks` 跳过）
- 可选更新 `CLAUDE.md`、语言规则

常用参数：`--full` 强制覆盖、`--hooks` 仅重装 hook、`--workflow` 仅重装功能工作流、`--domain <name>` 只刷新某一 domain。

跨仓库 init 时请在 **harness 根**执行，或设置 `HARNESS_ROOT` 指向 harness。详见 [`.agents/skills/init-knowledge/SKILL.md`](.agents/skills/init-knowledge/SKILL.md) 与 [`docs/agents/knowledge.md`](docs/agents/knowledge.md)。

---

## Agent 工作流

### 功能需求（先读工作流）

做**新功能、重构、修 bug、领取 issue** 前，Agent 应阅读 **[`.harness/workflow/feature-workflow.md`](.harness/workflow/feature-workflow.md)**：

1. 按文首判断属于哪一种工作类型（四选一）
2. **只跟对应章节的完整流程**，不要混用
3. 遵守全局硬规则：开工前有 scope（issue / PRD / 当轮验收标准）；收工前跑 `index.yaml` 中的 `validate_commands`

### 实现阶段（写代码）

见 [`CLAUDE.md`](CLAUDE.md) 中「实现阶段」：读 `index.yaml` → `query/<domain>.md` → 在 `domain.path` 下找相似实现 → 跑验证命令。

### Issue 与分拣（本仓可选）

- Issue：`.scratch/<feature>/`，约定见 [`docs/agents/issue-tracker.md`](docs/agents/issue-tracker.md)
- 分拣标签：[`docs/agents/triage-labels.md`](docs/agents/triage-labels.md)
- 技能：**triage**、**to-issues**、**to-prd**

首次使用工程类技能前，可运行 **setup-matt-pocock-skills** 配置 `docs/agents/` 与 `CLAUDE.md` 中的 Agent skills 块。

---

## 团队技能一览

技能真源在 **`.agents/skills/<name>/SKILL.md`**。在 Cursor 中通过 Skill 工具或 `@` 引用加载。

| 技能 | 用途摘要 |
|------|----------|
| **init-knowledge** | 初始化/重建 `.harness/knowledge`、workflow、hook、`CLAUDE.md` |
| **learn** | 开发后按 git 变更增量更新知识文档 |
| **setup-matt-pocock-skills** | 配置 issue 跟踪器、分拣标签、领域文档布局 |
| **triage** | Issue 分拣状态机 |
| **to-issues** / **to-prd** | 计划拆 issue、写 PRD |
| **grill-me** / **grill-with-docs** | 拷问式澄清方案；后者对照 CONTEXT/ADR |
| **improve-codebase-architecture** | 结合领域文档寻找架构改进点 |
| **prototype** | 可丢弃原型（终端或 UI 变体） |
| **diagnose** | 难解 bug / 性能回退诊断循环 |
| **tdd** | 测试驱动开发 |
| **landscape** / **product-council** | 竞品摸底、产品议会（探索向） |
| **handoff** | 压缩对话为交接文档 |
| **caveman** | 极简回复模式 |
| **write-a-skill** | 撰写新技能 |
| **zoom-out** | （见各技能 SKILL 说明） |

安装副本与对外分发见 [`skills/README.md`](skills/README.md)。

---

## 知识体系（三层）

面向 Agent 的渐进式阅读顺序（本仓已 init 时）：

| 层级 | 路径 | 内容 |
|------|------|------|
| 1 | `.cursor/rules/`、`.claude/rules/` | 编辑器规则（语言、Git、文档语言、Windows 终端等） |
| 2 | `.harness/knowledge/query/<domain>.md` | 该 domain 应先读哪些知识文件 |
| 3 | `.harness/knowledge/domains/<type>/` | 架构、业务域、基础设施模式、开发指南（各 4 份） |

索引：**[`.harness/knowledge/index.yaml`](.harness/knowledge/index.yaml)**  
学习规程：**[`.harness/knowledge/learner-workflow.md`](.harness/knowledge/learner-workflow.md)**

---

## CLI 开发

| 命令 | 说明 |
|------|------|
| `pnpm --dir cli dev` | 以开发入口运行 CLI（`src/cli.ts`） |
| `pnpm --dir cli test` | Vitest |
| `pnpm --dir cli run type-check` | `tsc --noEmit` |
| `pnpm --dir cli run build` | 构建发布产物 |

架构与模块说明见 [`cli/AGENTS.md`](cli/AGENTS.md)、[`cli/README.md`](cli/README.md)。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [CLAUDE.md](CLAUDE.md) | Agent 项目总览与实现约定 |
| [docs/agents/knowledge.md](docs/agents/knowledge.md) | 知识库与 init/learn |
| [docs/agents/issue-tracker.md](docs/agents/issue-tracker.md) | 本地 markdown issue |
| [docs/agents/triage-labels.md](docs/agents/triage-labels.md) | 分拣标签 |
| [docs/agents/domain.md](docs/agents/domain.md) | CONTEXT / ADR 布局 |
| [docs/cli-rules.md](docs/cli-rules.md) | `rules add` 与 `INIT_CWD` / `--cwd` |
| [docs/adr/0001-skills-cli-rules-install.md](docs/adr/0001-skills-cli-rules-install.md) | 规则安装 ADR |
| [rules/README.md](rules/README.md) | 规则规范源与跨项目安装 |
| [skills/README.md](skills/README.md) | 技能副本与 `skills add` 示例 |

---

## 约定与贡献

- **提交信息**：`<type>(<scope>): <subject>`，详见 [`.cursor/rules/git.mdc`](.cursor/rules/git.mdc)
- **终端（Windows）**：PowerShell 5.x 勿用 `&&` 连接命令，用 `;` 或分多条执行，见 [`.cursor/rules/windows-shell.mdc`](.cursor/rules/windows-shell.mdc)
- **禁止**在未经用户要求时由 Agent 自动 `git commit`（团队约定，见 `CLAUDE.md`）

修改 harness 本体功能后，建议在合并前运行 `vitest` 与 `type-check`；结构性变更可运行 **learn** 更新 `.harness/knowledge/domains/`。

---

## 相关链接

- CLI 详细用法：[cli/README.md](cli/README.md)
- 上游 skills 生态概念与 `npx skills`：[skills.sh](https://skills.sh)（本仓库 CLI 在其基础上扩展了 rules 安装等能力）
