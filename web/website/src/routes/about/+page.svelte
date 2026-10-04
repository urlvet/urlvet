<script lang="ts">
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { LINK } from "$lib/ui/text";
  import { GITHUB_LOGO } from "$lib/ui/icons";
  import { PILL_OUTLINE } from "$lib/ui/buttons";

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
      title: "Checked live",
      desc: "Phishing pages appear and vanish within hours. url.vet checks the link the moment you ask, so it catches pages that blocklists haven't added yet.",
    },
    {
      title: "Privacy first",
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
    <!-- A note from a person: a card with who's speaking up top, and a signature at the end. -->
    <figure
      class="relative mt-8 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 px-6 py-7 md:px-10 md:py-9 shadow-sm shadow-black/[0.03]"
    >
      <figcaption class="flex items-center gap-3.5">
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
              <path d={GITHUB_LOGO} />
            </svg>
            @abhizaik
          </a>
        </div>
      </figcaption>
      <span
        class="hidden sm:block pointer-events-none absolute right-6 top-2 md:right-9 md:top-3 font-serif text-[5.5rem] leading-none text-accent-light/25 dark:text-accent-dark/25 select-none"
        aria-hidden="true">&ldquo;</span
      >

      <blockquote
        class="mt-7 max-w-[36rem] space-y-5 text-[17px] md:text-lg leading-relaxed text-gray-800 dark:text-gray-200"
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
          Also I’m a pretty private person, and I don’t like how most apps treat their users today:
          tracking everything, collecting all your data, selling it to advertisers. I found all of
          it irritating and wanted to build something that does things the right way.
        </p>
        <p>That’s how the idea for url.vet started.</p>
      </blockquote>
      <!-- Signed by hand: the name with a pen flourish under it. -->
      <p class="mt-8 flex items-center gap-3 text-gray-500">
        <span class="w-6 h-px bg-current opacity-60" aria-hidden="true"></span>
        <span
          class="relative inline-block font-serif italic text-[1.5rem] leading-none text-gray-900 dark:text-gray-100"
          aria-label="Signed, abhizaik"
          >abhizaik<svg
            class="absolute -left-2 -bottom-3 w-[calc(100%+1rem)] h-4 overflow-visible text-gray-400 dark:text-gray-500"
            viewBox="0 0 100 16"
            preserveAspectRatio="none"
            aria-hidden="true"
            ><path
              d="M3 11C22 6 48 4 72 5.5C86 6.5 95 8 97 5.5C98.5 3.5 94 2.5 90 5"
              fill="none"
              stroke="currentColor"
              stroke-width="1.3"
              stroke-linecap="round"
              vector-effect="non-scaling-stroke"
            /></svg
          ></span
        >
      </p>
    </figure>
  </section>

  <!-- Principles -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Philosophy</h2>
    <ul class="mt-8">
      {#each PRINCIPLES as item}
        <li
          class="py-5 border-t border-gray-200 dark:border-gray-800 sm:grid sm:grid-cols-[13rem_1fr] sm:gap-8"
        >
          <p class="font-serif text-[1.45rem] leading-tight">{item.title}</p>
          <p class="mt-2 sm:mt-1 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {item.desc}
          </p>
        </li>
      {/each}
    </ul>
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
        <path d={GITHUB_LOGO} />
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

  <PageCta />

  <p class="mt-16 font-mono text-[11px] text-gray-400">
    url.vet is also known as URLvet, and was previously called SafeSurf.
  </p>
</div>
