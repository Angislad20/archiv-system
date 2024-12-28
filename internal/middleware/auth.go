package middleware

import (
	"archiv-system/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware checks the JWT and validates permissions dynamically from the database
func AuthMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Extract the token from the Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// get user role from claims
		role := claims.RoleName

		// Dynamically load user role permissions
		permissions, err := utils.LoadPermissions(role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load permissions"})
			c.Abort()
			return
		}

		// Check if the role has the required permission
		if !hasPermission(permissions, requiredPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission '%s' required", requiredPermission)})
			c.Abort()
			return
		}

		// Add user ID to context
		c.Set("UserID", claims.UserID)
		c.Set("Role", role)
		c.Next()
	}
}

// Helper pour vérifier si une permission est dans la liste
func hasPermission(permissions []string, requiredPermission string) bool {
	for _, perm := range permissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}
