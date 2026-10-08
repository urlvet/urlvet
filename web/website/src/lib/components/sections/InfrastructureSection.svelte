<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let infrastructure: any;
  export let isHostingPlatform: boolean = false;
</script>

{#if infrastructure}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if infrastructure.ip_addresses?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>
              Server IP Address{infrastructure.ip_addresses.length > 1 ? "es" : ""}:
            </span>
            <TooltipIcon
              text="The actual network address where the website is hosted. Each IP points to a specific physical or cloud server."
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.ip_addresses as ip}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{ip}</span
              >
            {/each}
          </div>
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>DNS Configuration:</span>
          <TooltipIcon
            text="Nameservers control where your domain points. They act like the internet's 'address book', linking your domain name to the right hosting provider."
          />
        </div>
        {#if infrastructure.nameservers_valid}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="ok" /> Detected</span
          >
        {:else if isHostingPlatform}
          <span class="text-gray-500 dark:text-gray-400 font-medium flex items-center gap-1.5"
            >Managed by hosting platform</span
          >
        {:else}
          <span class="text-red-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="bad" /> Not Detected</span
          >
        {/if}
      </div>

      {#if infrastructure.ns_hosts?.length > 0}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>Nameserver Hosts:</span>
            <TooltipIcon
              text="The servers responsible for managing your domain's DNS settings. These typically belong to your registrar or hosting provider."
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.ns_hosts as ns_host}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{ns_host}</span
              >
            {/each}
          </div>
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Email Server Records:</span>
          <TooltipIcon
            text="MX records define where emails for this domain are delivered, essential for sending and receiving mail securely."
          />
        </div>
        {#if infrastructure.mx_records_valid}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="ok" /> Detected</span
          >
        {:else if isHostingPlatform}
          <span class="text-gray-500 dark:text-gray-400 font-medium flex items-center gap-1.5"
            >Not applicable</span
          >
        {:else}
          <span class="text-red-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="bad" /> Not Detected</span
          >
        {/if}
      </div>

      {#if infrastructure.mx_hosts?.length > 0}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>MX Hosts:</span>
            <TooltipIcon
              text="The mail servers responsible for handling your domain's email traffic."
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.mx_hosts as mx_host}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{mx_host}</span
              >
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
