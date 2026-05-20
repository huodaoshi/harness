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
默认聊天走 NextChat 兼容聊天代理（脚手架期仅 **托管单模型**）；历史存浏览器本地（IndexedDB）。
_Avoid_: 称为关怀模式或 wellness 路径；或在 Scaffold 期要求用户在前端配置多厂商 Key/模型。

**托管单模型（Hosted single model）**:
模型 ID 与上游凭证仅由 **backend** 配置（`ARK_*`）；前端不展示、不持久化多厂商/多模型选择；用户无法在 UI 切换 Provider 或填写 API Key。
_Avoid_: 称 Scaffold 需实现 ADR-0001 所列全部 `/api/{provider}` 代理。

**NextChat 聊天代理路径（Scaffold，临时）**:
客户端仍请求 **`/api/bytedance`**；Go 用 `ARK_*` 转发火山方舟；**不经 SafetyGate**。属 **兼容占位**，后续可改为统一路径（如 `/v1/chat` 或 `/api/chat`），不在 Scaffold 锁死。
_Avoid_: 文档把 `bytedance` 路径名当作长期领域术语；在 Scaffold 期要求 `/api/*` 与 Wellness 同等拦截。

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

**前端脚手架阶段（Frontend Scaffold）**:
当前优先级：将 **NextChat** 业务代码迁入 `frontend/` 并接通 `backend` `/api/*`，先跑通通用聊天壳；**产品化**（关怀首屏、洪峰 UX、豆包式定制）后置优化。
_Avoid_: 在脚手架未完成前以 PRD 洪峰流程作为前端阻塞验收。

**前端第一期（Frontend P1）**:
（暂缓细化）脚手架之后的关怀向交付：洪峰入口、Wellness 会话路径、档案与危机 UX。具体范围待 Scaffold 完成后再定。
_Avoid_: 与 **Frontend Scaffold** 混为同一里程碑。

**Frontend Scaffold 实现策略**:
从 `NextChat/` **按目录迁移**至 `frontend/`（见 `frontend/MIGRATION.md`）；不复制 `app/api/`、`src-tauri/`；`NextChat/` 只读对照。范围 **仅浏览器 Web**（不含 Tauri 桌面/原生 App、不含 NextChat `export` 静态 App 包）。
_Avoid_: 在 `frontend/` 内从零重做聊天 UI 替代 NextChat 迁移（2026-05-20 已否决）；在 Scaffold 期承诺与上游 NextChat 同等多端。

**Frontend Scaffold 验收（最小 E2E）**:
`frontend` dev 可开聊天页；`GET /api/config` 下发单模型 + `hideUserApiKey`；经 rewrite **一条**聊天代理流式一问一答（模型由服务端决定）。Mask 可选，不阻塞。
_Avoid_: 验收多厂商代理、前端模型选择器或用户自填 Key。

**Frontend Scaffold 实施顺序**:
竖切 **C**：① `GET /api/config` 前后端通 → ② 迁聊天核心 + 豆包一条流式消息 → ③ 补 Mask / locales 等。非「先全量迁前端」或「先后端全部代理」。
_Avoid_: 在未完成 ① 时 bulk 复制整个 `NextChat/app/`。

**NextChat 代理用 LLM 配置**:
脚手架与 wellness **共用** `ARK_API_KEY`、`ARK_MODEL_ID`、`ARK_BASE_URL`（及既有 `LLM_*` 覆盖规则）。**不**再引入独立的 `BYTEDANCE_*` 作为正式配置名。
_Avoid_: 文档要求同时配置 `ARK_*` 与 `BYTEDANCE_*`。

**NextChat 访问口令（CODE）**:
未配置 `CODE` 时 `/api/config` 返回 `needCode=false`（本地脚手架默认可直接进入）。配置 `CODE` 后启用口令校验（与 NextChat 行为一致）。
_Avoid_: 本地脚手架强制要求口令。

**NextChat 设置 UI（Scaffold）**:
多厂商 / API Key / 模型选择 **仅通过 `/api/config` 隐藏**（如 `hideUserApiKey`、单条 `customModels`），**不**在迁移时删除相关前端代码。
_Avoid_: Scaffold 验收以「删掉 `platforms/*`」为必要条件。

**里程碑 B / 里程碑 C**:
见 ADR-0001：原 B = NextChat 聊天主路径 + 部分 Go `/api/*`；C = 全厂商代理、插件、关怀 UI 等。**Frontend P1 已 supersede B 的前端验收**（见 Flagged）。

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
- ~~Frontend P1 关怀优先~~：**已调整**——当前以 **Frontend Scaffold**（NextChat Web 壳 + 最小 Go `/api/*`）为先；关怀产品化属 **Frontend P1**。
- **ADR-0001 多厂商代理**：Scaffold **不**按 ADR 全量实现；以 **托管单模型** + 临时 `/api/bytedance` 为准（见 Language）。ADR 正文待实现时修订或加 supersede 说明。
- **`web/` 删除时机**：NextChat 脚手架 E2E（`frontend` dev + `backend` 流式一条消息）通过后再删；删除前可暂由 backend 继续托管。
- ~~`BYTEDANCE_*` 环境变量~~：已 supersede——正式配置统一为 **`ARK_*`**（YAML `llm:` 或进程环境变量，见 `backend/config/app/README.md`）。
- **backend 不用 `.env` 文件**：`config/app/*.yaml` + 可选 `*.secrets.yaml` + 环境变量覆盖。
- ~~ADR-0001 多厂商 `/api/*`~~：**Scaffold 已收窄**——仅 **托管单模型** + 必要 `GET /api/config`；其余 provider 路由**不在** Scaffold 范围。
- **NextChat 代理路径演进**：Scaffold 用 `/api/bytedance`（见上）；替换目标路径 **待 Scaffold 后** 再定。
- **Scaffold 期 `/api/bytedance` 与 SafetyGate**：**不过 Gate**（透明代理）；危机/医疗拦截仍在 **Wellness 会话路径**。换统一聊天路径时再决定是否合并 Gate。

## 架构决策

- `docs/adr/0001-frontend-backend-split.md` — 前后端分离与 `/api/*`
- `docs/adr/0002-backend-one-eino-alignment.md` — backend 平台对齐与 P1/P2 边界

## 本地开发

```text
backend/   → :8080  （Go API）
frontend/  → :3000  （Next.js；/api 可经 rewrite 转发到 backend）
```
