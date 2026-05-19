# 内测会话反馈（「更贴关系」）

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

会话 `done` 后 Web 展示可选 **1–5 分或是/否**：「是否比通用 AI 更贴你的关系处境？」；`POST /v1/feedback` 写入 Mongo（`session_id`、`score`、`created_at`），无 PII 扩采。

用于内测 H3/H4 观察，不做公开排行榜。

## 验收标准

- [ ] 用户可跳过反馈
- [ ] 提交后持久化且可按 session 查询（维护者脚本即可）
- [ ] 不在反馈中收集自由文本（MVP 防 PII 扩散）
- [ ] 隐私政策占位提及反馈用途

## 阻塞于

- [07-session-persist-and-summary-card.md](./07-session-persist-and-summary-card.md)

## 覆盖的用户故事

#15
