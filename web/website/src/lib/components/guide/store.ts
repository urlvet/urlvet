import { browser } from '$app/environment';
import { writable } from 'svelte/store';
import type { AnalyzeResult } from '../../types';

// ── Scan state, reported by the home page so Vetty can react to it ──────────
export type ScanState =
  | { status: 'idle' }
  | { status: 'scanning' }
  | { status: 'done'; result: AnalyzeResult }
  | { status: 'error'; message: string };

export const scanState = writable<ScanState>({ status: 'idle' });

// ── Things Vetty remembers per visitor (localStorage) ──────────────────────
type Memory = {
  hidden: boolean; // user tucked him away
  introduced: boolean; // has said hello once
  nudged: boolean; // the one-time "need a hand?" bubble was shown
};

const KEY = 'vetty';
const DEFAULTS: Memory = { hidden: false, introduced: false, nudged: false };

function load(): Memory {
  if (!browser) return DEFAULTS;
  try {
    return { ...DEFAULTS, ...JSON.parse(localStorage.getItem(KEY) ?? '{}') };
  } catch {
    return DEFAULTS;
  }
}

function createMemory() {
  const store = writable<Memory>(load());
  store.subscribe((m) => {
    if (!browser) return;
    try {
      localStorage.setItem(KEY, JSON.stringify(m));
    } catch {
      /* storage blocked: Vetty just forgets */
    }
  });
  return {
    subscribe: store.subscribe,
    set: (patch: Partial<Memory>) => store.update((m) => ({ ...m, ...patch })),
  };
}

export const vettyMemory = createMemory();
