<script lang="ts">
  import { onMount } from "svelte";
  import type { AdminClient } from "../../admin/client";
  import { relativeTime } from "../../admin/format";
  import type { ErrorRecord, Guard } from "../../admin/types";

  export let client: AdminClient;
  export let guard: Guard;
  /** Number of errors, shown as a badge on the tab. */
  export let count = 0;

  let errors: ErrorRecord[] = [];
  let loading = true;

  onMount(async () => {
    errors = (await guard(client.errors())) ?? [];
    count = errors.length;
    loading = false;
  });
</script>

{#if loading}
  <div class="space-y-2">
    {#each Array(4) as _}
      <div
        class="h-20 bg-white dark:bg-gray-900 rounded-2xl border border-gray-200 dark:border-gray-800 animate-pulse"
      ></div>
    {/each}
  </div>
{:else if errors.length === 0}
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
        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
    <p class="text-sm">No errors recorded. All good!</p>
  </div>
{:else}
  <div class="space-y-2">
    {#each errors as err}
      <div
        class="bg-white dark:bg-gray-900 border border-red-200 dark:border-red-900/50 rounded-2xl px-4 py-3"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span
                class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-red-50 dark:bg-red-500/10 text-red-700 dark:text-red-400 border border-red-200 dark:border-red-500/20"
                >{err.task}</span
              >
              <span class="text-xs text-gray-400 dark:text-gray-600">{relativeTime(err.time)}</span>
            </div>
            <p class="text-sm text-red-600 dark:text-red-300 font-mono break-all">
              {err.error}
            </p>
            {#if err.url}
              <p class="text-xs text-gray-400 dark:text-gray-600 mt-1 truncate">
                {err.url}
              </p>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
