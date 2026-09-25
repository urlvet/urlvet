<script lang="ts">
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import { page } from "$app/stores";
  import Guide from "$lib/components/guide/Guide.svelte";
  import { vettyMemory } from "$lib/components/guide/store";
  import { theme } from "$lib/theme";
  import { onMount } from "svelte";
  import "../app.css";

  const year = new Date().getFullYear();

  const FOOTER_LINKS = [
    { label: "How it works", href: "/how-it-works" },
    { label: "About", href: "/about" },
    { label: "Privacy", href: "/privacy" },
    { label: "GitHub", href: "https://github.com/urlvet/urlvet", external: true },
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
    sameAs: ["https://github.com/urlvet/urlvet"],
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
          href="https://github.com/urlvet/urlvet"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-gray-300 dark:border-gray-800 text-sm text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path
              d="M12 .5C5.648.5.5 5.648.5 12c0 5.084 3.292 9.387 7.872 10.905.576.106.784-.25.784-.556 0-.274-.01-1-.016-1.974-3.2.694-3.876-1.544-3.876-1.544-.522-1.33-1.276-1.683-1.276-1.683-1.042-.714.08-.699.08-.699 1.152.08 1.756 1.183 1.756 1.183 1.024 1.754 2.688 1.246 3.344.952.104-.742.4-1.246.728-1.532-2.554-.292-5.238-1.276-5.238-5.674 0-1.252.448-2.274 1.184-3.076-.12-.292-.512-1.468.112-3.064 0 0 .964-.308 3.16 1.176.916-.254 1.9-.382 2.876-.388.976.006 1.96.134 2.876.388 2.192-1.484 3.156-1.176 3.156-1.176.624 1.596.232 2.772.112 3.064.736.802 1.184 1.824 1.184 3.076 0 4.408-2.69 5.38-5.25 5.664.412.354.78 1.05.78 2.118 0 1.532-.014 2.768-.014 3.144 0 .308.2.672.788.556C20.708 21.385 24 17.084 24 12c0-6.352-5.148-11.5-12-11.5z"
            />
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
          href="https://github.com/urlvet/urlvet/blob/main/LICENSE"
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
