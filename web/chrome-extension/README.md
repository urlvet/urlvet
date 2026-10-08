# url.vet for Chrome

Check links and pages with url.vet without leaving the tab.

- **Stops repeat visits to risky pages.** A page url.vet found Risky earlier in the session is replaced by a warning before it loads again.
- **Warns on risky pages** (on by default): other pages are checked as they open, and Risky ones get a warning card (Suspicious too, if you tick it). The 10,000 best-known sites are never sent; for the rest, only the address is. This is why the extension asks for access to all sites when installed.
- **Guards passwords** (on by default): before you type a password on a page that isn't Safe, you get a warning.
- **Popup:** check the page you're on, or paste a link or a whole message (several links are offered as a choice). Report a wrong result from there.
- **Right-click** a link or selected text: *Check this link with url.vet*.
- **Toolbar dot:** green, yellow or red per tab, with **!** or **?** on it so it doesn't rely on colour alone; grey when a check couldn't run.
- **Trusted sites:** *Trust this site* on a warning, or add one in Settings, stops checks and warnings there.
- **Server:** one-click *Use url.vet* or *Use localhost* (`make dev-up`) in Settings, or any self-hosted address.

## Try it locally

```bash
cd web/chrome-extension
npm install
npm run build        # or: npm run watch
```

1. Open `chrome://extensions`, turn on **Developer mode**, click **Load unpacked** and choose `web/chrome-extension/dist`.
2. It uses the public url.vet API by default. For a local backend (`make dev-up`), open the popup's settings and click **Use localhost**.

After a rebuild, press reload on the extension's card in `chrome://extensions` and refresh open tabs. Errors from the popup show under **Inspect popup** (right-click the toolbar icon); from the background worker under **service worker** on the extension's card.

## Tests

```bash
npm test             # unit tests (vitest)
npm run test:e2e     # loads dist/ in Chromium against the local backend; needs make dev-up
```

The end-to-end run checks every feature above and writes screenshots to `e2e/shots/`. It needs a full Chrome or Chromium (headless-shell can't load extensions): set `CHROME_PATH`, or it uses Playwright's Chrome for Testing if installed. From the repo root: `make test-extension` and `make test-extension-e2e`.

## Layout

| Path | What it is |
|---|---|
| `src/background.js` | Checks (cached for the session), the well-known sites list, stopping repeat visits to risky pages, the toolbar dot, the right-click menu |
| `src/lists.js` | Whether a host is one of the well-known sites (and not a host anyone can publish on) |
| `src/popup.js`, `src/blocked.js` | The popup, and the page shown in place of a stopped one |
| `src/content/` | The on-page card (closed shadow root, constructed stylesheets so a page's CSP can't block them) and the page script for warnings and the password guard |
| `src/render.js` | Verdict, score ring, flags and notes, built from text nodes only |
| `src/shared.js` | Settings, plus verdict wording, link extraction and notes imported from the website's source |
| `public/` | Manifest, HTML, CSS, icons and fonts (Geist, Geist Mono, Instrument Serif; OFL), copied into `dist/` |

The well-known sites come from the server: `GET /api/v1/lists/known-sites` (`server/internal/handler/lists.go`).
