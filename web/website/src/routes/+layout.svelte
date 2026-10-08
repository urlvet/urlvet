<script lang="ts">
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import { page } from "$app/stores";
  import Guide from "$lib/components/guide/Guide.svelte";
  import { vettyMemory } from "$lib/components/guide/store";
  import { theme } from "$lib/theme";
  import { REPO } from "$lib/site";
  import { GITHUB_LOGO } from "$lib/ui/icons";
  import { onMount } from "svelte";
  import "../app.css";

  const year = new Date().getFullYear();

  const FOOTER_LINKS = [
    { label: "How it works", href: "/how-it-works" },
    { label: "About", href: "/about" },
    { label: "Privacy", href: "/privacy" },
    { label: "GitHub", href: REPO, external: true },
  ];

  $: isHome = $page.url.pathname === "/";
  $: isAdmin = $page.url.pathname.startsWith("/admin");

  const schemaWebSite = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: "url.vet",
    alternateName: ["URLvet", "urlvet", "url vet"],
    url: "https://url.vet",
    description:
      "Free real-time URL scanner. Paste any link and get an instant phishing and safety verdict, no signup needed.",
    potentialAction: {
      "@type": "SearchAction",
      target: {
        "@type": "EntryPoint",
        urlTemplate: "https://url.vet/?q={search_term_string}",
      },
      "query-input": "required name=search_term_string",
    },
  };

  const schemaOrganization = {
    "@context": "https://schema.org",
    "@type": "Organization",
    name: "url.vet",
    alternateName: ["URLvet", "urlvet", "url vet"],
    url: "https://url.vet",
    logo: "https://url.vet/favicon.ico",
    sameAs: [REPO],
  };

  onMount(() => {
    theme.init();
  });
</script>

<svelte:head>
  {@html `<script type="application/ld+json">${JSON.stringify(schemaWebSite)}</script>`}
  {@html `<script type="application/ld+json">${JSON.stringify(schemaOrganization)}</script>`}
</svelte:head>

<div class="relative min-h-screen flex flex-col">
  <nav class="relative z-20 w-full max-w-5xl mx-auto px-6 pt-6 flex items-center justify-between">
    <!-- Home link: the wordmark. Hidden on the home page, which shows the big one. -->
    <a
      href="/"
      on:click={() => (location.href = "/")}
      class:invisible={isHome}
      aria-label="url.vet home"
      class="group font-serif text-[1.65rem] leading-none tracking-[-0.02em]"
    >
      <span
        class="text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400 transition-colors"
        >url</span
      ><span class="-mx-[0.025em] text-accent-light dark:text-accent-dark">.</span><span
        class="text-gray-900 dark:text-gray-100">vet</span
      >
    </a>
    <div class="flex items-center gap-2">
      {#if !isAdmin}
        <a
          href={REPO}
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-gray-300 dark:border-gray-800 text-sm text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path d={GITHUB_LOGO} />
          </svg>
          <span class="hidden sm:inline">GitHub</span>
        </a>
      {/if}
      <ThemeToggle />
    </div>
  </nav>
  <main class="flex-1">
    <slot />
  </main>

  {#if !isAdmin}
    <Guide />
  {/if}

  <footer class="text-gray-500 dark:text-gray-400 py-8">
    <div
      class="max-w-5xl mx-auto px-6 pt-8 border-t border-gray-200 dark:border-gray-800 flex flex-col md:flex-row md:justify-between items-center gap-4 text-sm"
    >
      <!-- Left: site links, then GitHub, then (if hidden) Vetty -->
      <nav
        aria-label="Footer"
        data-guide="learn-more"
        class="flex flex-wrap items-center justify-center md:justify-start gap-x-5 sm:gap-x-4 gap-y-2"
      >
        {#each FOOTER_LINKS as link, i}
          {#if i > 0}<span
              class="hidden sm:inline text-gray-300 dark:text-gray-700"
              aria-hidden="true">·</span
            >{/if}
          <a
            href={link.href}
            target={link.external ? "_blank" : undefined}
            rel={link.external ? "noopener noreferrer" : undefined}
            aria-current={$page.url.pathname === link.href ? "page" : undefined}
            class="transition-colors hover:text-gray-900 dark:hover:text-gray-100 {$page.url
              .pathname === link.href
              ? 'text-gray-900 dark:text-gray-100'
              : 'text-gray-500 dark:text-gray-400'}"
          >
            {link.label}
          </a>
        {/each}
        {#if $vettyMemory.hidden}
          <span class="hidden sm:inline text-gray-300 dark:text-gray-700" aria-hidden="true">·</span
          >
          <button
            type="button"
            on:click={() => vettyMemory.set({ hidden: false })}
            class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          >
            Show Vetty 📎
          </button>
        {/if}
      </nav>

      <!-- Right: License and author -->
      <p class="text-gray-500 dark:text-gray-400 text-center md:text-right">
        <a
          href="{REPO}/blob/main/LICENSE"
          target="_blank"
          class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          >AGPL-3.0</a
        >
        © 2021–{year}
        <a
          href="https://abhizaik.com"
          target="_blank"
          class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          >abhizaik</a
        >
      </p>
    </div>
  </footer>
</div>
