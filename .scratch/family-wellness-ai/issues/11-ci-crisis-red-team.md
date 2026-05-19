# CI 危机回归 + 红队 50 条门禁

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — S4 / 上线门禁

## 要构建什么

将 #02 的 10 条危机剧本 + #05 关键用例纳入 **CI 必跑**（PR 触发，目标 ≤2min）。扩展 **红队 50 条** 数据集（含 prompt 注入样本），本地/CI 可跑报告：穿透率、危机漏判率。

门禁目标：危机 **0 漏判**；穿透 ≤5%（超标不阻塞合并但报告标红，由维护者解读）。

## 验收标准

- [ ] CI workflow 在 PR 上运行危机子集
- [ ] 50 条红队可一键运行并输出 JSON/Markdown 摘要
- [ ] 危机漏判任一则 CI 失败
- [ ] 文档说明如何追加新剧本

## 阻塞于

- [05-safety-medical-block-audit.md](./05-safety-medical-block-audit.md)

## 覆盖的用户故事

#17、#18
