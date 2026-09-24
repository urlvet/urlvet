<script lang="ts">
  // Outline SVG icon from a 24x24 path (see $lib/ui/icons).
  export let path: string;
  let className = "w-4 h-4";
  export { className as class };
  export let strokeWidth = 2;
</script>

<svg
  class={className}
  fill="none"
  stroke="currentColor"
  stroke-width={strokeWidth}
  viewBox="0 0 24 24"
  aria-hidden="true"
>
  <path stroke-linecap="round" stroke-linejoin="round" d={path} />
</svg>
