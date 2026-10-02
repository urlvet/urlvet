// Plain words for checks that didn't finish, keyed by backend task name.
const CHECK_NAMES: Record<string, string> = {
  phishtank_check: 'the phishing database (PhishTank)',
  content_check: 'the page itself',
  http_combined_check: 'redirects and HTTPS',
  tls_combined_check: 'the security certificate',
  whois_lookup: "the domain's registration",
  dns_validity_check: 'DNS records',
  ip_resolution: 'DNS records',
  domain_rank: 'how popular the site is',
};

/** One sentence saying what wasn't checked, or null if everything ran. */
export function incompleteNote(
  checks: string[] | undefined,
  incomplete: boolean | undefined
): string | null {
  if (!checks?.length)
    return incomplete
      ? "Some checks couldn't finish, so this result may be missing signals. Try again in a minute."
      : null;

  const names = [...new Set(checks.map((c) => CHECK_NAMES[c] ?? 'a few other checks'))];
  const list =
    names.length === 1
      ? names[0]
      : `${names.slice(0, -1).join(', ')} and ${names[names.length - 1]}`;
  return incomplete
    ? `Couldn't check ${list} this time, so this result may be missing signals. Try again in a minute.`
    : `Couldn't check ${list} this time. Everything else ran normally.`;
}
