import { describe, it, expect } from 'vitest';
import {
  parseRulesAddOptions,
  getRulesPackSubdir,
  RULES_AGENT_TYPES,
} from '../src/rules-add.ts';

describe('parseRulesAddOptions', () => {
  it('parses source, -a, --to', () => {
    const { source, options } = parseRulesAddOptions([
      'owner/repo',
      '-a',
      'cursor',
      '--to',
      './out/rules',
    ]);
    expect(source).toEqual(['owner/repo']);
    expect(options.agent).toEqual(['cursor']);
    expect(options.to).toBe('./out/rules');
  });

  it('requires positional source separate from flags', () => {
    const { source, options } = parseRulesAddOptions(['./local-pack', '-a', 'claude-code']);
    expect(source).toEqual(['./local-pack']);
    expect(options.agent).toEqual(['claude-code']);
  });
});

describe('getRulesPackSubdir', () => {
  it('maps agents to package subfolders', () => {
    expect(getRulesPackSubdir('cursor')).toBe('rules/cursor');
    expect(getRulesPackSubdir('claude-code')).toBe('rules/claude');
  });
});

describe('RULES_AGENT_TYPES', () => {
  it('lists MVP agents', () => {
    expect(RULES_AGENT_TYPES).toContain('cursor');
    expect(RULES_AGENT_TYPES).toContain('claude-code');
  });
});
