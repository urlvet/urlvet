package checks

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type RedirectionResult struct {
	IsRedirected  bool     `json:"is_redirected"`
	ChainLength   int      `json:"chain_length"`
	Chain         []string `json:"chain"`
	FinalURL      string   `json:"final_url"`
	FinalURLHost  string   `json:"final_url_domain"`
	HasDomainJump bool     `json:"has_domain_jump"`
}

func CheckRedirects(rawURL string) (RedirectionResult, error) {
	var redirects []string

	// Keep cookies across hops, like a browser (see CheckHTTPCombined).
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: newSafeTransport(),
		Jar:       jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects = append(redirects, req.URL.String())
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return RedirectionResult{}, err
	}
	defer resp.Body.Close()

	chain := append([]string{rawURL}, redirects...)
	finalURL := chain[len(chain)-1]
	finalURLHost, _ := GetHost(finalURL)

	// Detect domain jumps

	origDomain, _ := GetDomain(rawURL)
	hasJump := false

	for _, u := range chain[1:] {
		urlDomain, _ := GetDomain(u)
		if !sameSite(urlDomain, origDomain) {
			hasJump = true
			break
		}
	}

	return RedirectionResult{
		IsRedirected:  len(redirects) > 0,
		Chain:         chain,
		FinalURL:      finalURL,
		FinalURLHost:  finalURLHost,
		ChainLength:   len(chain),
		HasDomainJump: hasJump,
	}, nil
}
