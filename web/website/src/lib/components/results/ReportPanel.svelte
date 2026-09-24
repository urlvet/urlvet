<script lang="ts">
  import { onDestroy } from "svelte";
  import { fly, slide } from "svelte/transition";
  import { api } from "../../api";
  import { PILL_SOLID } from "../../ui/buttons";
  import { ICON } from "../../ui/icons";
  import Icon from "../Icon.svelte";

  // "Report this result" form, plus the toast shown after a successful send.
  export let open = false;
  export let url: string;
  export let domain: string;
  export let verdict: string | undefined;
  export let score: number | undefined;

  type ReportVerdict = "Safe" | "Suspicious" | "Risky";

  const OPTIONS: { value: ReportVerdict; active: string }[] = [
    {
      value: "Safe",
      active:
        "bg-emerald-100 dark:bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 border-emerald-300 dark:border-emerald-500/30",
    },
    {
      value: "Suspicious",
      active:
        "bg-yellow-100 dark:bg-yellow-500/20 text-yellow-700 dark:text-yellow-300 border-yellow-300 dark:border-yellow-500/30",
    },
    {
      value: "Risky",
      active:
        "bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border-red-300 dark:border-red-500/30",
    },
  ];

  let expected: ReportVerdict | "" = "";
  let comment = "";
  let status: "idle" | "sending" | "error" = "idle";
  let errorMessage = "";
  let toast = false;
  let toastTimer: ReturnType<typeof setTimeout> | undefined;

  // A new scan result starts with a clean form.
  $: (url, reset());

  function reset() {
    open = false;
    expected = "";
    comment = "";
    status = "idle";
    errorMessage = "";
  }

  function toggleExpected(v: ReportVerdict) {
    expected = expected === v ? "" : v;
  }

  async function submit() {
    if (status === "sending") return;
    status = "sending";
    const res = await api.report({
      url,
      verdict,
      score,
      expected_verdict: expected || undefined,
      comment: comment.trim() || undefined,
    });
    if (res.error) {
      status = "error";
      errorMessage = res.error;
      return;
    }
    reset();
    toast = true;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => (toast = false), 2000);
  }

  onDestroy(() => clearTimeout(toastTimer));
</script>

{#if open}
  <div
    id="report-panel"
    transition:slide={{ duration: 200 }}
    class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900"
  >
    <form class="p-4 sm:p-5 flex flex-col gap-4" on:submit|preventDefault={submit}>
      <div class="flex items-start justify-between gap-3">
        <div>
          <p class="font-serif text-2xl text-gray-900 dark:text-gray-100">Report this result</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
            We rated <span class="font-mono text-gray-700 dark:text-gray-300">{domain}</span>
            as
            <span class="font-semibold text-gray-700 dark:text-gray-300">{verdict}</span>. Tell us
            what's off.
          </p>
        </div>
        <button
          type="button"
          class="p-1.5 -m-1.5 rounded-md text-gray-400 dark:text-gray-500 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-800 transition-colors"
          aria-label="Close report"
          on:click={() => (open = false)}
        >
          <Icon path={ICON.close} />
        </button>
      </div>

      <div class="flex flex-col gap-2">
        <span
          class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
          >Should be</span
        >
        <div class="flex flex-wrap gap-2" role="group" aria-label="Expected verdict">
          {#each OPTIONS as opt}
            <button
              type="button"
              aria-pressed={expected === opt.value}
              on:click={() => toggleExpected(opt.value)}
              class="px-3.5 py-1.5 rounded-full border text-xs font-semibold transition-all duration-150 active:scale-95 {expected ===
              opt.value
                ? opt.active
                : 'bg-white dark:bg-gray-900 border-gray-300 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-900 dark:hover:text-gray-200'}"
            >
              {opt.value}
            </button>
          {/each}
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <label
          for="report-comment"
          class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
          >Details</label
        >
        <textarea
          id="report-comment"
          bind:value={comment}
          maxlength="1000"
          rows="3"
          placeholder="e.g. This is my bank's official site"
          class="w-full resize-none rounded-xl bg-gray-50 dark:bg-gray-950 border border-gray-200 dark:border-gray-800 px-4 py-3 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:border-gray-400 dark:focus:border-gray-600 transition-colors"
        ></textarea>
      </div>

      <div class="flex flex-col-reverse sm:flex-row sm:items-center justify-between gap-3">
        {#if status === "error"}
          <p class="text-xs text-red-600 dark:text-red-400">
            Couldn't send report: {errorMessage}
          </p>
        {:else}
          <span class="text-xs text-gray-400 dark:text-gray-500 tabular-nums"
            >{comment.length}/1000</span
          >
        {/if}
        <button
          type="submit"
          disabled={status === "sending"}
          class="self-end {PILL_SOLID} px-5 py-2.5"
        >
          {#if status === "sending"}
            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24" aria-hidden="true">
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
              ></path>
            </svg>
            Sending…
          {:else}
            Send report
          {/if}
        </button>
      </div>
    </form>
  </div>
{/if}

{#if toast}
  <div
    role="status"
    aria-live="polite"
    transition:fly={{ y: 16, duration: 250 }}
    class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 pl-3 pr-5 py-3 rounded-full border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shadow-xl shadow-black/10"
    style="margin-bottom: env(safe-area-inset-bottom, 0px)"
  >
    <span
      class="flex items-center justify-center w-7 h-7 rounded-full bg-emerald-100 dark:bg-emerald-500/20"
    >
      <Icon
        path={ICON.check}
        strokeWidth={2.5}
        class="w-4 h-4 text-emerald-600 dark:text-emerald-400"
      />
    </span>
    <div>
      <p class="text-sm font-medium text-gray-900 dark:text-gray-100">Report sent</p>
      <p class="text-xs text-gray-500 dark:text-gray-400">Thanks for helping improve url.vet</p>
    </div>
  </div>
{/if}
