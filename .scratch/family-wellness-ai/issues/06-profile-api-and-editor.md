# 纵向切片：关系档案 API + Web 编辑页

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

暴露 `GET /v1/profile`、`PUT /v1/profile`（全量 upsert）：字段含 `self`、`people[]`、`current_issue`（snake_case JSON）。Web 增加档案编辑页：保存后返回聊天；聊天无需档案亦可开始。

保存后下一条消息应能在 #03 的 inject 中体现新档案（与 Graph 联调）。

## 验收标准

- [ ] GET 无档案时返回空对象或 404 语义明确，不崩溃
- [ ] PUT 后 GET 回读一致
- [ ] Web 表单可增删「重要他人」条目
- [ ] 集成测试：PUT → stream 消息 → 断言 context 含新 `current_issue`
- [ ] 鉴权与 `user_id` 隔离（可与 #12 简化 token 并存）

## 阻塞于

- [03-spike-s3-profile-inject.md](./03-spike-s3-profile-inject.md)

## 覆盖的用户故事

#4、#5、#6
