package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	RoleID    uint      `gorm:"not null"`
	Role      Role      `gorm:"foreignKey:RoleID"` // Associe Role avec User
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Role struct {
	ID          uint          `gorm:"primaryKey"`
	Name        string        `gorm:"unique;not null"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}
