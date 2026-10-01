package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
)

// CheckUrlLengthHandler checks whether the URL exceeds the safe length threshold.
//
//	@Summary		URL length check
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Router			/length [get]
func CheckUrlLengthHandler(c *gin.Context) {
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

	tooLong := checks.TooLongUrl(rawURL)

	c.JSON(http.StatusOK, gin.H{
		"too_long": tooLong,
	})
}

// CheckUrlDepthHandler checks whether the URL path is suspiciously deep.
//
//	@Summary		URL depth check
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Router			/depth [get]
func CheckUrlDepthHandler(c *gin.Context) {
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

	tooDeep := checks.TooDeepUrl(rawURL)

	c.JSON(http.StatusOK, gin.H{
		"too_deep": tooDeep,
	})
}
