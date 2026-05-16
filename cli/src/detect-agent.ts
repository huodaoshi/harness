import { determineAgent, type AgentResult } from '@vercel/detect-agent';
import { setDetectedAgent } from './telemetry.ts';
import type { AgentType } from './types.ts';

let cachedResult: AgentResult | null = null;

/**
 * 将 @vercel/detect-agent 的名称映射为 skills-cli 的 AgentType 标识。
 * 仅包含两套系统中都存在的 agent。
 */
const agentNameToType: Record<string, AgentType> = {
  cursor: 'cursor',
  'cursor-cli': 'cursor',
  claude: 'claude-code',
  cowork: 'claude-code',
  devin: 'universal', // Devin 不在 skills-cli agent 列表中，回退到 universal
  replit: 'replit',
  gemini: 'gemini-cli',
  codex: 'codex',
  antigravity: 'antigravity',
  'augment-cli': 'augment',
  opencode: 'opencode',
  'github-copilot': 'github-copilot',
};

/**
 * 检测 CLI 是否在 AI agent 环境中运行。
 * 首次调用后缓存结果，并用 agent 名称更新遥测。
 */
export async function detectAgent(): Promise<AgentResult> {
  if (cachedResult) return cachedResult;
  cachedResult = await determineAgent();
  if (cachedResult.isAgent) {
    setDetectedAgent(cachedResult.agent.name);
  }
  return cachedResult;
}

/**
 * 若 CLI 在已检测到的 AI agent 内运行则返回 true。
 * 此时应跳过交互式提示并使用合理默认值。
 */
export async function isRunningInAgent(): Promise<boolean> {
  const result = await detectAgent();
  return result.isAgent;
}

/** 返回检测到的 agent 名称；若不在 agent 环境中则返回 null。 */
export async function getAgentName(): Promise<string | null> {
  const result = await detectAgent();
  return result.isAgent ? result.agent.name : null;
}

/**
 * 将检测到的 agent 名称映射为对应的 skills-cli AgentType。
 * 无法映射到具体 skills-cli agent 时返回 null。
 */
export function getAgentType(agentName: string): AgentType | null {
  return agentNameToType[agentName] ?? null;
}
