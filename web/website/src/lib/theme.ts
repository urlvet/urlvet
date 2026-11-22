import { browser } from '$app/env';
import { get, writable } from 'svelte/store';

type Theme = 'light' | 'dark';

const store = writable<Theme>('light');

export const theme = {
  subscribe: store.subscribe,

  init(): void {
    if (!browser) return;
    const saved = localStorage.getItem('theme') as Theme | null;
    const t = saved ?? 'light';
    store.set(t);
    document.documentElement.classList.toggle('dark', t === 'dark');
  },

  toggle(): void {
    if (!browser) return;
    const next: Theme = get(store) === 'dark' ? 'light' : 'dark';
    store.set(next);
    localStorage.setItem('theme', next);
    document.documentElement.classList.toggle('dark', next === 'dark');
  },
};
