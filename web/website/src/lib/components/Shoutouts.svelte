<script lang="ts">
  import { onMount } from "svelte";
  import { quintOut } from "svelte/easing";
  import { fly } from "svelte/transition";

  import { SHOUTOUTS as items } from "../data/shoutouts";
  import ShoutoutCard from "./shoutouts/ShoutoutCard.svelte";

  const loop = [...items, ...items];

  // --- Ticker (desktop) ---
  let track: HTMLElement;
  let offset = 0;
  let halfWidth = 0;
  let isDragging = false;
  let isHovered = false;
  let dragStartX = 0;
  let dragStartOffset = 0;

  // --- Slideshow (mobile) ---
  let isMobile = false;
  let currentIndex = 0;
  let direction = 1; // 1 = forward (right→left), -1 = backward (left→right)
  let touchStartX = 0;
  let touchStartOffset = 0;

  function goTo(index: number) {
    direction = index > currentIndex ? 1 : -1;
    currentIndex = index;
  }

  onMount(() => {
    isMobile = window.innerWidth < 768;

    if (!isMobile) {
      halfWidth = track.scrollWidth / 2;
    }

    let slideTimer: ReturnType<typeof setInterval>;

    function startSlideshow() {
      slideTimer = setInterval(() => {
        direction = 1;
        currentIndex = (currentIndex + 1) % items.length;
      }, 4000);
    }

    if (isMobile) startSlideshow();

    let frame: number;
    function tick() {
      if (!isMobile && !isDragging && !isHovered) {
        offset += 0.5;
        if (offset >= halfWidth) offset -= halfWidth;
        track.style.transform = `translateX(-${offset}px)`;
      }
      frame = requestAnimationFrame(tick);
    }
    frame = requestAnimationFrame(tick);

    const onResize = () => {
      const nowMobile = window.innerWidth < 768;
      if (nowMobile === isMobile) return;
      isMobile = nowMobile;
      if (isMobile) {
        startSlideshow();
      } else {
        clearInterval(slideTimer);
        setTimeout(() => {
          if (track) halfWidth = track.scrollWidth / 2;
        }, 50);
      }
    };
    window.addEventListener("resize", onResize);

    return () => {
      cancelAnimationFrame(frame);
      clearInterval(slideTimer);
      window.removeEventListener("resize", onResize);
    };
  });

  // Ticker mouse handlers
  function onMouseEnter() {
    isHovered = true;
  }
  function onMouseLeave() {
    isHovered = false;
    isDragging = false;
  }
  function onMouseDown(e: MouseEvent) {
    isDragging = true;
    dragStartX = e.clientX;
    dragStartOffset = offset;
  }
  function onMouseUp() {
    isDragging = false;
  }
  function onMouseMove(e: MouseEvent) {
    if (!isDragging) return;
    e.preventDefault();
    offset = dragStartOffset + (dragStartX - e.clientX);
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }

  // Touch handlers (swipe on mobile, drag on desktop)
  function onTouchStart(e: TouchEvent) {
    touchStartX = e.touches[0].clientX;
    if (!isMobile) {
      touchStartOffset = offset;
      isDragging = true;
    }
  }
  function onTouchMove(e: TouchEvent) {
    if (isMobile || !isDragging) return;
    offset = touchStartOffset + (touchStartX - e.touches[0].clientX);
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }
  function onTouchEnd(e: TouchEvent) {
    isDragging = false;
    if (!isMobile) return;
    const diff = touchStartX - (e.changedTouches[0]?.clientX ?? touchStartX);
    if (Math.abs(diff) > 40) {
      if (diff > 0) {
        direction = 1;
        currentIndex = (currentIndex + 1) % items.length;
      } else {
        direction = -1;
        currentIndex = (currentIndex - 1 + items.length) % items.length;
      }
    }
  }
</script>

<div class="pb-8 md:pb-0">
  <p class="font-mono text-xs text-gray-500 uppercase tracking-wider mb-6 text-center">
    What people are saying
  </p>

  {#if isMobile}
    <!-- Slideshow (mobile) -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="relative px-4"
      on:touchstart={onTouchStart}
      on:touchend={onTouchEnd}
      role="region"
      aria-label="Testimonials"
    >
      <div class="grid overflow-hidden">
        {#key currentIndex}
          <div
            class="col-start-1 row-start-1"
            in:fly={{ x: direction * 300, duration: 700, easing: quintOut }}
            out:fly={{ x: direction * -300, duration: 700, easing: quintOut }}
          >
            <ShoutoutCard item={items[currentIndex]} variant="slide" />
          </div>
        {/key}
      </div>

      <!-- Dot indicators -->
      <div class="flex justify-center gap-2 mt-4">
        {#each items as _, i}
          <!-- svelte-ignore a11y_consider_explicit_label -->
          <button
            class="w-1.5 h-1.5 rounded-full transition-colors {i === currentIndex
              ? 'bg-gray-500 dark:bg-gray-400'
              : 'bg-gray-200 dark:bg-gray-700'}"
            on:click={() => goTo(i)}
            aria-label="Go to slide {i + 1}"
          ></button>
        {/each}
      </div>
    </div>
  {:else}
    <!-- Ticker (desktop) -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="ticker-fade overflow-x-hidden overflow-y-visible relative py-2"
      on:mouseenter={onMouseEnter}
      on:mouseleave={onMouseLeave}
      on:mousedown={onMouseDown}
      on:mouseup={onMouseUp}
      on:mousemove={onMouseMove}
      on:touchstart={onTouchStart}
      on:touchmove={onTouchMove}
      on:touchend={onTouchEnd}
      role="region"
      aria-label="Testimonials"
    >
      <div
        bind:this={track}
        class="flex gap-4 will-change-transform {isDragging ? 'cursor-grabbing' : 'cursor-grab'}"
        style="width: max-content;"
      >
        {#each loop as item}
          <ShoutoutCard {item} />
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  /* Fade cards out at both edges instead of hard-clipping them. */
  .ticker-fade {
    -webkit-mask-image: linear-gradient(to right, transparent, #000 15%, #000 85%, transparent);
    mask-image: linear-gradient(to right, transparent, #000 15%, #000 85%, transparent);
  }
</style>
