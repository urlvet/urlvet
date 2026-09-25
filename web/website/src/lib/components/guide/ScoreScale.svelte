<script lang="ts">
  import { SCORE_BANDS, SCORES_FOOTNOTE } from "./messages";

  // "What does my score mean?": the 0–100 scale with this result marked on it.
  export let score: number;
  export let verdict: string | undefined;

  $: clamped = Math.max(0, Math.min(100, score));
  $: current = SCORE_BANDS.find((b) => b.verdict === verdict);
</script>

<p class="font-serif text-xl leading-snug text-gray-900 dark:text-gray-100">
  This one scored <span class={current?.text}>{score}</span> out of 100.
</p>

<!-- scale -->
<div class="mt-6 mb-1">
  <div class="relative">
    <div class="flex h-2 rounded-full overflow-hidden gap-0.5">
      {#each SCORE_BANDS as band}
        <div
          class="{band.dot} {band.verdict === verdict ? '' : 'opacity-30'}"
          style="width: {band.to - band.from + 1}%"
        ></div>
      {/each}
    </div>
    <!-- marker -->
    <div
      class="absolute -top-5 flex flex-col items-center"
      style="left: {clamped}%; transform: translateX(-50%)"
    >
      <span class="font-mono text-[10px] text-gray-500 leading-none">{score}</span>
      <span class="mt-1 w-0.5 h-5 rounded-full bg-gray-900 dark:bg-gray-100"></span>
    </div>
  </div>
  <div class="relative mt-1.5 h-3 font-mono text-[10px] text-gray-400">
    {#each [0, 30, 65, 100] as tick}
      <span
        class="absolute"
        style="left: {tick}%; transform: translateX({tick === 0
          ? '0'
          : tick === 100
            ? '-100%'
            : '-50%'})">{tick}</span
      >
    {/each}
  </div>
</div>

<!-- bands -->
<ul class="mt-3 space-y-2">
  {#each [...SCORE_BANDS].reverse() as band}
    <li
      class="rounded-xl border px-3 py-2 {band.verdict === verdict
        ? 'border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50'
        : 'border-transparent'}"
    >
      <p class="flex items-center gap-2 text-sm font-medium text-gray-900 dark:text-gray-100">
        <span class="w-1.5 h-1.5 rounded-full {band.dot}"></span>
        {band.verdict}
        <span class="font-mono text-[11px] font-normal text-gray-500">{band.from}–{band.to}</span>
        {#if band.verdict === verdict}
          <span class="ml-auto font-mono text-[10px] uppercase tracking-wider text-gray-500"
            >this link</span
          >
        {/if}
      </p>
      <p class="mt-0.5 text-[13px] leading-snug text-gray-600 dark:text-gray-400">{band.meaning}</p>
    </li>
  {/each}
</ul>

<p class="mt-3 text-xs text-gray-500">{SCORES_FOOTNOTE}</p>
