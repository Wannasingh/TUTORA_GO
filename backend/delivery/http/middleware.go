package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Save claims to context
		userID, ok := claims["sub"].(float64) // jwt numeric claims are float64 by default
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		c.Set("userID", int(userID))
		c.Set("email", claims["email"].(string))

		var roles []string
		if claimsRoles, ok := claims["roles"].([]interface{}); ok {
			for _, r := range claimsRoles {
				if rStr, ok := r.(string); ok {
					roles = append(roles, rStr)
				}
			}
		}
		c.Set("roles", roles)

		if len(roles) > 0 {
			c.Set("role", roles[0])
		} else {
			c.Set("role", "")
		}

		c.Next()
	}
}
