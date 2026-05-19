# PRD：关系情绪 AI — MVP 实现规格

**Status:** ready-for-agent  
**特性：** family-wellness-ai  
**版本：** MVP（C1 + C2 最小 + C3）  
**最后更新：** 2026-05-19  
**上游文档：** [项目建议书](./项目建议书.md)、[核心用户](./核心用户.md)、[references/](./references/README.md)

---

## 问题陈述

关注原生家庭与亲密关系的用户，在吵架、被催婚、翻旧账等**情绪洪峰**时，需要立刻倾诉与被理解。通用对话 AI（如豆包）往往「说得对但接不住具体关系脉络」；泛心理 App 常把关系场景套在测试、玄学或社区框架里。多数用户尚未准备好付费真人咨询（羞耻、价格、时间），但仍需要**安全、非评判、记得关系上下文**的 AI 倾诉出口。

## 解决方案

提供垂直于「关系情绪」的 **Web 优先 AI 倾诉伙伴**：一键进入洪峰模式、流式对话体验对齐主流 AI 产品；用户可编辑**关系档案**（我 / 重要他人 / 当前议题），系统注入档案与**最近 1 场会话摘要**（记忆 Level 1–2）；**SafetyGate** 在进模型前拦截危机与违规，自伤/家暴/医疗边界走固定流程且**不调用 LLM**。明确非医疗、非玄学、非虚拟恋爱。

---

## 用户故事

1. As a **情绪洪峰中的用户**, I want to tap「我现在很难受」and start a session immediately, so that I can vent without navigating menus or explaining the product first.
2. As a **用户**, I want the AI to respond in a streaming fashion like familiar chat products, so that I feel listened to in real time rather than waiting for a block of text.
3. As a **用户**, I want the AI to acknowledge my relationship context (family, partner, boundaries), so that replies feel grounded in my situation rather than generic advice.
4. As a **用户**, I want to create and edit a **relationship profile** (about me, important others, current issue), so that I control what the system remembers about my relationships.
5. As a **用户**, I want to skip profile setup and still chat, so that I am not blocked on my worst night.
6. As a **用户**, I want the system to use my profile and the **last session summary** in new conversations, so that I do not have to repeat the same backstory every time.
7. As a **用户**, I want to see a short **3-sentence session recap** when a conversation ends (optional card), so that I can leave with a light sense of closure without a full therapy module.
8. As a **用户**, I want clear **disclaimers** that this is not medical care or diagnosis, so that I understand the product boundary.
9. As a **用户 expressing self-harm intent**, I want **fixed crisis guidance** (empathy, hotlines, encourage human help) without AI-generated clinical advice, so that I am not harmed by model hallucination in a crisis.
10. As a **用户 describing domestic violence or imminent harm**, I want a **crisis branch** with safety resources, so that the product does not role-play through danger.
11. As a **用户 asking for medication or diagnosis**, I want a **fixed boundary response** (MVP: template, no LLM), so that the product does not practice medicine.
12. As a **用户 sending prohibited content**, I want the message blocked with a clear refusal, so that the service stays within policy.
13. As a **用户**, I want my crisis interactions **not logged in plain text** in operational logs, so that I can trust privacy on my worst day.
14. As a **用户**, I want to **delete my account and associated data** (profile, sessions, summaries), so that I can exercise privacy rights.
15. As a **内测用户**, I want to rate whether a session felt「更贴关系」than a generic AI, so that the team can validate the wedge (H3/H4).
16. As a **运营/客服（未来）**, I want audit events for gate outcomes without raw message bodies, so that we can investigate incidents compliantly.
17. As a **开发者**, I want W1 spikes (stream, crisis zero-LLM, profile inject) to pass before feature expansion, so that we do not build on a broken orchestration stack.
18. As a **开发者**, I want CI to run a **crisis regression subset** on each PR, so that safety refactors cannot silently remove the crisis branch.
19. As a **用户 on mobile browser**, I want a usable Web/PWA layout for chat and the distress button, so that I do not need an App Store install for MVP validation.
20. As a **用户**, I want **rate limiting** so that abuse does not exhaust API quota, so that the service remains available for others.
21. As a **用户 in 洪峰模式**, I want the AI tone to prioritize listening and emotion naming over lecturing, so that the product matches「隐忍型高敏感」expectations.
22. As a **用户**, I want **普通聊天模式** when not in distress, so that I can continue lighter conversations without forced crisis UX.
23. As a **用户**, I want session history for the **current session** visible in the UI, so that I can read back what was said tonight.
24. As a **用户**, I want errors (network, provider, gate block) surfaced clearly in the chat UI, so that I know whether to retry or stop.
25. As a **合规评审者**, I want hotline numbers and crisis copy **versioned in config**, so that they can be updated without code deploy ambiguity.

---

## 实现决策

### 深模块划分（建议实现边界）

| 模块 | 职责 | 对外接口（概念） | 深/浅 |
| ---- | ---- | ---------------- | ----- |
| **SafetyGate** | L1 规则 + L2 云审；输出 gate 枚举 | `Evaluate(ctx, userMessage) → GateResult` | **深** |
| **CrisisResponder** | 危机/医疗/拦截固定文案 | `Respond(gateResult) → FixedContent` | 深 |
| **RelationshipSessionGraph** | Eino 编排：Gate → Router → Inject → Model → PostProcess | `RunStream(ctx, SessionInput) → StreamReader` | **深** |
| **ProfileStore** | 关系档案 CRUD | `GetProfile`, `UpsertProfile`, `DeleteByUser` | 深 |
| **SessionStore** | 会话元数据与消息追加（含轮次 cap） | `CreateSession`, `AppendMessage`, `GetSession` | 深 |
| **SummaryStore** | 会话摘要读写；最近 N=1 | `SaveSummary`, `GetLatestForUser` | 浅（可并入 Session 域） |
| **ChatModelGateway** | 国内 LLM API + failover 接口 | `Stream(ctx, messages) → token stream` | 深 |
| **SessionHTTP** | Hertz：鉴权、限流、SSE 映射 | `POST /v1/sessions/stream` 等 | 中 |
| **WebClient** | 洪峰入口、聊天、档案编辑、免责声明 | 调用 Session HTTP + Profile API | 中 |

**依赖方向：** WebClient → SessionHTTP → Graph → {SafetyGate, Stores, ChatModelGateway}；Graph **编译期**不将 crisis/medical/block 边接到 ChatModelGateway。

### 架构与栈

- **HTTP：** Hertz；流式会话单一主入口 `POST /v1/sessions/stream`（SSE）。
- **编排：** Eino `compose.Graph`（关系会话流水线）；遵循仓库 `eino.md` rule：Graph 用于分支会话，非手写 goroutine 拼 SSE。
- **存储：** MongoDB **自建**（同 VPC）；Redis 限流与会话热数据。
- **记忆 MVP：** Level 1 用户档案 + Level 2 最近 **1** 场 `summary3`；无向量 RAG。
- **模型：** 火山方舟 / 豆包等国内 API；`ChatModel` 抽象预留第二供应商。

### SafetyGate 状态机（来自技术专项）

```text
Input(message)
  → L1 RuleGate
  → crisis_self_harm | crisis_violence → CrisisBranch (NO LLM)
  → medical_boundary → FixedTemplate (NO LLM, MVP)
  → block → Refusal (NO LLM)
  → pass → L2 CloudModeration (staging/prod: fail-closed)
  → pass → ModeRouter
```

- 规则外置 YAML；危机文案外置 JSON  locale `zh-CN`。
- `audit_events`（可选）：记录 `gate_result`, `session_id`, `timestamp`；**不存用户原文**（至多 hash / 20 字截断）。

### ModeRouter

- **洪峰模式：** 系统 prompt 偏倾听、情绪命名、关系议题；由客户端 `mode=distress` 或等价字段触发。
- **普聊模式：** 默认较轻语气；仍注入档案与摘要。
- MVP 不做自动模式检测（无 ML 分类器）。

### RelationshipSessionGraph 节点

1. **SafetyGate** — 见上。  
2. **ModeRouter** — 选择 prompt 模板集。  
3. **ProfileInject** — 从 Mongo 读取 `relationship_profiles` + 最近 1 条 `session_summaries`，组装 prompt 块。  
4. **ChatModel** — 流式 token；仅 `pass` 路径。  
5. **PostProcess** — 会话结束时生成 `summary3`（3 句）；可同步或异步；MVP 允许「仅洪峰会话生成」。

### MongoDB 文档形状（MVP）

**relationship_profiles**

```json
{
  "user_id": "string",
  "self": { "note": "string" },
  "people": [{ "label": "string", "relation": "string", "note": "string" }],
  "current_issue": "string",
  "updated_at": "ISODate"
}
```

**sessions**

```json
{
  "session_id": "string",
  "user_id": "string",
  "mode": "distress|normal",
  "gate_results": ["pass"],
  "messages": [{ "role": "user|assistant", "content": "string", "at": "ISODate" }],
  "created_at": "ISODate"
}
```

- 单会话 **消息轮次硬 cap**（建议 50）；防大文档。
- 消息亦可拆 `session_messages` 集合（实现自选）；PRD 不强制，但须满足 cap 与按 user 删除。

**session_summaries**

```json
{
  "session_id": "string",
  "user_id": "string",
  "summary3": ["string", "string", "string"],
  "created_at": "ISODate"
}
```

JSON 字段命名：**snake_case**（见仓库 naming rule）。

### HTTP / SSE 契约

**创建/继续流式会话**

- `POST /v1/sessions/stream`
- Request（JSON）：`session_id?`, `user_id`（或从 auth 推导）, `mode`, `message`
- Response：`text/event-stream`

| SSE event | payload 要点 | 何时 |
| --------- | ------------ | ---- |
| `token` | `{ "text": "..." }` | 模型流式片段 |
| `done` | `{ "session_id": "..." }` | 正常结束 |
| `crisis` | `{ "template_id": "...", "body": "..." }` | Gate 危机支路 |
| `error` | `{ "code": "...", "message": "..." }` | 阻断、_provider 失败、限流 |

**关系档案**

- `GET /v1/profile` — 当前用户档案  
- `PUT /v1/profile` — 全量 upsert（MVP 可不做 PATCH 粒度）  
- `DELETE /v1/account` — 级联删 profile、sessions、summaries、audit（Spike S7）

**鉴权 MVP：** 允许简化（如内测 token / 魔法链接）；须能在多租户逻辑上隔离 `user_id`；具体机制实现时选定，不阻塞 W1 spike。

### Web 客户端交互

- 首屏：免责声明勾选 + 「我现在很难受」主 CTA + 次要「先看看 / 编辑档案」。
- 聊天页：消息列表 + 输入框；SSE 驱动 assistant 气泡增量渲染。
- 危机事件：替换为不可继续输入的危机卡片（热线、求助指引）；**不**展示模型续写。
- 档案页：表单编辑 self / people[] / current_issue；可保存后返回聊天。
- 可选：会话结束展示 `summary3` 卡片（可关闭）。

### 配置与运维

- 环境：dev / staging / prod；staging 起接真实云审。
- Mongo：**自建**；`mongodump` 定时备份；监控磁盘与连接。
- CI：PR 运行危机剧本子集（S2 核心用例）。
- LLM 与云审密钥：环境变量 / 密钥管理，禁止入库。

### W1 Spike 门禁（先于大规模功能开发）

| ID | 必须通过 |
| ---- | -------- |
| S1 | Hertz → Eino → 假 ChatModel → SSE `token` 流 |
| S2 | 10 危机剧本：`crisis` 事件 + ChatModel **调用次数 0** |
| S3 | 档案 + 最近摘要 inject：20 条自动化 |

S4 红队 50 条不阻塞 W1 启动，阻塞 MVP 上线门禁。

### 命名与跨层约定

- API / Mongo：**snake_case**
- Go 字段：**camelCase** + json tag snake_case
- TypeScript（Web）：**camelCase** 属性，边界映射 API

---

## 测试决策

### 何谓好测试

- 测**对外行为与契约**：HTTP 状态、SSE 事件序列、Gate 枚举、Mongo 读写结果、**危机路径零 LLM 调用**。
- **不测** Graph 内部节点调用顺序、私有 prompt 拼接字符串全文（除非 golden 快照用于回归且稳定）。
- 使用 **fake / mock ChatModel** 与 **fake Moderation** 做确定性集成测试。

### 必须覆盖的模块

| 模块 | 测试类型 | 要点 |
| ---- | -------- | ---- |
| SafetyGate | 表驱动单元 + 契约 | 10 危机剧本、医疗模板、block、误杀样本（后续扩） |
| RelationshipSessionGraph | 集成 | 假模型；assert SSE；crisis 不挂载 model |
| SessionHTTP | 集成 | SSE 解析；错误码；限流 429 |
| ProfileStore / SessionStore | 集成（Mongo test container 或 ephemeral） | CRUD、cap、级联删 |
| ChatModelGateway | 契约 + 可选录制 | 超时、failover 切换（staging） |
| WebClient | 组件/E2E（可选 MVP 后期） | 危机卡片的输入禁用 |

### 先例

- 本仓库 `cli/tests/`：**表驱动**、子进程/HTTP 黑盒测试风格，可类比后端 API 测试组织方式。
- 危机回归：**CI 必跑**子集，与 Spike S2 用例同源。

### 不建议 MVP 做的测试

- 端到端依赖真实 LLM 的 CI（贵、 flaky）；放 nightly 或手动。
- 向量检索、支付、推送。

---

## 范围外

- 工具库、冥想库、心理测试、玄学、虚拟恋爱、UGC 社区、真人咨询撮合。
- 付费、订阅、IAP（商业会议后再议）。
- 向量全文 RAG、摘要向量检索（记忆 Level 3–4）。
- 实时语音通话、微信小程序为主战场、App 上架（v1.1）。
- 自动诊断、用药建议、病因判定、Replika 式无边界顺从。
- Kitex 微服务拆分、自研基座模型、运营可读全文永久存档。
- v1.1 能力：多关系人强化、近 3 次摘要、轻复盘模块、男性/妈妈分层运营。

---

## 补充说明

### 与《项目建议书》关系

本 PRD 为实现规格；商业定价、记忆 Level 3–4 细化、云审/LLM 厂商合同以建议书「待议」项为准，**不阻塞** MVP 编码，但须在上线前完成 S6/S8/S9。

### 上线门禁（MVP 退出）

- 红队穿透 ≤5%，危机剧本 **100%**（含零 LLM）。  
- 隐私评审通过（含删号、日志策略）。  
- Web 可测链接；内测 n≥30 的留存/盲测为**观察指标**，不作对外承诺。

### Agent 领取说明

1. 先完成 **S1–S3** spike，再纵向切片实现。  
2. 新建代码预期在 **`backend/`** 与 **`web/`**（或等价目录）；遵守 `.cursor/rules/eino.md`、`naming.md`、`git.md`。  
3. 实现 issue 见 [`issues/`](./issues/)（`01`–`13` 纵向切片，2026-05-19 发布）。

### 成功指标（内测参考）

- ≥60% 盲测认为比「关系向 prompt 的豆包」更贴关系。  
- 7 日留存正向信号（n≥30）。  
- 定性：用户愿下周再来 / 愿推荐给密友。

---

## 评论

### 2026-05-19 · to-prd

由 product-council 技术专项（4 轮）与《项目建议书》整理发布；`Status: ready-for-agent`。

### 2026-05-19 · to-issues

已拆分为 `issues/01`–`13`；建议领取顺序：01 → 02 → 03 → 04 → 05 → 06 → 07 → …；#08 为 `ready-for-human`（LLM 厂商配置）。
