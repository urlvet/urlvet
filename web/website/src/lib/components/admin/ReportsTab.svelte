<script lang="ts">
  import { onMount } from "svelte";
  import type { AdminClient } from "../../admin/client";
  import {
    copyWithFeedback,
    hostOf,
    relativeTime,
    verdictDot,
    verdictTextColor,
  } from "../../admin/format";
  import type { Guard, ReportRecord } from "../../admin/types";

  export let client: AdminClient;
  export let guard: Guard;
  /** Number of reports, shown as a badge on the tab. */
  export let count = 0;

  let reports: ReportRecord[] = [];
  let loading = true;
  let filterText = "";
  let filterExpected = "";
  let newestFirst = true;
  let copiedUrl: string | null = null;
  /** Report waiting for delete confirmation, and the one being deleted. */
  let confirmId: string | null = null;
  let deletingId: string | null = null;

  onMount(async () => {
    reports = (await guard(client.reports())) ?? [];
    count = reports.length;
    loading = false;
  });

  const SUGGESTIONS = [
    { value: "", label: "All" },
    { value: "Safe", label: "Should be Safe" },
    { value: "Suspicious", label: "Should be Suspicious" },
    { value: "Risky", label: "Should be Risky" },
    { value: "none", label: "No suggestion" },
  ];

  const matches = (r: ReportRecord, value: string) =>
    !value || (value === "none" ? !r.expected_verdict : r.expected_verdict === value);

  // Report times are stored in IST as "YYYY-MM-DD HH:MM:SS"; this layout sorts as a string.
  const age = (t: string) => relativeTime(`${t.replace(" ", "T")}+05:30`);

  $: withComment = reports.filter((r) => r.comment).length;

  $: filtered = reports
    .filter((r) => {
      const q = filterText.toLowerCase();
      const matchText =
        !q || r.url.toLowerCase().includes(q) || (r.comment ?? "").toLowerCase().includes(q);
      return matches(r, filterExpected) && matchText;
    })
    .sort((a, b) => (newestFirst ? b.time.localeCompare(a.time) : a.time.localeCompare(b.time)));

  function reportScanLink(r: ReportRecord): string {
    return `${window.location.origin}/?q=${encodeURIComponent(r.url)}`;
  }

  const copyUrl = (url: string) => copyWithFeedback(url, (v) => (copiedUrl = v));

  async function resolve(id: string) {
    deletingId = id;
    const done = await guard(client.deleteReport(id).then(() => true));
    if (done) {
      reports = reports.filter((r) => r.id !== id);
      count = reports.length;
    }
    deletingId = null;
    confirmId = null;
  }

  const ACTION =
    "inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs transition-colors";
  const ACTION_OUTLINE = `${ACTION} border border-gray-300 dark:border-gray-800 text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-900 dark:hover:text-gray-100`;
</script>

{#if loading}
  <div class="space-y-3">
    {#each Array(4) as _}
      <div class="h-28 rounded-2xl bg-gray-200/50 dark:bg-gray-800/40 animate-pulse"></div>
    {/each}
  </div>
{:else if reports.length === 0}
  <div class="text-center py-24">
    <p class="font-serif text-3xl text-gray-900 dark:text-gray-100">All clear.</p>
    <p class="mt-2 text-sm text-gray-500">No reports waiting. Resolved reports are deleted.</p>
  </div>
{:else}
  <!-- Toolbar: search, sort, then suggestion filters -->
  <div class="flex gap-2 sm:gap-3">
    <input
      bind:value={filterText}
      type="search"
      placeholder="Search URL or comment…"
      class="flex-1 min-w-0 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full px-4 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:border-gray-400 dark:focus:border-gray-600"
    />
    <button
      type="button"
      on:click={() => (newestFirst = !newestFirst)}
      class="{ACTION_OUTLINE} flex-shrink-0 justify-center px-4 py-2 text-sm"
      aria-label="Toggle sort order"
    >
      {newestFirst ? "Newest" : "Oldest"}<span class="hidden sm:inline -ml-1">&nbsp;first</span>
      <svg
        class="w-3.5 h-3.5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
        aria-hidden="true"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M3 7.5L7.5 3m0 0L12 7.5M7.5 3v13.5m13.5 0L16.5 21m0 0L12 16.5m4.5 4.5V7.5"
        />
      </svg>
    </button>
  </div>

  <div class="mt-4 flex flex-wrap items-center gap-x-1 gap-y-2 text-xs">
    {#each SUGGESTIONS as s}
      {@const n = reports.filter((r) => matches(r, s.value)).length}
      <button
        type="button"
        on:click={() => (filterExpected = s.value)}
        aria-pressed={filterExpected === s.value}
        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full transition-colors {filterExpected ===
        s.value
          ? 'bg-gray-200 dark:bg-gray-800 text-gray-900 dark:text-gray-100'
          : 'text-gray-500 hover:text-gray-900 dark:hover:text-gray-100'}"
      >
        {s.label}
        <span class="font-mono text-[10px] text-gray-400">{n}</span>
      </button>
    {/each}
    <span class="ml-auto font-mono text-[11px] text-gray-400">{withComment} with comments</span>
  </div>

  {#if filtered.length === 0}
    <p class="text-center py-16 text-sm text-gray-500">No reports match these filters.</p>
  {:else}
    <div class="mt-5 space-y-3">
      {#each filtered as r (r.id)}
        <article
          class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-5 transition-opacity {deletingId ===
          r.id
            ? 'opacity-50'
            : ''}"
        >
          <header class="flex items-start justify-between gap-4">
            <div class="min-w-0">
              <h3 class="text-[15px] font-medium text-gray-900 dark:text-gray-100 break-all">
                {hostOf(r.url)}
              </h3>
              <p class="mt-0.5 font-mono text-xs text-gray-500 break-all">{r.url}</p>
            </div>
            <time
              class="flex-shrink-0 font-mono text-[11px] text-gray-500 whitespace-nowrap pt-0.5"
              title="{r.time} IST">{age(r.time)}</time
            >
          </header>

          <p class="mt-3 flex flex-wrap items-center gap-x-2 gap-y-1 text-[13px] text-gray-500">
            <span>We said</span>
            {#if r.verdict}
              <span
                class="inline-flex items-center gap-1.5 font-medium {verdictTextColor(r.verdict)}"
              >
                <span class="w-1.5 h-1.5 rounded-full {verdictDot(r.verdict)}"></span>{r.verdict}
              </span>
              <span class="font-mono text-[11px]">{r.score}</span>
            {:else}
              <span>unknown</span>
            {/if}
            {#if r.expected_verdict}
              <span class="text-gray-300 dark:text-gray-700" aria-hidden="true">·</span>
              <span>user suggests {r.expected_verdict}</span>
            {/if}
          </p>

          {#if r.comment}
            <blockquote
              class="mt-4 border-l-2 border-gray-300 dark:border-gray-700 pl-4 text-[15px] leading-relaxed text-gray-800 dark:text-gray-200 whitespace-pre-wrap break-words"
            >
              {r.comment}
            </blockquote>
          {/if}

          <footer class="mt-4 flex flex-wrap items-center gap-2">
            <a
              href={reportScanLink(r)}
              target="_blank"
              rel="noopener noreferrer"
              class={ACTION_OUTLINE}
            >
              Re-scan
              <svg
                class="w-3.5 h-3.5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </a>
            <button type="button" on:click={() => copyUrl(r.url)} class={ACTION_OUTLINE}>
              {copiedUrl === r.url ? "Copied" : "Copy URL"}
            </button>

            <span class="ml-auto flex items-center gap-2">
              {#if confirmId === r.id}
                <span class="text-xs text-gray-500">Delete this report?</span>
                <button
                  type="button"
                  on:click={() => (confirmId = null)}
                  class="{ACTION} text-gray-500 hover:text-gray-900 dark:hover:text-gray-100"
                  >Cancel</button
                >
                <button
                  type="button"
                  on:click={() => resolve(r.id)}
                  disabled={deletingId === r.id}
                  class="{ACTION} bg-gray-900 dark:bg-gray-100 text-gray-50 dark:text-gray-900 hover:bg-gray-700 dark:hover:bg-white"
                  >Delete report</button
                >
              {:else}
                <button
                  type="button"
                  on:click={() => (confirmId = r.id)}
                  class={ACTION_OUTLINE}
                  title="Mark as fixed. The report is deleted, as the privacy page promises."
                >
                  <svg
                    class="w-3.5 h-3.5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="2"
                    aria-hidden="true"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M4.5 12.75l6 6 9-13.5"
                    />
                  </svg>
                  Resolve
                </button>
              {/if}
            </span>
          </footer>
        </article>
      {/each}
    </div>
  {/if}
{/if}
