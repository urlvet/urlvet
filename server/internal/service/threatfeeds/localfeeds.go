package threatfeeds

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/urlvet/urlvet/internal/logger"
)

// Local threat feeds: lists of known phishing and malware URLs downloaded on
// a schedule and matched in memory, so every scan is checked without a
// per-scan API call (PhishTank's API rate-limits almost every scan).
//
// Licences differ, so only feeds that allow commercial use are on by default:
//
//   - PhishTank dump: free, commercial use allowed. Without an app key only a
//     few downloads a day are allowed, so it's refreshed every 6 hours.
//   - OpenPhish community feed: non-commercial only. FEED_OPENPHISH=1.
//   - URLhaus (abuse.ch): malware URLs, free for not-for-profit use with an
//     Auth-Key (no card); commercial use needs a Spamhaus subscription.
//     URLHAUS_AUTH_KEY=….

// feed is one downloadable list.
type feed struct {
	name    string
	url     func() string
	every   time.Duration
	header  map[string]string
	parse   func(io.Reader) ([]string, error)
	enabled func() bool
}

func feeds() []feed {
	return []feed{
		{
			name: "PhishTank",
			url: func() string {
				if key := os.Getenv("PHISHTANK_API_KEY"); key != "" {
					return "http://data.phishtank.com/data/" + key + "/online-valid.json.bz2"
				}
				return "http://data.phishtank.com/data/online-valid.json.bz2"
			},
			every:   phishTankEvery(),
			parse:   parsePhishTankDump,
			enabled: func() bool { return os.Getenv("FEED_PHISHTANK") != "0" },
		},
		{
			name:    "OpenPhish",
			url:     func() string { return "https://openphish.com/feed.txt" },
			every:   12 * time.Hour,
			parse:   parseURLLines,
			enabled: func() bool { return os.Getenv("FEED_OPENPHISH") == "1" },
		},
		{
			// abuse.ch wants the Auth-Key in the export's path.
			name: "URLhaus",
			url: func() string {
				return "https://urlhaus-api.abuse.ch/v2/files/exports/" + url.PathEscape(os.Getenv("URLHAUS_AUTH_KEY")) + "/recent.csv"
			},
			every:   time.Hour,
			parse:   parseURLhausCSV,
			enabled: func() bool { return os.Getenv("URLHAUS_AUTH_KEY") != "" },
		},
	}
}

// phishTankEvery is hourly with an app key (what PhishTank allows) and every
// 6 hours without one, inside the "few downloads a day" keyless limit.
func phishTankEvery() time.Duration {
	if os.Getenv("PHISHTANK_API_KEY") != "" {
		return time.Hour
	}
	return 6 * time.Hour
}

// FeedMatch is a scan's result against the local feeds.
type FeedMatch struct {
	Listed  bool     `json:"listed"`
	Match   string   `json:"match,omitempty"`   // "url": this exact link; "host": another page on the same site
	Sources []string `json:"sources,omitempty"` // feeds that list it
	Checked []string `json:"checked,omitempty"` // feeds that were loaded and checked
}

// store holds every loaded feed, keyed by feed name, so a refresh replaces
// one feed's entries without touching the others.
type store struct {
	mu    sync.RWMutex
	urls  map[string]map[string]struct{} // feed → normalized URLs
	hosts map[string]map[string]struct{} // feed → hosts
	gen   int                            // bumped on every change, for ExportPrefixes' cache
}

var local = &store{urls: map[string]map[string]struct{}{}, hosts: map[string]map[string]struct{}{}}

func (s *store) set(name string, list []string) {
	urls := make(map[string]struct{}, len(list))
	hosts := make(map[string]struct{}, len(list))
	for _, raw := range list {
		key, host := feedKey(raw)
		if key == "" {
			continue
		}
		urls[key] = struct{}{}
		hosts[host] = struct{}{}
	}
	s.mu.Lock()
	s.urls[name] = urls
	s.hosts[name] = hosts
	s.gen++
	s.mu.Unlock()
}

// feedKey normalizes a URL for matching: no scheme, lowercase host without
// "www.", no trailing slash or fragment. Returns the key and the host.
func feedKey(raw string) (key, host string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ""
	}
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return "", ""
	}
	host = strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	key = host + strings.TrimSuffix(u.EscapedPath(), "/")
	if u.RawQuery != "" {
		key += "?" + u.RawQuery
	}
	return key, host
}

// LookupLocal checks a URL against the loaded feeds. A host match only counts
// when sharedHost is false: on a host many people publish under (github.com,
// a CDN) one reported page says nothing about the rest.
func LookupLocal(rawURL string, sharedHost bool) FeedMatch {
	key, host := feedKey(rawURL)
	if key == "" {
		return FeedMatch{}
	}
	local.mu.RLock()
	defer local.mu.RUnlock()

	var m FeedMatch
	for name, urls := range local.urls {
		if len(urls) > 0 {
			m.Checked = append(m.Checked, name)
		}
	}
	sort.Strings(m.Checked)
	for name, urls := range local.urls {
		if _, ok := urls[key]; ok {
			m.Sources = append(m.Sources, name)
		}
	}
	if len(m.Sources) > 0 {
		m.Listed, m.Match = true, "url"
		return m
	}
	if sharedHost {
		return m
	}
	for name, hosts := range local.hosts {
		if _, ok := hosts[host]; ok {
			m.Sources = append(m.Sources, name)
		}
	}
	if len(m.Sources) > 0 {
		m.Listed, m.Match = true, "host"
	}
	return m
}

// LocalFeedsLoaded reports whether any feed has entries, so callers can tell
// "not listed" from "no feeds to check".
func LocalFeedsLoaded() bool {
	local.mu.RLock()
	defer local.mu.RUnlock()
	for _, urls := range local.urls {
		if len(urls) > 0 {
			return true
		}
	}
	return false
}

// FeedLoaded reports whether the named feed ("PhishTank") has entries.
func FeedLoaded(name string) bool {
	local.mu.RLock()
	defer local.mu.RUnlock()
	return len(local.urls[name]) > 0
}

func feedsDir() string {
	if d := os.Getenv("FEEDS_DIR"); d != "" {
		return d
	}
	return filepath.Join("data", "feeds")
}

// StartLocalFeeds loads each enabled feed from its cached copy on disk, then
// keeps it fresh in the background until ctx is done. Downloads happen off
// the startup path, so a slow or failing feed never delays the server.
func StartLocalFeeds(ctx context.Context) {
	dir := feedsDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Warn("feeds dir unavailable, feeds stay in memory only", "dir", dir, "err", err)
	}
	if LookupsOff() {
		logger.Warn("threat lookups are off (THREAT_LOOKUPS=0): no lists, PhishTank, Safe Browsing or Web Risk")
		return
	}
	for _, f := range feeds() {
		if !f.enabled() {
			continue
		}
		path := filepath.Join(dir, strings.ToLower(f.name)+".txt")
		age := loadCached(f, path)
		go keepFresh(ctx, f, path, age)
	}
}

// loadCached loads a feed's saved list and returns how old it is (or a very
// large age when there's no copy).
func loadCached(f feed, path string) time.Duration {
	info, err := os.Stat(path)
	if err != nil {
		return 1<<63 - 1
	}
	file, err := os.Open(path)
	if err != nil {
		return 1<<63 - 1
	}
	defer file.Close()
	list, err := parseURLLines(file)
	if err != nil {
		return 1<<63 - 1
	}
	local.set(f.name, list)
	logger.Info("threat feed loaded from disk", "feed", f.name, "entries", len(list))
	return time.Since(info.ModTime())
}

func keepFresh(ctx context.Context, f feed, path string, age time.Duration) {
	wait := f.every - age
	if wait < 0 {
		wait = 0
	}
	timer := time.NewTimer(wait)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		next := f.every
		if err := refresh(ctx, f, path); err != nil {
			logger.Warn("threat feed refresh failed", "feed", f.name, "err", err)
			next = 30 * time.Minute // retry sooner, but not hammer the source
		}
		timer.Reset(next)
	}
}

// feedClient allows feed servers a slow handshake: abuse.ch regularly takes
// longer than Go's default 10 seconds, and a miss means a 30-minute wait.
var feedClient = func() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSHandshakeTimeout = 30 * time.Second
	t.ResponseHeaderTimeout = 60 * time.Second
	return &http.Client{Transport: t}
}()

// refresh downloads a feed, swaps it in, and saves it as one URL per line.
func refresh(ctx context.Context, f feed, path string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", f.url(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent())
	for k, v := range f.header {
		if v != "" {
			req.Header.Set(k, v)
		}
	}
	resp, err := feedClient.Do(req)
	if err != nil {
		// Feed URLs can carry an API key; log the cause, not the URL.
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	list, err := f.parse(io.LimitReader(resp.Body, 200<<20))
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("empty feed")
	}
	local.set(f.name, list)
	logger.Info("threat feed refreshed", "feed", f.name, "entries", len(list))

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(list, "\n")+"\n"), 0o644); err != nil {
		return nil // in memory is what matters; the disk copy only speeds up restarts
	}
	_ = os.Rename(tmp, path)
	return nil
}

func userAgent() string {
	if ua := os.Getenv("PHISHTANK_USER_AGENT"); ua != "" {
		return ua
	}
	return "phishtank/urlvet"
}

// parseURLLines reads one URL per line, skipping blanks and # comments.
func parseURLLines(r io.Reader) ([]string, error) {
	var list []string
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		list = append(list, line)
	}
	return list, sc.Err()
}

// parsePhishTankDump reads PhishTank's bzip2'd JSON array, streaming it so
// the ~40 MB of JSON is never held in memory at once.
func parsePhishTankDump(r io.Reader) ([]string, error) {
	dec := json.NewDecoder(bzip2.NewReader(r))
	if _, err := dec.Token(); err != nil { // opening [
		return nil, err
	}
	var list []string
	for dec.More() {
		var entry struct {
			URL string `json:"url"`
		}
		if err := dec.Decode(&entry); err != nil {
			return nil, err
		}
		if entry.URL != "" {
			list = append(list, entry.URL)
		}
	}
	return list, nil
}

// parseURLhausCSV reads URLhaus's recent export (the last 30 days) and keeps
// the URLs still online. Columns: id, dateadded, url, url_status, last_online,
// threat, tags, urlhaus_link, reporter; lines starting with # are comments.
func parseURLhausCSV(r io.Reader) ([]string, error) {
	var lines strings.Builder
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1<<20)
	for sc.Scan() {
		if line := sc.Text(); !strings.HasPrefix(line, "#") {
			lines.WriteString(line)
			lines.WriteByte('\n')
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	cr := csv.NewReader(strings.NewReader(lines.String()))
	cr.FieldsPerRecord = -1
	var list []string
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) > 3 && rec[3] == "online" {
			list = append(list, rec[2])
		}
	}
	return list, nil
}

// Exported lists for clients (the browser extension), so they can block
// known-bad pages locally, before they load, without sending anything.

// exportable are the feeds whose terms allow passing their data on.
// OpenPhish's community feed forbids redistribution, so it never leaves the server.
var exportable = map[string]bool{"PhishTank": true, "URLhaus": true}

// PrefixLen is the length of the SHA-256 prefix of each exported key.
// 8 bytes keeps collisions negligible at this list size.
const PrefixLen = 8

// PrefixExport is the hashed list of known-bad links.
type PrefixExport struct {
	Prefixes []byte   // sorted, PrefixLen bytes each
	Count    int      // number of prefixes
	Sources  []string // feeds included
	Version  string   // changes whenever the list does
}

var (
	exportMu    sync.Mutex
	exportCache PrefixExport
	exportGen   = -1
)

// FeedKey is the normalized form of a URL used for matching. Clients hash the
// same form: no scheme, lowercase host without "www.", no trailing slash or
// fragment, query kept.
func FeedKey(rawURL string) string {
	key, _ := feedKey(rawURL)
	return key
}

// ExportPrefixes returns the SHA-256 prefixes of every exportable feed's
// links, sorted and deduplicated. It's recomputed only after a feed changes.
func ExportPrefixes() PrefixExport {
	local.mu.RLock()
	gen := local.gen
	local.mu.RUnlock()

	exportMu.Lock()
	defer exportMu.Unlock()
	if gen == exportGen {
		return exportCache
	}

	local.mu.RLock()
	var sources []string
	seen := map[[PrefixLen]byte]struct{}{}
	for name, urls := range local.urls {
		if !exportable[name] || len(urls) == 0 {
			continue
		}
		sources = append(sources, name)
		for key := range urls {
			sum := sha256.Sum256([]byte(key))
			var p [PrefixLen]byte
			copy(p[:], sum[:PrefixLen])
			seen[p] = struct{}{}
		}
	}
	local.mu.RUnlock()

	list := make([][PrefixLen]byte, 0, len(seen))
	for p := range seen {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool { return bytes.Compare(list[i][:], list[j][:]) < 0 })
	flat := make([]byte, 0, len(list)*PrefixLen)
	for _, p := range list {
		flat = append(flat, p[:]...)
	}
	sort.Strings(sources)
	v := sha256.Sum256(flat)

	exportCache = PrefixExport{
		Prefixes: flat,
		Count:    len(list),
		Sources:  sources,
		Version:  hex.EncodeToString(v[:8]),
	}
	exportGen = gen
	return exportCache
}
