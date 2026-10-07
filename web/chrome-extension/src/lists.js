// The list of well-known sites from the server's /lists/known-sites endpoint.
// Pages on these sites aren't sent for automatic checks, except on the hosts
// where anyone can publish.

/**
 * A well-known site that needn't be sent for an automatic check: the host or a
 * parent is in the top sites, and neither is a host anyone can publish on
 * (docs.google.com, storage.googleapis.com, someone.github.io).
 */
export function isKnownSite(host, sites, userContent) {
  const labels = host
    .toLowerCase()
    .replace(/^www\./, "")
    .split(".");
  let known = false;
  for (let i = 0; i < labels.length - 1; i++) {
    const suffix = labels.slice(i).join(".");
    if (userContent.has(suffix)) return false;
    if (sites.has(suffix)) known = true;
  }
  return known;
}
