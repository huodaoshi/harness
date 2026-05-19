# 接入真实国内 LLM（ChatModelGateway）

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** HITL  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

将假 ChatModel 替换为 **eino-ext ChatModel**（火山方舟 / 豆包等，由维护者提供 API 配置）。实现 `ChatModelGateway`：流式 token、超时、可配置 failover 接口（第二供应商可先返回明确错误）。

**ModeRouter** 注入洪峰/普聊 system 模板（倾听、情绪命名 vs 较轻语气）。dev 可用 env；staging 接真实密钥。

## 验收标准

- [x] distress / normal 模式产生可区分风格的回复（fake 带【洪峰】/【普聊】前缀；ark 靠 system 模板区分，需人工抽检）
- [x] P95 首 token 延迟在 staging 可测并记录（`llm_first_token` 日志 + `LLM_FIRST_TOKEN_TARGET_MS`）
- [x] API 密钥不入库；从环境变量读取
- [x] 提供商失败时 SSE `error` 且不泄露密钥
- [x] 危机/医疗/block 路径仍零 LLM（回归 #02、#05）

## 实现说明

- 包：`backend/internal/chatmodel`（`Gateway`、`ark`、`fake`、failover 占位）
- Pass 路径：`StreamPassTurn` 真流式 SSE；门禁支路不调用 Gateway
- 默认 `LLM_PROVIDER=fake`；配置 `ARK_API_KEY` + `ARK_MODEL_ID` 后切 `ark`

## 阻塞于

- [07-session-persist-and-summary-card.md](./07-session-persist-and-summary-card.md)

## 覆盖的用户故事

#2、#3、#21、#22

## 评论

### HITL 说明

需维护者选定：**厂商、模型名、合同/数据不出境条款**后再由 Agent 或人工接代码。密钥与计费账号由人配置。
