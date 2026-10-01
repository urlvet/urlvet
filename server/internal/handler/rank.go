package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/rank"
)

type domainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// GetDomainRankHandler returns the global popularity rank of the URL's domain.
//
//	@Summary		Domain popularity rank
//	@Tags			DNS
//	@Produce		json
//	@Param			url	query		string	true	"URL whose domain to rank"
//	@Success		200	{object}	map[string]int
//	@Failure		400	{object}	map[string]string
//	@Router			/rank [get]
func GetDomainRankHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	_, isValid, err := checks.IsValidURL(url)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	domain, err := checks.GetDomain(url)
	if err != nil {
		logger.Error("domain extraction failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain from url"})
		return
	}

	rankValue := rank.DomainRankLookup(domain)

	c.JSON(http.StatusOK, gin.H{
		"rank": rankValue,
	})
}
