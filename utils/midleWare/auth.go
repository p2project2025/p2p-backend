package midleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware checks for valid JWT in Authorization header or cookie

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// secretKey := config.Cfg.JWTSecret
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("Authorization header or token cookie missing")
			tokenCookie, err := c.Cookie("token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header or token cookie missing"})
				log.Println("Error retrieving token from cookie:", err)
				return
			}
			log.Println("Using token from cookie:", tokenCookie)
			tokenString = tokenCookie
		} else {
			// Always try to remove the Bearer prefix
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set user info in context
		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		// Local dev: domain is empty so it works on localhost
		domain := "" // Empty = works for both localhost & production if needed

		// Set cookies
		c.SetCookie("token", tokenString, 3600, "/", domain, true, false)

		userID := claims.UserID
		role := claims.Role
		email := claims.Email

		c.SetCookie("userID", userID, 3600, "/", domain, true, false)
		c.SetCookie("role", role, 3600, "/", domain, true, false)
		c.SetCookie("email", email, 3600, "/", domain, true, false)

		// Also set token in the response header
		c.Header("Authorization", "Bearer "+tokenString)

		c.Next()
	}
}
