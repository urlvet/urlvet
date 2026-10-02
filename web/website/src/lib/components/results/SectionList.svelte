<script lang="ts">
  import { SECTION_ICONS, computeExpanded, sectionStatuses } from "../../results/sections";
  import type { AnalyzeResult } from "../../types";
  import AccordionSection from "../AccordionSection.svelte";
  import ContentSection from "../sections/ContentSection.svelte";
  import DomainInfoSection from "../sections/DomainInfoSection.svelte";
  import InfrastructureSection from "../sections/InfrastructureSection.svelte";
  import RedirectionSection from "../sections/RedirectionSection.svelte";
  import SecuritySection from "../sections/SecuritySection.svelte";
  import ThreatIntelSection from "../sections/ThreatIntelSection.svelte";
  import URLSignalsSection from "../sections/URLSignalsSection.svelte";

  // The grouped list of detail sections, each with a one-glance status.
  export let data: AnalyzeResult;

  let expanded: Record<string, boolean> = {};
  $: expanded = computeExpanded(data);
  $: statuses = sectionStatuses(data);

  function toggle(id: string) {
    expanded = { ...expanded, [id]: !expanded[id] };
  }
</script>

<div
  data-guide="sections"
  class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800 overflow-hidden animate-fadeIn delay-300"
>
  {#if data.domain_info}
    <AccordionSection
      id="domain"
      learnMore="/how-it-works#check-domain"
      title="Domain Info"
      icon={SECTION_ICONS.domain}
      status={statuses.domain}
      expanded={expanded.domain}
      onToggle={() => toggle("domain")}
    >
      <DomainInfoSection domainInfo={data.domain_info} rank={data.features?.rank} />
    </AccordionSection>
  {/if}

  {#if data.analysis}
    <AccordionSection
      id="analysis"
      learnMore="/how-it-works#check-network"
      title="Redirection"
      icon={SECTION_ICONS.analysis}
      status={statuses.analysis}
      expanded={expanded.analysis}
      onToggle={() => toggle("analysis")}
    >
      <RedirectionSection
        analysis={data.analysis}
        domain={data.domain}
        verdict={data.result?.verdict}
      />
    </AccordionSection>
  {/if}

  {#if data.phishing}
    <AccordionSection
      id="threatintel"
      learnMore="/how-it-works#check-threats"
      title="Threat Intel"
      icon={SECTION_ICONS.threatintel}
      status={statuses.threatintel}
      expanded={expanded.threatintel}
      onToggle={() => toggle("threatintel")}
    >
      <ThreatIntelSection phishing={data.phishing} />
    </AccordionSection>
  {/if}

  {#if data.ssl_info || data.tls_info}
    <AccordionSection
      id="security"
      learnMore="/how-it-works#check-tls"
      title="Security & SSL"
      icon={SECTION_ICONS.security}
      status={statuses.security}
      expanded={expanded.security}
      onToggle={() => toggle("security")}
    >
      <SecuritySection sslInfo={data.ssl_info} tlsInfo={data.tls_info} />
    </AccordionSection>
  {/if}

  {#if data.content_data}
    <AccordionSection
      id="content"
      learnMore="/how-it-works#check-content"
      title="Page Content"
      icon={SECTION_ICONS.content}
      status={statuses.content}
      expanded={expanded.content}
      onToggle={() => toggle("content")}
    >
      <ContentSection contentData={data.content_data} />
    </AccordionSection>
  {/if}

  {#if data.features}
    <AccordionSection
      id="features"
      learnMore="/how-it-works#check-url"
      title="URL Signals"
      icon={SECTION_ICONS.features}
      status={statuses.features}
      expanded={expanded.features}
      onToggle={() => toggle("features")}
    >
      <URLSignalsSection
        features={data.features}
        domainRandomness={data.domain_randomness}
        typosquatResult={data.typosquat_result}
      />
    </AccordionSection>
  {/if}

  {#if data.infrastructure}
    <AccordionSection
      id="infrastructure"
      learnMore="/how-it-works#check-dns"
      title="Hosting & Server"
      icon={SECTION_ICONS.infrastructure}
      status={statuses.infrastructure}
      expanded={expanded.infrastructure}
      onToggle={() => toggle("infrastructure")}
    >
      <InfrastructureSection
        infrastructure={data.infrastructure}
        isHostingPlatform={!!data.features?.tld?.is_hosting_platform}
      />
    </AccordionSection>
  {/if}
</div>
