<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { onDestroy, onMount, tick } from "svelte";
  import { fly, scale } from "svelte/transition";
  import { EXAMPLES } from "../../data/examples";
  import { ICON } from "../../ui/icons";
  import Icon from "../Icon.svelte";
  import CharacterClip from "./CharacterClip.svelte";
  import {
    GREETINGS,
    HIDE_BYE,
    INTRO,
    NAME,
    NUDGE,
    SCAN_FAILED,
    TIPS,
    TRY_EXAMPLE,
    explain,
    pick,
  } from "./messages";
  import type { Mood } from "./mood";
  import ScoreScale from "./ScoreScale.svelte";
  import { scanState, vettyMemory } from "./store";
  import { availableSteps, type TourStep } from "./tour";
  import TourOverlay from "./TourOverlay.svelte";

  // Vetty: waits in the corner, opens a speech bubble on click, runs the tour.
  type View = "menu" | "scores" | "explain" | "failed" | "examples" | "bye";
  type Item = { label: string; run: () => void };

  const NUDGE_AFTER_MS = 5000;

  let open = false;
  let view: View = "menu";
  let tourSteps: TourStep[] | null = null;
  let nudge = false;
  let greeting = "";
  let tip = "";
  let waving = false;
  let bubbleEl: HTMLDivElement;
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

  // Expression follows what's happening on the page.
  let mood: Mood = "idle";
  $: mood = tourSteps
    ? "point"
    : waving
      ? "wave"
      : scan.status === "scanning"
        ? "thinking"
        : scan.status === "error"
          ? "worried"
          : verdict === "Safe"
            ? "happy"
            : verdict === "Risky"
              ? "worried"
              : verdict === "Suspicious"
                ? "thinking"
                : "idle";

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
  $: items = ((): Item[] => {
    const hideItem: Item = { label: `Hide ${NAME}`, run: hide };
    const example: Item = { label: "Try an example", run: () => (view = "examples") };
    const howItWorks: Item = { label: "How does url.vet work?", run: () => go("/how-it-works") };
    const whoMade: Item = { label: "Who made this?", run: () => go("/about") };
    switch (context) {
      case "home":
        return [
          { label: "Show me around", run: startTour },
          example,
          howItWorks,
          whoMade,
          hideItem,
        ];
      case "scanning":
        return [hideItem];
      case "result":
        return [
          { label: "Explain this result", run: () => (view = "explain") },
          { label: "What does my score mean?", run: () => (view = "scores") },
          { label: "Show me around", run: startTour },
          { label: "Report a wrong verdict", run: () => clickTarget("report") },
          { label: "Scan another link", run: () => clickTarget("scan-another") },
          howItWorks,
          whoMade,
          hideItem,
        ];
      case "error":
        return [
          { label: "Why did my scan fail?", run: () => (view = "failed") },
          example,
          hideItem,
        ];
      default: {
        const nav: Item[] = [{ label: "Scan a link", run: () => go("/") }, example];
        if (context !== "how-it-works") nav.push(howItWorks);
        if (context !== "about") nav.push(whoMade);
        return [...nav, hideItem];
      }
    }
  })();

  $: explanation = result ? explain(result) : null;

  // ── open / close ───────────────────────────────────────────────────────────
  async function toggle() {
    nudge = false;
    if (open) return close();
    view = "menu";
    greeting = $vettyMemory.introduced ? pick(GREETINGS[context] ?? GREETINGS.home) : "";
    tip = pick(TIPS);
    open = true;
    waving = true;
    later(() => (waving = false), 1600);
    if (!$vettyMemory.introduced) vettyMemory.set({ introduced: true, nudged: true });
    await tick();
    bubbleEl?.querySelector<HTMLElement>("ul button")?.focus();
  }

  function close() {
    open = false;
    view = "menu";
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
    // One gentle nudge on a first visit, then never again.
    if (!$vettyMemory.nudged && !$vettyMemory.hidden) {
      later(() => {
        if (!open && !tourSteps && !$vettyMemory.hidden && !$vettyMemory.nudged) {
          nudge = true;
          vettyMemory.set({ nudged: true });
        }
      }, NUDGE_AFTER_MS);
    }
  });
  onDestroy(() => timers.forEach(clearTimeout));
</script>

<svelte:window on:keydown={onKey} />

{#if !$vettyMemory.hidden}
  <div
    class="fixed right-4 bottom-4 z-[70] flex flex-col items-end gap-2 pointer-events-none"
    style="margin-bottom: env(safe-area-inset-bottom, 0px)"
  >
    {#if open}
      <div
        bind:this={bubbleEl}
        role="dialog"
        aria-label="{NAME}, the url.vet helper"
        transition:fly={{ y: 12, duration: 200 }}
        class="bubble pointer-events-auto relative w-[min(330px,calc(100vw-32px))] max-h-[calc(100vh-140px)] overflow-y-auto rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shadow-2xl shadow-black/15 text-left"
      >
        <div class="flex items-start justify-between gap-3 px-5 pt-4">
          <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">{NAME}</p>
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
              <p class="font-serif text-xl leading-snug text-gray-900 dark:text-gray-100">
                {greeting}
              </p>
            {:else}
              <p class="font-serif text-xl leading-snug text-gray-900 dark:text-gray-100">
                {INTRO[0]}
              </p>
              <p class="mt-2">{INTRO[1]}</p>
            {/if}
            <ul
              class="mt-4 -mx-2 divide-y divide-gray-100 dark:divide-gray-800 border-y border-gray-100 dark:border-gray-800"
            >
              {#each items as item}
                <li>
                  <button
                    type="button"
                    class="group w-full flex items-center justify-between px-2 py-2.5 text-left text-[15px] text-gray-900 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800/50 rounded-lg focus:outline-none focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800"
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
            </ul>
            {#if context === "home" || context === "result"}
              <p class="mt-3 text-xs text-gray-500">{tip}</p>
            {/if}
          {:else if view === "scores" && result}
            <ScoreScale score={result.result?.final_score ?? 0} {verdict} />
          {:else if view === "failed"}
            {#each SCAN_FAILED as line, i}
              <p class={i === 0 ? "" : "mt-1"}>{line}</p>
            {/each}
          {:else if view === "explain" && explanation}
            <p class="font-serif text-xl leading-snug text-gray-900 dark:text-gray-100">
              {explanation.lead}
            </p>
            {#each explanation.lists as list}
              <p class="mt-3 font-mono text-[11px] uppercase tracking-wider text-gray-500">
                {list.title}
              </p>
              <ul class="mt-1.5 space-y-1">
                {#each list.items as item}
                  <li class="flex gap-2">
                    <span class="mt-[0.6em] w-1 h-1 rounded-full bg-gray-400 flex-shrink-0"></span>
                    <span>{item}</span>
                  </li>
                {/each}
              </ul>
            {/each}
            <p class="mt-3 text-gray-600 dark:text-gray-400">{explanation.advice}</p>
          {:else if view === "examples"}
            <p>{TRY_EXAMPLE}</p>
            <ul
              class="mt-3 -mx-2 divide-y divide-gray-100 dark:divide-gray-800 border-y border-gray-100 dark:border-gray-800"
            >
              {#each EXAMPLES as ex}
                <li>
                  <button
                    type="button"
                    class="w-full flex items-center gap-2.5 px-2 py-2.5 text-left rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 focus:outline-none focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800"
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
            <p class="font-serif text-xl leading-snug text-gray-900 dark:text-gray-100">
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
    {:else if nudge}
      <div
        transition:fly={{ y: 8, duration: 200 }}
        class="bubble pointer-events-auto relative max-w-[260px] rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shadow-xl shadow-black/10 px-4 py-3 text-sm text-gray-700 dark:text-gray-300"
      >
        <p>{NUDGE}</p>
        <div class="mt-2.5 flex gap-3 text-sm">
          <button
            type="button"
            class="font-medium text-gray-900 dark:text-gray-100"
            on:click={toggle}>Sure</button
          >
          <button type="button" class="text-gray-500" on:click={() => (nudge = false)}
            >No thanks</button
          >
        </div>
      </div>
    {/if}

    <button
      type="button"
      class="vetty pointer-events-auto rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400"
      aria-label={open ? `Close ${NAME}` : `Open ${NAME}, the url.vet helper`}
      aria-expanded={open}
      title="{NAME} (press ?)"
      on:click={toggle}
      in:scale={{ duration: 300, start: 0.6 }}
    >
      <CharacterClip {mood} size={68} />
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
