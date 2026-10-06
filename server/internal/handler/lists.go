package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/constants"
	"github.com/urlvet/urlvet/internal/service/rank"
	"github.com/urlvet/urlvet/internal/service/threatfeeds"
)

// Lists for clients (the browser extension) to keep locally: known-bad links
// to block before a page loads, and well-known sites that needn't be sent for
// checking at all.

// knownSitesCount is how many of the most popular sites clients may skip.
const knownSitesCount = 10000

// ThreatListResponse is the hashed list of known phishing and malware links.
type ThreatListResponse struct {
	Version   string   `json:"version"`
	Hash      string   `json:"hash" example:"sha256"`
	PrefixLen int      `json:"prefix_len" example:"8"`
	Count     int      `json:"count"`
	Sources   []string `json:"sources"`
	Prefixes  string   `json:"prefixes"` // base64 of the sorted prefixes, back to back
}

// KnownSitesResponse lists well-known sites, and the hosts on them where
// anyone can publish, which must still be checked.
type KnownSitesResponse struct {
	Version     string   `json:"version"`
	Sites       []string `json:"sites"`
	UserContent []string `json:"user_content"`
}

// notModified answers a conditional request for an unchanged list.
func notModified(c *gin.Context, version string) bool {
	etag := `"` + version + `"`
	c.Header("ETag", etag)
	c.Header("Cache-Control", "public, max-age=3600")
	if strings.Contains(c.GetHeader("If-None-Match"), etag) {
		c.Status(http.StatusNotModified)
		return true
	}
	return false
}

// ThreatListHandler returns known phishing and malware links as SHA-256
// prefixes of their normalized form (see threatfeeds.FeedKey).
//
//	@Summary		Hashed threat list
//	@Description	Known phishing and malware links (PhishTank, URLhaus) as sorted SHA-256 prefixes of each link's normalized form: no scheme, lowercase host without "www.", no trailing slash or fragment, query kept. Clients hash a link the same way and look the prefix up locally, so nothing is sent. Supports If-None-Match.
//	@Tags			Lists
//	@Produce		json
//	@Success		200	{object}	ThreatListResponse
//	@Success		304
//	@Router			/lists/threats [get]
func ThreatListHandler(c *gin.Context) {
	exp := threatfeeds.ExportPrefixes()
	if notModified(c, exp.Version) {
		return
	}
	c.JSON(http.StatusOK, ThreatListResponse{
		Version:   exp.Version,
		Hash:      "sha256",
		PrefixLen: threatfeeds.PrefixLen,
		Count:     exp.Count,
		Sources:   exp.Sources,
		Prefixes:  base64.StdEncoding.EncodeToString(exp.Prefixes),
	})
}

// KnownSitesHandler returns the most popular sites, which clients needn't
// send for automatic checks, and the hosts where anyone can publish, which
// they must check anyway (sites.google.com, storage.googleapis.com, ...).
//
//	@Summary		Well-known sites
//	@Description	The 10,000 most popular sites, and the user-content hosts on them that still need checking. A host matches a site if it is the site or a subdomain of it.
//	@Tags			Lists
//	@Produce		json
//	@Success		200	{object}	KnownSitesResponse
//	@Success		304
//	@Router			/lists/known-sites [get]
func KnownSitesHandler(c *gin.Context) {
	sites := rank.TopDomains(knownSitesCount)

	userContent := map[string]struct{}{}
	for h := range constants.UserUploadHosts {
		userContent[h] = struct{}{}
	}
	for h := range constants.ProviderServiceHosts {
		userContent[h] = struct{}{}
	}
	for h := range constants.TrustedHostingPlatforms {
		userContent[h] = struct{}{}
	}
	for h := range constants.CustomerSiteHosts {
		userContent[h] = struct{}{}
	}
	for h := range constants.UserPageHosts {
		userContent[h] = struct{}{}
	}
	for h := range constants.ShortLinkPaths {
		userContent[h] = struct{}{}
	}
	// Short links lead anywhere. Only the ones that would otherwise count as
	// well known need listing (bit.ly); the rest are checked anyway.
	top := make(map[string]struct{}, len(sites))
	for _, s := range sites {
		top[s] = struct{}{}
	}
	for h := range constants.URLShorteners {
		if _, ok := top[h]; ok {
			userContent[h] = struct{}{}
		}
	}
	hosts := make([]string, 0, len(userContent))
	for h := range userContent {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	sum := sha256.New()
	for _, s := range append(append([]string{}, sites...), hosts...) {
		sum.Write([]byte(s + "\n"))
	}
	version := hex.EncodeToString(sum.Sum(nil)[:8])
	if notModified(c, version) {
		return
	}
	c.JSON(http.StatusOK, KnownSitesResponse{Version: version, Sites: sites, UserContent: hosts})
}
