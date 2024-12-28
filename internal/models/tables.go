package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	RoleID    uint   `gorm:"not null"`          // Référence au rôle
	Role      Role   `gorm:"foreignKey:RoleID"` // Relation avec Role
	CreatedAt time.Time
}

type Document struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	URL       string `gorm:"not null"`
	Tags      string `gorm:"not null"`
	OwnerID   uint   `gorm:"not null"`           // Référence à l'utilisateur propriétaire
	Owner     User   `gorm:"foreignKey:OwnerID"` // Relation avec User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Role struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
}

type Permission struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null;unique"` // Nom unique pour chaque permission
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RolePermission struct {
	ID           uint `gorm:"primary_key"`
	RoleID       uint `gorm:"not null"` // Clé étrangère vers Role
	PermissionID uint `gorm:"not null"` // Clé étrangère vers Permission
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
