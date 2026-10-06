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
		return rdapData, setAge(rdapData)
	}

	// RDAP failed, fall back to WHOIS
	whoisData, err := GetWhoisData(domain)
	if err != nil {
		return nil, err
	}

	return whoisData, setAge(whoisData)
}

// setAge fills in the age from the creation date, leaving it unset when the
// registry didn't publish one.
func setAge(d *RegistrationData) error {
	if d.CreatedDate.IsZero() {
		return nil
	}
	ageHuman, ageDays, err := checks.GetDomainAge(d.CreatedDate)
	if err != nil {
		return err
	}
	d.AgeKnown = true
	d.AgeHuman = ageHuman
	d.AgeDays = &ageDays
	return nil
}
