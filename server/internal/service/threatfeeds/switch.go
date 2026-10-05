package threatfeeds

import "os"

// LookupsOff reports whether every threat lookup is switched off
// (THREAT_LOOKUPS=0): no lists are downloaded or loaded, and PhishTank, Safe
// Browsing and Web Risk aren't asked. For testing the scanner's own checks
// without calling anyone, whatever keys are set.
func LookupsOff() bool { return os.Getenv("THREAT_LOOKUPS") == "0" }

// PhishTankAPIAllowed reports whether a scan may ask PhishTank's API. It's
// only a stand-in while PhishTank's list is enabled but not loaded yet;
// FEED_PHISHTANK=0 means no PhishTank at all, key or no key.
func PhishTankAPIAllowed() bool {
	return !LookupsOff() && os.Getenv("FEED_PHISHTANK") != "0" && !FeedLoaded("PhishTank")
}
