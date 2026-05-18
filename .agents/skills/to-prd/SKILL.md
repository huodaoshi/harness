---
name: to-prd
description: 将当前对话上下文整理为 PRD 并发布到项目 issue 跟踪器。适用于用户希望基于当前上下文创建 PRD 时。
---

本技能根据当前对话上下文与对代码库的理解产出 PRD。**不要**采访用户——直接综合你已知信息即可。

issue 跟踪器与分拣标签词汇表应已提供给你；若尚未配置，请运行 `/setup-matt-pocock-skills`。

## 流程

1. 若尚未了解代码库，先探索仓库以把握现状。PRD 全文使用项目的领域术语表用语，并尊重所涉区域内的既有 ADR。

2. 勾勒为实现所需新建或修改的主要模块。主动寻找可单独测试的**深模块**（deep module）抽取机会。

深模块（相对浅模块）指：用简单、可测、极少变更的接口封装大量功能。

与用户确认这些模块是否符合预期，并确认用户希望为哪些模块编写测试。

3. 按下方模板撰写 PRD，并发布到项目 issue 跟踪器。应用 `ready-for-agent` 分拣标签——无需额外分拣。

<prd-template>

## 问题陈述

从用户视角描述用户面临的问题。

## 解决方案

从用户视角描述针对该问题的解决方案。

## 用户故事

一份**详尽**的编号用户故事列表。每条采用格式：

1. As an <actor>, I want a <feature>, so that <benefit>

<user-story-example>
1. As a mobile bank customer, I want to see balance on my accounts, so that I can make better informed decisions about my spending
</user-story-example>

列表应尽可能全面，覆盖该特性的各个方面。

## 实现决策

已做出的实现决策列表，可包括：

- 将新建/修改的模块
- 将修改的模块接口
- 来自开发者的技术澄清
- 架构决策
- Schema 变更
- API 契约
- 具体交互

**不要**写入具体文件路径或代码片段——它们可能很快过时。

例外：若原型中的片段比 prose 更精确地表达某项决策（状态机、reducer、schema、类型形状），可在对应决策下内联，并简短注明来自原型。只保留与决策相关的部分——不是可运行 demo，而是要点。

## 测试决策

已做出的测试决策列表，包括：

- 何谓好测试的说明（只测外部行为，不测实现细节）
- 将覆盖测试的模块
- 测试的先例（代码库中类似测试）

## 范围外

本 PRD 明确不包含的范围。

## 补充说明

关于该特性的其他说明。

</prd-template>
