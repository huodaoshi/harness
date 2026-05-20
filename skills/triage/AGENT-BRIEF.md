# 撰写 Agent 简报

Agent 简报是在 issue 进入 `ready-for-agent` 时发布在 GitHub issue 上的结构化评论。它是 AFK Agent 工作的权威规格。原始 issue 正文与讨论是背景——Agent 简报才是契约。

## 原则

### 耐久性优于精确路径

issue 可能在 `ready-for-agent` 停留数天或数周，期间代码库会变。撰写简报时应使其在文件重命名、移动或重构后仍有用。

- **应**描述接口、类型与行为契约
- **应**点明 Agent 需查找或修改的具体类型、函数签名或配置形状
- **勿**引用文件路径——会过时
- **勿**引用行号
- **勿**假设当前实现结构保持不变

### 描述行为，而非步骤

说明系统**应做什么**，而非**如何实现**。Agent 会重新探索代码库并自行做实现决策。

- **好：**「`SkillConfig` 类型应接受可选的 `schedule` 字段，类型为 `CronExpression`」
- **差：**「打开 src/types/skill.ts 在第 42 行加 schedule 字段」
- **好：**「用户无参数运行 `/triage` 时，应看到待关注 issue 的摘要」
- **差：**「在主处理函数里加 switch」

### 完整的验收标准

Agent 需知何时算完成。每条 Agent 简报须有具体、可测的验收标准，且每条可独立验证。

- **好：**「运行 `gh issue list --label needs-triage` 返回已完成初部分类的 issue」
- **差：**「分拣应正常工作」

### 明确的范围边界

写明范围外内容，避免 Agent 镀金或对相邻特性做假设。

## 模板

```markdown
## Agent Brief

**Category:** bug / enhancement
**Summary:** 一行说明需要完成什么

**Current behavior:**
描述当前行为。对 bug，即损坏行为；对 enhancement，即特性所基于的现状。

**Desired behavior:**
描述 Agent 完成后应有的行为。明确边界情况与错误条件。

**Key interfaces:**
- `TypeName` — 需改什么及原因
- `functionName()` 返回类型 — 当前返回值 vs 应有返回值
- Config shape — 所需新配置项

**Acceptance criteria:**
- [ ] 具体、可测的标准 1
- [ ] 具体、可测的标准 2
- [ ] 具体、可测的标准 3

**Out of scope:**
- 本 issue 不应改动或处理的内容
- 看似相关但属另一特性的相邻功能
```

## 示例

### 好的 Agent 简报（bug）

```markdown
## Agent Brief

**Category:** bug
**Summary:** 技能 description 截断在词中间断开，输出损坏

**Current behavior:**
当技能 description 超过 1024 字符时，在恰好 1024 处截断，不考虑词边界，导致 description 在词中间结束（如 "Use when the user wants to confi"）。

**Desired behavior:**
截断应在 1024 字符前最后一个词边界处断开，并追加 "..." 表示已截断。

**Key interfaces:**
- `SkillMetadata` 的 `description` 字段 — 类型无需改，但填充/校验逻辑须尊重词边界
- 读取 SKILL.md frontmatter 并提取 description 的任意函数

**Acceptance criteria:**
- [ ] 不足 1024 字符的 description 不变
- [ ] 超过 1024 字符的在最后词边界处截断
- [ ] 截断后以 "..." 结尾
- [ ] 含 "..." 的总长度不超过 1024 字符

**Out of scope:**
- 修改 1024 字符上限本身
- 多行 description 支持
```

### 好的 Agent 简报（enhancement）

```markdown
## Agent Brief

**Category:** enhancement
**Summary:** 支持 `.out-of-scope/` 目录，记录已拒绝的功能请求

**Current behavior:**
功能请求被拒绝时，issue 以 `wontfix` 关闭并附评论，无决策与理由的持久记录。日后类似请求需维护者回忆或搜索旧讨论。

**Desired behavior:**
已拒绝的功能请求应记录在 `.out-of-scope/<concept>.md`，含决策、理由及所有相关 issue 链接。分拣新 issue 时应检查是否匹配。

**Key interfaces:**
- `.out-of-scope/` 下 Markdown 格式 — 每文件含 `# Concept Name` 标题、`**Decision:**`、`**Reason:**`、`**Prior requests:**` 及 issue 链接列表
- 分拣流程应尽早读取全部 `.out-of-scope/*.md`，按概念相似度匹配新进 issue

**Acceptance criteria:**
- [ ] 以 wontfix 关闭 enhancement 时创建/更新 `.out-of-scope/` 中的文件
- [ ] 文件含决策、理由及已关闭 issue 链接
- [ ] 若已有匹配的 `.out-of-scope/` 文件，将新 issue 追加到其 Prior requests，而非新建重复文件
- [ ] 分拣时检查 `.out-of-scope/` 并在新 issue 匹配既往驳回时告知维护者

**Out of scope:**
- 自动匹配（由人确认）
- 重新打开既往已拒绝功能
- Bug 报告（仅 enhancement 驳回写入 `.out-of-scope/`）
```

### 差的 Agent 简报

```markdown
## Agent Brief

**Summary:** 修分拣 bug

**What to do:**
分拣坏了。看主文件修一下，大约 150 行附近的函数有问题。

**Files to change:**
- src/triage/handler.ts (line 150)
- src/types.ts (line 42)
```

差的原因：
- 无 Category
- 描述模糊（「分拣坏了」）
- 引用会过时的路径与行号
- 无验收标准
- 无范围边界
- 未区分当前 vs 期望行为
