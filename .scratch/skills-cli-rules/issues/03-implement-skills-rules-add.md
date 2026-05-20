Status: ready-for-agent

# 实现 skills rules add

## 要构建什么

实现子命令 **`skills rules add <source>`**（命名与根 CLI 路由以现有 `cli/src/cli.ts` 为准）：从 Git 或本地路径取得包内容，仅将 **`rules/cursor/` 或 `rules/claude/`**（由 `-a` 决定）同步到目标目录。

行为与现有 **`skills add`** 对齐的部分：**clone/缓存策略**、默认 **symlink**、**`--copy`** 复制、**`-g` 全局**若适用则沿用同一语义（若全局规则目录与项目级不同，在实现中明确并在 issue 01 文档补一句）。

缺失对应子目录时：**清晰错误**，不静默空装。

## 验收标准

- [ ] 针对典型包（可 fixture 小型本地目录）有自动化测试覆盖「拷对子树」「`--to` 铺在目标下」。
- [ ] 与 issue 02 的路径解析集成。
- [ ] README 或 docs 中有端到端示例命令（可与 issue 05 合并文档改动，但本 issue 合并时应可运行）。

## 阻塞于

- `.scratch/skills-cli-rules/issues/02-agent-rules-target-paths.md`
