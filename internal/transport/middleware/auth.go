package middleware

import (
	"net/http"
	"strings"

	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/gin-gonic/gin"
)

func AuthRequired(tokens application.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		raw := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			raw = strings.TrimPrefix(authHeader, "Bearer ")
		}

		userID, err := tokens.Validate(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
