# skills

开放 Agent 技能生态系统的命令行工具。

<!-- agent-list:start -->

支持 **OpenCode**、**Claude Code**、**Codex**、**Cursor** 以及[另外 51 个](#supported-agents)。

<!-- agent-list:end -->

[![skills.sh](https://skills.sh/b/vercel-labs/skills)](https://skills.sh/vercel-labs/skills)

## 安装编辑器规则（Cursor / Claude Code）

自本 fork 起的 CLI 支持从 Git 规则包安装项目规则（包内需含 `rules/cursor/` 或 `rules/claude/`）。完整约定见 **[`docs/cli-rules.md`](../docs/cli-rules.md)** 与 ADR `docs/adr/0001-skills-cli-rules-install.md`。

```bash
pnpm --dir cli dev rules add <owner/repo> -a cursor
pnpm --dir cli dev rules add <owner/repo> -a claude-code
pnpm --dir cli dev rules experimental_install
```

安装成功后写入/更新项目根 **`rules-lock.json`**（与 `skills-lock.json` 独立）。

## 安装技能

```bash
npx skills add vercel-labs/agent-skills
```

### 来源格式

```bash
# GitHub 简写（owner/repo）
npx skills add vercel-labs/agent-skills

# 完整 GitHub URL
npx skills add https://github.com/vercel-labs/agent-skills

# 仓库内技能的直接路径
npx skills add https://github.com/vercel-labs/agent-skills/tree/main/skills/web-design-guidelines

# GitLab URL
npx skills add https://gitlab.com/org/repo

# 任意 git URL
npx skills add git@github.com:vercel-labs/agent-skills.git

# 本地路径
npx skills add ./my-local-skills
```

### 选项

| 选项                      | 说明                                                                                                                                        |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `-g, --global`            | 安装到用户目录，而非项目目录                                                                                                                |
| `-a, --agent <agents...>` | <!-- agent-names:start -->指定 Agent（如 `claude-code`、`codex`）。见[支持的 Agent](#supported-agents)<!-- agent-names:end --> |
| `-s, --skill <skills...>` | 按名称安装指定技能（全部技能用 `'*'`）                                                                                                      |
| `-l, --list`              | 仅列出可用技能，不安装                                                                                                                      |
| `--copy`                  | 复制文件到各 Agent 目录，而非符号链接                                                                                                       |
| `-y, --yes`               | 跳过所有确认提示                                                                                                                            |
| `--all`                   | 无提示地将全部技能安装到全部 Agent                                                                                                          |

### 示例

```bash
# 列出仓库中的技能
npx skills add vercel-labs/agent-skills --list

# 安装指定技能
npx skills add vercel-labs/agent-skills --skill frontend-design --skill skill-creator

# 安装名称含空格的技能（须加引号）
npx skills add owner/repo --skill "Convex Best Practices"

# 安装到指定 Agent
npx skills add vercel-labs/agent-skills -a claude-code -a opencode

# 非交互安装（适合 CI/CD）
npx skills add vercel-labs/agent-skills --skill frontend-design -g -a claude-code -y

# 将仓库中全部技能安装到全部 Agent
npx skills add vercel-labs/agent-skills --all

# 将全部技能安装到指定 Agent
npx skills add vercel-labs/agent-skills --skill '*' -a claude-code

# 将指定技能安装到全部 Agent
npx skills add vercel-labs/agent-skills --agent '*' --skill frontend-design
```

### 安装范围

| 范围       | 标志      | 位置                | 适用场景                         |
| ---------- | --------- | ------------------- | -------------------------------- |
| **项目**   | （默认）  | `./<agent>/skills/` | 随项目提交，与团队共享           |
| **全局**   | `-g`      | `~/<agent>/skills/` | 在所有项目中可用                 |

### 安装方式

交互安装时可选择：

| 方式                 | 说明                                                                                 |
| -------------------- | ------------------------------------------------------------------------------------ |
| **符号链接**（推荐） | 各 Agent 指向规范副本的符号链接。单一事实来源，便于更新。                            |
| **复制**             | 为每个 Agent 创建独立副本。在不支持符号链接时使用。                                  |

## 其他命令

| 命令                         | 说明                                   |
| ---------------------------- | -------------------------------------- |
| `npx skills list`            | 列出已安装技能（别名：`ls`）           |
| `npx skills find [query]`    | 交互式或按关键词搜索技能               |
| `npx skills remove [skills]` | 从各 Agent 移除已安装技能              |
| `npx skills update [skills]` | 将已安装技能更新到最新版本             |
| `npx skills init [name]`     | 创建新的 SKILL.md 模板                 |

### `skills list`

列出所有已安装技能，类似 `npm ls`。

```bash
# 列出全部已安装技能（项目 + 全局）
npx skills list

# 仅列出全局技能
npx skills ls -g

# 按 Agent 筛选
npx skills ls -a claude-code -a cursor
```

### `skills find`

交互式或按关键词搜索技能。

```bash
# 交互搜索（类 fzf）
npx skills find

# 按关键词搜索
npx skills find typescript
```

### `skills update`

```bash
# 更新全部技能（交互式选择范围）
npx skills update

# 按名称更新单个技能
npx skills update my-skill

# 更新多个指定技能
npx skills update frontend-design web-design-guidelines

# 仅更新全局或项目技能
npx skills update -g
npx skills update -p

# 非交互（自动检测范围：在项目目录则为项目，否则为全局）
npx skills update -y
```

| 选项            | 说明                                                                       |
| --------------- | -------------------------------------------------------------------------- |
| `-g, --global`  | 仅更新全局技能                                                             |
| `-p, --project` | 仅更新项目技能                                                             |
| `-y, --yes`     | 跳过范围提示（自动检测：在项目目录则为项目，否则为全局）                   |
| `[skills...]`   | 按名称更新指定技能，而非全部                                                 |

### `skills init`

```bash
# 在当前目录创建 SKILL.md
npx skills init

# 在子目录中创建新技能
npx skills init my-skill
```

### `skills remove`

从各 Agent 移除已安装技能。

```bash
# 交互移除（从已安装列表中选择）
npx skills remove

# 按名称移除指定技能
npx skills remove web-design-guidelines

# 移除多个技能
npx skills remove frontend-design web-design-guidelines

# 从全局范围移除
npx skills remove --global web-design-guidelines

# 仅从指定 Agent 移除
npx skills remove --agent claude-code cursor my-skill

# 无确认移除全部已安装技能
npx skills remove --all

# 从指定 Agent 移除全部技能
npx skills remove --skill '*' -a cursor

# 从全部 Agent 移除指定技能
npx skills remove my-skill --agent '*'

# 使用 rm 别名
npx skills rm my-skill
```

| 选项           | 说明                                   |
| -------------- | -------------------------------------- |
| `-g, --global` | 从全局范围（~/）移除，而非项目         |
| `-a, --agent`  | 从指定 Agent 移除（全部用 `'*'`）      |
| `-s, --skill`  | 指定要移除的技能（全部用 `'*'`）       |
| `-y, --yes`    | 跳过确认提示                           |
| `--all`        | 等同于 `--skill '*' --agent '*' -y`    |

## 什么是 Agent 技能？

Agent 技能是可复用的指令集，用于扩展编程 Agent 的能力。它们定义在带 YAML frontmatter 的 `SKILL.md` 文件中，需包含 `name` 与 `description`。

技能让 Agent 能执行例如：

- 根据 git 历史生成发布说明
- 按团队规范创建 PR
- 与外部工具集成（Linear、Notion 等）

在 **[skills.sh](https://skills.sh)** 发现更多技能。

## 支持的 Agent

技能可安装到以下任一 Agent：

<!-- supported-agents:start -->

| Agent                                 | `--agent`                                | 项目路径                 | 全局路径                        |
| ------------------------------------- | ---------------------------------------- | ------------------------ | ------------------------------- |
| AiderDesk                             | `aider-desk`                             | `.aider-desk/skills/`    | `~/.aider-desk/skills/`         |
| Amp, Kimi Code CLI, Replit, Universal | `amp`, `kimi-cli`, `replit`, `universal` | `.agents/skills/`        | `~/.config/agents/skills/`      |
| Antigravity                           | `antigravity`                            | `.agents/skills/`        | `~/.gemini/antigravity/skills/` |
| Augment                               | `augment`                                | `.augment/skills/`       | `~/.augment/skills/`            |
| IBM Bob                               | `bob`                                    | `.bob/skills/`           | `~/.bob/skills/`                |
| Claude Code                           | `claude-code`                            | `.claude/skills/`        | `~/.claude/skills/`             |
| OpenClaw                              | `openclaw`                               | `skills/`                | `~/.openclaw/skills/`           |
| Cline, Dexto, Warp                    | `cline`, `dexto`, `warp`                 | `.agents/skills/`        | `~/.agents/skills/`             |
| CodeArts Agent                        | `codearts-agent`                         | `.codeartsdoer/skills/`  | `~/.codeartsdoer/skills/`       |
| CodeBuddy                             | `codebuddy`                              | `.codebuddy/skills/`     | `~/.codebuddy/skills/`          |
| Codemaker                             | `codemaker`                              | `.codemaker/skills/`     | `~/.codemaker/skills/`          |
| Code Studio                           | `codestudio`                             | `.codestudio/skills/`    | `~/.codestudio/skills/`         |
| Codex                                 | `codex`                                  | `.agents/skills/`        | `~/.codex/skills/`              |
| Command Code                          | `command-code`                           | `.commandcode/skills/`   | `~/.commandcode/skills/`        |
| Continue                              | `continue`                               | `.continue/skills/`      | `~/.continue/skills/`           |
| Cortex Code                           | `cortex`                                 | `.cortex/skills/`        | `~/.snowflake/cortex/skills/`   |
| Crush                                 | `crush`                                  | `.crush/skills/`         | `~/.config/crush/skills/`       |
| Cursor                                | `cursor`                                 | `.agents/skills/`        | `~/.cursor/skills/`             |
| Deep Agents                           | `deepagents`                             | `.agents/skills/`        | `~/.deepagents/agent/skills/`   |
| Devin for Terminal                    | `devin`                                  | `.devin/skills/`         | `~/.config/devin/skills/`       |
| Droid                                 | `droid`                                  | `.factory/skills/`       | `~/.factory/skills/`            |
| Firebender                            | `firebender`                             | `.agents/skills/`        | `~/.firebender/skills/`         |
| ForgeCode                             | `forgecode`                              | `.forge/skills/`         | `~/.forge/skills/`              |
| Gemini CLI                            | `gemini-cli`                             | `.agents/skills/`        | `~/.gemini/skills/`             |
| GitHub Copilot                        | `github-copilot`                         | `.agents/skills/`        | `~/.copilot/skills/`            |
| Goose                                 | `goose`                                  | `.goose/skills/`         | `~/.config/goose/skills/`       |
| Hermes Agent                          | `hermes-agent`                           | `.hermes/skills/`        | `~/.hermes/skills/`             |
| Junie                                 | `junie`                                  | `.junie/skills/`         | `~/.junie/skills/`              |
| iFlow CLI                             | `iflow-cli`                              | `.iflow/skills/`         | `~/.iflow/skills/`              |
| Kilo Code                             | `kilo`                                   | `.kilocode/skills/`      | `~/.kilocode/skills/`           |
| Kiro CLI                              | `kiro-cli`                               | `.kiro/skills/`          | `~/.kiro/skills/`               |
| Kode                                  | `kode`                                   | `.kode/skills/`          | `~/.kode/skills/`               |
| MCPJam                                | `mcpjam`                                 | `.mcpjam/skills/`        | `~/.mcpjam/skills/`             |
| Mistral Vibe                          | `mistral-vibe`                           | `.vibe/skills/`          | `~/.vibe/skills/`               |
| Mux                                   | `mux`                                    | `.mux/skills/`           | `~/.mux/skills/`                |
| OpenCode                              | `opencode`                               | `.agents/skills/`        | `~/.config/opencode/skills/`    |
| OpenHands                             | `openhands`                              | `.openhands/skills/`     | `~/.openhands/skills/`          |
| Pi                                    | `pi`                                     | `.pi/skills/`            | `~/.pi/agent/skills/`           |
| Qoder                                 | `qoder`                                  | `.qoder/skills/`         | `~/.qoder/skills/`              |
| Qwen Code                             | `qwen-code`                              | `.qwen/skills/`          | `~/.qwen/skills/`               |
| Rovo Dev                              | `rovodev`                                | `.rovodev/skills/`       | `~/.rovodev/skills/`            |
| Roo Code                              | `roo`                                    | `.roo/skills/`           | `~/.roo/skills/`                |
| Tabnine CLI                           | `tabnine-cli`                            | `.tabnine/agent/skills/` | `~/.tabnine/agent/skills/`      |
| Trae                                  | `trae`                                   | `.trae/skills/`          | `~/.trae/skills/`               |
| Trae CN                               | `trae-cn`                                | `.trae/skills/`          | `~/.trae-cn/skills/`            |
| Windsurf                              | `windsurf`                               | `.windsurf/skills/`      | `~/.codeium/windsurf/skills/`   |
| Zencoder                              | `zencoder`                               | `.zencoder/skills/`      | `~/.zencoder/skills/`           |
| Neovate                               | `neovate`                                | `.neovate/skills/`       | `~/.neovate/skills/`            |
| Pochi                                 | `pochi`                                  | `.pochi/skills/`         | `~/.pochi/skills/`              |
| AdaL                                  | `adal`                                   | `.adal/skills/`          | `~/.adal/skills/`               |

<!-- supported-agents:end -->

> [!NOTE]
> **Kiro CLI 用户：** 默认 Agent 会自动从 `.kiro/skills/` 与 `~/.kiro/skills/` 加载技能，**无需额外配置**。若使用**自定义 Agent**，请在其 `.kiro/agents/<agent>.json` 的 `resources` 中加入技能：
>
> ```json
> {
>   "resources": ["skill://.kiro/skills/**/SKILL.md"]
> }
> ```

CLI 会自动检测已安装的编程 Agent。若未检测到任何 Agent，会提示你选择要安装到的 Agent。

## 创建技能

技能是包含 `SKILL.md` 的目录，文件带 YAML frontmatter：

```markdown
---
name: my-skill
description: What this skill does and when to use it
---

# My Skill

Instructions for the agent to follow when this skill is activated.

## When to Use

Describe the scenarios where this skill should be used.

## Steps

1. First, do this
2. Then, do that
```

### 必填字段

- `name`：唯一标识（小写，可用连字符）
- `description`：技能用途的简要说明

### 可选字段

- `metadata.internal`：设为 `true` 可在常规发现中隐藏该技能。内部技能仅在设置 `INSTALL_INTERNAL_SKILLS=1` 时可见、可安装。适用于进行中的技能或仅内部工具使用的技能。

```markdown
---
name: my-internal-skill
description: An internal skill not shown by default
metadata:
  internal: true
---
```

### 技能发现

CLI 会在仓库内以下位置搜索技能：

<!-- skill-discovery:start -->

- 根目录（若包含 `SKILL.md`）
- `skills/`
- `skills/.curated/`
- `skills/.experimental/`
- `skills/.system/`
- `.aider-desk/skills/`
- `.agents/skills/`
- `.augment/skills/`
- `.bob/skills/`
- `.claude/skills/`
- `.codeartsdoer/skills/`
- `.codebuddy/skills/`
- `.codemaker/skills/`
- `.codestudio/skills/`
- `.commandcode/skills/`
- `.continue/skills/`
- `.cortex/skills/`
- `.crush/skills/`
- `.devin/skills/`
- `.factory/skills/`
- `.forge/skills/`
- `.goose/skills/`
- `.hermes/skills/`
- `.junie/skills/`
- `.iflow/skills/`
- `.kilocode/skills/`
- `.kiro/skills/`
- `.kode/skills/`
- `.mcpjam/skills/`
- `.vibe/skills/`
- `.mux/skills/`
- `.openhands/skills/`
- `.pi/skills/`
- `.qoder/skills/`
- `.qwen/skills/`
- `.rovodev/skills/`
- `.roo/skills/`
- `.tabnine/agent/skills/`
- `.trae/skills/`
- `.windsurf/skills/`
- `.zencoder/skills/`
- `.neovate/skills/`
- `.pochi/skills/`
- `.adal/skills/`
<!-- skill-discovery:end -->

### 插件清单发现

若存在 `.claude-plugin/marketplace.json` 或 `.claude-plugin/plugin.json`，其中声明的技能也会被一并发现：

```json
// .claude-plugin/marketplace.json
{
  "metadata": { "pluginRoot": "./plugins" },
  "plugins": [
    {
      "name": "my-plugin",
      "source": "my-plugin",
      "skills": ["./skills/review", "./skills/test"]
    }
  ]
}
```

可与 [Claude Code 插件市场](https://code.claude.com/docs/en/plugin-marketplaces) 生态兼容。

若在标准位置未找到技能，会执行递归搜索。

## 兼容性

技能通常可在各 Agent 间通用，因其遵循统一的 [Agent Skills 规范](https://agentskills.io)。部分功能可能仅特定 Agent 支持：

| 功能            | OpenCode | OpenHands | Claude Code | Cline | CodeBuddy | Codex | Command Code | Kiro CLI | Cursor | Antigravity | Roo Code | Github Copilot | Amp | OpenClaw | Neovate | Pi  | Qoder | Zencoder |
| --------------- | -------- | --------- | ----------- | ----- | --------- | ----- | ------------ | -------- | ------ | ----------- | -------- | -------------- | --- | -------- | ------- | --- | ----- | -------- |
| 基础技能        | 是       | 是        | 是          | 是    | 是        | 是    | 是           | 是       | 是     | 是          | 是       | 是             | 是  | 是       | 是      | 是  | 是    | 是       |
| `allowed-tools` | 是       | 是        | 是          | 是    | 是        | 是    | 是           | 否       | 是     | 是          | 是       | 是             | 是  | 是       | 是      | 是  | 是    | 否       |
| `context: fork` | 否       | 否        | 是          | 否    | 否        | 否    | 否           | 否       | 否     | 否          | 否       | 否             | 否  | 否       | 否      | 否  | 否    | 否       |
| Hooks           | 否       | 否        | 是          | 是    | 否        | 否    | 否           | 是       | 否     | 否          | 否       | 否             | 否  | 否       | 否      | 否  | 否    | 否       |

## 故障排除

### 「No skills found」（未找到技能）

请确认仓库中存在有效的 `SKILL.md`，且 frontmatter 中同时包含 `name` 与 `description`。

### Agent 未加载技能

- 确认技能已安装到正确路径
- 查阅该 Agent 文档中关于技能加载的说明
- 确认 `SKILL.md` 的 frontmatter 为合法 YAML

### 权限错误

请确认对目标目录有写入权限。

## 环境变量

| 变量                      | 说明                                                                       |
| ------------------------- | -------------------------------------------------------------------------- |
| `INSTALL_INTERNAL_SKILLS` | 设为 `1` 或 `true` 以显示并安装标记为 `internal: true` 的技能              |
| `DISABLE_TELEMETRY`       | 禁用匿名使用统计                                                           |
| `DO_NOT_TRACK`            | 另一种禁用统计的方式                                                       |

```bash
# 安装内部技能
INSTALL_INTERNAL_SKILLS=1 npx skills add vercel-labs/agent-skills --list
```

## 遥测

本 CLI 会收集匿名使用数据以改进工具，不收集个人信息。

在 CI 环境中会自动禁用遥测。

## 相关链接

- [Agent Skills 规范](https://agentskills.io)
- [技能目录](https://skills.sh)
- [Amp Skills 文档](https://ampcode.com/manual#agent-skills)
- [Antigravity Skills 文档](https://antigravity.google/docs/skills)
- [Factory AI / Droid Skills 文档](https://docs.factory.ai/cli/configuration/skills)
- [Claude Code Skills 文档](https://code.claude.com/docs/en/skills)
- [OpenClaw Skills 文档](https://docs.openclaw.ai/tools/skills)
- [Cline Skills 文档](https://docs.cline.bot/features/skills)
- [CodeBuddy Skills 文档](https://www.codebuddy.ai/docs/ide/Features/Skills)
- [Codex Skills 文档](https://developers.openai.com/codex/skills)
- [Command Code Skills 文档](https://commandcode.ai/docs/skills)
- [Crush Skills 文档](https://github.com/charmbracelet/crush?tab=readme-ov-file#agent-skills)
- [Cursor Skills 文档](https://cursor.com/docs/context/skills)
- [Firebender Skills 文档](https://docs.firebender.com/multi-agent/skills)
- [Gemini CLI Skills 文档](https://geminicli.com/docs/cli/skills/)
- [GitHub Copilot Agent Skills](https://docs.github.com/en/copilot/concepts/agents/about-agent-skills)
- [iFlow CLI Skills 文档](https://platform.iflow.cn/en/cli/examples/skill)
- [Kimi Code CLI Skills 文档](https://moonshotai.github.io/kimi-cli/en/customization/skills.html)
- [Kiro CLI Skills 文档](https://kiro.dev/docs/cli/custom-agents/configuration-reference/#skill-resources)
- [Kode Skills 文档](https://github.com/shareAI-lab/kode/blob/main/docs/skills.md)
- [OpenCode Skills 文档](https://opencode.ai/docs/skills)
- [Qwen Code Skills 文档](https://qwenlm.github.io/qwen-code-docs/en/users/features/skills/)
- [OpenHands Skills 文档](https://docs.openhands.ai/modules/usage/how-to/using-skills)
- [Pi Skills 文档](https://github.com/badlogic/pi-mono/blob/main/packages/coding-agent/docs/skills.md)
- [Qoder Skills 文档](https://docs.qoder.com/cli/Skills)
- [Replit Skills 文档](https://docs.replit.com/replitai/skills)
- [Roo Code Skills 文档](https://docs.roocode.com/features/skills)
- [Trae Skills 文档](https://docs.trae.ai/ide/skills)
- [Vercel Agent Skills 仓库](https://github.com/vercel-labs/agent-skills)

## 许可证

MIT
