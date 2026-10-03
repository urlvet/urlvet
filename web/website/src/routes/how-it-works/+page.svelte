<script lang="ts">
  import { env } from "$env/dynamic/public";
  import { PILL_SOLID } from "$lib/ui/buttons";

  const LINK =
    "text-gray-900 dark:text-gray-100 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:decoration-current transition-colors";

  const REPO = "https://github.com/urlvet/urlvet";
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
      label: "Try our API first",
      note: "Send a few test requests to url.vet's server. Rate-limited, so not for production",
      href: `${API_ORIGIN}/swagger/index.html`,
    },
  ];

  // A made-up scan, worded exactly as the engine words its findings, with the
  // points each finding is worth in the real scoring (analyzer/result.go).
  const EXAMPLE = {
    domain: "paypa1-login.top",
    red: [
      {
        text: "Typosquatting detected: domain closely resembles 'paypal.com' (1 character difference).",
        pts: 40,
      },
      { text: "Newly created domain (3 days old). High Risk.", pts: 25 },
      { text: "High-risk domain extension detected (often associated with spam).", pts: 20 },
      { text: "Sensitive security keywords found in URL: login", pts: 10 },
      { text: "Very low traffic volume.", pts: 10 },
    ],
    green: [{ text: "Enforces strict HTTPS security (HSTS Enabled).", pts: 20 }],
  };
  const risk = EXAMPLE.red.reduce((n, f) => n + f.pts, 0);
  const trust = EXAMPLE.green.reduce((n, f) => n + f.pts, 0);
  // Same formula as the backend: Go's math.Round rounds halves away from zero.
  const half = (trust - risk) / 2;
  const shift = Math.sign(half) * Math.round(Math.abs(half));
  const score = Math.max(0, Math.min(100, 50 + shift));

  const BANDS = [
    { v: "Risky", from: 0, to: 29, dot: "bg-red-500" },
    { v: "Suspicious", from: 30, to: 64, dot: "bg-yellow-500" },
    { v: "Safe", from: 65, to: 100, dot: "bg-emerald-500" },
  ];

  const PARTS = [
    {
      title: "Score and verdict",
      desc: "A number from 0 to 100, higher is safer, and the band it falls in: Risky, Suspicious or Safe.",
    },
    {
      title: "Red and green flags",
      desc: "Every finding that moved the score, in plain words. Red flags add risk, green flags add trust.",
    },
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

  const CHECKS = [
    {
      id: "url",
      label: "URL structure",
      desc: "Inspects the link before making any network request. Checks for IP addresses used as hostnames, URL shorteners, suspicious keywords in the path, lookalike characters from other alphabets, and unusually deep subdomains.",
    },
    {
      id: "network",
      label: "HTTP / Network",
      desc: "Makes one real request and follows every redirect. Checks HSTS, status code, and whether the final destination differs from the link you pasted.",
    },
    {
      id: "dns",
      label: "DNS",
      desc: "Verifies NS and MX records exist and that the domain resolves to a real IP.",
    },
    {
      id: "tls",
      label: "TLS / SSL",
      desc: "Checks certificate validity, expiry, issuer, Certificate Transparency log inclusion, and known-bad fingerprints.",
    },
    {
      id: "domain",
      label: "Domain intelligence",
      desc: "Looks up domain age via WHOIS, global traffic rank, TLD classification, DNSSEC status, character randomness in the domain name, and lookalike spellings of the 500 most-visited sites.",
    },
    {
      id: "content",
      label: "Content analysis",
      desc: "Fetches and parses the page. Detects login and payment forms on suspicious domains, hidden iframes, brand impersonation, and forms that submit data to external servers.",
    },
    {
      id: "threats",
      label: "Threat intelligence",
      desc: "Checks the URL against PhishTank's databases of confirmed and reported phishing links.",
    },
  ];

  // Worded the way people search, since these also feed the FAQ rich result.
  const FAQ = [
    {
      q: "How do I check if a link is safe?",
      a: "Paste it into the box on the url.vet home page and press Scan. A few seconds later you get a verdict (Safe, Suspicious or Risky), a score from 0 to 100, and the reasons behind it. You never have to open the link yourself.",
    },
    {
      q: "Is it safe to check a dangerous link here?",
      a: "Yes. Our server visits the link, not your device. Nothing from the page runs in your browser; you only see a screenshot.",
    },
    {
      q: "Will the website know I checked it?",
      a: "It sees a visit from our server, not from you. Your IP address is never sent to it.",
    },
    {
      q: "How can I tell if a link is a phishing scam?",
      a: "Common signs: the address misspells a real brand (paypa1.com), swaps in lookalike letters from another alphabet, was registered only days ago, uses an unusual ending like .top or .zip, or the page asks for a password or card details on a site the brand doesn't own. url.vet checks for all of these and lists any it finds.",
    },
    {
      q: "Can I check a link from WhatsApp, SMS, email or a QR code?",
      a: "Yes. Copy the link and paste it here. You can paste the whole message too; url.vet finds the links in it.",
    },
    {
      q: "Can it tell where a shortened link really goes?",
      a: "Yes. url.vet follows every redirect and shows the full path, from the link you were sent to the page you would land on, and flags it when the link jumps to a different site.",
    },
    {
      q: "I already clicked a suspicious link. What should I do?",
      a: "Don't enter anything on the page. If you already typed a password, change it on the real site and turn on two-step verification. If you entered card or bank details, call your bank right away. Then scan the link here to see what it was.",
    },
    {
      q: "What does Suspicious mean?",
      a: "The score landed between 30 and 64: there are some warning signs, or not enough evidence either way. That's common for new or little-known sites. Read the flags, and don't enter personal details unless you're sure.",
    },
    {
      q: "Why did a site I trust come back Suspicious?",
      a: "New or little-known sites have fewer signals in their favour, so they can score lower. The flags show exactly why. If you think it's wrong, use the Report button on the result. A human reads every report.",
    },
    {
      q: "Does Safe mean the site is honest?",
      a: "It means no warning signs turned up. A real shop can still sell bad products, and a brand-new scam may not show any signs yet. Use url.vet as one layer of defense, not the only one.",
    },
    {
      q: "Do you keep the links I check?",
      a: "Results are cached for up to 24 hours so a repeat check loads instantly, then deleted. We don't log who checked what, and there are no accounts. The privacy page has the details.",
    },
    {
      q: "Is url.vet free? Do I need an account?",
      a: "It's free, with no account, no signup and no ads.",
    },
  ];

  const schemaFAQ = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: FAQ.map((f) => ({
      "@type": "Question",
      name: f.q,
      acceptedAnswer: { "@type": "Answer", text: f.a },
    })),
  };
</script>

<svelte:head>
  <title>How It Works: Check If a Link Is Safe — url.vet (URLvet)</title>
  <meta
    name="description"
    content="How url.vet checks whether a link is safe or a phishing scam: the checks it runs, how the score is worked out, and answers to common questions."
  />
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

  <!-- Example result -->
  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">What you get back</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Paste a full link, a shortened one, or just a domain like example.com, and press Scan. url.vet
      checks the page itself, live, instead of looking it up in a blocklist or static database. A
      few seconds later you get something like this:
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
            >{score}<span class="text-base text-gray-400">/100</span></span
          >
        </span>
      </div>

      {#each [{ label: "Red flags", items: EXAMPLE.red, dot: "bg-red-500", sign: "+", kind: "risk" }, { label: "Green flags", items: EXAMPLE.green, dot: "bg-emerald-500", sign: "+", kind: "trust" }] as group}
        <p class="mt-6 font-mono text-[11px] uppercase tracking-wider text-gray-500">
          {group.label}
        </p>
        <ul class="mt-2 divide-y divide-gray-100 dark:divide-gray-800">
          {#each group.items as flag}
            <li class="flex gap-3 py-2.5 text-[15px] leading-snug text-gray-700 dark:text-gray-300">
              <span class="mt-[0.45em] w-1.5 h-1.5 flex-shrink-0 rounded-full {group.dot}"></span>
              <span class="flex-1">{flag.text}</span>
              <span class="flex-shrink-0 font-mono text-[11px] text-gray-400 pt-0.5"
                >{group.sign}{flag.pts} {group.kind}</span
              >
            </li>
          {/each}
        </ul>
      {/each}
      <figcaption class="mt-4 font-mono text-[11px] text-gray-400">
        Example result, simplified. Not a real site. Points are shown here to explain the score.
      </figcaption>
    </figure>

    <dl class="mt-10 grid sm:grid-cols-2 gap-x-10 gap-y-7">
      {#each PARTS as part, i}
        <!-- An odd one out at the end spans both columns instead of sitting alone. -->
        <div
          class="pt-4 border-t border-gray-200 dark:border-gray-800 {i === PARTS.length - 1 &&
          PARTS.length % 2
            ? 'sm:col-span-2'
            : ''}"
        >
          <dt class="font-serif text-[1.35rem] leading-tight">{part.title}</dt>
          <dd class="mt-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {part.desc}
          </dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- Scoring -->
  <section class="mt-24" id="score">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">How the score is calculated</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        Every finding is worth a set number of points. Red flags add risk points, green flags add
        trust points. The score starts at 50, the neutral middle, and moves by half the difference
        between the two, staying between 0 and 100.
      </p>
    </div>

    <p
      class="mt-6 overflow-x-auto rounded-xl bg-gray-100 dark:bg-gray-900 px-5 py-4 font-mono text-sm text-gray-800 dark:text-gray-200"
    >
      score = 50 + (trust − risk) ÷ 2
    </p>

    <p class="mt-6 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      For the example above: {trust} trust − {risk} risk = {trust - risk}. Half of that, rounded, is
      {shift}. So 50 {shift < 0 ? "−" : "+"}
      {Math.abs(shift)} =
      <strong class="font-medium text-gray-900 dark:text-gray-100">{score}</strong>, which is Risky.
    </p>

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
          style="left: {score}%; transform: translateX(-50%)"
        >
          <span class="font-mono text-[10px] text-gray-500 leading-none">{score}</span>
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

  <!-- Checks -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">The checks</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      18 checks start at the same time the moment you submit. Each one is independent, so a slow or
      failed check never holds up the rest, and the score is worked out once they're all back.
      Curious what gets sent where during a scan? See the <a href="/privacy" class={LINK}
        >privacy page</a
      >.
    </p>

    <dl class="mt-8">
      {#each CHECKS as item}
        <div
          id="check-{item.id}"
          class="scroll-mt-20 flex flex-col sm:flex-row gap-1 sm:gap-6 py-5 border-t border-gray-200 dark:border-gray-800"
        >
          <dt
            class="sm:w-44 flex-shrink-0 font-mono text-xs uppercase tracking-wider text-gray-500 sm:pt-1"
          >
            {item.label}
          </dt>
          <dd class="text-[15px] leading-relaxed text-gray-700 dark:text-gray-300">{item.desc}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- For developers -->
  <section class="mt-24" id="developers">
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">For developers</p>
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">
      Want url.vet in your own app?
    </h2>
    <p class="mt-4 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
      url.vet is <a href={REPO} class={LINK} target="_blank" rel="noopener noreferrer"
        >open source</a
      >, so you can run your own copy and connect your app to it. Same checks as this site, on a
      server you control.
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

  <!-- FAQ -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Frequently asked questions</h2>
    <dl class="mt-8">
      {#each FAQ as item}
        <div class="py-5 border-t border-gray-200 dark:border-gray-800">
          <dt class="font-serif text-[1.35rem] leading-tight">{item.q}</dt>
          <dd class="mt-2 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
            {item.a}
          </dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- Back to the product -->
  <section
    class="mt-24 pt-12 border-t border-gray-200 dark:border-gray-800 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6"
  >
    <p class="font-serif text-3xl md:text-4xl tracking-[-0.01em] leading-tight">
      Got a link you're <span class="italic">unsure about?</span>
    </p>
    <a href="/" class="{PILL_SOLID} px-6 py-3 self-start sm:self-auto flex-shrink-0">
      Check it on url.vet
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
          d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3"
        />
      </svg>
    </a>
  </section>
</div>
