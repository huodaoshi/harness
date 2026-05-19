# 删号与数据级联删除

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — Spike S7

## 要构建什么

实现 `DELETE /v1/account`（或等价）：按 `user_id` 删除 `relationship_profiles`、`sessions`、`session_summaries`、`audit_events`（若存在）。返回明确成功/失败；Web 可提供设置页入口（可选）。

提供自动化测试：种子多集合数据 → 删除 → 断言全部为空。

## 验收标准

- [ ] 删除后 GET profile 与 stream 均无法访问该用户历史
- [ ] 自动化级联测试通过
- [ ] 隐私说明文档列出删除范围与备份保留策略（占位可「待法务」）
- [ ] 误删需二次确认（Web 或 API 文档约定）

## 阻塞于

- [06-profile-api-and-editor.md](./06-profile-api-and-editor.md)
- [07-session-persist-and-summary-card.md](./07-session-persist-and-summary-card.md)

## 覆盖的用户故事

#14
