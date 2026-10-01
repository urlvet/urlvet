package checks

import (
	"regexp"
	"strings"

	"github.com/urlvet/urlvet/internal/constants"
)

func CheckURLKeywords(url string) (bool, []string, map[string][]string) {
	lowerURL := strings.ToLower(url)

	// Split by non-alphanumeric characters
	re := regexp.MustCompile(`[^a-z0-9]+`)
	words := re.Split(lowerURL, -1)

	matches := []string{}
	categories := make(map[string][]string)
	keywordPresent := false

	for _, word := range words {
		if word == "" {
			continue
		}
		if category, exists := constants.UrlKeywords[word]; exists {
			keywordPresent = true
			matches = append(matches, word)
			categories[category] = append(categories[category], word)
		}
	}

	return keywordPresent, matches, categories
}
