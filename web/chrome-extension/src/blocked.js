// Shown in place of a page url.vet found Risky earlier this session, before it
// loads again.
import * as R from "./render.js";
import { getSettings, reportUrl } from "./shared.js";

const url = new URLSearchParams(location.search).get("u") ?? "";
const $ = (id) => document.getElementById(id);

$("mark").replaceChildren(R.wordmark("Risky", 34));
$("link").textContent = url;

chrome.runtime.sendMessage({ type: "blocked" });

getSettings().then(({ siteBase }) => {
  $("report").href = reportUrl(siteBase, url);
});

// A link opened in a new tab has nowhere to go back to.
if (history.length <= 1) $("back").textContent = "Close this tab";
$("back").addEventListener("click", async () => {
  if (history.length > 1) history.back();
  else chrome.tabs.remove((await chrome.tabs.getCurrent()).id);
});

$("continue").addEventListener("click", async () => {
  await chrome.runtime.sendMessage({ type: "allow-once", url });
  location.replace(url);
});
