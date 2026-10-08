package checks

import (
	neturl "net/url"
	"strings"
)

// maxURLLength is where a URL starts to look padded. Ordinary links with a
// few tracking parameters run past 100 characters, so the bar sits well above.
const maxURLLength = 150

// maxPathSegments is where a path starts to look like it's burying the real
// page. Release downloads and docs pages routinely reach five or six.
const maxPathSegments = 6

func TooLongUrl(url string) bool {
	return len(url) > maxURLLength
}

// TooDeepUrl counts path segments only, so the slashes in "https://" and in
// query strings don't count towards depth.
func TooDeepUrl(url string) bool {
	u, err := neturl.Parse(url)
	if err != nil {
		return false
	}
	depth := 0
	for _, seg := range strings.Split(u.Path, "/") {
		if seg != "" {
			depth++
		}
	}
	return depth > maxPathSegments
}
