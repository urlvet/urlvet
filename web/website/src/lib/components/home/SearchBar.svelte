<script lang="ts">
  import { ICON } from "../../ui/icons";
  import { formatUrl, stripTrackers } from "../../utils";
  import Icon from "../Icon.svelte";

  // The pill search bar: URL input, Paste, Scan, plus tracker clean-up and length hints.
  export let input = "";
  export let loading = false;
  export let formError: string | null = null;
  export let isLanding = true;
  export let onSubmit: (value: string) => void;
  export let onPaste: (() => void) | undefined = undefined;

  let justPasted = false;
  let trackerCopied = false;
  let trackerDismissed = false;

  $: trackers = input ? stripTrackers(formatUrl(input)) : { cleaned: "", removed: [] as string[] };
  $: hasTrackers = trackers.removed.length > 0;
  $: if (input === "") trackerDismissed = false;

  async function pasteFromClipboard() {
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        input = text.trim();
        formError = null;
        onPaste?.();
        justPasted = true;
        setTimeout(() => (justPasted = false), 1500);
      }
    } catch {
      /* clipboard access denied */
    }
  }

  async function copyCleanUrl() {
    try {
      await navigator.clipboard.writeText(trackers.cleaned);
      trackerCopied = true;
      setTimeout(() => {
        trackerCopied = false;
        trackerDismissed = true;
      }, 1200);
    } catch {}
  }
</script>

<form
  class={isLanding ? "w-full max-w-2xl" : "w-full max-w-4xl mx-auto px-4"}
  on:submit|preventDefault={() => onSubmit(input)}
>
  <div
    data-guide="search"
    class="flex items-center gap-1 rounded-full border bg-white dark:bg-gray-900 p-1.5 pl-5 transition-colors focus-within:border-gray-500 dark:focus-within:border-gray-600 {formError
      ? 'border-red-400 dark:border-red-600/70'
      : 'border-gray-300 dark:border-gray-800'}"
  >
    <input
      id="url-input"
      type="text"
      class="flex-1 min-w-0 bg-transparent py-2.5 text-base text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none"
      placeholder="Paste a link"
      bind:value={input}
      on:input={() => {
        if (formError) formError = null;
      }}
      autocomplete="url"
      inputmode="url"
      aria-invalid={formError ? "true" : undefined}
      aria-describedby={formError ? "url-error" : undefined}
      required
    />
    <button
      type="button"
      data-guide="paste"
      on:click={pasteFromClipboard}
      class="flex-shrink-0 inline-flex items-center gap-1.5 px-3 py-2 rounded-full text-xs font-medium transition-colors {justPasted
        ? 'text-emerald-700 dark:text-emerald-400'
        : 'text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-800'}"
      aria-label="Paste from clipboard"
      title="Paste from clipboard"
    >
      {#if justPasted}
        <Icon path={ICON.check} strokeWidth={2.5} class="w-3.5 h-3.5" />
        <span class="hidden sm:inline">Pasted</span>
      {:else}
        <Icon path={ICON.clipboard} class="w-3.5 h-3.5" />
        <span class="hidden sm:inline">Paste</span>
      {/if}
    </button>
    <button
      type="submit"
      class="flex-shrink-0 inline-flex items-center justify-center gap-2 px-5 sm:px-6 py-2.5 rounded-full bg-gray-900 dark:bg-gray-100 text-gray-50 dark:text-gray-900 text-sm font-medium hover:bg-gray-700 dark:hover:bg-white transition-colors disabled:opacity-60 focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-gray-400 dark:focus-visible:ring-offset-gray-900"
      disabled={loading}
      aria-busy={loading}
      aria-label={loading ? "Scanning URL, please wait" : "Scan"}
    >
      {#if loading}
        <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24" aria-hidden="true">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
          ></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
          ></path>
        </svg>
        Scanning
      {:else}
        Scan
        <Icon path={ICON.arrowRight} />
      {/if}
    </button>
  </div>

  {#if formError}
    <p
      id="url-error"
      class="text-red-600 dark:text-red-400 text-xs mt-2 pl-5 text-left"
      role="alert"
    >
      {formError}
    </p>
  {:else if input.length > 1848 && !/\s/.test(input.trim())}
    <p
      class="text-[11px] mt-1.5 pr-5 text-right {input.length >= 2048
        ? 'text-red-500'
        : 'text-gray-500'}"
    >
      {2048 - input.length} chars remaining
    </p>
  {/if}

  {#if hasTrackers && !trackerDismissed}
    <div class="mt-2 px-1 text-xs space-y-1.5 text-center">
      <p class="text-gray-500">
        <span class="text-gray-600 dark:text-gray-400 font-medium"
          >{trackers.removed.length} tracker{trackers.removed.length > 1 ? "s" : ""} found —</span
        >
        <span class="break-all"> {trackers.removed.join(", ")}</span>
      </p>
      <p class="text-gray-500">Click to copy the clean link</p>
      <button
        type="button"
        on:click={copyCleanUrl}
        class="w-full font-mono text-xs break-all transition-colors cursor-copy {trackerCopied
          ? 'text-emerald-600 dark:text-emerald-400'
          : 'text-gray-500 hover:text-gray-800 dark:hover:text-gray-300'}"
      >
        {trackerCopied ? "Copied!" : trackers.cleaned}
      </button>
    </div>
  {/if}
</form>
