<script lang="ts">
  import { page } from "$app/stores";
</script>

<svelte:head>
  <title>{$page.status} — url.vet</title>
  <meta name="robots" content="noindex" />
</svelte:head>

<div class="flex items-center justify-center px-6 py-24 md:py-32">
  <div class="text-center max-w-sm w-full">
    <p class="font-serif italic text-9xl leading-none text-gray-900 dark:text-gray-100 mb-4">
      {$page.status}
    </p>

    <div class="w-12 h-px bg-gray-300 dark:bg-gray-700 mx-auto mb-6"></div>

    <h1 class="font-serif text-3xl text-gray-900 dark:text-gray-100 mb-3">
      {$page.status === 404 ? "Page not found" : "Something went wrong"}
    </h1>

    <p class="text-[15px] text-gray-600 dark:text-gray-400 mb-10 leading-relaxed">
      {$page.status === 404
        ? "This page doesn't exist or was moved. Head back and scan a link."
        : ($page.error?.message ?? "An unexpected error occurred.")}
    </p>

    <a
      href="/"
      class="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-gray-900 dark:bg-gray-100 text-gray-50 dark:text-gray-900 text-sm font-medium hover:bg-gray-700 dark:hover:bg-white transition-colors"
    >
      <svg
        class="w-3.5 h-3.5"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      Back to url.vet
    </a>
  </div>
</div>
