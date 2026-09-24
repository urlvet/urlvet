<script lang="ts">
  import { theme } from "$lib/theme";
</script>

<button
  class="inline-flex items-center justify-center w-9 h-9 rounded-full border border-gray-300 dark:border-gray-800 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:border-gray-400 dark:hover:border-gray-600 transition-colors"
  on:click={theme.toggle}
  aria-label={$theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
  title={$theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
>
  {#if $theme === "dark"}
    <svg
      class="w-4 h-4"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8a4 4 0 100 8 4 4 0 000-8z"
      />
    </svg>
  {:else}
    <svg
      class="w-4 h-4"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
      />
    </svg>
  {/if}
</button>
