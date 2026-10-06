package checks

import (
	"net"
	"regexp"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
)

// dashedIP finds an IPv4 address written with dots or dashes inside a name
// (1-2-3-4, 1.2.3.4), as IP hostname services read them.
var dashedIP = regexp.MustCompile(`(?:^|[.-])(\d{1,3}[.-]\d{1,3}[.-]\d{1,3}[.-]\d{1,3})(?:[.-]|$)`)

// IPHostname reports whether host is on a service that turns any name with an
// IP address in it into that address (webmail.50-6-22-93.nip.io is
// 50.6.22.93). It returns the service and the address, or "" for the address
// when none can be read from the name.
func IPHostname(host string) (service, ip string, ok bool) {
	host = strings.ToLower(host)
	for s := range constants.IPHostnameServices {
		if !strings.HasSuffix(host, "."+s) {
			continue
		}
		rest := strings.TrimSuffix(host, "."+s)
		for _, m := range dashedIP.FindAllStringSubmatch(rest, -1) {
			if parsed := net.ParseIP(strings.ReplaceAll(m[1], "-", ".")); parsed != nil {
				return s, parsed.String(), true
			}
		}
		return s, "", true
	}
	return "", "", false
}
