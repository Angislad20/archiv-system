package utils

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"log"
)

// HasPermission checks if a role has the required permission
func HasPermission(roleName, permissionName string) bool {
	var role models.Role
	if err := database.DB.Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
		log.Printf("Role not found: %s", roleName)
		return false
	}

	for _, perm := range role.Permissions {
		if perm.Name == permissionName {
			log.Printf("Permission granted: %s for role %s", permissionName, roleName)
			return true
		}
	}
	log.Printf("Permission denied: %s for role %s", permissionName, roleName)
	return false
}
