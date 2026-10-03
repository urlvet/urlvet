import { describe, expect, it } from 'vitest';
import type { AnalyzeResult } from '../types';
import { realSite } from './realsite';

const scan = (over: Partial<AnalyzeResult>, verdict = 'Risky') =>
  ({ url: '', domain: 'paypa1.com', result: { verdict }, ...over }) as unknown as AnalyzeResult;

describe('realSite', () => {
  it('points a brand impersonator to the brand', () => {
    const r = scan({
      content_data: {
        brand_check: {
          brand_found: 'PayPal',
          official_domain: 'paypal.com',
          is_mismatch: true,
          detected_names: [],
        },
      },
    } as unknown as Partial<AnalyzeResult>);
    expect(realSite(r)).toEqual({ domain: 'paypal.com', reason: 'brand' });
  });

  it('points a lookalike address to the site it copies', () => {
    const r = scan({ typosquat_result: { is_suspicious: true, matched_domain: 'paypal.com' } });
    expect(realSite(r)).toEqual({ domain: 'paypal.com', reason: 'lookalike' });
  });

  it('undoes letters swapped in from other alphabets', () => {
    expect(realSite(scan({ domain: 'аррӏе.com' }))).toEqual({
      domain: 'apple.com',
      reason: 'lookalike',
    });
    expect(realSite(scan({ domain: 'example.com' }))).toBeNull();
  });

  it('says nothing for safe results or without a match', () => {
    expect(
      realSite(
        scan({ typosquat_result: { is_suspicious: true, matched_domain: 'paypal.com' } }, 'Safe')
      )
    ).toBeNull();
    expect(realSite(scan({}))).toBeNull();
  });
});
