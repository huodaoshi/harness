---
name: setup-matt-pocock-skills
description: 在 AGENTS.md/CLAUDE.md 中配置 `## Agent skills` 块，并生成 `docs/agents/`，使工程类技能了解本仓库的 issue 跟踪器（GitHub 或本地 markdown）、分拣标签词汇表与领域文档布局。在首次使用 `to-issues`、`to-prd`、`triage`、`diagnose`、`tdd`、`improve-codebase-architecture` 或 `zoom-out` 之前运行；或当这些技能似乎缺少 issue 跟踪器、分拣标签或领域文档上下文时运行。
disable-model-invocation: true
---

# 配置 Matt Pocock 技能

搭建工程类技能所依赖的**每仓库**配置：

- **Issue 跟踪器** — issue 存放位置（默认 GitHub；亦原生支持本地 markdown）
- **分拣标签** — 五个规范分拣角色对应的字符串
- **领域文档** — `CONTEXT.md` 与 ADR 的位置及读取规则

这是**由提示驱动**的技能，非确定性脚本。先探索、呈现发现、与用户确认，再写入。

## 流程

### 1. 探索

查看当前仓库的初始状态。阅读已有内容，勿臆测：

- `git remote -v` 与 `.git/config` — 是否为 GitHub 仓库？哪个？
- 根目录 `AGENTS.md`、`CLAUDE.md` — 是否存在？是否已有 `## Agent skills` 小节？
- 根目录 `CONTEXT.md`、`CONTEXT-MAP.md`
- `docs/adr/` 及任意 `src/*/docs/adr/`
- `docs/agents/` — 本技能先前输出是否已存在？
- `.scratch/` — 是否已在使用本地 markdown issue 跟踪器约定

### 2. 呈现发现并询问

总结已有与缺失项。然后**逐项**带用户做三项决策——先呈现一节，得到回答后再下一节，勿一次抛出全部。

假定用户可能不懂这些术语。每节开头用简短说明（是什么、为何这些技能需要、选不同项会如何变化），再给出选项与默认值。

**A 节 — Issue 跟踪器。**

> 说明：「issue 跟踪器」是本仓库存放 issue 的地方。`to-issues`、`triage`、`to-prd`、`qa` 等技能会读写它——需知是调用 `gh issue create`、在 `.scratch/` 下写 markdown，还是你描述的其他工作流。请选择你**实际**用来跟踪本仓库工作的方式。

默认倾向：这些技能面向 GitHub 设计。若 `git remote` 指向 GitHub，建议选 GitHub。若指向 GitLab（`gitlab.com` 或自托管），建议选 GitLab。否则（或用户偏好）提供：

- **GitHub** — issue 在仓库的 GitHub Issues（使用 `gh` CLI）
- **GitLab** — issue 在仓库的 GitLab Issues（使用 [`glab`](https://gitlab.com/gitlab-org/cli) CLI）
- **本地 markdown** — issue 以 `.scratch/<feature>/` 下文件形式存放（适合个人项目或无远程仓库）
- **其他**（Jira、Linear 等）— 请用户用一段话描述工作流；技能将其记录为自由文本

**B 节 — 分拣标签词汇表。**

> 说明：`triage` 技能处理新进 issue 时，会将其在状态机中流转——待评估、等报告者、可供 AFK Agent 领取、需人工实现、不予处理。为此须在 issue 跟踪器上应用**你已配置**的标签（或等价物）。若仓库已用不同标签名（如 `bug:triage` 而非 `needs-triage`），在此映射，避免技能创建重复标签。

五个规范角色：

- `needs-triage` — 维护者需评估
- `needs-info` — 等待报告者
- `ready-for-agent` — 规格完整，可供 AFK Agent 领取（Agent 无需额外人工上下文即可开工）
- `ready-for-human` — 需人工实现
- `wontfix` — 不予处理

默认：各角色字符串与角色名相同。询问用户是否要覆盖。若跟踪器尚无标签，默认值即可。

**C 节 — 领域文档。**

> 说明：部分技能（`improve-codebase-architecture`、`diagnose`、`tdd`）会读 `CONTEXT.md` 了解项目领域语言，并读 `docs/adr/` 了解既往架构决策。需确认仓库是单一全局上下文还是多上下文（如前后端分离的 monorepo），以便在正确位置查找。

确认布局：

- **单上下文** — 仓库根目录一个 `CONTEXT.md` + `docs/adr/`。多数仓库如此。
- **多上下文** — 根目录 `CONTEXT-MAP.md` 指向各上下文的 `CONTEXT.md`（常见于 monorepo）。

### 3. 确认并编辑

向用户展示草稿：

- 将写入 `CLAUDE.md` / `AGENTS.md` 的 `## Agent skills` 块（见步骤 4 选择规则）
- `docs/agents/issue-tracker.md`、`docs/agents/triage-labels.md`、`docs/agents/domain.md` 的内容

写入前允许用户修改。

### 4. 写入

**选择要编辑的文件：**

- 若存在 `CLAUDE.md`，编辑它。
- 否则若存在 `AGENTS.md`，编辑它。
- 若二者皆无，询问用户要创建哪一个——勿替用户决定。

当 `CLAUDE.md` 已存在时勿再创建 `AGENTS.md`（反之亦然）——始终编辑已有文件。

若所选文件中已有 `## Agent skills` 块，就地更新其内容，勿追加重复块。勿覆盖用户对周围小节的编辑。

块内容：

```markdown
## Agent skills

### Issue tracker

[一行摘要：issue 存放在哪]. 详见 `docs/agents/issue-tracker.md`。

### Triage labels

[一行摘要：标签词汇表]. 详见 `docs/agents/triage-labels.md`。

### Domain docs

[一行摘要：single-context 或 multi-context 布局]. 详见 `docs/agents/domain.md`。
```

然后以本技能目录下的种子模板为起点，写入三个 `docs/agents/` 文件：

- [issue-tracker-github.md](./issue-tracker-github.md) — GitHub issue 跟踪器
- [issue-tracker-gitlab.md](./issue-tracker-gitlab.md) — GitLab issue 跟踪器
- [issue-tracker-local.md](./issue-tracker-local.md) — 本地 markdown issue 跟踪器
- [triage-labels.md](./triage-labels.md) — 标签映射
- [domain.md](./domain.md) — 领域文档消费规则与布局

对「其他」跟踪器，根据用户描述从头撰写 `docs/agents/issue-tracker.md`。

### 5. 完成

告知用户配置已完成，以及哪些工程类技能将从此读取。说明日后可直接编辑 `docs/agents/*.md`——仅当要切换 issue 跟踪器或从头重来时才需重新运行本技能。
