---
description: Eino / ADK / eino-ext——backend Go 服务编排与 RAG 规范
paths:
  - "backend/**/*.go"
---

# Eino / ADK / eino-ext（harness `backend/`）

> **生效条件：** 仅当仓库存在 `backend/**/*.go` 时由 Claude Code / Cursor 按路径启用。  
> **产品上下文：** family-wellness-ai 等项目的 **关系会话流水线** 优先用 `compose.Graph`（如 `SafetyGate → ModeRouter → ProfileInject → ChatModel`），见 `.scratch/family-wellness-ai/项目建议书.md`。本 rule 与 PRD 冲突时 **以 PRD + 用户拍板为准**。

## 核心原则

Eino / ADK / eino-ext **已提供的能力必须直接用**，禁止手写等价逻辑（Agent 执行、流式、Checkpoint、Session KV、工具调用、RAG 组件链、向量检索）。  
**业务自建**（不受此约束）：鉴权、Mongo 档案 CRUD、入库任务状态机、审计、SafetyGate 规则表等。

## 编排选型（与 PRD 对齐）

| 场景 | 用什么 | 禁止 |
| ---- | ------ | ---- |
| **关系会话 Graph**（安全分支、模式路由、档案注入、流式回复） | `compose.Graph` / `compose.Chain` | 裸 goroutine 拼 SSE；绕过 Graph 直调模型 |
| **对话 Agent**（多轮工具、ADK 事件流） | `adk.NewChatModelAgent` + `adk.Runner` | 裸 `react.NewAgent`；手写流消费循环 |
| 固定步骤流水线（知识入库） | `compose.Graph` 线性拓扑 | for 循环手串组件 |
| 并行 / 条件 / 有环 | `compose.Graph` | — |
| Graph 作为 Agent 工具 | `graphtool.NewInvokableGraphTool` | — |
| Human-in-the-Loop / Checkpoint | `ApprovableTool` + `Runner.ResumeWithParams` | 生产用 `InMemoryStore` 作唯一 Checkpoint |

**危机 / 医疗边界支路：** Graph **编译期**不连接 `ChatModel`；须可测试「LLM 调用次数 = 0」（见项目建议书 C3 / Spike S2）。

## ADK

- **`adk.Runner`** 为 Agent 唯一入口；不直接调 `agent.Generate` / `agent.Stream`
- **`CheckpointStore` 持久化：** 实现 `adk.CheckpointStore`（项目建议：**MongoDB 自建**，与业务库同运维）；禁止生产依赖 `adkstore.NewInMemoryStore`
- 工具间请求级状态：`adk.AddSessionValue` / `adk.GetSessionValue`，不写 Tool 结构体字段
- 工具错误：`safeToolMiddleware`；LLM 限流：`adk.ModelRetryConfig`，不手写 retry 循环
- 每个 `ChatModelAgent` 的 `Handlers` 须注册 `safeToolMiddleware`

## RAG（若有知识库模块）

- 入库：`compose.Graph`（Loader → Transformer → Indexer），不手 for 串联
- Parser：HTML/PDF 用 eino-ext；Markdown/JSON/CSV/纯文本等无官方包时 **实现 `parser.Parser` 接口**
- Embedder：`eino-ext/components/embedding/ark`（或项目选定 provider），不直调 HTTP
- Redis 向量索引：`InitRedisIndex` 在服务启动完成；`KeyPrefix` / `VectorField` / `DIM` 与 Indexer 一致

## Retriever（若启用）

- Redis Client `Protocol: 2`
- `VectorField` 与 Indexer `EmbedKey` 一致（常用 `"embedding"`）
- Hash 字段勿用保留名 `"distance"`
- 分数转换写在 `DocumentConverter`，业务层不重算

## Tool

- 实现 `tool.InvokableTool`（`Info` + `InvokableRun`）
- `InvokableRun` 不 panic；可恢复错误返回字符串；`compose.IsInterruptRerunError` 须原样 `return "", err`
- Schema 优先 `utils.InferTool`；复杂约束再手写 `schema.NewParamsOneOfByParams`

## StreamReader

必须 `defer sr.Close()`；禁止多 goroutine 并发 `Recv()`：

```go
sr, err := runnable.Stream(ctx, input)
if err != nil {
    return fmt.Errorf("stream: %w", err)
}
defer sr.Close()
for {
    chunk, err := sr.Recv()
    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        return fmt.Errorf("recv: %w", err)
    }
    // handle chunk
}
```

## eino-ext 优先

新组件先查 [cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) 是否已有实现。

**引入检查：**

1. 版本与 `go.mod` 内其他 `eino-ext` 包同 release tag
2. `APIKey` / `BaseURL` / `Model` 从配置读取，禁止硬编码
3. `loader/url`、`loader/s3` 须配 SSRF / Host 白名单

**常用包（按项目选用）：**

| 组件 | 包路径 |
| ---- | ------ |
| ChatModel | `eino-ext/components/model/ark` 等 |
| Embedder | `eino-ext/components/embedding/ark` |
| Indexer / Retriever | `eino-ext/components/indexer/redis`、`retriever/redis` |

## 示例代码（本地参考）

官方示例在 **CloudWeGo `eino-examples`** 仓库。路径因机器而异：

- 优先读环境变量 **`EINO_EXAMPLES_ROOT`**（指向本地 clone 根目录）
- 未设置时，开发前向用户确认路径；**禁止在 rule 或代码中写死 `D:\...` 盘符**

| 场景 | 典型路径（相对 `EINO_EXAMPLES_ROOT`） |
| ---- | ------------------------------------- |
| ADK Agent + Runner | `quickstart/chatwitheino/cmd/ch02` |
| 工具 + safeToolMiddleware | `quickstart/chatwitheino/cmd/ch05` |
| HITL + Resume | `quickstart/chatwitheino/cmd/ch07` |
| Checkpoint（参考接口，生产换 Mongo） | `quickstart/chatwitheino/cmd/ch03` |
| RAG 入库 Graph | `quickstart/eino_assistant/eino/knowledgeindexing/` |
| Tool 定义 | `quickstart/todoagent/` |

参考示例只学**模式**；密钥与模型名从本项目 `conf` / 环境配置读取。

## 与本仓其他 rule

- 命名：`rules/naming.md`（Go `sessionID` + `json:"session_id"`）
- Git：`rules/git.md`（`backend/` 开发分支）
- 文档语言：`docs-zh.md`
