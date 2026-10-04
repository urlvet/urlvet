import { describe, expect, it } from 'vitest';
import { FAQ_GROUPS, faqId, faqText } from './faq';

const all = FAQ_GROUPS.flatMap((g) => g.items);

describe('FAQ', () => {
  it('gives every question its own link', () => {
    const ids = all.map((f) => faqId(f.q));
    expect(new Set(ids).size).toBe(ids.length);
    expect(faqId('What is phishing?')).toBe('faq-what-is-phishing');
  });

  it('includes list points in the answer text', () => {
    const tips = all.find((f) => f.points);
    expect(tips && faqText(tips)).toContain('two-step verification');
  });
});
