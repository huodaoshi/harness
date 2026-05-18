# CONTEXT.md 格式

## 结构

```md
# {上下文名称}

{一两句话说明该上下文是什么、为何存在。}

## Language

**Order**:
{对该术语的简明定义}
_Avoid_: Purchase, transaction

**Invoice**:
向客户在交付后发出的付款请求。
_Avoid_: Bill, payment request

**Customer**:
下单的个人或组织。
_Avoid_: Client, buyer, account

## Relationships

- 一笔 **Order** 可产生一条或多条 **Invoice**
- 每条 **Invoice** 恰好属于一个 **Customer**

## Example dialogue

> **Dev:** 「**Customer** 下 **Order** 时，是否立即创建 **Invoice**？」
> **Domain expert:** 「否——仅在 **Fulfillment** 确认后才生成 **Invoice**。」

## Flagged ambiguities

- 「account」曾混指 **Customer** 与 **User**——已厘清：二者为不同概念。
```

## 规则

- **要有立场。** 同一概念存在多个词时，选定最佳用词，其余列入应避免的同义词。
- **显式标注冲突。** 术语用法含糊时，在「Flagged ambiguities」中写明及最终决议。
- **定义要短。** 每条最多一句；定义它**是什么**，而非它**做什么**。
- **写清关系。** 术语名加粗，在显而易见处标明基数。
- **只收录本项目上下文特有的概念。** 通用编程概念（超时、错误类型、工具模式等）即便项目大量使用也不应写入。新增术语前自问：这是本上下文独有概念，还是通用编程概念？仅前者收录。
- **自然聚类时用子标题分组。** 若所有术语同属一块领域，平铺列表亦可。
- **写一段示例对话。** 开发与设计/领域专家之间的对话，自然展示术语如何配合，并厘清相近概念的边界。

## 单上下文与多上下文仓库

**单上下文（多数仓库）：** 仓库根目录一个 `CONTEXT.md`。

**多上下文：** 根目录 `CONTEXT-MAP.md` 列出各上下文、路径及彼此关系：

```md
# Context Map

## Contexts

- [Ordering](./src/ordering/CONTEXT.md) — 接收并跟踪客户订单
- [Billing](./src/billing/CONTEXT.md) — 生成发票并处理付款
- [Fulfillment](./src/fulfillment/CONTEXT.md) — 管理仓内拣货与发货

## Relationships

- **Ordering → Fulfillment**：Ordering 发出 `OrderPlaced` 事件；Fulfillment 消费以开始拣货
- **Fulfillment → Billing**：Fulfillment 发出 `ShipmentDispatched` 事件；Billing 消费以生成发票
- **Ordering ↔ Billing**：共享 `CustomerId`、`Money` 等类型
```

技能按以下方式判断结构：

- 存在 `CONTEXT-MAP.md` → 读取以定位各上下文
- 仅根目录有 `CONTEXT.md` → 单上下文
- 二者皆无 → 首个术语敲定后懒创建根目录 `CONTEXT.md`

多上下文时，推断当前话题属于哪一块；不明确则询问用户。
