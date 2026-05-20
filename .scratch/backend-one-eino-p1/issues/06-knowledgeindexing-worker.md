# P1-06：knowledgeindexing 独立 Worker

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

新增 `cmd/knowledgeindexing`，在 `mq.provider=rocketmq`（或等价生产配置）时独立消费 `knowledge-ingest`。与 `cmd/server` 内嵌 local 消费互斥文档化。

## 验收标准

- [x] `go build ./cmd/knowledgeindexing` 成功
- [x] Worker 与 API 进程可同时连 Mongo/Redis/MQ（见 backend/README MQ 说明）
- [x] 与 05 相同 ingest payload 在 Worker 路径可完成（rocketmq 或 server 内嵌 local）

## 阻塞于

- [05-knowledge-ingest-e2e.md](./05-knowledge-ingest-e2e.md)

## 覆盖的用户故事

—
