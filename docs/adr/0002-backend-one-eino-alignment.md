# ADR-0002：Backend 对齐 one-eino 平台架构

- **状态**：已接受
- **日期**：2026-05-20
- **决策者**：产品与工程（grill-with-docs 共识）
- **关联**：`CONTEXT.md`、`docs/adr/0001-frontend-backend-split.md`；参考实现 `D:\one-eino\backend`（复制改写，非 runtime 依赖）

## 背景

`harness/backend/` 当前为 family-wellness-ai MVP：扁平 `internal/*`、`POST /v1/sessions/stream`（SafetyGate + RelationshipSessionGraph）、`?user_id=` 鉴权、无知识库与正式用户体系。

`one-eino/backend` 为已落地的 Clean Architecture 平台： `modules/*` 三层、`api/` + `infra/` + `conf/`、JWT/游客、知识 Space 与入库 Worker、ADK `/v1/chat`（游戏垂直，含 score/tag 等）。

目标：将 harness backend **重构为 one-eino 式平台骨架**，纳入**本产品**知识库与用户鉴权；**去掉游戏垂直**；**保留并分期演进**关怀路径（SafetyGate、关系档案、summary3）。与 ADR-0001 正交：0001 管前后端分离与 `/api/*` 代理；本 ADR 管 Go 后端内部形态与能力边界。

## 决策

### 1. 范围：平台对齐，非游戏业务移植

- **做**：目录分层、auth、knowledge + knowledgeindexing、共享 infra/pkg；知识内容为产品与关系情绪领域（政策、FAQ、话术示例等）。
- **不做（P1）**：游戏向入库 LLM 评分/打标；`adminorg`；微信/扫码登录（模块可搬，配置默认关）；A2UI `/v1/chat`；ADR-0001 中尚未完成的 NextChat `/api/*` 全量代理。

### 2. 代码集成方式

从 `one-eino/backend` **复制改写到** `harness/backend`，`go.mod` 模块路径为 `github.com/huodaoshi/harness/backend`。one-eino 仓库保持参考实现，**不**采用 go workspace 双仓联动或 P1 抽取公共 module。

### 3. 对话编排：分两期（D → B）

| 期 | 用户聊天 | 编排 |
|----|----------|------|
| **P1** | 仍为 `POST /v1/sessions/stream`（SSE 契约不变） | 现有 RelationshipSessionGraph + SafetyGate；危机/医疗/拦截**编译期不接 LLM** |
| **P2** | 统一 `POST /v1/chat`（ADK Runner） | SafetyGate 作为 ADK 前置或图内首节点；可选 A2UI；Wellness 路径并入平台对话路径 |

P1 **不**为提前合流而改动 stream _handler 与 crisis 回归测试面。

### 4. P1 交付包

- 目录：`conf/`、`infra/`、`api/`、`pkg/`、`modules/`（含自 `internal/*` 迁出的 **`modules/wellness`**）。
- **auth**：短信 SMS + JWT refresh；**游客** `X-Anon-ID`（`user_id`=`anon:{uuid}`）可访问 stream/profile；微信等 P1 配置关闭。
- **knowledge**：Space 机制 + admin 入库 API；**默认仅运营一个 Space**；`doc_type` 语义：**1=产品与边界，2=FAQ，3=运营话术/示例**（非攻略/设定/剧情）。
- **knowledgeindexing**：独立 Worker（`cmd/knowledgeindexing`）；本地 `mq.provider=local` 时可进程内嵌消费（对齐 one-eino）；流水线：**解析 → 分块 → assignID → 向量索引**，**无** per-chunk LLM score/tag。
- **鉴权废弃**：query `user_id` / `X-User-Id` 作为租户键；改为 JWT 或游客 header。
- **游客关系档案**：允许 `GET/PUT /v1/profile`（key=`anon:{uuid}`）；登录后档案合并 **P2**。

### 5. 知识库与用户对话（RAG）

- **P1**：知识库用于**运营入库与检索验收**；Wellness prompt **仍仅**关系档案（Level 1）+ 最近会话摘要 summary3（Level 2）；**不向** `/v1/sessions/stream` 注入向量 RAG。PRD「MVP 无向量 RAG」在「用户聊天路径」上仍成立，与「存在知识 Space」不矛盾。
- **P2**：经平台对话路径 / ADK tool 启用 RAG；具体 topK、Space 绑定在 P2 设计时定。
- **可选（P1）**：仅内测/管理端用的检索 HTTP API，不进模型——若工期允许再加，不阻塞 P1 主线。

### 6. 危机与配置文案

危机/医疗/拦截**固定文案**继续外置 `config/safety_rules_v1.yaml`、`config/crisis_templates/` 等（SafetyGate）；**不强制**迁入向量库。知识 Space 与危机模板职责分离。

### 7. 与 ADR-0001 的关系

- NextChat `/api/*` Go 代理、里程碑 B/C 范围**不变**；因 P1 聚焦平台骨架，**`/api/*` 实现顺延**至 P1 之后（可与 P2 或里程碑 B 并行排期）。
- `frontend/` 迁移不依赖本 ADR 完成，但 P1 鉴权变更后，关怀前端需改发 JWT/游客 header。

## 曾考虑的选项

| 选项 | 未选原因 |
|------|----------|
| 仅目录对齐（A），不搬 auth/knowledge | 无法满足「本产品知识库」与多用户隔离 |
| P1 即统一 `/v1/chat`（跳过 D） | 同时改编排 + 鉴权 + 知识库，危机回归与前端风险过大 |
| go workspace 依赖 one-eino module | 双仓同步与游戏向代码泄漏成本高 |
| P1 在 Wellness Graph 内接 RAG | 违背 D→B；扰动已验证 SafetyGate 邻接逻辑 |
| 全量搬 one-eino 入库 score/tag | 游戏向、chunk 成本高，与运营选材策略不符 |

## 后果

### 正面

- 后端可与 one-eino 技能文档、模块边界对齐，降低后续搬 ADK/知识检索成本；
- P1 用户路径稳定，安全回归面可控；
- 知识库与关怀配置职责清晰。

### 负面 / 风险

- P1 为大范围文件迁移与接线，需分 PR 或 tracer bullet 纵向切片；
- 鉴权切换破坏现有 `?user_id=` 联调脚本与测试，需一并更新；
- 复制代码后需删减游戏/adminorg 引用，避免死代码与错误 doc_type 假设；
- P2 合流时存在短期双聊天入口，需在 CONTEXT 与 API 文档中标注弃用时间线。

## 参考

- 术语：`CONTEXT.md`
- one-eino 架构说明：`D:\one-eino\.claude\skills\backend\backend\01-architecture.md`
- 产品规格：`.scratch/family-wellness-ai/PRD.md`
- 入库精简：相对 `docs/specs/2026-05-07-知识入库Eino重构-设计.md`（one-eino），Harness P1 **不实现** score/tag BatchNode
