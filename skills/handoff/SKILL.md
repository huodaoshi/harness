---
name: handoff
description: 将当前对话压缩成交接文档，供其他 Agent 接续工作。
argument-hint: "下一会话将用于什么？"
---

撰写交接文档，概括当前对话，使全新 Agent 能继续工作。保存到 `mktemp -t handoff-XXXXXX.md` 生成的路径（写入前先读取该文件）。

若下一会话需用技能，在文档中建议。

勿重复 PRD、计划、ADR、issue、提交、diff 等已有产物中的内容；用路径或 URL 引用即可。

若用户传入参数，视为下一会话重点描述，据此定制文档。
