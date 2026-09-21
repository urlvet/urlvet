<script lang="ts">
  import { browser } from "$app/environment";
  import { replaceState } from "$app/navigation";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
  import { api } from "../lib/api";
  import ResultSection from "../lib/components/ResultSection.svelte";
  import ScanProgress from "../lib/components/ScanProgress.svelte";
  import Shoutouts from "../lib/components/Shoutouts.svelte";
  import type { AnalyzeResult } from "../lib/types";
  import { encodeVerdict, formatUrl, isValidUrl, stripTrackers } from "../lib/utils";

  const schemaSoftwareApp = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "url.vet",
    alternateName: ["URLvet", "urlvet", "url vet"],
    url: "https://url.vet",
    description:
      "Free real-time URL scanner for detecting phishing, typosquatting, and malicious links. Instant verdict with full signal breakdown — no signup required.",
    applicationCategory: "SecurityApplication",
    applicationSubCategory: "URL Scanner",
    operatingSystem: "Web",
    browserRequirements: "Requires JavaScript",
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
    },
    creator: {
      "@type": "Person",
      name: "abhizaik",
      url: "https://abhizaik.com",
    },
    license: "https://github.com/urlvet/urlvet/blob/main/LICENSE",
    codeRepository: "https://github.com/urlvet/urlvet",
    featureList: [
      "Real-time URL scanning",
      "Phishing detection",
      "WHOIS and domain age analysis",
      "Typosquatting and homoglyph detection",
      "Redirect chain tracing",
      "TLS certificate verification",
      "IP and infrastructure analysis",
      "No signup required",
      "Free to use",
      "Open source (AGPL-3.0)",
    ],
  };

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
  $: trackers = input ? stripTrackers(formatUrl(input)) : { cleaned: "", removed: [] };
  $: hasTrackers = trackers.removed.length > 0;

  let trackerCopied = false;
  let trackerDismissed = false;
  $: if (input === "") trackerDismissed = false;
  async function copyCleanUrl() {
    try {
      await navigator.clipboard.writeText(trackers.cleaned);
      trackerCopied = true;
      setTimeout(() => {
        trackerCopied = false;
        trackerDismissed = true;
      }, 1200);
    } catch {}
  }

  let justPasted = false;

  function clearResult() {
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

  async function pasteFromClipboard() {
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        input = text.trim();
        error = null;
        formError = null;
        justPasted = true;
        setTimeout(() => (justPasted = false), 1500);
      }
    } catch {
      /* clipboard access denied */
    }
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
      } else {
        scanResult = res.data as AnalyzeResult;
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
    } finally {
      loading = false;
      scanDone = true;
      setTimeout(() => (scanDone = false), 1200);
    }
  }

  onMount(() => {
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
        if (scanResult || error) clearResult();
        else input = "";
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });
</script>

<svelte:head>
  {#if shareDomain}
    {@const ogVerdict = scanResult?.result?.verdict || data.verdict}
    {@const ogScore =
      scanResult?.result?.final_score ?? (data.score ? Number(data.score) : undefined)}
    {@const desc = ogVerdict
      ? `${ogVerdict} — url.vet scanned ${shareDomain}. See the full breakdown.`
      : `Sketchy link? url.vet scanned ${shareDomain} — check if it's actually safe to click.`}
    {@const ogImage = `https://url.vet/og?domain=${encodeURIComponent(shareDomain)}${ogVerdict ? `&v=${encodeURIComponent(ogVerdict)}` : ""}${ogScore !== undefined ? `&s=${ogScore}` : ""}`}
    <title>url.vet — is {shareDomain} sus?</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content="url.vet — is {shareDomain} sus?" />
    <meta property="og:description" content={desc} />
    <meta property="og:type" content="website" />
    <link rel="canonical" href="https://url.vet" />
    <meta
      property="og:url"
      content={currentUrl || `https://url.vet/?q=${encodeURIComponent(data.queryUrl)}`}
    />
    <meta property="og:image" content={ogImage} />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="url.vet — is {shareDomain} sus?" />
    <meta name="twitter:description" content={desc} />
    <meta name="twitter:image" content={ogImage} />
  {:else}
    <title>url.vet (URLvet) — sketchy link? just url.vet it</title>
    <meta
      name="description"
      content="Got a sketchy link? url.vet it (URLvet) — free, instant phishing verdict with no signup needed."
    />
    <meta
      name="keywords"
      content="URLvet, urlvet, url.vet, URL scanner, phishing detection, link checker, safe link checker, phishing URL checker, link safety checker"
    />
    <meta property="og:title" content="url.vet (URLvet) — just url.vet it" />
    <meta
      property="og:description"
      content="Sketchy link? Paste it. Get a verdict in seconds — safe, suspicious, or risky. Free & transparent."
    />
    <link rel="canonical" href="https://url.vet" />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://url.vet" />
    <meta property="og:image" content="https://url.vet/urlvet-light.png" />
    <meta name="twitter:card" content="summary" />
    <meta name="twitter:title" content="url.vet (URLvet) — just url.vet it" />
    <meta
      name="twitter:description"
      content="Sketchy link? url.vet it. Free, instant, no signup."
    />
    {@html `<script type="application/ld+json">${JSON.stringify(schemaSoftwareApp)}</script>`}
  {/if}
</svelte:head>

<section class="relative overflow-hidden">
  <div
    class={`max-w-5xl mx-auto px-6 ${isLanding ? "flex flex-col items-center text-center pt-16 md:pt-20 pb-12" : "py-12"}`}
  >
    <header class="relative mb-14 flex flex-col items-center text-center">
      <h1 class="relative text-5xl md:text-6xl font-extrabold tracking-tight z-10">
        <a href="/" on:click={() => (location.href = "/")} class="group">
          <span
            class="text-gray-500 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors"
            >url</span
          ><span class="text-blue-500 dark:text-blue-400">.</span><span
            class="text-gray-800 dark:text-white group-hover:text-gray-900 dark:group-hover:text-gray-100 transition-colors"
            >vet</span
          >
        </a>
        <span class="sr-only">&nbsp;(also known as URLvet)</span>
      </h1>

      <p class="relative mt-3 text-gray-500 dark:text-gray-400 text-sm md:text-base z-10">
        sketchy link? just url.vet it
      </p>
      <p class="relative mt-2 text-gray-400 dark:text-gray-500 text-xs md:text-sm z-10">
        Vet any URL for phishing, malware & suspicious redirects
      </p>
    </header>

    <form
      class={isLanding ? "w-full max-w-2xl" : "w-full"}
      on:submit|preventDefault={() => runAnalyze(input)}
    >
      <div class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            id="url-input"
            type="text"
            class={`w-full rounded-xl bg-white dark:bg-gray-900 border px-4 py-3.5 pr-24 text-base placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500/30 transition-all duration-200 ${formError ? "border-red-400 dark:border-red-600/70" : "border-gray-300 dark:border-gray-700/80"}`}
            placeholder="Paste URL (e.g. example.com)"
            bind:value={input}
            on:input={() => {
              if (formError) formError = null;
            }}
            autocomplete="url"
            inputmode="url"
            aria-invalid={formError ? "true" : undefined}
            aria-describedby={formError ? "url-error" : undefined}
            required
          />
          <div class="absolute right-0 top-0 bottom-0 flex items-center">
            <button
              type="button"
              on:click={pasteFromClipboard}
              class="h-full flex items-center gap-1.5 px-3.5 rounded-r-xl border-l text-[11px] font-medium transition-all duration-200 {justPasted
                ? 'border-gray-300 dark:border-white/5 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'border-gray-300 dark:border-white/5 bg-transparent hover:bg-gray-100 dark:hover:bg-white/5 text-gray-500 hover:text-gray-800 dark:text-gray-500 dark:hover:text-gray-200'}"
              aria-label="Paste from clipboard"
              title="Paste from clipboard"
            >
              {#if justPasted}
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                Pasted
              {:else}
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                  />
                </svg>
                Paste
              {/if}
            </button>
          </div>
        </div>

        <button
          type="submit"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white text-base font-semibold shadow-lg shadow-blue-900/30 transition-all duration-200 disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-500 active:scale-95"
          disabled={loading}
          aria-busy={loading}
          aria-label={loading ? "Scanning URL, please wait" : "Scan Now"}
        >
          {#if loading}
            <svg
              class="w-4 h-4 animate-spin"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 4v4m0 8v4m8-8h-4M4 12H0"
              />
            </svg>
            Scanning..
          {:else}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z"
              />
            </svg>
            Scan Now
          {/if}
        </button>
      </div>

      {#if formError}
        <p id="url-error" class="text-red-400 text-xs mt-2 text-left" role="alert">{formError}</p>
      {:else if input.length > 1848}
        <p
          class="text-[11px] mt-1.5 text-right {input.length >= 2048
            ? 'text-red-500'
            : 'text-gray-500'}"
        >
          {2048 - input.length} chars remaining
        </p>
      {/if}

      {#if hasTrackers && !trackerDismissed}
        <div class="mt-2 px-1 text-xs space-y-1.5 text-center">
          <p class="text-gray-500">
            <span class="text-gray-400 font-medium"
              >{trackers.removed.length} tracker{trackers.removed.length > 1 ? "s" : ""} found —</span
            >
            <span class="break-all"> {trackers.removed.join(", ")}</span>
          </p>
          <p class="text-gray-600">Click to copy the clean link</p>
          <button
            type="button"
            on:click={copyCleanUrl}
            class="w-full font-mono text-xs break-all transition-colors cursor-copy {trackerCopied
              ? 'text-emerald-400'
              : 'text-gray-500 hover:text-gray-300'}"
          >
            {trackerCopied ? "Copied!" : trackers.cleaned}
          </button>
        </div>
      {/if}
    </form>

    {#if isLanding}
      <div class="mt-5 flex flex-wrap justify-center items-center gap-2">
        <span class="text-[11px] text-gray-500 mr-0.5">Try:</span>
        {#each [{ label: "google.com", url: "google.com", dot: "bg-emerald-500", hint: "Safe" }, { label: "аррӏе.com", url: "аррӏе.com", dot: "bg-red-500", hint: "Risky" }, { label: "wikipedia.org", url: "wikipedia.org", dot: "bg-emerald-500", hint: "Safe" }, { label: "pаypal.com", url: "pаypal.com", dot: "bg-red-500", hint: "Risky" }] as example}
          <button
            type="button"
            on:click={() => {
              input = example.url;
              runAnalyze(example.url);
            }}
            title="Expected: {example.hint}"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-white dark:bg-gray-900 border border-gray-300 dark:border-gray-800 hover:border-gray-400 dark:hover:border-gray-600 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 text-xs transition-all"
          >
            {example.label}
            <span
              class={`text-xs font-semibold ${example.hint === "Safe" ? "text-emerald-500" : "text-red-500"}`}
              >{example.hint}</span
            >
          </button>
        {/each}
      </div>

      <div class="mt-5 flex flex-wrap justify-center gap-x-5 gap-y-2">
        {#each ["Free & open source", "No signup required", "Explains every verdict", "See it before you click"] as pill}
          <span
            class="flex items-center gap-1.5 text-sm font-medium text-gray-600 dark:text-gray-400"
          >
            <svg
              class="w-3.5 h-3.5 text-emerald-500 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
            {pill}
          </span>
        {/each}
      </div>
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
