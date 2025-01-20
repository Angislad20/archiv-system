package models

import (
	"time"
)

type Document struct {
	ID                uint      `gorm:"primary_key"`
	Name              string    `gorm:"not null"`
	Type              string    `gorm:"not null"`
	URL               string    `gorm:"not null"`
	Tags              *[]Tag    `gorm:"many2many:document_tags;"`
	OwnerID           uint      `gorm:"not null"`           // Référence à l'utilisateur propriétaire
	Owner             User      `gorm:"foreignKey:OwnerID"` // Relation avec User
	Version           int       `gorm:"default:1"`
	PreviousVersionID uint      `gorm:"default:0"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

type DocumentTag struct {
	DocumentID uint `gorm:"primaryKey"`
	TagID      uint `gorm:"primaryKey"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
