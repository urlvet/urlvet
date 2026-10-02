<script lang="ts">
  import { onMount, tick } from "svelte";
  import { EXAMPLES } from "../../data/examples";
  // Example links to try, and the short list of promises under the search bar.
  export let onTry: (url: string) => void;

  const POINTS = ["Open source", "No signup", "Explains every verdict", "See it before you click"];

  // The chips should sit in two rows. Screen width, fonts and phone text-size
  // settings all change how wide they are, so measure instead of guessing a
  // breakpoint: drop the legit/fake labels, then the "Try" label, until they fit.
  // 0 = everything, 1 = no labels, 2 = no labels and no "Try".
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
  class="mt-10 pt-6 w-full max-w-2xl border-t border-gray-200 dark:border-gray-800 grid grid-cols-2 gap-x-4 gap-y-2.5 sm:flex sm:flex-wrap sm:justify-center sm:gap-x-6 sm:gap-y-2"
>
  {#each POINTS as point}
    <span
      class="flex items-start gap-2 text-left leading-snug text-[13px] min-[360px]:text-sm text-gray-600 dark:text-gray-400"
    >
      <span
        class="mt-[0.6em] w-1 h-1 flex-shrink-0 rounded-full bg-accent-light dark:bg-accent-dark"
      ></span>
      {point}
    </span>
  {/each}
</div>
