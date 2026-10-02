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
  import {
    encodeVerdict,
    extractLinks,
    formatUrl,
    getDomainFromUrl,
    isValidUrl,
  } from "../lib/utils";

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
  /** Links found in a pasted message, when there's more than one to choose from. */
  let linkChoices: string[] = [];

  $: isLanding = !scanResult && !loading && !error && !formError;
  // The browser tab follows the scan too: "Checking …" while it runs, then the
  // verdict with a coloured dot as the icon.
  const VERDICT_TAB: Record<string, string> = {
    Safe: "#10b981",
    Suspicious: "#eab308",
    Risky: "#ef4444",
  };
  $: verdictNow = scanResult?.result?.verdict;
  $: liveTitle = loading
    ? `Checking ${getDomainFromUrl(formatUrl(input)) || "link"}… · url.vet`
    : scanResult && verdictNow && VERDICT_TAB[verdictNow]
      ? `${verdictNow} — ${scanResult.domain} · url.vet`
      : undefined;
  $: verdictDot = !loading && verdictNow ? VERDICT_TAB[verdictNow] : undefined;

  // The logo dot mirrors the scan: blinking while it runs, then the verdict's colour.
  $: heroDot = loading
    ? ("scanning" as const)
    : ((scanResult?.result?.verdict as "Safe" | "Suspicious" | "Risky" | undefined) ??
      ("idle" as const));
  // Choices belong to the text they came from; editing it dismisses them.
  let choicesFor = "";
  $: if (linkChoices.length && input !== choicesFor) linkChoices = [];
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

  // Accepts a link, or a whole message with links in it (e.g. forwarded from WhatsApp).
  function submit(raw: string) {
    linkChoices = [];
    const text = raw.trim();
    if (/\s/.test(text) || !isValidUrl(formatUrl(text))) {
      const links = extractLinks(text);
      if (links.length > 1) {
        choicesFor = raw;
        linkChoices = links;
        return;
      }
      if (links.length === 1) {
        input = links[0];
        return runAnalyze(links[0]);
      }
      if (/\s/.test(text)) {
        formError = "No link found in that text";
        return;
      }
    }
    runAnalyze(text);
  }

  function chooseLink(link: string) {
    linkChoices = [];
    input = link;
    runAnalyze(link);
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
        // Drop anything the share menu passed in: the original message stays private.
        for (const k of ["url", "text", "title"]) share.searchParams.delete(k);
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
    const params = new URLSearchParams(window.location.search);
    const q = params.get("q");
    // Shared from another app via the phone's share menu (share_target in manifest.webmanifest).
    const shared = ["url", "text", "title"]
      .map((k) => params.get(k))
      .filter(Boolean)
      .join(" ");
    if (q) {
      input = q;
      runAnalyze(q);
    } else if (shared) {
      // The scan swaps these params for ?q= once it finishes.
      input = shared;
      submit(shared);
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
  {liveTitle}
  liveDot={verdictDot}
/>

<section class="relative overflow-x-clip">
  {#if isLanding}
    <Backdrop />
  {/if}
  <div
    class={`relative max-w-5xl mx-auto px-6 ${isLanding ? "flex flex-col items-center text-center pt-8 min-[400px]:pt-12 sm:pt-16 md:pt-24 pb-12" : "pt-10 md:pt-12 pb-12"}`}
  >
    <Hero {isLanding} dot={heroDot} />

    <SearchBar
      bind:input
      bind:formError
      {loading}
      {isLanding}
      onSubmit={submit}
      onPaste={() => (error = null)}
    />

    {#if linkChoices.length}
      <div class="mt-5 w-full max-w-2xl mx-auto text-left" role="group" aria-label="Links found">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          Found {linkChoices.length} links in that message. Which one should we check?
        </p>
        <ul
          class="mt-3 divide-y divide-gray-200 dark:divide-gray-800 border-y border-gray-200 dark:border-gray-800"
        >
          {#each linkChoices as link}
            <li>
              <button
                type="button"
                on:click={() => chooseLink(link)}
                class="group w-full flex items-center gap-3 py-3 text-left font-mono text-[13px] text-gray-800 dark:text-gray-200 hover:text-gray-900 dark:hover:text-white"
              >
                <span class="flex-1 min-w-0 truncate">{link}</span>
                <span
                  class="flex-shrink-0 font-sans text-xs text-gray-400 group-hover:text-gray-900 dark:group-hover:text-gray-100 transition-colors"
                  >Check →</span
                >
              </button>
            </li>
          {/each}
        </ul>
      </div>
    {/if}

    {#if isLanding && !linkChoices.length}
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
