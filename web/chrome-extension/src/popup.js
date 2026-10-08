import * as R from "./render.js";
import {
  DEFAULTS,
  extractLinks,
  formatUrl,
  getSettings,
  isCheckable,
  reportUrl,
} from "./shared.js";

const $ = (id) => document.getElementById(id);
const send = (msg) => chrome.runtime.sendMessage(msg);

let settings = DEFAULTS;
let tab = null;

function setMark(state) {
  $("mark").replaceChildren(R.wordmark(state, 30));
}

// ── This page ───────────────────────────────────────────────────────────────

function renderPage() {
  const box = $("page");
  if (!tab || !isCheckable(tab.url)) {
    box.replaceChildren(
      R.el("p", {
        class: "note",
        text: "Browser pages can't be checked. Paste a link below instead.",
      }),
    );
    return;
  }
  const u = new URL(tab.url);
  const path = u.pathname + u.search === "/" ? "" : u.pathname + u.search;
  box.replaceChildren(
    R.el(
      "div",
      { class: "page-info" },
      R.el("span", { class: "eyebrow", text: "This page" }),
      R.el(
        "p",
        { class: "page-url", title: tab.url },
        u.hostname,
        path ? R.el("span", { class: "path", text: path }) : null,
      ),
    ),
    R.el("button", {
      id: "checkPage",
      class: "pill pill-solid",
      type: "button",
      text: "Check this page",
      onclick: () => check(tab.url, tab.id),
    }),
  );
}

// ── Checking and results ────────────────────────────────────────────────────

function busy(on) {
  for (const id of ["checkBtn", "checkPage"]) {
    const b = $(id);
    if (b) b.disabled = on;
  }
}

async function check(url, tabId) {
  setMark("scanning");
  busy(true);
  $("result").replaceChildren(
    R.el(
      "div",
      { class: "checking" },
      R.el("p", { class: "lead", text: `Checking ${R.host(url)}…` }),
      R.el("p", { class: "hint", text: "This usually takes a few seconds." }),
    ),
  );
  const result = await send({ type: "scan", url, tabId }).catch(() => null);
  busy(false);
  if (!result || result.error) {
    setMark("idle");
    showError(
      result?.message ?? "Something went wrong. Try again.",
      url,
      tabId,
    );
    return;
  }
  showResult(result);
}

function showResult(result) {
  setMark(result.verdict);
  $("result").replaceChildren(
    ...R.notes(result),
    R.verdictCard(result),
    R.flags(result, 3),
    R.el(
      "div",
      { class: "result-actions" },
      R.el("a", {
        class: "pill pill-solid",
        href: result.report,
        target: "_blank",
        rel: "noopener noreferrer",
        text: "See the full report",
      }),
    ),
    R.el(
      "div",
      { class: "minor" },
      R.el("button", {
        class: "text-btn",
        type: "button",
        text: "Wrong result? Tell us",
        onclick: (e) => e.target.parentElement.replaceWith(reportForm(result)),
      }),
    ),
  );
}

// ── Report a wrong result (read by a person, like reports on the website) ───

function reportForm(result) {
  const options = ["Safe", "Suspicious", "Risky"].filter(
    (v) => v !== result.verdict,
  );
  const choice = (v, i) =>
    R.el(
      "label",
      {},
      R.el("input", {
        type: "radio",
        name: "expected",
        value: v,
        checked: i === 0,
      }),
      R.el("span", { text: v }),
    );
  const comment = R.el("textarea", {
    maxlength: "1000",
    placeholder: "What makes you think so? (optional)",
    "aria-label": "Comment",
  });
  const status = R.el("p", { class: "hint", role: "status" });
  const sendBtn = R.el("button", {
    class: "pill pill-solid",
    type: "button",
    text: "Send report",
  });
  const form = R.el(
    "div",
    { class: "report" },
    R.el("p", {
      class: "report-q",
      text: `We said ${result.verdict}. What should it be?`,
    }),
    R.el(
      "div",
      { class: "segmented", role: "radiogroup" },
      options.map(choice),
    ),
    comment,
    R.el("div", { class: "result-actions" }, sendBtn),
    status,
  );
  sendBtn.addEventListener("click", async () => {
    const expected = form.querySelector(
      'input[name="expected"]:checked',
    )?.value;
    sendBtn.disabled = true;
    status.textContent = "Sending…";
    const ok = await send({
      type: "report",
      payload: {
        url: result.url,
        verdict: result.verdict,
        score: result.score ?? 0,
        expected_verdict: expected,
        comment: comment.value.trim(),
      },
    }).catch(() => false);
    if (ok) {
      form.replaceChildren(
        R.el("p", {
          class: "note",
          text: "Thanks. A person reads every report, and it helps make url.vet better.",
        }),
      );
    } else {
      sendBtn.disabled = false;
      status.textContent = "Couldn't send that. Try again in a moment.";
    }
  });
  return form;
}

function showError(message, url, tabId) {
  $("result").replaceChildren(
    R.el("p", { class: "note warn", text: message }),
    R.el(
      "div",
      { class: "result-actions" },
      R.el("button", {
        class: "pill pill-outline",
        type: "button",
        text: "Try again",
        onclick: () => check(url, tabId),
      }),
    ),
  );
}

// A pasted link, a bare "example.com/path", or a whole message: links are found
// the way the website finds them, and several are offered as a choice.
function checkLink(url) {
  // If it's the page you're on, the toolbar dot can show the verdict too.
  check(url, url === tab?.url ? tab.id : null);
}

function chooseLink(links) {
  setMark("idle");
  $("result").replaceChildren(
    R.el(
      "div",
      { class: "choose" },
      R.el("p", {
        class: "lead",
        text: `Found ${links.length} links. Which one?`,
      }),
      R.el(
        "ul",
        { class: "choose-list" },
        links.map((l) =>
          R.el(
            "li",
            {},
            R.el("button", {
              type: "button",
              text: l,
              onclick: () => checkLink(l),
            }),
          ),
        ),
      ),
    ),
  );
}

$("form").addEventListener("submit", (e) => {
  e.preventDefault();
  const raw = $("url").value.trim();
  const links = raw ? extractLinks(raw).map(formatUrl) : [];
  $("form").classList.toggle("invalid", !links.length);
  $("formError").hidden = !!links.length;
  if (!links.length) {
    $("formError").textContent = raw
      ? "There's no link in that."
      : "Paste a link to check.";
    return;
  }
  if (links.length === 1) checkLink(links[0]);
  else chooseLink(links);
});
$("url").addEventListener("input", () => {
  $("form").classList.remove("invalid");
  $("formError").hidden = true;
});

// ── Settings ────────────────────────────────────────────────────────────────

const AUTO_ORIGINS = ["http://*/*", "https://*/*"];

// The dev stack from make dev-up.
const LOCAL = {
  apiBase: "http://localhost:8080/api/v1",
  siteBase: "http://localhost:5173",
};

function openSettings(open) {
  $("main").hidden = open;
  $("settings").hidden = !open;
  if (open) {
    $("autoCheck").checked = settings.autoCheck;
    $("warnSuspicious").checked = settings.warnSuspicious;
    $("warnSuspiciousRow").hidden = !settings.autoCheck;
    $("passwordGuard").checked = settings.passwordGuard;
    renderTrusted();
    $("apiBase").value = settings.apiBase;
    $("siteBase").value = settings.siteBase;
    $("serverStatus").textContent = "";
    $("closeSettings").focus();
  } else {
    $("openSettings").focus();
  }
}
$("openSettings").addEventListener("click", () => openSettings(true));
$("closeSettings").addEventListener("click", () => openSettings(false));

$("autoCheck").addEventListener("change", async (e) => {
  const on = e.target.checked;
  if (on) {
    // Access to all sites is granted at install; ask again only if it was removed in Chrome's settings.
    const granted = await chrome.permissions
      .request({ origins: AUTO_ORIGINS })
      .catch(() => false);
    if (!granted) {
      e.target.checked = false;
      return;
    }
  }
  settings.autoCheck = on;
  $("warnSuspiciousRow").hidden = !on;
  await chrome.storage.sync.set({ autoCheck: on });
  await send({ type: "sync-page" });
});

$("warnSuspicious").addEventListener("change", async (e) => {
  settings.warnSuspicious = e.target.checked;
  await chrome.storage.sync.set({ warnSuspicious: e.target.checked });
});

$("passwordGuard").addEventListener("change", async (e) => {
  settings.passwordGuard = e.target.checked;
  await chrome.storage.sync.set({ passwordGuard: e.target.checked });
  await send({ type: "sync-page" });
});

async function renderTrusted() {
  settings = await getSettings();
  const list = $("trustedList");
  if (!settings.trustedSites.length) {
    list.replaceChildren(R.el("li", { class: "empty", text: "None yet" }));
    return;
  }
  list.replaceChildren(
    ...settings.trustedSites.map((h) =>
      R.el(
        "li",
        {},
        h,
        R.el("button", {
          type: "button",
          "aria-label": `Stop trusting ${h}`,
          text: "×",
          onclick: async () => {
            await chrome.storage.sync.set({
              trustedSites: settings.trustedSites.filter((x) => x !== h),
            });
            renderTrusted();
          },
        }),
      ),
    ),
  );
}

// Add a site by hand: a host, or any link on it.
$("trustForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const raw = $("trustInput").value.trim();
  const [link] = raw ? extractLinks(raw) : [];
  let host = "";
  try {
    host = new URL(formatUrl(link ?? "")).hostname.replace(/^www\./, "");
  } catch {
    /* not a site */
  }
  $("trustError").hidden = !!host;
  if (!host) {
    $("trustError").textContent = raw
      ? "That doesn't look like a site."
      : "Type a site, like example.com.";
    return;
  }
  const { trustedSites } = await getSettings();
  if (!trustedSites.includes(host)) {
    await chrome.storage.sync.set({
      trustedSites: [...trustedSites, host].sort(),
    });
  }
  $("trustInput").value = "";
  renderTrusted();
});

function cleanBase(value) {
  try {
    const u = new URL(value.trim());
    if (!/^https?:$/.test(u.protocol)) return null;
    return u.toString().replace(/\/+$/, "");
  } catch {
    return null;
  }
}

$("saveServer").addEventListener("click", async () => {
  const status = $("serverStatus");
  const apiBase = cleanBase($("apiBase").value);
  const siteBase = cleanBase($("siteBase").value);
  status.classList.remove("bad");
  if (!apiBase || !siteBase) {
    status.classList.add("bad");
    status.textContent =
      "Both need to be full addresses, starting with http:// or https://.";
    return;
  }
  // A server other than url.vet's needs permission to be reached.
  const known = [DEFAULTS.apiBase, LOCAL.apiBase].includes(apiBase);
  if (!known) {
    const origin = new URL(apiBase).origin + "/*";
    const granted = await chrome.permissions
      .request({ origins: [origin] })
      .catch(() => false);
    if (!granted) {
      status.classList.add("bad");
      status.textContent = "url.vet needs permission to reach that server.";
      return;
    }
  }
  settings = { ...settings, apiBase, siteBase };
  await chrome.storage.sync.set({ apiBase, siteBase });
  await chrome.storage.session.clear(); // results from the old server no longer apply
  $("privacy").href = `${siteBase}/privacy`;
  status.textContent = "Saved.";
});

// One click to a server: fills both addresses and saves them.
async function useServer(apiBase, siteBase, message) {
  $("apiBase").value = apiBase;
  $("siteBase").value = siteBase;
  settings = { ...settings, apiBase, siteBase };
  await chrome.storage.sync.set({ apiBase, siteBase });
  await chrome.storage.session.clear(); // results from the old server no longer apply
  $("privacy").href = `${siteBase}/privacy`;
  $("serverStatus").classList.remove("bad");
  $("serverStatus").textContent = message;
}

$("resetServer").addEventListener("click", () =>
  useServer(DEFAULTS.apiBase, DEFAULTS.siteBase, "Using url.vet."),
);
$("localServer").addEventListener("click", () =>
  useServer(
    LOCAL.apiBase,
    LOCAL.siteBase,
    "Using your local server (make dev-up).",
  ),
);

// ── Start ───────────────────────────────────────────────────────────────────

(async () => {
  setMark("idle");
  settings = await getSettings();
  $("privacy").href = `${settings.siteBase.replace(/\/+$/, "")}/privacy`;
  [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  renderPage();
  // A page checked earlier in this session shows its result straight away.
  if (tab && isCheckable(tab.url)) {
    const hit = await send({ type: "cached", url: tab.url }).catch(() => null);
    if (hit)
      showResult({ ...hit, report: reportUrl(settings.siteBase, tab.url) });
  }
  if (!$("result").childElementCount) $("url").focus();
})();
