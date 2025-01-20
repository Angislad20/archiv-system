package database

import (
	"archiv-system/internal/models"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=root dbname=archiv_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	// Assign global database
	DB = db

	// Table migration
	if err := db.AutoMigrate(
		&models.Document{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
	); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// Seed roles and permissions
	if err := SeedRolesAndPermissions(); err != nil {
		panic("failed to seed roles and permissions: " + err.Error())
	}

	var admin models.User
	if err := DB.Where("role_id = ?", 1).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a default admin
			defaultAdmin := models.User{
				Username: "Angislad",
				Password: HashPassword("viceash"), // Use a secure hashing function
				RoleID:   1,                       // 1 corresponds to the Admin role
			}
			DB.Create(&defaultAdmin)
			log.Println("Default admin created: username=Angislad, password=viceash")
		}
	}

	// Log admin information
	log.Printf("Admin user: %+v\n", admin)

	return DB
}

func SeedRolesAndPermissions() error {
	// Define roles and permissions
	permissions := []string{"read_document", "update_document", "delete_document", "upload_document"}

	// Create permissions
	var createdPermissions []models.Permission
	for _, permName := range permissions {
		perm := models.Permission{Name: permName}
		if err := DB.FirstOrCreate(&perm, "name = ?", permName).Error; err != nil {
			return fmt.Errorf("failed to seed permission '%s': %v", permName, err)
		}
		createdPermissions = append(createdPermissions, perm)
	}

	// Create roles and assign permissions
	rolePermissions := map[string][]string{
		"admin": {"read_document", "update_document", "delete_document", "upload_document"},
		"user":  {"read_document", "upload_document"},
	}

	for roleName, permNames := range rolePermissions {
		var role models.Role
		// Use FirstOrCreate to ensure no duplicates are created
		if err := DB.Where("name = ?", roleName).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// If the role does not exist, create it
				role = models.Role{Name: roleName}
				if err := DB.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to seed role '%s': %v", roleName, err)
				}
			} else {
				return fmt.Errorf("failed to check if role '%s' exists: %v", roleName, err)
			}
		}

		// Assign permissions
		for _, permName := range permNames {
			for _, perm := range createdPermissions {
				if perm.Name == permName {
					if err := DB.Model(&role).Association("Permissions").Append(&perm); err != nil {
						return fmt.Errorf("failed to assign permission '%s' to role '%s': %v", permName, roleName, err)
					}
				}
			}
		}

		// Log role and permissions information
		log.Printf("Role: %+v\n", role)
		log.Printf("Permissions: %+v\n", role.Permissions)
	}

	return nil
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

// CheckPassword verifies if a password matches its hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
