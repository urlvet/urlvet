<script lang="ts">
  import { login } from "../../admin/client";

  // Password form; hands the session token back to the page.
  export let onToken: (token: string) => void;
  export let error: string | null = null;

  let passwordInput = "";
  let loggingIn = false;

  async function submit() {
    const password = passwordInput.trim();
    if (!password || loggingIn) return;
    loggingIn = true;
    error = null;
    const res = await login(password);
    loggingIn = false;
    if (res.token) {
      passwordInput = "";
      onToken(res.token);
    } else {
      error = res.error ?? "Login failed. Please try again.";
    }
  }
</script>

<div class="flex items-center justify-center px-6 py-24 md:py-32">
  <div class="w-full max-w-sm">
    <div class="text-center mb-8">
      <p class="font-mono text-xs uppercase tracking-wider text-gray-500">Admin</p>
      <h1 class="mt-3 font-serif text-5xl text-gray-900 dark:text-gray-100">Sign in</h1>
      <p class="text-[15px] text-gray-600 dark:text-gray-400 mt-3">
        Enter your admin password to continue
      </p>
    </div>
    <form
      on:submit|preventDefault={submit}
      class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-6 space-y-4"
    >
      {#if error}
        <p
          class="text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl px-3 py-2"
        >
          {error}
        </p>
      {/if}
      <div>
        <label
          for="admin-password"
          class="block font-mono text-[11px] uppercase tracking-wider text-gray-500 mb-2"
          >Password</label
        >
        <input
          id="admin-password"
          type="password"
          bind:value={passwordInput}
          placeholder="Admin password"
          autocomplete="current-password"
          disabled={loggingIn}
          class="w-full bg-gray-50 dark:bg-gray-950 border border-gray-200 dark:border-gray-800 rounded-full px-4 py-2.5 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:outline-none focus:border-gray-400 dark:focus:border-gray-600 disabled:opacity-50"
        />
      </div>
      <button
        type="submit"
        disabled={!passwordInput.trim() || loggingIn}
        class="w-full py-2.5 rounded-full bg-gray-900 dark:bg-gray-100 text-gray-50 dark:text-gray-900 text-sm font-medium hover:bg-gray-700 dark:hover:bg-white transition-colors disabled:opacity-40"
      >
        {loggingIn ? "Signing in…" : "Sign In"}
      </button>
    </form>
  </div>
</div>
