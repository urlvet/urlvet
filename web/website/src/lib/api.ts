import { env } from '$env/dynamic/public';
import type {
  ApiResponse,
  DepthResponse,
  HstsResponse,
  IpCheckResponse,
  IpResolveResponse,
  LengthResponse,
  PunycodeResponse,
  RankResponse,
  RedirectsResponse,
  RiskyTldResponse,
  StatusCodeResponse,
  TrustedTldResponse,
  UrlShortenerResponse,
  WhoisResponse,
} from './types';

const PUBLIC_BASE_URL = env.PUBLIC_BASE_URL || 'http://localhost:8080/api/v1';

async function makeRequest<T>(endpoint: string, url: string): Promise<ApiResponse<T>> {
  try {
    // Ensure the URL is properly formatted before encoding
    const formattedUrl = url.trim();
    const encodedUrl = encodeURIComponent(formattedUrl);
    const fullUrl = `${PUBLIC_BASE_URL}${endpoint}?url=${encodedUrl}`;

    const response = await fetch(fullUrl);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      console.error(`API Error:`, errorData);
      return { error: errorData.error || `HTTP ${response.status}` };
    }

    const data = await response.json();
    return { data };
  } catch (error) {
    console.error(`API Error for ${endpoint}:`, error);
    return { error: error instanceof Error ? error.message : 'Network error' };
  }
}

async function makeScreenshotRequest(url: string): Promise<ApiResponse<string | null>> {
  try {
    const formattedUrl = url.trim();
    const encodedUrl = encodeURIComponent(formattedUrl);
    const fullUrl = `${PUBLIC_BASE_URL}/screenshot?url=${encodedUrl}`;

    const response = await fetch(fullUrl);

    if (!response.ok) {
      console.error(`Screenshot API Error: HTTP ${response.status}`);
      return { error: `HTTP ${response.status}` };
    }

    // Check if response is an image
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.startsWith('image/')) {
      // Convert blob to object URL
      const blob = await response.blob();
      const blobUrl = URL.createObjectURL(blob);
      return { data: blobUrl };
    }

    // If not an image, return null
    return { data: null };
  } catch (error) {
    console.error(`Screenshot API Error:`, error);
    return { error: error instanceof Error ? error.message : 'Network error' };
  }
}

export interface ReportPayload {
  url: string;
  verdict?: string;
  score?: number;
  expected_verdict?: 'Safe' | 'Suspicious' | 'Risky';
  comment?: string;
}

async function submitReport(payload: ReportPayload): Promise<ApiResponse<{ status: string }>> {
  try {
    const response = await fetch(`${PUBLIC_BASE_URL}/report`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
      return { error: data.error || `HTTP ${response.status}` };
    }
    return { data };
  } catch (error) {
    return { error: error instanceof Error ? error.message : 'Network error' };
  }
}

export const api = {
  async report(payload: ReportPayload): Promise<ApiResponse<{ status: string }>> {
    return submitReport(payload);
  },
  async analyze(url: string): Promise<ApiResponse<any>> {
    return makeRequest<any>('/analyze', url);
  },
  async screenshot(url: string): Promise<ApiResponse<string | null>> {
    return makeScreenshotRequest(url);
  },
  async getRank(url: string): Promise<ApiResponse<RankResponse>> {
    return makeRequest<RankResponse>('/rank', url);
  },

  async getIpCheck(url: string): Promise<ApiResponse<IpCheckResponse>> {
    return makeRequest<IpCheckResponse>('/ip/check', url);
  },

  async getIpResolve(url: string): Promise<ApiResponse<IpResolveResponse>> {
    return makeRequest<IpResolveResponse>('/ip/resolve', url);
  },

  async getLength(url: string): Promise<ApiResponse<LengthResponse>> {
    return makeRequest<LengthResponse>('/length', url);
  },

  async getDepth(url: string): Promise<ApiResponse<DepthResponse>> {
    return makeRequest<DepthResponse>('/depth', url);
  },

  async getHsts(url: string): Promise<ApiResponse<HstsResponse>> {
    return makeRequest<HstsResponse>('/hsts', url);
  },

  async getRedirects(url: string): Promise<ApiResponse<RedirectsResponse>> {
    return makeRequest<RedirectsResponse>('/redirects', url);
  },

  async getPunycode(url: string): Promise<ApiResponse<PunycodeResponse>> {
    return makeRequest<PunycodeResponse>('/punycode', url);
  },

  async getTrustedTld(url: string): Promise<ApiResponse<TrustedTldResponse>> {
    return makeRequest<TrustedTldResponse>('/trusted-tld', url);
  },

  async getRiskyTld(url: string): Promise<ApiResponse<RiskyTldResponse>> {
    return makeRequest<RiskyTldResponse>('/risky-tld', url);
  },

  async getUrlShortener(url: string): Promise<ApiResponse<UrlShortenerResponse>> {
    return makeRequest<UrlShortenerResponse>('/url-shortener', url);
  },

  async getStatusCode(url: string): Promise<ApiResponse<StatusCodeResponse>> {
    return makeRequest<StatusCodeResponse>('/status-code', url);
  },

  async getWhois(url: string): Promise<ApiResponse<WhoisResponse>> {
    return makeRequest<WhoisResponse>('/whois', url);
  },
};
