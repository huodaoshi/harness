/**
 * 从远程宿主解析得到的技能。
 * 不同宿主识别技能的方式可能不同。
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
  /** frontmatter 中的其他元数据 */
  metadata?: Record<string, unknown>;
}

/**
 * URL 与 provider 匹配尝试的结果。
 */
export interface ProviderMatch {
  /** URL 是否匹配该 provider */
  matches: boolean;
  /** 遥测/存储用 source 标识（如 "mintlify/bun.com"、"huggingface/hf-skills/hf-jobs"） */
  sourceIdentifier?: string;
}

/**
 * 远程 SKILL.md 宿主 provider 接口。
 * 每个 provider 需能：
 * - 判断 URL 是否属于本宿主
 * - 拉取并解析 SKILL.md
 * - 将 URL 转为 raw 内容 URL
 * - 提供遥测用 source 标识
 */
export interface HostProvider {
  /** 唯一 ID（如 "mintlify"、"huggingface"、"github"） */
  readonly id: string;

  /** 显示名称 */
  readonly displayName: string;

  /**
   * 判断 URL 是否匹配本 provider。
   * @param url - 待检查的 URL
   * @returns 匹配结果，可含 source 标识
   */
  match(url: string): ProviderMatch;

  /**
   * 从给定 URL 拉取并解析 SKILL.md。
   * @param url - SKILL.md 的 URL
   * @returns 解析后的技能；无效或未找到时返回 null
   */
  fetchSkill(url: string): Promise<RemoteSkill | null>;

  /**
   * 将面向用户的 URL 转为 raw 内容 URL。
   * 例如 GitHub blob URL → raw.githubusercontent.com。
   * @param url - 待转换 URL
   * @returns raw 内容 URL
   */
  toRawUrl(url: string): string;

  /**
   * 获取遥测/存储用 source 标识。
   * 应为稳定标识，便于按来源分组技能。
   * @param url - 原始 URL
   * @returns source 标识（如 "mintlify/bun.com"）
   */
  getSourceIdentifier(url: string): string;
}

/** 宿主 provider 注册表。 */
export interface ProviderRegistry {
  /** 注册新 provider。 */
  register(provider: HostProvider): void;

  /**
   * 查找匹配给定 URL 的 provider。
   * @param url - 待匹配 URL
   * @returns 匹配的 provider，无则 null
   */
  findProvider(url: string): HostProvider | null;

  /** 获取所有已注册 provider。 */
  getProviders(): HostProvider[];
}
