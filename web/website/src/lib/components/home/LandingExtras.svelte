<script lang="ts">
  import { onMount, tick } from "svelte";
  import { EXAMPLES } from "../../data/examples";
  // Example links to try, and the short list of promises under the search bar.
  export let onTry: (url: string) => void;

  const POINTS = ["Free and open source", "No tracking", "Explains every result"];

  // Keep the chips to two rows by measuring, not by breakpoint: drop the legit/fake
  // labels (level 1), then "Try" too (level 2), until they fit.
  let level = 0;
  let row: HTMLDivElement;

  const rows = () =>
    new Set(
      // By vertical centre: "Try" is shorter than the chips beside it.
      [...row.children].map((el) => {
        const e = el as HTMLElement;
        return Math.round((e.offsetTop + e.offsetHeight / 2) / 8);
      })
    ).size;

  async function fit() {
    if (!row) return;
    for (level = 0; level < 2; level++) {
      await tick();
      if (rows() <= 2) return;
    }
  }

  onMount(() => {
    fit();
    document.fonts?.ready.then(fit);
    let width = row.clientWidth;
    const ro = new ResizeObserver(() => {
      if (row.clientWidth !== width) {
        width = row.clientWidth;
        fit();
      }
    });
    ro.observe(row);
    return () => ro.disconnect();
  });
</script>

<div
  bind:this={row}
  data-guide="examples"
  class="mt-5 flex flex-wrap justify-center items-center gap-2"
>
  {#if level < 2}
    <span class="text-xs text-gray-500 mr-1">Try</span>
  {/if}
  {#each EXAMPLES as example}
    <button
      type="button"
      data-example={example.url}
      on:click={() => onTry(example.url)}
      title={example.hint === "Legit"
        ? "The real site"
        : "A lookalike, spelled with letters from another alphabet"}
      class="inline-flex items-center py-1.5 {level < 2
        ? 'gap-2 px-3.5'
        : 'gap-1.5 px-2.5'} rounded-full border border-gray-300 dark:border-gray-800 hover:border-gray-400 dark:hover:border-gray-600 text-gray-700 dark:text-gray-300 text-xs transition-colors"
    >
      <span
        class="w-1.5 h-1.5 rounded-full {example.hint === 'Legit'
          ? 'bg-emerald-500'
          : 'bg-red-500'}"
      ></span>
      <span class="font-mono">{example.label}</span>
      {#if level === 0}
        <span class="text-[11px] text-gray-400 dark:text-gray-500"
          >{example.hint.toLowerCase()}</span
        >
      {/if}
    </button>
  {/each}
</div>

<div
  class="mt-10 pt-6 w-full max-w-2xl border-t border-gray-200 dark:border-gray-800 flex justify-center"
>
  <!-- Short sentences on a centred line; each stays whole when the line wraps. -->
  <p
    class="max-w-xl text-center text-sm leading-relaxed text-gray-600 dark:text-gray-400 [text-wrap:balance]"
  >
    {#each POINTS as point}
      <span class="whitespace-nowrap">{point}.</span>{" "}
    {/each}
  </p>
</div>
