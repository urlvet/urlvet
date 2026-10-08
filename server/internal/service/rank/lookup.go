package rank

func DomainRankLookup(domain string) int {
	if rank, exists := domainRankMap[domain]; exists {
		return rank
	}
	return 0
}

// TopDomains returns up to n of the most popular domains, most popular first.
func TopDomains(n int) []string {
	if n > len(topDomains) {
		n = len(topDomains)
	}
	return topDomains[:n]
}
