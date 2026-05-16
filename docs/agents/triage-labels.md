# 分拣标签（Triage）

技能内部使用五个标准分拣**角色名**。本文件将这些角色映射到本仓库 issue tracker 中实际使用的字符串。

本仓库使用**本地 Markdown** issue：`Status:` 行的值使用下表「本仓库 tracker」列（与技能角色名 1:1）。**不要**对 GitHub 调用 `gh issue edit --add-label`。

| 技能中的角色名 | 本仓库 tracker（`Status:` 值） | 含义 |
| -------------- | ------------------------------ | ---- |
| `needs-triage` | `needs-triage` | 维护者需要先评估 |
| `needs-info` | `needs-info` | 等待报告人补充信息 |
| `ready-for-agent` | `ready-for-agent` | 规格完整，可供 AFK agent 接手 |
| `ready-for-human` | `ready-for-human` | 需要人工实现 |
| `wontfix` | `wontfix` | 不予处理 |

当技能提到某个角色（例如「应用 AFK-ready 分拣标签」）时，在 issue 文件的 `Status:` 行写入表中对应的字符串。

issue 文件中的说明文字、分拣备注、agent brief 等**人类可读内容使用简体中文**；`Status:` 值保持上表英文标识符，以便与 mattpocock/skills 一致。

若你实际使用的词汇不同，可修改「本仓库 tracker」列；技能角色名列请勿改动。
