package checks

import (
	"mime"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// ScriptRedirect describes a page that sends visitors elsewhere on load, by
// script or <meta http-equiv="refresh">, instead of an HTTP redirect.
type ScriptRedirect struct {
	Target      string `json:"target,omitempty"` // literal destination, when the page names one
	ToIP        bool   `json:"to_ip"`            // destination is a raw IP address
	CrossDomain bool   `json:"cross_domain"`     // destination is on another site
}

var (
	// location = …, location.href = …, location.replace(…), location.assign(…)
	jsRedirectRe = regexp.MustCompile(`location(?:\.href)?\s*=[^=]|location\.(?:replace|assign)\s*\(`)
	// The same, capturing a quoted literal destination.
	jsRedirectTargetRe = regexp.MustCompile(`location(?:\.href)?\s*=\s*["']([^"']+)["']|location\.(?:replace|assign)\s*\(\s*["']([^"']+)["']`)
	// A quoted or URL-embedded IPv4 address: "185.80.128.4", '//1.2.3.4/'.
	quotedIPv4Re = regexp.MustCompile(`["'/](\d{1,3}(?:\.\d{1,3}){3})(?::\d+)?["'/:?]`)
	// <meta http-equiv="refresh" content="0; url=…">
	metaRefreshURLRe = regexp.MustCompile(`(?i)url\s*=\s*['"]?([^'";\s]+)`)
)

// redirectorTextLimit is how much visible text a page can have and still read
// as a bare redirector rather than a page with links that happen to use script.
const redirectorTextLimit = 200

// detectScriptRedirect looks for a load-time redirect in inline scripts and a
// meta refresh URL. Spam redirectors often build the destination from pieces
// ('http://' + ip + '/?' + params), so a raw IP anywhere in a redirecting
// script counts, not only a literal URL. A script sending visitors to another
// domain only counts on a page with next to no text of its own: on a full
// page it's usually a button handler.
func detectScriptRedirect(scripts []string, metaRefresh string, pageHost string, visibleText int) *ScriptRedirect {
	var found *ScriptRedirect
	thinPage := visibleText < redirectorTextLimit
	consider := func(target string, fromScript bool) {
		u, err := url.Parse(strings.TrimSpace(target))
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "") {
			return
		}
		host := u.Hostname()
		r := &ScriptRedirect{Target: target}
		r.ToIP = isPublicIP(host)
		r.CrossDomain = !sameHost(host, pageHost)
		if !r.ToIP && (!r.CrossDomain || (fromScript && !thinPage)) {
			return
		}
		if found == nil || (r.ToIP && !found.ToIP) {
			found = r
		}
	}

	if metaRefresh != "" {
		if m := metaRefreshURLRe.FindStringSubmatch(metaRefresh); m != nil {
			consider(m[1], false)
		}
	}

	for _, js := range scripts {
		if !jsRedirectRe.MatchString(js) {
			continue
		}
		for _, m := range jsRedirectTargetRe.FindAllStringSubmatch(js, -1) {
			consider(m[1]+m[2], true)
		}
		if found != nil && found.ToIP {
			continue
		}
		for _, m := range quotedIPv4Re.FindAllStringSubmatch(js, -1) {
			if isPublicIP(m[1]) {
				found = &ScriptRedirect{Target: m[1], ToIP: true, CrossDomain: true}
				break
			}
		}
	}
	return found
}

func isPublicIP(host string) bool {
	ip := net.ParseIP(host)
	return ip != nil && ip.To4() != nil && !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsUnspecified()
}

// isHTMLResponse reports whether a response is a web page worth parsing.
// Servers that omit Content-Type are judged by sniffing the first bytes.
func isHTMLResponse(resp *http.Response, head []byte) bool {
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(head)
	}
	mt, _, _ := mime.ParseMediaType(ct)
	return mt == "text/html" || mt == "application/xhtml+xml"
}

// downloadFileName returns the file name a download would be saved under:
// the Content-Disposition filename, or the last segment of the URL path.
func downloadFileName(resp *http.Response) string {
	if _, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); err == nil && params["filename"] != "" {
		return params["filename"]
	}
	if resp.Request != nil {
		return path.Base(resp.Request.URL.Path)
	}
	return ""
}

// executableExts are file types that run code when opened. ".js", ".com" and
// ".app" are left out: they turn up in ordinary links (CDN scripts,
// "/go/example.com") far more often than in downloads.
var executableExts = map[string]struct{}{
	".exe": {}, ".msi": {}, ".scr": {}, ".bat": {}, ".cmd": {}, ".pif": {},
	".vbs": {}, ".jse": {}, ".wsf": {}, ".ps1": {}, ".hta": {}, ".jar": {},
	".apk": {}, ".xapk": {}, ".dmg": {}, ".pkg": {}, ".deb": {}, ".rpm": {},
	".appimage": {}, ".lnk": {}, ".iso": {},
}

// ExecutableExt returns the executable extension of a file name or URL
// ("VMware-Workstation.exe" → ".exe"), or "" when it isn't a program.
func ExecutableExt(name string) string {
	if u, err := url.Parse(name); err == nil && u.Path != "" {
		name = u.Path
	}
	ext := strings.ToLower(path.Ext(name))
	if _, ok := executableExts[ext]; ok {
		return ext
	}
	return ""
}

// ProviderBlock records that the host or CDN in front of a site has blocked it.
type ProviderBlock struct {
	Provider string `json:"provider"`
	Reason   string `json:"reason"` // "phishing", "malware" or "blocked"
}

// cloudflareWarningRe matches the title of Cloudflare's interstitial for sites
// it has flagged: "Suspected Phishing | Cloudflare", "Suspected Malware | Cloudflare".
var cloudflareWarningRe = regexp.MustCompile(`(?i)^suspected (phishing|malware)\s*\|\s*cloudflare$`)

// detectProviderBlock recognises a host's or CDN's own warning page in place
// of the site. These come from the provider's abuse team, so they're a
// stronger signal than anything we can infer from the page.
func detectProviderBlock(resp *http.Response, title string) *ProviderBlock {
	if m := cloudflareWarningRe.FindStringSubmatch(title); m != nil {
		return &ProviderBlock{Provider: "Cloudflare", Reason: strings.ToLower(m[1])}
	}
	// 451 is what Vercel and others serve for a deployment taken down for abuse.
	if resp.StatusCode == http.StatusUnavailableForLegalReasons {
		provider := "The hosting provider"
		if strings.Contains(strings.ToLower(resp.Header.Get("Server")), "vercel") {
			provider = "Vercel"
		}
		return &ProviderBlock{Provider: provider, Reason: "blocked"}
	}
	return nil
}

// personalWords are the words in a field name, placeholder or label that ask
// for personal details.
var personalWords = map[string]bool{
	"address": true, "street": true, "phone": true, "telephone": true, "mobile": true,
	"ssn": true, "dob": true, "birth": true, "birthday": true, "birthdate": true,
	"city": true, "zip": true, "zipcode": true, "postal": true, "postcode": true, "state": true,
}

// addressOf are words that make "address" something other than a home
// address: email address, IP address, web address.
var addressOf = map[string]bool{"email": true, "e": true, "mail": true, "ip": true, "web": true, "wallet": true}

// fieldWord splits "dateOfBirth", "billing_zip" and "Your email address" into words.
var fieldWord = regexp.MustCompile(`[A-Z]?[a-z]+|[A-Z]+(?:[a-z]+)?|[0-9]+`)

// asksPersonal reports whether a field name, placeholder or label asks for
// personal details, judged by whole words.
func asksPersonal(s string) bool {
	words := fieldWord.FindAllString(s, -1)
	for i, w := range words {
		w = strings.ToLower(w)
		if !personalWords[w] {
			continue
		}
		if w == "address" && i > 0 && addressOf[strings.ToLower(words[i-1])] {
			continue
		}
		return true
	}
	return false
}

// IsMarketingFormService reports whether a form's action is an email
// marketing platform: newsletter sign-ups post there from sites of every size.
// General form backends (Formspree, Getform, Jotform) aren't on the list:
// phishing kits use them to collect what they steal.
func IsMarketingFormService(action string) bool {
	u, err := url.Parse(action)
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Hostname())
	for _, s := range marketingFormHosts {
		if host == s || strings.HasSuffix(host, "."+s) {
			return true
		}
	}
	return false
}

var marketingFormHosts = []string{
	"list-manage.com", "mailchimp.com", // Mailchimp
	"hsforms.com", "hsforms.net", "hubspot.com", // HubSpot
	"eloqua.com",                 // Oracle Eloqua
	"marketo.com", "mktoweb.com", // Adobe Marketo
	"pardot.com",                // Salesforce Pardot
	"klaviyo.com",               // Klaviyo
	"sibforms.com", "brevo.com", // Brevo (Sendinblue)
	"mailerlite.com", "ml.email", // MailerLite
	"convertkit.com", "ck.page", "kit.com", // Kit (ConvertKit)
	"substack.com", "beehiiv.com", "buttondown.email",
	"constantcontact.com", "ctctcdn.com", // Constant Contact
	"getresponse.com", "aweber.com", "campaignmonitor.com", "createsend.com",
	"activehosted.com", "activecampaign.com",
}
