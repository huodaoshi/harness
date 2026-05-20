# Harness 领域上下文

关系情绪 AI 产品（family-wellness-ai）的领域用语；实现细节见 PRD 与 ADR，不在此展开。

## Language

**backend**:
Go + Hertz API 服务；托管关怀路径 `/v1/sessions/*`、（规划中的）平台对话 `/v1/chat`、以及 NextChat 兼容 `/api/*` 代理。
_Avoid_: 用 backend 指前端或 NextChat 目录本身。

**Wellness 会话路径**:
用户关怀聊天主入口 `POST /v1/sessions/stream`；经 SafetyGate 与 RelationshipSessionGraph 编排；一期重构中保持此契约不变。
_Avoid_: 与 one-eino 的 `/v1/chat` 混称为同一路由。

**平台对话路径**:
（二期）统一聊天入口 `POST /v1/chat`；ADK Runner + 可选知识检索；由 Wellness 路径逐步并入。
_Avoid_: 在一期方案中默认已替换 `/v1/sessions/stream`。

**关怀模式**:
客户端以 `X-Harness-Mode: wellness` 或等价约定进入 Wellness 会话路径与 SafetyGate。
_Avoid_: 与「通用模式」共用同一存储或编排假设。

**通用模式**:
默认聊天走 `/api/*` 透明代理；历史存浏览器本地（IndexedDB）。
_Avoid_: 称为关怀模式或 wellness 路径。

**知识 Space**:
可版本化的文档集合（入库、向量检索）；内容为**本产品**相关，非游戏攻略垂直。机制沿用 one-eino（`Space` + 自声明 `doc_types`）；**P1 仅运营一个默认 Space**，后续可拆多库。
_Avoid_: 与危机固定模板 YAML/JSON 混为同一存储（危机文案仍走 SafetyGate 配置，非强制入向量库）。

**doc_type（Harness 默认三类，Space 级配置）**:
**1** = 产品与边界；**2** = FAQ；**3** = 运营话术/示例。仅在该 Space 的 `doc_types` 内合法。
_Avoid_: 沿用 one-eino 攻略/设定/剧情 语义。

**关系档案（Profile）**:
用户可编辑的「我 / 重要他人 / 当前议题」结构化记忆（Level 1）。
_Avoid_: 与知识 Space 全文检索结果混称为同一类「记忆」。

**会话摘要（summary3）**:
单场会话结束时的三句 recap（Level 2）；注入下一场 prompt，非向量 RAG。
_Avoid_: 称为知识 chunk 或 RAG 命中。

**里程碑 B / 里程碑 C**:
见 ADR-0001：B = 聊天主路径 + 部分 Go 代理；C = 全厂商代理、插件、关怀 UI 等。

**游客（Guest）**:
未登录用户；请求带合法 `X-Anon-ID`（UUID）；wellness 可聊天；`user_id` 形如 `anon:{uuid}`。
_Avoid_: 与已登录 **User** 混用同一账号体系（绑定/合并策略属二期）。

**知识入库流水线（P1）**:
解析 → 分块 → 赋 chunk ID → 向量索引；**不做** per-chunk LLM 评分/打标。
_Avoid_: 称为与 one-eino 游戏向 score/tag 批处理相同。

**Backend 重构 P1**:
目录与 one-eino 对齐（`conf/`、`infra/`、`api/`、`modules/`）+ auth（SMS JWT + Guest；微信默认关）+ 知识入库/检索；不改 Wellness Graph 与 stream SSE 契约。
_Avoid_: 把 P1 说成已包含 ADK 合流或 `/api/*` 全量代理。

**Backend 重构 P2**:
Wellness 并入平台对话路径（SafetyGate 进 ADK 编排）；可选 A2UI；ADR 中尚未完成的 `/api/*` 等。

## Relationships

- 一个 **Customer**（未来 JWT 用户）可有多场 **Session**（wellness 会话元数据）
- 每个 **Session** 至多绑定一份当前有效的 **关系档案** 注入快照（按 user 读取，非按 session 复制存储）
- 一个 **知识 Space** 包含多条入库文档；检索命中为 **Chunk**，经 RAG 进入模型上下文（时机见 Flagged）
- **会话摘要** 按 user 保留最近 N 场（MVP：N=1），与 **知识 Space** 并行存在

## Example dialogue

> **Dev:** 「运营上传的危机热线文案进 **知识 Space** 后，洪峰模式会自动 RAG 吗？」
> **Domain expert:** 「一期只保证能入库和检索验证；是否注入 **Wellness 会话路径** 的 prompt 单独决策，默认二期跟平台对话路径一起接。」

## Flagged ambiguities

- PRD 原写「MVP 无向量 RAG」——已 supersede：有 **知识 Space**，但 **P1 不注入** Wellness prompt；**P2** 经平台对话路径 / ADK tool 启用。可选 **P1** 仅管理端/内测检索 API 验收，不进模型。
- ~~`user_id` query 鉴权~~：**P1** 起 stream/profile 走 **JWT 或游客 `X-Anon-ID`**；废弃 query `user_id` 作租户键。知识 admin 需 JWT + Admin。
- ~~游客 profile~~：**P1** 游客可读写 **关系档案**（key=`anon:{uuid}`）；登录后档案合并属 P2。
- **NextChat** `/api/*` 代理：P1 不做，仍归里程碑 B（可顺延 P2）。

## 架构决策

- `docs/adr/0001-frontend-backend-split.md` — 前后端分离与 `/api/*`
- `docs/adr/0002-backend-one-eino-alignment.md` — backend 平台对齐与 P1/P2 边界

## 本地开发

```text
backend/   → :8080  （Go API）
frontend/  → :3000  （Next.js；/api 可经 rewrite 转发到 backend）
```
