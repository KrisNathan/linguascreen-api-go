package middleware

import (
	"net/http"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/application/port"
	"github.com/gin-gonic/gin"
)

const userIDKey = "user_id"

func RequiresAuth(tokenService port.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(strings.TrimSpace(c.GetHeader("Authorization")), "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := tokenService.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(userIDKey, userID)
		c.Next()
	}
}

// returns the authenticated user's ID. Only valid behind RequiresAuth.
func UserID(c *gin.Context) string {
	return c.GetString(userIDKey)
}
