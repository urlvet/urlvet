package threatfeeds

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	PhishTankAPIURL = "https://checkurl.phishtank.com/checkurl/"
)

// phishTankRaw is the wire format from the PhishTank API.
// Field notes:
//   - phish_id        : number (int) always
//   - phish_detail_page : link to the human-readable report page
//   - verified        : bool — community confirmed phishing
//   - valid           : bool — URL is currently active / still a threat
//   - verified_at     : ISO 8601 timestamp
//   - target          : brand being impersonated (may be absent)
type phishTankRaw struct {
	Meta struct {
		Status string `json:"status"`
	} `json:"meta"`
	Results struct {
		URL             string          `json:"url"`
		InDatabase      bool            `json:"in_database"`
		PhishID         int64           `json:"phish_id"`
		PhishDetailPage string          `json:"phish_detail_page"`
		Verified        json.RawMessage `json:"verified"`
		VerifiedAt      string          `json:"verified_at"`
		Valid           json.RawMessage `json:"valid"`
		Target          string          `json:"target"`
	} `json:"results"`
}

// parseBoolField understands JSON bool and legacy string "y"/"n".
func parseBoolField(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var b bool
	if json.Unmarshal(raw, &b) == nil {
		return b
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return strings.EqualFold(s, "y") || strings.EqualFold(s, "true")
	}
	return false
}

type PhishTankResult struct {
	InDatabase      bool            `json:"in_database"`
	PhishID         int64           `json:"phish_id"`
	PhishDetailPage string          `json:"phish_detail_page"`
	Verified        bool            `json:"verified"`
	VerifiedAt      string          `json:"verified_at"`
	Valid           bool            `json:"valid"`
	Target          string          `json:"target"`
	FromCache       bool            `json:"from_cache"`
	RawResponse     json.RawMessage `json:"raw_response,omitempty"`
}

// ErrRateLimited means PhishTank refused the lookup because we've made too
// many recently. The link simply wasn't checked against PhishTank this time.
var ErrRateLimited = errors.New("phishtank: rate limited")

func CheckPhishTank(targetURL string) (*PhishTankResult, error) {
	apiKey := os.Getenv("PHISHTANK_API_KEY")
	userAgent := os.Getenv("PHISHTANK_USER_AGENT")
	if userAgent == "" {
		userAgent = "phishtank/urlvet"
	}

	data := url.Values{}
	data.Set("url", targetURL)
	data.Set("format", "json")
	if apiKey != "" {
		data.Set("app_key", apiKey)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("POST", PhishTankAPIURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("phishtank api returned status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("phishtank response read failed: %w", err)
	}

	var raw phishTankRaw
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&raw); err != nil {
		return nil, fmt.Errorf("phishtank response decode failed: %w", err)
	}

	// Reject API error envelopes (status 200 with error payload).
	if raw.Meta.Status != "" && raw.Meta.Status != "success" {
		return nil, fmt.Errorf("phishtank api error status: %s", raw.Meta.Status)
	}

	// Compact for clean inline JSON in the API response.
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, bodyBytes); err == nil {
		bodyBytes = compacted.Bytes()
	}

	result := &PhishTankResult{
		InDatabase:  raw.Results.InDatabase,
		RawResponse: json.RawMessage(bodyBytes),
	}

	// Only populate phishing details when the URL is actually in the database.
	if result.InDatabase {
		result.PhishID = raw.Results.PhishID
		result.PhishDetailPage = raw.Results.PhishDetailPage
		result.Verified = parseBoolField(raw.Results.Verified)
		result.VerifiedAt = raw.Results.VerifiedAt
		result.Valid = parseBoolField(raw.Results.Valid)
		result.Target = raw.Results.Target
	}

	return result, nil
}
