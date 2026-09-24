<script lang="ts">
  import { onDestroy } from "svelte";
  import { PILL_OUTLINE } from "../ui/buttons";
  import { ICON } from "../ui/icons";
  import type { AnalyzeResult } from "../types";
  import Icon from "./Icon.svelte";
  import PerformanceToggle from "./results/PerformanceToggle.svelte";
  import ReportPanel from "./results/ReportPanel.svelte";
  import ResultHeader from "./results/ResultHeader.svelte";
  import ScanError from "./results/ScanError.svelte";
  import SectionList from "./results/SectionList.svelte";
  import VerdictRow from "./results/VerdictRow.svelte";
  import FlagsGrid from "./sections/FlagsGrid.svelte";

  export let data: AnalyzeResult | null = null;
  export let screenshotUrl: string | null = null;
  export let screenshotLoading = false;
  export let screenshotFailed = false;
  export let loading = false;
  export let error: string | null = null;
  export let onScanAnother: (() => void) | undefined = undefined;

  let reportOpen = false;

  $: primary = data?.result;

  onDestroy(() => {
    if (screenshotUrl) URL.revokeObjectURL(screenshotUrl);
  });

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => {
      (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
    }, 400);
  }
</script>

{#if error}
  <ScanError message={error} />
{:else if loading}
  <div class="max-w-3xl mx-auto space-y-4">
    <div
      class="animate-pulse rounded-2xl border border-gray-200 dark:border-gray-800 bg-gray-100/60 dark:bg-gray-900/60 h-32"
    ></div>
    <div
      class="animate-pulse rounded-2xl border border-gray-200 dark:border-gray-800 bg-gray-100/60 dark:bg-gray-900/60 h-24"
    ></div>
  </div>
{:else if data}
  <section class="max-w-4xl mx-auto space-y-8 px-4">
    <ResultHeader domain={data.domain} bind:reportOpen />

    <ReportPanel
      bind:open={reportOpen}
      url={data.url}
      domain={data.domain}
      verdict={primary?.verdict}
      score={primary?.final_score}
    />

    <VerdictRow
      verdict={primary?.verdict}
      score={primary?.final_score}
      httpStatusCode={data.analysis?.http_status?.code ?? null}
      {screenshotUrl}
      {screenshotLoading}
      {screenshotFailed}
    />

    <div data-guide="flags" class="animate-fadeIn delay-200">
      <FlagsGrid reasons={primary?.reasons} />
    </div>

    <SectionList {data} />

    <div class="flex justify-center pt-2">
      <button
        type="button"
        data-guide="scan-another"
        class="group {PILL_OUTLINE} px-5 py-2.5"
        on:click={() => (onScanAnother ? onScanAnother() : scrollToTop())}
      >
        <Icon
          path={ICON.arrowUp}
          class="w-4 h-4 group-hover:-translate-y-0.5 transition-transform duration-200"
        />
        Scan another link
      </button>
    </div>

    {#if data.performance}
      <PerformanceToggle performance={data.performance} />
    {/if}
  </section>
{/if}
