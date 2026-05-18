# 分拣标签

技能使用五个规范分拣角色表述。本文件将这些角色映射到本仓库 issue 跟踪器中实际使用的字符串。

本仓库使用**本地 Markdown** issue：`Status:` 行的值使用下表「本仓库 tracker」列（与技能角色名 1:1）。**不要**对 GitHub 调用 `gh issue edit --add-label`。

| mattpocock/skills 中的角色 | 本仓库 tracker（`Status:` 值） | 含义 |
| -------------------------- | ------------------------------ | ---- |
| `needs-triage` | `needs-triage` | 维护者需评估 |
| `needs-info` | `needs-info` | 等待报告者补充信息 |
| `ready-for-agent` | `ready-for-agent` | 规格完整，可供 AFK Agent 领取 |
| `ready-for-human` | `ready-for-human` | 需人工实现 |
| `wontfix` | `wontfix` | 不予处理 |

当技能提到某角色（如「应用可供 AFK 领取的分拣标签」）时，在 issue 文件的 `Status:` 行写入表中对应的字符串。

issue 文件中的说明文字、分拣备注、Agent Brief 等**人类可读内容使用简体中文**；`Status:` 值保持上表英文标识符。

若实际词汇不同，可修改「本仓库 tracker」列；左侧技能角色名列请勿改动。
