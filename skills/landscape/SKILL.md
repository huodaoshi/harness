---
name: landscape
description: 根据用户想法上网搜索竞品与参考产品，在 .scratch/<feature>/references/ 生成索引与分文件摘要，并请用户删减或补充后结束。在用户不知道竞品是谁、要准备 references、市场摸底、landscape 扫描时使用；不写项目建议书、不开 product-council。
---

# Landscape（市场摸底）

**只做一件事：** 想法 → 搜索 → 写入 `references/` → 请你删/补 → **结束**。

不写《项目建议书》，不跑议会。完成后可手动再开 `/product-council`。

## 快速开始

用户：「我想做 XXX，不知道竞品有哪些。」

→ 收集 3 句想法 → 搜索 → 写 `references/` → 展示清单 → 等用户删/补 → 结束。

## 流程

### 1. 收集想法（尽量不问）

从对话提取或请用户用三行说明（缺项由 Agent 合理推断，**勿问卷**）：

- **帮谁**
- **解决什么**
- **像什么 / 品类关键词**（可写「不确定」）

确定 `feature-slug`（与用户确认或用短语 kebab-case，缺省 `product-discovery`）。

若尚无 `user-notes.md`，写入 `.scratch/<feature-slug>/references/user-notes.md`（见 [NOTES-TEMPLATE.md](./NOTES-TEMPLATE.md)）。

### 2. 搜索（主动执行，不必请示）

用 Web 搜索，至少覆盖：

- 品类 + 目标用户 + 问题关键词（中文与英文各搜一轮若相关）
- 「X alternative」「X 竞品」「best tools for …」

目标找到 **5～8 个**候选，分为：

- **直接竞品** — 用户可能二选一
- **间接竞品** — 部分重叠
- **参考品** — 仅灵感，非对手

搜不到时扩大关键词或改品类表述，**勿**把「搜不到」抛给用户去猜词。

### 3. 写入 references/

目录：`.scratch/<feature-slug>/references/`

| 文件 | 内容 |
|------|------|
| `README.md` | 按 [README-TEMPLATE.md](./README-TEMPLATE.md) 写索引 |
| `<slug>.md` 每个产品一个 | 按 [REFERENCE-TEMPLATE.md](./REFERENCE-TEMPLATE.md) |
| `user-notes.md` | 用户想法（已存在则更新「最后同步」） |

文件名：`notion.md`、`linear.md` 等小写 kebab，勿用空格。

每条文件顶部 **来源：Agent 搜索**、**置信度：高/中/低**；不确定处写「待核实」。

### 4. 请你删/补（本技能唯一必等用户的一步）

用表格展示：

| 类型 | 名称 | 文件 | 一句话 |
|------|------|------|--------|

然后问：

> 请指出要**删掉**的、要**补充**的（名称或链接）、或回复「可以」结束。

- **删掉** — 删文件并更新 `README.md`
- **补充** — 按名搜索后追加 md
- **可以 / 够了** — 进入结束

用户要求「再搜 N 家」→ 只做补充。

### 5. 结束

- 确认 `references/README.md` 与文件一致
- 一句话总结：找到几家、缺什么仍待核实
- **可选**一句：「下一步可 `/product-council`，请先读 `references/`。」
- **禁止**自动调用 product-council、to-prd，或撰写建议书

## 规则

- 正文**简体中文**；产品名、URL 可保留原文
- 不编造定价/功能；查不到写「未找到公开信息」
- 不把 references 合并进 `CONTEXT.md`（除非用户另要求）
- 本技能结束即停，勿展开产品讨论

## 检查清单

- [ ] `references/README.md` 索引完整
- [ ] 每个产品独立 md，含类型（直接/间接/参考）
- [ ] 已请用户删/补并落实修改
- [ ] 未写项目建议书、未开议会
