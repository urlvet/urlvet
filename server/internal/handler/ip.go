package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/logger"
	"github.com/urlvet/urlvet/internal/service/checks"
)

// CheckIfUsesIPHandler checks whether the URL uses a raw IP address instead of a domain.
//
//	@Summary		IP address detection
//	@Tags			URL
//	@Produce		json
//	@Param			url	query		string	true	"URL to check"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Router			/ip/check [get]
func CheckIfUsesIPHandler(c *gin.Context) {
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

	isIP, err := checks.UsesIPInsteadOfDomain(url)
	if err != nil {
		logger.Error("IP check failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uses_ip": isIP,
	})
}

// ResolveIPHandler resolves the domain in the URL to its IP addresses.
//
//	@Summary		DNS IP resolution
//	@Tags			DNS
//	@Produce		json
//	@Param			url	query		string	true	"URL whose domain to resolve"
//	@Success		200	{object}	map[string][]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/ip/resolve [get]
func ResolveIPHandler(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	ips, err := checks.GetIPAddress(domain)
	if err != nil {
		logger.Error("IP resolution failed", "domain", domain, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ip_addresses": ips,
	})
}
