<script lang="ts">
  import { browser } from "$app/environment";
  import { onMount } from "svelte";

  // The redirect chain drawn as a path: where you click, each hop, where you land.
  export let chain: string[];
  /** The scanned site's registrable domain, e.g. "paypal.com". */
  export let domain: string;
  /** Overall verdict, which describes the page you land on. */
  export let verdict: string | undefined;

  type Hop = {
    url: string;
    scheme: string;
    host: string;
    rest: string;
    offSite: boolean;
    jump: boolean;
  };

  const sameSite = (host: string) =>
    host === domain || host.endsWith(`.${domain}`) || domain.endsWith(`.${host}`);

  // Full host and path, so hops that differ only by "www." or a path don't look
  // identical; "http://" is shown (dimmed) because an unencrypted hop matters.
  function split(url: string): { scheme: string; host: string; rest: string } {
    try {
      const u = new URL(url);
      return {
        scheme: u.protocol === "http:" ? "http://" : "",
        host: u.hostname,
        rest: `${u.pathname}${u.search}`,
      };
    } catch {
      return { scheme: "", host: url, rest: "" };
    }
  }

  $: hops = chain.map((url, i): Hop => {
    const { scheme, host, rest } = split(url);
    const prevHost = i > 0 ? split(chain[i - 1]).host : host;
    return {
      url,
      scheme,
      host,
      rest,
      offSite: !sameSite(host),
      jump: i > 0 && host !== prevHost && !sameSite(host),
    };
  });

  const VERDICT_DOT: Record<string, string> = {
    Safe: "bg-emerald-500",
    Suspicious: "bg-yellow-500",
    Risky: "bg-red-500",
  };

  // Hops appear one after another the first time the section opens.
  let shown = 0;
  onMount(() => {
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      shown = hops.length;
      return;
    }
    const t = setInterval(() => (shown < hops.length ? shown++ : clearInterval(t)), 140);
    return () => clearInterval(t);
  });

  function scanInNewTab(url: string) {
    if (!browser) return;
    const u = new URL(window.location.origin);
    u.searchParams.set("q", url);
    window.open(u.toString(), "_blank", "noopener");
  }
</script>

<ol class="relative" aria-label="Redirect path">
  {#each hops as hop, i}
    {@const last = i === hops.length - 1}
    <li
      class="relative flex gap-3 pb-5 last:pb-0 transition-all duration-300 {i < shown
        ? 'opacity-100 translate-y-0'
        : 'opacity-0 translate-y-1'}"
    >
      <!-- connector to the next hop -->
      {#if !last}
        <span
          class="absolute left-[5px] top-4 bottom-0 w-px {hops[i + 1].jump
            ? 'bg-yellow-500/60'
            : 'bg-gray-300 dark:bg-gray-700'}"
          aria-hidden="true"
        ></span>
      {/if}

      <span
        class="relative mt-1.5 w-[11px] h-[11px] flex-shrink-0 rounded-full {last
          ? (VERDICT_DOT[verdict ?? ''] ?? 'bg-gray-500')
          : hop.offSite
            ? 'bg-yellow-500'
            : i === 0
              ? 'bg-white dark:bg-gray-900 border-2 border-gray-400 dark:border-gray-500'
              : 'bg-gray-400 dark:bg-gray-500'}"
        aria-hidden="true"
      ></span>

      <div class="min-w-0 flex-1">
        <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">
          {#if i === 0}
            You click
          {:else if last}
            You land here
          {:else}
            Then
          {/if}
          {#if hop.jump}
            <span
              class="ml-1.5 normal-case tracking-normal font-sans text-yellow-700 dark:text-yellow-400"
              >· jumps to another site</span
            >
          {/if}
        </p>
        <p class="mt-0.5 break-all">
          <span class="font-mono text-[13px] text-gray-400 dark:text-gray-500">{hop.scheme}</span
          ><span class="font-mono text-[15px] font-medium text-gray-900 dark:text-gray-100"
            >{hop.host}</span
          ><span class="font-mono text-[13px] text-gray-400 dark:text-gray-500">{hop.rest}</span>
        </p>
        {#if hop.offSite}
          <button
            type="button"
            class="mt-1 text-xs text-gray-500 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:text-gray-900 dark:hover:text-gray-100 hover:decoration-current"
            on:click={() => scanInNewTab(hop.url)}
            aria-label={`Check ${hop.url} in a new tab`}>Check this one ↗</button
          >
        {/if}
      </div>
    </li>
  {/each}
</ol>
