# Scaffold-03：迁 NextChat 聊天核心并完成 E2E

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

按 `frontend/MIGRATION.md`，从 `NextChat/` 迁入 **聊天主路径** 所需目录（至少）：

- `app/components/`（聊天 UI）
- `app/store/`
- `app/client/`（`ApiPath` 请求经 `lib/harness/api-base.ts` / 同域 `/api`）
- `app/utils/`、`app/constant.ts`、`app/typing.ts`、`app/styles/`、`public/`
- 合并 `NextChat/package.json` 运行时依赖；`next.config.mjs` 保持 SVG/Sass 与 `/api` rewrite

**不复制** `app/api/`、`app/config/server.ts`（密钥仅 backend）。运行时通过 **`GET /api/config`** 拉配置；`hideUserApiKey` + 单条 `customModels` 隐藏多厂商/Key 设置（**不删除** settings 代码）。

本地同时启动 `backend`（:8080）与 `frontend`（:3000），在浏览器完成：**打开聊天页 → 发送一条消息 → 经 `/api/bytedance` 流式收到回复**。

## 验收标准

- [ ] `frontend` `yarn dev` 可进入 NextChat 聊天界面（非仅占位页）
- [ ] 一问一答流式正常，`done` 或等价结束
- [ ] 设置页或启动流程能读到 `/api/config`（单模型、隐藏用户 Key）
- [ ] 不依赖用户在前端填写 API Key 或多模型选择
- [ ] `yarn build` 通过（或 documented 已知限制）

## 阻塞于

- [01-api-config-e2e.md](./01-api-config-e2e.md)
- [02-api-bytedance-proxy.md](./02-api-bytedance-proxy.md)

## 覆盖的用户故事

通用聊天壳（技术）；PRD 洪峰/档案属后续 Frontend P1
