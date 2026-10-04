<script lang="ts">
  import { PILL_SOLID } from "../ui/buttons";
  import { ICON } from "../ui/icons";
  import Icon from "./Icon.svelte";

  // The closing nudge at the bottom of the content pages: back to checking a link.
</script>

<section
  class="mt-24 pt-12 border-t border-gray-200 dark:border-gray-800 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6"
>
  <p class="font-serif text-3xl md:text-4xl tracking-[-0.01em] leading-tight">
    Got a link you're <span class="italic">unsure about?</span>
  </p>
  <a href="/" class="{PILL_SOLID} px-6 py-3 self-start sm:self-auto flex-shrink-0">
    Check it now
    <Icon path={ICON.arrowRight} />
  </a>
</section>
