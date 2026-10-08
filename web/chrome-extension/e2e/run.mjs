// Loads dist/ into Chromium and checks the extension end to end against a
// local url.vet backend (make dev-up). Screenshots go to e2e/shots/.
//
//   CHROME_PATH=/path/to/chrome npm run test:e2e
//
// CHROME_PATH must be a full Chrome or Chromium (not headless-shell): only
// those load extensions. Defaults to Playwright's Chrome for Testing.
import { chromium } from "playwright-core";
import { existsSync, mkdirSync, mkdtempSync, readdirSync } from "node:fs";
import { homedir, tmpdir } from "node:os";
import path from "node:path";

const EXT = path.resolve("dist");
const SHOTS = path.resolve("e2e/shots");
const API = process.env.API_BASE ?? "http://localhost:8080/api/v1";
const SITE = process.env.SITE_BASE ?? "http://localhost:5173";
const PHISHING = "http://testsafebrowsing.appspot.com/s/phishing.html"; // Google's test page
mkdirSync(SHOTS, { recursive: true });

function chromePath() {
  if (process.env.CHROME_PATH) return process.env.CHROME_PATH;
  const cache = path.join(homedir(), ".cache/ms-playwright");
  const dir =
    existsSync(cache) &&
    readdirSync(cache).find((d) => /^chromium-\d+$/.test(d));
  if (!dir) throw new Error("Set CHROME_PATH to a Chrome or Chromium binary.");
  return path.join(cache, dir, "chrome-linux64/chrome");
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
let failures = 0;
function check(name, ok, detail = "") {
  console.log(`${ok ? "✓" : "✗"} ${name}${detail ? `  (${detail})` : ""}`);
  if (!ok) failures++;
}

// Our redirect aborts the browser's own navigation, so poll for where the tab ends up.
async function landsOnBlocked(page, url) {
  await page.goto(url, { waitUntil: "commit" }).catch(() => {});
  for (let i = 0; i < 40; i++) {
    if (page.url().includes("blocked.html")) return true;
    await sleep(250);
  }
  return false;
}

const ctx = await chromium.launchPersistentContext(
  mkdtempSync(path.join(tmpdir(), "urlvet-e2e-")),
  {
    executablePath: chromePath(),
    headless: true,
    viewport: { width: 1100, height: 760 },
    args: [
      `--disable-extensions-except=${EXT}`,
      `--load-extension=${EXT}`,
      "--headless=new",
    ],
  },
);
const sw = ctx.serviceWorkers()[0] ?? (await ctx.waitForEvent("serviceworker"));
const id = sw.url().split("/")[2];
const tabTitle = (host) =>
  sw.evaluate(async (h) => {
    const [t] = await chrome.tabs.query({ url: `*://${h}/*` });
    return t ? chrome.action.getTitle({ tabId: t.id }) : "no tab";
  }, host);

try {
  // Point at the local backend; this also downloads the well-known sites from it.
  await sw.evaluate(
    ([apiBase, siteBase]) => chrome.storage.sync.set({ apiBase, siteBase }),
    [API, SITE],
  );
  for (let i = 0; i < 40; i++) {
    const { known } = await sw.evaluate(() =>
      chrome.storage.local.get("known"),
    );
    if (known?.sites?.length) break;
    await sleep(500);
  }
  const { known } = await sw.evaluate(() => chrome.storage.local.get("known"));
  check(
    "downloads the well-known sites",
    (known?.sites?.length ?? 0) > 0,
    `${known?.sites?.length} sites`,
  );
  check(
    "keeps no threat list",
    !(await sw.evaluate(() => chrome.storage.local.get("threats"))).threats,
  );
  const scripts = await sw.evaluate(async () =>
    (await chrome.scripting.getRegisteredContentScripts()).map((s) => s.id),
  );
  check("page script registered on install", scripts.includes("urlvet-page"));

  // Well-known sites aren't sent anywhere.
  const wiki = await ctx.newPage();
  await wiki.goto("https://www.wikipedia.org/", { waitUntil: "load" });
  await sleep(3000);
  check(
    "well-known site isn't checked",
    (await tabTitle("www.wikipedia.org")).startsWith(
      "url.vet: check this page",
    ),
  );
  await wiki.close();

  // A risky page gets the automatic warning.
  const risky = await ctx.newPage();
  await risky.goto(PHISHING, { waitUntil: "load" });
  await risky
    .waitForSelector("urlvet-card", { state: "attached", timeout: 60000 })
    .catch(() => {});
  await sleep(2200);
  check(
    "automatic warning on a risky page",
    await risky.evaluate(() => !!document.querySelector("urlvet-card")),
  );
  check(
    "toolbar shows the verdict",
    (await tabTitle("testsafebrowsing.appspot.com")).includes("Risky"),
  );
  await risky.screenshot({ path: `${SHOTS}/auto-warning.png` });
  await risky.close();

  // Found Risky earlier: stopped before it loads next time.
  const again = await ctx.newPage();
  check(
    "page found Risky earlier is stopped next time",
    await landsOnBlocked(again, PHISHING),
  );
  await sleep(600);
  await again.screenshot({ path: `${SHOTS}/blocked.png` });
  await again.close();

  // The password guard, with automatic warnings off.
  await sw.evaluate(() => chrome.storage.session.clear());
  await sw.evaluate(() => chrome.storage.sync.set({ autoCheck: false }));
  await sleep(500);
  const login = await ctx.newPage();
  await login.goto(PHISHING, { waitUntil: "load" });
  await sleep(1500);
  check(
    "no automatic warning when turned off",
    !(await login.evaluate(() => !!document.querySelector("urlvet-card"))),
  );
  await login
    .focus('input[type="password"], input[name*="pass" i]')
    .catch(() => {});
  await login
    .waitForSelector("urlvet-card", { state: "attached", timeout: 60000 })
    .catch(() => {});
  await sleep(2200);
  check(
    "password guard warns on focus",
    await login.evaluate(() => !!document.querySelector("urlvet-card")),
  );
  await login.screenshot({ path: `${SHOTS}/password-guard.png` });
  await login.close();
  await sw.evaluate(() => chrome.storage.sync.set({ autoCheck: true }));

  // Trusting a site stops checks and warnings there.
  await sw.evaluate(() => chrome.storage.session.clear());
  await sw.evaluate(() =>
    chrome.storage.sync.set({ trustedSites: ["testsafebrowsing.appspot.com"] }),
  );
  const trusted = await ctx.newPage();
  await trusted.goto(PHISHING, { waitUntil: "load" });
  await sleep(4000);
  check(
    "trusted site isn't warned about",
    !(await trusted.evaluate(() => !!document.querySelector("urlvet-card"))),
  );
  await trusted.close();
  await sw.evaluate(() => chrome.storage.sync.set({ trustedSites: [] }));

  // Popup: a pasted message with two links, a result, the report form.
  const pop = await ctx.newPage();
  await pop.setViewportSize({ width: 380, height: 700 });
  await pop.goto(`chrome-extension://${id}/popup.html`);
  await pop.fill(
    "#url",
    "Hi! Your parcel: https://goo.su/TIcsAp or track it at https://dhl-tracking.example/now.",
  );
  await pop.click("#checkBtn");
  check(
    "popup offers a choice of links",
    (await pop.locator(".choose-list button").count()) === 2,
  );
  await sleep(500); // let the fade-in finish
  await pop.screenshot({ path: `${SHOTS}/popup-choose.png` });
  await pop.locator(".choose-list button").first().click();
  await pop.waitForSelector(".verdict", { timeout: 120000 });
  await sleep(2000);
  await pop.screenshot({ path: `${SHOTS}/popup-result.png`, fullPage: true });
  await pop.click("text=Wrong result? Tell us");
  check("report form opens", await pop.isVisible(".report"));
  await sleep(500);
  await pop.screenshot({ path: `${SHOTS}/popup-report.png`, fullPage: true });

  // Settings: no threat list panel, trusting a site by hand, one-click servers.
  await pop.click("#openSettings");
  await sleep(500);
  check("no threat list panel", !(await pop.$("#listsStatus")));
  await pop.fill("#trustInput", "https://www.Example.org/some/page");
  await pop.click("#trustForm button");
  await sleep(300);
  check(
    "trust form adds the site",
    (await pop.textContent("#trustedList")).includes("example.org"),
  );
  const { trustedSites } = await sw.evaluate(() =>
    chrome.storage.sync.get("trustedSites"),
  );
  check("trusted site is saved", trustedSites.includes("example.org"));
  await pop.click("#resetServer");
  await sleep(200);
  check(
    "Use url.vet fills the public server",
    (await pop.inputValue("#apiBase")) === "https://api.url.vet/api/v1",
  );
  await pop.click("#localServer");
  await sleep(200);
  const saved = await sw.evaluate(() =>
    chrome.storage.sync.get(["apiBase", "siteBase"]),
  );
  check(
    "Use localhost fills and saves both",
    saved.apiBase === "http://localhost:8080/api/v1" &&
      saved.siteBase === "http://localhost:5173",
  );
  await pop.screenshot({ path: `${SHOTS}/popup-settings.png`, fullPage: true });
} finally {
  await ctx.close();
}

console.log(failures ? `\n${failures} check(s) failed` : "\nAll checks passed");
process.exit(failures ? 1 : 0);
