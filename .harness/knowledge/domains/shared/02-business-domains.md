# harness — 业务域（骨架能力）

本仓库无传统业务 API，「域」指 Agent 工作流与工具链能力分区。

## 核心能力模块

| 模块 | 路径 | 说明 |
|------|------|------|
| Skills CLI | `cli/` | 安装/列出/更新技能与规则；锁文件与跨平台路径 |
| 团队技能 | `.agents/skills/*/` | init-knowledge、learn、triage、to-prd、to-issues 等 |
| 可分发技能副本 | `skills/*/` | 与 `.agents/skills` 同步，供 `skills add` |
| 规则规范源 | `rules/cursor/`、`rules/claude/` | 语言与项目约定；`rules add` 写入消费端 |
| 项目知识库 | `.harness/knowledge/` | index、domain 四件套、query、learner-workflow |
| 会话引导 | `.harness/session/` | session-bootstrap，由 sessionStart hook 注入 |
| Issue 分拣 | `.scratch/<feature>/` | 本地 markdown issue（可选） |
| 协作文档 | `docs/agents/` | issue 跟踪器、分拣标签、知识库说明 |

## Agent 技能清单（.agents/skills）

约 18 个技能目录，含：init-knowledge、learn、caveman、diagnose、triage、to-issues、to-prd、setup-matt-pocock-skills、grill-me、grill-with-docs、handoff、improve-codebase-architecture、landscape、product-council、prototype、tdd、write-a-skill、zoom-out 等。

## 关键流程

1. **消费端初始化**：`init-knowledge` → A 指纹 → B 学 domain → C 生成资产 → D 下发 hook
2. **规则/技能安装**：`pnpm --dir cli dev add <pkg>` / `rules add <harness> -a cursor`
3. **开发后沉淀**：`learn` 按 git diff 增量更新 domain 知识
