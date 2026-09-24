<script lang="ts">
  export let performance: any;
</script>

{#if performance}
  <div
    class="bg-white dark:bg-gray-900/80 border border-gray-300 dark:border-gray-800 rounded-lg shadow-md hover:shadow-lg hover:scale-[1.01] transition-all"
  >
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 px-5 py-4"
    >
      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <span class="text-gray-400">Total time:</span>
        <span class="font-medium text-gray-800 dark:text-white">{performance.total_time}</span>
      </div>

      {#if performance.timings && performance.timings.length > 0}
        {#each performance.timings as timing}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <span class="text-gray-400">{timing.task}</span>
            <span class="font-medium text-gray-800 dark:text-white">{timing.time}</span>
          </div>
        {/each}
      {/if}
    </div>
  </div>
{/if}
