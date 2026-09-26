// Shapes returned by the /api/v1/admin endpoints.

export interface CacheEntry {
  key: string;
  prefix: string;
  ttl_seconds: number;
  value: unknown;
}

export interface ScanRecord {
  url: string;
  domain: string;
  verdict: string;
  score: number;
  duration: string;
  time: string;
  cached: boolean;
}

export interface ErrorRecord {
  task: string;
  error: string;
  url: string;
  time: string;
}

export interface ReportRecord {
  /** Derived from the stored report; used to delete it. */
  id: string;
  url: string;
  verdict: string;
  score: number;
  expected_verdict?: string;
  comment?: string;
  time: string;
}

export interface DomainCount {
  domain: string;
  count: number;
}

export interface Stats {
  total_scans_today: number;
  total_scans_all: number;
  cache_hits: number;
  cache_misses: number;
  cache_hit_rate: number;
  avg_duration_ms: number;
  top_domains: DomainCount[];
  verdict_counts: Record<string, number>;
}

export type Tab = 'overview' | 'recent' | 'errors' | 'reports' | 'cache';

/** Runs an admin request, handling auth expiry and errors centrally. Undefined on failure. */
export type Guard = <T>(request: Promise<T>) => Promise<T | undefined>;
