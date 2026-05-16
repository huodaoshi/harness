import { parse as parseYaml } from 'yaml';

/**
 * 精简的 frontmatter 解析器。仅支持 YAML（`---` 分隔符）。
 * 不支持 `---js` / `---javascript`，以避免 gray-matter 内置 JS 引擎
 * 中基于 eval() 的远程代码执行（RCE）风险。
 */
export function parseFrontmatter(raw: string): {
  data: Record<string, unknown>;
  content: string;
} {
  const match = raw.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n?([\s\S]*)$/);
  if (!match) return { data: {}, content: raw };
  const data = (parseYaml(match[1]!) as Record<string, unknown>) ?? {};
  return { data, content: match[2] ?? '' };
}
