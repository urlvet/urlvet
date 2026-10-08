<script lang="ts">
  import { onMount } from "svelte";
  import type { AdminClient } from "../../admin/client";
  import { verdictTextColor } from "../../admin/format";
  import type { Guard, Stats } from "../../admin/types";

  export let client: AdminClient;
  export let guard: Guard;

  let stats: Stats | null = null;
  let loading = true;

  onMount(async () => {
    stats = (await guard(client.stats())) ?? null;
    loading = false;
  });
</script>

{#if loading}
  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
    {#each Array(4) as _}
      <div
        class="h-28 bg-white dark:bg-gray-900 rounded-2xl border border-gray-200 dark:border-gray-800 animate-pulse"
      ></div>
    {/each}
  </div>
{:else if stats}
  <!-- Stat cards -->
  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <p class="font-mono text-[11px] text-gray-500 uppercase tracking-wider mb-2">Scans Today</p>
      <p class="font-serif text-5xl font-normal leading-none text-gray-900 dark:text-gray-100">
        {stats.total_scans_today}
      </p>
    </div>
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <p class="font-mono text-[11px] text-gray-500 uppercase tracking-wider mb-2">Total Scans</p>
      <p class="font-serif text-5xl font-normal leading-none text-gray-900 dark:text-gray-100">
        {stats.total_scans_all}
      </p>
      <p class="text-xs text-gray-400 dark:text-gray-600 mt-1">in memory</p>
    </div>
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <p class="font-mono text-[11px] text-gray-500 uppercase tracking-wider mb-2">
        Cache Hit Rate
      </p>
      <p class="font-serif text-5xl font-normal leading-none text-gray-900 dark:text-gray-100">
        {stats.cache_hit_rate.toFixed(1)}<span class="text-lg text-gray-400 dark:text-gray-500"
          >%</span
        >
      </p>
      <p class="text-xs text-gray-400 dark:text-gray-600 mt-1">
        {stats.cache_hits} hits · {stats.cache_misses} misses
      </p>
    </div>
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <p class="font-mono text-[11px] text-gray-500 uppercase tracking-wider mb-2">Avg Duration</p>
      <p class="font-serif text-5xl font-normal leading-none text-gray-900 dark:text-gray-100">
        {stats.avg_duration_ms.toFixed(0)}<span class="text-lg text-gray-400 dark:text-gray-500"
          >ms</span
        >
      </p>
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Verdict breakdown -->
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">Verdict Breakdown</h2>
      {#if Object.keys(stats.verdict_counts).length === 0}
        <p class="text-sm text-gray-400 dark:text-gray-600">No data yet</p>
      {:else}
        {@const total = Object.values(stats.verdict_counts).reduce((a, b) => a + b, 0)}
        <div class="space-y-3">
          {#each Object.entries(stats.verdict_counts).sort((a, b) => b[1] - a[1]) as [verdict, count]}
            {@const pct = total > 0 ? (count / total) * 100 : 0}
            <div>
              <div class="flex justify-between items-center mb-1">
                <span class="text-sm font-medium {verdictTextColor(verdict)}">{verdict}</span>
                <span class="text-xs text-gray-500"
                  >{count}
                  <span class="text-gray-400 dark:text-gray-600">({pct.toFixed(0)}%)</span></span
                >
              </div>
              <div class="h-1.5 rounded-full bg-gray-200 dark:bg-gray-800 overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500 {verdict === 'Safe'
                    ? 'bg-emerald-500'
                    : verdict === 'Risky'
                      ? 'bg-red-500'
                      : 'bg-yellow-500'}"
                  style="width: {pct}%"
                ></div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Top domains -->
    <div
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5"
    >
      <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">
        Top Scanned Domains
      </h2>
      {#if !stats.top_domains || stats.top_domains.length === 0}
        <p class="text-sm text-gray-400 dark:text-gray-600">No data yet</p>
      {:else}
        {@const maxCount = stats.top_domains[0]?.count ?? 1}
        <div class="space-y-2.5">
          {#each stats.top_domains as item, i}
            <div class="flex items-center gap-3">
              <span class="text-xs text-gray-400 dark:text-gray-600 w-4 text-right shrink-0"
                >{i + 1}</span
              >
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between mb-0.5">
                  <span class="text-sm text-gray-700 dark:text-gray-300 font-mono truncate"
                    >{item.domain}</span
                  >
                  <span class="text-xs text-gray-500 ml-2 shrink-0">{item.count}</span>
                </div>
                <div class="h-1 rounded-full bg-gray-200 dark:bg-gray-800 overflow-hidden">
                  <div
                    class="h-full rounded-full bg-gray-900/70 dark:bg-gray-100/70"
                    style="width: {(item.count / maxCount) * 100}%"
                  ></div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{:else}
  <div class="text-center py-20 text-gray-400 dark:text-gray-600">
    <p class="text-sm">No stats available — start scanning URLs.</p>
  </div>
{/if}
