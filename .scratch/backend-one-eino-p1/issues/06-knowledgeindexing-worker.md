# P1-06：knowledgeindexing 独立 Worker

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

新增 `cmd/knowledgeindexing`，在 `mq.provider=rocketmq`（或等价生产配置）时独立消费 `knowledge-ingest`。与 `cmd/server` 内嵌 local 消费互斥文档化。

## 验收标准

- [ ] `go build ./cmd/knowledgeindexing` 成功
- [ ] Worker 与 API 进程可同时连 Mongo/Redis/MQ（集成测试或文档化手工步骤）
- [ ] 与 05 相同 ingest payload 在 Worker 路径可完成

## 阻塞于

- [05-knowledge-ingest-e2e.md](./05-knowledge-ingest-e2e.md)

## 覆盖的用户故事

—
