export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export interface AnalyzeResult {
  url: string;
  domain: string;
  result: {
    risk_score: number;
    trust_score: number;
    final_score: number;
    verdict: string;
    reasons?: {
      neutral_reasons?: string[];
      good_reasons?: string[];
      bad_reasons?: string[];
    };
  };
  analysis?: any;
  features?: any;
  infrastructure?: any;
  domain_info?: any;
  performance?: any;
  ssl_info?: SSLInfo;
  tls_info?: TLSInfo;
  content_data?: ContentData;
  domain_randomness?: DomainRandomness;
  typosquat_result?: TyposquatResult;
  phishing?: PhishingResult;
  incomplete?: boolean;
  /** Checks that didn't finish (task names), e.g. "phishtank_check". */
  incomplete_checks?: string[];
  errors?: any;
}

export interface SSLInfo {
  Domain: string;
  HasTLS: boolean;
  ChainValid: boolean;
  Issuer: string;
  NotBefore: string;
  NotAfter: string;
  AgeDays: number;
  Fingerprint: string;
  IsSuspicious: boolean;
  Reasons: string[];
  CTLogged: boolean;
  KnownBadChain: boolean;
}

export interface TLSInfo {
  Present: boolean;
  Issuer: string;
  AgeDays: number;
  HostnameMismatch: boolean;
}

export interface ContentData {
  url: string;
  title: string;
  has_forms: boolean;
  has_login_form: boolean;
  has_payment_form: boolean;
  has_personal_form: boolean;
  form_count: number;
  fetch_duration: number;
  forms?: Array<{
    action: string;
    method: string;
    inputs: string[];
    has_password: boolean;
    has_user_like: boolean;
    has_payment: boolean;
    has_personal: boolean;
    submit_texts: string[];
    is_external: boolean;
    is_hidden: boolean;
  }>;
  iframes?: Array<{
    src: string;
    is_hidden: boolean;
    width: string;
    height: string;
  }>;
  has_hidden_iframe: boolean;
  has_tracking: boolean;
  brand_check?: {
    brand_found: string;
    /** The brand's own site, when the page claims a brand this domain isn't. */
    official_domain?: string;
    is_mismatch: boolean;
    detected_names: string[];
  };
}

export interface TyposquatResult {
  is_suspicious: boolean;
  matched_domain?: string;
  matched_brand?: string;
  distance?: number;
  is_combo_squat?: boolean;
}

export interface DomainRandomness {
  entropy: number;
}

export interface PhishingResult {
  in_database: boolean;
  phish_id: number;
  phish_detail_page: string;
  verified: boolean;
  verified_at: string;
  valid: boolean;
  target: string;
  source: string;
  from_cache: boolean;
  raw_response?: unknown;
}

export interface ScreenshotResponse {
  status: string;
  msg: string;
  file: string;
}

export interface RankResponse {
  rank: number;
}

export interface IpCheckResponse {
  uses_ip: boolean;
}

export interface IpResolveResponse {
  ip_addresses: string[];
}

export interface LengthResponse {
  too_long: boolean;
}

export interface DepthResponse {
  too_deep: boolean;
}

export interface HstsResponse {
  has_hsts: boolean;
}

export interface RedirectsResponse {
  redirect_count: number;
}

export interface PunycodeResponse {
  uses_punycode: boolean;
}

export interface TrustedTldResponse {
  is_trusted_tld: boolean;
  is_icann: boolean;
  is_hosting_platform: boolean;
}

export interface RiskyTldResponse {
  is_risky_tld: boolean;
  is_icann: boolean;
  is_hosting_platform: boolean;
}

export interface UrlShortenerResponse {
  is_url_shortener: boolean;
}

export interface StatusCodeResponse {
  status_code: number;
}

export interface WhoisResponse {
  domain: string;
  age: number;
  created_at: string;
  expires_at: string;
  registrar: string;
  raw_data: any;
}

export interface PhishingCheckResult {
  rank?: RankResponse;
  ipCheck?: IpCheckResponse;
  ipResolve?: IpResolveResponse;
  length?: LengthResponse;
  depth?: DepthResponse;
  hsts?: HstsResponse;
  redirects?: RedirectsResponse;
  punycode?: PunycodeResponse;
  trustedTld?: TrustedTldResponse;
  riskyTld?: RiskyTldResponse;
  urlShortener?: UrlShortenerResponse;
  statusCode?: StatusCodeResponse;
  whois?: WhoisResponse;
  [key: string]: any; // Allow string indexing
}

export interface CheckStatus {
  loading: boolean;
  error?: string;
  completed: boolean;
}
