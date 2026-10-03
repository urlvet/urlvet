import type { AnalyzeResult } from '../types';

// Letters from other alphabets that look like Latin ones, as used in lookalike
// addresses (аррӏе.com is Cyrillic а, р, р, ӏ, е).
const LOOKALIKES: Record<string, string> = {
  а: 'a',
  е: 'e',
  о: 'o',
  р: 'p',
  с: 'c',
  у: 'y',
  х: 'x',
  і: 'i',
  ј: 'j',
  ӏ: 'l',
  ԁ: 'd',
  ѕ: 's',
  ԛ: 'q',
  ԝ: 'w',
  һ: 'h',
  ɡ: 'g',
  α: 'a',
  ο: 'o',
  ν: 'v',
  ρ: 'p',
};

/** The plain-Latin address a lookalike imitates, or null if it isn't one. */
function unmask(domain: string): string | null {
  const plain = [...domain].map((ch) => LOOKALIKES[ch] ?? ch).join('');
  return plain !== domain && /^[a-z0-9.-]+$/.test(plain) ? plain : null;
}

export type RealSite = { domain: string; reason: 'brand' | 'lookalike' };

/** The site a flagged page seems to imitate, so Vetty can point to the real one. */
export function realSite(r: AnalyzeResult): RealSite | null {
  if (!r.result?.verdict || r.result.verdict === 'Safe') return null;
  const own = r.domain?.toLowerCase();
  const brand = r.content_data?.brand_check;
  if (brand?.is_mismatch && brand.official_domain && brand.official_domain !== own) {
    return { domain: brand.official_domain, reason: 'brand' };
  }
  const typo = r.typosquat_result;
  if (typo?.is_suspicious && typo.matched_domain && typo.matched_domain !== own) {
    return { domain: typo.matched_domain, reason: 'lookalike' };
  }
  // Lookalikes that swap many letters (аррӏе.com) are too far from the real
  // name for the spelling check, but undoing the swap gives it back.
  const plain = own ? unmask(own) : null;
  if (plain) return { domain: plain, reason: 'lookalike' };
  return null;
}
