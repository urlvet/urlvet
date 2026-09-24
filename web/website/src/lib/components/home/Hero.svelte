<script lang="ts">
  // The url.vet wordmark and tagline; slightly smaller (and without the descriptor) on results.
  export let isLanding: boolean;
</script>

<header class="relative flex flex-col items-center text-center {isLanding ? 'mb-12' : 'mb-8'}">
  <h1
    class="font-serif font-normal tracking-[-0.02em] leading-none {isLanding
      ? 'text-7xl md:text-8xl'
      : 'text-4xl md:text-[2.75rem]'}"
  >
    <a href="/" on:click={() => (location.href = "/")} class="group">
      <span
        class="text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400 transition-colors"
        >url</span
      ><span class="-mx-[0.025em] text-accent-light dark:text-accent-dark">.</span><span
        class="text-gray-900 dark:text-gray-100">vet</span
      >
    </a>
    <span class="sr-only">&nbsp;(also known as URLvet)</span>
  </h1>

  <p
    class="font-serif italic [word-spacing:0.08em] [text-wrap:balance] text-gray-700 dark:text-gray-300 {isLanding
      ? 'mt-5 text-xl md:text-2xl'
      : 'mt-3 text-xl'}"
  >
    sketchy link? just
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
  </p>
  {#if isLanding}
    <p
      class="mt-3 font-serif text-gray-600 dark:text-gray-400 text-base md:text-[17px] [text-wrap:balance]"
    >
      Vet any URL for phishing, scams & suspicious redirects
    </p>
  {/if}
</header>
