package domaininfo

import (
	"context"

	"github.com/urlvet/urlvet/internal/service/checks"
)

// Lookup tries RDAP first, falls back to WHOIS if RDAP fails.
// Uses context for timeout/cancellation support.
func Lookup(domain string) (*RegistrationData, error) {
	return LookupWithContext(context.Background(), domain)
}

// LookupWithContext tries RDAP first with timeout, falls back to WHOIS if RDAP fails.
func LookupWithContext(ctx context.Context, domain string) (*RegistrationData, error) {
	// Try RDAP first with timeout
	rdapData, err := fetchRDAPWithContext(ctx, domain)
	if err == nil && rdapData != nil {
		ageHuman, ageDays, err := checks.GetDomainAge(rdapData.CreatedDate)
		if err != nil {
			return rdapData, err
		}
		rdapData.AgeHuman = ageHuman
		rdapData.AgeDays = ageDays
		return rdapData, nil
	}

	// RDAP failed, fall back to WHOIS
	whoisData, err := GetWhoisData(domain)
	if err != nil {
		return nil, err
	}

	ageHuman, ageDays, err := checks.GetDomainAge(whoisData.CreatedDate)
	if err != nil {
		return whoisData, err
	}
	whoisData.AgeHuman = ageHuman
	whoisData.AgeDays = ageDays

	return whoisData, nil
}
