<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { onDestroy, onMount, tick } from "svelte";
  import { fly, scale } from "svelte/transition";
  import { EXAMPLES } from "../../data/examples";
  import { easterEgg } from "../../results/eastereggs";
  import { realSite } from "../../results/realsite";
  import { PILL_OUTLINE, PILL_SOLID } from "../../ui/buttons";
  import { ICON } from "../../ui/icons";
  import { formatUrl, stripTrackers } from "../../utils";
  import Icon from "../Icon.svelte";
  import Bullets from "./Bullets.svelte";
  import CharacterClip from "./CharacterClip.svelte";
  import {
    CHECK_FOR_ME,
    CLEAN_LINK,
    GREETINGS,
    HIDE_BYE,
    INTRO,
    MEET,
    NAME,
    NUDGE,
    NUDGE_NO,
    NUDGE_YES,
    POKES,
    REAL_SITE,
    SCAN_FAILED,
    TIPS,
    TRY_EXAMPLE,
    WARN_INTRO,
    explain,
    pick,
    warnMessage,
  } from "./messages";
  import type { Mood } from "./mood";
  import ScoreScale from "./ScoreScale.svelte";
  import { nudgeDeclined, scanState, seenRecently, vettyMemory } from "./store";
  import { availableSteps, type TourStep } from "./tour";
  import TourOverlay from "./TourOverlay.svelte";

  // Vetty: waits in the corner, opens a speech bubble on click, runs the tour.
  type View =
    | "menu"
    | "scores"
    | "explain"
    | "failed"
    | "examples"
    | "check"
    | "meet"
    | "real"
    | "warn"
    | "clean"
    | "bye";
  type Item = { label: string; run: () => void };

  const NUDGE_AFTER_MS = 5000;

  // Styles shared by the speech bubbles and what's inside them.
  const BUBBLE =
    "rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900";
  const LEAD = "font-serif text-xl leading-snug text-gray-900 dark:text-gray-100";
  const LABEL = "font-mono text-[11px] uppercase tracking-wider text-gray-500";
  const FIELD =
    "w-full resize-none rounded-xl border border-gray-300 dark:border-gray-700 bg-transparent px-3 py-2";

  let open = false;
  let view: View = "menu";
  let tourSteps: TourStep[] | null = null;
  let nudge = false;
  /** An easter-egg line about the link just scanned (see results/eastereggs.ts). */
  let quip: string | null = null;
  let quipFor = "";
  let greeting = "";
  let tip = "";
  let waving = false;
  let bubbleEl: HTMLDivElement;
  let checkInput = "";
  /** Back from a view opened under "More" returns to "More". */
  const MORE_ROW =
    "w-full flex items-center justify-between px-2 py-2.5 text-left text-[15px] text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800/50 rounded-lg focus:outline-hidden focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800";
  /** Whether "More" is open below the first items. */
  let showMore = false;
  let lessRow: HTMLLIElement;

  // Opens the rest of the menu below the first items. Where the bubble can't
  // fit it all (phones), scroll so the new items come into view.
  async function expandMore() {
    showMore = true;
    await tick();
    if (bubbleEl && lessRow && bubbleEl.scrollHeight > bubbleEl.clientHeight) {
      bubbleEl.scrollTo({ top: lessRow.offsetTop - 8, behavior: "smooth" });
    }
  }
  let checkEl: HTMLTextAreaElement;
  /** The warning message being edited, and the link in the search bar minus its trackers. */
  let warnText = "";
  let clean: { cleaned: string; removed: string[] } | null = null;
  let copied = "";
  const timers: ReturnType<typeof setTimeout>[] = [];
  const later = (fn: () => void, ms: number) => timers.push(setTimeout(fn, ms));

  // Where are we? The scan only matters on the home page.
  $: path = $page.url.pathname;
  $: onHome = path === "/";
  $: scan = onHome ? $scanState : ({ status: "idle" } as const);
  $: result = scan.status === "done" ? scan.result : null;
  $: verdict = result?.result?.verdict;

  /** Which situation we're in, for greetings and the menu. */
  $: context = !onHome
    ? path.startsWith("/how-it-works")
      ? "how-it-works"
      : path.startsWith("/about")
        ? "about"
        : "other"
    : scan.status === "done"
      ? "result"
      : scan.status === "scanning"
        ? "scanning"
        : scan.status === "error"
          ? "error"
          : "home";

  /** A short-lived reaction (a paste, a poke) that wins over the page's mood. */
  let reaction: Mood | null = null;
  function react(m: Mood, ms = 1200) {
    reaction = m;
    later(() => reaction === m && (reaction = null), ms);
  }

  // Expression follows what's happening on the page.
  let mood: Mood = "idle";
  $: mood = tourSteps
    ? "point"
    : waving
      ? "wave"
      : reaction
        ? reaction
        : scan.status === "scanning"
          ? "scanning"
          : scan.status === "error"
            ? "worried"
            : verdict === "Safe"
              ? "happy"
              : verdict === "Risky"
                ? "alarmed"
                : verdict === "Suspicious"
                  ? "suspicious"
                  : "idle";

  // ── eyes follow the pointer, and the search bar while typing ───────────────
  let vettyEl: HTMLButtonElement;
  let gaze: [number, number] | null = null;
  let gazeTimer: ReturnType<typeof setTimeout> | undefined;
  let gazeFrame = 0;
  const GAZE_REACH = 2; // how far the pupils move, in drawing units

  function lookAt(x: number, y: number) {
    if (!vettyEl) return;
    const box = vettyEl.getBoundingClientRect();
    // The eyes sit a little above the middle of the drawing.
    const dx = x - (box.left + box.width * 0.5);
    const dy = y - (box.top + box.height * 0.33);
    const dist = Math.hypot(dx, dy);
    gaze = dist < 24 ? [0, 0.4] : [(dx / dist) * GAZE_REACH, (dy / dist) * GAZE_REACH];
    // Drift back to the idle glances once nothing has moved for a while.
    clearTimeout(gazeTimer);
    gazeTimer = setTimeout(() => (gaze = null), 2500);
  }

  function onPointerMove(e: PointerEvent) {
    if (gazeFrame) return;
    gazeFrame = requestAnimationFrame(() => {
      gazeFrame = 0;
      lookAt(e.clientX, e.clientY);
    });
  }

  function lookAtSearch() {
    const bar = document.querySelector<HTMLElement>('[data-guide="search"]');
    if (!bar) return;
    const box = bar.getBoundingClientRect();
    lookAt(box.left + box.width / 2, box.top + box.height / 2);
  }

  // Typing in the search bar draws his eyes to it; a pasted link gets an "ooh".
  function onInput(e: Event) {
    if ((e.target as HTMLElement | null)?.id !== "url-input") return;
    lookAtSearch();
    if ((e as InputEvent).inputType === "insertFromPaste") react("ooh");
  }
  function onClick(e: MouseEvent) {
    if ((e.target as HTMLElement | null)?.closest('[data-guide="paste"]')) {
      lookAtSearch();
      react("ooh");
    }
  }

  // ── phones: a bit smaller ──────────────────────────────────────────────────
  let small = false;

  // ── pokes: tapping him over and over gets escalating complaints ────────────
  let pokeTimes: number[] = [];
  let pokeStep = 0;
  let pokeReset: ReturnType<typeof setTimeout> | undefined;
  function onVettyClick() {
    const now = Date.now();
    pokeTimes = [...pokeTimes.filter((t) => now - t < 1500), now];
    clearTimeout(pokeReset);
    pokeReset = setTimeout(() => (pokeStep = 0), 5000);
    if (pokeTimes.length < 3) return toggle();
    close();
    nudge = false;
    const line = POKES[Math.min(pokeStep, POKES.length - 1)];
    pokeStep++;
    quip = line;
    react(pokeStep > 2 ? "worried" : "ooh", 1500);
    later(() => quip === line && (quip = null), 4000);
  }

  // ── actions ────────────────────────────────────────────────────────────────
  function startTour() {
    close();
    const steps = availableSteps();
    if (steps.length) tourSteps = steps;
  }

  function clickTarget(target: string) {
    close();
    document.querySelector<HTMLElement>(`[data-guide="${target}"]`)?.click();
  }

  function tryExample(url: string) {
    close();
    const chip = document.querySelector<HTMLElement>(
      `[data-guide="examples"] [data-example="${url}"]`
    );
    if (chip)
      chip.click(); // landing page: same as clicking the chip
    else if (onHome)
      location.href = `/?q=${encodeURIComponent(url)}`; // home without chips
    else goto(`/?q=${encodeURIComponent(url)}`); // other pages: home scans ?q= on load
  }

  // Puts the link in the real search bar and submits it there, so the visitor sees
  // where links go and the page's own handling (e.g. finding links in a message) applies.
  function checkForMe() {
    const text = checkInput.trim();
    if (!text) return;
    checkInput = "";
    putInSearchBar(text);
  }

  function putInSearchBar(text: string) {
    close();
    const field = document.getElementById("url-input") as HTMLInputElement | null;
    if (!field?.form) {
      goto(`/?q=${encodeURIComponent(text)}`);
      return;
    }
    field.scrollIntoView({ block: "center", behavior: "smooth" });
    field.value = text;
    field.dispatchEvent(new Event("input", { bubbles: true }));
    field.focus({ preventScroll: true });
    later(() => field.form?.requestSubmit(), 400);
  }

  $: if (view === "check" && checkEl) tick().then(() => checkEl?.focus());

  $: real = result ? realSite(result) : null;

  function openWarn() {
    if (!result?.result) return;
    warnText = warnMessage(
      result.result.verdict,
      result.domain,
      result.result.final_score,
      window.location.href
    );
    view = "warn";
  }

  async function copy(text: string, what: string) {
    try {
      await navigator.clipboard.writeText(text);
      copied = what;
      later(() => copied === what && (copied = ""), 1500);
    } catch {
      /* clipboard blocked */
    }
  }

  const canShare = typeof navigator !== "undefined" && "share" in navigator;

  function shareNative() {
    navigator.share?.({ text: warnText }).catch(() => {});
  }

  function go(to: string) {
    close();
    goto(to);
  }

  function hide() {
    view = "bye";
    later(() => {
      close();
      vettyMemory.set({ hidden: true });
    }, 2200);
  }

  // ── menu: only what works on this page ─────────────────────────────────────
  // A few choices for what people most likely want right now; everything else
  // is one tap away under "More", so the first list never gets overwhelming.
  type Menu = { main: Item[]; more: Item[] };
  $: menu = ((): Menu => {
    const hideItem: Item = { label: `Hide ${NAME}`, run: hide };
    const meetItem: Item = { label: `Who is ${NAME}?`, run: () => (view = "meet") };
    const example: Item = { label: "Try it with a sample link", run: () => (view = "examples") };
    const howItWorks: Item = {
      label: "How does url.vet check a link?",
      run: () => go("/how-it-works"),
    };
    const whoMade: Item = { label: "Who made this website?", run: () => go("/about") };
    const cleanItem = (label: string): Item[] =>
      clean ? [{ label, run: () => (view = "clean") }] : [];
    switch (context) {
      case "home":
        return {
          // In the order a visitor goes: find their way round, try it, then
          // check their own link (cleaning it first, when it has trackers).
          main: [
            { label: "Show me around", run: startTour },
            example,
            ...cleanItem("Remove trackers from my link"),
            { label: "Check a link for me", run: () => (view = "check") },
          ],
          more: [howItWorks, whoMade, meetItem, hideItem],
        };
      case "scanning":
        return { main: [meetItem, hideItem], more: [] };
      case "result":
        return {
          // In the order a visitor goes: find their way round the result,
          // understand it, act on it, or say it's wrong. The rest is under More.
          main: [
            { label: "Show me around", run: startTour },
            { label: "Explain this result simply", run: () => (view = "explain") },
            ...(real
              ? [{ label: `Take me to the real ${real.domain}`, run: () => (view = "real") }]
              : []),
            { label: "Let the sender know", run: openWarn },
            { label: "I think this result is wrong", run: () => clickTarget("report") },
          ],
          more: [
            { label: "Check a different link", run: () => clickTarget("scan-another") },
            ...cleanItem("Copy the link without trackers"),
            howItWorks,
            whoMade,
            meetItem,
            hideItem,
          ],
        };
      case "error":
        return {
          main: [{ label: "Why didn't the check work?", run: () => (view = "failed") }, example],
          more: [meetItem, hideItem],
        };
      default: {
        const nav: Item[] = [];
        if (context !== "how-it-works") nav.push(howItWorks);
        if (context !== "about") nav.push(whoMade);
        return {
          main: [{ label: "I want to check a link", run: () => go("/") }, example],
          more: [...nav, meetItem, hideItem],
        };
      }
    }
  })();

  $: explanation = result ? explain(result) : null;

  function declineNudge() {
    nudge = false;
    vettyMemory.set({ nudgeDeclinedAt: Date.now() });
  }

  // ── open / close ───────────────────────────────────────────────────────────
  async function toggle() {
    nudge = false;
    quip = null;
    if (open) return close();
    view = "menu";
    showMore = false;
    // Trackers in whatever is in the search bar (the scanned link, on a result).
    const typed = (document.getElementById("url-input") as HTMLInputElement | null)?.value.trim();
    const stripped = typed ? stripTrackers(formatUrl(typed)) : null;
    clean = stripped?.removed.length ? stripped : null;
    greeting = seenRecently($vettyMemory) ? pick(GREETINGS[context] ?? GREETINGS.home) : "";
    tip = pick(TIPS);
    open = true;
    waving = true;
    later(() => (waving = false), 1600);
    vettyMemory.set({ openedAt: Date.now() });
    await tick();
    bubbleEl?.querySelector<HTMLElement>("ul button")?.focus();
  }

  function close() {
    open = false;
    view = "menu";
  }

  // Once per result, special links get a quip. Not Risky ones: a joke, or any
  // pop-up, is the wrong note there (the demo fakes still get their joke). If
  // Vetty is hidden, the page shows the quip instead.
  $: if (!result) {
    quip = null;
  }
  $: if (result && result.url !== quipFor && !$vettyMemory.hidden) {
    quipFor = result.url;
    const line = easterEgg(result.url);
    // The fake example chips are demos, not real threats: they get their joke.
    const demo = EXAMPLES.some((e) => e.hint === "Fake" && e.url === result.domain);
    if (line && (result.result?.verdict !== "Risky" || demo)) {
      later(() => {
        if (open || tourSteps || $vettyMemory.hidden) return;
        quip = line;
        later(() => quip === line && (quip = null), 8000);
      }, 1100);
    }
  }

  // Changing page closes the bubble.
  $: (path, close());

  function onKey(e: KeyboardEvent) {
    const typing =
      e.target instanceof HTMLInputElement ||
      e.target instanceof HTMLTextAreaElement ||
      (e.target as HTMLElement)?.isContentEditable;
    if (e.key === "?" && !typing && !tourSteps) {
      e.preventDefault();
      toggle();
    } else if (e.key === "Escape" && open) {
      close();
    }
  }

  onMount(() => {
    const narrow = window.matchMedia("(max-width: 639px)");
    small = narrow.matches;
    const onNarrow = () => (small = narrow.matches);
    narrow.addEventListener("change", onNarrow);
    document.addEventListener("input", onInput, true);
    document.addEventListener("click", onClick, true);

    // Until Vetty has been opened (and again a day later), a nudge waits beside
    // him, with two replies: open him, or "No thanks" (no nudge for a month).
    const nudgeWanted = () =>
      !$vettyMemory.hidden && !seenRecently($vettyMemory) && !nudgeDeclined($vettyMemory);
    if (nudgeWanted()) {
      later(() => {
        if (!open && !tourSteps && nudgeWanted()) nudge = true;
      }, NUDGE_AFTER_MS);
    }

    return () => {
      narrow.removeEventListener("change", onNarrow);
      document.removeEventListener("input", onInput, true);
      document.removeEventListener("click", onClick, true);
    };
  });
  onDestroy(() => {
    timers.forEach(clearTimeout);
    clearTimeout(gazeTimer);
    clearTimeout(pokeReset);
    if (gazeFrame) cancelAnimationFrame(gazeFrame);
  });
</script>

<svelte:window on:keydown={onKey} on:pointermove={onPointerMove} />

{#if !$vettyMemory.hidden}
  <div
    class="fixed right-3 sm:right-4 bottom-3 sm:bottom-4 z-[70] flex flex-col items-end gap-2 pointer-events-none"
    style="margin-bottom: env(safe-area-inset-bottom, 0px)"
  >
    {#if open}
      <div
        bind:this={bubbleEl}
        role="dialog"
        aria-label="{NAME}, the url.vet helper"
        transition:fly={{ y: 12, duration: 200 }}
        class="bubble pointer-events-auto relative w-[min(330px,calc(100vw-32px))] max-h-[calc(100vh-140px)] overflow-y-auto {BUBBLE} shadow-2xl shadow-black/15 text-left"
      >
        <div class="flex items-start justify-between gap-3 px-5 pt-4">
          <p class={LABEL}>{NAME}</p>
          <button
            type="button"
            class="-mr-2 -mt-1 p-1.5 rounded-md text-gray-400 hover:text-gray-900 dark:hover:text-gray-100"
            aria-label="Close {NAME}"
            on:click={close}
          >
            <Icon path={ICON.close} class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="px-5 pb-5 pt-1 text-[15px] leading-relaxed text-gray-700 dark:text-gray-300">
          {#if view === "menu"}
            {#if greeting}
              <p class={LEAD}>
                {greeting}
              </p>
            {:else}
              <p class={LEAD}>
                {INTRO[0]}
              </p>
              <p class="mt-2">{INTRO[1]}</p>
            {/if}
            <ul
              class="mt-4 -mx-2 divide-y divide-gray-100 dark:divide-gray-800 border-y border-gray-100 dark:border-gray-800"
            >
              {#each showMore ? [...menu.main, ...menu.more] : menu.main as item, i}
                {#if i === menu.main.length}
                  <!-- "More" stays where it was, now as "Less", with the rest below it. -->
                  <li bind:this={lessRow}>
                    <button
                      type="button"
                      class={MORE_ROW}
                      aria-expanded="true"
                      on:click={() => (showMore = false)}
                    >
                      Less
                      <Icon path={ICON.chevronDown} class="w-3.5 h-3.5 text-gray-400 rotate-180" />
                    </button>
                  </li>
                {/if}

                <li>
                  <button
                    type="button"
                    class="group w-full flex items-center justify-between px-2 py-2.5 text-left text-[15px] text-gray-900 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800/50 rounded-lg focus:outline-hidden focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800"
                    on:click={item.run}
                  >
                    {item.label}
                    <Icon
                      path={ICON.arrowRight}
                      class="w-3.5 h-3.5 text-gray-400 group-hover:translate-x-0.5 transition-transform"
                    />
                  </button>
                </li>
              {/each}
              {#if !showMore && menu.more.length}
                <li>
                  <button
                    type="button"
                    class={MORE_ROW}
                    aria-expanded="false"
                    on:click={expandMore}
                  >
                    More
                    <Icon path={ICON.chevronDown} class="w-3.5 h-3.5 text-gray-400" />
                  </button>
                </li>
              {/if}
            </ul>
            {#if !showMore && (context === "home" || context === "result")}
              <p class="mt-3 text-xs text-gray-500">{tip}</p>
            {/if}
          {:else if view === "scores" && result}
            <ScoreScale score={result.result?.final_score ?? 0} {verdict} />
          {:else if view === "failed"}
            {#each SCAN_FAILED as line, i}
              <p class={i === 0 ? "" : "mt-1"}>{line}</p>
            {/each}
          {:else if view === "explain" && explanation}
            <p class={LEAD}>
              {explanation.lead}
            </p>
            {#each explanation.lists as list}
              <p class="mt-3 {LABEL}">
                {list.title}
              </p>
              <Bullets items={list.items} />
            {/each}
            <p class="mt-3 text-gray-600 dark:text-gray-400">{explanation.advice}</p>
            <button
              type="button"
              class="mt-3 block text-sm text-gray-900 dark:text-gray-100 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:decoration-current"
              on:click={() => (view = "scores")}>What does the score mean?</button
            >
          {:else if view === "meet"}
            <p class={LEAD}>
              {MEET.lead}
            </p>
            <p class="mt-2">{MEET.why}</p>
            <p class="mt-4 {LABEL}">
              {MEET.canTitle}
            </p>
            <Bullets items={MEET.can} />
            <p class="mt-4 text-sm text-gray-500">{MEET.note}</p>
          {:else if view === "check"}
            <form on:submit|preventDefault={checkForMe}>
              <label for="vetty-check" class="block">{CHECK_FOR_ME}</label>
              <textarea
                id="vetty-check"
                bind:this={checkEl}
                bind:value={checkInput}
                rows="3"
                placeholder="https://…"
                class="mt-3 {FIELD} font-mono text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-hidden focus:border-gray-500 dark:focus:border-gray-500"
                on:keydown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    checkForMe();
                  }
                }}
              ></textarea>
              <button
                type="submit"
                disabled={!checkInput.trim()}
                class="mt-2 {PILL_SOLID} px-3.5 py-2"
              >
                Check it
                <Icon path={ICON.arrowRight} class="w-3.5 h-3.5" />
              </button>
            </form>
          {:else if view === "real" && real}
            <p class={LEAD}>
              {REAL_SITE[real.reason](real.domain)}
            </p>
            <p class="mt-2">{REAL_SITE.advice(real.domain)}</p>
            <div class="mt-4 flex flex-wrap gap-2">
              <a
                href="https://{real.domain}"
                target="_blank"
                rel="noopener noreferrer"
                class="{PILL_SOLID} px-3.5 py-2">Open {real.domain} ↗</a
              >
              <button
                type="button"
                class="{PILL_OUTLINE} px-3.5 py-2"
                on:click={() => real && putInSearchBar(real.domain)}>Check it first</button
              >
            </div>
          {:else if view === "warn"}
            <label for="vetty-warn" class="block">{WARN_INTRO[verdict ?? "Risky"]}</label>
            <textarea
              id="vetty-warn"
              bind:value={warnText}
              rows="6"
              class="mt-3 {FIELD} text-sm leading-relaxed text-gray-900 dark:text-gray-100 focus:outline-hidden focus:border-gray-500"
            ></textarea>
            <div class="mt-2 flex flex-wrap gap-2">
              <a
                href="https://wa.me/?text={encodeURIComponent(warnText)}"
                target="_blank"
                rel="noopener noreferrer"
                class="{PILL_SOLID} px-3.5 py-2">Send on WhatsApp</a
              >
              <button
                type="button"
                class="{PILL_OUTLINE} px-3.5 py-2"
                on:click={() => copy(warnText, "warn")}
                >{copied === "warn" ? "Copied" : "Copy"}</button
              >
              {#if canShare}
                <button type="button" class="{PILL_OUTLINE} px-3.5 py-2" on:click={shareNative}
                  >Other apps…</button
                >
              {/if}
            </div>
          {:else if view === "clean" && clean}
            <p>{CLEAN_LINK.lead}</p>
            <p class="mt-3 text-xs text-gray-500">
              Removed: <span class="font-mono">{clean.removed.join(", ")}</span>
            </p>
            <p class="mt-3">{CLEAN_LINK.done}</p>
            <p
              class="mt-1.5 rounded-xl border border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/40 px-3 py-2 font-mono text-[13px] break-all text-gray-900 dark:text-gray-100"
            >
              {clean.cleaned}
            </p>
            <div>
              <button
                type="button"
                class="mt-3 {PILL_SOLID} px-3.5 py-2"
                on:click={() => clean && copy(clean.cleaned, "clean")}
                >{copied === "clean" ? "Copied" : "Copy clean link"}</button
              >
            </div>
          {:else if view === "examples"}
            <p>{TRY_EXAMPLE}</p>
            <ul
              class="mt-3 -mx-2 divide-y divide-gray-100 dark:divide-gray-800 border-y border-gray-100 dark:border-gray-800"
            >
              {#each EXAMPLES as ex}
                <li>
                  <button
                    type="button"
                    class="w-full flex items-center gap-2.5 px-2 py-2.5 text-left rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 focus:outline-hidden focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800"
                    on:click={() => tryExample(ex.url)}
                  >
                    <span
                      class="w-1.5 h-1.5 rounded-full {ex.hint === 'Legit'
                        ? 'bg-emerald-500'
                        : 'bg-red-500'}"
                    ></span>
                    <span class="font-mono text-sm text-gray-900 dark:text-gray-100"
                      >{ex.label}</span
                    >
                    <span class="ml-auto font-mono text-[11px] text-gray-500">{ex.hint}</span>
                  </button>
                </li>
              {/each}
            </ul>
          {:else if view === "bye"}
            <p class={LEAD}>
              {HIDE_BYE}
            </p>
          {/if}

          {#if view !== "menu" && view !== "bye"}
            <button
              type="button"
              class="mt-4 inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-gray-900 dark:hover:text-gray-100"
              on:click={() => (view = "menu")}
            >
              ← Back
            </button>
          {/if}
        </div>
      </div>
    {:else if quip}
      <!-- An aside, not a question: no buttons. Tap it (or wait) to make it go away. -->
      <div role="status" transition:fly={{ y: 8, duration: 200 }} class="pointer-events-auto">
        <button
          type="button"
          class="bubble relative max-w-[280px] {BUBBLE} shadow-xl shadow-black/10 px-4 py-3 text-left text-sm text-gray-700 dark:text-gray-300"
          title="Dismiss"
          on:click={() => (quip = null)}
        >
          {quip}
        </button>
      </div>
    {:else if nudge}
      <!-- Stays until answered: open Vetty, or "No thanks". -->
      <div
        role="status"
        transition:fly={{ y: 8, duration: 200 }}
        class="bubble pointer-events-auto relative max-w-[260px] {BUBBLE} shadow-xl shadow-black/10 px-4 py-3 text-sm text-gray-700 dark:text-gray-300"
      >
        <p>{NUDGE}</p>
        <div class="mt-3 flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-full bg-gray-900 dark:bg-gray-100 px-3 py-1.5 text-xs font-medium text-white dark:text-gray-900 hover:opacity-90 transition-opacity"
            on:click={toggle}>{NUDGE_YES}</button
          >
          <button
            type="button"
            class="rounded-full border border-gray-300 dark:border-gray-600 px-3 py-1.5 text-xs font-medium text-gray-700 dark:text-gray-300 hover:border-gray-500 dark:hover:border-gray-400 transition-colors"
            on:click={declineNudge}>{NUDGE_NO}</button
          >
        </div>
      </div>
    {/if}

    <button
      type="button"
      data-guide="vetty"
      class="vetty pointer-events-auto rounded-full focus:outline-hidden focus-visible:ring-2 focus-visible:ring-gray-400"
      aria-label={open ? `Close ${NAME}` : `Open ${NAME}, the url.vet helper`}
      aria-expanded={open}
      title="{NAME} (press ?)"
      bind:this={vettyEl}
      on:click={onVettyClick}
      in:scale={{ duration: 300, start: 0.6 }}
    >
      <CharacterClip {mood} {gaze} size={small ? 56 : 68} />
    </button>
  </div>
{/if}

{#if tourSteps}
  <TourOverlay steps={tourSteps} onClose={() => (tourSteps = null)} />
{/if}

<style>
  .vetty {
    transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
    filter: drop-shadow(0 6px 10px rgba(0, 0, 0, 0.12));
  }
  .vetty:hover {
    transform: translateY(-3px) rotate(-3deg);
  }
  /* little tail pointing down at Vetty */
  .bubble::after {
    content: "";
    position: absolute;
    right: 26px;
    bottom: -7px;
    width: 12px;
    height: 12px;
    transform: rotate(45deg);
    background: inherit;
    border-right: inherit;
    border-bottom: inherit;
  }
</style>
