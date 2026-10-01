package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/admintoken"
)

// BearerAuth validates the Authorization: Bearer <token> header.
// Tokens are issued by AdminLoginHandler and signed with ADMIN_JWT_SECRET.
func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || !admintoken.Validate(parts[1]) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
