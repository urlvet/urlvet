<script lang="ts">
  import type { Shoutout } from "../../data/shoutouts";
  import { ICON } from "../../ui/icons";
  import Icon from "../Icon.svelte";
  import SourceIcon from "./SourceIcon.svelte";

  export let item: Shoutout;
  /** "ticker": fixed-width card in the desktop strip. "slide": full-width mobile card. */
  export let variant: "ticker" | "slide" = "ticker";

  const SOURCE_LABEL: Record<Shoutout["type"], string> = {
    tweet: "Posted on X",
    linkedin: "Posted on LinkedIn",
    paper: "Research paper",
    package: "Published package",
    newsletter: "Newsletter",
    email: "Email",
  };
</script>

<a
  href={item.url}
  target="_blank"
  rel="noopener noreferrer"
  draggable="false"
  class="group flex flex-col gap-3 p-4 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 text-left cursor-pointer {variant ===
  'ticker'
    ? 'shrink-0 w-72 hover:border-gray-400 dark:hover:border-gray-700 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-black/10 dark:hover:shadow-black/30 transition-all duration-200'
    : 'w-full h-[230px]'}"
>
  <div class="flex items-center gap-3">
    <!-- avatar: initial in a ringed circle, source badge on the corner -->
    <div class="relative shrink-0">
      <div class="w-10 h-10 rounded-full p-[3px] ring-1 ring-gray-200 dark:ring-gray-800">
        <div
          class="w-full h-full rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center font-serif text-lg leading-none text-gray-700 dark:text-gray-300"
        >
          {item.name[0]}
        </div>
      </div>
      <span
        class="absolute -bottom-0.5 -right-1 w-[18px] h-[18px] rounded-full bg-white dark:bg-gray-900 ring-1 ring-gray-200 dark:ring-gray-800 flex items-center justify-center text-gray-600 dark:text-gray-400"
        title={SOURCE_LABEL[item.type]}
        aria-label={SOURCE_LABEL[item.type]}
      >
        <SourceIcon type={item.type} class="w-2.5 h-2.5" />
      </span>
    </div>
    <div class="leading-tight min-w-0">
      <p class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{item.name}</p>
      <p class="text-xs text-gray-500 truncate">
        {item.handle} <span class="ml-0.5">{item.flag}</span>
      </p>
    </div>
  </div>

  <p
    class="text-[15px] text-gray-700 dark:text-gray-300 leading-relaxed {variant === 'slide'
      ? 'line-clamp-4'
      : ''}"
  >
    "{item.text}"
  </p>
  <div class="flex items-center justify-between mt-auto">
    <p class="font-mono text-[11px] text-gray-400 dark:text-gray-500">{item.date}</p>
    <Icon
      path={ICON.external}
      class="w-4 h-4 text-gray-300 dark:text-gray-700 group-hover:text-gray-700 dark:group-hover:text-gray-300 transition-colors"
    />
  </div>
</a>
