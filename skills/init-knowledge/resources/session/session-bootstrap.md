# Harness 会话引导（由 sessionStart hook 注入）

本仓库的 Agent 技能真源在 **`.agents/skills/<name>/SKILL.md`**（简体中文）。通过 Cursor / Claude 的 Skill 工具加载时，也应优先使用该路径，而非已废弃的 Superpowers 全局目录。

## 默认沟通模式：caveman

**本会话默认启用 caveman**（极简沟通，省 token，技术内容保持准确）。

1. 开工前先 **Read** **`.agents/skills/caveman/SKILL.md`**
2. **每条回复**按 caveman 规则，直至用户明确说「stop caveman」「normal mode」「退出 caveman」「正常模式」
3. 安全警告、不可逆操作确认、多步易误解顺序、用户要求澄清时，可**暂时**不用 caveman，说清楚后恢复

## 优先阅读

- 项目总览：**`CLAUDE.md`**
- 知识库索引：**`.harness/knowledge/index.yaml`**（若存在）
- 按域查询：**`.harness/knowledge/query/<domain>.md`**
- 功能/重构流程：**`.harness/workflow/feature-workflow.md`**（做需求相关工作时）

## 常用技能（按需加载，勿臆造）

| 场景 | 技能 |
|------|------|
| **默认沟通（本会话已开）** | **caveman**（见上；关闭见 caveman 技能「持久性」） |
| 首次建立 / 重建项目知识库 | **init-knowledge** |
| 开发后更新知识文档 | **learn** |
| Issue 分拣 | **triage** |
| 计划拆 issue | **to-issues** |
| 写 PRD | **to-prd** |

完整列表以 **`.agents/skills/*/SKILL.md`** 为准；对外安装副本在 **`skills/`**（`skills add` 用）。

## 规则与 CLI

- 编辑器规则：**`.cursor/rules/`**（`*.mdc`）、**`.claude/rules/`**（`*.md`）
- 规则/技能规范源：**`rules/`**、**`skills/`**
- CLI：`pnpm --dir cli dev …`（安装技能/规则时注意 **`--cwd`** 指向目标项目根）

## 纪律

- 有 1% 可能适用某技能时，先读取对应 **`SKILL.md`** 再动手。
- 除非用户明确要求，不要执行 **`git commit`**。
- Windows 终端勿用 **`&&`** 连接命令（见 **windows-shell** 规则）。
