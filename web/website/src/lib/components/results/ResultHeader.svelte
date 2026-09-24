<script lang="ts">
  import { browser } from "$app/environment";
  import { PILL_OUTLINE, PILL_SOLID } from "../../ui/buttons";
  import { ICON } from "../../ui/icons";
  import Icon from "../Icon.svelte";

  // "Analysis Summary" title, the scanned domain, and the Report / Share actions.
  export let domain: string;
  export let reportOpen = false;

  let shareCopied = false;

  async function shareLink() {
    if (!browser) return;
    const url = window.location.href;
    if (navigator.share) {
      try {
        await navigator.share({ title: document.title, url });
        return;
      } catch {
        // user cancelled or share failed — fall through to clipboard
      }
    }
    try {
      await navigator.clipboard.writeText(url);
      shareCopied = true;
      setTimeout(() => (shareCopied = false), 2000);
    } catch (err) {
      console.error("Clipboard copy failed:", err);
    }
  }
</script>

<div class="animate-fadeIn">
  <h2
    class="font-serif text-4xl md:text-5xl font-normal tracking-[-0.01em] text-gray-900 dark:text-gray-100"
    id="analysis-summary"
  >
    Analysis Summary
  </h2>
  <div class="mt-3 flex items-center justify-between gap-3">
    <div class="flex items-center gap-2 min-w-0">
      <span class="hidden sm:inline text-gray-600 dark:text-gray-400 text-sm whitespace-nowrap"
        >Security profile for</span
      >
      <span
        class="min-w-0 inline-flex items-center gap-2 px-3 py-1 rounded-full border border-gray-300 dark:border-gray-800 font-mono text-[13px] text-gray-800 dark:text-gray-200"
      >
        <Icon path={ICON.globe} class="w-3.5 h-3.5 text-gray-400 flex-shrink-0" />
        <span class="truncate">{domain}</span>
      </span>
    </div>

    <div class="flex-shrink-0 flex items-center gap-2">
      <button
        type="button"
        class="{PILL_OUTLINE} p-2.5 sm:px-4 sm:py-2 {reportOpen
          ? 'border-gray-500 dark:border-gray-500 text-gray-900 dark:text-gray-100'
          : ''}"
        aria-expanded={reportOpen}
        aria-controls="report-panel"
        data-guide="report"
        aria-label="Report this result"
        title="Think this verdict is wrong? Let us know."
        on:click={() => (reportOpen = !reportOpen)}
      >
        <Icon path={ICON.flag} />
        <span class="hidden sm:inline">Report</span>
      </button>
      <button
        type="button"
        aria-label="Share this result"
        data-guide="share"
        class="{PILL_SOLID} p-2.5 sm:px-4 sm:py-2"
        on:click={shareLink}
        disabled={shareCopied}
      >
        {#if shareCopied}
          <Icon path={ICON.check} strokeWidth={2.5} />
          <span class="hidden sm:inline">Copied!</span>
        {:else}
          <Icon path={ICON.share} />
          <span class="hidden sm:inline">Share</span>
        {/if}
      </button>
    </div>
  </div>
</div>
