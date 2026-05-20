#!/usr/bin/env bash
# 人在环中的复现循环。
# 复制本文件，编辑下方步骤后运行。
# Agent 运行脚本；用户在终端按提示操作。
#
# 用法：
#   bash hitl-loop.template.sh
#
# 两个辅助函数：
#   step "<说明>"          → 显示说明，等待 Enter
#   capture VAR "<问题>"   → 显示问题，读入回答到 VAR
#
# 结束时以 KEY=VALUE 打印捕获值，供 Agent 解析。

set -euo pipefail

step() {
  printf '\n>>> %s\n' "$1"
  read -r -p "    [完成后按 Enter] " _
}

capture() {
  local var="$1" question="$2" answer
  printf '\n>>> %s\n' "$question"
  read -r -p "    > " answer
  printf -v "$var" '%s' "$answer"
}

# --- 在下方编辑 ---------------------------------------------------------

step "在浏览器打开 http://localhost:3000 并登录。"

capture ERRORED "点击「Export」按钮。是否抛出错误？(y/n)"

capture ERROR_MSG "粘贴错误信息（或填 'none'）："

# --- 在上方编辑 ---------------------------------------------------------

printf '\n--- 已捕获 ---\n'
printf 'ERRORED=%s\n' "$ERRORED"
printf 'ERROR_MSG=%s\n' "$ERROR_MSG"
