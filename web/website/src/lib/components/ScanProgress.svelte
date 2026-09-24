<script lang="ts">
  import { onDestroy } from "svelte";
  import { fade, fly } from "svelte/transition";

  export let loading = false;
  export let done = false;

  const STEPS = [
    "Resolving domain & IP",
    "URL structure & patterns",
    "TLD & DNS records",
    "Redirections & HTTP status",
    "WHOIS & domain age",
    "Homoglyphs & typosquatting",
    "SSL & TLS certificate",
    "Threat intelligence feeds",
    "Page content & forms",
    "Computing trust score",
  ];

  let activeStep = -1;
  let completedSteps: Set<number> = new Set();
  let elapsed = 0;
  let startTime = 0;
  let intervalTimer: ReturnType<typeof setInterval>;
  let elapsedTimer: ReturnType<typeof setInterval>;

  $: progress =
    completedSteps.size === STEPS.length
      ? 100
      : Math.round((completedSteps.size / STEPS.length) * 100);

  $: elapsedStr = (elapsed / 10).toFixed(1) + "s";

  function startProgress() {
    activeStep = 0;
    completedSteps = new Set();
    elapsed = 0;
    startTime = performance.now();

    elapsedTimer = setInterval(() => {
      elapsed = Math.round((performance.now() - startTime) / 100);
    }, 100);

    let stepIdx = 0;
    intervalTimer = setInterval(() => {
      if (stepIdx < STEPS.length - 2) {
        completedSteps = new Set([...completedSteps, stepIdx]);
        stepIdx++;
        activeStep = stepIdx;
      }
    }, 1000);
  }

  function finishProgress() {
    clearInterval(intervalTimer);
    clearInterval(elapsedTimer);
    let i = activeStep;
    const flush = setInterval(() => {
      if (i < STEPS.length) {
        completedSteps = new Set([...completedSteps, i]);
        i++;
        activeStep = i;
      } else {
        clearInterval(flush);
        activeStep = STEPS.length;
      }
    }, 60);
  }

  $: if (loading) startProgress();
  $: if (done) finishProgress();

  onDestroy(() => {
    clearInterval(intervalTimer);
    clearInterval(elapsedTimer);
  });
</script>

{#if loading || done}
  <div class="max-w-lg mx-auto mt-6 select-none" in:fade={{ duration: 200 }}>
    <div
      class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 overflow-hidden"
    >
      <!-- Header -->
      <div
        class="flex items-center justify-between px-5 py-3.5 border-b border-gray-200 dark:border-gray-800"
      >
        <div class="flex items-center gap-2">
          {#if activeStep < STEPS.length}
            <span class="relative flex h-2 w-2">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent-light dark:bg-accent-dark opacity-60"
              ></span>
              <span
                class="relative inline-flex rounded-full h-2 w-2 bg-accent-light dark:bg-accent-dark"
              ></span>
            </span>
            <span class="font-mono text-[11px] text-gray-500 tracking-wider uppercase"
              >Scanning</span
            >
          {:else}
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            <span
              class="font-mono text-[11px] text-emerald-600 dark:text-emerald-400 tracking-wider uppercase"
              >Complete</span
            >
          {/if}
        </div>
        <span class="text-xs font-mono text-gray-500">{elapsedStr}</span>
      </div>

      <!-- Active step (prominent) -->
      <div class="px-5 pt-5 pb-4 min-h-[3.5rem] flex items-center">
        {#key activeStep}
          <p
            class="font-serif text-2xl text-gray-900 dark:text-gray-100 leading-snug"
            in:fly={{ y: 6, duration: 180 }}
          >
            {#if activeStep >= 0 && activeStep < STEPS.length}
              {STEPS[activeStep]}
            {:else if activeStep === STEPS.length}
              Analysis complete
            {:else}
              Initializing...
            {/if}
          </p>
        {/key}
      </div>

      <!-- Progress bar -->
      <div class="px-5 pb-5">
        <div class="flex justify-end mb-1.5">
          <span class="text-[10px] font-mono text-gray-500">{progress}%</span>
        </div>
        <div class="h-1 bg-gray-200 dark:bg-gray-800 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-500 ease-out {progress === 100
              ? 'bg-emerald-500'
              : 'bg-gray-900 dark:bg-gray-100'}"
            style="width: {progress}%"
          ></div>
        </div>
      </div>
    </div>
  </div>
{/if}
