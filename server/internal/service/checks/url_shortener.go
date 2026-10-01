package checks

import "github.com/urlvet/urlvet/internal/constants"

func IsUrlShortener(domain string) bool {
	_, ok := constants.URLShorteners[domain]
	return ok
}
