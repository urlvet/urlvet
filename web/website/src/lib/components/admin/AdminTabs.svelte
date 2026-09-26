<script lang="ts">
  import type { Tab } from "../../admin/types";

  export let activeTab: Tab;
  export let errorCount = 0;
  export let reportCount = 0;
  export let onSelect: (tab: Tab) => void;

  const TABS: { id: Tab; label: string; short?: string }[] = [
    { id: "overview", label: "Overview" },
    { id: "recent", label: "Recent scans", short: "Scans" },
    { id: "reports", label: "Reports" },
    { id: "errors", label: "Errors" },
    { id: "cache", label: "Cache" },
  ];

  $: counts = { reports: reportCount, errors: errorCount } as Partial<Record<Tab, number>>;
</script>

<!-- Underline tabs: all five fit on a phone; scrolls sideways only as a last resort. -->
<nav
  aria-label="Admin sections"
  class="-mb-px flex justify-between sm:justify-start gap-1.5 min-[360px]:gap-2 sm:gap-6 overflow-x-auto no-scrollbar"
>
  {#each TABS as tab}
    {@const n = counts[tab.id] ?? 0}
    <button
      type="button"
      on:click={() => onSelect(tab.id)}
      aria-current={activeTab === tab.id ? "page" : undefined}
      class="flex-shrink-0 inline-flex items-center gap-1 sm:gap-1.5 py-3 border-b-2 text-[13px] sm:text-sm whitespace-nowrap transition-colors {activeTab ===
      tab.id
        ? 'border-gray-900 dark:border-gray-100 text-gray-900 dark:text-gray-100 font-medium'
        : 'border-transparent text-gray-500 hover:text-gray-900 dark:hover:text-gray-100'}"
    >
      {#if tab.short}<span class="sm:hidden">{tab.short}</span><span class="hidden sm:inline"
          >{tab.label}</span
        >{:else}{tab.label}{/if}
      {#if n > 0}
        <span
          class="min-w-[1.1rem] px-1 py-px rounded-full font-mono text-[10px] text-center {tab.id ===
          'errors'
            ? 'bg-red-100 dark:bg-red-500/20 text-red-600 dark:text-red-400'
            : 'bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-300'}">{n}</span
        >
      {/if}
    </button>
  {/each}
</nav>

<style>
  .no-scrollbar {
    scrollbar-width: none;
  }
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }
</style>
