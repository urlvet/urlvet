<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { TAGLINE_LINKS } from "../../data/tagline";

  // The url.vet wordmark and tagline; slightly smaller (and without the descriptor) on results.
  export let isLanding: boolean;
  /** Drives the logo dot: a slow pulse when idle, a quick blink while scanning, the verdict's colour after. */
  export let dot: "idle" | "scanning" | "Safe" | "Suspicious" | "Risky" = "idle";

  // Typed into the tagline, each followed by "?".
  const EXAMPLES = TAGLINE_LINKS.map((l) => `${l}?`);

  /** What the tagline's first slot shows; null means the plain "sketchy link?". */
  let typed: string | null = null;
  /** Plays the one-time "vetting" sweep over the wordmark. */
  let intro = false;
  /** Rendered widths of "sketchy link?" and of each example, measured once (hidden copies below). */
  let restingWidth = 0;
  let exampleWidths: number[] = [];
  /** The example being shown. The slot widens to fit it before typing starts, so the
      rest of the line glides once instead of shifting with every letter. */
  let target: string | null = null;
  // On narrow screens a long example may not fit beside "just url.vet it"; it's then
  // set smaller (scale) rather than wrapping, so the whole link stays on one line.
  let headerWidth = 0;
  let tailWidth = 0;
  $: room = headerWidth - tailWidth - 12;
  $: natural = target && isLanding ? (exampleWidths[EXAMPLES.indexOf(target)] ?? 0) : 0;
  $: scale = natural > room && room > 0 ? room / natural : 1;
  $: slotWidth = Math.max(restingWidth, Math.min(natural, room > 0 ? room : natural));

  let timer: ReturnType<typeof setTimeout> | undefined;
  const wait = (ms: number) => new Promise<void>((r) => (timer = setTimeout(r, ms)));
  let stopped = false;

  // Type an example out, hold it, erase it, show "sketchy link?" again, next example.
  async function cycle() {
    // A fresh order each visit, so returning visitors don't always start with the same one.
    const order = [...EXAMPLES];
    for (let i = order.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [order[i], order[j]] = [order[j], order[i]];
    }
    for (let i = 0; !stopped; i = (i + 1) % order.length) {
      await wait(2600);
      // Only the landing page rotates; on results the header stays still.
      if (!isLanding) continue;
      const text = order[i];
      target = text;
      await wait(350);
      for (let n = 1; n <= text.length && !stopped; n++) {
        typed = text.slice(0, n);
        await wait(45);
      }
      await wait(2200);
      for (let n = text.length - 1; n >= 0 && !stopped; n--) {
        typed = text.slice(0, n);
        await wait(22);
      }
      typed = null;
      target = null;
    }
  }

  onMount(() => {
    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    if (reduced || !isLanding) return;
    try {
      if (!sessionStorage.getItem("hero-intro")) {
        intro = true;
        sessionStorage.setItem("hero-intro", "1");
      }
    } catch {
      /* storage blocked: skip the intro */
    }
    cycle();
  });

  onDestroy(() => {
    stopped = true;
    clearTimeout(timer);
  });

  const DOT_COLOR: Record<typeof dot, string> = {
    idle: "text-accent-light dark:text-accent-dark",
    scanning: "text-accent-light dark:text-accent-dark",
    Safe: "text-emerald-500",
    Suspicious: "text-yellow-500",
    Risky: "text-red-500",
  };
</script>

<header
  bind:clientWidth={headerWidth}
  class="relative w-full flex flex-col items-center text-center {isLanding ? 'mb-12' : 'mb-8'}"
>
  <h1
    class="font-serif font-normal tracking-[-0.02em] leading-none {isLanding
      ? 'text-7xl md:text-8xl'
      : 'text-4xl md:text-[2.75rem]'}"
    class:intro
  >
    <a href="/" on:click={() => (location.href = "/")} class="group">
      <span
        class="url relative text-[#b3a998] dark:text-[#6e665a] group-hover:text-gray-500 dark:group-hover:text-gray-400 transition-colors"
        >url</span
      ><span
        class="dot inline-block -mx-[0.045em] transition-colors duration-500 {DOT_COLOR[dot]}"
        data-state={dot}>.</span
      ><span class="text-gray-900 dark:text-gray-100">vet</span>
    </a>
    <span class="sr-only">&nbsp;(also known as URLvet)</span>
  </h1>

  <p
    class="font-serif italic [word-spacing:0.08em] text-gray-700 dark:text-gray-300 {isLanding
      ? 'mt-5 text-xl md:text-2xl'
      : 'mt-3 text-xl'}"
  >
    <span class="sr-only">sketchy link? just url.vet it</span>
    <span
      aria-hidden="true"
      class="relative inline-flex flex-wrap items-baseline justify-center gap-x-[0.3em]"
    >
      {#if isLanding}
        <span class="absolute invisible whitespace-nowrap pointer-events-none">
          <span bind:offsetWidth={restingWidth}>sketchy link?</span>
          {#each EXAMPLES as example, i}
            <span bind:offsetWidth={exampleWidths[i]}
              ><span class="not-italic font-mono text-[0.72em] tracking-tight">{example}</span><span
                class="caret"
              ></span></span
            >
          {/each}
        </span>
      {/if}
      <!-- Right-aligned at a set width (see slotWidth), so "just url.vet it" stays put while typing. -->
      <span
        class="slot inline-block text-right whitespace-nowrap"
        style={slotWidth ? `width: ${slotWidth}px` : undefined}
      >
        {#if typed === null || !isLanding}
          sketchy link?
        {:else}
          <span
            class="not-italic font-mono tracking-tight text-gray-900 dark:text-gray-100"
            style="font-size: {0.72 * scale}em">{typed}</span
          ><span class="caret" aria-hidden="true"></span>
        {/if}
      </span>
      <span class="whitespace-nowrap" bind:offsetWidth={tailWidth}>
        just
        <span class="relative inline-block"
          >url.vet<svg
            class="absolute left-0 -bottom-0.5 w-full h-[0.22em] overflow-visible text-accent-light dark:text-accent-dark"
            viewBox="0 0 100 8"
            preserveAspectRatio="none"
            aria-hidden="true"
            ><path
              d="M1 5.5C20 2.5 45 1.8 99 4.2"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              vector-effect="non-scaling-stroke"
            /></svg
          ></span
        >
        it
      </span>
    </span>
  </p>
  {#if isLanding}
    <p
      class="mt-3 font-serif text-gray-600 dark:text-gray-400 text-base md:text-[17px] [text-wrap:balance]"
    >
      Vet any URL for phishing, scams & suspicious redirects
    </p>
  {/if}
</header>

<style>
  /* The dot sits on the baseline, so scale it from there. */
  .dot {
    transform-origin: 50% 85%;
  }
  .dot[data-state="idle"] {
    animation: breathe 3.2s ease-in-out infinite;
  }
  .dot[data-state="scanning"] {
    animation: blink 0.7s ease-in-out infinite;
  }
  @keyframes breathe {
    0%,
    100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.18);
      opacity: 0.8;
    }
  }
  @keyframes blink {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.25;
    }
  }

  /* One-time "vetting" moment: a line sweeps under "url", which darkens as if
     checked, then the dot drops in. */
  .intro .url::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0.08em;
    height: 0.035em;
    background: currentColor;
    transform-origin: left;
    animation: sweep 0.9s cubic-bezier(0.6, 0, 0.2, 1) 0.2s both;
  }
  .intro .url {
    animation: checked 1.4s ease 0.2s;
  }
  .intro .dot {
    animation:
      drop 0.5s cubic-bezier(0.3, 1.6, 0.5, 1) 1.05s both,
      breathe 3.2s ease-in-out 1.6s infinite;
  }
  @keyframes sweep {
    0% {
      transform: scaleX(0);
      opacity: 1;
    }
    70% {
      transform: scaleX(1);
      opacity: 1;
    }
    100% {
      transform: scaleX(1);
      opacity: 0;
    }
  }
  @keyframes checked {
    45% {
      color: var(--color-gray-700);
    }
  }
  :global(.dark) .intro .url {
    animation-name: checked-dark;
  }
  @keyframes checked-dark {
    45% {
      color: var(--color-gray-300);
    }
  }
  @keyframes drop {
    0% {
      transform: translateY(-0.6em);
      opacity: 0;
    }
    100% {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .slot {
    transition: width 0.35s ease;
  }

  .caret {
    display: inline-block;
    width: 0.06em;
    height: 0.9em;
    margin-left: 0.04em;
    vertical-align: -0.1em;
    background: currentColor;
    animation: blink 1s steps(1) infinite;
  }

  @media (prefers-reduced-motion: reduce) {
    .dot,
    .intro .dot,
    .intro .url,
    .intro .url::after,
    .caret {
      animation: none !important;
    }
    .slot {
      transition: none;
    }
  }
</style>
