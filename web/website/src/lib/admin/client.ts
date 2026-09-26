import { env } from '$env/dynamic/public';
import type { CacheEntry, ErrorRecord, ReportRecord, ScanRecord, Stats } from './types';

const BASE = env.PUBLIC_BASE_URL || 'http://localhost:8080/api/v1';

/** Thrown when the admin token is missing, expired or rejected. */
export class UnauthorizedError extends Error {}

export type AdminClient = ReturnType<typeof createAdminClient>;

export function createAdminClient(token: string) {
  async function request(path: string, init: RequestInit = {}): Promise<Response> {
    const res = await fetch(`${BASE}${path}`, {
      ...init,
      headers: { ...init.headers, Authorization: `Bearer ${token}` },
    });
    if (res.status === 401) throw new UnauthorizedError();
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res;
  }

  const getJSON = async <T>(path: string): Promise<T> => (await request(path)).json();

  return {
    stats: () => getJSON<Stats>('/admin/stats'),
    recent: async () => (await getJSON<{ scans?: ScanRecord[] }>('/admin/recent')).scans ?? [],
    errors: async () => (await getJSON<{ errors?: ErrorRecord[] }>('/admin/errors')).errors ?? [],
    reports: async () =>
      (await getJSON<{ reports?: ReportRecord[] }>('/admin/reports')).reports ?? [],
    cache: async () => (await getJSON<{ keys?: CacheEntry[] }>('/admin/cache')).keys ?? [],
    deleteKey: async (key: string) => {
      await request(`/admin/cache/${encodeURIComponent(key)}`, { method: 'DELETE' });
    },
    deleteReport: async (id: string) => {
      await request(`/admin/reports/${encodeURIComponent(id)}`, { method: 'DELETE' });
    },
    flushCache: async () => {
      await request('/admin/cache', { method: 'DELETE' });
    },
  };
}

/** Exchanges the admin password for a session token. */
export async function login(password: string): Promise<{ token?: string; error?: string }> {
  try {
    const res = await fetch(`${BASE}/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    });
    if (res.status === 401) return { error: 'Incorrect password.' };
    if (res.status === 503) return { error: 'Admin access is disabled on this server.' };
    if (!res.ok) return { error: 'Login failed. Please try again.' };
    const data = await res.json();
    return { token: data.token };
  } catch {
    return { error: 'Unable to reach the server.' };
  }
}
