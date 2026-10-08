<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { DomainRandomness, TyposquatResult } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let features: any;
  export let domainRandomness: DomainRandomness | undefined;
  export let typosquatResult: TyposquatResult | undefined;
</script>

{#if features}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if features.tld}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Domain Ending (TLD):</span>
            <TooltipIcon
              text="The last part of a domain name (like .com, .org, .io). It can hint at the site's trust level or purpose."
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white">.{features.tld.tld}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>High Trust Domain Ending:</span>
            <TooltipIcon
              text="Indicates whether this domain ending (TLD) is widely recognized and commonly used by highly trusted entities like government and other institutions."
            />
          </div>
          {#if features.tld.is_trusted_tld}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> Yes</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Is Risky TLD:</span>
            <TooltipIcon
              text="Some TLDs are frequently abused by scammers or malicious sites. 'Yes' suggests caution."
            />
          </div>
          {#if features.tld.is_risky_tld}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>TLD Recognized by ICANN:</span>
            <TooltipIcon
              text="ICANN oversees global domain names. Recognition means this TLD is officially managed and monitored."
            />
          </div>
          {#if features.tld.is_icann}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> Yes</span
            >
          {:else if features.tld.is_hosting_platform}
            <span class="text-yellow-600 dark:text-yellow-400 font-medium flex items-center gap-1.5"
              >Hosting Platform</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> No</span
            >
          {/if}
        </div>
      {/if}

      {#if features.url}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Uses Link Shortener:</span>
            <TooltipIcon
              text="Shows if the URL uses a shortening service (like bit.ly). Shortened links can hide a site's real destination."
            />
          </div>
          {#if features.url.url_shortener}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Uses Direct IP Address:</span>
            <TooltipIcon
              text="Some malicious sites use IP addresses instead of domain names to avoid detection. Legitimate websites rarely do this."
            />
          </div>
          {#if features.url.uses_ip}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Contains Punycode Characters:</span>
            <TooltipIcon
              text="Punycode allows special or non-Latin characters in domains (like xn--example). Scammers sometimes exploit this for lookalike attacks."
            />
          </div>
          {#if features.url.contains_punycode}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>URL Too Long:</span>
            <TooltipIcon
              text="Very long URLs can be used to obscure malicious content or trick users into trusting a fake site."
            />
          </div>
          {#if features.url.too_long}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>URL Too Deep (Many Slashes):</span>
            <TooltipIcon
              text="A URL with many nested paths may indicate redirections or hidden content, often seen in phishing attempts."
            />
          </div>
          {#if features.url.too_deep}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Has Look-Alike Letters (Homoglyph):</span>
            <TooltipIcon
              text="Detects if the domain name includes characters that look like others (e.g., 'go0gle.com'). Often used for impersonation."
            />
          </div>
          {#if features.url.has_homoglyph}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Yes</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Subdomain Count:</span>
            <TooltipIcon
              text="Shows how many subdomains (like shop.example.com) are used. Too many can hint at suspicious or temporary setups."
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white"
            >{features.url.subdomain_count}</span
          >
        </div>
      {/if}

      {#if domainRandomness && domainRandomness.entropy !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Domain Randomness (Entropy):</span>
            <TooltipIcon
              text="Measures the unpredictability of the domain name. High entropy often indicates randomly generated domains used by malware (DGAs)."
            />
          </div>
          <div class="flex items-center gap-2">
            <span class="font-medium text-gray-800 dark:text-white"
              >{domainRandomness.entropy.toFixed(2)}</span
            >
            {#if domainRandomness.entropy > 3.8}
              <span
                class="text-xs px-1.5 py-0.5 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded-sm"
                >High Entropy</span
              >
            {/if}
          </div>
        </div>
      {/if}

      {#if typosquatResult}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Typosquatting Check:</span>
            <TooltipIcon
              text="Checks if this domain closely resembles a well-known brand by looking for character substitutions, additions, or deletions (e.g. paypa1.com vs paypal.com), or if a brand name is embedded inside a longer domain."
            />
          </div>
          {#if typosquatResult.is_suspicious}
            {#if typosquatResult.is_combo_squat}
              <span class="text-red-400 font-medium flex flex-col gap-0.5">
                <span class="inline-flex items-center gap-1"
                  ><StatusIcon kind="bad" /> Combo-squatting — contains brand name</span
                >
                <span class="text-xs text-red-300 font-normal">
                  Contains "<span class="font-mono">{typosquatResult.matched_brand}</span>"
                  (official: {typosquatResult.matched_domain})
                </span>
              </span>
            {:else}
              <span class="text-red-400 font-medium flex flex-col gap-0.5">
                <span class="inline-flex items-center gap-1"
                  ><StatusIcon kind="bad" /> Typosquat, {typosquatResult.distance} character{typosquatResult.distance ===
                  1
                    ? ""
                    : "s"} away from the known brand</span
                >
                <span class="text-xs text-red-300 font-normal">
                  <span class="font-mono">{typosquatResult.matched_domain}</span>
                </span>
              </span>
            {/if}
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> No match found</span
            >
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
