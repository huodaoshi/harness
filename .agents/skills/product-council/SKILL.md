---
name: product-council
description: 通过可配置的推进席/质疑席视角、Analyst 调研与 Default-forward 探索新产品；会读取 landscape 产出的 references/；开场确认议会席与讨论轮次，默认至少 10 轮后才提议成文。在用户要从零探索产品、要立项建议书、或提到 product council 时使用；竞品摸底请先用 landscape。
---

# 产品议会（Product Council）

**发现阶段**：对抗式对话 + Agent 调研 + 可配置**视角席位**（非写死职业）→《项目建议书》。

**推荐前置：** [`landscape`](../landscape/SKILL.md) 产出 `.scratch/<feature>/references/`。议会**不**重复做全面竞品扫描。

## 职能 vs 席位

| 职能 | 说明 |
|------|------|
| **Analyst** | 调研事实。见 [RESEARCH.md](RESEARCH.md) |
| **Advocate** | 推进；每轮戴 **1 顶**推进席帽（来自你的菜单） |
| **Skeptic** | 质疑；每轮 **1～3** 个质疑席，各 ≤2 条反驳 |
| **Facilitator** | 计轮次、假设登记、默认推进、少问你 |

**席位菜单（你可改）：** [SEATS.md](./SEATS.md)  
项目覆盖：`.scratch/<feature-slug>/council-seats.md`（有则优先）

## 默认模式：Default-forward

用户不确定时：Analyst 先查 → Advocate 推荐 → Skeptic 多视角质疑 → Facilitator「除非你反对…」

## 何时才问用户

席位与轮次在入场配好。议会中**禁止每轮必问**；仅品味/红线/无法调研的阻塞项可一问。连续两轮不得提问（除非你插话）。

## 每轮输出格式

```markdown
### Analyst
…

### Advocate · <推进席标签>
…

### Skeptic · <质疑席标签 1>
…

### Skeptic · <质疑席标签 2>
（仅启用多席时）

### Facilitator
**进度：** 第 N / 目标 轮
**本轮席位：** 推进=<…>；质疑=<…>
**假设登记册：** …
**默认推进：** …
```

## 流程

### 0. 入场（配置议会）

1. 复述用户想法；读 `CONTEXT.md` / ADR / `.scratch`（若有）
2. **读取 references/（竞品摸底）：**
   - 路径：`.scratch/<feature-slug>/references/`（`README.md`、`user-notes.md`、各参考品 md）
   - **有内容** → 本会 Analyst 优先引用，勿整库重搜竞品；可针对空白点补搜
   - **无或仅空目录** → 告知：「建议先 `/landscape`」；用户坚持开会则 Analyst 入场时补做竞品搜索（仍不写建议书）
   - `feature-slug` 与 `references/` 一致；若无则与 landscape 相同规则拟定（缺省 `product-discovery`）
3. **读取并确认席位：** `council-seats.md`（项目）或 [SEATS.md](./SEATS.md)；展示清单（启用、标签、审视重点）。用户可改文件或口述。确认后锁定 `推进席列表`、`质疑席列表`
4. **约定轮次（一问）：** 「计划讨论几轮？默认 **10** 轮后再提议成文。」→ `目标轮次`
5. 产物路径：`.scratch/<feature-slug>/项目建议书.md`

**禁止**在席位与轮次未确认前开始第 1 轮议会。

### 1. 议会轮次

**Analyst（按需）→ Advocate（本轮推进席）→ 各启用质疑席 → Facilitator**

- 推进席按 [RESEARCH.md](RESEARCH.md) 议题维度轮换（勿每轮问用户选帽）
- 质疑席：对**已启用**席位各发言；勿新增未在清单中的席

**成文门槛：**

- `当前轮次 < 目标轮次` → **禁止**提议成文
- `≥ 目标轮次` → 可问是否写建议书；未表态则继续
- 用户明确「成文 / 够了」→ 可提前成文

### 2. 成文

[PROPOSAL-TEMPLATE.md](./PROPOSAL-TEMPLATE.md)；建议书「议会记录」须列出本场**最终席位配置**。

### 3. 下一步

| 状态 | 建议 |
|------|------|
| 待验证多 | 继续本会话 / `/prototype` |
| 可写规格 | `/to-prd` |

## 规则

- 席位标签用户可改；**技能不写死职业**
- Skeptic 须可证伪 + 验证实验；技术 **1 推荐 + 1 备选**
- 轮次默认 10；未满禁止提议成文

## 快速开始

`landscape` → 读 `references/` → 展示 [SEATS.md](./SEATS.md) → 确认席位 → 确认轮次 → 第 1 轮（**不成文**）→ …
