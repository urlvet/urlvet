<script lang="ts">
  import { PILL_OUTLINE, PILL_SOLID } from "$lib/ui/buttons";

  const REPO = "https://github.com/urlvet/urlvet";
  const GITHUB_PATH =
    "M12 .5C5.648.5.5 5.648.5 12c0 5.084 3.292 9.387 7.872 10.905.576.106.784-.25.784-.556 0-.274-.01-1-.016-1.974-3.2.694-3.876-1.544-3.876-1.544-.522-1.33-1.276-1.683-1.276-1.683-1.042-.714.08-.699.08-.699 1.152.08 1.756 1.183 1.756 1.183 1.024 1.754 2.688 1.246 3.344.952.104-.742.4-1.246.728-1.532-2.554-.292-5.238-1.276-5.238-5.674 0-1.252.448-2.274 1.184-3.076-.12-.292-.512-1.468.112-3.064 0 0 .964-.308 3.16 1.176.916-.254 1.9-.382 2.876-.388.976.006 1.96.134 2.876.388 2.192-1.484 3.156-1.176 3.156-1.176.624 1.596.232 2.772.112 3.064.736.802 1.184 1.824 1.184 3.076 0 4.408-2.69 5.38-5.25 5.664.412.354.78 1.05.78 2.118 0 1.532-.014 2.768-.014 3.144 0 .308.2.672.788.556C20.708 21.385 24 17.084 24 12c0-6.352-5.148-11.5-12-11.5z";

  const LINK =
    "text-gray-900 dark:text-gray-100 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:decoration-current transition-colors";

  const PRINCIPLES = [
    {
      title: "Explainable results",
      desc: "A verdict you can't explain can't be trusted. Every signal is shown, along with how it led to the final verdict.",
    },
    {
      title: "Built for everyone",
      desc: "Detailed enough for security folks, plain enough for your parents. Same scan, read at whatever depth you like.",
    },
    {
      title: "Fully open source",
      desc: "Security tools should be open to inspection. The detection engine, scoring and website are all on GitHub, free to read, audit and self-host.",
    },
    {
      title: "Checked live, not from a list",
      desc: "Phishing pages appear and vanish within hours. url.vet checks the link the moment you ask, so it catches pages that blocklists haven't added yet.",
    },
    {
      title: "Private by design",
      desc: "No account, no ads, no trackers, no data sold. We don't log who you are or what you scan.",
    },
  ];
</script>

<svelte:head>
  <title>About — url.vet (URLvet)</title>
  <meta
    name="description"
    content="url.vet (URLvet) tells you whether a link is safe, and why. Open source, private by design, and built for experts and everyone else."
  />
  <link rel="canonical" href="https://url.vet/about" />
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">About</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    Doubt, <span class="italic">made useful.</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    Paste any link and url.vet tells you whether it's safe, and why. Detailed enough for experts,
    plain enough for everyone else.
  </p>

  <!-- Maker: the why, in first person -->
  <section class="mt-20">
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">Why I built it</h2>
    <figure class="mt-6">
      <!-- A letter: body text in the page's sans, with only the sign-off in serif. -->
      <blockquote
        class="max-w-[36rem] space-y-5 text-[17px] md:text-lg leading-relaxed text-gray-800 dark:text-gray-200"
      >
        <p>
          I’d often get a link and wonder if it was safe to open. I tried a few online link
          scanners. Some just said “safe” or “risky” and nothing else. Others gave me pages of
          technical data to dig through myself. I wasn’t really happy with either.
        </p>
        <p>
          Then I started noticing it at home. My parents would open and forward links from WhatsApp
          without really knowing if they were real (and most of the time they weren’t). Friends
          would send me “hot deals” that didn’t even go to the shopping website they claimed to.
        </p>
        <p>
          I wanted something I could use myself and also hand to my parents. A straight answer for
          them, and the reasons behind it for anyone who wants to look.
        </p>
        <p>
          I’m a pretty private person, and I don’t like how most apps treat their users today:
          tracking everything, collecting all your data, selling it to advertisers. I found all of
          it irritating and wanted something that does things the right way.
        </p>
        <p
          class="pt-1 font-serif italic text-2xl md:text-[1.7rem] leading-snug text-gray-900 dark:text-gray-100"
        >
          That’s how the idea for url.vet started.
        </p>
      </blockquote>
      <figcaption class="mt-8 flex items-center gap-4">
        <img
          src="/team/abhizaik.jpg"
          alt=""
          class="w-11 h-11 rounded-full flex-shrink-0 grayscale-[30%]"
          loading="lazy"
        />
        <div class="min-w-0">
          <p class="text-[15px]">
            <a
              href="https://abhizaik.com"
              target="_blank"
              rel="noopener noreferrer"
              class="font-medium hover:underline underline-offset-4">abhizaik</a
            ><span class="text-gray-500">, creator &amp; maintainer</span>
          </p>
          <a
            href="https://github.com/abhizaik"
            target="_blank"
            rel="noopener noreferrer"
            class="mt-0.5 inline-flex items-center gap-1.5 text-[13px] text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          >
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path d={GITHUB_PATH} />
            </svg>
            @abhizaik
          </a>
        </div>
      </figcaption>
    </figure>
  </section>

  <!-- Principles -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Philosophy</h2>
    <dl class="mt-8 grid sm:grid-cols-2 gap-x-10 gap-y-8">
      {#each PRINCIPLES as item}
        <div class="pt-4 border-t border-gray-200 dark:border-gray-800">
          <dt class="font-serif text-[1.45rem] leading-tight">{item.title}</dt>
          <dd class="mt-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {item.desc}
          </dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- Open source + help -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">
      Open, <span class="italic">all of it.</span>
    </h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        Plenty of projects call themselves open source but only publish the website. Here the
        detection engine, the scoring, the API and the website are all public under <a
          href="{REPO}/blob/main/LICENSE"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>AGPL-3.0</a
        >. Read it, audit it, or host it yourself.
      </p>
      <p>
        Want to help? Start with a <a
          href="{REPO}/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>good first issue</a
        >
        or read the
        <a
          href="{REPO}/blob/main/.github/CONTRIBUTING.md"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>contributing guide</a
        >.
      </p>
    </div>
    <a
      href={REPO}
      target="_blank"
      rel="noopener noreferrer"
      class="mt-7 {PILL_OUTLINE} px-5 py-2.5"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path d={GITHUB_PATH} />
      </svg>
      View on GitHub
    </a>
  </section>

  <!-- Contact -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Get in touch</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        Feedback, ideas, or just to say hi: <a
          href="mailto:hi@url.vet"
          class="font-mono text-base {LINK}">hi@url.vet</a
        >.
      </p>
      <p>
        Think a verdict is wrong? The Report button on any result is the fastest way. A human reads
        every report. Found a security issue? Please follow the <a
          href="{REPO}/blob/main/.github/SECURITY.md"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>security policy</a
        >.
      </p>
    </div>
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

  <p class="mt-16 font-mono text-[11px] text-gray-400">
    url.vet is also known as URLvet, and was previously called SafeSurf.
  </p>
</div>
