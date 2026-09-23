package handler

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
)

const (
	defaultReportsFile = "data/reports.jsonl"
	reportTimeLayout   = "2006-01-02 15:04:05"
)

// ist is India Standard Time. A fixed zone (India has no DST) avoids depending
// on tzdata being present in the container.
var ist = time.FixedZone("IST", 5*60*60+30*60)

var reportsMu sync.Mutex

// ReportRequest is a user report that a scan result looks wrong.
type ReportRequest struct {
	URL             string `json:"url" binding:"required,max=2048"`
	Verdict         string `json:"verdict" binding:"max=32"`
	Score           int    `json:"score"`
	ExpectedVerdict string `json:"expected_verdict" binding:"omitempty,oneof=Safe Suspicious Risky"`
	Comment         string `json:"comment" binding:"max=1000"`
}

// ReportRecord is a stored report. Time is in IST, formatted as "2006-01-02 15:04:05".
// ID is not stored: it's derived from the stored line when reading (see reportID).
type ReportRecord struct {
	ReportRequest
	Time string `json:"time"`
	ID   string `json:"id,omitempty"`
}

// reportID identifies a report by a hash of its stored line, so reports written
// before IDs existed can still be addressed and deleted.
func reportID(line []byte) string {
	sum := sha256.Sum256(bytes.TrimSpace(line))
	return hex.EncodeToString(sum[:8])
}

// reportsFilePath reads REPORTS_FILE from the environment, falling back to data/reports.jsonl.
func reportsFilePath() string {
	if p := strings.TrimSpace(os.Getenv("REPORTS_FILE")); p != "" {
		return p
	}
	return defaultReportsFile
}

// ReportResultHandler stores a user report about an incorrect scan result.
//
//	@Summary		Report an incorrect result
//	@Tags			Analysis
//	@Accept			json
//	@Produce		json
//	@Param			report	body		ReportRequest	true	"Report details"
//	@Success		201		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Router			/report [post]
func ReportResultHandler(c *gin.Context) {
	var req ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report"})
		return
	}
	req.URL = strings.TrimSpace(req.URL)
	req.Comment = strings.TrimSpace(req.Comment)

	line, err := json.Marshal(ReportRecord{ReportRequest: req, Time: time.Now().In(ist).Format(reportTimeLayout)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save report"})
		return
	}

	if err := appendReport(line); err != nil {
		logger.Error("failed to save report", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

// AdminDeleteReportHandler deletes one report, once the issue it points to is
// resolved. Reports are only kept until then (see the privacy page).
func AdminDeleteReportHandler(c *gin.Context) {
	found, err := deleteReport(c.Param("id"))
	if err != nil {
		logger.Error("failed to delete report", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete report"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// AdminReportsHandler returns all stored reports, newest first.
func AdminReportsHandler(c *gin.Context) {
	reports, err := readReports()
	if err != nil {
		logger.Error("failed to read reports", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read reports"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

// readReports parses the reports file, skipping malformed lines. A missing file means no reports yet.
func readReports() ([]ReportRecord, error) {
	reportsMu.Lock()
	data, err := os.ReadFile(reportsFilePath())
	reportsMu.Unlock()
	if errors.Is(err, os.ErrNotExist) {
		return []ReportRecord{}, nil
	}
	if err != nil {
		return nil, err
	}

	reports := []ReportRecord{}
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var r ReportRecord
		if err := json.Unmarshal(line, &r); err != nil {
			continue
		}
		r.ID = reportID(line)
		reports = append(reports, r)
	}
	// Newest first. The time layout sorts correctly as a string; the reverse
	// keeps same-second reports in newest-written-first order.
	slices.Reverse(reports)
	slices.SortStableFunc(reports, func(a, b ReportRecord) int { return cmp.Compare(b.Time, a.Time) })
	return reports, nil
}

// deleteReport rewrites the reports file without the report with this ID.
// It reports whether a matching report was found.
func deleteReport(id string) (bool, error) {
	reportsMu.Lock()
	defer reportsMu.Unlock()

	path := reportsFilePath()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	var kept bytes.Buffer
	found := false
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if !found && reportID(line) == id {
			found = true
			continue
		}
		kept.Write(line)
		kept.WriteByte('\n')
	}
	if !found {
		return false, nil
	}

	// Write to a temp file and rename, so a crash can't leave a half-written file.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, kept.Bytes(), 0o644); err != nil {
		return false, err
	}
	return true, os.Rename(tmp, path)
}

func appendReport(line []byte) error {
	reportsMu.Lock()
	defer reportsMu.Unlock()

	path := reportsFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}
