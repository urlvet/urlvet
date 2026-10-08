import type { Tone } from '../ui/tone';
import type { AnalyzeResult } from '../types';

export type SectionId =
  | 'domain'
  | 'analysis'
  | 'threatintel'
  | 'security'
  | 'content'
  | 'features'
  | 'infrastructure';

export type Status = { tone: Tone; label: string };

// Sections open by default only when they hold a real problem.
export function computeExpanded(d: AnalyzeResult | null): Record<string, boolean> {
  if (!d)
    return {
      domain: false,
      analysis: false,
      threatintel: false,
      security: false,
      content: false,
      features: false,
      infrastructure: false,
      performance: false,
    };
  return {
    // age_days is null when the registry publishes no creation date, and
    // null < 365 is true in JS, so check for it explicitly.
    domain: d.domain_info?.age_days != null && d.domain_info.age_days < 365,
    analysis:
      !!d.analysis?.redirection_result?.has_domain_jump ||
      (d.analysis?.redirection_result?.chain_length ?? 0) > 3,
    threatintel:
      !!d.phishing?.valid ||
      !!d.threat_feeds?.listed ||
      !!d.web_risk?.listed ||
      !!d.safe_browsing?.listed,
    security: !!(
      d.ssl_info?.IsSuspicious ||
      d.ssl_info?.KnownBadChain ||
      (d.ssl_info && !d.ssl_info.HasTLS) ||
      (d.ssl_info?.HasTLS && !d.ssl_info.ChainValid) ||
      d.tls_info?.HostnameMismatch
    ),
    content: !!(
      d.content_data?.brand_check?.is_mismatch ||
      d.content_data?.has_hidden_iframe ||
      d.content_data?.forms?.some((f) => f.is_external)
    ),
    features: !!(
      d.features?.url?.has_homoglyph ||
      d.features?.url?.contains_punycode ||
      d.features?.url?.uses_ip ||
      d.features?.url?.url_shortener ||
      d.typosquat_result?.is_suspicious ||
      (d.domain_randomness?.entropy ?? 0) > 3.8 ||
      d.features?.tld?.is_risky_tld ||
      (d.features?.tld && !d.features?.tld?.is_icann && !d.features?.tld?.is_hosting_platform)
    ),
    infrastructure: !!(
      d.infrastructure &&
      !d.infrastructure.nameservers_valid &&
      !d.features?.tld?.is_hosting_platform
    ),
    performance: false,
  };
}

// One-glance summary shown on each collapsed section header.
export function sectionStatuses(d: AnalyzeResult): Record<SectionId, Status> {
  const days: number | undefined = d.domain_info?.age_days ?? undefined;
  const redir = d.analysis?.redirection_result;
  const ssl = d.ssl_info;
  const c = d.content_data;
  const url = d.features?.url;
  const tld = d.features?.tld;
  const ips: string[] = d.infrastructure?.ip_addresses ?? [];

  return {
    domain:
      days === undefined
        ? { tone: 'neutral', label: 'Age unknown' }
        : days <= 30
          ? { tone: 'bad', label: `New · ${d.domain_info.age_human}` }
          : days < 365
            ? { tone: 'warn', label: `Young · ${d.domain_info.age_human}` }
            : { tone: 'good', label: d.domain_info.age_human },
    analysis: redir?.has_domain_jump
      ? { tone: 'bad', label: 'Leaves domain' }
      : (redir?.chain_length ?? 0) > 3
        ? { tone: 'warn', label: `${redir.chain_length} hops` }
        : redir?.is_redirected
          ? { tone: 'neutral', label: 'Same-site redirect' }
          : { tone: 'good', label: 'No redirects' },
    threatintel:
      d.web_risk?.listed || d.safe_browsing?.listed
        ? { tone: 'bad', label: 'Listed by Google' }
        : d.threat_feeds?.listed
          ? d.threat_feeds.match === 'url'
            ? { tone: 'bad', label: 'Listed as a threat' }
            : { tone: 'bad', label: 'Site has listed pages' }
          : d.phishing?.valid
            ? d.phishing.verified
              ? { tone: 'bad', label: 'Confirmed phishing' }
              : { tone: 'bad', label: 'Reported phishing' }
            : { tone: 'good', label: 'Not listed' },
    security:
      ssl && !ssl.HasTLS
        ? { tone: 'bad', label: 'No HTTPS' }
        : ssl?.KnownBadChain || ssl?.IsSuspicious
          ? { tone: 'bad', label: 'Suspicious cert' }
          : (ssl?.HasTLS && !ssl.ChainValid) || d.tls_info?.HostnameMismatch
            ? { tone: 'bad', label: 'Invalid cert' }
            : { tone: 'good', label: 'Valid certificate' },
    content: !c
      ? { tone: 'neutral', label: 'Not fetched' }
      : c.brand_check?.is_mismatch
        ? { tone: 'bad', label: 'Brand mismatch' }
        : c.forms?.some((f) => f.is_external)
          ? { tone: 'bad', label: 'External form' }
          : c.has_hidden_iframe
            ? { tone: 'bad', label: 'Hidden iframe' }
            : c.has_payment_form
              ? { tone: 'warn', label: 'Payment form' }
              : c.has_login_form
                ? { tone: 'neutral', label: 'Login form' }
                : { tone: 'good', label: 'Nothing unusual' },
    features: url?.has_homoglyph
      ? { tone: 'bad', label: 'Homoglyph' }
      : url?.contains_punycode
        ? { tone: 'bad', label: 'Punycode' }
        : url?.uses_ip
          ? { tone: 'bad', label: 'Raw IP' }
          : d.typosquat_result?.is_suspicious
            ? { tone: 'bad', label: 'Lookalike domain' }
            : url?.url_shortener
              ? { tone: 'bad', label: 'URL shortener' }
              : (d.domain_randomness?.entropy ?? 0) > 3.8
                ? { tone: 'bad', label: 'Random-looking' }
                : tld?.is_risky_tld
                  ? { tone: 'warn', label: 'Risky TLD' }
                  : tld && !tld.is_icann && !tld.is_hosting_platform
                    ? { tone: 'warn', label: 'Unregulated TLD' }
                    : { tone: 'good', label: 'Clean' },
    infrastructure:
      d.infrastructure && !d.infrastructure.nameservers_valid && !tld?.is_hosting_platform
        ? { tone: 'bad', label: 'DNS issue' }
        : ips.length
          ? { tone: 'good', label: `${ips.length} IP${ips.length === 1 ? '' : 's'}` }
          : { tone: 'neutral', label: 'DNS OK' },
  };
}

// Outline icons (24x24) for each section header.
export const SECTION_ICONS: Record<SectionId, string> = {
  domain:
    'M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3zM6 6h.008v.008H6V6z',
  analysis: 'M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5',
  threatintel:
    'M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z',
  security:
    'M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z',
  content:
    'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z',
  features:
    'M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244',
  infrastructure:
    'M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z',
};
