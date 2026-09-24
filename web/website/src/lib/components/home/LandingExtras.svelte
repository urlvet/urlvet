<script lang="ts">
  import { EXAMPLES } from "../../data/examples";
  // Example links to try, and the short list of promises under the search bar.
  export let onTry: (url: string) => void;

  const POINTS = ["Open source", "No signup", "Explains every verdict", "See it before you click"];
</script>

<div data-guide="examples" class="mt-5 flex flex-wrap justify-center items-center gap-2">
  <span class="hidden min-[360px]:inline text-xs text-gray-500 mr-1">Try</span>
  {#each EXAMPLES as example}
    <button
      type="button"
      data-example={example.url}
      on:click={() => onTry(example.url)}
      title={example.hint === "Legit"
        ? "The real site"
        : "A lookalike, spelled with letters from another alphabet"}
      class="inline-flex items-center gap-1.5 min-[360px]:gap-2 px-2.5 min-[360px]:px-3.5 py-1.5 rounded-full border border-gray-300 dark:border-gray-800 hover:border-gray-400 dark:hover:border-gray-600 text-gray-700 dark:text-gray-300 text-xs transition-colors"
    >
      <span
        class="w-1.5 h-1.5 rounded-full {example.hint === 'Legit'
          ? 'bg-emerald-500'
          : 'bg-red-500'}"
      ></span>
      <span class="font-mono">{example.label}</span>
      <span class="hidden min-[360px]:inline text-[11px] text-gray-400 dark:text-gray-500"
        >{example.hint.toLowerCase()}</span
      >
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
