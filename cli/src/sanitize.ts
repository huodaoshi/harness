/**
 * 在输出到终端前清理不可信字符串。
 *
 * 剥离字符串中的全部终端转义序列，包括：
 *   - CSI 序列  (ESC [ ... final_byte)    — 光标移动、清屏、SGR 颜色
 *   - OSC 序列  (ESC ] ... BEL/ST)         — 窗口标题、超链接
 *   - 简单转义 (ESC + 单字符)              — 如 ESC 7（保存光标）
 *   - C1 控制码 (0x80–0x9F)
 *   - 原始控制字符 (BEL、BS 等)            — 保留安全的 \t 与 \n
 *
 * 用于防御 CWE-150（终端转义注入）：不可信数据（如 SKILL.md frontmatter
 * 或远程 API 中的技能名/描述）可能清屏、移动光标、改窗口标题，
 * 或渲染看似合法 CLI 输出的攻击者文本。
 */

// CSI：ESC[ + 参数字节 (0x30-0x3F) + 中间字节 (0x20-0x2F) + 结束字节 (0x40-0x7E)
const CSI_RE = /\x1b\[[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]/g;

// OSC：ESC] ... 以 BEL (\x07) 或 ST (ESC\) 结束
const OSC_RE = /\x1b\][\s\S]*?(?:\x07|\x1b\\)/g;

// DCS、PM、APC：ESC P|^|_ ... 以 ST (ESC\) 结束
const DCS_PM_APC_RE = /\x1b[P^_][\s\S]*?(?:\x1b\\)/g;

// 简单两字节转义：ESC + 0x20-0x7E 单字符
// 含 ESC 7 (DECSC)、ESC 8 (DECRC)、ESC c (RIS)、ESC M (RI) 等
const SIMPLE_ESC_RE = /\x1b[\x20-\x7e]/g;

// C1 控制码 (0x80-0x9F) — ESC 序列的 8 位等价形式
const C1_RE = /[\x80-\x9f]/g;

// 除制表符 (\x09)、换行 (\x0a) 外的原始控制字符
// 含 BEL (\x07)、BS (\x08)、CR (\x0d) 等
const CONTROL_RE = /[\x00-\x06\x07\x08\x0b\x0c\x0d-\x1a\x1c-\x1f\x7f]/g;

/**
 * 剥离全部终端转义序列及危险控制字符。
 *
 * 在将不可信输入打印到终端前使用。
 */
export function stripTerminalEscapes(str: string): string {
  return str
    .replace(OSC_RE, '') // 先处理 OSC（最长匹配）
    .replace(DCS_PM_APC_RE, '') // DCS/PM/APC
    .replace(CSI_RE, '') // CSI
    .replace(SIMPLE_ESC_RE, '') // ESC+字符
    .replace(C1_RE, '') // C1
    .replace(CONTROL_RE, ''); // 原始控制字符（保留 \t \n）
}

/**
 * 清理技能元数据字符串（名称、描述等）以便安全显示在终端。
 *
 * 除剥离转义外，还会 trim，并将内部换行折叠为空格
 * （列表中的名称/描述应为单行）。
 */
export function sanitizeMetadata(str: string): string {
  return stripTerminalEscapes(str)
    .replace(/[\r\n]+/g, ' ')
    .trim();
}
