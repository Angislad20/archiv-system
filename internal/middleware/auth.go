package middleware

import (
	"archiv-system/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware verifies the validity of the JWT and adds user information to the context
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Split the token type (Bearer) and the token itself
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Verify and parse the token
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token", "error": err.Error()})
			c.Abort()
			return
		}

		// Check if the role is missing in the token
		if claims.RoleName == "" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Role is missing in token"})
			c.Abort()
			return
		}

		// Add token information to the context
		c.Set("userID", claims.UserID)
		c.Set("roleName", claims.RoleName)

		// Log user information
		log.Printf("User ID: %d, Role: %s", claims.UserID, claims.RoleName)

		// Continue the request
		c.Next()
	}
}

// AuthMiddleware verifies if the user has the required permission
func AuthMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve token information from the context (set by JWTAuthMiddleware)
		userID, exists := c.Get("userID")
		roleName, existsRole := c.Get("roleName")
		if !exists || !existsRole {
			utils.RespondError(c, http.StatusInternalServerError, "User information is missing in context", nil)
			c.Abort()
			return
		}

		// Check if the user has the required permission
		if !utils.HasPermission(roleName.(string), requiredPermission) {
			utils.RespondError(c, http.StatusForbidden, "You don't have permission to access this resource", nil)
			log.Printf("Permission denied: %s for user ID: %v, Role: %s", requiredPermission, userID, roleName)
			c.Abort()
			return
		}

		// Log permission check
		log.Printf("Permission granted: %s for user ID: %v, Role: %s", requiredPermission, userID, roleName)

		// Continue the request
		c.Next()
	}
}
