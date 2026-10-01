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
