# Scaffold-04：Mask 与 locales（可选增强）

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

在 #03 基础上继续迁移：

- `app/masks/` 及 `package.json` 中 `mask` / `mask:watch` 脚本
- `app/locales/`

确保 Mask 可切换且不影响托管单模型聊天。本 issue **不阻塞** Scaffold 最小 E2E 闭环（#03 完成即 epic 可演示）。

## 验收标准

- [ ] `yarn mask`（或 `mask:watch`）可生成面具资源
- [ ] UI 中可切换 Mask，聊天仍走 `/api/bytedance`
- [ ] 中文 locale 基本可用（无大面积 key 裸露）

## 阻塞于

- [03-nextchat-chat-core-e2e.md](./03-nextchat-chat-core-e2e.md)

## 覆盖的用户故事

无（体验增强）
