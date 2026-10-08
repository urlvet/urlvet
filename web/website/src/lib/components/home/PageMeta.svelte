<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy } from "svelte";
  import { REPO } from "../../site";
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

  const SITE_TITLE = "url.vet (URLvet): Is this link safe?";
  const SITE_DESC =
    "Paste any link to see if it's safe, suspicious or risky, and why. Free, no signup.";

  const dotIcon = (color: string) =>
    `data:image/svg+xml,${encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="7" fill="${color}"/></svg>`)}`;

  const schemaSoftwareApp = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "url.vet",
    alternateName: ["URLvet", "urlvet", "url vet"],
    url: "https://url.vet",
    description:
      "A free link checker that spots phishing, lookalike addresses and other risky links, and shows the reasons behind every verdict.",
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
    license: `${REPO}/blob/main/LICENSE`,
    codeRepository: REPO,
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
    {@const title = `Is ${shareDomain} safe? | url.vet`}
    {@const desc = verdict
      ? `url.vet rated ${shareDomain} ${verdict}. See why, or check another link for free.`
      : `See if ${shareDomain} is safe to open, and why. Free, no signup.`}
    {@const ogImage = `https://url.vet/og?domain=${encodeURIComponent(shareDomain)}${verdict ? `&v=${encodeURIComponent(verdict)}` : ""}${score !== undefined ? `&s=${score}` : ""}`}
    <title>{liveTitle ?? title}</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content={title} />
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
    <meta name="twitter:title" content={title} />
    <meta name="twitter:description" content={desc} />
    <meta name="twitter:image" content={ogImage} />
  {:else}
    <title>{liveTitle ?? SITE_TITLE}</title>
    <meta name="description" content={SITE_DESC} />
    <meta property="og:title" content={SITE_TITLE} />
    <meta property="og:description" content={SITE_DESC} />
    <link rel="canonical" href="https://url.vet" />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://url.vet" />
    <meta property="og:image" content="https://url.vet/og" />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:image" content="https://url.vet/og" />
    <meta name="twitter:title" content={SITE_TITLE} />
    <meta name="twitter:description" content={SITE_DESC} />
    {@html `<script type="application/ld+json">${JSON.stringify(schemaSoftwareApp)}</script>`}
  {/if}
</svelte:head>
