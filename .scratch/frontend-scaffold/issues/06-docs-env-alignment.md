# Scaffold-06：文档与环境变量对齐

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

与 Scaffold 共识对齐文档与示例配置：

- **`backend/.env.example`**：移除或标注废弃 `BYTEDANCE_*`；NextChat 代理仅文档化 **`ARK_*`**、`CODE`、`CUSTOM_MODELS` / `DEFAULT_MODEL`（若由 backend 读取）
- **`frontend/MIGRATION.md`**：更新第一期范围为托管单模型、竖切顺序、不过 SafetyGate、仅 Web
- **`docs/adr/0001-frontend-backend-split.md`**：增加简短 **Supersede 说明** 或链接到 `CONTEXT.md`（不强制重写全文）
- **`frontend/README.md`**：本地启动步骤（`backend` + `frontend`）

## 验收标准

- [ ] 新开发者仅读 README + `.env.example` 能启动 #03 的 E2E
- [ ] 文档不再要求实现多厂商 `/api/*` 作为 Scaffold 验收
- [ ] ADR 或 CONTEXT 交叉引用一致

## 阻塞于

- [03-nextchat-chat-core-e2e.md](./03-nextchat-chat-core-e2e.md)

## 覆盖的用户故事

无（文档）
