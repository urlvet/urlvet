<script lang="ts">
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { LINK } from "$lib/ui/text";

  // Plain-language privacy notice. Keep it in sync with what the server actually does:
  // logging (server/internal/logger, handler/middleware/request_log.go), rate limiting
  // (handler/middleware/rate_limit.go), caching (server/internal/constants/ttls.go),
  // reports (handler/report.go).

  // The short version: what we don't have, and the few things we keep for a while.
  const NONE = [
    "accounts",
    "cookies",
    "ads",
    "trackers",
    "profiles of you",
    "IP addresses in our logs",
    "data sold",
  ];
  // What we keep, how long for, and what's in it.
  const KEEP = [
    {
      t: "Scan results",
      for: "24 hours",
      d: "Cached by link, so a repeat or shared check loads instantly. Some individual checks expire sooner.",
    },
    {
      t: "Page screenshots",
      for: "24 hours",
      d: "Stored on our server, then deleted automatically.",
    },
    {
      t: "Recent searches",
      for: "Until restart",
      d: "The last 100 scans (link, verdict, score, time), held in the server's memory for the admin dashboard. Never written to disk.",
    },
    {
      t: "Reports you send",
      for: "Until fixed",
      d: "The link, our verdict and score, your suggestion, your comment and the time. No IP address, no account.",
    },
  ];

  const THIRD_PARTIES = [
    {
      who: "Vercel, our host",
      what: "Hosts this website, so it sees your IP address and the page you open (for a share link, that includes the scanned link). It keeps request logs briefly.",
      href: "https://vercel.com/legal/privacy-policy",
    },
    {
      who: "The website you scan",
      what: "Our server visits it and takes a screenshot. The site, and any trackers on it, see our server, not you.",
    },
    {
      who: "PhishTank",
      what: "Receives the link, to check it against reported phishing.",
    },
    {
      who: "DNS resolvers",
      what: "Receive the domain name only, to find the servers behind it.",
    },
    {
      who: "Domain registries",
      what: "Receive the domain name only, so we can look up its age.",
    },
  ];
</script>

<svelte:head>
  <title>Privacy — url.vet (URLvet)</title>
  <meta
    name="description"
    content="What url.vet does with the links you scan: no accounts, no cookies, no analytics, no data sold, and no IP addresses in our logs."
  />
  <link rel="canonical" href="https://url.vet/privacy" />
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">Privacy</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    We check the link, <span class="italic">not you.</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    url.vet needs the link you paste, and nothing else. Here is exactly what happens to it.
  </p>

  <!-- At a glance: what we don't have, as one statement; what we keep, as three figures -->
  <section class="mt-14" aria-label="The short version">
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">In short</p>
    <p class="mt-4 font-serif text-[1.75rem] md:text-[2.1rem] leading-[1.3] tracking-[-0.01em]">
      {#each NONE as item}
        <span class="whitespace-nowrap"
          ><span class="text-gray-400 dark:text-gray-500">No</span>
          {item}.</span
        >{" "}
      {/each}
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">When you scan a link</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        The link goes to our server, which runs the checks and sends back the result. Our server
        sees your IP address while it answers, and uses it for one thing: limiting how many requests
        one address can make per minute. It keeps only a scrambled version that can't be turned
        back, and deletes it after a minute.
      </p>
      <p>
        <span class="text-gray-900 dark:text-gray-100 font-medium"
          >Our logs never contain your IP address or the link you scanned.</span
        >
        They record one line per request: which page, whether it worked, how long it took. If an error
        mentions a web address, everything after the
        <span class="font-mono text-[15px]">?</span> is cut before it's written.
      </p>
    </div>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">
      What we keep, and for how long
    </h2>
    <dl class="mt-6">
      {#each KEEP as item}
        <div
          class="flex flex-col sm:flex-row gap-1 sm:gap-6 py-5 border-t border-gray-200 dark:border-gray-800"
        >
          <dt class="sm:w-44 flex-shrink-0">
            <span class="block font-serif text-[1.3rem] leading-tight">{item.t}</span>
            <span class="mt-1 block font-mono text-[11px] uppercase tracking-wider text-gray-500"
              >{item.for}</span
            >
          </dt>
          <dd class="text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">{item.d}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Who else sees the link</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Some checks need outside help. These are the only others involved:
    </p>
    <dl class="mt-6">
      {#each THIRD_PARTIES as p}
        <div
          class="flex flex-col sm:flex-row gap-1 sm:gap-6 py-5 border-t border-gray-200 dark:border-gray-800"
        >
          <dt class="sm:w-44 flex-shrink-0 font-serif text-[1.3rem] leading-tight">{p.who}</dt>
          <dd class="text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">
            {p.what}
            {#if p.href}
              <a href={p.href} target="_blank" rel="noopener noreferrer" class={LINK}
                >See Vercel's privacy policy</a
              >.
            {/if}
          </dd>
        </div>
      {/each}
    </dl>
    <p class="mt-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Everything this website loads (fonts, images, code) comes from url.vet itself.
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Sharing a result</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      A share link contains the scanned link, the verdict and the score, so anyone you send it to
      (and any app that shows a preview of it) can see those. Nothing about you is in it.
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">In your browser</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Your browser remembers your light or dark theme and whether you've met or hidden Vetty. That
      stays on your device, and clearing your site data removes it.
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Your rights</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      url.vet is run by <a
        href="https://abhizaik.com"
        target="_blank"
        rel="noopener noreferrer"
        class={LINK}>abhizaik</a
      >. We hold nothing tied to your name. To ask what we hold about a link you scanned or
      reported, or to have a report deleted, email
      <a href="mailto:hi@url.vet" class="font-mono text-base {LINK}">hi@url.vet</a>.
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Verify it</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Everything above is in the <a
        href={REPO}
        target="_blank"
        rel="noopener noreferrer"
        class={LINK}>source code</a
      >, and you can run your own copy if you'd rather not use ours. Every change to this page is
      public too, in its
      <a
        href="{REPO}/commits/main/web/website/src/routes/privacy/+page.svelte"
        target="_blank"
        rel="noopener noreferrer"
        class={LINK}>history on GitHub</a
      >.
    </p>
  </section>

  <PageCta />
</div>
