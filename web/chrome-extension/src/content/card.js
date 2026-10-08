// The url.vet card shown on a page: a right-click result, an automatic warning,
// or the password guard. It lives in a closed shadow root so the page's styles
// can't reach it, and it's built only from text nodes (see render.js).
import uiCss from "../../public/ui.css";
import * as R from "../render.js";

const CARD_CSS = `
  :host { all: initial; }
  .card {
    position: fixed; right: 16px; bottom: 16px; z-index: 2147483647;
    width: 344px; max-width: calc(100vw - 32px); max-height: calc(100vh - 32px); overflow: auto;
    box-sizing: border-box; padding: 16px;
    display: flex; flex-direction: column; gap: 14px;
    background: var(--bg); color: var(--text);
    border: 1px solid var(--line); border-radius: 20px;
    box-shadow: 0 18px 48px rgb(15 13 11 / 0.18), 0 2px 6px rgb(15 13 11 / 0.08);
    font: 14px/1.45 var(--sans);
    -webkit-font-smoothing: antialiased;
    animation: urlvet-fade-in 0.3s ease-out both;
  }
  .card * { box-sizing: border-box; }
  .top { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
  .close {
    display: inline-flex; align-items: center; justify-content: center;
    width: 28px; height: 28px; border-radius: 999px; border: 0; padding: 0;
    background: transparent; color: var(--muted); cursor: pointer; font: 20px/1 var(--sans);
  }
  .close:hover { color: var(--text); background: var(--line); }
  .close:focus-visible, .text-btn:focus-visible { outline: 2px solid var(--faint); outline-offset: 2px; }
  .link {
    margin: 0; font: 12px/1.4 var(--mono); color: var(--muted);
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  }
  .checking { margin: 0; color: var(--soft); }
  .lead { margin: 0; color: var(--text); font-weight: 500; }
  .actions { display: flex; flex-wrap: wrap; gap: 8px; }
  .minor { display: flex; flex-wrap: wrap; gap: 4px 14px; margin-top: -4px; }
  .text-btn {
    border: 0; padding: 0; background: none; cursor: pointer;
    font: 12.5px/1.4 var(--sans); color: var(--muted);
    text-decoration: underline; text-underline-offset: 3px; text-decoration-color: var(--line-strong);
  }
  .text-btn:hover { color: var(--text); }
  .verdict-word { font-size: 34px; }
  .ring { width: 58px; height: 58px; }
`;

const asset = (p) => chrome.runtime.getURL(p);

// Fonts go under url.vet's own names, so they can't clash with the page's.
function loadFonts() {
  const fonts = [
    ["urlvet Geist", "fonts/geist-normal-latin.woff2", { weight: "300 700" }],
    [
      "urlvet Geist Mono",
      "fonts/geist-mono-normal-latin.woff2",
      { weight: "300 700" },
    ],
    [
      "urlvet Instrument Serif",
      "fonts/instrument-serif-normal-latin.woff2",
      {},
    ],
  ];
  for (const [family, file, desc] of fonts) {
    try {
      const face = new FontFace(family, `url(${asset(file)})`, {
        display: "swap",
        ...desc,
      });
      document.fonts.add(face);
      face.load().catch(() => {});
    } catch {
      /* system fonts are the fallback */
    }
  }
}

export function installCard() {
  if (window.__urlvetCard) return window.__urlvetCard;
  loadFonts();

  const hostEl = document.createElement("urlvet-card");
  const root = hostEl.attachShadow({ mode: "closed" });
  // Constructed stylesheets, not <link> or <style>: a page's Content Security
  // Policy can block those, and the card would show up unstyled.
  const sheet = (css) => {
    const s = new CSSStyleSheet();
    s.replaceSync(css);
    return s;
  };
  root.adoptedStyleSheets = [sheet(uiCss), sheet(CARD_CSS)];
  const card = R.el("div", {
    class: "card",
    role: "dialog",
    "aria-live": "polite",
    "aria-label": "url.vet",
  });
  root.append(card);

  let current = null;
  const dismiss = () => {
    hostEl.remove();
    current = null;
  };
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && hostEl.isConnected) dismiss();
  });

  const top = (state) =>
    R.el(
      "div",
      { class: "top" },
      R.wordmark(state, 22),
      R.el("button", {
        class: "close",
        type: "button",
        "aria-label": "Dismiss",
        text: "×",
        onclick: dismiss,
      }),
    );

  const textBtn = (text, onclick) =>
    R.el("button", { class: "text-btn", type: "button", text, onclick });

  function show(p) {
    current = p;
    card.replaceChildren();
    if (p.state === "checking") {
      card.append(
        top("scanning"),
        R.el("p", { class: "checking", text: "Checking this link…" }),
        R.el("p", { class: "link", text: p.url, title: p.url }),
      );
    } else if (p.state === "error") {
      card.append(
        top("idle"),
        R.el("p", { class: "checking", text: p.message }),
        p.url ? R.el("p", { class: "link", text: p.url, title: p.url }) : null,
      );
    } else {
      // A warning about the page itself (automatic check or password guard),
      // or the result of checking some other link (right-click).
      const aboutPage = p.mode === "auto" || p.mode === "guard";
      const lead =
        p.mode === "guard"
          ? "Careful with your password here."
          : p.mode === "auto"
            ? p.verdict === "Risky"
              ? "This page looks risky."
              : "Something about this page is off."
            : null;
      const leaving =
        p.mode === "auto" && p.verdict === "Risky" && history.length > 1;
      card.append(
        top(p.verdict),
        lead
          ? R.el("p", { class: "lead", text: lead })
          : R.el("p", { class: "link", text: p.url, title: p.url }),
        ...R.notes(p),
        R.verdictCard(p),
        R.flags(p, 3),
        R.el(
          "div",
          { class: "actions" },
          R.el("a", {
            class: "pill pill-solid",
            href: p.report,
            target: "_blank",
            rel: "noopener noreferrer",
            text: "See the full report",
          }),
          leaving
            ? R.el("button", {
                class: "pill pill-outline",
                type: "button",
                text: "Take me back",
                onclick: () => history.back(),
              })
            : null,
        ),
        R.el(
          "div",
          { class: "minor" },
          aboutPage
            ? textBtn("Trust this site", () => {
                chrome.runtime.sendMessage({
                  type: "trust",
                  url: location.href,
                });
                dismiss();
              })
            : null,
          // The full report has the website's report form.
          textBtn("Wrong result?", () =>
            window.open(p.report, "_blank", "noopener"),
          ),
        ),
      );
    }
    if (!hostEl.isConnected) document.documentElement.append(hostEl);
  }

  chrome.runtime.onMessage.addListener((msg) => {
    if (msg?.type === "urlvet:card") show(msg.payload);
  });

  window.__urlvetCard = { show, dismiss, current: () => current };
  return window.__urlvetCard;
}
