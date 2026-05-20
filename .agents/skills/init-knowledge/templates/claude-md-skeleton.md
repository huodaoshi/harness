# {{PROJECT_NAME}}

> {{PROJECT_ONE_LINER}}

## 仓库结构

```
{{REPO_TREE}}
```

## 技术栈

{{TECH_STACK_TABLE}}

## 服务架构

{{SERVICE_ARCH}}

## 域清单

{{DOMAIN_LIST}}

## 常用命令

```bash
{{COMMON_COMMANDS}}
```

## 全局约定

- **禁止自动提交**：除非用户明确要求，不要执行 `git commit`
- **提交格式**：`<type>(<scope>): <subject>`，详见 `.cursor/rules/git.mdc` 或 `.claude/rules/git.md`（若存在）
- **文档语言**：人类可读文档默认**简体中文**（见 `.cursor/rules/docs-zh.mdc` 或等价规则）
- **Windows 终端**：勿用 `&&` 连接命令，见 `windows-shell` 规则
- 其他约定见 **`.cursor/rules/`**、**`.claude/rules/`**、**`.harness/knowledge/`** 与 **`.harness/workflow/`**

## Agent 工作流（harness）

### Issue 与分拣（若使用 `.scratch/`）

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
- **全量 / 首次**：在 **harness** 仓库根对 `$TARGET` 运行 **init-knowledge**（可传目标仓绝对路径）
- **项目索引**：**`.harness/knowledge/index.yaml`**；domain 划分变更后重跑 **init-knowledge**（如 `--domain <name>`）
- **规则包同步**（可选）：`skills rules add <harness-repo> -a cursor` 等，见 `rules/README.md`
