# 领域文档

工程类技能在探索代码库时，应如何消费本仓库的领域文档。

## 探索前先读

- 仓库根目录的 **`CONTEXT.md`**，或
- 若存在根目录 **`CONTEXT-MAP.md`** — 它指向各上下文的 `CONTEXT.md`。阅读与当前话题相关的那些。
- **`docs/adr/`** — 阅读与即将改动区域相关的 ADR。多上下文仓库中，另查 `src/<context>/docs/adr/` 中的上下文内决策。

若上述文件不存在，**静默继续**。勿强调缺失；勿主动建议预先创建。产出方技能（`/grill-with-docs`）在术语或决策实际敲定后懒创建。

## 文件结构

单上下文仓库（多数）：

```
/
├── CONTEXT.md
├── docs/adr/
│   ├── 0001-event-sourced-orders.md
│   └── 0002-postgres-for-write-model.md
└── src/
```

多上下文仓库（根目录存在 `CONTEXT-MAP.md`）：

```
/
├── CONTEXT-MAP.md
├── docs/adr/                          ← 全系统级决策
└── src/
    ├── ordering/
    │   ├── CONTEXT.md
    │   └── docs/adr/                  ← 上下文内决策
    └── billing/
        ├── CONTEXT.md
        └── docs/adr/
```

## 使用术语表的用语

输出中命名领域概念时（issue 标题、重构建议、假设、测试名），使用 `CONTEXT.md` 中的定义。勿改用术语表明确避免的同义词。

若所需概念不在术语表中，这是信号——要么你在发明项目未用的语言（应重新考虑），要么确有缺口（记下供 `/grill-with-docs` 处理）。

## 标注 ADR 冲突

若输出与既有 ADR 矛盾，应显式指出，勿静默覆盖：

> _与 ADR-0007（event-sourced orders）矛盾——但值得重新讨论，因为…_
