Status: ready-for-agent

# 测试补强与 README 面向使用者

## 要构建什么

- 为 **rules add** 与 **lock 恢复** 补全关键路径测试（与现有 Vitest 风格一致），避免仅 happy path。
- 更新 **`cli/README.md`**（或本会话已定稿的文档位置）：安装规则包、`rules-lock.json`、与 **skills** 命令对照表；指向 ADR-0001。

不扩大 MVP 范围（不写 list/remove）。

## 验收标准

- [ ] `pnpm test`（或仓库既定命令）在 `cli/` 下通过。
- [ ] 新用户仅凭 README 能完成「装一次 + 锁恢复」的最小流程。
- [ ] 文档中说明超出 MVP 的能力**未**承诺。

## 阻塞于

- `.scratch/skills-cli-rules/issues/03-implement-skills-rules-add.md`
- `.scratch/skills-cli-rules/issues/04-implement-rules-lock-install.md`
