# 竞品与参考品索引

**feature：** doubao-like-open-frontend
**扫描日期：** 2026-05-19
**品类（Agent 归纳）：** 自托管 / 可二开的 AI 对话 Web（及桌面）前端，体验对标豆包类 C 端助手

## 怎么用

供 `product-council` 或人工阅读。Agent 搜索内容标置信度，决策前请核实。

## 清单

| 类型 | 名称 | 文件 | 置信度 | 一句话 |
| ---- | ---- | ---- | ------ | ------ |
| 直接竞品 | Lobe Chat | [lobe-chat.md](./lobe-chat.md) | 高 | 现代助手 UI，功能全，最接近「豆包式」壳 |
| 直接竞品 | NextChat | [nextchat.md](./nextchat.md) | 高 | 轻量跨端聊天壳，部署快，偏 ChatGPT 布局 |
| 直接竞品 | Open WebUI | [open-webui.md](./open-webui.md) | 高 | 聊天+平台一体，功能重，Docker 友好 |
| 直接竞品 | Mini-Doubao | [mini-doubao.md](./mini-doubao.md) | 中 | 对标豆包的全栈示例，体量小待核实 |
| 间接竞品 | Dify | [dify.md](./dify.md) | 高 | 带 Chat 的 LLM 平台，非纯前端库 |
| 间接竞品 | ChatOllama | [chat-ollama.md](./chat-ollama.md) | 高 | Nuxt 聊天壳，本地模型 + 知识库/Agent |
| 参考品 | Hugging Face Chat UI | [chat-ui-huggingface.md](./chat-ui-huggingface.md) | 高 | HuggingChat 架构参考，非豆包审美 |

## 刻意未纳入

| 名称 | 原因 |
| ---- | ---- |
| Better Doubao | 浏览器扩展，非独立开源前端 |
| 2025doubao-free-api 等 | API 代理，无对话 UI |
| 豆包微信机器人插件 | 后端插件，非 Web 前端 |

## 待补充

- [ ] 用户是否要「仅 UI 组件」vs「全栈平台」
- [ ] 是否必须内置火山/豆包 API
- [ ] 移动端 / 小程序需求
