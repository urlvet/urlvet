<script lang="ts">
  import { onMount } from "svelte";
  import { UnauthorizedError, createAdminClient } from "../../lib/admin/client";
  import type { Guard, Tab } from "../../lib/admin/types";
  import AdminLogin from "../../lib/components/admin/AdminLogin.svelte";
  import AdminTabs from "../../lib/components/admin/AdminTabs.svelte";
  import CacheTab from "../../lib/components/admin/CacheTab.svelte";
  import ErrorsTab from "../../lib/components/admin/ErrorsTab.svelte";
  import OverviewTab from "../../lib/components/admin/OverviewTab.svelte";
  import RecentTab from "../../lib/components/admin/RecentTab.svelte";
  import ReportsTab from "../../lib/components/admin/ReportsTab.svelte";

  let token: string | null = null;
  let loginError: string | null = null;
  let activeTab: Tab = "overview";
  let globalError: string | null = null;
  let errorCount = 0;
  let reportCount = 0;
  // Bumping this remounts the active tab, which reloads its data.
  let refreshKey = 0;

  $: client = token ? createAdminClient(token) : null;
  let updatedAt = "";

  // Tab badges are filled up front, so they show before a tab is opened.
  async function loadCounts() {
    if (!client) return;
    const [reports, errors] = await Promise.all([guard(client.reports()), guard(client.errors())]);
    if (reports) reportCount = reports.length;
    if (errors) errorCount = errors.length;
    updatedAt = new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }

  $: if (client) loadCounts();

  function setToken(t: string) {
    token = t;
    sessionStorage.setItem("admin_token", t);
  }

  function clearAuth(message?: string) {
    token = null;
    sessionStorage.removeItem("admin_token");
    loginError = message ?? null;
  }

  function logout() {
    clearAuth();
    errorCount = 0;
    reportCount = 0;
  }

  // Every tab request goes through here: expired sessions sign out, other failures show a banner.
  const guard: Guard = async (request) => {
    try {
      return await request;
    } catch (e) {
      if (e instanceof UnauthorizedError) clearAuth("Invalid token. Please sign in again.");
      else globalError = e instanceof Error ? e.message : "Request failed";
      return undefined;
    }
  };

  function switchTab(tab: Tab) {
    globalError = null;
    if (tab === activeTab) refreshKey++;
    activeTab = tab;
    window.scrollTo({ top: 0 });
  }

  function refresh() {
    globalError = null;
    refreshKey++;
    loadCounts();
  }

  onMount(() => {
    token = sessionStorage.getItem("admin_token");
  });
</script>

<svelte:head>
  <title>Admin — url.vet</title>
  <meta name="robots" content="noindex, nofollow, noarchive" />
</svelte:head>

{#if !token || !client}
  <AdminLogin onToken={setToken} bind:error={loginError} />
{:else}
  <div class="text-gray-900 dark:text-gray-200">
    <!-- Header: title and actions, then a sticky tab bar -->
    <div
      class="max-w-5xl mx-auto px-4 min-[360px]:px-6 pt-8 md:pt-12 flex items-end justify-between gap-4"
    >
      <div>
        <p class="font-mono text-xs uppercase tracking-wider text-gray-500">Admin</p>
        <h1 class="mt-2 font-serif text-4xl md:text-5xl text-gray-900 dark:text-gray-100">
          Dashboard
        </h1>
      </div>
      <div class="flex items-center gap-1 pb-1">
        {#if updatedAt}
          <span class="hidden sm:inline mr-2 font-mono text-[11px] text-gray-400"
            >Updated {updatedAt}</span
          >
        {/if}
        <button
          on:click={refresh}
          aria-label="Refresh"
          title="Refresh"
          class="p-2 rounded-full text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-200/60 dark:hover:bg-gray-800/60 transition-colors"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
        <button
          on:click={logout}
          class="px-3 py-1.5 rounded-full text-sm text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-200/60 dark:hover:bg-gray-800/60 transition-colors"
        >
          Sign out
        </button>
      </div>
    </div>

    <div
      class="sticky top-[env(safe-area-inset-top,0px)] z-30 mt-6 bg-gray-50 dark:bg-gray-950 border-b border-gray-200 dark:border-gray-800"
    >
      <div class="max-w-5xl mx-auto px-4 min-[360px]:px-6">
        <AdminTabs {activeTab} {errorCount} {reportCount} onSelect={switchTab} />
      </div>
    </div>

    <div class="max-w-5xl mx-auto px-4 min-[360px]:px-6 py-8">
      <!-- Global error -->
      {#if globalError}
        <div
          class="flex items-center gap-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-4 py-3 text-sm text-red-700 dark:text-red-300 mb-6"
        >
          <svg class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
              clip-rule="evenodd"
            />
          </svg>
          {globalError}
          <button
            on:click={() => (globalError = null)}
            class="ml-auto text-red-500 dark:text-red-500 hover:text-red-700 dark:hover:text-red-300"
            >✕</button
          >
        </div>
      {/if}

      {#key `${activeTab}:${refreshKey}`}
        {#if activeTab === "overview"}
          <OverviewTab {client} {guard} />
        {:else if activeTab === "recent"}
          <RecentTab {client} {guard} />
        {:else if activeTab === "errors"}
          <ErrorsTab {client} {guard} bind:count={errorCount} />
        {:else if activeTab === "reports"}
          <ReportsTab {client} {guard} bind:count={reportCount} />
        {:else if activeTab === "cache"}
          <CacheTab {client} {guard} />
        {/if}
      {/key}
    </div>
  </div>
{/if}
