import { isKnownSite } from "./lists.js";
import {
  extractLinks,
  formatUrl,
  getSettings,
  hostIn,
  hostOf,
  isCheckable,
  reportUrl,
  verdictCopy,
} from "./shared.js";

// Results are kept for the browser session, so reopening the popup or revisiting
// a page doesn't scan again (and doesn't eat into the API's rate limit).
const CACHE_MS = 60 * 60 * 1000;

// Automatic checks stay well under the public API's 20 scans a minute.
const AUTO_PER_MINUTE = 12;
let autoTimes = [];

// The well-known sites list refreshes this often; the server sends it only when changed.
const LISTS_EVERY_MIN = 6 * 60;

const PAGE_SCRIPT_ID = "urlvet-page";

const icons = (name) => ({
  16: `icons/dot-${name}-16.png`,
  32: `icons/dot-${name}-32.png`,
});
const ICONS = {
  idle: icons("brand"),
  error: icons("grey"),
  Safe: icons("safe"),
  Suspicious: icons("suspicious"),
  Risky: icons("risky"),
};
// A letter as well as a colour, so the verdict doesn't rely on colour alone.
const BADGES = {
  Risky: { text: "!", color: "#ef4444", textColor: "#ffffff" },
  Suspicious: { text: "?", color: "#eab308", textColor: "#1a1714" },
};

// ── Scanning ────────────────────────────────────────────────────────────────

async function cached(url) {
  const key = "scan:" + url;
  const hit = (await chrome.storage.session.get(key))[key];
  return hit && Date.now() - hit.at < CACHE_MS ? hit.result : null;
}

async function scan(url) {
  const hit = await cached(url);
  if (hit) return hit;

  const { apiBase, siteBase } = await getSettings();
  let res;
  try {
    res = await fetch(
      `${apiBase.replace(/\/+$/, "")}/analyze?url=${encodeURIComponent(url)}`,
    );
  } catch {
    return {
      error: "offline",
      message: "Couldn't reach url.vet. Check your connection and try again.",
    };
  }
  if (res.status === 429) {
    return {
      error: "rate",
      message:
        "That's a lot of checks in a minute. Give it a moment and try again.",
    };
  }
  if (res.status === 400) {
    return {
      error: "invalid",
      message: "That doesn't look like a link url.vet can check.",
    };
  }
  if (!res.ok) {
    return {
      error: "server",
      message: `url.vet couldn't check this link right now (error ${res.status}).`,
    };
  }

  const d = await res.json();
  const r = d.result ?? {};
  const verdict = ["Safe", "Suspicious", "Risky"].includes(r.verdict)
    ? r.verdict
    : "Suspicious";
  const copy = verdictCopy(verdict);
  const result = {
    url,
    scannedUrl: d.url ?? url,
    verdict,
    score: typeof r.final_score === "number" ? r.final_score : null,
    label: copy.label,
    quip: copy.quip,
    bad: r.reasons?.bad_reasons ?? [],
    good: r.reasons?.good_reasons ?? [],
    shortLink: d.short_link ?? null,
    redirectedFrom: d.redirected_from ?? null,
    incomplete: !!d.incomplete,
    incompleteChecks: d.incomplete_checks ?? [],
    report: reportUrl(siteBase, url),
  };
  // Incomplete results aren't kept, so the next look retries the missing checks.
  if (!result.incomplete) {
    await chrome.storage.session.set({
      ["scan:" + url]: { at: Date.now(), result },
    });
  }
  return result;
}

// ── The toolbar dot: the verdict's colour and letter on that tab ────────────

async function markTab(tabId, result) {
  if (tabId == null) return;
  let state = "idle";
  let title = "url.vet: check this page or a link";
  if (result?.error) {
    state = "error";
    title = `url.vet couldn't check this page. ${result.message}`;
  } else if (result) {
    state = result.verdict;
    title = `url.vet: ${result.verdict} (trust score ${result.score ?? "?"} of 100)`;
  }
  const badge = BADGES[state];
  await Promise.all([
    chrome.action.setIcon({ tabId, path: ICONS[state] }),
    chrome.action.setTitle({ tabId, title }),
    chrome.action.setBadgeText({ tabId, text: badge?.text ?? "" }),
    badge &&
      chrome.action.setBadgeBackgroundColor({ tabId, color: badge.color }),
    badge && chrome.action.setBadgeTextColor({ tabId, color: badge.textColor }),
  ]).catch(() => {});
}

// ── Well-known sites ────────────────────────────────────────────────────────
// The server's list of the 10,000 best-known sites (and the hosts on them
// where anyone can publish). Pages on those sites aren't sent for automatic
// checks; see shouldCheck below.

let known = null; // { sites: Set, userContent: Set }

async function loadKnown() {
  if (known) return known;
  const { known: stored } = await chrome.storage.local.get("known");
  known = {
    sites: new Set(stored?.sites ?? []),
    userContent: new Set(stored?.userContent ?? []),
  };
  return known;
}

/** Downloads the list if it changed. A failure keeps the copy already held. */
async function refreshKnown() {
  const { apiBase } = await getSettings();
  const { known: stored } = await chrome.storage.local.get("known");
  try {
    const res = await fetch(
      `${apiBase.replace(/\/+$/, "")}/lists/known-sites`,
      {
        headers: stored?.etag ? { "If-None-Match": stored.etag } : {},
      },
    );
    if (res.status === 304 && stored) {
      await chrome.storage.local.set({ known: { ...stored, at: Date.now() } });
    } else if (res.ok) {
      const body = await res.json();
      await chrome.storage.local.set({
        known: {
          etag: res.headers.get("ETag"),
          sites: body.sites,
          userContent: body.user_content,
          at: Date.now(),
        },
      });
    }
  } catch {
    /* kept the old list */
  }
  known = null; // reload on next use
}

async function isWellKnown(url) {
  const { sites, userContent } = await loadKnown();
  return isKnownSite(hostOf(url), sites, userContent);
}

chrome.alarms.onAlarm.addListener((a) => {
  if (a.name === "urlvet-lists") refreshKnown();
});

// ── Stop pages found Risky earlier, before they load again ─────────────────

async function allowedOnce(url) {
  const key = "allow:" + url;
  return !!(await chrome.storage.session.get(key))[key];
}

chrome.webNavigation.onBeforeNavigate.addListener(async (nav) => {
  if (nav.frameId !== 0 || !isCheckable(nav.url)) return;
  if ((await cached(nav.url))?.verdict !== "Risky") return;
  const { trustedSites } = await getSettings();
  if (hostIn(hostOf(nav.url), trustedSites) || (await allowedOnce(nav.url)))
    return;

  const blocked = chrome.runtime.getURL(
    `blocked.html?u=${encodeURIComponent(nav.url)}`,
  );
  chrome.tabs.update(nav.tabId, { url: blocked }).catch(() => {});
});

// ── On-page card (right-click checks) ───────────────────────────────────────

async function showCard(tabId, payload) {
  try {
    await chrome.scripting.executeScript({
      target: { tabId },
      files: ["card.js"],
    });
    await chrome.tabs.sendMessage(tabId, { type: "urlvet:card", payload });
    return true;
  } catch {
    // Browser pages and the Web Store can't show the card.
    return false;
  }
}

// The first link in a selection, found the way the website finds links in a pasted message.
function extractLink(text) {
  const [first] = extractLinks(text ?? "");
  return first ? formatUrl(first) : null;
}

chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.removeAll(() => {
    chrome.contextMenus.create({
      id: "check-link",
      title: "Check this link with url.vet",
      contexts: ["link"],
    });
    chrome.contextMenus.create({
      id: "check-selection",
      title: "Check selected link with url.vet",
      contexts: ["selection"],
    });
  });
  chrome.alarms.create("urlvet-lists", { periodInMinutes: LISTS_EVERY_MIN });
  chrome.storage.local.remove("threats"); // left by earlier versions, which kept a threat list
  refreshKnown();
  syncPageScript();
});

chrome.runtime.onStartup.addListener(async () => {
  syncPageScript();
  const { known: stored } = await chrome.storage.local.get("known");
  if (!stored || Date.now() - stored.at > LISTS_EVERY_MIN * 60_000)
    refreshKnown();
});

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  const url =
    info.menuItemId === "check-link"
      ? info.linkUrl
      : extractLink(info.selectionText);
  if (!tab?.id) return;

  if (!url || !isCheckable(url)) {
    await showCard(tab.id, {
      state: "error",
      message: "There's no web link in that selection to check.",
    });
    return;
  }
  const shown = await showCard(tab.id, { state: "checking", url });
  const result = await scan(url);
  if (!shown) {
    // No card on this page: open the full report instead.
    const { siteBase } = await getSettings();
    chrome.tabs.create({ url: reportUrl(siteBase, url) });
    return;
  }
  await showCard(
    tab.id,
    result.error
      ? { state: "error", url, message: result.message }
      : { state: "result", ...result },
  );
});

// ── Page script: automatic warnings and the password guard ──────────────────

async function syncPageScript() {
  const { autoCheck, passwordGuard } = await getSettings();
  const want = autoCheck || passwordGuard;
  const allowed = await chrome.permissions.contains({
    origins: ["http://*/*", "https://*/*"],
  });
  const existing = await chrome.scripting.getRegisteredContentScripts({
    ids: [PAGE_SCRIPT_ID],
  });
  if (want && allowed && !existing.length) {
    await chrome.scripting.registerContentScripts([
      {
        id: PAGE_SCRIPT_ID,
        matches: ["http://*/*", "https://*/*"],
        js: ["page.js"],
        runAt: "document_idle",
        allFrames: false,
        persistAcrossSessions: true,
      },
    ]);
  } else if ((!want || !allowed) && existing.length) {
    await chrome.scripting.unregisterContentScripts({ ids: [PAGE_SCRIPT_ID] });
  }
}

chrome.storage.onChanged.addListener((changes, area) => {
  if (area === "sync" && ("autoCheck" in changes || "passwordGuard" in changes))
    syncPageScript();
  if (area === "sync" && ("apiBase" in changes || "siteBase" in changes)) {
    chrome.storage.local.remove("known").then(refreshKnown);
  }
});
chrome.permissions.onRemoved.addListener(syncPageScript);
chrome.permissions.onAdded.addListener(syncPageScript);

function underAutoLimit() {
  const now = Date.now();
  autoTimes = autoTimes.filter((t) => now - t < 60_000);
  if (autoTimes.length >= AUTO_PER_MINUTE) return false;
  autoTimes.push(now);
  return true;
}

/**
 * Whether a page should be checked automatically: not our own site, not one
 * the user trusts, and not a well-known site (those are never sent).
 */
async function shouldCheck(url) {
  if (!isCheckable(url)) return false;
  const { apiBase, siteBase, trustedSites } = await getSettings();
  const host = hostOf(url);
  if ([hostOf(apiBase), hostOf(siteBase)].includes(host)) return false;
  if (hostIn(host, trustedSites)) return false;
  return !(await isWellKnown(url));
}

async function checkPage(url, tabId) {
  const hit = await cached(url);
  if (hit) {
    await markTab(tabId, hit);
    return hit;
  }
  if (!underAutoLimit()) return null;
  const result = await scan(url);
  await markTab(tabId, result);
  return result;
}

async function autoCheck(url, tabId) {
  const { autoCheck, warnSuspicious } = await getSettings();
  if (!autoCheck || !(await shouldCheck(url))) return null;
  const result = await checkPage(url, tabId);
  if (!result || result.error) return null;
  // Only Risky pages interrupt (Suspicious too if asked); others just colour the dot.
  return result.verdict === "Risky" ||
    (warnSuspicious && result.verdict === "Suspicious")
    ? result
    : null;
}

/** About to type a password: warn unless the page checks out as Safe. */
async function guard(url, tabId) {
  const { passwordGuard } = await getSettings();
  if (!passwordGuard || !(await shouldCheck(url))) return null;
  const result = await checkPage(url, tabId);
  if (!result || result.error || result.verdict === "Safe") return null;
  return result;
}

async function trustSite(url) {
  const host = hostOf(url).replace(/^www\./, "");
  const { trustedSites } = await getSettings();
  if (host && !trustedSites.includes(host)) {
    await chrome.storage.sync.set({
      trustedSites: [...trustedSites, host].sort(),
    });
  }
}

async function sendReport(payload) {
  const { apiBase } = await getSettings();
  try {
    const res = await fetch(`${apiBase.replace(/\/+$/, "")}/report`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return res.ok;
  } catch {
    return false;
  }
}

// ── Messages from the popup, the blocked page and pages ─────────────────────

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  (async () => {
    switch (msg?.type) {
      case "scan": {
        const result = await scan(msg.url);
        if (msg.tabId != null) await markTab(msg.tabId, result);
        return result;
      }
      case "cached":
        return cached(msg.url);
      case "auto":
        return autoCheck(msg.url, sender.tab?.id);
      case "guard":
        return guard(msg.url, sender.tab?.id);
      case "trust":
        await trustSite(msg.url);
        return true;
      case "blocked":
        // The blocked page stands in for a Risky one: show that on the toolbar.
        await markTab(sender.tab?.id, { verdict: "Risky", score: null });
        return true;
      case "allow-once":
        await chrome.storage.session.set({ ["allow:" + msg.url]: true });
        return true;
      case "report":
        return sendReport(msg.payload);
      case "sync-page":
        await syncPageScript();
        return true;
      default:
        return null;
    }
  })().then(sendResponse, () => sendResponse(null));
  return true; // answered asynchronously
});

// A new page starts with the plain dot again.
chrome.tabs.onUpdated.addListener((tabId, info) => {
  if (info.status === "loading") markTab(tabId, null);
});
