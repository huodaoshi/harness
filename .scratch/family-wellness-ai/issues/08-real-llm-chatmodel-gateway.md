# 接入真实国内 LLM（ChatModelGateway）

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** HITL  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

将假 ChatModel 替换为 **eino-ext ChatModel**（火山方舟 / 豆包等，由维护者提供 API 配置）。实现 `ChatModelGateway`：流式 token、超时、可配置 failover 接口（第二供应商可先返回明确错误）。

**ModeRouter** 注入洪峰/普聊 system 模板（倾听、情绪命名 vs 较轻语气）。dev 可用 env；staging 接真实密钥。

## 验收标准

- [ ] distress / normal 模式产生可区分风格的回复（人工抽检 ≥5 轮）
- [ ] P95 首 token 延迟在 staging 可测并记录（目标值由维护者写入配置）
- [ ] API 密钥不入库；从环境变量读取
- [ ] 提供商失败时 SSE `error` 且不泄露密钥
- [ ] 危机/医疗/block 路径仍零 LLM（回归 #02、#05）

## 阻塞于

- [07-session-persist-and-summary-card.md](./07-session-persist-and-summary-card.md)

## 覆盖的用户故事

#2、#3、#21、#22

## 评论

### HITL 说明

需维护者选定：**厂商、模型名、合同/数据不出境条款**后再由 Agent 或人工接代码。密钥与计费账号由人配置。
