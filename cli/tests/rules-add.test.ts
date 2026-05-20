import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { resolve } from 'path';
import {
  parseRulesAddOptions,
  getRulesPackSubdir,
  RULES_AGENT_TYPES,
  resolveRulesProjectRoot,
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

  it('parses --cwd and -C', () => {
    expect(parseRulesAddOptions(['.', '-a', 'cursor', '--cwd', 'D:\\proj']).options.cwd).toBe('D:\\proj');
    expect(parseRulesAddOptions(['.', '-C', '/tmp/x', '-a', 'cursor']).options.cwd).toBe('/tmp/x');
  });
});

describe('resolveRulesProjectRoot', () => {
  const hadInit = 'INIT_CWD' in process.env;
  const prevInit = process.env.INIT_CWD;

  beforeEach(() => {
    delete process.env.INIT_CWD;
  });

  afterEach(() => {
    if (hadInit) process.env.INIT_CWD = prevInit;
    else delete process.env.INIT_CWD;
  });

  it('uses explicit path when given', () => {
    process.env.INIT_CWD = '/from-env';
    expect(resolveRulesProjectRoot('/explicit')).toBe(resolve('/explicit'));
  });

  it('uses INIT_CWD when no explicit path', () => {
    process.env.INIT_CWD = 'D:\\other-project';
    expect(resolveRulesProjectRoot()).toBe(resolve('D:\\other-project'));
  });

  it('falls back to process.cwd when INIT_CWD unset', () => {
    expect(resolveRulesProjectRoot()).toBe(process.cwd());
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
