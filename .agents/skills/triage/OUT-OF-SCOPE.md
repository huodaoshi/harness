# 范围外知识库

仓库中的 `.out-of-scope/` 目录存放已拒绝功能请求的持久记录，用途有二：

1. **组织记忆** — 记录功能为何被拒绝，避免 issue 关闭后理由丢失
2. **去重** — 新进 issue 若与既往驳回同类，可点出先前决策，避免重复争论

## 目录结构

```
.out-of-scope/
├── dark-mode.md
├── plugin-system.md
└── graphql-api.md
```

每个**概念**一个文件，而非每个 issue 一个文件。多个 issue 请求同一事物时，归在同一文件下。

## 文件格式

宜采用轻松可读的风格——更像短设计说明，而非数据库条目。用段落、代码示例阐明理由，便于首次阅读者理解。

```markdown
# Dark Mode

本项目不支持暗色模式或面向用户的主题切换。

## 为何在范围外

渲染管线假定 `ThemeConfig` 中定义单一调色板。支持多主题需要：

- 包裹整棵组件树的主题 Context Provider
- 各组件按主题解析样式
- 用户主题偏好的持久层

这是与项目「内容创作」重心不符的重大架构变更。主题化应由嵌入或再分发产出的下游消费者负责。

```ts
// 当前 ThemeConfig 未设计运行时切换：
interface ThemeConfig {
  colors: ColorPalette; // 单一调色板，构建时解析
  fonts: FontStack;
}
```

## Prior requests

- #42 — "Add dark mode support"
- #87 — "Night theme for accessibility"
- #134 — "Dark theme option"
```

### 文件命名

用简短、描述性的 kebab-case 概念名：`dark-mode.md`、`plugin-system.md`、`graphql-api.md`。浏览目录时应能一眼看出被拒概念，无需打开文件。

### 撰写理由

理由应有实质——不是「我们不想要」，而是**为何**。好的理由可引用：

- 项目范围或理念（「本项目聚焦 X；主题是下游关切」）
- 技术约束（「支持此项需要 Y，与 Z 架构冲突」）
- 战略决策（「因 … 选择 A 而非 B」）

理由应耐久。避免引用临时情况（「现在太忙」）——那是延期，不是真正的拒绝。

## 何时检查 `.out-of-scope/`

分拣时（步骤 1：收集上下文），读取 `.out-of-scope/` 下全部文件。评估新 issue 时：

- 检查请求是否匹配既有范围外概念
- 按概念相似度匹配，非关键词——「night theme」可匹配 `dark-mode.md`
- 若匹配，告知维护者：「这与 `.out-of-scope/dark-mode.md` 类似——此前因 [理由] 拒绝。您是否仍持相同看法？」

维护者可：

- **确认** — 将新 issue 追加到既有文件的 Prior requests，然后关闭
- **重新考虑** — 删除或更新范围外文件，issue 走正常分拣
- **不同意** — issue 相关但不同，走正常分拣

## 何时写入 `.out-of-scope/`

仅当 **enhancement**（非 bug）以 `wontfix` 拒绝时。流程：

1. 维护者认定功能请求在范围外
2. 检查是否已有匹配的 `.out-of-scope/` 文件
3. 若有：将新 issue 追加到 Prior requests
4. 若无：以概念名新建文件，含决策、理由与首条 prior request
5. 在 issue 评论说明决策并提及 `.out-of-scope/` 文件
6. 以 `wontfix` 标签关闭 issue

## 更新或删除范围外文件

若维护者改变对既往拒绝概念的看法：

- 删除对应 `.out-of-scope/` 文件
- 技能无需重新打开旧 issue——它们是历史记录
- 触发重新考虑的新 issue 走正常分拣
