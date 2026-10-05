package threatfeeds

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Google Web Risk (https://cloud.google.com/web-risk) is Safe Browsing's
// commercial-use API: 100,000 lookups a month free, then billed. It runs only
// when WEBRISK_API_KEY is set. (Safe Browsing itself is free but
// non-commercial only.)

const webRiskURL = "https://webrisk.googleapis.com/v1/uris:search"

// ErrWebRiskDisabled means no API key is configured.
var ErrWebRiskDisabled = errors.New("web risk: no API key")

// ErrWebRiskRateLimited means the monthly quota or rate limit was hit.
var ErrWebRiskRateLimited = errors.New("web risk: rate limited")

// GoogleThreatResult is a Google list's verdict on one URL (Web Risk or Safe Browsing).
type GoogleThreatResult struct {
	Listed      bool       `json:"listed"`
	ThreatTypes []string   `json:"threat_types,omitempty"` // SOCIAL_ENGINEERING, MALWARE, UNWANTED_SOFTWARE
	ExpireTime  *time.Time `json:"expire_time,omitempty"`  // how long Google says the answer holds (Web Risk only)
}

// WebRiskEnabled reports whether lookups are configured.
func WebRiskEnabled() bool { return os.Getenv("WEBRISK_API_KEY") != "" && !LookupsOff() }

// CheckWebRisk looks a URL up in Google's phishing, malware and unwanted
// software lists.
func CheckWebRisk(rawURL string) (*GoogleThreatResult, error) {
	key := os.Getenv("WEBRISK_API_KEY")
	if key == "" {
		return nil, ErrWebRiskDisabled
	}
	q := url.Values{}
	q.Set("uri", rawURL)
	q.Add("threatTypes", "SOCIAL_ENGINEERING")
	q.Add("threatTypes", "MALWARE")
	q.Add("threatTypes", "UNWANTED_SOFTWARE")
	q.Set("key", key)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(webRiskURL + "?" + q.Encode())
	if err != nil {
		// The error would quote the request URL, API key included, and scan
		// errors are shown to users. Keep only the cause.
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return nil, fmt.Errorf("web risk: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrWebRiskRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("web risk: HTTP %d", resp.StatusCode)
	}

	var body struct {
		Threat *struct {
			ThreatTypes []string  `json:"threatTypes"`
			ExpireTime  time.Time `json:"expireTime"`
		} `json:"threat"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("web risk: %w", err)
	}
	if body.Threat == nil {
		return &GoogleThreatResult{}, nil
	}
	expire := body.Threat.ExpireTime
	return &GoogleThreatResult{Listed: true, ThreatTypes: body.Threat.ThreatTypes, ExpireTime: &expire}, nil
}
