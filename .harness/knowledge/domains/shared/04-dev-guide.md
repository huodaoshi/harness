# harness — 开发指南

## 新增 CLI 能力

1. 在 `cli/src/` 增加模块，于 `cli.ts` 注册子命令
2. 若涉及安装路径，复用 `project-root.ts` 的 `resolveCliProjectRoot`
3. 在 `cli/tests/` 增加 vitest 用例
4. 更新 `cli/README.md`、`docs/cli-rules.md`（若影响 rules/skills 契约）

## 新增或修改 Agent 技能

1. 在 **`.agents/skills/<name>/SKILL.md`** 撰写（简体中文，本仓维护）
2. 可选：同步副本到 **`skills/<name>/`** 供 `skills add` 分发
3. 大改流程时更新 **init-knowledge** / **learn** 相关说明

## 修改编辑器规则

1. 只改 **`rules/cursor/`** 或 **`rules/claude/`**
2. 需要 IDE 生效时：`pnpm --dir cli dev rules add . -a cursor`
3. 保持 `.cursor` 与 `.claude` 镜像一致（团队约定）

## 验证命令（本 domain）

```powershell
pnpm --dir cli exec vitest run
pnpm --dir cli run type-check
```

## 知识库与 session 维护

| 操作 | 方式 |
|------|------|
| 首次/重建知识库 | **init-knowledge**（A+B+C+D，D 默认下发 hook） |
| 开发后增量 | **learn**（依赖 `index.yaml` + `learner-workflow.md`） |
| 改学习规程 | 编辑 `.harness/knowledge/learner-workflow.md` → 复制到 `init-knowledge/resources/` |
| 改 session / hook | 编辑 `.harness/session/`、`.cursor/hooks/` → 同步 `resources/` → `init-knowledge --hooks` |

## 错误处理约定

- CLI 用户可见错误用 `@clack/prompts` 输出；致命错误 `process.exit(1)`
- Git 克隆失败：`GitCloneError`（`git.ts`）

## 提交

- 不自动 `git commit`；用户明确要求再提交
- 消息格式见 `.cursor/rules/git.mdc`

## 调试 CLI

```powershell
cd D:\harness
pnpm --dir cli dev add . --list
pnpm --dir cli dev rules add . -a cursor
```
