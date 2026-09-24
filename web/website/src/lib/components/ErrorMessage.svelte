<script lang="ts">
  export let message: string;
</script>

<div class="max-w-2xl mx-auto mb-6" role="alert">
  <div
    class="rounded-2xl border border-red-200 dark:border-red-900/70 bg-red-50/70 dark:bg-red-950/25 p-4"
  >
    <div class="flex">
      <div class="flex-shrink-0">
        <svg
          class="h-5 w-5 text-red-500"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <div class="ml-3">
        <p class="text-sm text-red-700 dark:text-red-300">{message}</p>
      </div>
    </div>
  </div>
</div>
