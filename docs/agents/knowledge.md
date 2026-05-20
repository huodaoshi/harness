# 项目知识库（`.harness/knowledge/`）

harness 用 **`.harness/knowledge/`** 存放可被 Agent 渐进读取的项目知识：索引、按域四件套、查询说明与学习规程。

## 技能

| 技能 | 路径 | 用途 |
|------|------|------|
| **init-knowledge** | `.agents/skills/init-knowledge/` | 首次初始化或重建索引、domain 文档、query、`CLAUDE.md`、语言规则；**默认**下发 sessionStart hook 与 `session-bootstrap.md`（`--no-hooks` 可跳过，`--hooks` 仅重装 hook） |
| **learn** | `.agents/skills/learn/` | 开发后按 git 变更增量更新 domain 四件套 |

二者运行时读 **`$TARGET/.harness/knowledge/learner-workflow.md`**。该文件**不必**随仓库模板整包同步——**首次 init/learn** 时由技能 **A0** 从已安装的 **`init-knowledge` 技能包**（`resources/learner-workflow.md`）或 harness 真源复制落地。

### `learner-workflow.md` 三份拷贝（别混）

| 位置 | 角色 |
|------|------|
| `.harness/knowledge/learner-workflow.md`（harness 仓） | **维护真源**，改流程只改这里 |
| `.agents/skills/init-knowledge/resources/learner-workflow.md` | **随 `skills add` 分发**，与其它项目同步的是技能，不是 `.harness/` 目录 |
| `<消费端>/.harness/knowledge/learner-workflow.md` | **运行时副本**，Agent 实际 Read 的路径 |

改真源后请执行（PowerShell 示例）：

```powershell
Copy-Item .harness\knowledge\learner-workflow.md .agents\skills\init-knowledge\resources\learner-workflow.md
```

## 在 harness 根执行（本仓）

```powershell
# 初始化/更新本仓知识库（$TARGET = 当前目录）
# 在 Cursor 中加载 init-knowledge 技能后按 SKILL.md 执行

# 开发后增量学习
# 加载 learn 技能后执行（需已有 index.yaml）
```

跨仓库：在 **harness 根**对其它项目 init 时，传入目标仓绝对路径，并确保 **`$HARNESS_ROOT`** 指向 harness（默认 `pwd` 即为 harness 根）。

## 目录约定

| 路径 | 维护者 |
|------|--------|
| `index.yaml` | init-knowledge（A 阶段） |
| `init-detection.json` | init-knowledge（调试用） |
| `learner-workflow.md` | harness 真源；A0 复制到消费端 |
| `domains/<type>/01-…04-*.md` | init（INIT）/ learn（INCREMENTAL、FULL） |
| `query/<domain>.md` | init-knowledge（C1） |

## 模板与规则源（仅 init 读取）

- 模板：`.agents/skills/init-knowledge/templates/`
- Session / Hook：`resources/session/session-bootstrap.md` → `.harness/session/`；`resources/hooks/` → `.cursor/hooks*`、`.claude/hooks*`（阶段 **D**）
- 规则：`rules/cursor/_lang/*.mdc`、`rules/claude/_lang/*.md`（回退 `.cursor/rules/_lang/`）

## 与 `rules/`、`skills/` 的关系

- **`rules/`**：Cursor/Claude **编辑器规则**规范源；可用 `skills rules add` 安装到消费端
- **`skills/`**：**`skills add`** 可安装的技能包副本
- **`.harness/knowledge/`**：**架构与业务知识**（供 Agent 深读），与上述二者互补
