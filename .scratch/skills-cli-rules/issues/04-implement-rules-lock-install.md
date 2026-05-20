Status: ready-for-agent

# 实现 rules-lock.json 写入与 experimental_install（恢复）

## 要构建什么

- **`rules add` 成功后**（或显式 flag，按实现择一）能更新项目根 **`rules-lock.json`**，记录足够信息以供恢复（字段形状遵循 issue 01）。
- 提供与现有命名风格一致的恢复命令（如 **`skills rules experimental_install`**），行为类比 **`skills experimental_install`**：**只**根据 `rules-lock.json` 恢复规则，**不**修改 `skills-lock.json`。

并发或多包策略 MVP：**单一来源一条锁**或可文档化的最简规则；若多条目，需在 issue 01 已说明。

## 验收标准

- [ ] 从空目录或干净状态：根据锁文件可重复得到同一规则树（测试可用本地 file URL/fixture）。
- [ ] 与 `skills-lock.json` 互不覆盖同一文件时的字段冲突无（两文件独立）。
- [ ] 错误输入（坏锁、缺目录）有可读错误信息。

## 阻塞于

- `.scratch/skills-cli-rules/issues/01-spec-rules-package-and-lock-schema.md`
- `.scratch/skills-cli-rules/issues/03-implement-skills-rules-add.md`
