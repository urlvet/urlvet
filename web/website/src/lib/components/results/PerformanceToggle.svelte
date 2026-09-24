<script lang="ts">
  import { CHEVRON_DOWN } from "../../ui/icons";
  import PerformanceSection from "../sections/PerformanceSection.svelte";

  // Collapsed-by-default task timings, for power users.
  export let performance: any;

  let open = false;
</script>

<div class="flex flex-col items-center gap-2 pt-2">
  <button
    type="button"
    class="inline-flex items-center gap-1.5 font-mono text-[11px] text-gray-500 hover:text-gray-800 dark:hover:text-gray-300 transition-colors focus:outline-none"
    aria-expanded={open}
    on:click={() => (open = !open)}
  >
    <svg
      class="w-3 h-3 transition-transform duration-200 {open ? 'rotate-180' : ''}"
      viewBox="0 0 20 20"
      fill="currentColor"
      aria-hidden="true"><path fill-rule="evenodd" d={CHEVRON_DOWN} clip-rule="evenodd" /></svg
    >
    Performance timings
  </button>
  {#if open}
    <div class="w-full perf-body">
      <PerformanceSection {performance} />
    </div>
  {/if}
</div>

<style>
  /* The section brings its own card; here it sits flat under the toggle. */
  .perf-body > :global(section),
  .perf-body > :global(div) {
    border: none !important;
    border-radius: 0 !important;
    box-shadow: none !important;
    background: transparent !important;
    margin: 0 !important;
  }
</style>
