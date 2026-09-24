<script lang="ts">
  // One faint warm light from the top of the landing page.
  const TINT = "214 170 120";
</script>

<div
  class="backdrop pointer-events-none absolute inset-x-0 -top-24 h-[44rem]"
  style="--tint: {TINT}"
  aria-hidden="true"
></div>

<style>
  .backdrop {
    background: radial-gradient(
      ellipse 55% 45% at 50% 0%,
      rgb(var(--tint) / 0.07),
      transparent 70%
    );
    transition: background 0.8s ease;
  }
  :global(.dark) .backdrop {
    background: radial-gradient(
      ellipse 55% 45% at 50% 0%,
      rgb(var(--tint) / 0.09),
      transparent 70%
    );
  }
</style>
