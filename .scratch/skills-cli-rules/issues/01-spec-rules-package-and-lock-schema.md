Status: ready-for-agent

# 规格：规则包布局与 rules-lock.json 草案

## 要构建什么

在代码实现前，把「远程规则包长什么样、锁文件长什么样」写成**可实施的**约定文档（含字段说明与示例），并与 ADR-0001 一致。产出放在 `cli/` 的 README 增补段落，或 `docs/` 下一页短文（二选一或都写，不重复堆砌）。

文档需明确：

- 包内必选目录：`rules/cursor/`、`rules/claude/`（与 `-a` 的对应关系）。
- `--to` 行为：将选中子树**以内**的文件铺到目标路径，不额外嵌套一层 `rules/cursor`。
- `rules-lock.json` 建议字段：至少能锁定 **来源（Git URL 或路径约定）**、**技能 CLI 可解析的 revision/commit 或等价物**、**针对哪一 Agent（cursor / claude-code）**、可选 **安装时选项**（如 `--copy`）。

不要求最终实现锁文件解析，但 schema **须能被 issue 03/04 直接实现**，避免返工。

## 验收标准

- [ ] 文档为简体中文，且与 `docs/adr/0001-skills-cli-rules-install.md` 无矛盾。
- [ ] 含一份 **最小示例** `rules-lock.json`（可标注为草案）。
- [ ] 明确 MVP **不包含** list/remove 与双 `-a`。

## 阻塞于

无——可立即开始。
