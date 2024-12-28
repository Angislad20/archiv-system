package utils

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"fmt"
)

func LoadPermissions(role string) ([]string, error) {
	var permissions []models.Permission

	err := database.DB.
		Raw(`SELECT p.name 
                 FROM permissions AS p
                 JOIN role_permissions AS rp ON rp.permission_id = p.id
                 JOIN roles AS r ON r.id = rp.role_id
                 WHERE r.name = ?`, role).
		Scan(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("error loading permissions: %w", err)
	}

	var perms []string
	for _, p := range permissions {
		perms = append(perms, p.Name)
	}

	return perms, nil

}
