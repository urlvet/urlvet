<script lang="ts">
  export let reasons: { bad_reasons?: string[]; good_reasons?: string[] } | undefined;

  $: groups = [
    {
      label: "Red flags",
      items: reasons?.bad_reasons ?? [],
      dot: "bg-red-500",
      title: "text-red-600 dark:text-red-400",
    },
    {
      label: "Green flags",
      items: reasons?.good_reasons ?? [],
      dot: "bg-emerald-500",
      title: "text-emerald-700 dark:text-emerald-400",
    },
  ];
</script>

{#if reasons}
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#each groups as g}
      <div
        class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-5 sm:p-6"
      >
        <h3
          class="flex items-center gap-2 font-mono text-[11px] uppercase tracking-wider {g.items
            .length
            ? g.title
            : 'text-gray-500'}"
        >
          <span class="w-1.5 h-1.5 rounded-full {g.items.length ? g.dot : 'bg-gray-400'}"></span>
          {g.label}
        </h3>
        {#if g.items.length}
          <ul class="mt-4 space-y-2.5 text-[15px] leading-snug text-gray-800 dark:text-gray-200">
            {#each g.items as r}
              <li class="flex items-start gap-3">
                <span class="mt-[0.55em] h-1 w-1 rounded-full shrink-0 {g.dot}"></span>
                <span class="break-words min-w-0">{r}</span>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="mt-4 text-[15px] text-gray-500 dark:text-gray-400">None found.</p>
        {/if}
      </div>
    {/each}
  </div>
{/if}
