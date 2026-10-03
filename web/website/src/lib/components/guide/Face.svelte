<svelte:options namespace="svg" />

<script lang="ts">
  import FaceLines from "./FaceLines.svelte";
  import type { Mood } from "./mood";

  // Eyes, brows and mouth, shared by both characters so expressions match.
  export let lx: number; // left eye centre x
  export let rx: number; // right eye centre x
  export let y: number; // eye centre y
  export let mouthY: number;
  export let mood: Mood = "idle";
  /** Show arched brows even in calm moods (the paperclip's signature look). */
  export let brows = false;
  /** Cut a background-coloured gap behind brows and mouth, for bodies drawn through the face. */
  export let knockout = false;
  /** Eye size; the paperclip uses bigger eyes. */
  export let eyeRx = 5.6;
  export let eyeRy = 6.6;
  export let pupilR = 2.5;
  /** Brow stroke width and half-length. */
  export let browWidth = 2.2;
  export let browSpan = 5;
  /** Where to look, as a pupil offset; null lets the idle glances run. */
  export let gaze: [number, number] | null = null;

  // Pupils glance up-left while thinking or pointing (the tour target is up and to the
  // left), down at the magnifying glass while scanning, and sideways when suspicious.
  // Calm moods follow the gaze instead.
  const FIXED: Partial<Record<Mood, [number, number]>> = {
    thinking: [-1.6, -2.2],
    point: [-1.8, -1.8],
    scanning: [2, 1.8],
    suspicious: [2.2, 0.2],
  };
  $: fixed = FIXED[mood];
  $: look = fixed ?? gaze ?? [0, 0.4];
  $: wide = mood === "ooh" || mood === "alarmed" ? 0.9 : 0;
</script>

{#if mood === "happy"}
  <!-- closed, smiling eyes -->
  {#if knockout}
    <g class="knockout">
      <path d="M{lx - 5} {y + 1} Q{lx} {y - 5} {lx + 5} {y + 1}" class="ink" />
      <path d="M{rx - 5} {y + 1} Q{rx} {y - 5} {rx + 5} {y + 1}" class="ink" />
    </g>
  {/if}
  <path d="M{lx - 5} {y + 1} Q{lx} {y - 5} {lx + 5} {y + 1}" class="ink" stroke-width="2.4" />
  <path d="M{rx - 5} {y + 1} Q{rx} {y - 5} {rx + 5} {y + 1}" class="ink" stroke-width="2.4" />
{:else}
  <g class="blink">
    {#each [lx, rx] as cx}
      <ellipse
        {cx}
        cy={y}
        rx={eyeRx + wide * 0.5}
        ry={eyeRy + wide}
        class="eye"
        stroke-width="1.8"
      />
      <g class:look-around={mood === "idle" && !gaze}>
        <g class="gaze" style="transform: translate({look[0]}px, {look[1]}px)">
          <circle {cx} cy={y} r={pupilR} class="pupil" />
          <circle cx={cx + pupilR * 0.36} cy={y - pupilR * 0.4} r={pupilR * 0.32} class="glint" />
        </g>
      </g>
    {/each}
  </g>
{/if}

{#if knockout}
  <g class="knockout"><FaceLines {lx} {rx} {y} {mouthY} {mood} {brows} {browWidth} {browSpan} /></g>
{/if}
<FaceLines {lx} {rx} {y} {mouthY} {mood} {brows} {browWidth} {browSpan} />
