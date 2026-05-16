import { homedir } from 'os';
import type { AgentType } from './types.ts';
import { agents } from './agents.ts';
import { listInstalledSkills, type InstalledSkill } from './installer.ts';
import { sanitizeMetadata } from './sanitize.ts';
import { getAllLockedSkills } from './skill-lock.ts';

const RESET = '\x1b[0m';
const BOLD = '\x1b[1m';
const DIM = '\x1b[38;5;102m';
const TEXT = '\x1b[38;5;145m';
const CYAN = '\x1b[36m';
const YELLOW = '\x1b[33m';

interface ListOptions {
  global?: boolean;
  agent?: string[];
  json?: boolean;
}

/** 缩短路径用于显示：将 homedir 换为 ~，cwd 换为 . */
function shortenPath(fullPath: string, cwd: string): string {
  const home = homedir();
  if (fullPath.startsWith(home)) {
    return fullPath.replace(home, '~');
  }
  if (fullPath.startsWith(cwd)) {
    return '.' + fullPath.slice(cwd.length);
  }
  return fullPath;
}

/** 格式化列表；过多时截断显示 */
function formatList(items: string[], maxShow: number = 5): string {
  if (items.length <= maxShow) {
    return items.join(', ');
  }
  const shown = items.slice(0, maxShow);
  const remaining = items.length - maxShow;
  return `${shown.join(', ')} +${remaining} more`;
}

export function parseListOptions(args: string[]): ListOptions {
  const options: ListOptions = {};

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '-g' || arg === '--global') {
      options.global = true;
    } else if (arg === '--json') {
      options.json = true;
    } else if (arg === '-a' || arg === '--agent') {
      options.agent = options.agent || [];
      // 收集后续参数直至下一个 flag
      while (i + 1 < args.length && !args[i + 1]!.startsWith('-')) {
        options.agent.push(args[++i]!);
      }
    }
  }

  return options;
}

export async function runList(args: string[]): Promise<void> {
  const options = parseListOptions(args);

  // 默认仅项目（本地）；-g 表示全局
  const scope = options.global === true ? true : false;

  // 校验 agent 过滤参数
  let agentFilter: AgentType[] | undefined;
  if (options.agent && options.agent.length > 0) {
    const validAgents = Object.keys(agents);
    const invalidAgents = options.agent.filter((a) => !validAgents.includes(a));

    if (invalidAgents.length > 0) {
      console.log(`${YELLOW}Invalid agents: ${invalidAgents.join(', ')}${RESET}`);
      console.log(`${DIM}Valid agents: ${validAgents.join(', ')}${RESET}`);
      process.exit(1);
    }

    agentFilter = options.agent as AgentType[];
  }

  const installedSkills = await listInstalledSkills({
    global: scope,
    agentFilter,
  });

  // JSON 模式：结构化、无 ANSI、agent 列表不截断
  if (options.json) {
    const jsonOutput = installedSkills.map((skill) => ({
      name: skill.name,
      path: skill.canonicalPath,
      scope: skill.scope,
      agents: skill.agents.map((a) => agents[a].displayName),
    }));
    console.log(JSON.stringify(jsonOutput, null, 2));
    return;
  }

  // 读取 lock 条目以获取插件分组信息
  const lockedSkills = await getAllLockedSkills();

  const cwd = process.cwd();
  const scopeLabel = scope ? 'Global' : 'Project';

  if (installedSkills.length === 0) {
    if (options.json) {
      console.log('[]');
      return;
    }
    console.log(`${DIM}No ${scopeLabel.toLowerCase()} skills found.${RESET}`);
    if (scope) {
      console.log(`${DIM}Try listing project skills without -g${RESET}`);
    } else {
      console.log(`${DIM}Try listing global skills with -g${RESET}`);
    }
    return;
  }

  function printSkill(skill: InstalledSkill, indent: boolean = false): void {
    const prefix = indent ? '  ' : '';
    const shortPath = shortenPath(skill.canonicalPath, cwd);
    const agentNames = skill.agents.map((a) => agents[a].displayName);
    const agentInfo =
      skill.agents.length > 0 ? formatList(agentNames) : `${YELLOW}not linked${RESET}`;
    console.log(
      `${prefix}${CYAN}${sanitizeMetadata(skill.name)}${RESET} ${DIM}${shortPath}${RESET}`
    );
    console.log(`${prefix}  ${DIM}Agents:${RESET} ${agentInfo}`);
  }

  console.log(`${BOLD}${scopeLabel} Skills${RESET}`);
  console.log();

  // 按插件分组
  const groupedSkills: Record<string, InstalledSkill[]> = {};
  const ungroupedSkills: InstalledSkill[] = [];

  for (const skill of installedSkills) {
    const lockEntry = lockedSkills[skill.name];
    if (lockEntry?.pluginName) {
      const group = lockEntry.pluginName;
      if (!groupedSkills[group]) {
        groupedSkills[group] = [];
      }
      groupedSkills[group].push(skill);
    } else {
      ungroupedSkills.push(skill);
    }
  }

  const hasGroups = Object.keys(groupedSkills).length > 0;

  if (hasGroups) {
    // 按字母序打印分组
    const sortedGroups = Object.keys(groupedSkills).sort();
    for (const group of sortedGroups) {
      // 将 kebab-case 转为标题式显示名
      const title = group
        .split('-')
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(' ');

      console.log(`${BOLD}${title}${RESET}`);
      const skills = groupedSkills[group];
      if (skills) {
        for (const skill of skills) {
          printSkill(skill, true);
        }
      }
      console.log();
    }

    // 打印未分组的技能（若有）
    if (ungroupedSkills.length > 0) {
      console.log(`${BOLD}General${RESET}`);
      for (const skill of ungroupedSkills) {
        printSkill(skill, true);
      }
      console.log();
    }
  } else {
    // 无分组时按扁平列表打印
    for (const skill of installedSkills) {
      printSkill(skill);
    }
    console.log();
  }
}
