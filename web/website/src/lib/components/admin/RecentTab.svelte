<script lang="ts">
  import { onMount } from "svelte";
  import type { AdminClient } from "../../admin/client";
  import {
    copyWithFeedback,
    relativeTime,
    verdictColor,
    verdictDot,
    verdictTextColor,
  } from "../../admin/format";
  import type { Guard, ScanRecord } from "../../admin/types";
  import { encodeVerdict } from "../../utils";

  export let client: AdminClient;
  export let guard: Guard;

  let scans: ScanRecord[] = [];
  let loading = true;
  let copiedUrl: string | null = null;

  onMount(async () => {
    scans = (await guard(client.recent())) ?? [];
    loading = false;
  });

  function shareLink(scan: ScanRecord): string {
    const params = new URLSearchParams();
    params.set("q", scan.url);
    params.set("v", encodeVerdict(scan.verdict));
    params.set("s", String(scan.score));
    return `${window.location.origin}/?${params.toString()}`;
  }

  const copyScanLink = (scan: ScanRecord) => copyWithFeedback(scan.url, (v) => (copiedUrl = v));
</script>

{#if loading}
  <div class="space-y-2">
    {#each Array(6) as _}
      <div
        class="h-16 bg-white dark:bg-gray-900 rounded-2xl border border-gray-200 dark:border-gray-800 animate-pulse"
      ></div>
    {/each}
  </div>
{:else if scans.length === 0}
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
        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
    <p class="text-sm">No scans recorded yet.</p>
  </div>
{:else}
  <div class="space-y-2">
    {#each scans as scan}
      <div
        class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl px-4 py-3 hover:border-gray-300 dark:hover:border-gray-700 transition-colors"
      >
        <div class="flex items-center gap-3">
          <!-- Verdict dot -->
          <span class="w-2 h-2 rounded-full shrink-0 {verdictDot(scan.verdict)}"></span>

          <!-- Domain -->
          <p class="flex-1 min-w-0 text-sm text-gray-800 dark:text-gray-200 font-mono truncate">
            {scan.domain}
          </p>

          <!-- Verdict badge -->
          <span
            class="shrink-0 text-[11px] font-semibold px-2 py-0.5 rounded-full border {verdictColor(
              scan.verdict
            )}"
          >
            {scan.verdict}
          </span>

          <!-- Score -->
          <span class="shrink-0 text-sm font-bold tabular-nums {verdictTextColor(scan.verdict)}"
            >{scan.score}</span
          >

          <!-- Copy link -->
          <button
            on:click={() => copyScanLink(scan)}
            class="shrink-0 p-1.5 rounded-md text-gray-400 dark:text-gray-500 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            aria-label="Copy URL"
            title="Copy URL"
          >
            {#if copiedUrl === scan.url}
              <svg
                class="w-4 h-4 text-emerald-600 dark:text-emerald-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            {:else}
              <svg
                class="w-4 h-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                />
              </svg>
            {/if}
          </button>

          <!-- Open in new tab -->
          <a
            href={shareLink(scan)}
            target="_blank"
            rel="noopener noreferrer"
            class="shrink-0 p-1.5 rounded-md text-gray-400 dark:text-gray-500 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            aria-label="Open scan in new tab"
            title="Open scan in new tab"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
              />
            </svg>
          </a>
        </div>

        <div class="flex items-center gap-3 mt-1.5 pl-5 text-xs">
          <!-- URL -->
          <p class="flex-1 min-w-0 truncate text-gray-400 dark:text-gray-600">
            {scan.url}
          </p>

          <!-- Cached badge -->
          {#if scan.cached}
            <span
              class="shrink-0 text-[10px] font-semibold px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 border border-gray-200 dark:border-gray-700 font-mono"
              >cached</span
            >
          {/if}

          <!-- Duration -->
          <span class="shrink-0 text-gray-500 dark:text-gray-400 tabular-nums">{scan.duration}</span
          >

          <!-- Time -->
          <span class="shrink-0 text-gray-400 dark:text-gray-600 tabular-nums"
            >{relativeTime(scan.time)}</span
          >
        </div>
      </div>
    {/each}
  </div>
{/if}
