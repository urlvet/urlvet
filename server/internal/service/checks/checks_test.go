package checks

import "testing"

func TestTooLongUrl(t *testing.T) {
	url := "https://google.com"
	result := TooLongUrl(url)
	if result == true {
		t.Errorf("Error while testing TooLongUrl(%v)", url)
	}
}

func TestIsValidURL_LowercasesHost(t *testing.T) {
	tests := map[string]string{
		"http://Google.com/aa":      "http://google.com/aa",
		"Google.com/Hh":             "https://google.com/Hh",
		"HTTPS://WWW.Example.COM":   "https://www.example.com",
		"https://Example.com/A?Q=B": "https://example.com/A?Q=B",
	}
	for in, want := range tests {
		u, ok, err := IsValidURL(in)
		if err != nil || !ok {
			t.Fatalf("IsValidURL(%q) = invalid, err=%v", in, err)
		}
		if got := u.String(); got != want {
			t.Errorf("IsValidURL(%q) = %q, want %q", in, got, want)
		}
	}
}
