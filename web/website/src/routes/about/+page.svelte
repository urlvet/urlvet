<svelte:head>
  <title>About — url.vet (URLvet)</title>
  <meta
    name="description"
    content="url.vet (URLvet) is an open-source phishing detection engine built to make the web a little safer. Free, transparent, and fully explainable."
  />
  <link rel="canonical" href="https://url.vet/about" />
</svelte:head>

<div class="min-h-screen bg-white dark:bg-gray-950 text-gray-900 dark:text-gray-100">
  <div class="max-w-3xl mx-auto px-6 py-16">
    <!-- Back link -->
    <a
      href="/"
      class="inline-flex items-center gap-2 text-base text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors mb-10"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      url.vet
    </a>

    <h1 class="text-3xl font-bold tracking-tight mb-3">About</h1>
    <p class="text-gray-500 dark:text-gray-400 text-lg mb-12">
      A fully open-source phishing detection engine that runs multiple checks in real time and
      explains exactly what it found. No accounts, no hidden algorithms. Built for people who want
      to understand why a link is safe, not just be told it is.
    </p>

    <!-- Origin -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-4">Why this exists</h2>
      <div class="space-y-4 text-gray-600 dark:text-gray-400 text-base leading-relaxed">
        <p>
          Phishing links show up in emails, DMs, QR codes, and shortened URLs. Most people either
          click and hope, or avoid the link entirely. Neither is good. Existing tools don't close
          that gap: quick scanners return a pass/fail with no explanation, thorough crawlers take
          too long and bury results in raw technical data. Nothing is both fast and legible.
        </p>
        <p>
          url.vet (also known as URLvet, previously SafeSurf) is that middle ground. Paste a link,
          get a trust score, a verdict, and a full breakdown in seconds. No guesswork, no signup,
          nothing hidden.
        </p>
        <p>
          It runs 18 checks: URL structure, DNS, TLS, domain age, typosquatting, page content, and
          live threat feeds. Each produces a plain-language reason, so you always know what moved
          the score.
        </p>
      </div>
    </section>

    <!-- Philosophy -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-5">Philosophy</h2>
      <div class="grid sm:grid-cols-2 gap-4">
        {#each [{ title: "Explainability over security theater", desc: "A verdict you can't explain is a verdict you can't trust. Every signal is shown, whether it helped or hurt the score." }, { title: "Live over cached", desc: "Phishing pages spin up and disappear in hours. Scanning at request time catches what static databases miss." }, { title: "Open source by default", desc: "Security tools should be auditable. The entire detection engine is on GitHub under AGPL-3.0." }, { title: "No friction", desc: "No account, no ads, no API key. Paste a URL, get an answer." }] as item}
          <div
            class="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4"
          >
            <p class="font-medium text-base mb-1.5">{item.title}</p>
            <p class="text-gray-500 dark:text-gray-400 text-base leading-relaxed">{item.desc}</p>
          </div>
        {/each}
      </div>
    </section>

    <!-- Team -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-6">Team</h2>
      <div
        class="flex items-start gap-4 p-4 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl mb-4"
      >
        <img
          src="https://avatars.githubusercontent.com/u/65501471?v=4"
          alt="abhizaik"
          class="w-14 h-14 rounded-full flex-shrink-0"
          loading="lazy"
        />
        <div>
          <a
            href="https://abhizaik.com"
            target="_blank"
            rel="noopener noreferrer"
            class="font-semibold hover:underline">abhizaik</a
          >
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 mb-2">
            Creator &amp; maintainer
          </p>
          <p class="text-base text-gray-600 dark:text-gray-400 leading-relaxed">
            Got one too many suspicious links with no good way to check them, so built url.vet.
          </p>
          <div class="flex items-center gap-4 mt-3">
            <a
              href="https://github.com/abhizaik"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 .5C5.648.5.5 5.648.5 12c0 5.084 3.292 9.387 7.872 10.905.576.106.784-.25.784-.556 0-.274-.01-1-.016-1.974-3.2.694-3.876-1.544-3.876-1.544-.522-1.33-1.276-1.683-1.276-1.683-1.042-.714.08-.699.08-.699 1.152.08 1.756 1.183 1.756 1.183 1.024 1.754 2.688 1.246 3.344.952.104-.742.4-1.246.728-1.532-2.554-.292-5.238-1.276-5.238-5.674 0-1.252.448-2.274 1.184-3.076-.12-.292-.512-1.468.112-3.064 0 0 .964-.308 3.16 1.176.916-.254 1.9-.382 2.876-.388.976.006 1.96.134 2.876.388 2.192-1.484 3.156-1.176 3.156-1.176.624 1.596.232 2.772.112 3.064.736.802 1.184 1.824 1.184 3.076 0 4.408-2.69 5.38-5.25 5.664.412.354.78 1.05.78 2.118 0 1.532-.014 2.768-.014 3.144 0 .308.2.672.788.556C20.708 21.385 24 17.084 24 12c0-6.352-5.148-11.5-12-11.5z"
                />
              </svg>
              @abhizaik
            </a>
            <a
              href="mailto:hi@url.vet"
              class="inline-flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <svg
                class="w-3.5 h-3.5"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                />
              </svg>
              hi@url.vet
            </a>
          </div>
        </div>
      </div>
      <p class="text-base text-gray-500 dark:text-gray-400">
        With help from
        <a
          href="https://github.com/urlvet/urlvet/graphs/contributors"
          target="_blank"
          rel="noopener noreferrer"
          class="text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white underline underline-offset-2 transition-colors"
          >friends</a
        > on GitHub.
      </p>
    </section>

    <!-- Open source -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-4">Open source</h2>
      <p class="text-base text-gray-600 dark:text-gray-400 mb-3 leading-relaxed">
        url.vet is fully open source. A lot of projects call themselves open source but only publish
        the frontend while keeping the actual engine private. Here, everything is publicly
        available. The detection engine, scoring logic, API, and frontend. Nothing hidden.
      </p>
      <p class="text-base text-gray-600 dark:text-gray-400 mb-5 leading-relaxed">
        The community edition is licensed under
        <a
          href="https://github.com/urlvet/urlvet/blob/main/LICENSE"
          target="_blank"
          rel="noopener noreferrer"
          class="text-gray-900 dark:text-white underline underline-offset-2">AGPL-3.0</a
        >. Free to use, self-host, and modify. Found a bug or want to contribute? Pull requests are
        welcome.
      </p>
      <div class="flex flex-wrap gap-3">
        <a
          href="https://github.com/urlvet/urlvet"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 bg-gray-900 dark:bg-white text-white dark:text-gray-900 px-4 py-2 rounded-lg text-base font-medium hover:bg-gray-700 dark:hover:bg-gray-200 transition-colors"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M12 .5C5.648.5.5 5.648.5 12c0 5.084 3.292 9.387 7.872 10.905.576.106.784-.25.784-.556 0-.274-.01-1-.016-1.974-3.2.694-3.876-1.544-3.876-1.544-.522-1.33-1.276-1.683-1.276-1.683-1.042-.714.08-.699.08-.699 1.152.08 1.756 1.183 1.756 1.183 1.024 1.754 2.688 1.246 3.344.952.104-.742.4-1.246.728-1.532-2.554-.292-5.238-1.276-5.238-5.674 0-1.252.448-2.274 1.184-3.076-.12-.292-.512-1.468.112-3.064 0 0 .964-.308 3.16 1.176.916-.254 1.9-.382 2.876-.388.976.006 1.96.134 2.876.388 2.192-1.484 3.156-1.176 3.156-1.176.624 1.596.232 2.772.112 3.064.736.802 1.184 1.824 1.184 3.076 0 4.408-2.69 5.38-5.25 5.664.412.354.78 1.05.78 2.118 0 1.532-.014 2.768-.014 3.144 0 .308.2.672.788.556C20.708 21.385 24 17.084 24 12c0-6.352-5.148-11.5-12-11.5z"
            />
          </svg>
          View on GitHub
        </a>
        <a
          href="https://github.com/urlvet/urlvet/issues"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 border border-gray-300 dark:border-gray-700 px-4 py-2 rounded-lg text-base font-medium hover:border-gray-500 dark:hover:border-gray-500 transition-colors"
        >
          Report an issue
        </a>
      </div>
    </section>

    <!-- Contact -->
    <section>
      <h2 class="text-xl font-semibold mb-4">Get in touch</h2>
      <p class="text-base text-gray-600 dark:text-gray-400 leading-relaxed">
        Feedback, false positives, or just to say hi. Reach out at
        <a
          href="mailto:hi@url.vet"
          class="text-gray-900 dark:text-white underline underline-offset-2 hover:no-underline transition-all"
          >hi@url.vet</a
        >.
      </p>
    </section>
  </div>
</div>
