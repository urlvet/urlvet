package threatfeeds

import (
	"reflect"
	"sort"
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
)

// Examples from Safe Browsing's URL canonicalization spec.
func TestSBCanonicalize(t *testing.T) {
	tests := map[string]string{
		"http://host/%25%32%35":          "host/%25",
		"http://host/%25%32%35%25%32%35": "host/%25%25",
		"http://host/%2525252525252525":  "host/%25",
		"http://host/asdf%25%32%35asd":   "host/asdf%25asd",
		"http://host/%%%25%32%35asd%%":   "host/%25%25%25asd%25%25",
		"http://www.google.com/":         "www.google.com/",
		"http://%31%36%38%2e%31%38%38%2e%39%39%2e%32%36/%2E%73%65%63%75%72%65/%77%77%77%2E%65%62%61%79%2E%63%6F%6D/": "168.188.99.26/.secure/www.ebay.com/",
		"http://www.google.com/blah/..":             "www.google.com/",
		"www.google.com/":                           "www.google.com/",
		"www.google.com":                            "www.google.com/",
		"http://www.evil.com/blah#frag":             "www.evil.com/blah",
		"http://www.GOOgle.com/":                    "www.google.com/",
		"http://www.google.com.../":                 "www.google.com/",
		"http://www.google.com/foo\tbar\rbaz\n2":    "www.google.com/foobarbaz2",
		"http://www.google.com/q?":                  "www.google.com/q?",
		"http://www.google.com/q?r?":                "www.google.com/q?r?",
		"http://www.google.com/q?r?s":               "www.google.com/q?r?s",
		"http://evil.com/foo#bar#baz":               "evil.com/foo",
		"http://evil.com/foo;":                      "evil.com/foo;",
		"http://evil.com/foo?bar;":                  "evil.com/foo?bar;",
		"http://notrailingslash.com":                "notrailingslash.com/",
		"http://www.gotaport.com:1234/":             "www.gotaport.com/",
		"  http://www.google.com/  ":                "www.google.com/",
		"http://%20leadingspace.com/":               "%20leadingspace.com/",
		"https://www.securesite.com/":               "www.securesite.com/",
		"http://host.com/ab%23cd":                   "host.com/ab%23cd",
		"http://host.com//twoslashes?more//slashes": "host.com/twoslashes?more//slashes",
	}
	for in, want := range tests {
		host, path, query, err := sbCanonicalize(in)
		if got := host + path + query; err != nil || got != want {
			t.Errorf("sbCanonicalize(%q) = %q (%v), want %q", in, got, err, want)
		}
	}
}

func TestSBExpressions(t *testing.T) {
	tests := map[string][]string{
		"http://a.b.c/1/2.html?param=1": {
			"a.b.c/1/2.html?param=1", "a.b.c/1/2.html", "a.b.c/", "a.b.c/1/",
			"b.c/1/2.html?param=1", "b.c/1/2.html", "b.c/", "b.c/1/",
		},
		"http://a.b.c.d.e.f.g/1.html": {
			"a.b.c.d.e.f.g/1.html", "a.b.c.d.e.f.g/",
			"c.d.e.f.g/1.html", "c.d.e.f.g/", "d.e.f.g/1.html", "d.e.f.g/",
			"e.f.g/1.html", "e.f.g/", "f.g/1.html", "f.g/",
		},
		"http://1.2.3.4/1/": {"1.2.3.4/1/", "1.2.3.4/"},
	}
	for in, want := range tests {
		got, err := sbExpressions(in)
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(got)
		sort.Strings(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("sbExpressions(%q)\n got  %v\n want %v", in, got, want)
		}
	}
}

func TestDecodeSearchHashesResponse(t *testing.T) {
	hash := make([]byte, 32)
	hash[0] = 0xAB
	detail := func(threat uint64) []byte {
		var d []byte
		d = protowire.AppendTag(d, 1, protowire.VarintType)
		d = protowire.AppendVarint(d, threat)
		d = protowire.AppendTag(d, 2, protowire.VarintType) // an attribute, ignored
		d = protowire.AppendVarint(d, 1)
		return d
	}
	var fh []byte
	fh = protowire.AppendTag(fh, 1, protowire.BytesType)
	fh = protowire.AppendBytes(fh, hash)
	fh = protowire.AppendTag(fh, 2, protowire.BytesType)
	fh = protowire.AppendBytes(fh, detail(2))
	fh = protowire.AppendTag(fh, 2, protowire.BytesType)
	fh = protowire.AppendBytes(fh, detail(0)) // unspecified: disregarded
	var msg []byte
	msg = protowire.AppendTag(msg, 1, protowire.BytesType)
	msg = protowire.AppendBytes(msg, fh)
	msg = protowire.AppendTag(msg, 2, protowire.BytesType) // cache_duration, skipped
	msg = protowire.AppendBytes(msg, []byte{0x08, 0x2c})

	got, err := decodeSearchHashesResponse(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].hash[0] != 0xAB || !reflect.DeepEqual(got[0].threatTypes, []string{"SOCIAL_ENGINEERING"}) {
		t.Errorf("decoded %+v", got)
	}
	if got, err := decodeSearchHashesResponse(nil); err != nil || len(got) != 0 {
		t.Errorf("empty response = %v, %v", got, err)
	}
}
