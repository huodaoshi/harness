---
name: init-knowledge
description: "在目标 git 仓库初始化或更新项目知识库：写入 .harness/knowledge/index.yaml 与 domain 文档等；默认下发 .harness/workflow/feature-workflow.md、session hook 与 session-bootstrap；可选双写规则与 CLAUDE.md。在 Cursor 中加载本技能后执行。"
---

# Init Knowledge — 项目知识库初始化

把消费端 **`$TARGET`** 的项目知识库初始化到 **`.harness/knowledge/`**（索引 + domain 文档 + 按域查询说明），**默认**写入 **`.harness/workflow/feature-workflow.md`** 与 **sessionStart hook**（`.harness/session/`、`.cursor/hooks*`、`.claude/hooks*`），并可选写入 **`.claude/` / `.cursor/`** 语言规则与 **`CLAUDE.md`**。

## 参数

| 参数 | 行为 |
|---|---|
| `<目标仓库绝对路径>`（可选位置参数） | 指定要扫描和写入的目标仓库；不传则用当前 `pwd` |
| 无参数 | 完整 init：跑 A+B+C+D（**D 为 hook**，默认开启），已存在文件保留（按需合并） |
| `--full` | 强制完全重跑 A+B+C+D（**含 D**；hook 脚本覆盖，settings/hooks.json 仍按合并规则） |
| `--domain <name>` | 仅对指定 domain 跑 **B+C**（该 domain 的深度学习 + 生成其 query 说明等）；**不跑 D** |
| `--commands` | 仅跑 A4+C4（重新推断验证命令并扩充 `permissions.allow`）；**不跑 D** |
| `--hooks` | 仅跑 **D**（重装 sessionStart hook 与 `.harness/session/session-bootstrap.md`） |
| `--workflow` | 仅跑 **C5**（重装 `.harness/workflow/feature-workflow.md`） |
| `--no-hooks` | 完整 init 时**跳过 D**（其余阶段不变） |

## 目标路径解析（$TARGET）

执行任何阶段前先确定 **`$TARGET`**：

1. 若 ARGUMENTS 中存在以 `/` 开头的位置参数，取该绝对路径作为 `$TARGET`
2. 否则 `$TARGET = $(pwd)`

校验：
- `$TARGET` 必须是已存在的目录，否则报错退出
- 若 `$TARGET != $(pwd)`，用 `git -C "$TARGET" rev-parse --show-toplevel` 确认是 git 仓库（非 git 仓库则警告但继续）

**约定**：本文档中后续 **`.harness/knowledge/...`**、**`.claude/...`**、**`.cursor/...`**、`CLAUDE.md`、`README.md`、`go.mod`、`package.json` 等**项目文件路径**均**相对 `$TARGET`** 解析。所有 Glob / Read / Write / Edit 必须用 `$TARGET/...` 绝对路径；Bash 用 `git -C "$TARGET" ...` 或 `cd "$TARGET" && ...`。

## Harness 根目录（`$HARNESS_ROOT`）

模板与 **`learner-workflow.md` 真源**在 **harness** 维护仓（本技能位于 `.agents/skills/init-knowledge/`），**不**从 `$TARGET` 读取。跨仓库 init 时须在 harness 根执行，或设置环境变量。

解析顺序：

1. 环境变量 **`HARNESS_ROOT`** 已设置且为存在目录 → 使用该路径
2. 否则若 **`$(pwd)/.agents/skills/init-knowledge/SKILL.md`** 存在 → **`$HARNESS_ROOT = $(pwd)`**
3. 否则退出并提示：在 harness 仓库根运行，或 `export HARNESS_ROOT=/path/to/harness`

| 资源 | 路径（相对 `$HARNESS_ROOT`） |
|------|------------------------------|
| query 模板 | `.agents/skills/init-knowledge/templates/_query-domain.md` |
| CLAUDE 骨架 | `.agents/skills/init-knowledge/templates/claude-md-skeleton.md` |
| Cursor 语言规则源 | `rules/cursor/_lang/<lang>.mdc`（优先）；否则 `.cursor/rules/_lang/<lang>.mdc` |
| Claude 语言规则源 | `rules/claude/_lang/<lang>.md`（优先）；否则 `.claude/rules/_lang/<lang>.md` |
| learner-workflow（维护真源） | `.harness/knowledge/learner-workflow.md` |
| learner-workflow（随技能安装） | `.agents/skills/init-knowledge/resources/learner-workflow.md` |
| session-bootstrap（随技能安装） | `.agents/skills/init-knowledge/resources/session/session-bootstrap.md` |
| feature-workflow（随技能安装） | `.agents/skills/init-knowledge/resources/workflow/feature-workflow.md` |
| sessionStart hook 包（随技能安装） | `.agents/skills/init-knowledge/resources/hooks/`（优先）；回退 `$HARNESS_ROOT/.cursor/hooks/`、`.claude/hooks/` |

从 harness 根对 **`$TARGET` 绝对路径`** 初始化消费端仓库时：业务路径写 `$TARGET`，模板读 `$HARNESS_ROOT`。  
**其它项目**只 `skills add` 了本技能时：无 `$HARNESS_ROOT` 也能靠 **resources/** 下发 workflow（见 A0、C5）与 hook（见 D）。

### 维护者：改 workflow 后同步两处

1. 编辑 **`.harness/knowledge/learner-workflow.md`**（真源）
2. 复制到 **`.agents/skills/init-knowledge/resources/learner-workflow.md`**（随 `skills add` 分发）
3. 已 init 过的消费端：可删其 `.harness/knowledge/learner-workflow.md` 后重跑 A0，或手动覆盖

### 维护者：改 feature-workflow 后同步 resources

1. 编辑 harness 真源：**`.harness/workflow/feature-workflow.md`**
2. 复制到 **`.agents/skills/init-knowledge/resources/workflow/feature-workflow.md`**
3. 已 init 过的消费端：重跑 **`init-knowledge --workflow`**（或 `--full`）

### 维护者：改 hook / session 后同步 resources

1. 编辑 harness 真源：**`.harness/session/session-bootstrap.md`**、**`.cursor/hooks/*`**、**`.claude/hooks/*`**
2. 同步到 **`.agents/skills/init-knowledge/resources/`**（`session/session-bootstrap.md`、`hooks/cursor/`、`hooks/claude/`）
3. 已 init 过的消费端：在目标仓重跑 **`init-knowledge --hooks`**（或 `--full`）

## 进度追踪

使用 **TodoWrite** 为「阶段 A / B / C / D」及各 domain 创建待办，随阶段更新勾选状态（`--hooks` 仅 D；`--commands` 无 D）。

## 代理注意（Cursor）

- 所有写入 **`$TARGET`** 的路径须用 **Read/Write/Edit** 的绝对路径。
- **模板 / 规则源 / learner-workflow 真源**仅从 **`$HARNESS_ROOT`** 读取（见上表），**从不**从 `$TARGET` 读模板。
- 阶段 B 前确保 **`$TARGET/.harness/knowledge/learner-workflow.md`** 存在（见 **A0**）；再 **Read** 该文件并生成四份 domain md（**不要求**子代理）。

### A0. 下发 learner-workflow（目标仓缺失时）

消费端运行时要读的是 **`$TARGET/.harness/knowledge/learner-workflow.md`**。若不存在，按顺序取**第一份可读来源**复制过去（父目录不存在则创建）：

| 优先级 | 来源（相对路径解析基准） |
|--------|--------------------------|
| 1 | **`$TARGET/.agents/skills/init-knowledge/resources/learner-workflow.md`**（`skills add` 已装本技能） |
| 2 | **`$HARNESS_ROOT/.harness/knowledge/learner-workflow.md`**（在 harness 根执行 init） |
| 3 | **`$HARNESS_ROOT/.agents/skills/init-knowledge/resources/learner-workflow.md`**（与 1 内容应一致） |

复制成功后日志：`[A0] 已下发 learner-workflow.md（来源: <路径>）`。  
三者皆不存在 → 报错退出，提示安装 **init-knowledge** 技能或设置 **HARNESS_ROOT**。

**说明**：`.harness/knowledge/` 通常**不会**随 rules/skills 同步进其它项目；靠 A0 在**首次 init/learn** 时落地即可，无需单独「同步 harness 目录」。

---

## 阶段 A：项目指纹识别（Detect）

由本 skill 主体完成，**不调用任何子 agent**。

### A1. 仓库形态判断

检查以下信号文件（用 Glob，存在即记录）：

- 工作区文件：`go.work`、`pnpm-workspace.yaml`、`lerna.json`、`turbo.json`、`Cargo.toml`（含 `[workspace]`）、`pnpm-lock.yaml` 顶层
- 顶层目录：`app/`、`apps/`、`packages/`、`services/`、`cmd/`、`src/`

判断规则：
- 出现工作区文件或 `app/<N>/`、`packages/<N>/`、`services/<N>/` 多子目录 → `repo_layout: monorepo`
- 单一根级源码 → `repo_layout: single`

### A2. 技术栈识别

按下表识别（可多语言并存，遍历整棵仓库的根目录与一级子目录）：

| 信号文件 | 推断栈 | 规则源（`$HARNESS_ROOT`） | `index.yaml` 中 rules 引用（相对 `$TARGET`） |
|---|---|---|---|
| `go.mod` | Go | `rules/cursor/_lang/go.mdc` 等 | `.cursor/rules/_lang/go.mdc`、`.claude/rules/_lang/go.md` |
| `package.json` 含 `"react"` | React | `react` + `typescript` | 同上两路径各语言各一条 |
| `package.json` 含 `"vue"` | Vue | `vue` + `typescript` | 同上 |
| `package.json` 含 `"typescript"` 且无 react/vue | TS | `typescript` | 同上 |
| `Cargo.toml` | Rust | `rust` | 同上 |
| `pyproject.toml` / `requirements.txt` / `setup.py` | Python | `python` | 同上 |
| `pom.xml` / `build.gradle*` | JVM | `jvm` | 同上 |
| `Gemfile` | Ruby | — | 无内置模板，跳过或提示用户自编 |

### A3. 域（Domain）划分

划分逻辑（由仓库形态决定）：

**monorepo**：
- 把 `app/<X>/`、`apps/<X>/`、`services/<X>/`、`packages/<X>/` 下每个直接子目录作为一个候选 domain
- 每个候选 domain 的 `lang` 取该目录下识别到的栈（如目录内同时含 `go.mod` 与 `package.json`，按主导文件类型判断）
- `type` 推断：
  - 含 `cmd/main.*`、`server.*`、`grpc`、`http/controller` → `backend`
  - 含 `src/pages/`、`src/components/`、`vite.config.*`、`webpack.config.*` → `frontend`
  - 含 `pkg/`、`packages/`、纯 lib 类型 → `shared`
  - 含 `Dockerfile`、`k8s/`、`terraform/`、`ansible/` 等 → `infra`

**single**：
- 把仓库本身视为一个 domain
- `name` 取仓库目录名
- `path` 设为 `./`
- `type` 按 A2 推断的栈与目录特征综合判断

**知识路径约定**：四件套目录为 **`.harness/knowledge/domains/<type>/`**（不含 `name` 段）。`index.yaml` 中每个 `type` 在磁盘上对应唯一目录；若 monorepo 存在多个同 `type` 的候选 domain，阶段 A 须合并为单一 domain 条目，或为各 domain 分配互斥的 `type`（如 `backend-account`）。

### A4. 验证命令推断

按检测到的栈，扫描根目录与每个 domain.path 下的构建文件：

| 文件 | 提取动作 |
|---|---|
| `Makefile` | 解析 target；提取 `build`、`test`、`lint`、`verify`、`check` 等可能命令 |
| `package.json` | 读取 `scripts.{build,test,lint,typecheck,check}`，组合为 `pnpm run <script>` 或 `npm run <script>`（按 lockfile 选择） |
| `Cargo.toml` | 默认 `cargo build`、`cargo test`、`cargo clippy --all-targets -- -D warnings` |
| `pyproject.toml` | 检测 `[tool.pytest.ini_options]` → `pytest`；`[tool.ruff]` → `ruff check .`；`[tool.mypy]` → `mypy .` |
| `pom.xml` | 默认 `mvn compile`、`mvn test`、`mvn verify` |
| `build.gradle*` | 默认 `gradle build`、`gradle test`、`gradle check` |

每个 domain 的 `validate_commands` 是一个有序列表：编译/类型检查 → lint → 单元测试。

如检测到的 monorepo 顶层 Makefile 含 `make build <service>` 模式（如 oneone 项目），优先用顶层命令；否则进入 domain.path 后执行（如 `cd <path> && pnpm run build`）

### A5. 产出

写入两份文件：

**`.harness/knowledge/index.yaml`** — 主索引：

```yaml
project_name: <从根 Makefile / package.json / Cargo.toml / 仓库目录名推断>
repo_layout: <monorepo | single>
detected_languages: [<lang1>, <lang2>, ...]
domains:
  - name: <name>
    type: <backend | frontend | shared | infra>
    lang: <go | typescript | react | vue | python | rust | jvm>
    path: <relative/path/>
    description: <一句话简介，从 README / 入口注释提取，可空>
    skills:
      - .harness/knowledge/domains/<type>/01-architecture.md
      - .harness/knowledge/domains/<type>/02-business-domains.md
      - .harness/knowledge/domains/<type>/03-infra-patterns.md
      - .harness/knowledge/domains/<type>/04-dev-guide.md
    rules:
      - .cursor/rules/_lang/<lang>.mdc
      - .claude/rules/_lang/<lang>.md
    validate_commands:
      - <command 1>
      - <command 2>
    test_glob: <optional, glob pattern for test files, e.g. "*_test.go" or "*.test.ts"; defaults to language-specific convention if omitted>
shared_skills: []
cross_domain_skills: []
```

**`.harness/knowledge/init-detection.json`** — 中间产物（人可读，用于 debug）：

```json
{
  "detected_at": "<ISO8601>",
  "signals": {
    "monorepo_files": ["go.work"],
    "lang_files": {
      "go.mod": ["./go.mod", "./app/backend/account/go.mod"],
      "package.json": ["./app/frontend/client/package.json"]
    }
  },
  "domain_candidates": [
    { "path": "app/backend/account", "evidence": ["cmd/main.go", "internal/service"] }
  ]
}
```

阶段 A 完成后，将「阶段 A」对应待办标记为完成，输出阶段总结：

```
[阶段 A 完成]
仓库形态: monorepo
检测到栈: [go, typescript, react]
划分 domain: 5 个
  - account (backend, go) — app/backend/account/
  - community (backend, go) — app/backend/community/
  - client (frontend, react) — app/frontend/client/
  - admin (frontend, react) — app/frontend/admin/
  - seo-server (frontend, typescript) — app/frontend/seo-server/
```

---

## 阶段 B：深度学习每个 Domain（Learn）

将「阶段 B」对应待办标记为进行中。

**串行**遍历 **`.harness/knowledge/index.yaml`** 中 `domains[]`，对每个 domain：

1. 为当前 domain 创建或更新子待办（如「学习 domain: <name>」），标记为进行中。
2. 确保目录存在：`$TARGET/.harness/knowledge/domains/<type>/`（**`<type>`** 与该 domain 的 `type` 一致）。
3. **在当前会话中执行**（**不要求**子代理 / Task）：先 **Read** **`$TARGET/.harness/knowledge/learner-workflow.md`**，再按其 **INIT** 节与下列参数，在 **`$TARGET/<domain.path>`** 下扫描源码，生成 4 份 skill 并写入上述目录：

```
模式：INIT
目标仓库根（$TARGET）: <绝对路径>
目标域: <name>
路径: <$TARGET>/<domain.path>
检测到的栈: <lang>
domain type: <type>

产出文件（与 learner-workflow 中 INIT 节一致）：
  - 01-architecture.md
  - 02-business-domains.md
  - 03-infra-patterns.md
  - 04-dev-guide.md

写入路径：$TARGET/.harness/knowledge/domains/<type>/0X-*.md
完成后输出 [LEARNER 完成] 摘要。
```

4. 完成后将该 domain 子待办标记为完成。

**串行的理由**：避免 context 爆炸；失败时易定位到具体 domain。

如使用了 `--domain <name>` 参数，仅对指定 domain 执行本阶段。

阶段 B 完成后，将「阶段 B」对应待办标记为完成。

---

## 阶段 C：生成项目特定资产（Generate）

将「阶段 C」对应待办标记为进行中。

### C1. 生成按域「知识查询」说明（Markdown，非 Agent）

为每个 domain 生成 **`$TARGET/.harness/knowledge/query/<domain.name>.md`**（**禁止**写入 **`.claude/agents/`**）。

读取 **`$HARNESS_ROOT/.agents/skills/init-knowledge/templates/_query-domain.md`**，按以下变量填充后写入上述路径：

| 占位符 | 来源 |
|---|---|
| `{{DOMAIN_NAME}}` | `domain.name` |
| `{{DOMAIN_PATH}}` | `domain.path` |
| `{{DOMAIN_TYPE}}` | `domain.type` |
| `{{DOMAIN_LANG}}` | `domain.lang` |
| `{{ONE_LINE_SUMMARY}}` | 从 `01-architecture.md` 第一段提取 |
| `{{PROJECT_NAME}}` | `index.yaml.project_name` |
| `{{SKILLS_LIST}}` | `domain.skills[]` 渲染为列表 |
| `{{SHARED_SKILLS_LIST}}` | `index.yaml.shared_skills[]` 渲染为列表；若为空，删除占位符及其上方的「涉及共享代码时，按需读取：」标题行 |

### C2. 复制栈特定 rules

对每个 domain 推断出的 **`<lang>`**（去重后）：

1. **读取源**（相对 `$HARNESS_ROOT`，按优先级）：
   - Cursor：`rules/cursor/_lang/<lang>.mdc`，不存在则 `.cursor/rules/_lang/<lang>.mdc`
   - Claude：`rules/claude/_lang/<lang>.md`，不存在则 `.claude/rules/_lang/<lang>.md`
2. **写入目标**（相对 `$TARGET`）：
   - **`$TARGET/.cursor/rules/_lang/<lang>.mdc`**
   - **`$TARGET/.claude/rules/_lang/<lang>.md`**
   - 父目录不存在则创建；**若目标文件已存在则跳过**（不覆盖），并在报告中注明 `skipped`
3. 写入时更新 frontmatter 的 **`paths`**：合并所有同 `lang` 的 `domain.path` 下该语言源文件的 glob（例如 `cli/**/*.ts`、`**/*.go`）

**harness 骨架消费端**若已有 `rules/` 规范源，也可在 init 后让用户用 `skills rules add` 同步，本阶段仅处理 **`_lang`** 语言规则。

### C3. 生成 CLAUDE.md

读取 **`$HARNESS_ROOT/.agents/skills/init-knowledge/templates/claude-md-skeleton.md`**，按以下变量填充并写入 `$TARGET/CLAUDE.md`：

| 占位符 | 来源 |
|---|---|
| `{{PROJECT_NAME}}` | `index.yaml.project_name` |
| `{{PROJECT_ONE_LINER}}` | 从根 README.md 第一行 / package.json 的 description / 用户输入 |
| `{{REPO_TREE}}` | 基于 `domain.path[]` 与扫描到的顶层目录生成 ASCII 目录树 |
| `{{TECH_STACK_TABLE}}` | 基于 `detected_languages` 与各 domain 的依赖文件版本生成表格 |
| `{{SERVICE_ARCH}}` | 基于 domain 间引用关系生成（无明确关系时省略本节） |
| `{{COMMON_COMMANDS}}` | 聚合所有 `validate_commands[]` |
| `{{DOMAIN_LIST}}` | `index.yaml.domains[]` 列表 |

骨架中 **「知识体系（三层）」「Agent 工作流」「自学习与索引」** 等与 harness 相关的段落原样保留；若消费端无 `.scratch/`、无 `docs/agents/`，可在生成后由用户删减无关小节。

如项目根已存在 CLAUDE.md：
- 默认行为：备份原文件为 `CLAUDE.md.bak.<timestamp>` 后写入新内容
- `--full` 模式下：直接覆盖，不备份（用户负责 git）

### C4. 扩充 settings.json

若 **`$TARGET/.claude/settings.json`** 不存在：可创建仅含 `permissions` 骨架的最小文件，或**跳过** C4（纯 Cursor 消费端无 `.claude/` 时跳过）。若文件存在，将 A4 推断到的命令前缀（如 `make *`、`go vet *`、`go test *`、`pnpm run *`、`cargo *`）追加到 `permissions.allow` 列表（去重）。

**不动 `permissions.deny`**；**`hooks` 字段由 D 合并**（C4 仅处理 `permissions.allow`）。

如果 `permissions.allow` 中已存在某条目，跳过不重复添加。

### C5. 下发 feature-workflow

**何时执行：**

- 完整 init、**`--domain`**（含阶段 C）、**`--workflow`** → 执行
- 仅 **`--workflow`** → **只执行本步**（仍须先解析 `$TARGET`）
- **`--commands`**、**`--hooks`** → **跳过**
- **`--no-hooks`** 不影响本步（完整 init 仍跑 C5）

**源路径**（禁止从 `$TARGET` 读模板）：

| 优先级 1 | 优先级 2（harness 根 init） |
|----------|---------------------------|
| `resources/workflow/feature-workflow.md` | `.harness/workflow/feature-workflow.md` |

二者皆不存在 → 报错退出。

**写入：**

- 目标：**`$TARGET/.harness/workflow/feature-workflow.md`**
- 父目录不存在则创建
- 已存在且非 `--full` / `--workflow` → **跳过**（报告 `skipped`）
- **`--full` 或 `--workflow`** → **覆盖**

**可选裁剪（阶段 A 未发现 `.scratch/` 且消费端无 `docs/agents/issue-tracker.md` 时）：**

- 下发后可用 Edit 删除文中「PRD / to-issues / triage / ready-for-agent」相关行，或在文首增加说明：「本仓库未启用本地 issue 跟踪，§ 新功能开发 ②③、§ 领 issue 实现 整节不适用。」
- 默认仍下发完整版，便于用户日后启用 `.scratch/`

日志：`[C5] feature-workflow.md: created | updated | skipped`

将「阶段 C」对应待办标记为完成。

---

## 阶段 D：下发 sessionStart hook（默认开启）

**何时执行：**

- 完整 init（无 `--no-hooks`、无 `--domain`、无 `--commands`）→ **必须执行**
- 仅 `--hooks` → **只执行本阶段**（仍须先解析 `$TARGET`）
- `--no-hooks`、`--domain`、`--commands` → **跳过**

**源路径解析**（相对 `$HARNESS_ROOT` 或已安装的 init-knowledge 技能包；**禁止**从 `$TARGET` 读模板）：

| 用途 | 优先级 1 | 优先级 2（harness 根 init 时） |
|------|----------|--------------------------------|
| `session-bootstrap.md` | `resources/session/session-bootstrap.md` | `.harness/session/session-bootstrap.md` |
| Cursor `hooks.json`、脚本 | `resources/hooks/cursor/` | `.cursor/hooks/` + 根 `.cursor/hooks.json` |
| Claude 脚本 | `resources/hooks/claude/` | `.claude/hooks/` |
| Claude `hooks` 片段 | `resources/hooks/claude/settings-hooks.fragment.json` | 从 `$HARNESS_ROOT/.claude/settings.json` 提取 `hooks` 对象 |

技能包 resources 与 HARNESS_ROOT 回退皆不可用 → 报错退出，提示设置 **HARNESS_ROOT** 或 `skills add` **init-knowledge**。

### D1. 写入 session-bootstrap

- 目标：**`$TARGET/.harness/session/session-bootstrap.md`**
- 父目录不存在则创建
- 已存在且非 `--full` → **跳过**（报告 `skipped`）
- `--full` 或 `--hooks` 且源更新 → **覆盖**

### D2. Cursor hook

| 源 | 目标 |
|----|------|
| `hooks/cursor/hooks.json` | `$TARGET/.cursor/hooks.json` |
| `hooks/cursor/run-hook.cmd` | `$TARGET/.cursor/hooks/run-hook.cmd` |
| `hooks/cursor/session-start` | `$TARGET/.cursor/hooks/session-start` |

- **`hooks.json`**：若不存在则写入；若已存在但无 `hooks.sessionStart` 数组 → **合并**该键（保留其它 hook 事件）；若已有 `sessionStart` 且命令指向本仓 `run-hook.cmd session-start` → 跳过；若已有不同 `sessionStart` → 跳过并在报告注明 `conflict-skipped`
- **脚本**：目标不存在 → 写入；已存在 → 默认跳过；**`--full` 或 `--hooks`** → 覆盖

### D3. Claude hook

| 源 | 目标 |
|----|------|
| `hooks/claude/run-hook.cmd` | `$TARGET/.claude/hooks/run-hook.cmd` |
| `hooks/claude/session-start` | `$TARGET/.claude/hooks/session-start` |

脚本覆盖规则同 D2。

**`$TARGET/.claude/settings.json`：**

1. 读取 **`settings-hooks.fragment.json`**（或 harness 真源中的 `hooks` 对象）
2. 若 **`settings.json` 不存在**：写入仅含 `hooks` 的 JSON（**不要**在此步写入 `permissions`；若 C4 已写入 `permissions` 则保留）
3. 若已存在：
   - 无 `hooks.SessionStart`（或为空）→ **合并** `hooks` 对象（保留 `permissions` 与其它顶层键）
   - 已有 `SessionStart` 且已注册 `run-hook.cmd session-start` → 跳过
   - 已有不同 `SessionStart` 配置 → 跳过，报告 `conflict-skipped`
4. **禁止**删除或覆盖 `permissions.deny`、`permissions.allow`（C4 专责）

纯 Cursor 消费端无 `.claude/` 时：可跳过 D3（报告 `claude-skipped`）。

### D4. 日志

```
[D] sessionStart hook
  - .harness/session/session-bootstrap.md: created | updated | skipped
  - .cursor/hooks.json: created | merged | skipped | conflict-skipped
  - .cursor/hooks/*: created | updated (N files) | skipped
  - .claude/hooks/*: ...
  - .claude/settings.json hooks: merged | skipped | conflict-skipped | claude-skipped
```

**注意：** Windows 需 **bash**（Git Bash）；无 bash 时 hook **fail-open**，新会话可能无注入——在报告中提示安装 Git for Windows。

阶段 D 完成后，将子待办「阶段 D」标记为完成（若使用 TodoWrite）。

---

## 最终报告（C5）

```
[INIT-KNOWLEDGE 完成]

项目: <project_name>
仓库形态: <monorepo | single>
检测到 <N> 个 domain：
  - <name> (<type>, <lang>)        → 4 份 domain md + 1 份 query 说明
  - ...

生成的文件：
  - CLAUDE.md (新建 / 更新，原文件已备份为 CLAUDE.md.bak.<ts>)
  - .harness/knowledge/index.yaml
  - .harness/knowledge/init-detection.json
  - .harness/knowledge/query/<N> 份按域查询说明
  - .harness/knowledge/domains/<type>/0X-*.md  (共 <4N> 份)
  - .cursor/rules/_lang/*.mdc 与 .claude/rules/_lang/*.md 栈规则（若写入）
  - .claude/settings.json (扩充 permissions.allow，若存在)
  - .harness/workflow/feature-workflow.md（C5）
  - .harness/session/session-bootstrap.md（D）
  - .cursor/hooks.json、.cursor/hooks/*（D）
  - .claude/hooks/*、.claude/settings.json 的 hooks.SessionStart（D，若适用）

下一步：
  1. 检查 git diff .harness/knowledge/ .cursor/rules/ .claude/ CLAUDE.md 确认生成内容
  2. 必要时手动修正 .harness/knowledge/index.yaml 中错误的域划分
  3. 重跑特定 domain：init-knowledge --domain <name>
  4. 仅更新 hook：init-knowledge --hooks
  5. 仅更新功能工作流：init-knowledge --workflow
  6. 新开 Agent 会话验证 sessionStart 是否注入 bootstrap
  7. 做功能前先 Read .harness/workflow/feature-workflow.md 并选对章节
  8. 全部满意后提交（示例）：
     git add .harness/knowledge .harness/workflow .harness/session .cursor/rules .claude .cursor/hooks CLAUDE.md
     git commit -m "docs(knowledge): 初始化项目知识库与 session hooks"
```

将「阶段 C」对应待办标记为完成。

---

## 失败处理

任何阶段失败时：
- 已完成的产物保留在磁盘
- 标记当前阶段对应待办为完成（避免卡住），同时在描述中追加 `[FAILED]`
- 输出错误摘要 + 已完成进度 + 重跑建议（用 `--domain` 重跑特定 domain、`--commands` 仅重跑命令推断、`--hooks` 仅重跑 hook、`--workflow` 仅重跑功能工作流）

## 重跑语义

| 场景 | 推荐 |
|---|---|
| 首次初始化 | 调用 **init-knowledge** 技能（无额外参数，`$TARGET` 为当前仓库根） |
| 大改后完全重建 | **init-knowledge** + `--full` |
| 单个 domain 知识过时 | **init-knowledge** + `--domain <name>` |
| 项目构建命令变了 | **init-knowledge** + `--commands` |
| harness 升级后同步 hook | **init-knowledge** + `--hooks` |
| 只更新功能工作流 | **init-knowledge** + `--workflow` |
| 不要 hook（纯 CI/无 IDE） | **init-knowledge** + `--no-hooks` |
| 日常增量学习 | **learn** 技能（不走 init；不修改 hook） |

## 约束

- 读 / 写 **`$TARGET/.harness/knowledge/**`**、**`$TARGET/.harness/workflow/**`**、**`$TARGET/.harness/session/**`**、**`$TARGET/CLAUDE.md`**（及其备份）；可选读 / 写 **`$TARGET/.claude/settings.json`**（C4：`permissions.allow`；D：合并 `hooks`）、**`$TARGET/.claude/hooks/**`**、**`$TARGET/.cursor/hooks.json`**、**`$TARGET/.cursor/hooks/**`**、**`$TARGET/.claude/rules/_lang/**`**、**`$TARGET/.cursor/rules/_lang/**`**
- 只读 **`$HARNESS_ROOT`** 下模板、规则源、**resources/workflow**、**resources/hooks**（见「Harness 根目录」表）；可写 **`$TARGET/.harness/knowledge/learner-workflow.md`**（A0 复制）
- 不修改 `$TARGET` 的业务源代码
- 不修改 `$TARGET/.git/`
- 不删除已有 **`$TARGET/.harness/knowledge/`** 下文件（除 `--full` 模式且用户已确认覆盖策略）
- 不引入任何外部运行时依赖（仅用 Read / Write / Edit / Glob / Grep / Bash；进度用 **TodoWrite** 或宿主等价能力）
- 不调用网络

## 从旧路径迁移（一次性）

若 **`$TARGET`** 仍存在 **`.claude/index.yaml`** 或 **`.claude/skills/<type>/<name>/`** 下旧四件套：可读旧内容 → 写入 **`.harness/knowledge/`** 新路径 → 更新 **`domains[].skills[]`** → 删除或归档旧路径（与业务源码目录勿混淆）。
