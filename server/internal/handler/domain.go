package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/domaininfo"
)

// DomainInfoHandler returns WHOIS / RDAP registration data for the URL's domain.
//
//	@Summary		Domain WHOIS / RDAP info
//	@Tags			DNS
//	@Produce		json
//	@Param			url	query		string	true	"URL whose domain to look up"
//	@Success		200	{object}	map[string]any
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/domain-info [get]
func DomainInfoHandler(c *gin.Context) {
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

	domainInfo, err := domaininfo.Lookup(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "domainInfo lookup failed",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":     domain,
		"domainInfo": domainInfo,
	})
}
