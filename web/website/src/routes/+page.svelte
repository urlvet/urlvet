<script lang="ts">
  import { browser } from "$app/environment";
  import { replaceState } from "$app/navigation";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
  import { api } from "../lib/api";
  import Backdrop from "../lib/components/home/Backdrop.svelte";
  import Hero from "../lib/components/home/Hero.svelte";
  import LandingExtras from "../lib/components/home/LandingExtras.svelte";
  import PageMeta from "../lib/components/home/PageMeta.svelte";
  import SearchBar from "../lib/components/home/SearchBar.svelte";
  import { scanState } from "../lib/components/guide/store";
  import ResultSection from "../lib/components/ResultSection.svelte";
  import ScanProgress from "../lib/components/ScanProgress.svelte";
  import Shoutouts from "../lib/components/Shoutouts.svelte";
  import type { AnalyzeResult } from "../lib/types";
  import { encodeVerdict, formatUrl, isValidUrl } from "../lib/utils";

  // Page load data from +page.ts — runs server-side so bots get correct OG meta tags.
  export let data: {
    queryDomain: string;
    queryUrl: string;
    formattedQueryUrl: string;
    verdict: string;
    score: string;
  };

  let input = "";
  let loading = false;
  let error: string | null = null;
  let formError: string | null = null;
  let scanResult: AnalyzeResult | null = null;
  let screenshotUrl: string | null = null;
  let screenshotLoading = false;
  let screenshotFailed = false;
  let scanDone = false;

  $: isLanding = !scanResult && !loading && !error && !formError;
  $: currentUrl = browser ? window.location.href : "";
  $: shareDomain = scanResult?.domain || data.queryDomain;

  function clearResult() {
    scanState.set({ status: "idle" });
    scanResult = null;
    input = "";
    error = null;
    formError = null;
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
      screenshotUrl = null;
    }
    screenshotLoading = false;
    screenshotFailed = false;
    if (browser) replaceState(window.location.pathname, {});
    setTimeout(() => {
      (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
    }, 50);
  }

  async function runAnalyze(q: string) {
    const url = formatUrl(q);
    if (!isValidUrl(url)) {
      formError = "Please enter a valid URL";
      return;
    }

    scanDone = false;
    error = null;
    formError = null;
    scanResult = null;
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
      screenshotUrl = null;
    }
    screenshotLoading = true;
    screenshotFailed = false;
    scanState.set({ status: "scanning" });

    try {
      api
        .screenshot(url)
        .then((res) => {
          if (res.data) screenshotUrl = res.data as string;
          else screenshotFailed = true;
        })
        .catch(() => {
          screenshotFailed = true;
        })
        .finally(() => {
          screenshotLoading = false;
        });

      loading = true;
      const res = await api.analyze(url);
      if (res.error) {
        error = res.error;
        scanState.set({ status: "error", message: res.error });
      } else {
        scanResult = res.data as AnalyzeResult;
        scanState.set({ status: "done", result: scanResult });
        const share = new URL(window.location.href);
        share.searchParams.set("q", url);
        if (scanResult.result?.verdict)
          share.searchParams.set("v", encodeVerdict(scanResult.result.verdict));
        if (scanResult.result?.final_score !== undefined)
          share.searchParams.set("s", String(scanResult.result.final_score));
        replaceState(share.toString(), {});
      }
    } catch {
      error = "Analyze request failed";
      scanState.set({ status: "error", message: error });
    } finally {
      loading = false;
      scanDone = true;
      setTimeout(() => (scanDone = false), 1200);
    }
  }

  onMount(() => {
    // Start fresh so Vetty doesn't react to a result from an earlier visit to this page.
    scanState.set({ status: "idle" });
    const q = new URLSearchParams(window.location.search).get("q");
    if (q) {
      input = q;
      runAnalyze(q);
    } else {
      (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
    }

    const onKey = (e: KeyboardEvent) => {
      if (
        e.key === "/" &&
        !(e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement)
      ) {
        e.preventDefault();
        (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
      }
      if (e.key === "Escape") {
        // Escape belongs to an open dialog (Vetty, the tour) before it clears the page.
        if (document.querySelector('[role="dialog"]')) return;
        if (scanResult || error) clearResult();
        else input = "";
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });
</script>

<PageMeta
  {shareDomain}
  verdict={scanResult?.result?.verdict || data.verdict || undefined}
  score={scanResult?.result?.final_score ?? (data.score ? Number(data.score) : undefined)}
  queryUrl={data.queryUrl}
  {currentUrl}
/>

<section class="relative overflow-x-clip">
  {#if isLanding}
    <Backdrop />
  {/if}
  <div
    class={`relative max-w-5xl mx-auto px-6 ${isLanding ? "flex flex-col items-center text-center pt-8 min-[400px]:pt-12 sm:pt-16 md:pt-24 pb-12" : "pt-10 md:pt-12 pb-12"}`}
  >
    <Hero {isLanding} />

    <SearchBar
      bind:input
      bind:formError
      {loading}
      {isLanding}
      onSubmit={runAnalyze}
      onPaste={() => (error = null)}
    />

    {#if isLanding}
      <LandingExtras
        onTry={(url) => {
          input = url;
          runAnalyze(url);
        }}
      />
    {/if}

    <div class={`w-full ${isLanding ? "mt-12" : "mt-8"}`} aria-live="polite">
      {#if loading || scanDone}
        <ScanProgress {loading} done={scanDone} />
      {:else}
        <div in:fade={{ duration: 180 }}>
          <ResultSection
            data={scanResult}
            {loading}
            {error}
            {screenshotUrl}
            {screenshotLoading}
            {screenshotFailed}
            onScanAnother={clearResult}
          />
        </div>
      {/if}
    </div>

    {#if isLanding}
      <div class="mt-8 w-full">
        <Shoutouts />
      </div>
    {/if}
  </div>
</section>
