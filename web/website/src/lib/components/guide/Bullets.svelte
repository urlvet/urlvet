<script lang="ts">
  // A short dotted list inside a speech bubble.
  export let items: readonly string[];
</script>

<ul class="mt-1.5 space-y-1">
  {#each items as item}
    <li class="flex gap-2">
      <span class="mt-[0.6em] w-1 h-1 rounded-full bg-gray-400 flex-shrink-0"></span>
      <span>{item}</span>
    </li>
  {/each}
</ul>
