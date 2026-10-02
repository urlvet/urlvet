<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy } from "svelte";
  // <head> tags: per-scan share card when a domain is known, the site card otherwise.
  export let shareDomain: string;
  export let verdict: string | undefined;
  export let score: number | undefined;
  export let queryUrl: string;
  export let currentUrl: string;
  /** Tab title for the visitor's own scan (client-side only; shared links keep their SEO title). */
  export let liveTitle: string | undefined = undefined;
  /** Colour of the verdict dot, shown as the tab's icon after a scan. */
  export let liveDot: string | undefined = undefined;

  // Browsers keep showing the first icon they loaded, so swap the existing
  // <link rel="icon"> in place (and put the original back) instead of adding one.
  let originalIcon: string | null = null;
  $: if (browser) {
    const link = document.querySelector<HTMLLinkElement>('link[rel="icon"]');
    if (link) {
      originalIcon ??= link.href;
      link.href = liveDot ? dotIcon(liveDot) : originalIcon;
    }
  }
  onDestroy(() => {
    const link = browser ? document.querySelector<HTMLLinkElement>('link[rel="icon"]') : null;
    if (link && originalIcon) link.href = originalIcon;
  });

  const dotIcon = (color: string) =>
    `data:image/svg+xml,${encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="7" fill="${color}"/></svg>`)}`;

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
</script>

<svelte:head>
  {#if shareDomain}
    {@const ogVerdict = verdict}
    {@const ogScore = score}
    {@const desc = ogVerdict
      ? `${ogVerdict} — url.vet scanned ${shareDomain}. See the full breakdown.`
      : `Sketchy link? url.vet scanned ${shareDomain} — check if it's actually safe to click.`}
    {@const ogImage = `https://url.vet/og?domain=${encodeURIComponent(shareDomain)}${ogVerdict ? `&v=${encodeURIComponent(ogVerdict)}` : ""}${ogScore !== undefined ? `&s=${ogScore}` : ""}`}
    <title>{liveTitle ?? `url.vet — is ${shareDomain} sus?`}</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content="url.vet — is {shareDomain} sus?" />
    <meta property="og:description" content={desc} />
    <meta property="og:type" content="website" />
    <link rel="canonical" href="https://url.vet" />
    <meta
      property="og:url"
      content={currentUrl || `https://url.vet/?q=${encodeURIComponent(queryUrl)}`}
    />
    <meta property="og:image" content={ogImage} />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="url.vet — is {shareDomain} sus?" />
    <meta name="twitter:description" content={desc} />
    <meta name="twitter:image" content={ogImage} />
  {:else}
    <title>{liveTitle ?? "url.vet (URLvet) — sketchy link? just url.vet it"}</title>
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
    <meta property="og:image" content="https://url.vet/og" />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:image" content="https://url.vet/og" />
    <meta name="twitter:title" content="url.vet (URLvet) — just url.vet it" />
    <meta
      name="twitter:description"
      content="Sketchy link? url.vet it. Free, instant, no signup."
    />
    {@html `<script type="application/ld+json">${JSON.stringify(schemaSoftwareApp)}</script>`}
  {/if}
</svelte:head>
