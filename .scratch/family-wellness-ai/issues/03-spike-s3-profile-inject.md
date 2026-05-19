# Spike S3：Mongo 关系档案 + 最近摘要注入 Graph

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — W1 门禁 S3

## 要构建什么

接入 **MongoDB（自建，dev 可用 compose）**：实现 `relationship_profiles`、`session_summaries` 读写（形状见 PRD）。在 Graph 增加 **ProfileInject** 节点：对 `pass` 路径，将档案与**最近 1 条** `summary3` 编入 prompt 上下文（假 ChatModel 可记录收到的 context 供断言）。

提供 **20 条自动化**用例：种子档案/摘要 → 发消息 → 断言 inject 块包含关键字段（非全文 golden 字符串）。

## 验收标准

- [x] Mongo 集合与 PRD snake_case 字段一致（`store` 模型 + `MongoStore`）
- [x] 无档案时 inject 为空，不 500
- [x] 有档案 + 摘要时，假 chat 回复含 inject 关键字段
- [x] 20 条自动化（`profile_inject_test.go`）+ distress 补充用例
- [x] `docker-compose.yml` + README

## 评论

### 2026-05-19 · 实现完成

- Graph：`profile_inject` → `fake_chat`（危机路径不经过 inject）
- `USE_MEMORY_STORE=true` 用于无 Mongo 测试；默认连 `MONGODB_URI`

## 阻塞于

- [01-spike-s1-stream-sse.md](./01-spike-s1-stream-sse.md)

## 覆盖的用户故事

#3、#4、#6、#17
