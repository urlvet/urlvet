<script lang="ts">
  import { onDestroy, onMount, tick } from "svelte";
  import { fade } from "svelte/transition";
  import { PILL_OUTLINE, PILL_SOLID } from "../../ui/buttons";
  import type { TourStep } from "./tour";

  // Dims the page around one element at a time and shows a step card beside it.
  export let steps: TourStep[];
  export let onClose: () => void;

  let index = 0;
  let rect: DOMRect | null = null;
  let cardHeight = 0;
  let cardEl: HTMLDivElement;

  const DEFAULT_PAD = 8;
  const CARD_W = 320;
  const GAP = 14;

  $: step = steps[index];
  $: PAD = step?.pad ?? DEFAULT_PAD;
  $: last = index === steps.length - 1;

  function targetEl(): HTMLElement | null {
    return step ? document.querySelector(`[data-guide="${step.target}"]`) : null;
  }

  function measure() {
    rect = targetEl()?.getBoundingClientRect() ?? null;
  }

  async function show(i: number) {
    index = i;
    await tick();
    const el = targetEl();
    if (!el) return;
    el.scrollIntoView({ block: "center", behavior: "smooth" });
    // Re-measure while the smooth scroll settles.
    for (const ms of [0, 150, 350, 600]) setTimeout(measure, ms);
    cardEl?.focus();
  }

  function next() {
    if (last) onClose();
    else show(index + 1);
  }
  function back() {
    if (index > 0) show(index - 1);
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === "Escape") onClose();
    else if (e.key === "ArrowRight") next();
    else if (e.key === "ArrowLeft") back();
  }

  // Card goes below the target when there's room, otherwise above; kept on screen.
  $: card = (() => {
    if (!rect) return { top: 24, left: 24 };
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    const below = rect.bottom + PAD + GAP;
    const top =
      below + cardHeight < vh - 16 ? below : Math.max(16, rect.top - PAD - GAP - cardHeight);
    const left = Math.min(Math.max(16, rect.left), vw - Math.min(CARD_W, vw - 32) - 16);
    return { top, left };
  })();

  onMount(() => {
    show(0);
    window.addEventListener("resize", measure);
    window.addEventListener("scroll", measure, true);
  });
  onDestroy(() => {
    window.removeEventListener("resize", measure);
    window.removeEventListener("scroll", measure, true);
  });
</script>

<svelte:window on:keydown={onKey} />

<div class="fixed inset-0 z-[60]" transition:fade={{ duration: 150 }}>
  <!-- click-catcher: clicking the dimmed page ends the tour -->
  <button
    class="absolute inset-0 w-full h-full cursor-default"
    aria-label="End tour"
    on:click={onClose}
  ></button>

  {#if rect}
    <div
      class="spotlight absolute rounded-2xl pointer-events-none"
      style="top:{rect.top - PAD}px;left:{rect.left - PAD}px;width:{rect.width +
        PAD * 2}px;height:{rect.height + PAD * 2}px"
    ></div>
  {/if}

  <div
    bind:this={cardEl}
    bind:clientHeight={cardHeight}
    tabindex="-1"
    role="dialog"
    aria-label="Tour step {index + 1} of {steps.length}"
    class="absolute rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-2xl shadow-black/20 dark:shadow-black/60 p-5 focus:outline-hidden"
    style="top:{card.top}px;left:{card.left}px;width:min({CARD_W}px, calc(100vw - 32px));transition:top .25s, left .25s"
  >
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">
      Step {index + 1} of {steps.length}
    </p>
    <p class="mt-2 font-serif text-2xl leading-tight text-gray-900 dark:text-gray-100">
      {step?.title}
    </p>
    <p class="mt-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">{step?.text}</p>
    <div class="mt-5 flex items-center justify-between gap-2">
      <button
        type="button"
        class="text-sm text-gray-500 hover:text-gray-900 dark:hover:text-gray-100"
        on:click={onClose}>Skip</button
      >
      <div class="flex gap-2">
        {#if index > 0}
          <button
            type="button"
            class="{PILL_OUTLINE} dark:!border-gray-600 px-4 py-2"
            on:click={back}>Back</button
          >
        {/if}
        <button type="button" class="{PILL_SOLID} px-4 py-2" on:click={next}>
          {last ? "Done" : "Next"}
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .spotlight {
    box-shadow: 0 0 0 9999px rgba(15, 13, 11, 0.55);
    transition:
      top 0.25s,
      left 0.25s,
      width 0.25s,
      height 0.25s;
  }
  /* On a dark page a light dim barely shows, so dim harder and outline the
     spotlit area so it reads as lifted out. */
  :global(.dark) .spotlight {
    box-shadow:
      0 0 0 1.5px rgba(255, 255, 255, 0.22),
      0 0 0 9999px rgba(0, 0, 0, 0.72);
  }
</style>
