package checks

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/urlvet/urlvet/internal/logger"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// FormInfo describes a discovered form
type FormInfo struct {
	Action           string   `json:"action"` // resolved absolute action URL (may be empty -> current URL)
	Method           string   `json:"method"` // GET/POST etc.
	Inputs           []string `json:"inputs"` // raw input names/types/placeholder snippets
	ContainsPassword bool     `json:"has_password"`
	ContainsUserLike bool     `json:"has_user_like"` // username/email-like inputs
	ContainsPayment  bool     `json:"has_payment"`   // credit card, money, etc.
	ContainsPersonal bool     `json:"has_personal"`  // address, phone, ssn, etc.
	SubmitTexts      []string `json:"submit_texts"`
	ExternalAction   bool     `json:"is_external"` // action host != page host (true if external)
	IsHidden         bool     `json:"is_hidden"`
}

type IframeInfo struct {
	Source   string `json:"src"`
	IsHidden bool   `json:"is_hidden"`
	Width    string `json:"width"`
	Height   string `json:"height"`
}

// PageFormResult summarizes page-level findings
type PageFormResult struct {
	URL             string        `json:"url"`
	Title           string        `json:"title"`
	HasForms        bool          `json:"has_forms"`
	HasLoginForm    bool          `json:"has_login_form"`
	HasPaymentForm  bool          `json:"has_payment_form"`
	HasPersonalForm bool          `json:"has_personal_form"`
	FormCount       int           `json:"form_count"`
	Forms           []FormInfo    `json:"forms"`
	Iframes         []IframeInfo  `json:"iframes"`
	HasHiddenIframe bool          `json:"has_hidden_iframe"`
	HasTracking     bool          `json:"has_tracking"`                         // detection of 1x1 pixels etc.
	FetchDuration   time.Duration `json:"fetch_duration" swaggertype:"integer"` // nanoseconds
	BrandCheck      BrandResult   `json:"brand_check"`
}

// hostOnly returns scheme://host from a URL, used to keep full URLs out of logs.
func hostOnly(rawURL string) string {
	if u, err := url.Parse(rawURL); err == nil {
		return u.Scheme + "://" + u.Host
	}
	return rawURL
}

// GetPageFormInfo fetches the page (with timeout), parses HTML and returns info.
func GetPageFormInfo(pageURL string) (*PageFormResult, error) {
	start := time.Now()
	logger.Debug("starting content analysis", "host", hostOnly(pageURL))

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		logger.Error("failed to create request", "host", hostOnly(pageURL), "err", err)
		return nil, err
	}
	// Use a common browser User-Agent to avoid blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := newSafeHTTPClient(10 * time.Second)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("HTTP request failed", "host", hostOnly(pageURL), "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("non-200 status code", "host", hostOnly(pageURL), "status", resp.StatusCode)
	}

	// read body (limited)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024)) // 5MB cap
	if err != nil {
		logger.Error("failed to read body", "host", hostOnly(pageURL), "err", err)
		return nil, err
	}

	logger.Debug("fetched page content", "bytes", len(body), "host", hostOnly(pageURL))

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		logger.Error("failed to parse HTML", "host", hostOnly(pageURL), "err", err)
		return nil, err
	}

	uParsed, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	pageHost := uParsed.Hostname()

	var results []FormInfo
	var iframes []IframeInfo
	var pageTitle string
	var hasTracking bool
	var bodyText strings.Builder

	// traverse nodes to find <form>, <title>, <iframe> and collect body text
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "form":
				f := extractFormInfo(n, uParsed, pageHost)
				results = append(results, f)
			case "iframe":
				src := getAttr(n, "src")
				w := getAttr(n, "width")
				h := getAttr(n, "height")
				style := getAttr(n, "style")
				hidden := isHidden(n, style, w, h)
				iframes = append(iframes, IframeInfo{
					Source:   src,
					IsHidden: hidden,
					Width:    w,
					Height:   h,
				})
			case "img":
				w := getAttr(n, "width")
				h := getAttr(n, "height")
				if w == "1" && h == "1" || w == "0" && h == "0" {
					hasTracking = true
				}
			case "title":
				// Inline SVGs carry their own <title> (e.g. logo names); only the
				// document title counts, and the first one wins.
				if pageTitle == "" && n.FirstChild != nil && !hasAncestor(n, "svg") {
					pageTitle = n.FirstChild.Data
				}
			}
		} else if n.Type == html.TextNode {
			bodyText.WriteString(n.Data + " ")
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	hasHiddenIframe := false
	for _, ifr := range iframes {
		if ifr.IsHidden {
			hasHiddenIframe = true
			break
		}
	}

	res := &PageFormResult{
		URL:             pageURL,
		Title:           strings.TrimSpace(pageTitle),
		HasForms:        len(results) > 0,
		HasLoginForm:    false,
		HasPaymentForm:  false,
		HasPersonalForm: false,
		FormCount:       len(results),
		Forms:           results,
		Iframes:         iframes,
		HasHiddenIframe: hasHiddenIframe,
		HasTracking:     hasTracking,
		FetchDuration:   time.Since(start),
	}

	// mark if any form looks like a login form
	for _, f := range results {
		if f.ContainsPassword || f.ContainsUserLike {
			res.HasLoginForm = true
		}
		if f.ContainsPayment {
			res.HasPaymentForm = true
		}
		if f.ContainsPersonal {
			res.HasPersonalForm = true
		}
	}

	// Brand Check
	res.BrandCheck = CheckBrandMismatch(pageHost, res.Title)

	return res, nil
}

// extractFormInfo inspects a <form> node and returns a FormInfo
func extractFormInfo(form *html.Node, base *url.URL, pageHost string) FormInfo {
	style := getAttr(form, "style")
	w := getAttr(form, "width")
	h := getAttr(form, "height")

	info := FormInfo{
		Method:   strings.ToUpper(getAttr(form, "method")),
		IsHidden: isHidden(form, style, w, h),
	}
	if info.Method == "" {
		info.Method = "GET"
	}
	rawAction := getAttr(form, "action")
	actionResolved := resolveAction(base, rawAction)
	info.Action = actionResolved

	// determine if action host differs
	if actionResolved != "" {
		if parsed, err := url.Parse(actionResolved); err == nil {
			info.ExternalAction = !sameHost(parsed.Hostname(), pageHost)
		}
	}

	// collect inputs and buttons inside the form
	var ftr func(*html.Node)
	ftr = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "input":
				it := strings.ToLower(getAttr(n, "type"))
				name := getAttr(n, "name")
				placeholder := getAttr(n, "placeholder")
				aria := getAttr(n, "aria-label")
				id := getAttr(n, "id")
				info.Inputs = append(info.Inputs, fmtInputSummary(it, name, placeholder, aria, id))

				// password detection
				if it == "password" || strings.Contains(strings.ToLower(name), "pass") || strings.Contains(strings.ToLower(id), "pass") || strings.Contains(strings.ToLower(placeholder), "pass") {
					info.ContainsPassword = true
				}
				// username/email-ish detection
				if it == "email" ||
					strings.Contains(strings.ToLower(name), "user") ||
					strings.Contains(strings.ToLower(name), "login") ||
					strings.Contains(strings.ToLower(name), "email") ||
					strings.Contains(strings.ToLower(id), "user") ||
					strings.Contains(strings.ToLower(placeholder), "user") ||
					strings.Contains(strings.ToLower(aria), "user") ||
					strings.Contains(strings.ToLower(placeholder), "email") {
					info.ContainsUserLike = true
				}

				// Payment/Money detection
				lName := strings.ToLower(name)
				lId := strings.ToLower(id)
				lPh := strings.ToLower(placeholder)
				if strings.Contains(lName, "card") || strings.Contains(lId, "card") || strings.Contains(lPh, "card") ||
					strings.Contains(lName, "cvv") || strings.Contains(lId, "cvv") ||
					strings.Contains(lName, "expiry") || strings.Contains(lId, "expiry") ||
					strings.Contains(lName, "credit") || strings.Contains(lId, "credit") ||
					strings.Contains(lName, "money") || strings.Contains(lName, "pay") || strings.Contains(lName, "billing") ||
					strings.Contains(lPh, "checkout") || strings.Contains(lPh, "payment") {
					info.ContainsPayment = true
				}

				// Personal Info detection
				if strings.Contains(lName, "address") || strings.Contains(lId, "address") || strings.Contains(lPh, "address") ||
					strings.Contains(lName, "phone") || strings.Contains(lId, "phone") || strings.Contains(lPh, "phone") ||
					strings.Contains(lName, "ssn") || strings.Contains(lId, "ssn") ||
					strings.Contains(lName, "dob") || strings.Contains(lName, "birth") ||
					strings.Contains(lName, "city") || strings.Contains(lName, "zip") || strings.Contains(lName, "state") {
					info.ContainsPersonal = true
				}
			case "button":
				// capture text content of button
				txt := strings.TrimSpace(nodeText(n))
				if txt != "" {
					info.SubmitTexts = append(info.SubmitTexts, txt)
					l := strings.ToLower(txt)
					if looksLikeLoginText(l) {
						info.ContainsUserLike = true
					}
					if strings.Contains(l, "pay") || strings.Contains(l, "buy") || strings.Contains(l, "checkout") || strings.Contains(l, "order") {
						info.ContainsPayment = true
					}
				}
			case "a":
				// sometimes login is an <a> styled as button
				txt := strings.TrimSpace(nodeText(n))
				if txt != "" {
					l := strings.ToLower(txt)
					if looksLikeLoginText(l) {
						info.SubmitTexts = append(info.SubmitTexts, txt)
					}
					if strings.Contains(l, "pay now") || strings.Contains(l, "pay $") ||
						strings.Contains(l, "complete payment") || strings.Contains(l, "checkout") {
						info.ContainsPayment = true
					}
				}
			case "label":
				// label text may indicate username/password fields
				txt := strings.ToLower(strings.TrimSpace(nodeText(n)))
				if strings.Contains(txt, "password") {
					info.ContainsPassword = true
				}
				if strings.Contains(txt, "username") || strings.Contains(txt, "email") || strings.Contains(txt, "sign in") {
					info.ContainsUserLike = true
				}
				if strings.Contains(txt, "card") || strings.Contains(txt, "cvv") || strings.Contains(txt, "expiry") || strings.Contains(txt, "credit") {
					info.ContainsPayment = true
				}
				if strings.Contains(txt, "address") || strings.Contains(txt, "phone") || strings.Contains(txt, "zip") {
					info.ContainsPersonal = true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ftr(c)
		}
	}
	ftr(form)

	// additional heuristic: if form has a password OR (text inputs + "sign in" in action or form text) -> login-like
	if !info.ContainsUserLike {
		// check action path / query for login keywords
		if strings.Contains(strings.ToLower(info.Action), "login") || strings.Contains(strings.ToLower(info.Action), "signin") || strings.Contains(strings.ToLower(info.Action), "auth") {
			info.ContainsUserLike = true
		}
	}
	return info
}

// helpers

func sameHost(a, b string) bool {
	if strings.EqualFold(a, b) {
		return true
	}
	// Compare by registered domain so www.nasa.gov and nasa.gov are treated the same.
	domA, errA := publicsuffix.EffectiveTLDPlusOne(strings.ToLower(a))
	domB, errB := publicsuffix.EffectiveTLDPlusOne(strings.ToLower(b))
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(domA, domB)
}

func resolveAction(base *url.URL, raw string) string {
	if raw == "" || raw == "#" {
		// empty action means submit to same URL
		return base.String()
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return base.ResolveReference(parsed).String()
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

// nodeText returns concatenated text of a node subtree
func nodeText(n *html.Node) string {
	var b strings.Builder
	var walker func(*html.Node)
	walker = func(nn *html.Node) {
		if nn.Type == html.ElementNode && (nn.Data == "style" || nn.Data == "script") {
			return
		}
		if nn.Type == html.TextNode {
			b.WriteString(nn.Data)
		}
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(n)
	return strings.TrimSpace(b.String())
}

func fmtInputSummary(typ, name, placeholder, aria, id string) string {
	parts := []string{}
	if typ != "" {
		parts = append(parts, "type="+typ)
	}
	if name != "" {
		parts = append(parts, "name="+name)
	}
	if id != "" {
		parts = append(parts, "id="+id)
	}
	if placeholder != "" {
		parts = append(parts, "ph="+placeholder)
	}
	if aria != "" {
		parts = append(parts, "aria="+aria)
	}
	return strings.Join(parts, "|")
}

func looksLikeLoginText(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	keywords := []string{"login", "log in", "sign in", "signin", "submit", "sign-on", "signon", "sign up", "signup"}
	for _, k := range keywords {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

func isHidden(n *html.Node, style string, w string, h string) bool {
	// Check hidden attribute
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, "hidden") || strings.EqualFold(a.Key, "type") && strings.EqualFold(a.Val, "hidden") {
			return true
		}
	}

	// Check style
	style = strings.ToLower(style)
	if strings.Contains(style, "display:none") || strings.Contains(style, "visibility:hidden") ||
		strings.Contains(style, "opacity:0") || strings.Contains(style, "width:0") || strings.Contains(style, "height:0") {
		return true
	}

	// Check dimensions
	if w == "0" || h == "0" {
		return true
	}

	return false
}

// hasAncestor reports whether n is nested inside an element with the given tag.
func hasAncestor(n *html.Node, tag string) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && p.Data == tag {
			return true
		}
	}
	return false
}
