# 领域文档（Domain Docs）

工程技能在探索代码库时应如何消费本仓库的领域文档。

**生成约定：** 新建或更新 `CONTEXT.md`、`docs/adr/*.md` 时，正文使用**简体中文**。术语表条目、ADR 标题与叙述均用中文；文件名可继续使用 `0001-简短英文-slug.md` 这类编号 + slug 形式。

## 探索前先读这些

- 仓库根目录的 **`CONTEXT.md`**，或
- 若存在根目录 **`CONTEXT-MAP.md`** — 它指向各上下文的 `CONTEXT.md`，按主题读取相关文件
- **`docs/adr/`** — 阅读与即将改动区域相关的 ADR；多上下文仓库还需查看 `src/<context>/docs/adr/`

若上述文件不存在，**静默继续**。不要强调缺失，也不要主动建议「先建一份」；生产者技能（`/grill-with-docs`）会在术语或决策真正落定后懒创建。

## 目录结构

单上下文仓库（本仓库采用）：

```
/
├── CONTEXT.md
├── docs/adr/
│   ├── 0001-event-sourced-orders.md
│   └── 0002-postgres-for-write-model.md
└── src/
```

多上下文仓库（根目录存在 `CONTEXT-MAP.md` 时）：

```
/
├── CONTEXT-MAP.md
├── docs/adr/                          ← 全系统级决策
└── src/
    ├── ordering/
    │   ├── CONTEXT.md
    │   └── docs/adr/                  ← 该上下文内的决策
    └── billing/
        ├── CONTEXT.md
        └── docs/adr/
```

## 使用词汇表中的术语

输出中若命名领域概念（issue 标题、重构建议、假设、测试名等），应使用 `CONTEXT.md` 中的定义，避免使用词汇表明确排斥的同义词。

若所需概念尚不在词汇表中，说明要么在用项目未采纳的说法（应重新考虑），要么存在真实缺口（可记一笔，供 `/grill-with-docs` 补充）。

## 标明与 ADR 的冲突

若你的结论与现有 ADR 矛盾，应明确写出，而不是悄悄覆盖：

> _与 ADR-0007（事件溯源订单）矛盾——但值得重新讨论，因为……_
