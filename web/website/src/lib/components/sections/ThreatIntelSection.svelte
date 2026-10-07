<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { FeedMatch, PhishingResult, WebRiskResult } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let phishing: PhishingResult | undefined;
  export let feeds: FeedMatch | undefined = undefined;
  export let webRisk: WebRiskResult | undefined = undefined;
  /** Which Google service gave the verdict in webRisk. */
  export let googleSource = "Google Web Risk";

  const THREAT_NAMES: Record<string, string> = {
    SOCIAL_ENGINEERING: "phishing",
    MALWARE: "malware",
    UNWANTED_SOFTWARE: "unwanted software",
    POTENTIALLY_HARMFUL_APPLICATION: "a harmful app",
  };
  $: listedBy = feeds?.sources?.join(" and ") ?? "";
  $: googleThreats = (webRisk?.threat_types ?? []).map((t) => THREAT_NAMES[t] ?? t).join(" and ");
</script>

{#if phishing || feeds || webRisk}
  <section class="p-4 sm:p-5 space-y-4">
    {#if feeds || webRisk}
      <div
        class="{feeds?.listed || webRisk?.listed
          ? 'border-red-700/60 bg-red-950/20'
          : 'border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/30'} rounded-lg border p-4"
      >
        <div class="flex items-center gap-2 mb-3">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">Threat lists</span>
          <TooltipIcon
            text="Lists of reported phishing and malware links, kept up to date on our server and checked on every scan."
          />
        </div>
        {#if feeds}
          {#if !feeds.listed}
            <p
              class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
            >
              <StatusIcon kind="ok" />
              {feeds.checked?.length
                ? `Not listed by ${feeds.checked.join(" or ")}.`
                : "Not on any phishing or malware list we hold."}
            </p>
          {:else if feeds.match === "url"}
            <p class="text-sm font-semibold text-red-400 flex items-center gap-2">
              <StatusIcon kind="bad" /> This link is listed by {listedBy}.
            </p>
          {:else}
            <p class="text-sm font-semibold text-red-400 flex items-center gap-2">
              <StatusIcon kind="bad" /> Other pages on this site are listed by {listedBy}.
            </p>
          {/if}
        {/if}
        {#if webRisk}
          <p
            class="text-sm font-medium flex items-center gap-2 mt-2 {webRisk.listed
              ? 'text-red-400 font-semibold'
              : 'text-emerald-700 dark:text-emerald-400'}"
          >
            {#if webRisk.listed}
              <StatusIcon kind="bad" /> Google lists this link as {googleThreats || "dangerous"}.
            {:else}
              <StatusIcon kind="ok" /> Not listed by {googleSource}.
            {/if}
          </p>
          {#if webRisk.listed}
            <!-- Required by the Safe Browsing and Web Risk terms wherever Google's verdict is shown. -->
            <p class="mt-2 text-xs leading-snug text-gray-500 dark:text-gray-400">
              <a
                href="https://developers.google.com/safe-browsing/v4/advisory"
                target="_blank"
                rel="noopener noreferrer"
                class="underline underline-offset-2">Advisory provided by Google</a
              >. Google works to provide the most accurate and up-to-date information about unsafe
              web resources. However, Google cannot guarantee that its information is comprehensive
              and error-free: some risky sites may not be identified, and some safe sites may be
              identified in error.
            </p>
          {/if}
        {/if}
      </div>
    {/if}

    {#if phishing}
      <div
        class="{phishing.in_database && phishing.valid
          ? 'border-red-700/60 bg-red-950/20'
          : 'border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/30'} rounded-lg border p-4"
      >
        <!-- Block header -->
        <div class="flex items-center gap-2 mb-3">
          <img src="https://phishtank.com/favicon.ico" alt="" class="w-4 h-4" />
          <span class="text-sm font-semibold text-gray-900 dark:text-white">PhishTank Lookup</span>
          <TooltipIcon text="Community-verified phishing database check. Runs on every scan." />
        </div>

        {#if !phishing.in_database}
          <!-- Not in database at all — cleanest result -->
          <p
            class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
          >
            <StatusIcon kind="ok" /> Not found in PhishTank database.
          </p>
        {:else if phishing.verified && !phishing.valid}
          <!-- In database, reviewed, and confirmed NOT phishing -->
          <p
            class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
          >
            <StatusIcon kind="ok" /> Reviewed by PhishTank community, confirmed not phishing.
          </p>
          {#if phishing.phish_id}
            <p class="text-xs text-gray-500 mt-1">
              Report
              {#if phishing.phish_detail_page}
                <a
                  href={phishing.phish_detail_page}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-gray-400 hover:text-gray-600 dark:text-gray-300 underline"
                  >#{phishing.phish_id}</a
                >
              {:else}
                #{phishing.phish_id}
              {/if}
            </p>
          {/if}
        {:else}
          <!-- In database and valid=true (phishing) or not yet reviewed -->
          <div
            class="space-y-0 divide-y divide-gray-100 dark:divide-gray-700/50 text-sm text-gray-800 dark:text-gray-200"
          >
            <!-- in_database -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 first:pt-0"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>In Database:</span>
                <TooltipIcon text="This URL has been submitted to or reported in PhishTank." />
              </div>
              <span class="text-amber-700 dark:text-amber-400 font-medium">Yes, on record</span>
            </div>

            <!-- valid — the actual phishing signal -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>Is Phishing:</span>
                <TooltipIcon
                  text="Whether this URL is confirmed as a phishing site. true = phishing, false = not phishing."
                />
              </div>
              {#if phishing.valid}
                <span class="font-semibold text-red-400">Yes, phishing confirmed</span>
              {:else}
                <span class="font-semibold text-gray-400">No</span>
              {/if}
            </div>

            <!-- verified — reviewed status -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>Community Reviewed:</span>
                <TooltipIcon
                  text="Whether the PhishTank community has reviewed this report (true = reviewed, false = still pending)."
                />
              </div>
              {#if phishing.verified}
                <span class="text-gray-900 dark:text-white font-medium">Yes</span>
              {:else}
                <span class="text-amber-700 dark:text-amber-400 font-medium"
                  >No, awaiting review</span
                >
              {/if}
            </div>

            <!-- verified_at -->
            {#if phishing.verified_at}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>Reviewed At:</span>
                  <TooltipIcon text="When the PhishTank community reviewed this report." />
                </div>
                <span class="font-medium text-gray-800 dark:text-white">{phishing.verified_at}</span
                >
              </div>
            {/if}

            <!-- target -->
            {#if phishing.target}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>Impersonation Target:</span>
                  <TooltipIcon text="The brand or service this phishing URL is impersonating." />
                </div>
                <span class="font-medium text-gray-800 dark:text-white">{phishing.target}</span>
              </div>
            {/if}

            <!-- phish_id + phish_detail_page -->
            {#if phishing.phish_id}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 last:pb-0"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>PhishTank Report:</span>
                  <TooltipIcon text="Unique PhishTank ID. Click to view the full report page." />
                </div>
                {#if phishing.phish_detail_page}
                  <a
                    href={phishing.phish_detail_page}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="font-mono text-sm text-gray-900 dark:text-gray-100 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:decoration-current"
                    >#{phishing.phish_id}</a
                  >
                {:else}
                  <span class="font-mono text-sm text-gray-600 dark:text-gray-300"
                    >#{phishing.phish_id}</span
                  >
                {/if}
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}
  </section>
{/if}
