package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		googleClientID := os.Getenv("GOOGLE_CLIENT_ID")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Extract Bearer token from the Authorization header
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token missing"})
			return
		}

		// Validate the token using Google ID Token validator
		payload, err := idtoken.Validate(context.Background(), bearerToken, googleClientID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		// Add user claims to the context if token is valid
		c.Set("userPayload", payload.Claims)
		c.Next()
	}
}
