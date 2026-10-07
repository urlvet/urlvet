// Shared by the popup, the blocked page and the background worker.
// Wording and link parsing come straight from the website's source, so the
// extension and url.vet always say the same thing.
export { VERDICT_COPY, verdictCopy } from "../../website/src/lib/verdict";
export { extractLinks, formatUrl } from "../../website/src/lib/utils";
export {
  redirectNote,
  shortLinkNote,
} from "../../website/src/lib/results/shortlink";
export { incompleteNote } from "../../website/src/lib/results/incomplete";

// The public instance. Self-hosted copies set their own in Settings.
export const DEFAULTS = {
  apiBase: "https://api.url.vet/api/v1",
  siteBase: "https://url.vet",
  // Warn on risky pages from the start; Settings turns it off.
  autoCheck: true,
  // Automatic warnings are for Risky pages unless this is on.
  warnSuspicious: false,
  // Warn before a password is typed on a site that isn't Safe.
  passwordGuard: true,
  // Hosts the user has chosen to trust: never warned about or checked automatically.
  trustedSites: [],
};

export async function getSettings() {
  const stored = await chrome.storage.sync.get(DEFAULTS);
  return { ...DEFAULTS, ...stored };
}

/** The full report for a link on the website. */
export function reportUrl(siteBase, url) {
  return `${siteBase.replace(/\/+$/, "")}/?q=${encodeURIComponent(url)}`;
}

/** Only web pages can be checked; not chrome://, file:// or the new tab page. */
export function isCheckable(url) {
  return /^https?:\/\//i.test(url ?? "");
}

export function hostOf(url) {
  try {
    return new URL(url).hostname.toLowerCase();
  } catch {
    return "";
  }
}

/** Whether host is one of the hosts listed, or a subdomain of one. */
export function hostIn(host, hosts) {
  host = host.replace(/^www\./, "");
  return hosts.some((h) => host === h || host.endsWith("." + h));
}
