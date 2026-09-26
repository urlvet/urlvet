const TRACKER_PARAMS = new Set([
  // UTM (universal)
  'utm_source',
  'utm_medium',
  'utm_campaign',
  'utm_term',
  'utm_content',
  'utm_id',
  // Google Ads / Google Analytics
  'gclid',
  'gbraid',
  'wbraid',
  'gclsrc',
  'gad_source',
  'dclid',
  // Facebook / Meta
  'fbclid',
  'fb_action_ids',
  'fb_action_types',
  'fb_source',
  // Microsoft / Bing Ads
  'msclkid',
  // Microsoft Shopping
  '_hvadid',
  '_hvdev',
  '_hvkeyword',
  '_hvnetw',
  '_hvqmt',
  '_hvbmt',
  '_hvpos',
  '_hvtargid',
  '_hvcamid',
  '_hvadgrpid',
  '_hvpid',
  // TikTok
  'ttclid',
  // Twitter / X
  'twclid',
  // LinkedIn
  'li_fat_id',
  // HubSpot (email & ads)
  '_hsenc',
  '_hsmi',
  'hsa_acc',
  'hsa_cam',
  'hsa_grp',
  'hsa_ad',
  'hsa_src',
  'hsa_tgt',
  'hsa_kw',
  'hsa_mt',
  'hsa_net',
  'hsa_ver',
  // Mailchimp
  'mc_eid',
  'mc_cid',
  // Marketo
  'mkt_tok',
  // Pardot (Salesforce)
  'pi_campaign_id',
  'pi_list_email_id',
  'pi_list_id',
  'pi_uid',
  'pi_utid',
  // Oracle Eloqua
  'elqtrackid',
  'elqaid',
  'elqat',
  'elqcampaignid',
  // Klaviyo
  '_kx',
  'klaviyo_id',
  // Iterable
  'iterableemailcampaignid',
  'iterabletemplateid',
  'iterablemessageid',
  // Sailthru
  'sailthru_mid',
  // Customer.io
  'cio_id',
  // ConvertKit
  'ck_subscriber_id',
  // Drip
  '__s',
  // ActiveCampaign
  'vgo_ee',
  // Brevo / Sendinblue
  'sib_id',
  // GetResponse
  'gr_pk',
  // AWeber
  'awc',
  // Omnisend
  'omnisendcontactid',
  // Instagram
  'igshid',
  'igsh',
  // Pinterest
  'epik',
  'e_t',
  // Snapchat
  'sc_channel',
  'sccid',
  // Reddit
  'rdt_cid',
  // Yandex
  'yclid',
  'ymclid',
  // Baidu
  'bd_vid',
  'bd_source',
  // Adobe Analytics
  's_kwcid',
  'ef_id',
  's_cid',
  // IBM / Acoustic
  'cm_mmc',
  'cm_mmca1',
  'cm_mmca2',
  // Criteo
  'cto_pld',
  'cto_tld',
  'cto_gum',
  // AppsFlyer (mobile deep links)
  'af_sub1',
  'af_sub2',
  'af_sub3',
  'af_sub4',
  'af_sub5',
  'af_siteid',
  // Adjust (mobile)
  'adj_t',
  // Branch.io (mobile deep links)
  '_branch_match_id',
  '_branch_referrer',
  // Kochava
  'ko_click_id',
  // Singular
  'singular_click_id',
  // Impact / Radius (affiliate)
  'irclickid',
  // Commission Junction (affiliate)
  'cjuid',
  'cjevent',
  // ShareASale (affiliate)
  'sscid',
  // Rakuten (affiliate)
  'ransiteid',
  'ranmid',
  'raneaid',
  'ranlinkid',
  // Tradedoubler (affiliate)
  'tduid',
  // Awin (affiliate)
  'awin_mid',
  // Partnerize (affiliate)
  'pf_id',
]);

export function stripTrackers(raw: string): { cleaned: string; removed: string[] } {
  try {
    const url = new URL(raw.startsWith('http') ? raw : `https://${raw}`);
    const removed: string[] = [];
    for (const key of [...url.searchParams.keys()]) {
      if (TRACKER_PARAMS.has(key.toLowerCase())) {
        removed.push(key);
        url.searchParams.delete(key);
      }
    }
    return { cleaned: url.toString(), removed };
  } catch {
    return { cleaned: raw, removed: [] };
  }
}

const VERDICT_ENCODE: Record<string, string> = { Safe: 'f3', Risky: '9c', Suspicious: '7b' };
const VERDICT_DECODE: Record<string, string> = { f3: 'Safe', '9c': 'Risky', '7b': 'Suspicious' };

export function encodeVerdict(v: string): string {
  return VERDICT_ENCODE[v] ?? v;
}

export function decodeVerdict(v: string): string {
  return VERDICT_DECODE[v] ?? v;
}

export function isValidUrl(url: string): boolean {
  try {
    const parsed = new URL(url);
    const hostname = parsed.hostname;
    // IPv4
    if (/^(\d{1,3}\.){3}\d{1,3}$/.test(hostname)) return true;
    // IPv6 (brackets stripped by URL parser, e.g. [::1] -> ::1)
    if (hostname.includes(':')) return true;
    // Must have a dot and a TLD of at least 2 chars
    const lastDot = hostname.lastIndexOf('.');
    return lastDot !== -1 && hostname.length - lastDot - 1 >= 2;
  } catch {
    return false;
  }
}

export function formatUrl(url: string): string {
  url = url.trim();
  // Add protocol if missing
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    return `http://${url}`;
  }
  return url;
}

export function getDomainFromUrl(url: string): string {
  try {
    const urlObj = new URL(url);
    return urlObj.hostname;
  } catch {
    return '';
  }
}

export function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  } catch {
    return dateString;
  }
}

/**
 * Formats a URL for safe sharing by stripping schema and replacing dots with [.]
 * Example: "https://example.com" -> "example[.]com"
 */
export function formatUrlForShare(url: string): string {
  try {
    // Remove schema (http://, https://)
    let formatted = url.replace(/^https?:\/\//i, '');
    // Remove trailing slash
    formatted = formatted.replace(/\/$/, '');
    // Replace dots with [.]
    formatted = formatted.replace(/\./g, '[.]');
    return formatted;
  } catch {
    // If parsing fails, just replace dots
    return url.replace(/\./g, '[.]');
  }
}

/** Characters that end a sentence rather than a link, e.g. "visit example.com." */
const LINK_TRAILING = /[.,;:!?'"’”)\]}>»]+$/u;
const LINK_PATTERN =
  /(?:https?:\/\/)?(?:[\p{L}\p{N}](?:[\p{L}\p{N}-]*[\p{L}\p{N}])?\.)+\p{L}{2,}(?::\d+)?(?:[/?#][^\s<>"]*)?/giu;

/**
 * Finds the links in a block of text, such as a forwarded WhatsApp message,
 * in order and without duplicates. Email addresses are skipped. If any link
 * has a scheme (or starts with www.), bare "word.word" matches are dropped,
 * since in a message those are more often typos than links.
 */
export function extractLinks(text: string): string[] {
  const found: { link: string; explicit: boolean }[] = [];
  const seen = new Set<string>();
  for (const m of text.matchAll(LINK_PATTERN)) {
    const before = m.index ? text[m.index - 1] : '';
    if (before && /[@\p{L}\p{N}_]/u.test(before)) continue; // email or part of a word
    const link = m[0].replace(LINK_TRAILING, '');
    const key = link.toLowerCase().replace(/^https?:\/\//, '');
    if (seen.has(key)) continue;
    seen.add(key);
    found.push({ link, explicit: /^(https?:\/\/|www\.)/i.test(link) });
  }
  const explicit = found.filter((f) => f.explicit);
  return (explicit.length ? explicit : found).map((f) => f.link);
}
