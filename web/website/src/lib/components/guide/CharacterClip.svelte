<script lang="ts">
  import "./character.css";
  import Face from "./Face.svelte";
  import type { Mood } from "./mood";

  // Vetty: a paperclip helper in the spirit of the classic office assistant.
  // Silver wire, big eyes on the wire, bushy brows, leaning in to help.
  // Drawn from scratch; not Microsoft's artwork.
  export let mood: Mood = "idle";
  export let size = 96;
  /** Pupil offset to look at something (the pointer, the search bar); null = idle glances. */
  export let gaze: [number, number] | null = null;

  // Unique gradient ids so several instances can share a page.
  const uid = `clip-${Math.random().toString(36).slice(2, 8)}`;

  // One continuous wire: inner loop, bottom loop, outer loop.
  const WIRE =
    "M47 50 L47 22 A8 8 0 0 0 31 22 L31 61 A10.5 10.5 0 0 0 52 61 L52 16 A14 14 0 0 0 24 16 L24 55";
</script>

<svg class="guide-char clip" width={size} height={size} viewBox="0 0 80 78" aria-hidden="true">
  <defs>
    <linearGradient id="{uid}-wire" x1="0" y1="0" x2="1" y2="0.25">
      <stop offset="0" style="stop-color: var(--wire-lo)" />
      <stop offset="0.45" style="stop-color: var(--wire-hi)" />
      <stop offset="1" style="stop-color: var(--wire-lo)" />
    </linearGradient>
  </defs>

  <!-- a rare stretch while idle, so he never looks frozen -->
  <g class:stretch={mood === "idle"}>
    {#key mood}
      <g class="hop">
        <g class="body" class:shiver={mood === "alarmed"}>
          <g transform="rotate(-7 40 50)">
            <!-- wire: dark edge, metallic body, thin shine -->
            <path d={WIRE} class="wire-edge" stroke-width="5.6" />
            <path d={WIRE} stroke="url(#{uid}-wire)" class="wire" stroke-width="3.8" />
            <path d={WIRE} class="wire-shine" stroke-width="1" />

            <!-- left arm -->
            {#if mood === "point"}
              <!-- raised, pointing up-left (Vetty sits bottom-right, the tour targets are above him) -->
              <path d="M24 44 Q16 38 9 29" class="ink" stroke-width="2.4" />
              <circle cx="8.2" cy="27.8" r="1.8" style="fill: var(--ink)" />
            {:else if mood === "alarmed"}
              <!-- holding up a little shield -->
              <path d="M24 47 Q18 49 15 52" class="ink" stroke-width="2.4" />
              <path
                d="M5 47 L17 47 L17 54 Q17 61 11 64 Q5 61 5 54 Z"
                class="ink prop-fill"
                stroke-width="2"
              />
              <path d="M11 51 L11 56" class="ink" stroke-width="2" />
              <circle cx="11" cy="59" r="0.9" style="fill: var(--ink)" />
            {:else if mood === "ooh"}
              <path d="M24 46 Q17 42 14 36" class="ink" stroke-width="2.4" />
            {:else}
              <path d="M24 48 Q16 50 13 57" class="ink" stroke-width="2.4" />
            {/if}

            <!-- right arm -->
            {#if mood === "wave"}
              <path d="M52 44 Q60 38 62 27" class="ink wave-arm" stroke-width="2.4" />
            {:else if mood === "thinking" || mood === "suspicious"}
              <path d="M52 46 Q58 44 54 37 Q51 34 46 37" class="ink" stroke-width="2.4" />
            {:else if mood === "happy"}
              <!-- thumbs up -->
              <g class="thumb">
                <path d="M52 46 Q59 43 61 36" class="ink" stroke-width="2.4" />
                <rect
                  x="58.6"
                  y="31.6"
                  width="5.6"
                  height="5"
                  rx="1.8"
                  class="ink prop-fill"
                  stroke-width="1.8"
                />
                <path d="M60.2 31.6 L60.2 27.6" class="ink" stroke-width="2.2" />
              </g>
            {:else if mood === "scanning"}
              <!-- sweeping a magnifying glass over the link -->
              <g class="sweep">
                <path d="M52 48 Q58 51 60 55" class="ink" stroke-width="2.4" />
                <path d="M60 55 L63.5 51.5" class="ink" stroke-width="2.6" />
                <circle cx="68.5" cy="46.5" r="6.8" class="ink lens" stroke-width="2.2" />
              </g>
            {:else if mood === "ooh"}
              <path d="M52 46 Q59 42 62 36" class="ink" stroke-width="2.4" />
            {:else}
              <path d="M52 48 Q60 50 63 57" class="ink" stroke-width="2.4" />
            {/if}

            <g class="face">
              <Face
                lx={31}
                rx={47}
                y={23}
                mouthY={38}
                {mood}
                brows
                knockout
                eyeRx={7}
                eyeRy={8.2}
                pupilR={3}
                browWidth={3.4}
                browSpan={5.5}
                {gaze}
              />
            </g>
          </g>
        </g>
      </g>
    {/key}
  </g>
</svg>
