<svelte:options namespace="svg" />

<script lang="ts">
  import type { Mood } from "./mood";

  // Brows and mouth. Split out so Face can draw them twice: once as a
  // background-coloured cut-out, then as the visible line.
  export let lx: number;
  export let rx: number;
  export let y: number;
  export let mouthY: number;
  export let mood: Mood = "idle";
  export let brows = false;
  export let browWidth = 2.2;
  export let browSpan = 5;

  $: mx = (lx + rx) / 2;
</script>

<!-- brows -->
{#if mood === "worried" || mood === "alarmed"}
  <path
    d="M{lx - browSpan} {y - 8} L{lx + browSpan - 1} {y - 11}"
    class="ink"
    stroke-width={browWidth}
  />
  <path
    d="M{rx + browSpan} {y - 8} L{rx - browSpan + 1} {y - 11}"
    class="ink"
    stroke-width={browWidth}
  />
{:else if mood === "suspicious"}
  <!-- one brow way up, the other pressed flat: "really?" -->
  <path
    d="M{lx - browSpan} {y - 12} Q{lx} {y - 17} {lx + browSpan} {y - 12}"
    class="ink"
    stroke-width={browWidth}
  />
  <path
    d="M{rx - browSpan} {y - 8} L{rx + browSpan} {y - 9}"
    class="ink"
    stroke-width={browWidth}
  />
{:else if mood === "ooh"}
  <path
    d="M{lx - browSpan} {y - 12} Q{lx} {y - 16} {lx + browSpan} {y - 12}"
    class="ink"
    stroke-width={browWidth}
  />
  <path
    d="M{rx - browSpan} {y - 12} Q{rx} {y - 16} {rx + browSpan} {y - 12}"
    class="ink"
    stroke-width={browWidth}
  />
{:else if mood === "thinking"}
  <path
    d="M{lx - browSpan} {y - 11} Q{lx} {y - 14} {lx + browSpan} {y - 11}"
    class="ink"
    stroke-width={browWidth}
  />
  <path
    d="M{rx - browSpan} {y - 9} L{rx + browSpan} {y - 9}"
    class="ink"
    stroke-width={browWidth}
  />
{:else if brows}
  <g class="brow-raise">
    <path
      d="M{lx - browSpan} {y - 10} Q{lx} {y - 14} {lx + browSpan} {y - 10}"
      class="ink"
      stroke-width={browWidth}
    />
    <path
      d="M{rx - browSpan} {y - 10} Q{rx} {y - 14} {rx + browSpan} {y - 10}"
      class="ink"
      stroke-width={browWidth}
    />
  </g>
{/if}

<g class="mouth">
  <!-- mouth -->
  {#if mood === "happy" || mood === "wave"}
    <path
      d="M{mx - 5} {mouthY} Q{mx} {mouthY + 5} {mx + 5} {mouthY}"
      class="ink"
      stroke-width="2.2"
    />
  {:else if mood === "ooh"}
    <ellipse cx={mx} cy={mouthY + 1} rx="2.2" ry="2.8" class="ink" stroke-width="2" />
  {:else if mood === "suspicious"}
    <path d="M{mx - 3} {mouthY + 1} L{mx + 4} {mouthY - 0.5}" class="ink" stroke-width="2.2" />
  {:else if mood === "worried" || mood === "alarmed"}
    <path
      d="M{mx - 5} {mouthY + 2} Q{mx - 2.5} {mouthY - 1} {mx} {mouthY + 1} Q{mx + 2.5} {mouthY +
        3} {mx + 5} {mouthY}"
      class="ink"
      stroke-width="2"
    />
  {:else if mood === "thinking" || mood === "scanning"}
    <path d="M{mx - 2} {mouthY + 1} L{mx + 3} {mouthY}" class="ink" stroke-width="2.2" />
  {:else}
    <path
      d="M{mx - 3.5} {mouthY} Q{mx} {mouthY + 2.5} {mx + 3.5} {mouthY}"
      class="ink"
      stroke-width="2.2"
    />
  {/if}
</g>
