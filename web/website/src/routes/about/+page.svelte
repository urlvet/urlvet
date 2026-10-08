<script lang="ts">
  import { vettyMemory } from "$lib/components/guide/store";
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { PILL_OUTLINE } from "$lib/ui/buttons";
  import { GITHUB_LOGO } from "$lib/ui/icons";
  import { LINK } from "$lib/ui/text";
  import { tick } from "svelte";

  const VETTY = [
    {
      title: "What Vetty does",
      dot: "bg-emerald-500",
      items: [
        "Shows you how to use url.vet",
        "Checks a link for you",
        "Explains a result in simple words",
        "Takes you to the real site if a page is fake",
      ],
    },
    {
      title: "What Vetty never does",
      dot: "bg-red-500",
      items: [
        "Collect your data",
        "Use what you type to train AI",
        "Ask you to sign up or pay",
        "Show ads or try to sell you anything",
      ],
    },
  ];

  // Opens Vetty's bubble, bringing him back first if he was hidden.
  async function sayHi() {
    vettyMemory.set({ hidden: false });
    await tick();
    if (!document.querySelector(".vetty[aria-expanded='true']")) {
      document.querySelector<HTMLElement>("button.vetty")?.click();
    }
  }

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
  <title>About url.vet</title>
  <meta
    name="description"
    content="url.vet tells you if a link is safe, and why. Free, open source and private by design."
  />
  <meta property="og:title" content="About url.vet" />
  <meta
    property="og:description"
    content="url.vet tells you if a link is safe, and why. Free, open source and private by design."
  />
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://url.vet/about" />
  <meta property="og:image" content="https://url.vet/og" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="About url.vet" />
  <meta
    name="twitter:description"
    content="url.vet tells you if a link is safe, and why. Free, open source and private by design."
  />
  <meta name="twitter:image" content="https://url.vet/og" />
  <link rel="canonical" href="https://url.vet/about" />
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">About</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    Built for my loved ones. <span class="italic sm:block">Open to everyone.</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    url.vet is a free, open-source tool that tells you whether a link is safe, and why. This page is
    about why it exists and what it stands for.
  </p>

  <!-- Maker: the why, in first person -->
  <section class="mt-20">
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">Why I built it</h2>
    <!-- A letter: straight into the words, signed once at the end. -->
    <figure
      class="relative mt-8 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 px-6 py-7 md:px-10 md:py-9 shadow-xs shadow-black/[0.03]"
    >
      <span
        class="block h-9 font-serif text-[4.5rem] leading-none text-accent-light/40 dark:text-accent-dark/40 select-none"
        aria-hidden="true">&ldquo;</span
      >
      <blockquote
        class="mt-3 max-w-[36rem] space-y-5 text-[17px] md:text-lg leading-relaxed text-gray-800 dark:text-gray-200"
      >
        <p>
          I’d often get a link and wonder if it was safe to open. I tried a few online link
          scanners. Some just said “safe” or “risky” and nothing else. Others gave me pages of
          technical data to dig through myself. I wasn’t really happy with either.
        </p>
        <p>
          Then I saw the same problem at home. My parents would open and forward links from WhatsApp
          without really knowing if they were real (and often they weren’t). Friends would send me
          “hot deals” that didn’t even lead to the shop they claimed to be from.
        </p>
        <p>
          Also, I’m a pretty private person, and I don’t like how most apps treat their users today:
          tracking everything, collecting all your data, selling it to advertisers. I found all of
          it irritating and wanted to build something that does none of that.
        </p>
        <p>
          I wanted something my family, my friends and I could all use. A straight answer for them,
          and the reasons behind it for anyone who wants to look deeper.
        </p>
        <p class="font-medium text-gray-900 dark:text-gray-100">That’s how url.vet started.</p>
      </blockquote>

      <!-- Sign-off: photo and name, both linking to the creator's site. -->
      <figcaption class="mt-8">
        <a
          href="https://abhizaik.com"
          target="_blank"
          rel="noopener noreferrer"
          class="group inline-flex items-center gap-4"
        >
          <img
            src="/team/abhizaik.jpg"
            alt=""
            class="w-12 h-12 rounded-full shrink-0 grayscale-[30%] group-hover:grayscale-0 transition"
            loading="lazy"
          />
          <span>
            <span
              class="block font-signature text-[1.75rem] leading-none text-gray-900 dark:text-gray-100"
              >abhizaik</span
            >
            <span class="mt-1.5 block text-[13px] text-gray-500">Creator of url.vet</span>
          </span>
        </a>
      </figcaption>
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

  <!-- Vetty -->
  <section class="mt-24" id="vetty">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">Meet Vetty</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Vetty is the paperclip in the corner. He helps you use url.vet. He is inspired by Clippy, the
      paperclip helper from old office software.
    </p>

    <div class="mt-8 grid sm:grid-cols-2 gap-x-10 gap-y-8">
      {#each VETTY as list}
        <div class="pt-4 border-t border-gray-200 dark:border-gray-800">
          <h3 class="font-serif text-[1.45rem] leading-tight">{list.title}</h3>
          <ul class="mt-3 space-y-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {#each list.items as item}
              <li class="flex gap-2.5">
                <span
                  class="mt-[0.55em] w-1.5 h-1.5 rounded-full shrink-0 {list.dot}"
                  aria-hidden="true"
                ></span>
                <span>{item}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/each}
    </div>

    <p class="mt-8 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Don't want him around? Hide him from his menu. He will wait at the bottom of the page.
    </p>
    <button type="button" class="mt-5 {PILL_OUTLINE} px-5 py-2.5" on:click={sayHi}>
      Say hi to Vetty
    </button>
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
        Have feedback or an idea, or just want to say hello? Email <a
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
      <p class="text-[15px] text-gray-500">
        url.vet is also written as URLvet. It used to be called SafeSurf.
      </p>
    </div>
  </section>

  <PageCta />
</div>
