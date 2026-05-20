---
name: learn
description: "在开发完成后增量或全量更新 .harness/knowledge 下 domain 知识文档；依赖 index.yaml 与 learner-workflow.md。在 Cursor 中使用。"
---

# 知识库学习工作流（learn）

分析最近的代码变更，将结构性知识更新到 **`.harness/knowledge/domains/`** 下各 domain 四件套，以及（若存在）**`.cursor/rules/_lang/`**、**`.claude/rules/_lang/`** 与项目根 **`CLAUDE.md`**。不修改业务源代码。

## 目标路径（`$TARGET`）

与 **init-knowledge** 相同：

1. 参数中有以 `/` 或 `D:\` / `C:\` 等开头的绝对路径 → 作为 `$TARGET`
2. 否则 `$TARGET` = 当前工作目录（应对应**消费端项目根**，而非 `cli/` 子目录）

所有 **`$TARGET/.harness/...`** 路径用绝对路径读写。

## 前置条件

### 1. `learner-workflow.md`

须存在 **`$TARGET/.harness/knowledge/learner-workflow.md`**。

若缺失，**先按 init-knowledge 的 A0 相同顺序自动下发**（无需子代理）：

| 优先级 | 复制来源 |
|--------|----------|
| 1 | `$TARGET/.agents/skills/init-knowledge/resources/learner-workflow.md` |
| 2 | `$HARNESS_ROOT/.harness/knowledge/learner-workflow.md`（需解析 `HARNESS_ROOT`，规则同 init-knowledge） |
| 3 | `$HARNESS_ROOT/.agents/skills/init-knowledge/resources/learner-workflow.md` |

→ 写入 **`$TARGET/.harness/knowledge/learner-workflow.md`**。成功则继续；否则：

```
[LEARN ERROR] 无法下发 learner-workflow.md
请在该项目执行: skills add <harness> --skill init-knowledge -y
或安装后重试 learn；亦可手动复制 resources/learner-workflow.md → .harness/knowledge/
```

**只装了 learn、未装 init-knowledge 时**：learn 无法单独提供 workflow，须先安装 **init-knowledge** 技能（workflow 打在该技能包里）。

### 2. `index.yaml`

须存在 **`$TARGET/.harness/knowledge/index.yaml`**。若缺失：

```
[LEARN ERROR] 未找到 .harness/knowledge/index.yaml
请先在 $TARGET 运行 init-knowledge 初始化知识库后再使用 learn。
```

## 参数解析

- **无参数**：增量模式（**INCREMENTAL**），学习自上次知识库相关文件提交以来的变更
- **domain 名**（如 `harness`）：**FULL** 模式仅针对该 domain（必须存在于 `index.yaml` 的 `domains[].name`）
- **`--full`**：全量模式，忽略 git 基准点，重建 **所有** domain 的知识文件

## 执行步骤

建议用 **TodoWrite** 标记「前置检查 → 执行 learner 流程 → 展示报告」三步。

### 阶段 1：按 learner-workflow 执行

**在当前会话中执行**（**不要求**子代理）。先 **Read** **`$TARGET/.harness/knowledge/learner-workflow.md`**，再按其中对应模式与下列 prompt 继续。

**增量模式（默认）**：

```
模式：INCREMENTAL
目标仓库根: <绝对路径 $TARGET>

请分析自上次学习以来的代码变更，更新知识库中的结构性信息。

步骤：
1. 通过 git log 推断基准点（见 learner-workflow；含 .harness/knowledge/domains/、query/、.cursor/rules/、.claude/rules/、CLAUDE.md）
2. 读取 $TARGET/.harness/knowledge/index.yaml 解析 domain 划分
3. 分析 git diff 中各 domain.path 下的变更
4. 对有结构性变更的 domain，更新 index.yaml.skills[] 所指向的四份 md
5. 检查 cross-cutting 规则与 CLAUDE.md 是否需要更新（见 learner-workflow）
```

**指定 domain（单域全量）**：

```
模式：FULL
目标仓库根: <绝对路径 $TARGET>
目标域: <DOMAIN>

请扫描指定 domain 的当前代码状态，更新对应的知识文件。

步骤：
1. 在 index.yaml 中查找 domains[].name == <DOMAIN>
2. 扫描 domain.path 下的完整结构
3. 对比现有四份 domain md
4. 更新或重建 01-architecture … 04-dev-guide
```

**`--full`（全库）**：

```
模式：FULL
目标仓库根: <绝对路径 $TARGET>

请全量扫描 index.yaml 中的所有 domain，重建知识库。

步骤：
1. 遍历 domains[] 中每个 domain
2. 对每个 domain 全量扫描 domain.path 并对比现有知识文件
3. 更新所有过时内容、补充缺失内容
4. 检查 cross-cutting 规则与 CLAUDE.md

注意：全量模式耗时较长，请逐个 domain 处理并输出中间进度。
```

### 阶段 2：展示报告

将 **learner-workflow** 执行结果中的 **`[LEARNER 完成]`** 或 **`[LEARNER SKIP]`** 全文展示给用户。

**完成时建议下一步**：

```
- git diff .harness/knowledge/ .cursor/rules/ .claude/rules/ CLAUDE.md
- 确认后可提交：docs(knowledge): 更新知识库
```

## 代理注意（Cursor）

- **Read** **`$TARGET/.harness/knowledge/learner-workflow.md`** 后再构造后续步骤；无需子代理。
- Git 命令对消费端仓库使用 **`git -C "$TARGET" ...`**（单条命令，勿用 `&&` 拼脚本，见 `windows-shell` 规则）。

## 约束

- 遵循 **`learner-workflow.md`** 中的硬性约束（不修改 **index.yaml** / **learner-workflow.md** / **init-detection.json** / **query/** 等由 init 维护的文件，除非该文档明确允许）
- 不修改业务源代码
- 全量模式仅在知识库严重过时时使用

## 与 init-knowledge 的分工

| 场景 | 使用 |
|------|------|
| 首次建库、改 domain 划分、重生成 query / CLAUDE 骨架 | **init-knowledge** |
| 日常开发后的知识增量 | **learn**（本技能） |
