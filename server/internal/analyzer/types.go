package analyzer

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
	"github.com/urlvet/urlvet/internal/service/typosquat"
)

// Response and related public types mirror the handler's previous structs
type Response struct {
	URL              string                        `json:"url"`
	Domain           string                        `json:"domain"`
	Features         Features                      `json:"features"`
	Infrastructure   Infrastructure                `json:"infrastructure"`
	DomainInfo       *domaininfo.RegistrationData  `json:"domain_info"`
	Analysis         Analysis                      `json:"analysis"`
	SSLInfo          checks.SSLCertResult          `json:"ssl_info"`
	TLSInfo          checks.TLSResult              `json:"tls_info"`
	ContentData      *checks.PageFormResult        `json:"content_data"`
	DomainRandomness checks.DomainRandomnessResult `json:"domain_randomness"`
	TyposquatResult  typosquat.TyposquatResult     `json:"typosquat_result"`
	Phishing         *PhishingResult               `json:"phishing"`
	Performance      Performance                   `json:"performance"`
	Result           Result                        `json:"result"`
	Incomplete       bool                          `json:"incomplete"`
	Errors           []string                      `json:"errors"`
}

// PhishingResult is the unified phishing-check output exposed in the API response.
type PhishingResult struct {
	InDatabase      bool            `json:"in_database"`
	PhishID         int64           `json:"phish_id"`
	PhishDetailPage string          `json:"phish_detail_page"`
	Verified        bool            `json:"verified"`
	VerifiedAt      string          `json:"verified_at"`
	Valid           bool            `json:"valid"`
	Target          string          `json:"target"`
	Source          string          `json:"source"`
	FromCache       bool            `json:"from_cache"`
	RawResponse     json.RawMessage `json:"raw_response,omitempty"`
}

type Features struct {
	Rank int       `json:"rank"`
	TLD  TLDInfo   `json:"tld"`
	URL  URLChecks `json:"url"`
}

type TLDInfo struct {
	TLD               string `json:"tld"`
	IsTrusted         bool   `json:"is_trusted_tld"`
	IsRisky           bool   `json:"is_risky_tld"`
	IsICANN           bool   `json:"is_icann"`
	IsHostingPlatform bool   `json:"is_hosting_platform"`
}

type Keywords struct {
	HasKeywords bool                `json:"has_keywords"`
	Found       []string            `json:"found"`
	Categories  map[string][]string `json:"categories"`
}

type URLChecks struct {
	IsURLShortener   bool     `json:"url_shortener"`
	UsesIP           bool     `json:"uses_ip"`
	ContainsPunycode bool     `json:"contains_punycode"`
	TooLong          bool     `json:"too_long"`
	TooDeep          bool     `json:"too_deep"`
	HasHomoglyph     bool     `json:"has_homoglyph"`
	SubdomainCount   int      `json:"subdomain_count"`
	Keywords         Keywords `json:"keywords"`
}

type Infrastructure struct {
	IPAddresses      []string `json:"ip_addresses"`
	NameserversValid bool     `json:"nameservers_valid"`
	NSHosts          []string `json:"ns_hosts"`
	MXRecordsValid   bool     `json:"mx_records_valid"`
	MXHosts          []string `json:"mx_hosts"`
}

type HTTPStatus struct {
	Code                 int    `json:"code"`
	Text                 string `json:"text"`
	Success              bool   `json:"success"`
	IsRedirectStatusCode bool   `json:"is_redirect"`
}

type Analysis struct {
	RedirectionResult checks.RedirectionResult `json:"redirection_result"`
	// RedirectionChain       []string   `json:"redirection_chain"`
	// RedirectionChainLength int        `json:"redirection_chain_length"`
	// RedirectionFinalURL    string     `json:"redirection_final_url"`
	HTTPStatus   HTTPStatus `json:"http_status"`
	SupportsHSTS bool       `json:"is_hsts_supported"`
}

type Result struct {
	RiskScore  int     `json:"risk_score"`
	TrustScore int     `json:"trust_score"`
	FinalScore int     `json:"final_score"`
	Verdict    string  `json:"verdict"`
	Reasons    Reasons `json:"reasons"`
}

type Reasons struct {
	NeutralReasons []string `json:"neutral_reasons"`
	GoodReasons    []string `json:"good_reasons"`
	BadReasons     []string `json:"bad_reasons"`
}

type Performance struct {
	TotalTime string        `json:"total_time"`
	Timings   []TimingEntry `json:"timings"`
}

// Internal inputs/outputs for analyzer pipeline
type Input struct {
	URL    string
	Domain string
	Cache  CacheInterface // Optional cache interface
}

// CacheInterface defines the cache operations needed by tasks
type CacheInterface interface {
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// Output collects task results. Tasks write to it concurrently, under mu.
type Output struct {
	mu sync.Mutex
	OutputData
}

// OutputData is the result data itself, kept apart from the lock so a
// consistent copy can be taken (see runTasks).
type OutputData struct {
	Timings map[string]string

	// features
	Rank                 int
	TLDTrusted           bool
	TLDRisky             bool
	TLDICANN             bool
	TLD                  string
	TLDIsHostingPlatform bool

	URLIsShortener     bool
	URLUsesIP          bool
	URLContainsPuny    bool
	URLTooLong         bool
	URLTooDeep         bool
	URLSubdomainCount  int
	URLKeywordsPresent bool
	URLKeywordMatches  []string
	URLKeywordCats     map[string][]string
	DomainRandomness   checks.DomainRandomnessResult
	HomoglyphPresent   bool

	// infra
	IPs     []string
	NSValid bool
	MXValid bool
	MXHosts []string
	NSHosts []string

	// analysis
	RedirectionResult checks.RedirectionResult
	SupportsHSTS      bool
	StatusCode        int
	StatusText        string
	StatusSuccess     bool
	StatusIsRedirect  bool

	DomainInfo      *domaininfo.RegistrationData
	SSLInfo         checks.SSLCertResult
	ContentData     *checks.PageFormResult
	TLSInfo         checks.TLSResult
	PhishTank       *threatfeeds.PhishTankResult
	TyposquatResult typosquat.TyposquatResult
}

func (o *Output) setTiming(name string, d time.Duration) {
	o.mu.Lock()
	if o.Timings == nil {
		o.Timings = map[string]string{}
	}
	o.Timings[name] = d.String()
	o.mu.Unlock()
}

// Task represents a single analyzer unit of work
type Task interface {
	Name() string
	Run(in *Input, out *Output) error
}
