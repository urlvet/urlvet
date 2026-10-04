<script lang="ts">
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { LINK } from "$lib/ui/text";

  // Plain-language privacy notice. Keep it in sync with what the server actually does:
  // logging (server/internal/logger, handler/middleware/request_log.go), rate limiting
  // (handler/middleware/rate_limit.go), caching (server/internal/constants/ttls.go),
  // reports (handler/report.go).

  const AT_A_GLANCE = [
    { label: "Accounts", value: "None" },
    { label: "Cookies", value: "None" },
    { label: "Analytics & ads", value: "None" },
    { label: "Data sold", value: "Never" },
    { label: "Your IP in our logs", value: "Never" },
    { label: "Cached scan results", value: "Up to 24 hours" },
  ];

  const NEVER = ["Sell your data", "Show you ads", "Ask you to sign up", "Build a profile of you"];

  const KEEP = [
    {
      t: "Scan results",
      d: "Cached by link for up to 24 hours, so a link that's checked twice (or shared) loads instantly. Individual checks expire sooner, some after 3 hours.",
    },
    {
      t: "Page screenshots",
      d: "Stored on our server for up to 24 hours, then deleted automatically.",
    },
    {
      t: "Recent activity",
      d: "The last 100 scans (link, verdict, score, time) are held in the server's memory for the maintainer's dashboard. Never written to disk; gone whenever the server restarts.",
    },
    {
      t: "Reports you send",
      d: "If you report a wrong verdict, we keep the link, our verdict and score, your suggestion, your comment and the time. It's kept until the issue it points to is fixed, then deleted. No IP address, no account.",
    },
  ];

  const THIRD_PARTIES = [
    {
      who: "Vercel, our host",
      what: "Hosts this website. Like any host, it handles your visit, so it sees your IP address and the address of the page you open (for a share link, that includes the scanned link). It keeps its request logs briefly.",
      href: "https://vercel.com/legal/privacy-policy",
    },
    {
      who: "The website you scan",
      what: "Our server visits it to check redirects, certificates and content, and takes a screenshot. The site, and any trackers on it, see our server, not you.",
    },
    {
      who: "PhishTank",
      what: "Receives the link, to check it against its database of reported phishing.",
    },
    {
      who: "DNS resolvers",
      what: "Receive the domain name only, to find the servers behind it (as any lookup on the internet does).",
    },
    {
      who: "Domain registries",
      what: "The public record of who registered a domain. They receive the domain name only, so we can look up its age.",
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
    url.vet needs the link you paste to check it. That's it. Here is exactly what happens to it, in
    plain words. And because the code is open source, you don't have to take our word for any of
    this.
  </p>

  <!-- At a glance -->
  <dl class="mt-12 grid sm:grid-cols-2 gap-x-10">
    {#each AT_A_GLANCE as item}
      <div
        class="flex items-baseline justify-between gap-4 py-3.5 border-b border-gray-200 dark:border-gray-800"
      >
        <dt class="text-[15px] text-gray-600 dark:text-gray-400">{item.label}</dt>
        <dd class="font-mono text-sm text-gray-900 dark:text-gray-100">{item.value}</dd>
      </div>
    {/each}
  </dl>

  <!-- Promises -->
  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">What we'll never do</h2>
    <ul class="mt-6 grid sm:grid-cols-2 gap-x-10">
      {#each NEVER as item}
        <li
          class="py-3 border-t border-gray-200 dark:border-gray-800 font-serif text-[1.45rem] leading-tight"
        >
          {item}
        </li>
      {/each}
    </ul>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">When you scan a link</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        The link goes to our server, which runs the checks and sends back the result. Like any
        website, our server can see your IP address while it answers you. We use it for one thing:
        limiting how many requests one address can make per minute, so the service stays up. Even
        then we don't store the address itself, only a scrambled version that can't be turned back
        into it, and the key used to scramble it is replaced every day. That counter is deleted
        after a minute.
      </p>
      <p>
        <span class="text-gray-900 dark:text-gray-100 font-medium"
          >Our logs never contain your IP address or the link you scanned.</span
        >
        They record one line per request (which page, whether it worked, how long it took), and if an
        error mentions a web address, everything after the
        <span class="font-mono text-[15px]">?</span> is removed before it's written, since that's where
        links hide things like password-reset tokens.
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
          <dt class="sm:w-44 flex-shrink-0 font-serif text-[1.3rem] leading-tight">{item.t}</dt>
          <dd class="text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">{item.d}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Who else sees the link</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      This website runs on Vercel, and some checks need outside help. These are the only parties
      involved:
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
      Everything this website loads (fonts, images, code) comes from url.vet itself. No trackers, no
      ad networks, no font or analytics services watching your visit.
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
      We don't set cookies. Your browser's local storage remembers your light or dark theme and
      whether you've met or hidden Vetty. It never leaves your device, and clearing your site data
      removes it.
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
      >. There are no accounts, so we hold nothing tied to you by name. To ask what we hold about a
      link you scanned or reported, or to have a report deleted, email
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
