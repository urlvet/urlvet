<script lang="ts">
  import ScreenshotViewer from "../sections/ScreenshotViewer.svelte";
  import VerdictCard from "../sections/VerdictCard.svelte";

  // Verdict card with the page screenshot beside it (stacked on phones).
  export let verdict: string | undefined;
  export let score: number | undefined;
  export let httpStatusCode: number | null;
  export let screenshotUrl: string | null = null;
  export let screenshotLoading = false;
  export let screenshotFailed = false;

  $: unreachable = httpStatusCode === 0 || httpStatusCode === null;
  $: unavailableReason = !screenshotFailed
    ? null
    : unreachable
      ? "No web content detected. Preview unavailable."
      : "Screenshot unavailable — the site may be blocking automated access or failed to respond.";
</script>

<div class="flex flex-col md:flex-row gap-4 items-stretch animate-fadeIn delay-100">
  <div data-guide="verdict" class="flex-1 min-w-0">
    <VerdictCard {verdict} finalScore={score} {unreachable} />
  </div>
  {#if screenshotLoading || screenshotUrl}
    <div
      data-guide="screenshot"
      class="md:w-56 md:flex-shrink-0 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-3 flex flex-col gap-2"
    >
      <p
        class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
      >
        Screenshot
      </p>
      <ScreenshotViewer
        {screenshotUrl}
        loading={screenshotLoading}
        failed={screenshotFailed}
        {unavailableReason}
        compact={true}
      />
    </div>
  {/if}
</div>
