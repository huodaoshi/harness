export type AgentType =
  | 'aider-desk'
  | 'amp'
  | 'antigravity'
  | 'augment'
  | 'bob'
  | 'claude-code'
  | 'openclaw'
  | 'cline'
  | 'codearts-agent'
  | 'codebuddy'
  | 'codemaker'
  | 'codestudio'
  | 'codex'
  | 'command-code'
  | 'continue'
  | 'cortex'
  | 'crush'
  | 'cursor'
  | 'deepagents'
  | 'devin'
  | 'dexto'
  | 'droid'
  | 'firebender'
  | 'forgecode'
  | 'gemini-cli'
  | 'github-copilot'
  | 'goose'
  | 'hermes-agent'
  | 'iflow-cli'
  | 'junie'
  | 'kilo'
  | 'kimi-cli'
  | 'kiro-cli'
  | 'kode'
  | 'mcpjam'
  | 'mistral-vibe'
  | 'mux'
  | 'neovate'
  | 'opencode'
  | 'openhands'
  | 'pi'
  | 'qoder'
  | 'qwen-code'
  | 'replit'
  | 'roo'
  | 'rovodev'
  | 'tabnine-cli'
  | 'trae'
  | 'trae-cn'
  | 'warp'
  | 'windsurf'
  | 'zencoder'
  | 'pochi'
  | 'adal'
  | 'universal';

export interface Skill {
  name: string;
  description: string;
  path: string;
  /** 用于哈希的原始 SKILL.md 内容 */
  rawContent?: string;
  /** 所属插件名称（若有） */
  pluginName?: string;
  metadata?: Record<string, unknown>;
}

export interface AgentConfig {
  name: string;
  displayName: string;
  skillsDir: string;
  /** 全局技能目录；若 agent 不支持全局安装则为 undefined。 */
  globalSkillsDir: string | undefined;
  detectInstalled: () => Promise<boolean>;
  /** 是否在通用 agent 列表中显示；默认为 true。 */
  showInUniversalList?: boolean;
}

export interface ParsedSource {
  type: 'github' | 'gitlab' | 'git' | 'local' | 'well-known';
  url: string;
  subpath?: string;
  localPath?: string;
  ref?: string;
  /** 从 @skill 语法提取的技能名（如 owner/repo@skill-name） */
  skillFilter?: string;
}

/**
 * 从远程宿主 provider 获取的技能。
 */
export interface RemoteSkill {
  /** 显示名称（来自 frontmatter） */
  name: string;
  /** 描述（来自 frontmatter） */
  description: string;
  /** 含 frontmatter 的完整 markdown 内容 */
  content: string;
  /** 安装目录名标识 */
  installName: string;
  /** 原始 source URL */
  sourceUrl: string;
  /** 拉取该技能的 provider */
  providerId: string;
  /** 遥测用 source 标识（如 "mintlify.com"） */
  sourceIdentifier: string;
  /** frontmatter 中的其他元数据 */
  metadata?: Record<string, unknown>;
}
