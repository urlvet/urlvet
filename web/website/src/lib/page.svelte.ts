import { page as pageState } from '$app/state';
import type { Readable } from 'svelte/store';

// `$app/stores` is gone in SvelteKit 3. These components still use Svelte 4-style
// `$:` statements, which can't track the new `page` state directly, so this wraps it
// as a store with the same shape they already use. The page is read inside `subscribe`
// (not at import time) because on the server it can only be read while rendering.
const snapshot = () => ({ url: pageState.url, status: pageState.status, error: pageState.error });

export const page: Readable<ReturnType<typeof snapshot>> = {
  subscribe(run) {
    run(snapshot());
    return $effect.root(() => {
      $effect(() => run(snapshot()));
    });
  },
};
