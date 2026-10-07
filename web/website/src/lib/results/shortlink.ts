import type { RedirectInfo, ShortLinkInfo } from '../types';

function host(url: string): string {
  try {
    return new URL(url).hostname;
  } catch {
    return url;
  }
}

/** One sentence saying where a short link led, or null for an ordinary link. */
export function shortLinkNote(link: ShortLinkInfo | undefined): string | null {
  if (!link) return null;
  const via = host(link.url);
  if (!link.resolved || !link.target) {
    return `This is a short link on ${via}. We couldn't see where it leads, so this result is for the short link itself and can't be better than Suspicious.`;
  }
  // Short links can hand off to other short links; name each one on the way.
  const hops = [...new Set(link.chain.slice(1).map(host))].filter(
    (h) => h !== via && h !== host(link.target!)
  );
  const through = hops.length ? ` (through ${hops.join(', ')})` : '';
  return `This is a short link on ${via}. It leads to ${host(link.target)}${through}, and this result is for that page.`;
}

/** One sentence for a link a well-known site redirected elsewhere (an open redirect), or null. */
export function redirectNote(info: RedirectInfo | undefined): string | null {
  if (!info?.target) return null;
  return `This link on ${host(info.url)} sends visitors on to ${host(info.target)}. This result is for that page, plus anything wrong with the link itself.`;
}
