<script lang="ts">
  import { slide } from "svelte/transition";
  import type { Tone } from "../ui/tone";

  // One row in the grouped section list. The parent provides the card and dividers.
  export let id: string;
  export let title: string;
  /** SVG path (24x24 outline icon). */
  export let icon: string;
  export let status: { tone: Tone; label: string } | null = null;
  export let expanded = false;
  export let onToggle: () => void;
  /** Where this group of checks is explained, on /how-it-works. */
  export let learnMore: string | undefined = undefined;

  const DOT: Record<Tone, string> = {
    good: "bg-emerald-500",
    warn: "bg-yellow-500",
    bad: "bg-red-500",
    neutral: "bg-gray-400 dark:bg-gray-500",
  };
  const TEXT: Record<Tone, string> = {
    good: "text-gray-600 dark:text-gray-400",
    warn: "text-yellow-700 dark:text-yellow-400",
    bad: "text-red-600 dark:text-red-400",
    neutral: "text-gray-600 dark:text-gray-400",
  };

  $: tone = status?.tone ?? "neutral";
</script>

<div id="section-{id}" class="scroll-mt-20">
  <button
    type="button"
    class="group w-full flex items-center gap-3.5 px-5 py-4 text-left hover:bg-gray-100/70 dark:hover:bg-gray-800/40 transition-colors focus:outline-none focus-visible:bg-gray-100 dark:focus-visible:bg-gray-800/60"
    aria-expanded={expanded}
    aria-controls="section-{id}-body"
    on:click={onToggle}
  >
    <svg
      class="w-[18px] h-[18px] flex-shrink-0 text-gray-400 dark:text-gray-500 group-hover:text-gray-700 dark:group-hover:text-gray-300 transition-colors"
      fill="none"
      stroke="currentColor"
      stroke-width="1.75"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path stroke-linecap="round" stroke-linejoin="round" d={icon} />
    </svg>
    <span class="flex-shrink-0 text-[15px] font-medium text-gray-900 dark:text-gray-100"
      >{title}</span
    >
    <span class="flex-1 min-w-0 flex justify-end">
      {#if status}
        <span
          class="min-w-0 inline-flex items-center gap-2 font-mono text-xs {TEXT[tone]}"
          title={status.label}
        >
          <span class="w-1.5 h-1.5 rounded-full flex-shrink-0 {DOT[tone]}"></span>
          <span class="truncate">{status.label}</span>
        </span>
      {/if}
    </span>
    <svg
      class="w-4 h-4 flex-shrink-0 text-gray-400 dark:text-gray-600 transition-transform duration-200 {expanded
        ? 'rotate-180'
        : ''}"
      viewBox="0 0 20 20"
      fill="currentColor"
      aria-hidden="true"
    >
      <path
        fill-rule="evenodd"
        d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z"
        clip-rule="evenodd"
      />
    </svg>
  </button>
  {#if expanded}
    <div
      id="section-{id}-body"
      transition:slide={{ duration: 200 }}
      class="acc-body border-t border-gray-200 dark:border-gray-800 bg-gray-50/60 dark:bg-gray-950/40"
    >
      <slot />
      {#if learnMore}
        <p class="px-5 pb-4">
          <a
            href={learnMore}
            class="font-mono text-[11px] text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
            >How this check works →</a
          >
        </p>
      {/if}
    </div>
  {/if}
</div>

<style>
  /* Section components bring their own card chrome; the list is the card here. */
  .acc-body > :global(section),
  .acc-body > :global(div) {
    border: none !important;
    border-radius: 0 !important;
    box-shadow: none !important;
    background: transparent !important;
    margin: 0 !important;
  }
  .acc-body > :global(section):hover,
  .acc-body > :global(div):hover {
    transform: none !important;
    box-shadow: none !important;
  }
</style>
