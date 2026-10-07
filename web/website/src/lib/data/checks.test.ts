import { describe, expect, it } from 'vitest';
import { CHECK_GROUPS, TOTAL_CHECKS } from './checks';

describe('check groups', () => {
  it('add up to the 21 checks the scanner runs', () => {
    expect(TOTAL_CHECKS).toBe(21);
  });

  it('keep the ids that result sections link to', () => {
    expect(CHECK_GROUPS.map((g) => g.id)).toEqual([
      'url',
      'network',
      'dns',
      'tls',
      'domain',
      'content',
      'threats',
    ]);
  });

  it('mark every finding as risk, trust or shown only', () => {
    for (const s of CHECK_GROUPS.flatMap((g) => g.signals)) {
      expect(['risk', 'trust', 'info']).toContain(s.kind);
    }
  });
});
