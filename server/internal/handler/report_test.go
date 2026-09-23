package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestReportResultHandler(t *testing.T) {
	file := filepath.Join(t.TempDir(), "reports.jsonl")
	t.Setenv("REPORTS_FILE", file)
	r := testR()

	post := func(body string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/report", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	if code := post(`{"url":"https://example.com","verdict":"Risky","score":20,"expected_verdict":"Safe","comment":"legit site"}`); code != http.StatusCreated {
		t.Fatalf("valid report: got %d, want 201", code)
	}
	if code := post(`{"verdict":"Risky"}`); code != http.StatusBadRequest {
		t.Errorf("missing url: got %d, want 400", code)
	}
	if code := post(`{"url":"https://example.com","expected_verdict":"Maybe"}`); code != http.StatusBadRequest {
		t.Errorf("bad expected_verdict: got %d, want 400", code)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("reading reports file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 stored report, got %d", len(lines))
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("stored line is not JSON: %v", err)
	}
	if rec["url"] != "https://example.com" || rec["expected_verdict"] != "Safe" {
		t.Errorf("unexpected stored record: %v", rec)
	}
	ts, _ := rec["time"].(string)
	if _, err := time.ParseInLocation(reportTimeLayout, ts, ist); err != nil {
		t.Errorf("time %q not in %q format: %v", ts, reportTimeLayout, err)
	}

	// Malformed lines are skipped; newest report comes first.
	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString("not json\n")
	f.Close()
	post(`{"url":"https://second.example"}`)

	reports, err := readReports()
	if err != nil {
		t.Fatalf("readReports: %v", err)
	}
	if len(reports) != 2 || reports[0].URL != "https://second.example" {
		t.Errorf("readReports = %+v, want 2 reports newest first", reports)
	}
}

func TestReadReports_MissingFile(t *testing.T) {
	t.Setenv("REPORTS_FILE", filepath.Join(t.TempDir(), "none.jsonl"))
	reports, err := readReports()
	if err != nil || len(reports) != 0 {
		t.Errorf("readReports on missing file = %v, %v; want empty, nil", reports, err)
	}
}

func TestReadReports_SortedByTime(t *testing.T) {
	file := filepath.Join(t.TempDir(), "reports.jsonl")
	t.Setenv("REPORTS_FILE", file)
	os.WriteFile(file, []byte(
		`{"url":"https://b.example","time":"2026-10-02 10:00:00"}`+"\n"+
			`{"url":"https://a.example","time":"2026-09-30 10:00:00"}`+"\n"+
			`{"url":"https://c.example","time":"2026-10-03 10:00:00"}`+"\n"), 0o644)

	reports, err := readReports()
	if err != nil {
		t.Fatalf("readReports: %v", err)
	}
	var got []string
	for _, r := range reports {
		got = append(got, r.URL)
		if r.ID == "" {
			t.Errorf("report %s has no ID", r.URL)
		}
	}
	if strings.Join(got, ",") != "https://c.example,https://b.example,https://a.example" {
		t.Errorf("order = %v, want newest first", got)
	}
}

func TestAdminDeleteReportHandler(t *testing.T) {
	file := filepath.Join(t.TempDir(), "reports.jsonl")
	t.Setenv("REPORTS_FILE", file)
	os.WriteFile(file, []byte(
		`{"url":"https://keep.example","time":"2026-10-01 10:00:00"}`+"\n"+
			`{"url":"https://drop.example","time":"2026-10-02 10:00:00"}`+"\n"), 0o644)

	reports, _ := readReports()
	dropID := reports[0].ID // newest: drop.example

	del := func(id string) int {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id}}
		AdminDeleteReportHandler(c)
		return w.Code
	}

	if code := del(dropID); code != http.StatusOK {
		t.Fatalf("delete: got %d, want 200", code)
	}
	if code := del(dropID); code != http.StatusNotFound {
		t.Errorf("delete again: got %d, want 404", code)
	}

	reports, _ = readReports()
	if len(reports) != 1 || reports[0].URL != "https://keep.example" {
		t.Errorf("after delete = %+v, want only keep.example", reports)
	}
}
