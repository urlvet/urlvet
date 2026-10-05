package threatfeeds

import "testing"

func TestNoPhishTankCallsWhenOff(t *testing.T) {
	t.Setenv("PHISHTANK_API_KEY", "key")
	t.Setenv("SAFE_BROWSING_API_KEY", "key")
	t.Setenv("WEBRISK_API_KEY", "key")

	t.Setenv("FEED_PHISHTANK", "0")
	if PhishTankAPIAllowed() {
		t.Error("FEED_PHISHTANK=0 still allows the PhishTank API")
	}

	t.Setenv("FEED_PHISHTANK", "")
	t.Setenv("THREAT_LOOKUPS", "0")
	if PhishTankAPIAllowed() || SafeBrowsingEnabled() || WebRiskEnabled() {
		t.Error("THREAT_LOOKUPS=0 left a lookup on")
	}

	t.Setenv("THREAT_LOOKUPS", "")
	if !SafeBrowsingEnabled() || !WebRiskEnabled() {
		t.Error("keys set but lookups off")
	}
}
