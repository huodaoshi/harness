const TELEMETRY_URL = 'https://add-skill.vercel.sh/t';
const AUDIT_URL = 'https://add-skill.vercel.sh/audit';

interface InstallTelemetryData {
  event: 'install';
  source: string;
  skills: string;
  agents: string;
  global?: '1';
  skillFiles?: string; // JSON 序列化：{ skillName: relativePath }
  /**
   * 不同宿主的 source 类型：
   * - 'github'：GitHub 仓库（默认，使用 raw.githubusercontent.com）
   * - 'raw'：SKILL.md 直链（通用 raw URL）
   * - 以及 'mintlify'、'huggingface' 等 provider ID
   */
  sourceType?: string;
}

interface RemoveTelemetryData {
  event: 'remove';
  source?: string;
  skills: string;
  agents: string;
  global?: '1';
  sourceType?: string;
}

interface UpdateTelemetryData {
  event: 'update';
  scope?: string;
  skillCount: string;
  successCount: string;
  failCount: string;
}

interface FindTelemetryData {
  event: 'find';
  query: string;
  resultCount: string;
  interactive?: '1';
}

interface SyncTelemetryData {
  event: 'experimental_sync';
  skillCount: string;
  successCount: string;
  agents: string;
}

type TelemetryData =
  | InstallTelemetryData
  | RemoveTelemetryData
  | UpdateTelemetryData
  | FindTelemetryData
  | SyncTelemetryData;

let cliVersion: string | null = null;
let detectedAgentName: string | null = null;

/**
 * 设置检测到的 AI agent 名称用于遥测。
 * 在 agent 检测时调用一次，随后包含在所有遥测事件中。
 */
export function setDetectedAgent(agentName: string | null): void {
  detectedAgentName = agentName;
}

function isCI(): boolean {
  return !!(
    process.env.CI ||
    process.env.GITHUB_ACTIONS ||
    process.env.GITLAB_CI ||
    process.env.CIRCLECI ||
    process.env.TRAVIS ||
    process.env.BUILDKITE ||
    process.env.JENKINS_URL ||
    process.env.TEAMCITY_VERSION
  );
}

function isEnabled(): boolean {
  return !process.env.DISABLE_TELEMETRY && !process.env.DO_NOT_TRACK;
}

export function setVersion(version: string): void {
  cliVersion = version;
}

// ─── 安全审计数据 ───

export interface PartnerAudit {
  risk: 'safe' | 'low' | 'medium' | 'high' | 'critical' | 'unknown';
  alerts?: number;
  score?: number;
  analyzedAt: string;
}

export type SkillAuditData = Record<string, PartnerAudit>;
export type AuditResponse = Record<string, SkillAuditData>;

/**
 * 从审计 API 获取技能安全审计结果。
 * 出错或超时时返回 null — 永不阻塞安装流程。
 */
export async function fetchAuditData(
  source: string,
  skillSlugs: string[],
  timeoutMs = 3000
): Promise<AuditResponse | null> {
  if (skillSlugs.length === 0) return null;

  try {
    const params = new URLSearchParams({
      source,
      skills: skillSlugs.join(','),
    });

    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), timeoutMs);

    const response = await fetch(`${AUDIT_URL}?${params.toString()}`, {
      signal: controller.signal,
    });
    clearTimeout(timeout);

    if (!response.ok) return null;
    return (await response.json()) as AuditResponse;
  } catch {
    return null;
  }
}

// 进行中的遥测请求 — CLI 退出前 await，避免丢数据，但不阻塞主流程。
const pendingTelemetry: Promise<void>[] = [];

export function track(data: TelemetryData): void {
  if (!isEnabled()) return;

  try {
    const params = new URLSearchParams();

    // 附加版本号
    if (cliVersion) {
      params.set('v', cliVersion);
    }

    // 在 CI 环境中附加 CI 标记
    if (isCI()) {
      params.set('ci', '1');
    }

    // 附加检测到的 AI agent 名称
    if (detectedAgentName) {
      params.set('agent', detectedAgentName);
    }

    // 附加事件数据
    for (const [key, value] of Object.entries(data)) {
      if (value !== undefined && value !== null) {
        params.set(key, String(value));
      }
    }

    // 工作流中 fire-and-forget，但记录 promise，
    // 以便 flushTelemetry() 在进程退出前 await。
    const p = fetch(`${TELEMETRY_URL}?${params.toString()}`)
      .catch(() => {})
      .then(() => {});
    pendingTelemetry.push(p);
  } catch {
    // 静默失败 — 遥测不得影响 CLI
  }
}

/**
 * 等待所有进行中的遥测请求结束。
 * CLI 退出时调用一次，避免在打开套接字上挂起，也避免过早退出丢数据。
 */
export async function flushTelemetry(timeoutMs = 5000): Promise<void> {
  if (pendingTelemetry.length === 0) return;
  const timeout = new Promise<void>((resolve) => setTimeout(resolve, timeoutMs));
  await Promise.race([Promise.all(pendingTelemetry), timeout]);
}
