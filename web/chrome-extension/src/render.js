// Builds the url.vet result UI as DOM nodes, for the popup, the blocked page
// and the on-page card. Text is always set with textContent: scanned URLs and
// reasons come from the network and must never be parsed as HTML.
import { incompleteNote, redirectNote, shortLinkNote } from "./shared.js";

const SVG = "http://www.w3.org/2000/svg";

export function el(tag, props = {}, ...children) {
  const node = document.createElement(tag);
  for (const [k, v] of Object.entries(props)) {
    if (v == null || v === false) continue;
    if (k === "class") node.className = v;
    else if (k === "text") node.textContent = v;
    else if (k.startsWith("on")) node.addEventListener(k.slice(2), v);
    else node.setAttribute(k, v === true ? "" : v);
  }
  for (const c of children.flat()) if (c) node.append(c);
  return node;
}

export function host(url) {
  try {
    return new URL(url).hostname;
  } catch {
    return url;
  }
}

/** The url.vet wordmark. state: "idle" | "scanning" | a verdict. */
export function wordmark(state = "idle", size = 26) {
  return el(
    "span",
    {
      class: "wordmark",
      style: `font-size:${size}px`,
      "aria-label": "url.vet",
    },
    el("span", { class: "url", text: "url", "aria-hidden": "true" }),
    el("span", {
      class: "dot",
      text: ".",
      "data-state": state,
      "aria-hidden": "true",
    }),
    el("span", { text: "vet", "aria-hidden": "true" }),
  );
}

/** The trust score ring, filling and counting up to the score. */
function ring(score) {
  const R = 36;
  const CIRC = 2 * Math.PI * R;
  const svg = document.createElementNS(SVG, "svg");
  svg.setAttribute("viewBox", "0 0 88 88");
  svg.setAttribute("aria-hidden", "true");
  const circle = (cls) => {
    const c = document.createElementNS(SVG, "circle");
    for (const [k, v] of Object.entries({
      cx: 44,
      cy: 44,
      r: R,
      fill: "none",
      "stroke-width": 5,
    }))
      c.setAttribute(k, v);
    c.setAttribute("class", cls);
    return c;
  };
  const arc = circle("arc");
  arc.setAttribute("stroke-linecap", "round");
  arc.setAttribute("stroke-dasharray", CIRC);
  arc.setAttribute("stroke-dashoffset", CIRC);
  svg.append(circle("track"), arc);

  const value = el("span", { class: "score", text: "0" });
  const box = el(
    "div",
    {
      class: "ring",
      role: "img",
      "aria-label": `Trust score ${score ?? "unknown"} out of 100`,
    },
    svg,
    el(
      "div",
      { class: "value", "aria-hidden": "true" },
      value,
      el("span", { class: "rule" }),
      el("span", { class: "of", text: "100" }),
    ),
  );

  if (typeof score === "number") {
    requestAnimationFrame(() =>
      requestAnimationFrame(() =>
        arc.setAttribute("stroke-dashoffset", CIRC - (score / 100) * CIRC),
      ),
    );
    const start = performance.now();
    const tick = (now) => {
      const t = Math.min((now - start) / 1600, 1);
      value.textContent = String(Math.round((1 - Math.pow(1 - t, 3)) * score));
      if (t < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  } else {
    value.textContent = "—";
  }
  return box;
}

/** Verdict word stamped in, its badge, the one-line quip, and the score ring. */
export function verdictCard(result) {
  return el(
    "div",
    { class: "verdict", "data-verdict": result.verdict },
    el(
      "div",
      { class: "verdict-main" },
      el("span", { class: "eyebrow", text: "Verdict" }),
      el(
        "div",
        { class: "verdict-line" },
        el("span", { class: "verdict-word", text: result.verdict }),
        el("span", { class: "badge", text: result.label }),
      ),
      el("p", { class: "quip", text: result.quip }),
    ),
    ring(result.score),
  );
}

/** Red flags for a warning, green flags for a safe result; the full list lives in the report. */
export function flags(result, max = 3) {
  const bad = result.verdict !== "Safe" && result.bad.length > 0;
  const items = bad ? result.bad : result.good;
  if (!items.length) return null;
  const shown = items.slice(0, max);
  const more = items.length - shown.length;
  return el(
    "section",
    {},
    el("h3", {
      class: `eyebrow flags-title ${bad ? "bad" : "good"}`,
      text: bad ? "Red flags" : "Green flags",
    }),
    el(
      "ul",
      { class: `flags ${bad ? "bad" : "good"}` },
      shown.map((r) => el("li", { text: r })),
    ),
    more > 0
      ? el("p", { class: "more", text: `+${more} more in the full report` })
      : null,
  );
}

/** Short-link and incomplete-scan notes, in the website's own words. */
export function notes(result) {
  const out = [];
  const short =
    shortLinkNote(result.shortLink ?? undefined) ??
    redirectNote(result.redirectedFrom ?? undefined);
  if (short) out.push(el("p", { class: "note", text: short }));
  const missing = result.incomplete
    ? incompleteNote(result.incompleteChecks, true)
    : null;
  if (missing) out.push(el("p", { class: "note warn", text: missing }));
  return out;
}
