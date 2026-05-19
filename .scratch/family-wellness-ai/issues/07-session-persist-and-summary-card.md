# 纵向切片：会话持久化 + summary3 + 结束卡片

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

实现 **sessions** 存储：创建会话、追加 user/assistant 消息、单会话轮次 **cap（50）**。Graph **PostProcess** 在会话结束生成 `summary3`（3 句）写入 `session_summaries`；下一会话 #03 inject 读取最近 1 条。

Web：会话正常结束后可选展示 **3 句小结卡片**（可关闭）。历史仅展示**当前会话**消息列表（不要求跨会话历史页）。

## 验收标准

- [x] 多轮对话后 Mongo `sessions.messages` 与 UI 一致
- [x] 超 cap 时行为明确（拒绝新消息或归档，须在 issue 实现说明中写清）
- [x] 结束后 `session_summaries` 有 3 条字符串
- [x] 新会话 inject 含上一场 summary3
- [x] Web 小结卡可关闭且不再自动弹出同会话

## Cap 策略（#07 实现）

单会话 `messages` 数组硬上限 **50 条**（含 user + assistant）。达到上限后 **拒绝新消息**：流式接口返回 HTTP 409 / SSE `error` code `session_message_cap`，提示结束本场并开始新会话。不做自动归档。

## API

- `GET /v1/sessions/:id?user_id=` — 读取当前会话消息
- `POST /v1/sessions/end` — 结束会话、生成 `summary3` 并写入 `session_summaries`

## 阻塞于

- [03-spike-s3-profile-inject.md](./03-spike-s3-profile-inject.md)
- [04-web-distress-chat-shell.md](./04-web-distress-chat-shell.md)

## 覆盖的用户故事

#6、#7、#23
