package utils

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func GetDocumentsByOwnerID(db *gorm.DB, ownerID uint) ([]models.Document, error) {
	var documents []models.Document
	err := db.Where("owner_id = ?", ownerID).Find(&documents).Error
	return documents, err
}

func GetDocumentUpdatedAt(docID string) (time.Time, error) {
	var document models.Document

	// Filtrer par ID et récupérer uniquement updated at
	err := database.DB.Select("updated_at").Where("id = ?", docID).First(&document).Error

	if err != nil {
		return time.Time{}, err
	}

	return document.UpdatedAt, nil
}

func IsDocumentOwner(docID int, userID uint) (bool, error) {
	// Variable pour stocker le `owner_id` récupéré
	var ownerID uint

	// Exécution de la requête avec GORM
	err := database.DB.Table("documents").
		Select("owner_id").
		Where("id = ?", docID).
		Scan(&ownerID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Retourne false si aucun document n’est trouvé
			return false, nil
		}
		// Retourner l’erreur sans arrêter brusquement
		return false, err
	}

	// Vérifie si l’utilisateur est le propriétaire
	return ownerID == userID, nil
}
