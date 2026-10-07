<script lang="ts">
  import { env } from "$env/dynamic/public";
  import Icon from "$lib/components/Icon.svelte";
  import PageCta from "$lib/components/PageCta.svelte";
  import { CHECK_GROUPS, TOTAL_CHECKS } from "$lib/data/checks";
  import { FAQ_GROUPS, faqId, faqText } from "$lib/data/faq";
  import { REPO } from "$lib/site";
  import { ICON } from "$lib/ui/icons";
  import { LINK } from "$lib/ui/text";
  import { onMount } from "svelte";

  // Swagger UI is served by the API server itself, next to /api/v1.
  const API_ORIGIN = (() => {
    try {
      return new URL(env.PUBLIC_BASE_URL || "http://localhost:8080/api/v1").origin;
    } catch {
      return "http://localhost:8080";
    }
  })();
  // Self-hosting comes first; the public API is only for a quick try.
  const DEV_LINKS = [
    {
      label: "Run it on your own server",
      note: "Docker setup, configuration and a reverse proxy, step by step",
      href: `${REPO}/blob/main/docs/deployment.md`,
    },
    {
      label: "API reference",
      note: "The analyze endpoint, what it returns, and its errors",
      href: `${REPO}/blob/main/docs/api.md`,
    },
    {
      label: "Try the API in Swagger UI",
      note: "Send a few test requests to url.vet's server from your browser. Rate-limited, so not for production",
      href: `${API_ORIGIN}/swagger/index.html`,
    },
  ];

  // A made-up scan, worded exactly as the engine words its findings.
  const EXAMPLE = {
    domain: "paypa1-login.top",
    score: 7,
    red: [
      "Typosquatting detected: domain closely resembles 'paypal.com' (1 character difference).",
      "Registered 3 days ago. Phishing sites are usually brand new.",
      "High-risk domain extension detected (often associated with spam).",
      "Sensitive security keywords found in URL: login",
      "Very low traffic volume.",
    ],
    green: ["Enforces strict HTTPS security (HSTS Enabled)."],
  };

  const STEPS = [
    {
      title: "Paste the link",
      desc: "A link or a whole message with links in it.",
    },
    {
      title: "Our server visits it",
      desc: "Nothing from the page reaches your device.",
    },
    {
      title: "You get a verdict, and why",
      desc: "Safe, Sus or Risky, with every reason listed.",
    },
  ];

  const BANDS = [
    { v: "Risky", from: 0, to: 29, dot: "bg-red-500" },
    { v: "Suspicious", from: 30, to: 64, dot: "bg-yellow-500" },
    { v: "Safe", from: 65, to: 100, dot: "bg-emerald-500" },
  ];

  // What a result has besides the verdict and flags shown in the example.
  const PARTS = [
    {
      title: "Page preview",
      desc: "A screenshot of the page, taken by our server. If it looks like your bank but the address has nothing to do with your bank, that tells you a lot.",
    },
    {
      title: "Details",
      desc: "Grouped sections (domain, redirects, security certificate, page content, threat feeds) for anyone who wants the raw facts.",
    },
    {
      title: "Share link",
      desc: "Every result has a link you can send to whoever sent you the original. Its preview shows the verdict and score at a glance.",
    },
  ];

  // A link to one question (#faq-...) or one check group (#check-...) opens it.
  function openFromHash() {
    const target = /^#(faq|check)-/.test(location.hash)
      ? document.getElementById(location.hash.slice(1))
      : null;
    const details =
      target instanceof HTMLDetailsElement ? target : target?.querySelector("details");
    if (details) {
      details.open = true;
      target?.scrollIntoView({ block: target === details ? "center" : "start" });
    }
  }
  onMount(openFromHash);

  const schemaFAQ = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: FAQ_GROUPS.flatMap((g) => g.items).map((f) => ({
      "@type": "Question",
      name: f.q,
      acceptedAnswer: { "@type": "Answer", text: faqText(f) },
    })),
  };
</script>

<svelte:window on:hashchange={openFromHash} />

<svelte:head>
  <title>How url.vet checks a link</title>
  <meta
    name="description"
    content="The checks url.vet runs on every link, how they add up to a verdict, and answers to common questions."
  />
  <meta property="og:title" content="How url.vet checks a link" />
  <meta
    property="og:description"
    content="The checks url.vet runs on every link, how they add up to a verdict, and answers to common questions."
  />
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://url.vet/how-it-works" />
  <meta property="og:image" content="https://url.vet/og" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="How url.vet checks a link" />
  <meta
    name="twitter:description"
    content="The checks url.vet runs on every link, how they add up to a verdict, and answers to common questions."
  />
  <meta name="twitter:image" content="https://url.vet/og" />
  <link rel="canonical" href="https://url.vet/how-it-works" />
  {@html `<script type="application/ld+json">${JSON.stringify(schemaFAQ)}</script>`}
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">How it works</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    Paste a link. <span class="italic">See why.</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    A verdict you can't explain can't be trusted, so every one comes with its reasons. Here's what
    happens when you scan, and how to read what comes back.
  </p>

  <!-- The whole thing in three steps: a row on desktop, a short timeline on phones -->
  <ol class="mt-12 grid sm:grid-cols-3 sm:gap-x-6">
    {#each STEPS as step, n}
      <li class="relative flex sm:block gap-4 pb-7 last:pb-0 sm:pb-0">
        <!-- line to the next step -->
        {#if n < STEPS.length - 1}
          <span
            class="absolute bg-gray-200 dark:bg-gray-800 left-[17px] top-11 bottom-1 w-px sm:left-12 sm:-right-3 sm:top-[17px] sm:bottom-auto sm:w-auto sm:h-px"
            aria-hidden="true"
          ></span>
        {/if}
        <span
          class="relative flex-shrink-0 flex items-center justify-center w-9 h-9 rounded-full border border-gray-300 dark:border-gray-700 font-mono text-[13px] text-gray-700 dark:text-gray-300"
          aria-hidden="true">{n + 1}</span
        >
        <div class="sm:mt-4 sm:pr-4">
          <p class="font-serif text-[1.35rem] leading-tight">{step.title}</p>
          <p class="mt-1.5 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {step.desc}
          </p>
        </div>
      </li>
    {/each}
  </ol>

  <!-- Example result -->
  <section class="mt-20 scroll-mt-20" id="result">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">What you get back</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      url.vet checks the page itself, live, not just a list of known bad links. A few seconds after
      you press Scan, you get something like this:
    </p>

    <figure
      class="mt-8 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-5 sm:p-6"
      aria-label="Example result"
    >
      <div class="flex flex-wrap items-center justify-between gap-3">
        <span
          class="px-3 py-1 rounded-full border border-gray-300 dark:border-gray-800 font-mono text-[13px] text-gray-800 dark:text-gray-200"
          >{EXAMPLE.domain}</span
        >
        <span class="flex items-baseline gap-3">
          <span
            class="inline-flex items-center gap-1.5 text-sm font-medium text-red-600 dark:text-red-400"
            ><span class="w-1.5 h-1.5 rounded-full bg-red-500"></span>Risky</span
          >
          <span class="font-serif text-3xl leading-none"
            >{EXAMPLE.score}<span class="text-base text-gray-400">/100</span></span
          >
        </span>
      </div>

      {#each [{ label: "Red flags", items: EXAMPLE.red, dot: "bg-red-500" }, { label: "Green flags", items: EXAMPLE.green, dot: "bg-emerald-500" }] as group}
        <p class="mt-6 font-mono text-[11px] uppercase tracking-wider text-gray-500">
          {group.label}
        </p>
        <ul class="mt-2 divide-y divide-gray-100 dark:divide-gray-800">
          {#each group.items as flag}
            <li class="flex gap-3 py-2.5 text-[15px] leading-snug text-gray-700 dark:text-gray-300">
              <span class="mt-[0.45em] w-1.5 h-1.5 flex-shrink-0 rounded-full {group.dot}"></span>
              <span class="flex-1">{flag}</span>
            </li>
          {/each}
        </ul>
      {/each}
      <figcaption class="mt-4 font-mono text-[11px] text-gray-400">
        Example result, simplified. Not a real site.
      </figcaption>
    </figure>

    <dl class="mt-10 grid sm:grid-cols-3 gap-x-8 gap-y-7">
      {#each PARTS as part}
        <div class="pt-4 border-t border-gray-200 dark:border-gray-800">
          <dt class="font-serif text-[1.35rem] leading-tight">{part.title}</dt>
          <dd class="mt-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {part.desc}
          </dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- Checks: what each group looks at, and which side of the score a finding counts towards -->
  <section class="mt-24 scroll-mt-20" id="checks">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">The checks</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      {TOTAL_CHECKS} checks start at the same moment you press Scan. Each one is independent, so a slow
      or failed check never holds up the rest. Together they answer seven questions, and every finding
      counts towards either the risk or the trust side of the score.
    </p>
    <p class="mt-3 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
      A few name-based rules are relaxed for well-known sites and verified endings like .gov, so a
      real bank's own login page isn't flagged. Curious what gets sent where during a scan? See the <a
        href="/privacy"
        class={LINK}>privacy page</a
      >.
    </p>

    <div class="mt-8">
      {#each CHECK_GROUPS as group, n}
        <article
          id="check-{group.id}"
          class="scroll-mt-20 py-7 border-t border-gray-200 dark:border-gray-800"
        >
          <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">
            {String(n + 1).padStart(2, "0")} · {group.label} · {group.checks}
            {group.checks === 1 ? "check" : "checks"}
          </p>
          <h3 class="mt-2 font-serif text-[1.6rem] leading-tight">{group.title}</h3>
          <p class="mt-2 text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">
            {group.desc}
          </p>

          <details class="group mt-3">
            <summary
              class="inline-flex items-center gap-1.5 cursor-pointer list-none [&::-webkit-details-marker]:hidden text-[14px] text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 rounded"
            >
              <span
                class="underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700"
                >See what it looks for ({group.signals.length})</span
              >
              <span class="transition-transform duration-200 group-open:rotate-180">
                <Icon path={ICON.chevronDown} class="w-3.5 h-3.5" />
              </span>
            </summary>
            <ul class="mt-4 grid sm:grid-cols-2 gap-x-10">
              {#each group.signals as signal}
                <li
                  class="flex items-baseline justify-between gap-4 py-2 border-b border-gray-100 dark:border-gray-800/70 text-[14px]"
                >
                  <span class="flex items-baseline gap-2.5 text-gray-700 dark:text-gray-300">
                    <span
                      class="relative top-[-1px] w-1.5 h-1.5 rounded-full flex-shrink-0 {signal.kind ===
                      'risk'
                        ? 'bg-red-500'
                        : signal.kind === 'trust'
                          ? 'bg-emerald-500'
                          : 'bg-gray-300 dark:bg-gray-600'}"
                      aria-hidden="true"
                    ></span>
                    {signal.text}
                  </span>
                  <span
                    class="flex-shrink-0 font-mono text-[11px] whitespace-nowrap {signal.kind ===
                    'risk'
                      ? 'text-red-600 dark:text-red-400'
                      : signal.kind === 'trust'
                        ? 'text-emerald-700 dark:text-emerald-400'
                        : 'text-gray-400 dark:text-gray-500'}"
                  >
                    {signal.kind === "risk"
                      ? "adds risk"
                      : signal.kind === "trust"
                        ? "adds trust"
                        : "shown only"}
                  </span>
                </li>
              {/each}
            </ul>

            {#if group.against}
              <p class="mt-3 text-[13px] text-gray-500">Compared against {group.against}</p>
            {/if}
          </details>
        </article>
      {/each}
    </div>
  </section>

  <!-- Scoring: the idea in words; no point values, since those get tuned -->
  <section class="mt-24 scroll-mt-20" id="score">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">How the score works</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        Every finding counts towards one of two sides: risk or trust. The score starts in the
        middle, at 50. Risk pulls it down and trust pulls it up, and it always stays between 0 and
        100.
      </p>
      <p>
        Serious findings weigh more than small ones, so one confirmed phishing report outweighs a
        handful of minor good signs. The number then falls into one of three bands. The example
        above landed at {EXAMPLE.score}.
      </p>
    </div>

    <!-- scale -->
    <div class="mt-10 mb-2">
      <div class="relative">
        <div class="flex h-2 rounded-full overflow-hidden gap-0.5">
          {#each BANDS as band}
            <div class={band.dot} style="width: {band.to - band.from + 1}%"></div>
          {/each}
        </div>
        <div
          class="absolute -top-6 flex flex-col items-center"
          style="left: {EXAMPLE.score}%; transform: translateX(-50%)"
        >
          <span class="font-mono text-[10px] text-gray-500 leading-none">{EXAMPLE.score}</span>
          <span class="mt-1 w-0.5 h-6 rounded-full bg-gray-900 dark:bg-gray-100"></span>
        </div>
      </div>
      <div class="mt-4 grid grid-cols-3 gap-4">
        {#each BANDS as band}
          <div>
            <p class="flex items-center gap-2 text-[15px] font-medium">
              <span class="w-1.5 h-1.5 rounded-full {band.dot}"></span>{band.v}
            </p>
            <p class="mt-0.5 font-mono text-xs text-gray-500">{band.from}–{band.to}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- FAQ: collapsed, one open at a time; every answer stays in the page for search engines -->
  <section class="mt-24 scroll-mt-20" id="faq">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Frequently asked questions</h2>
    <div class="mt-8">
      {#each FAQ_GROUPS as group}
        <div class="border-t border-gray-200 dark:border-gray-800">
          <h3 class="pt-6 font-mono text-[11px] uppercase tracking-wider text-gray-500">
            {group.title}
          </h3>
          <div class="mt-2 divide-y divide-gray-200 dark:divide-gray-800">
            {#each group.items as item}
              <details id={faqId(item.q)} name="faq" class="group scroll-mt-24">
                <summary
                  class="flex items-start justify-between gap-4 py-4 cursor-pointer list-none [&::-webkit-details-marker]:hidden text-gray-900 dark:text-gray-100 hover:text-gray-600 dark:hover:text-gray-300 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 rounded"
                >
                  <span class="font-serif text-[1.3rem] leading-snug">{item.q}</span>
                  <span
                    class="mt-1.5 flex-shrink-0 text-gray-400 transition-transform duration-200 group-open:rotate-180"
                  >
                    <Icon path={ICON.chevronDown} />
                  </span>
                </summary>
                <div class="pb-5 pr-8 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
                  <p>{item.a}</p>
                  {#if item.points}
                    <ul class="mt-3 space-y-2">
                      {#each item.points as point}
                        <li class="flex gap-2.5">
                          <span
                            class="mt-[0.6em] w-1 h-1 rounded-full bg-gray-400 flex-shrink-0"
                            aria-hidden="true"
                          ></span>
                          <span>{point}</span>
                        </li>
                      {/each}
                    </ul>
                  {/if}
                </div>
              </details>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  </section>

  <!-- For developers -->
  <section class="mt-24 scroll-mt-20" id="developers">
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">For developers</p>
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">
      Self-hosting and API docs
    </h2>
    <p class="mt-4 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
      url.vet is <a href={REPO} class={LINK} target="_blank" rel="noopener noreferrer"
        >open source</a
      >. You can run the whole thing on your server, or use its API in your app.
    </p>
    <ul class="mt-6">
      {#each DEV_LINKS as link}
        <li class="border-t border-gray-200 dark:border-gray-800">
          <a
            href={link.href}
            target="_blank"
            rel="noopener noreferrer"
            class="group flex items-center justify-between gap-4 py-4"
          >
            <span>
              <span class="block font-serif text-[1.35rem] leading-tight">{link.label}</span>
              <span class="mt-1 block text-[15px] text-gray-600 dark:text-gray-400"
                >{link.note}</span
              >
            </span>
            <span
              class="text-gray-400 group-hover:text-gray-900 dark:group-hover:text-gray-100 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform"
              aria-hidden="true">↗</span
            >
          </a>
        </li>
      {/each}
    </ul>
  </section>

  <PageCta />
</div>
