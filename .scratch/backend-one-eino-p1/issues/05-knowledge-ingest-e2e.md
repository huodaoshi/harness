# P1-05：知识库默认 Space + 精简入库 E2E

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

复制改写 `modules/knowledge`、`modules/knowledgeindexing`（**无** score/tag LLM 批处理）。Bootstrap **一个默认 Space**，`doc_types` 为 1=产品与边界、2=FAQ、3=运营话术/示例。

`mq.provider=local` 时在 API 进程**内嵌**消费 `knowledge-ingest`（对齐 one-eino）。管理员 JWT+Admin 可提交 inline/url 入库，job 完成后 Redis 向量可检索（管理端或测试用，**不**注入 Wellness stream）。

## 验收标准

- [x] 启动后存在默认 Space（`EnsureDefaultSpace` bootstrap）
- [x] Admin 提交 markdown ingest → job status=done → `FT.SEARCH`/Search 能命中 chunk（需 Redis Stack；单测在无 FT 时仅验 job done）
- [x] 流水线仅：resolve → parse → chunk → assignID → index
- [x] 危机文案仍走 YAML/JSON SafetyGate，不依赖本 Space

## 阻塞于

- [03-auth-jwt-guest.md](./03-auth-jwt-guest.md)

## 覆盖的用户故事

运营知识维护（PRD 外延）；ADR P1 知识库目标
