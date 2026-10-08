package threatfeeds

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestFeedKey(t *testing.T) {
	tests := map[string]string{
		"https://www.Evil.example/Login/": "evil.example/Login",
		"http://evil.example/login":       "evil.example/login",
		"evil.example/a?b=c":              "evil.example/a?b=c",
		"https://evil.example/x#frag":     "evil.example/x",
		"":                                "",
	}
	for in, want := range tests {
		if got, _ := feedKey(in); got != want {
			t.Errorf("feedKey(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLookupLocal(t *testing.T) {
	saved := local
	defer func() { local = saved }()
	local = &store{urls: map[string]map[string]struct{}{}, hosts: map[string]map[string]struct{}{}}
	local.set("PhishTank", []string{"https://evil.example/login", "http://bad.pages.dev/"})

	if m := LookupLocal("https://www.evil.example/login/", false); !m.Listed || m.Match != "url" || m.Sources[0] != "PhishTank" {
		t.Errorf("exact match = %+v", m)
	}
	if m := LookupLocal("https://evil.example/other", false); !m.Listed || m.Match != "host" {
		t.Errorf("same-site match = %+v", m)
	}
	if m := LookupLocal("https://evil.example/other", true); m.Listed {
		t.Errorf("host match on a shared host = %+v", m)
	}
	if m := LookupLocal("https://fine.example/", false); m.Listed {
		t.Errorf("unlisted site matched: %+v", m)
	}
	if !LocalFeedsLoaded() || !FeedLoaded("PhishTank") || FeedLoaded("OpenPhish") {
		t.Error("loaded-feed reporting is wrong")
	}
}

func TestParsePhishTankDump(t *testing.T) {
	if _, err := exec.LookPath("bzip2"); err != nil {
		t.Skip("bzip2 not installed")
	}
	cmd := exec.Command("bzip2", "-c")
	cmd.Stdin = strings.NewReader(`[{"phish_id":1,"url":"https://evil.example/a"},{"phish_id":2,"url":"http://bad.example/"}]`)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	list, err := parsePhishTankDump(&out)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0] != "https://evil.example/a" {
		t.Errorf("parsed %v", list)
	}
}

func TestParseURLLines(t *testing.T) {
	list, _ := parseURLLines(strings.NewReader("# URLhaus header\n\nhttps://a.example/x\n  http://b.example/  \n"))
	if len(list) != 2 || list[1] != "http://b.example/" {
		t.Errorf("parsed %v", list)
	}
}

func TestParseURLhausCSV(t *testing.T) {
	csv := `################################################################
# abuse.ch URLhaus Database Dump (CSV - recent URLs)              #
# id,dateadded,url,url_status,last_online,threat,tags,urlhaus_link,reporter
"3601234","2026-10-07 10:00:00","http://1.2.3.4:8080/bins/x.sh","online","2026-10-07 10:00:00","malware_download","elf,mirai","https://urlhaus.abuse.ch/url/3601234/","r"
"3601235","2026-10-07 09:00:00","http://gone.example/a.exe","offline","","malware_download","exe","https://urlhaus.abuse.ch/url/3601235/","r"
`
	list, err := parseURLhausCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0] != "http://1.2.3.4:8080/bins/x.sh" {
		t.Errorf("parsed %v, want only the online URL", list)
	}
}
