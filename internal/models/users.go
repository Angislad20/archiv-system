package models

import "time"

type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"default:'user'"`
	CreatedAt time.Time
}

type Document struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	URL       string `gorm:"not null"`
	Tags      string `gorm:"not null"`
	CreatedAt time.Time
}

type Role struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null"`
}
