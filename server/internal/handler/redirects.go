package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
)

// CheckRedirectsHandler follows the URL's redirect chain and reports findings.
//
//	@Summary		Redirect chain analysis
//	@Tags			Security
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]checks.RedirectionResult
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/redirects [get]
func CheckRedirectsHandler(c *gin.Context) {
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

	redir, err := checks.CheckRedirects(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"redirection_result": redir,
	})
}
