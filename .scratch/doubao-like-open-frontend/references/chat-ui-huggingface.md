# 参考：Hugging Face Chat UI

**类型：** 参考品
**链接：** https://github.com/huggingface/chat-ui
**来源：** Agent 搜索
**置信度：** 高
**搜索日期：** 2026-05-19

## 他们解决什么

HuggingChat 同款开源聊天前端，接 OpenAI 兼容端点与 HF Inference。

## 目标用户

想复刻 HuggingChat 体验、或需要可定制 Web Chat 的团队。

## 核心功能（与本次想法相关的）

- SvelteKit + TypeScript
- 多提供商、模型路由、多模态上传
- MCP tools、可选 OIDC
- MongoDB 持久化（开发可用内嵌库）

## 定价 / 商业模式

Apache 2.0；HF 云服务另计。

## 优点

- 代码结构清晰，适合作「聊天前端」架构参考
- 国际社区标准实现之一

## 缺点 / 局限

- UI 偏 HuggingChat，非豆包国内产品范式
- 依赖 MongoDB 等，国内栈团队可能不熟

## 与我们的关系

不适用直接竞争 — 学架构/路由/多提供商接入，非抄豆包 UI。

## 待核实

- [ ] 无 Mongo 生产部署的最佳实践
