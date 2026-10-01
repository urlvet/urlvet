package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
)

// CheckTrustedTLDHandler checks whether the URL uses a trusted (gov/edu) TLD.
//
//	@Summary		Trusted TLD check
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]any
//	@Failure		400	{object}	map[string]string
//	@Router			/trusted-tld [get]
func CheckTrustedTLDHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain"})
		return
	}

	trusted, icann, tld := checks.IsTrustedTld(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_trusted_tld": trusted,
		"is_icann":       icann,
		"tld":            tld,
	})
}

// CheckRiskyTLDHandler checks whether the URL uses a high-risk TLD.
//
//	@Summary		Risky TLD check
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]any
//	@Failure		400	{object}	map[string]string
//	@Router			/risky-tld [get]
func CheckRiskyTLDHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain"})
		return
	}

	risky, icann, tld := checks.IsRiskyTld(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_risky_tld": risky,
		"is_icann":     icann,
		"tld":          tld,
	})
}

// CheckUrlShortenerHandler checks whether the URL belongs to a known URL-shortener service.
//
//	@Summary		URL shortener detection
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Router			/url-shortener [get]
func CheckUrlShortenerHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	domain, err := checks.GetDomain(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain"})
		return
	}

	isShortener := checks.IsUrlShortener(domain)
	c.JSON(http.StatusOK, gin.H{
		"is_url_shortener": isShortener,
	})
}
