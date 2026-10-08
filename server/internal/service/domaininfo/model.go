package domaininfo

import "time"

// RegistrationData holds normalized domain registry info
// regardless of whether it came from WHOIS or RDAP.
type RegistrationData struct {
	Domain      string    `json:"domain"`
	Registrar   string    `json:"registrar"`
	CreatedDate time.Time `json:"created"`
	UpdatedDate time.Time `json:"updated"`
	ExpiryDate  time.Time `json:"expiry"`
	Nameservers []string  `json:"nameservers"`
	Status      []string  `json:"status"`
	DNSSEC      bool      `json:"dnssec"`
	// Some registries (.de, .eu, …) don't publish a creation date. Then
	// AgeKnown is false and AgeDays is null: a zero date would read as 2,025
	// years old, and 0 would read as "registered today".
	AgeKnown bool   `json:"age_known"`
	AgeHuman string `json:"age_human"` // e.g. "2 years 3 months"; empty when unknown
	AgeDays  *int   `json:"age_days"`  // total days since registration; null when unknown
	Raw      string `json:"raw"`       // Raw WHOIS or RDAP JSON for debugging
	Source   string `json:"source"`    // "whois" or "rdap"
}

// Age returns the days since registration, and false when it isn't known.
// It checks the date rather than AgeKnown, so WHOIS results cached before the
// flag existed keep their age, and a cached missing date doesn't read as ancient.
func (r *RegistrationData) Age() (int, bool) {
	if r == nil || r.AgeDays == nil || r.CreatedDate.IsZero() {
		return 0, false
	}
	return *r.AgeDays, true
}
