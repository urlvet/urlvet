<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { SSLInfo, TLSInfo } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let sslInfo: SSLInfo | undefined;
  export let tlsInfo: TLSInfo | undefined;
</script>

{#if sslInfo || tlsInfo}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if sslInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>SSL Support:</span>
            <TooltipIcon text="Checks if the website supports secure HTTPS connections." />
          </div>
          {#if sslInfo.HasTLS}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> Enabled</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Disabled</span
            >
          {/if}
        </div>

        {#if sslInfo.HasTLS}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Certificate Chain:</span>
              <TooltipIcon
                text="Verifies if the SSL certificate is issued by a trusted authority and the full chain is valid."
              />
            </div>
            {#if sslInfo.ChainValid}
              <span
                class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="ok" /> Valid</span
              >
            {:else}
              <span class="text-red-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="bad" /> Invalid / Self-signed</span
              >
            {/if}
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Certificate Issuer:</span>
              <TooltipIcon text="The organization that issued the SSL certificate." />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.Issuer || "-"}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Certificate Age:</span>
              <TooltipIcon
                text="How many days ago the certificate was issued. Recently issued certificates on new domains can be suspicious."
              />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.AgeDays} days</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Valid From:</span>
              <TooltipIcon text="The date this certificate first became active." />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.NotBefore}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Expiry Date:</span>
              <TooltipIcon text="When the current SSL certificate will expire." />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.NotAfter}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Certificate Risk Level:</span>
              <TooltipIcon text="Overall assessment of the certificate's technical integrity." />
            </div>
            {#if !sslInfo.IsSuspicious}
              <span
                class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="ok" /> Low Risk</span
              >
            {:else}
              <span class="text-amber-700 dark:text-amber-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="warn" /> Suspicious</span
              >
            {/if}
          </div>

          {#if sslInfo.Reasons && sslInfo.Reasons.length > 0}
            <div
              class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>Technical Warnings:</span>
                <TooltipIcon text="Specific technical reasons why this certificate is flagged." />
              </div>
              <ul class="text-xs text-amber-700 dark:text-amber-400 list-disc list-inside">
                {#each sslInfo.Reasons as reason}
                  <li>{reason}</li>
                {/each}
              </ul>
            </div>
          {/if}

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>Certificate Fingerprint:</span>
              <TooltipIcon
                text="A unique identifier (SHA-256 hash) for this specific certificate."
              />
            </div>
            <span class="font-mono text-[12px] text-gray-600 dark:text-gray-300 break-all"
              >{sslInfo.Fingerprint}</span
            >
          </div>
        {/if}
      {/if}

      {#if tlsInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>TLS Issuer (Connection):</span>
            <TooltipIcon text="The certificate issuer detected during the live connection." />
          </div>
          <span class="font-medium text-gray-800 dark:text-white">{tlsInfo.Issuer || "-"}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Hostname Match:</span>
            <TooltipIcon
              text="Ensures the certificate is actually issued for the domain you are visiting."
            />
          </div>
          {#if !tlsInfo.HostnameMismatch}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> Match</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> Mismatch</span
            >
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
