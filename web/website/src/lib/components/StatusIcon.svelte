<script lang="ts">
  // Small inline line icon used in place of ✅ / ❌ / ⚠️ in section tables.
  export let kind: "ok" | "bad" | "warn" = "ok";

  const PATHS = {
    ok: "M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
    bad: "M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
    warn: "M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z",
  };
</script>

<svg
  class="w-4 h-4 shrink-0"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  viewBox="0 0 24 24"
  aria-hidden="true"
>
  <path stroke-linecap="round" stroke-linejoin="round" d={PATHS[kind]} />
</svg>
