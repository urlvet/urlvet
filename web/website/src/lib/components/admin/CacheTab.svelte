<script lang="ts">
  import { onMount } from "svelte";
  import type { AdminClient } from "../../admin/client";
  import { formatJSON, formatTTL, ttlColor } from "../../admin/format";
  import type { CacheEntry, Guard } from "../../admin/types";

  export let client: AdminClient;
  export let guard: Guard;

  let cacheEntries: CacheEntry[] = [];
  let loading = true;
  let flushConfirm = false;
  let deletingKey: string | null = null;
  let expandedKeys = new Set<string>();
  let filterText = "";
  let filterPrefix = "";

  onMount(load);

  async function load() {
    loading = true;
    cacheEntries = (await guard(client.cache())) ?? [];
    loading = false;
  }

  $: prefixes = [...new Set(cacheEntries.map((e) => e.prefix))].sort();
  $: filteredCache = cacheEntries.filter((e) => {
    const matchPrefix = !filterPrefix || e.prefix === filterPrefix;
    const matchText = !filterText || e.key.toLowerCase().includes(filterText.toLowerCase());
    return matchPrefix && matchText;
  });

  async function deleteKey(key: string) {
    deletingKey = key;
    const ok = await guard(client.deleteKey(key).then(() => true));
    if (ok) {
      cacheEntries = cacheEntries.filter((e) => e.key !== key);
      expandedKeys.delete(key);
      expandedKeys = expandedKeys;
    }
    deletingKey = null;
  }

  async function flushAll() {
    flushConfirm = false;
    loading = true;
    const ok = await guard(client.flushCache().then(() => true));
    if (ok) {
      cacheEntries = [];
      expandedKeys = new Set();
    }
    loading = false;
  }

  function toggleExpand(key: string) {
    if (expandedKeys.has(key)) expandedKeys.delete(key);
    else expandedKeys.add(key);
    expandedKeys = expandedKeys;
  }
</script>

<!-- Cache header -->
<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
  <p class="text-sm text-gray-500">
    {cacheEntries.length} key{cacheEntries.length !== 1 ? "s" : ""} in store
    {#if filteredCache.length !== cacheEntries.length}
      <span class="text-gray-400 dark:text-gray-600"> · {filteredCache.length} shown</span>
    {/if}
  </p>
  <div class="flex items-center gap-2">
    {#if !flushConfirm}
      <button
        on:click={() => (flushConfirm = true)}
        disabled={loading || cacheEntries.length === 0}
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-red-50 dark:bg-red-900/30 hover:bg-red-100 dark:hover:bg-red-900/60 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 text-xs font-medium transition-colors disabled:opacity-40"
      >
        <svg
          class="w-3.5 h-3.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          />
        </svg>
        Flush All
      </button>
    {:else}
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-xs text-red-500 dark:text-red-400"
          >Delete all {cacheEntries.length} keys?</span
        >
        <button
          on:click={flushAll}
          class="px-3 py-1.5 rounded-lg bg-red-600 dark:bg-red-700 hover:bg-red-700 dark:hover:bg-red-600 text-white text-xs font-semibold transition-colors"
          >Confirm</button
        >
        <button
          on:click={() => (flushConfirm = false)}
          class="px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs font-semibold transition-colors"
          >Cancel</button
        >
      </div>
    {/if}
  </div>
</div>

<!-- Filters -->
{#if cacheEntries.length > 0}
  <div class="flex flex-col sm:flex-row gap-3 mb-4">
    <input
      bind:value={filterText}
      placeholder="Search keys…"
      class="flex-1 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full px-4 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:border-gray-400 dark:focus:border-gray-600"
    />
    <select
      bind:value={filterPrefix}
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full px-4 py-2 text-sm text-gray-900 dark:text-gray-200 focus:outline-none focus:border-gray-400 dark:focus:border-gray-600"
    >
      <option value="">All prefixes</option>
      {#each prefixes as p}
        <option value={p}>{p}</option>
      {/each}
    </select>
  </div>
{/if}

{#if loading && cacheEntries.length === 0}
  <div class="space-y-2">
    {#each Array(5) as _}
      <div
        class="h-14 bg-white dark:bg-gray-900 rounded-2xl border border-gray-200 dark:border-gray-800 animate-pulse"
      ></div>
    {/each}
  </div>
{:else if !loading && cacheEntries.length === 0}
  <div class="text-center py-24 text-gray-400 dark:text-gray-600">
    <svg
      class="w-10 h-10 mx-auto mb-3 opacity-30"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="1.5"
        d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
      />
    </svg>
    <p class="text-sm">Cache is empty.</p>
  </div>
{:else if filteredCache.length === 0 && cacheEntries.length > 0}
  <p class="text-center text-sm text-gray-400 dark:text-gray-600 py-10">
    No keys match your filter.
  </p>
{:else}
  <div class="space-y-2">
    {#each filteredCache as entry (entry.key)}
      {@const expanded = expandedKeys.has(entry.key)}
      <div
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl overflow-hidden hover:border-gray-300 dark:hover:border-gray-700 transition-colors"
      >
        <div class="px-4 py-3">
          <div class="flex items-center gap-3">
            <button
              on:click={() => toggleExpand(entry.key)}
              class="flex-shrink-0 text-gray-400 dark:text-gray-600 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              aria-label={expanded ? "Collapse" : "Expand"}
            >
              <svg
                class="w-4 h-4 transition-transform duration-150 {expanded ? 'rotate-90' : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fill-rule="evenodd"
                  d="M7.293 4.293a1 1 0 011.414 0l5 5a1 1 0 010 1.414l-5 5a1 1 0 01-1.414-1.414L11.586 10 7.293 5.707a1 1 0 010-1.414z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
            <button
              on:click={() => toggleExpand(entry.key)}
              class="flex-1 min-w-0 text-left font-mono text-sm text-gray-700 dark:text-gray-300 truncate hover:text-gray-900 dark:hover:text-white transition-colors"
              title={entry.key}>{entry.key}</button
            >
            <button
              on:click={() => deleteKey(entry.key)}
              disabled={deletingKey === entry.key}
              class="flex-shrink-0 ml-1 p-1.5 rounded-md text-gray-300 dark:text-gray-700 hover:text-red-500 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors disabled:opacity-40"
              aria-label="Delete key"
              title="Delete key"
            >
              <svg
                class="w-4 h-4 {deletingKey === entry.key ? 'animate-spin' : ''}"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>
          <div class="flex items-center gap-2 mt-1.5 pl-7">
            <span
              class="flex-shrink-0 text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-800 text-gray-500 border border-gray-200 dark:border-gray-700"
              >{entry.prefix}</span
            >
            <span
              class="flex-shrink-0 text-xs font-medium {ttlColor(entry.ttl_seconds)} tabular-nums"
              >{formatTTL(entry.ttl_seconds)}</span
            >
          </div>
        </div>
        {#if expanded}
          <div class="border-t border-gray-200 dark:border-gray-800 px-4 py-3">
            <pre
              class="text-xs text-emerald-700 dark:text-emerald-300 font-mono bg-gray-100 dark:bg-gray-950 rounded-lg p-4 overflow-x-auto whitespace-pre-wrap break-all leading-relaxed max-h-80 overflow-y-auto">{formatJSON(
                entry.value
              )}</pre>
          </div>
        {/if}
      </div>
    {/each}
  </div>
{/if}
