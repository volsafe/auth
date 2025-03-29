package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks for valid JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		tokenString := ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString, "access")
		if err != nil {
			if err == ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "access token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			}
			c.Abort()
			return
		}

		// Store user information in the context
		if claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

func ValidateAccessToken(tokenString string) (any, any) {
	panic("unimplemented")
}
