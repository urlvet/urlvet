// This handler is for development/testing purposes only.
// This code is intended for local experiments and feature prototyping.
// Do NOT commit permanent changes or production features to this handler.

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
)

func TestApiHandler(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not extract domain from url"})
		return
	}

	////////////////////////////////////////////////////////////////////
	// THIS IS A TEMPORARY CODE SECTION
	// MAKE CHANGES ONLY INSIDE THIS BOX
	// SAFE TO DELETE LATER IF PUSHED BY MISTAKE

	c.JSON(http.StatusOK, gin.H{
		"domain": domain,
	})

	////////////////////////////////////////////////////////////////////

}
