<script lang="ts">
  import { ICON } from "../../ui/icons";
  import Icon from "../Icon.svelte";

  export let message: string;
</script>

<div class="max-w-xl mx-auto flex flex-col items-center gap-4 py-10 text-center">
  <div class="w-12 h-12 rounded-full bg-red-500/10 flex items-center justify-center">
    <Icon path={ICON.warning} strokeWidth={1.5} class="w-6 h-6 text-red-500" />
  </div>
  <div>
    <p class="font-serif text-2xl text-gray-900 dark:text-gray-100 mb-1">Scan failed</p>
    <p class="text-gray-600 dark:text-gray-400 text-sm max-w-sm">{message}</p>
    <p class="text-gray-500 text-xs mt-2">
      The site may be unreachable, blocking automated access, or the URL may be invalid.
    </p>
  </div>
</div>
