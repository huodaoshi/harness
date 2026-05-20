Status: ready-for-agent

# Agent 规则目标路径与 CLI 参数（-a / --to）

## 要构建什么

在 `skills` CLI 的 Agent 定义中，为支持规则的 Agent 增加**项目内默认规则目录**（例如 Cursor → `.cursor/rules`，Claude Code → `.claude/rules`，具体键名与 `agents.ts` 现有结构一致）。解析 `skills rules add` 时：

- **必须**要求 `-a` 且 **MVP 仅允许单值**（与多次调用策略一致）。
- 支持 **`--to <path>`**：安装源为 `rules/<端>/` 时，目标为该路径下的**文件树内容**（与 issue 01 文档一致）。
- 未提供 `--to` 时，使用当前项目解析根下的默认规则目录（行为与现有 `skills add` 找项目根的方式对齐）。

本 issue **不要求**完成完整 `add` 克隆逻辑，但路径解析与**错误信息**（缺 `rules/cursor` 等）应可测或可演示。

## 验收标准

- [ ] 默认路径与 `--to` 覆盖在代码或单元测试中可验证。
- [ ] 未知或省略 `-a` 时有明确报错。
- [ ] 与 issue 01 文档一致；若有分歧，先改文档再合代码。

## 阻塞于

- `.scratch/skills-cli-rules/issues/01-spec-rules-package-and-lock-schema.md`（至少 schema 与路径约定已定稿）。
