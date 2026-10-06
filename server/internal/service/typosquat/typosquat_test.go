package typosquat

import "testing"

func TestCheckTyposquatting_ShortNamesNeedOneEdit(t *testing.T) {
	saved := topEntries
	defer func() { topEntries = saved }()
	topEntries = []topEntry{{domain: "htsc.com", sld: "htsc"}, {domain: "paypal.com", sld: "paypal"}}

	tests := []struct {
		domain string
		want   bool
	}{
		{"hdfc.bank.in", false}, // "hdfc" is two edits from "htsc": too loose for 4 letters
		{"htsx.com", true},      // one edit on a short name
		{"paypa1.com", true},    // one edit
		{"paypa11.com", true},   // two edits on a long name
	}
	for _, tt := range tests {
		if got := CheckTyposquatting(tt.domain).IsSuspicious; got != tt.want {
			t.Errorf("CheckTyposquatting(%q) = %v, want %v", tt.domain, got, tt.want)
		}
	}
}

func TestCheckTyposquatting_ComboSquatUsesCuratedBrands(t *testing.T) {
	saved := topEntries
	defer func() { topEntries = saved }()
	// service.gov.uk is a top site, but "service" is a plain word.
	topEntries = []topEntry{{domain: "service.gov.uk", sld: "service"}}

	tests := []struct {
		domain string
		want   bool
	}{
		{"takeoffservicesllc.com", false}, // contains "service"
		{"officesupplies.com", false},     // "office" is a brand name and a word
		{"paypal-secure-login.com", true}, // curated brand
		{"paypal.xyz", true},              // brand on someone else's TLD
		{"zohomail.com", false},           // the brand's own domain
	}
	for _, tt := range tests {
		got := CheckTyposquatting(tt.domain)
		if got.IsSuspicious != tt.want {
			t.Errorf("CheckTyposquatting(%q) = %+v, want suspicious=%v", tt.domain, got, tt.want)
		}
	}
}

func TestCheckURLParts(t *testing.T) {
	saved := lookupRank
	defer func() { lookupRank = saved }()
	lookupRank = func(d string) int {
		return map[string]int{"google.com": 1, "co.uk": 0}[d]
	}

	tests := []struct {
		url, domain                   string
		subBrand, embedded, pathBrand string
		exact                         bool
	}{
		{url: "https://horyzonix.paypal-login.antimoney-laundering.org", domain: "antimoney-laundering.org", subBrand: "paypal"},
		{url: "https://accounts.google.com.asso-entrautres.fr", domain: "asso-entrautres.fr", embedded: "google.com", subBrand: "google", exact: true}, // scored once, as embedded
		{url: "https://github.acme-corp.example/login", domain: "acme-corp.example", subBrand: "github", exact: true},
		{url: "https://1286524.us23.myftpupload.com/Paypal-dede2k25/x/", domain: "1286524.us23.myftpupload.com", pathBrand: "paypal"},
		{url: "https://nus.nbt.mybluehost.me/website_9f/832/Newcorreos/cc/login.php", domain: "nus.nbt.mybluehost.me", pathBrand: "correos"},
		{url: "https://www.paypal.com/myaccount/paypal-summary", domain: "paypal.com"}, // the brand's own site
		{url: "https://www.example.com/blog/how-to", domain: "example.com"},
		{url: "https://www.disneyland.com/", domain: "disneyland.com"}, // www. is not a subdomain brand
		{url: "https://onedrive.at-us.therelayservice.com/matpwp", domain: "therelayservice.com", subBrand: "onedrive", exact: true},
	}
	for _, tt := range tests {
		var got TyposquatResult
		CheckURLParts(tt.url, tt.domain, &got)
		if got.SubdomainBrand != tt.subBrand || got.EmbeddedDomain != tt.embedded || got.PathBrand != tt.pathBrand || got.SubdomainBrandExact != tt.exact {
			t.Errorf("%s: got sub=%q exact=%v embedded=%q path=%q", tt.url, got.SubdomainBrand, got.SubdomainBrandExact, got.EmbeddedDomain, got.PathBrand)
		}
	}
}
