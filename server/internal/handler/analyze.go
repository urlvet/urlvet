package handler

import (
	"context"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/analyzer"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/service/checks"
)

// AnalyzeURLHandler runs a full safety analysis on the given URL.
//
//	@Summary		Full URL analysis
//	@Description	Runs every check in parallel and returns one scored report. Most clients only need
//	@Description	result.verdict (Safe, Suspicious or Risky), result.final_score (0-100, higher is safer)
//	@Description	and result.reasons. Checks still running after 15 seconds are dropped and named in
//	@Description	incomplete_checks; incomplete is true when that could change the verdict. Complete
//	@Description	results are cached for 24 hours per URL; incomplete ones are not cached.
//	@Tags			Analysis
//	@Produce		json
//	@Param			url	query		string	true	"URL to analyse (max 2048 chars)"
//	@Success		200	{object}	analyzer.Response
//	@Failure		400	{object}	map[string]string
//	@Failure		422	{object}	map[string]string
//	@Router			/analyze [get]
func AnalyzeURLHandler(c *gin.Context) {
	url := strings.TrimSpace(c.Query("url"))
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(url)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": "invalid url"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, errs := analyzer.Analyze(ctx, url)
	if len(errs) > 0 {
		origin := url
		if u, err := neturl.Parse(url); err == nil {
			origin = u.Scheme + "://" + u.Host
		}
		for _, e := range errs {
			logger.Warn("analyzer error", "origin", origin, "err", e)
		}
	}
	// No verdict means the analysis never ran (e.g. the URL has no usable
	// domain). Say so, rather than send an empty result as a success.
	if resp.Result.Verdict == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "ERROR", "error": "could not analyze this URL"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
