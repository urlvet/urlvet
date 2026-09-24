<script lang="ts">
  import { VERDICT_COPY, verdictCopy } from "../../verdict";
  import { onMount } from "svelte";

  export let verdict: string | undefined;
  export let finalScore: number | undefined;
  export let unreachable = false;

  let displayedScore = 0;
  let prevScore: number | undefined = undefined;
  let ringReady = false;

  onMount(() => {
    setTimeout(() => {
      ringReady = true;
    }, 60);
  });

  $: if (finalScore !== undefined && finalScore !== prevScore) {
    prevScore = finalScore;
    const target = finalScore;
    const duration = 2000;
    const startTime = performance.now();
    const tick = (now: number) => {
      const t = Math.min((now - startTime) / duration, 1);
      const eased = 1 - Math.pow(1 - t, 3);
      displayedScore = Math.round(eased * target);
      if (t < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  }

  const STYLES: Record<
    string,
    {
      border: string;
      bg: string;
      shadow: string;
      badge: string;
      label: string;
      ringColor: string;
      scoreText: string;
      pulse: string;
    }
  > = {
    Safe: {
      border: "border-emerald-200 dark:border-emerald-900/70",
      bg: "bg-emerald-50/70 dark:bg-emerald-950/25",
      shadow: "shadow-emerald-500/10",
      badge:
        "bg-emerald-100 dark:bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-500/30",
      label: VERDICT_COPY.Safe.label,
      ringColor: "#10b981",
      scoreText: "text-emerald-600 dark:text-emerald-400",
      pulse: "verdict-safe",
    },
    Risky: {
      border: "border-red-200 dark:border-red-900/70",
      bg: "bg-red-50/70 dark:bg-red-950/25",
      shadow: "shadow-red-500/10",
      badge:
        "bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border border-red-300 dark:border-red-500/30",
      label: VERDICT_COPY.Risky.label,
      ringColor: "#ef4444",
      scoreText: "text-red-600 dark:text-red-400",
      pulse: "verdict-risky",
    },
    Suspicious: {
      border: "border-yellow-200 dark:border-yellow-900/70",
      bg: "bg-amber-50/70 dark:bg-yellow-950/25",
      shadow: "shadow-yellow-500/10",
      badge:
        "bg-yellow-100 dark:bg-yellow-500/20 text-yellow-700 dark:text-yellow-300 border border-yellow-300 dark:border-yellow-500/30",
      label: VERDICT_COPY.Suspicious.label,
      ringColor: "#eab308",
      scoreText: "text-yellow-600 dark:text-yellow-400",
      pulse: "verdict-suspicious",
    },
  };

  const R = 36;
  const CIRC = 2 * Math.PI * R;

  $: style = STYLES[verdict ?? ""] ?? STYLES.Suspicious;
  $: dashOffset = ringReady ? CIRC - ((finalScore ?? 0) / 100) * CIRC : CIRC;
  $: copy = verdictCopy(verdict);
  $: quip = copy.quip;
</script>

<div
  class={`flex flex-row items-center gap-4 p-5 sm:gap-6 sm:p-7 rounded-2xl border ${style.border} ${style.bg}`}
>
  <!-- Verdict -->
  <div class="flex-1 flex flex-col gap-1.5 min-w-0">
    <span
      class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
      >Verdict</span
    >
    <div class="flex items-center gap-3 flex-wrap">
      <span
        class="font-serif text-5xl sm:text-6xl font-normal text-gray-900 dark:text-gray-100 tracking-[-0.01em] leading-none"
        >{verdict ?? "—"}</span
      >
      <span
        class={`px-2.5 py-0.5 rounded-full font-mono text-[10px] sm:text-[11px] uppercase tracking-wider whitespace-nowrap ${style.badge}`}
      >
        {style.label}
      </span>
    </div>
    <p class="mt-1 text-[15px] sm:text-base text-gray-700 dark:text-gray-300">{quip}</p>
    {#if unreachable}
      <p class="text-[10px] sm:text-[11px] text-red-400/80">
        Site may be unreachable or returning no content.
      </p>
    {/if}
  </div>

  <!-- Circular Score Ring -->
  <div class="flex flex-col items-center gap-1 flex-shrink-0">
    <span
      class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest mb-1"
      >Trust Score</span
    >
    <div class="relative w-16 h-16 sm:w-20 sm:h-20 md:w-24 md:h-24">
      <svg class="w-full h-full -rotate-90" viewBox="0 0 88 88">
        <circle cx="44" cy="44" r={R} fill="none" stroke="var(--ring-track)" stroke-width="5" />
        <circle
          cx="44"
          cy="44"
          r={R}
          fill="none"
          stroke={style.ringColor}
          stroke-width="5"
          stroke-linecap="round"
          stroke-dasharray={CIRC}
          stroke-dashoffset={dashOffset}
          style="transition: stroke-dashoffset 2s cubic-bezier(0.16, 1, 0.3, 1)"
        />
      </svg>
      <div class="absolute inset-0 flex flex-col items-center justify-center">
        <span
          class={`${(finalScore ?? 0) >= 100 ? "text-lg sm:text-2xl" : "text-xl sm:text-2xl"} font-medium tabular-nums leading-none ${style.scoreText}`}
          >{finalScore !== undefined ? displayedScore : "—"}</span
        >
        <span class="block w-5 border-t border-gray-500 my-0.5"></span>
        <span class="text-[9px] sm:text-[10px] text-gray-600 dark:text-gray-500 font-medium"
          >100</span
        >
      </div>
    </div>
  </div>
</div>

<style>
  @keyframes pulse-safe {
    0% {
      box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.45);
    }
    100% {
      box-shadow: 0 0 0 16px rgba(16, 185, 129, 0);
    }
  }
  @keyframes pulse-risky {
    0% {
      box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.45);
    }
    100% {
      box-shadow: 0 0 0 16px rgba(239, 68, 68, 0);
    }
  }
  @keyframes pulse-suspicious {
    0% {
      box-shadow: 0 0 0 0 rgba(234, 179, 8, 0.45);
    }
    100% {
      box-shadow: 0 0 0 16px rgba(234, 179, 8, 0);
    }
  }
  .verdict-safe {
    animation: pulse-safe 0.7s ease-out 0.3s both;
  }
  .verdict-risky {
    animation: pulse-risky 0.7s ease-out 0.3s both;
  }
  .verdict-suspicious {
    animation: pulse-suspicious 0.7s ease-out 0.3s both;
  }
</style>
