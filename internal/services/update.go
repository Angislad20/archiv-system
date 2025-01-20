package services

import (
	"archiv-system/internal/database"
	"archiv-system/internal/models"
	"errors"
	"fmt"
)

type DocumentService struct{}

func (ds *DocumentService) ProcessFileUpdate(docID string, updateRequest models.UpdateRequest) (*models.Document, error) {
	var document models.Document

	// Charger le document
	if err := database.DB.First(&document, docID).Error; err != nil {
		return nil, errors.New("document not found")
	}

	// Appliquer les modifications
	document.Name = updateRequest.Name
	document.Type = updateRequest.Type

	// Convertir les noms de tags en modèles de tags
	tags, err := findOrCreateTags(updateRequest.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to process tags: %w", err)
	}

	document.Tags = &tags

	// Sauvegarder dans la base
	if err := database.DB.Save(&document).Error; err != nil {
		return nil, errors.New("failed to save document")
	}

	return &document, nil
}
