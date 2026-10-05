package threatfeeds

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
)

// Google Safe Browsing v5 (https://developers.google.com/safe-browsing):
// free, no billing account, non-commercial use only. Runs when
// SAFE_BROWSING_API_KEY is set.
//
// It uses hashes:search, not urls:search: Google receives 4-byte prefixes of
// SHA-256 hashes of the URL's host/path combinations, never the URL itself.
// Google returns every full hash sharing those prefixes, and the match is
// made here.

const safeBrowsingURL = "https://safebrowsing.googleapis.com/v5/hashes:search"

// ErrSafeBrowsingRateLimited means Google's quota for the key was hit.
var ErrSafeBrowsingRateLimited = errors.New("safe browsing: rate limited")

// SafeBrowsingEnabled reports whether lookups are configured.
func SafeBrowsingEnabled() bool { return os.Getenv("SAFE_BROWSING_API_KEY") != "" && !LookupsOff() }

// CheckSafeBrowsing looks a URL up in Google Safe Browsing.
func CheckSafeBrowsing(rawURL string) (*GoogleThreatResult, error) {
	key := os.Getenv("SAFE_BROWSING_API_KEY")
	if key == "" {
		return nil, errors.New("safe browsing: no API key")
	}
	expressions, err := sbExpressions(rawURL)
	if err != nil {
		return nil, fmt.Errorf("safe browsing: %w", err)
	}

	full := make(map[[32]byte]bool, len(expressions))
	seen := map[[4]byte]bool{}
	q := url.Values{}
	for _, e := range expressions {
		h := sha256.Sum256([]byte(e))
		full[h] = true
		var prefix [4]byte
		copy(prefix[:], h[:4])
		if !seen[prefix] {
			seen[prefix] = true
			q.Add("hashPrefixes", base64.StdEncoding.EncodeToString(prefix[:]))
		}
	}
	q.Set("key", key)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(safeBrowsingURL + "?" + q.Encode())
	if err != nil {
		var ue *url.Error // its message would include the API key
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return nil, fmt.Errorf("safe browsing: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrSafeBrowsingRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("safe browsing: HTTP %d", resp.StatusCode)
	}

	// v5 answers only in protobuf (JSON is "Unsupported Output Format").
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("safe browsing: %w", err)
	}
	hashes, err := decodeSearchHashesResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("safe browsing: %w", err)
	}

	res := &GoogleThreatResult{}
	types := map[string]bool{}
	for _, fh := range hashes {
		if !full[fh.hash] {
			continue // shares a prefix with our URL, but is someone else's
		}
		for _, t := range fh.threatTypes {
			if !types[t] {
				types[t] = true
				res.ThreatTypes = append(res.ThreatTypes, t)
			}
		}
	}
	res.Listed = len(res.ThreatTypes) > 0
	return res, nil
}

// sbThreatTypes are v5's ThreatType enum values, by number.
var sbThreatTypes = map[uint64]string{
	1: "MALWARE",
	2: "SOCIAL_ENGINEERING",
	3: "UNWANTED_SOFTWARE",
	4: "POTENTIALLY_HARMFUL_APPLICATION",
}

type sbFullHash struct {
	hash        [32]byte
	threatTypes []string
}

// decodeSearchHashesResponse reads the protobuf SearchHashesResponse:
//
//	message SearchHashesResponse { repeated FullHash full_hashes = 1; Duration cache_duration = 2; }
//	message FullHash { bytes full_hash = 1; repeated FullHashDetail full_hash_details = 2; }
//	message FullHashDetail { ThreatType threat_type = 1; repeated ThreatAttribute attributes = 2; }
//
// Fields it doesn't need are skipped.
func decodeSearchHashesResponse(b []byte) ([]sbFullHash, error) {
	var out []sbFullHash
	err := eachField(b, func(num protowire.Number, typ protowire.Type, v []byte, _ uint64) error {
		if num != 1 || typ != protowire.BytesType {
			return nil
		}
		var fh sbFullHash
		var hashLen int
		err := eachField(v, func(num protowire.Number, typ protowire.Type, v []byte, _ uint64) error {
			switch {
			case num == 1 && typ == protowire.BytesType:
				hashLen = copy(fh.hash[:], v)
				if len(v) != 32 {
					hashLen = 0
				}
			case num == 2 && typ == protowire.BytesType:
				return eachField(v, func(num protowire.Number, typ protowire.Type, _ []byte, n uint64) error {
					// Unknown or unspecified types are to be ignored, per the API.
					if t, ok := sbThreatTypes[n]; ok && num == 1 && typ == protowire.VarintType {
						fh.threatTypes = append(fh.threatTypes, t)
					}
					return nil
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
		if hashLen == 32 {
			out = append(out, fh)
		}
		return nil
	})
	return out, err
}

// eachField walks one protobuf message's fields, passing length-delimited
// values as bytes and varints as numbers.
func eachField(b []byte, fn func(num protowire.Number, typ protowire.Type, bytesVal []byte, varint uint64) error) error {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
		var bytesVal []byte
		var varint uint64
		switch typ {
		case protowire.BytesType:
			bytesVal, n = protowire.ConsumeBytes(b)
		case protowire.VarintType:
			varint, n = protowire.ConsumeVarint(b)
		default:
			n = protowire.ConsumeFieldValue(num, typ, b)
		}
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
		if err := fn(num, typ, bytesVal, varint); err != nil {
			return err
		}
	}
	return nil
}

// sbExpressions canonicalizes a URL and returns the host-suffix/path-prefix
// expressions Safe Browsing hashes ("a.b.example.com/1/2.html?x=1",
// "example.com/", …), at most 30, as its URL hashing spec describes.
func sbExpressions(rawURL string) ([]string, error) {
	host, path, query, err := sbCanonicalize(rawURL)
	if err != nil {
		return nil, err
	}

	hosts := []string{host}
	if net.ParseIP(host) == nil {
		labels := strings.Split(host, ".")
		// The last five components, then drop leading ones, never down to the TLD alone.
		start := len(labels) - 5
		if start < 1 {
			start = 1
		}
		for i := start; i < len(labels)-1; i++ {
			hosts = append(hosts, strings.Join(labels[i:], "."))
		}
	}

	var paths []string
	if query != "" {
		paths = append(paths, path+query)
	}
	paths = append(paths, path)
	// "/", then each leading directory, up to four in all.
	prefix := "/"
	paths = append(paths, prefix)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 0; i < len(parts)-1 && i < 3; i++ {
		if parts[i] == "" {
			break
		}
		prefix += parts[i] + "/"
		paths = append(paths, prefix)
	}

	seen := map[string]bool{}
	var out []string
	for _, h := range hosts {
		for _, p := range paths {
			e := h + p
			if !seen[e] && len(out) < 30 {
				seen[e] = true
				out = append(out, e)
			}
		}
	}
	return out, nil
}

// sbCanonicalize returns the host, path and query (with its "?") after Safe
// Browsing's canonicalization: no fragment,
// repeated unescaping, lowercase host with stray dots removed, IPs in
// dotted-decimal, "." and ".." resolved, repeated slashes collapsed, then
// escaping of control, non-ASCII, "#" and "%" characters.
func sbCanonicalize(rawURL string) (host, path, query string, err error) {
	s := strings.TrimSpace(rawURL)
	s = strings.NewReplacer("\t", "", "\r", "", "\n", "").Replace(s)
	if i := strings.Index(s, "#"); i >= 0 {
		s = s[:i]
	}
	if !strings.Contains(s, "://") {
		s = "http://" + s
	}
	s = sbUnescape(s)

	// Split by hand: after unescaping, the URL can hold literal "%" and "#",
	// which url.Parse would reject or treat as a fragment.
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	rest := ""
	if i := strings.IndexAny(s, "/?"); i >= 0 {
		s, rest = s[:i], s[i:]
	}
	if i := strings.LastIndex(s, "@"); i >= 0 {
		s = s[i+1:] // user:pass@
	}
	if h, _, err := net.SplitHostPort(s); err == nil {
		s = h
	}
	host = strings.Trim(strings.ToLower(strings.Trim(s, "[]")), ".")
	for strings.Contains(host, "..") {
		host = strings.ReplaceAll(host, "..", ".")
	}
	if ip := net.ParseIP(host); ip != nil && ip.To4() != nil {
		host = ip.To4().String()
	}
	if host == "" {
		return "", "", "", errors.New("no host")
	}

	path = rest
	if i := strings.Index(rest, "?"); i >= 0 {
		path, query = rest[:i], rest[i:] // keeps the "?", even with nothing after it
	}
	if path == "" {
		path = "/"
	}
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	trailing := strings.HasSuffix(path, "/")
	var segs []string
	for _, seg := range strings.Split(path, "/") {
		switch seg {
		case "", ".":
		case "..":
			if len(segs) > 0 {
				segs = segs[:len(segs)-1]
			}
		default:
			segs = append(segs, seg)
		}
	}
	path = "/" + strings.Join(segs, "/")
	if trailing && path != "/" {
		path += "/"
	}

	return sbEscape(host), sbEscape(path), sbEscape(query), nil
}

// sbUnescape percent-decodes until nothing changes (bounded). Invalid
// escapes such as "%%" are left as they are, as the spec requires; Go's own
// unescapers reject the whole string instead.
func sbUnescape(s string) string {
	for i := 0; i < 10; i++ {
		d := sbUnescapeOnce(s)
		if d == s {
			return s
		}
		s = d
	}
	return s
}

func sbUnescapeOnce(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
			b.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	}
	return c - 'A' + 10
}

// sbEscape escapes bytes <= 32, >= 127, "#" and "%".
func sbEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c <= 32 || c >= 127 || c == '#' || c == '%' {
			fmt.Fprintf(&b, "%%%02X", c)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}
