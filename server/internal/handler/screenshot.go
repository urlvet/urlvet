package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/checks"
	"github.com/urlvet/urlvet/internal/service/screenshot"
)

func ScreenshotHandler(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query parameter is required"})
		return
	}

	// Basic URL validation (additional validation happens in TakeScreenshot)
	_, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	imageBytes, err := screenshot.TakeScreenshot(rawURL)
	if err != nil {
		// Log error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to capture screenshot"})
		return
	}

	// Return the image as binary with proper content-type
	c.Data(http.StatusOK, "image/png", imageBytes)
}
