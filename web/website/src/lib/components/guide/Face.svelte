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

  // Pupils glance up-left while thinking or pointing (the tour target is up and to the left).
  $: look = mood === "thinking" ? [-1.6, -2.2] : mood === "point" ? [-1.8, -1.8] : [0, 0.4];
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
      <ellipse {cx} cy={y} rx={eyeRx} ry={eyeRy} class="eye" stroke-width="1.8" />
      <g class:look-around={mood === "idle"}>
        <circle cx={cx + look[0]} cy={y + look[1]} r={pupilR} class="pupil" />
        <circle
          cx={cx + look[0] + pupilR * 0.36}
          cy={y + look[1] - pupilR * 0.4}
          r={pupilR * 0.32}
          class="glint"
        />
      </g>
    {/each}
  </g>
{/if}

{#if knockout}
  <g class="knockout"><FaceLines {lx} {rx} {y} {mouthY} {mood} {brows} {browWidth} {browSpan} /></g>
{/if}
<FaceLines {lx} {rx} {y} {mouthY} {mood} {brows} {browWidth} {browSpan} />
