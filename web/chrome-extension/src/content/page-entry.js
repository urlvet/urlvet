// Runs on every page while automatic warnings or the password guard are on.
// The background worker decides what to check: well-known and trusted sites
// never are, and known-bad pages are stopped before they get here.
import { installCard } from "./card.js";

if (window.top === window) {
  const card = installCard();
  const ask = (type) =>
    new Promise((resolve) =>
      chrome.runtime.sendMessage({ type, url: location.href }, (result) =>
        resolve(chrome.runtime.lastError ? null : result),
      ),
    );

  // Automatic warning as the page opens.
  ask("auto").then((result) => {
    if (result) card.show({ state: "result", mode: "auto", ...result });
  });

  // The password guard: once per page, the first time a password field gets focus.
  let guarded = false;
  document.addEventListener(
    "focusin",
    async (e) => {
      const field = e.target;
      if (
        guarded ||
        !(field instanceof HTMLInputElement) ||
        field.type !== "password"
      )
        return;
      guarded = true;
      const result = await ask("guard");
      // Already warning about this page: no need for a second card.
      if (result && card.current()?.mode !== "auto")
        card.show({ state: "result", mode: "guard", ...result });
    },
    true,
  );
}
