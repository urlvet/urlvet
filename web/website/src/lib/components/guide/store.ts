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
  openedAt: number; // when he was last opened (ms); 0 = never
  nudgeDeclinedAt: number; // when the visitor said "No thanks" to his nudge (ms); 0 = never
};

const KEY = 'vetty';
const DEFAULTS: Memory = { hidden: false, openedAt: 0, nudgeDeclinedAt: 0 };

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

/** Vetty re-introduces himself, and nudges again, a day after he was last opened. */
const REINTRODUCE_AFTER_MS = 24 * 60 * 60 * 1000;

export function seenRecently(m: { openedAt: number }): boolean {
  return m.openedAt > 0 && Date.now() - m.openedAt < REINTRODUCE_AFTER_MS;
}

/** After a "No thanks", Vetty doesn't nudge again for a month. */
const NUDGE_DECLINED_FOR_MS = 30 * 24 * 60 * 60 * 1000;

export function nudgeDeclined(m: { nudgeDeclinedAt?: number }): boolean {
  const at = m.nudgeDeclinedAt ?? 0;
  return at > 0 && Date.now() - at < NUDGE_DECLINED_FOR_MS;
}
