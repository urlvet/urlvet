import type { PageLoad } from './$types';
import { formatUrl, getDomainFromUrl, formatUrlForShare, decodeVerdict } from '#lib/utils';

// Runs on the server for every request, giving bots the correct OG meta tags.
export const load: PageLoad = ({ url }) => {
  const raw = url.searchParams.get('q') ?? '';
  const q = raw === '{search_term_string}' ? '' : raw;
  const normalized = q ? formatUrl(q) : '';
  const domain = normalized ? getDomainFromUrl(normalized) : '';
  const formatted = domain ? formatUrlForShare(normalized) : '';
  const verdict = decodeVerdict(url.searchParams.get('v') ?? '');
  const score = url.searchParams.get('s') ?? '';

  return {
    queryDomain: domain,
    queryUrl: normalized,
    formattedQueryUrl: formatted,
    verdict,
    score,
  };
};
