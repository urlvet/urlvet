<script lang="ts">
  import { tick } from "svelte";
  export let text: string;

  let iconEl: HTMLElement;
  let tooltipEl: HTMLDivElement;
  let visible = false;
  let tipLeft = 0;
  let tipTop = 0;
  let arrowLeft = "50%";

  function portal(node: HTMLElement) {
    document.body.appendChild(node);
    return {
      destroy() {
        if (node.parentNode) node.parentNode.removeChild(node);
      },
    };
  }

  async function show() {
    visible = true;
    await tick();
    if (!iconEl || !tooltipEl) return;

    const icon = iconEl.getBoundingClientRect();
    const tw = tooltipEl.offsetWidth;
    const th = tooltipEl.offsetHeight;
    const vw = window.innerWidth;
    const padding = 10;

    const naturalLeft = icon.left + icon.width / 2 - tw / 2;
    const clampedLeft = Math.max(padding, Math.min(vw - padding - tw, naturalLeft));
    const shift = clampedLeft - naturalLeft;

    tipLeft = clampedLeft;
    tipTop = icon.top - th - 8;

    const arrowPos = tw / 2 - shift;
    arrowLeft = `${Math.max(12, Math.min(tw - 12, arrowPos))}px`;
  }

  function hide() {
    visible = false;
  }
</script>

<div class="relative inline-flex items-center" bind:this={iconEl}>
  <div
    class="w-4 h-4 flex items-center justify-center rounded-full
           border border-gray-300 dark:border-gray-700 text-gray-400 dark:text-gray-500 font-serif italic text-[11px] cursor-pointer
           hover:border-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors duration-150 select-none"
    role="button"
    tabindex="0"
    aria-label={text}
    on:mouseenter={show}
    on:mouseleave={hide}
    on:focusin={show}
    on:focusout={hide}
  >
    i
  </div>
</div>

{#if visible}
  <div
    use:portal
    bind:this={tooltipEl}
    role="tooltip"
    class="fixed w-max max-w-[85vw] md:max-w-xs
           bg-gray-900 text-gray-100 text-xs leading-relaxed px-3 py-2 rounded-xl shadow-xl shadow-black/20
           border border-gray-800 pointer-events-none z-[99999] whitespace-normal"
    style="left: {tipLeft}px; top: {tipTop}px;"
  >
    {text}
    <div
      class="absolute top-full w-2 h-2 bg-gray-900 border-b border-r border-gray-800"
      style="left: {arrowLeft}; transform: translateX(-50%) rotate(45deg);"
    ></div>
  </div>
{/if}
